// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhook

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/apache/dubbo-admin/pkg/authority/config"
	"github.com/apache/dubbo-admin/pkg/authority/logger"
	"github.com/mattbaird/jsonpatch"
	"io"
	admissionV1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type PodPatch func(*v1.Pod) (*v1.Pod, error)
type GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)

type Webhook struct {
	Patches        []PodPatch
	AllowOnErr     bool
	getCertificate GetCertificate
	server         *http.Server
	cert           string
	key            string
}

type AdmitError struct {
	error
	message string
}

func NewAdmitError(message string) *AdmitError {
	return &AdmitError{
		message: message,
	}
}

func (e *AdmitError) Error() string {
	return e.message
}

func NewWebhook(certificate GetCertificate) *Webhook {
	return &Webhook{
		getCertificate: certificate,
	}
}

func (wh *Webhook) NewServer(port int32) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", wh.ServeHealth)
	mux.HandleFunc("/mutating-services", wh.Mutate)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: wh.getCertificate,
		},
	}
}

func (wh *Webhook) Init(options *config.Options) {
	wh.server = wh.NewServer(options.WebhookPort)
}

func (wh *Webhook) Serve() {
	if wh.cert != "" && wh.key != "" {
		err := wh.server.ListenAndServeTLS("", "")
		if err != nil {
			logger.Sugar.Fatalf("[Webhook] Serve webhook server failed. %v", err.Error())
			return
		}
	} else {
		err := wh.server.ListenAndServe()
		if err != nil {
			logger.Sugar.Fatalf("[Webhook] Serve webhook server failed. %v", err.Error())
			return
		}
	}
}

func (wh *Webhook) Stop() {
	err := wh.server.Close()
	if err != nil {
		logger.Sugar.Fatalf("[Webhook] Stop webhook server failed. %v", err.Error())
		return
	}
}

// ServeHealth returns 200 when things are good
func (wh *Webhook) ServeHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (wh *Webhook) Mutate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	logger.Sugar.Infof("[Webhook] Mutation request: " + string(body))

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		outputLog := fmt.Sprintf("[Webhook] contentType=%s, expect application/json", contentType)
		logger.Sugar.Errorf(outputLog)
		return
	}

	var reviewResponse *admissionV1.AdmissionResponse
	ar := admissionV1.AdmissionReview{}
	if err := json.Unmarshal(body, &ar); err != nil {
		outputLog := fmt.Sprintf("[Webhook] json unmarshal err=%s", err)
		logger.Sugar.Errorf(outputLog)

		reviewResponse = &admissionV1.AdmissionResponse{
			Allowed: wh.AllowOnErr,
			Result: &metav1.Status{
				Status:  "Failure",
				Message: err.Error(),
				Reason:  metav1.StatusReason(err.Error()),
			},
		}
	} else {
		reviewResponse, err = wh.Admit(ar)
		if err != nil {
			logger.Sugar.Errorf(err.Error())

			reviewResponse = &admissionV1.AdmissionResponse{
				Allowed: wh.AllowOnErr,
				Result: &metav1.Status{
					Status:  "Failure",
					Message: err.Error(),
					Reason:  metav1.StatusReason(err.Error()),
				},
			}
		}
	}

	response := admissionV1.AdmissionReview{}
	response.TypeMeta.Kind = "AdmissionReview"
	response.TypeMeta.APIVersion = "admission.k8s.io/v1"
	response.Response = reviewResponse

	logger.Sugar.Infof("[Webhook] AdmissionReview response: %v", response)

	resp, err := json.Marshal(response)
	if err != nil {
		outputLog := fmt.Sprintf("[Webhook] response json unmarshal err=%s", err)
		logger.Sugar.Errorf(outputLog)
	}
	if _, err := w.Write(resp); err != nil {
		outputLog := fmt.Sprintf("[Webhook] write resp err=%s", err)
		logger.Sugar.Errorf(outputLog)
	}
}

func (wh *Webhook) Admit(ar admissionV1.AdmissionReview) (*admissionV1.AdmissionResponse, error) {
	var reviewResponse = &admissionV1.AdmissionResponse{
		Allowed: true,
		UID:     ar.Request.UID,
	}

	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	if ar.Request.Resource != podResource {
		outputLog := fmt.Sprintf("[Webhook] expect resource to be %s", podResource)
		return nil, NewAdmitError(outputLog)
	}

	raw := ar.Request.Object.Raw

	pod := v1.Pod{}

	if err := json.Unmarshal(raw, &pod); err != nil {
		outputLog := fmt.Sprintf("[Webhook] pod unmarshal error. %s", err)
		return nil, NewAdmitError(outputLog)
	}

	patchBytes, err := wh.PatchPod(&pod)
	if err != nil {
		outputLog := fmt.Sprintf("[Webhook] Patch error: %v", pod.ObjectMeta.Name)
		return nil, NewAdmitError(outputLog)
	}

	reviewResponse.Patch = patchBytes

	logger.Sugar.Infof("[Webhook] Patch after mutate : %s", string(patchBytes))
	pt := admissionV1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return reviewResponse, nil
}

func (wh *Webhook) PatchPod(pod *v1.Pod) ([]byte, error) {
	origin, originErr := json.Marshal(pod)

	if originErr == nil {
		logger.Sugar.Infof("[Webhook] Pod before mutate: %v", string(origin))
	} else {
		return nil, originErr
	}

	for _, patch := range wh.Patches {
		pod, _ = patch(pod)
	}

	after, afterErr := json.Marshal(pod)

	if afterErr == nil {
		logger.Sugar.Infof("[Webhook] Pod after mutate: %v", string(after))
	} else {
		return nil, afterErr
	}

	patch, patchErr := jsonpatch.CreatePatch(origin, after)
	if patchErr != nil {
		return nil, patchErr
	}

	return json.Marshal(patch)
}

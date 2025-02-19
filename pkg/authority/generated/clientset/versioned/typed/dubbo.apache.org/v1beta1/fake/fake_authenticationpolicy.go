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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"
	v1beta12 "github.com/apache/dubbo-admin/pkg/authority/apis/dubbo.apache.org/v1beta1"
	dubboapacheorgv1beta1 "github.com/apache/dubbo-admin/pkg/authority/generated/applyconfiguration/dubbo.apache.org/v1beta1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAuthenticationPolicies implements AuthenticationPolicyInterface
type FakeAuthenticationPolicies struct {
	Fake *FakeDubboV1beta1
	ns   string
}

var authenticationpoliciesResource = v1beta12.SchemeGroupVersion.WithResource("authenticationpolicies")

var authenticationpoliciesKind = v1beta12.SchemeGroupVersion.WithKind("AuthenticationPolicy")

// Get takes name of the authenticationPolicy, and returns the corresponding authenticationPolicy object, and an error if there is any.
func (c *FakeAuthenticationPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta12.AuthenticationPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(authenticationpoliciesResource, c.ns, name), &v1beta12.AuthenticationPolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta12.AuthenticationPolicy), err
}

// List takes label and field selectors, and returns the list of AuthenticationPolicies that match those selectors.
func (c *FakeAuthenticationPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1beta12.AuthenticationPolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(authenticationpoliciesResource, authenticationpoliciesKind, c.ns, opts), &v1beta12.AuthenticationPolicyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta12.AuthenticationPolicyList{ListMeta: obj.(*v1beta12.AuthenticationPolicyList).ListMeta}
	for _, item := range obj.(*v1beta12.AuthenticationPolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested authenticationPolicies.
func (c *FakeAuthenticationPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(authenticationpoliciesResource, c.ns, opts))

}

// Create takes the representation of a authenticationPolicy and creates it.  Returns the server's representation of the authenticationPolicy, and an error, if there is any.
func (c *FakeAuthenticationPolicies) Create(ctx context.Context, authenticationPolicy *v1beta12.AuthenticationPolicy, opts v1.CreateOptions) (result *v1beta12.AuthenticationPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(authenticationpoliciesResource, c.ns, authenticationPolicy), &v1beta12.AuthenticationPolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta12.AuthenticationPolicy), err
}

// Update takes the representation of a authenticationPolicy and updates it. Returns the server's representation of the authenticationPolicy, and an error, if there is any.
func (c *FakeAuthenticationPolicies) Update(ctx context.Context, authenticationPolicy *v1beta12.AuthenticationPolicy, opts v1.UpdateOptions) (result *v1beta12.AuthenticationPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(authenticationpoliciesResource, c.ns, authenticationPolicy), &v1beta12.AuthenticationPolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta12.AuthenticationPolicy), err
}

// Delete takes name of the authenticationPolicy and deletes it. Returns an error if one occurs.
func (c *FakeAuthenticationPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(authenticationpoliciesResource, c.ns, name, opts), &v1beta12.AuthenticationPolicy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAuthenticationPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(authenticationpoliciesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta12.AuthenticationPolicyList{})
	return err
}

// Patch applies the patch and returns the patched authenticationPolicy.
func (c *FakeAuthenticationPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta12.AuthenticationPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(authenticationpoliciesResource, c.ns, name, pt, data, subresources...), &v1beta12.AuthenticationPolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta12.AuthenticationPolicy), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied authenticationPolicy.
func (c *FakeAuthenticationPolicies) Apply(ctx context.Context, authenticationPolicy *dubboapacheorgv1beta1.AuthenticationPolicyApplyConfiguration, opts v1.ApplyOptions) (result *v1beta12.AuthenticationPolicy, err error) {
	if authenticationPolicy == nil {
		return nil, fmt.Errorf("authenticationPolicy provided to Apply must not be nil")
	}
	data, err := json.Marshal(authenticationPolicy)
	if err != nil {
		return nil, err
	}
	name := authenticationPolicy.Name
	if name == nil {
		return nil, fmt.Errorf("authenticationPolicy.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(authenticationpoliciesResource, c.ns, *name, types.ApplyPatchType, data), &v1beta12.AuthenticationPolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta12.AuthenticationPolicy), err
}

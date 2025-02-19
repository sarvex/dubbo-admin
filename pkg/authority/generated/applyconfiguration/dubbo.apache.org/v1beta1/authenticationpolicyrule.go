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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

// AuthenticationPolicyRuleApplyConfiguration represents an declarative configuration of the AuthenticationPolicyRule type for use
// with apply.
type AuthenticationPolicyRuleApplyConfiguration struct {
	From *AuthenticationPolicySourceApplyConfiguration `json:"from,omitempty"`
	To   *AuthenticationPolicyTargetApplyConfiguration `json:"to,omitempty"`
}

// AuthenticationPolicyRuleApplyConfiguration constructs an declarative configuration of the AuthenticationPolicyRule type for use with
// apply.
func AuthenticationPolicyRule() *AuthenticationPolicyRuleApplyConfiguration {
	return &AuthenticationPolicyRuleApplyConfiguration{}
}

// WithFrom sets the From field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the From field is set to the value of the last call.
func (b *AuthenticationPolicyRuleApplyConfiguration) WithFrom(value *AuthenticationPolicySourceApplyConfiguration) *AuthenticationPolicyRuleApplyConfiguration {
	b.From = value
	return b
}

// WithTo sets the To field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the To field is set to the value of the last call.
func (b *AuthenticationPolicyRuleApplyConfiguration) WithTo(value *AuthenticationPolicyTargetApplyConfiguration) *AuthenticationPolicyRuleApplyConfiguration {
	b.To = value
	return b
}

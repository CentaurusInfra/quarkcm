/*
Copyright 2022 quarkcm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package event

import (
	"fmt"

	"github.com/CentaurusInfra/quarkcm/pkg/utils"
	api_v1 "k8s.io/api/core/v1"
)

// Event represent an event got from k8s api server
type Event struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Name      string
}

// New create new Event
func New(obj interface{}, action string) Event {
	var namespace, kind, component, host, reason, name string

	objectMeta := utils.GetObjectMetaData(obj)
	namespace = objectMeta.Namespace
	name = objectMeta.Name
	reason = action

	switch object := obj.(type) {
	case *api_v1.Pod:
		kind = "pod"
		host = object.Spec.NodeName
	case *api_v1.Node:
		kind = "node"
	case Event:
		name = object.Name
		kind = object.Kind
		namespace = object.Namespace
	}

	kbEvent := Event{
		Namespace: namespace,
		Kind:      kind,
		Component: component,
		Host:      host,
		Reason:    reason,
		Name:      name,
	}
	return kbEvent
}

// Message returns event message in standard format.
// included as a part of event packege to enhance code resuablity across handlers.
func (e *Event) Message() (msg string) {
	// using switch over if..else, since the format could vary based on the kind of the object in future.
	switch e.Kind {
	case "pod":
		msg = fmt.Sprintf(
			"A pod `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	case "node":
		msg = fmt.Sprintf(
			"A node `%s` has been `%s`",
			e.Name,
			e.Reason,
		)
	default:
		msg = fmt.Sprintf(
			"A `%s` in namespace `%s` has been `%s`:\n`%s`",
			e.Kind,
			e.Namespace,
			e.Reason,
			e.Name,
		)
	}
	return msg
}

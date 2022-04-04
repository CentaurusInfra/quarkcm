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

package handlers

import (
	"github.com/CentaurusInfra/quarkcm/pkg/objects"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type NodeHandler struct {
}

func (d *NodeHandler) Handle(eventItem objects.EventItem) {
	klog.Infof("Handle node event %s:%s, tracking: %s", eventItem.Key, eventItem.EventType, eventItem.Id)
	d.handleInternal(eventItem.EventType, eventItem.Obj.(*v1.Node))
}

func (d *NodeHandler) handleInternal(eventType string, node *v1.Node) {
}

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
	"github.com/CentaurusInfra/quarkcm/pkg/constants"
	"github.com/CentaurusInfra/quarkcm/pkg/datastore"
	"github.com/CentaurusInfra/quarkcm/pkg/objects"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type NodeHandler struct {
}

func (d *NodeHandler) Handle(eventItem objects.EventItem) {
	klog.Infof("Handle node event %s:%s. Tracking Id: %s", eventItem.Key, eventItem.EventType, eventItem.Id)
	if eventItem.EventType == constants.EventType_Delete {
		datastore.DeleteNode(eventItem.Key, eventItem.Id)
	} else {
		handleNodeSet(eventItem, eventItem.Obj.(*v1.Node))
	}
}

func handleNodeSet(eventItem objects.EventItem, node *v1.Node) {
	var hostname string
	var nodeIP string
	for _, item := range node.Status.Addresses {
		if item.Type == "InternalIP" {
			nodeIP = item.Address
		} else if item.Type == "Hostname" {
			hostname = item.Address
		}
	}
	datastore.SetNode(eventItem.Key, hostname, nodeIP, node.ObjectMeta.CreationTimestamp.Unix(), eventItem.Id)
}

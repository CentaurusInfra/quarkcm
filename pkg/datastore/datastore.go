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

package datastore

import (
	"encoding/json"
	"sync"

	"github.com/CentaurusInfra/quarkcm/pkg/constants"
	"github.com/CentaurusInfra/quarkcm/pkg/objects"
	"k8s.io/klog"
)

type DataStore struct {
	NodeResourceVersion int
	NodeMap             map[string]*objects.NodeObject   // map[node name] => node object
	NodeEventMap        map[int]*objects.NodeEventObject // map[resource version] => node event object
	PodResourceVersion  int
	PodMap              map[string]*objects.PodObject   // map[key] => pod object
	PodEventMap         map[int]*objects.PodEventObject // map[resource version] => pod event object
}

var lock = &sync.Mutex{}
var dataStore *DataStore

func Instance() *DataStore {
	if dataStore == nil {
		lock.Lock()
		defer lock.Unlock()
		if dataStore == nil {
			dataStore = &DataStore{
				NodeResourceVersion: 0,
				NodeMap:             map[string]*objects.NodeObject{},
				NodeEventMap:        map[int]*objects.NodeEventObject{},
				PodResourceVersion:  0,
				PodMap:              map[string]*objects.PodObject{},
				PodEventMap:         map[int]*objects.PodEventObject{},
			}
		}
	}
	return dataStore
}

func calculateNextNodeResourceVersion() int {
	instance := Instance()
	lock.Lock()
	defer lock.Unlock()
	instance.NodeResourceVersion += 1
	return instance.NodeResourceVersion
}

func calculateNextPodResourceVersion() int {
	instance := Instance()
	lock.Lock()
	defer lock.Unlock()
	instance.PodResourceVersion += 1
	return instance.PodResourceVersion
}

func SetNode(name string, nodeHostname string, nodeIP string, creationTimestamp int64, trackingId string) {
	nodeMap := Instance().NodeMap
	node, exists := nodeMap[name]
	changed := false
	if exists {
		if node.Hostname != nodeHostname || node.IP != nodeIP {
			changed = true
		} else {
			klog.Infof("Handling node completed. Node %s is unchanged. Tracking Id: %s", name, trackingId)
		}
	} else {
		changed = true
	}

	if changed {
		resourceVersion := calculateNextNodeResourceVersion()
		newNode := &objects.NodeObject{
			Name:              name,
			Hostname:          nodeHostname,
			IP:                nodeIP,
			CreationTimestamp: creationTimestamp,
			ResourceVersion:   resourceVersion,
		}
		newNodeEvent := &objects.NodeEventObject{
			ResourceVersion: resourceVersion,
			EventType:       constants.EventType_Set,
			NodeObject:      *newNode,
		}
		nodeMap[name] = newNode
		Instance().NodeEventMap[resourceVersion] = newNodeEvent

		nodeStr, _ := json.Marshal(nodeMap[name])
		klog.Infof("Handling node completed. Node set as %s. Tracking Id: %s", nodeStr, trackingId)
	}
}

func DeleteNode(name string, trackingId string) {
	nodeMap := Instance().NodeMap
	node, exists := nodeMap[name]
	if exists {
		resourceVersion := calculateNextNodeResourceVersion()
		newNodeEvent := &objects.NodeEventObject{
			ResourceVersion: resourceVersion,
			EventType:       constants.EventType_Delete,
			NodeObject:      *node,
		}
		Instance().NodeEventMap[resourceVersion] = newNodeEvent
		delete(nodeMap, name)
		klog.Infof("Handling node completed. Node %s is deleted. Tracking Id: %s", name, trackingId)
	}
}

func ListNode(minResourceVersion int) []objects.NodeEventObject {
	maxResourceVersion := Instance().NodeResourceVersion
	nodeEventMap := Instance().NodeEventMap

	var nodeEvents []objects.NodeEventObject
	for i := minResourceVersion + 1; i <= maxResourceVersion; i++ {
		nodeEvents = append(nodeEvents, *nodeEventMap[i])
	}
	return nodeEvents
}

func SetPod(key string, podIP string, nodeName string, trackingId string) {
	podMap := Instance().PodMap
	pod, exists := podMap[key]
	changed := false
	if exists {
		if pod.IP != podIP || pod.NodeName != nodeName {
			changed = true
		} else {
			klog.Infof("Handling pod completed. Pod %s is unchanged. Tracking Id: %s", key, trackingId)
		}
	} else {
		changed = true
	}

	if changed {
		resourceVersion := calculateNextPodResourceVersion()
		newPod := &objects.PodObject{
			Key:             key,
			IP:              podIP,
			NodeName:        nodeName,
			ResourceVersion: resourceVersion,
		}
		newPodEvent := &objects.PodEventObject{
			ResourceVersion: resourceVersion,
			EventType:       constants.EventType_Set,
			PodObject:       *newPod,
		}
		podMap[key] = newPod
		Instance().PodEventMap[resourceVersion] = newPodEvent

		podStr, _ := json.Marshal(podMap[key])
		klog.Infof("Handling pod completed. Pod set as %s. Tracking Id: %s", podStr, trackingId)
	}
}

func DeletePod(key string, trackingId string) {
	podMap := Instance().PodMap
	pod, exists := podMap[key]
	if exists {
		resourceVersion := calculateNextPodResourceVersion()
		newPodEvent := &objects.PodEventObject{
			ResourceVersion: resourceVersion,
			EventType:       constants.EventType_Delete,
			PodObject:       *pod,
		}
		Instance().PodEventMap[resourceVersion] = newPodEvent
		delete(podMap, key)
		klog.Infof("Handling pod completed. Pod %s is deleted. Tracking Id: %s", key, trackingId)
	}
}

func ListPod(minResourceVersion int) []objects.PodEventObject {
	maxResourceVersion := Instance().PodResourceVersion
	podEventMap := Instance().PodEventMap

	var podEvents []objects.PodEventObject
	for i := minResourceVersion + 1; i <= maxResourceVersion; i++ {
		podEvents = append(podEvents, *podEventMap[i])
	}
	return podEvents
}

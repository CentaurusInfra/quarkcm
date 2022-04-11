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

	"github.com/CentaurusInfra/quarkcm/pkg/objects"
	"k8s.io/klog"
)

type DataStore struct {
	ResourceVersion int
	NodeMap         map[string]*objects.NodeObject // map[node name] => node object
	PodMap          map[string]*objects.PodObject  // map[key] => pod object
}

var lock = &sync.Mutex{}
var dataStore *DataStore

func Instance() *DataStore {
	if dataStore == nil {
		lock.Lock()
		defer lock.Unlock()
		if dataStore == nil {
			dataStore = &DataStore{
				ResourceVersion: 0,
				NodeMap:         map[string]*objects.NodeObject{},
				PodMap:          map[string]*objects.PodObject{},
			}
		}
	}
	return dataStore
}

func getNextResourceVersion() int {
	instance := Instance()
	lock.Lock()
	defer lock.Unlock()
	instance.ResourceVersion += 1
	return instance.ResourceVersion
}

func SetNode(key string, nodeHostname string, nodeIP string, creationTimestamp string, trackingId string) {
	nodeMap := Instance().NodeMap
	node, exists := nodeMap[key]
	changed := false
	if exists {
		if node.Hostname != nodeHostname || node.IP != nodeIP {
			node.Hostname = nodeHostname
			node.IP = nodeIP
			node.CreationTimestamp = creationTimestamp
			node.ResourceVersion = getNextResourceVersion()
			changed = true
		} else {
			klog.Infof("Handling node completed. Node %s is unchanged. Tracking Id: %s", key, trackingId)
		}
	} else {
		nodeMap[key] = &objects.NodeObject{
			Name:              key,
			Hostname:          nodeHostname,
			IP:                nodeIP,
			CreationTimestamp: creationTimestamp,
			ResourceVersion:   getNextResourceVersion(),
		}
		changed = true
	}
	if changed {
		nodeStr, _ := json.Marshal(nodeMap[key])
		klog.Infof("Handling node completed. Node set as %s. Tracking Id: %s", nodeStr, trackingId)
	}
}

func DeleteNode(key string, trackingId string) {
	nodeMap := Instance().NodeMap
	_, exists := nodeMap[key]
	if exists {
		delete(nodeMap, key)
		klog.Infof("Handling node completed. Node %s is deleted. Tracking Id: %s", key, trackingId)
	}
}

func SetPod(key string, podIP string, nodeName string, trackingId string) {
	podMap := Instance().PodMap
	pod, exists := podMap[key]
	changed := false
	if exists {
		if pod.IP != podIP || pod.NodeName != nodeName {
			pod.IP = podIP
			pod.NodeName = nodeName
			pod.ResourceVersion = getNextResourceVersion()
			changed = true
		} else {
			klog.Infof("Handling pod completed. Pod %s is unchanged. Tracking Id: %s", key, trackingId)
		}
	} else {
		podMap[key] = &objects.PodObject{
			Key:             key,
			IP:              podIP,
			NodeName:        nodeName,
			ResourceVersion: getNextResourceVersion(),
		}
		changed = true
	}
	if changed {
		podStr, _ := json.Marshal(podMap[key])
		klog.Infof("Handling pod completed. Pod set as %s. Tracking Id: %s", podStr, trackingId)
	}
}

func DeletePod(key string, trackingId string) {
	podMap := Instance().PodMap
	_, exists := podMap[key]
	if exists {
		delete(podMap, key)
		klog.Infof("Handling pod completed. Pod %s is deleted. Tracking Id: %s", key, trackingId)
	}
}

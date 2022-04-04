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

package controller

import (
	"github.com/CentaurusInfra/quarkcm/pkg/handlers"
	"github.com/CentaurusInfra/quarkcm/pkg/utils"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/klog"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func NewPodController(client kubernetes.Interface) *Controller {
	var eventHandler handlers.Handler = new(handlers.PodHandler)
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Pods("").List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Pods("").Watch(options)
			},
		},
		&api_v1.Pod{},
		0,
		cache.Indexers{},
	)
	return newResourceController(client, eventHandler, informer, "pod")
}

func NewNodeController(client kubernetes.Interface) *Controller {
	var eventHandler handlers.Handler = new(handlers.NodeHandler)
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return client.CoreV1().Nodes().List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.CoreV1().Nodes().Watch(options)
			},
		},
		&api_v1.Node{},
		0,
		cache.Indexers{},
	)
	return newResourceController(client, eventHandler, informer, "node")
}

func newResourceController(client kubernetes.Interface, eventHandler handlers.Handler, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var eventItem EventItem
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			eventItem.key, err = cache.MetaNamespaceKeyFunc(obj)
			eventItem.eventType = "create"
			eventItem.resourceType = resourceType
			klog.Infof("Processing add to %v: %s", resourceType, eventItem.key)
			if err == nil {
				queue.Add(eventItem)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			eventItem.key, err = cache.MetaNamespaceKeyFunc(old)
			eventItem.eventType = "update"
			eventItem.resourceType = resourceType
			klog.Infof("Processing update to %v: %s", resourceType, eventItem.key)
			if err == nil {
				queue.Add(eventItem)
			}
		},
		DeleteFunc: func(obj interface{}) {
			eventItem.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			eventItem.eventType = "delete"
			eventItem.resourceType = resourceType
			eventItem.namespace = utils.GetObjectMetaData(obj).Namespace
			klog.Infof("Processing delete to %v: %s", resourceType, eventItem.key)
			if err == nil {
				queue.Add(eventItem)
			}
		},
	})

	return &Controller{
		resourceType: resourceType,
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

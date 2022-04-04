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
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/CentaurusInfra/quarkcm/pkg/event"
	"github.com/CentaurusInfra/quarkcm/pkg/handlers"
	"github.com/CentaurusInfra/quarkcm/pkg/utils"
	"k8s.io/klog"

	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 5

var serverStartTime time.Time

type EventItem struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	resourceType string
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler handlers.Handler
}

// Start prepares watchers and run their controllers, then waits for process termination signals
func Start(eventHandler handlers.Handler) {
	var kubeClient kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		kubeClient = utils.GetClientOutOfCluster()
	} else {
		kubeClient = utils.GetClient()
	}

	podInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Pods("").List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Pods("").Watch(options)
			},
		},
		&api_v1.Pod{},
		0, //Skip resync
		cache.Indexers{},
	)

	podController := newResourceController(kubeClient, eventHandler, podInformer, "pod")
	podStopCh := make(chan struct{})
	defer close(podStopCh)

	go podController.Run(podStopCh)

	nodeInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				return kubeClient.CoreV1().Nodes().List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return kubeClient.CoreV1().Nodes().Watch(options)
			},
		},
		&api_v1.Node{},
		0, //Skip resync
		cache.Indexers{},
	)

	nodeController := newResourceController(kubeClient, eventHandler, nodeInformer, "node")
	nodeStopCh := make(chan struct{})
	defer close(nodeStopCh)

	go nodeController.Run(nodeStopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
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

// Run starts the quarkcm controller
func (c *Controller) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Infof("Starting quarkcm %s controller", c.resourceType)
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("%s controller timed out waiting for caches to sync", c.resourceType))
		return
	}

	klog.Infof("quarkcm %s controller synced and ready", c.resourceType)

	wait.Until(c.runWorker, time.Second, stopCh)
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	queueItem, quit := c.queue.Get()
	eventItem := queueItem.(EventItem)

	if quit {
		return false
	}
	defer c.queue.Done(queueItem)
	err := c.processItem(eventItem)
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(queueItem)
	} else if c.queue.NumRequeues(queueItem) < maxRetries {
		klog.Errorf("error processing %s (will retry): %v", eventItem.key, err)
		c.queue.AddRateLimited(queueItem)
	} else {
		// err != nil and too many retries
		klog.Errorf("error processing %s (giving up): %v", eventItem.key, err)
		c.queue.Forget(queueItem)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processItem(eventItem EventItem) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(eventItem.key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", eventItem.key, err)
	}
	// get object's metedata
	objectMeta := utils.GetObjectMetaData(obj)

	// namespace retrived from event key incase namespace value is empty
	if eventItem.namespace == "" && strings.Contains(eventItem.key, "/") {
		substring := strings.Split(eventItem.key, "/")
		eventItem.namespace = substring[0]
		eventItem.key = substring[1]
	}

	// process events based on its type
	switch eventItem.eventType {
	case "create":
		if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			kbEvent := event.Event{
				Name:      objectMeta.Name,
				Namespace: eventItem.namespace,
				Kind:      eventItem.resourceType,
				Reason:    "Created",
			}
			c.eventHandler.Handle(kbEvent)
			return nil
		}
	case "update":
		kbEvent := event.Event{
			Name:      eventItem.key,
			Namespace: eventItem.namespace,
			Kind:      eventItem.resourceType,
			Reason:    "Updated",
		}
		c.eventHandler.Handle(kbEvent)
		return nil
	case "delete":
		kbEvent := event.Event{
			Name:      eventItem.key,
			Namespace: eventItem.namespace,
			Kind:      eventItem.resourceType,
			Reason:    "Deleted",
		}
		c.eventHandler.Handle(kbEvent)
		return nil
	}
	return nil
}

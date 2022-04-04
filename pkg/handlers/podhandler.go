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
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	"github.com/CentaurusInfra/quarkcm/pkg/event"
)

type PodHandler struct {
}

func (d *PodHandler) Init() error {
	return nil
}

func (d *PodHandler) Handle(e event.Event) {
	d.handleInternal(e, e.Obj.(*v1.Pod))
}

func (d *PodHandler) handleInternal(e event.Event, pod *v1.Pod) {
	klog.Infof("Handle pod event %s", pod.Name)
}

/*
Copyright 2017 Caicloud authors. All rights reserved.

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

package syncqueue

import (
	"errors"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func TestSyncQueue_Enqueue(t *testing.T) {

	syncPods := func(obj interface{}) error {
		pod := obj.(*v1.Pod)
		pod.Name = pod.Name + "_synced"
		return nil
	}

	queue := NewSyncQueueForKeyFunc(&v1.Pod{}, syncPods, PassthroughKeyFunc)
	stopCh := make(chan struct{})
	defer func() {
		close(stopCh)
		queue.ShutDown()
	}()
	queue.Run(1, stopCh)

	tests := []struct {
		podName string
		want    string
	}{
		{"test", "test_synced"},
	}
	for _, tt := range tests {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: tt.podName,
			},
		}
		queue.Enqueue(pod)

		time.Sleep(1 * time.Millisecond)
		if pod.Name != tt.want {
			t.Errorf("SyncQueque.Enqueque() == %v, want %v", pod.Name, tt.want)
		}

	}
}

func TestSyncQueue_EnqueueError(t *testing.T) {
	syncError := func(obj interface{}) error {
		pod := obj.(*v1.Pod)
		if pod.Name == "test" {
			pod.Name = "test_1"
			return errors.New("error")
		}
		if pod.Name == "test_1" {
			pod.Name = "test_synced"
			return nil
		}
		return nil
	}
	queue := NewSyncQueueForKeyFunc(&v1.Pod{}, syncError, PassthroughKeyFunc)
	queue.SetMaxRetries(1)
	stopCh := make(chan struct{})
	defer func() {
		close(stopCh)
		queue.ShutDown()
	}()
	queue.Run(1, stopCh)

	tests := []struct {
		podName string
		want    string
	}{
		{"test", "test_synced"},
	}
	for _, tt := range tests {
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: tt.podName,
			},
		}
		queue.Enqueue(pod)

		time.Sleep(10 * time.Millisecond)
		if pod.Name != tt.want {
			t.Errorf("SyncQueque.Enqueque() == %v, want %v", pod.Name, tt.want)
		}

	}
}

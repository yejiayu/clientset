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
	"fmt"
	"reflect"
	"sync"
	"time"

	log "github.com/zoumo/logdog"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	maxRetries = 3
)

var (
	// KeyFunc is the default key function
	KeyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

type syncHandler func(obj interface{}) error

type keyFunc func(obj interface{}) (interface{}, error)

// PassthroughKeyFunc is a keyFunc which returns the original obj
func PassthroughKeyFunc(obj interface{}) (interface{}, error) {
	return obj, nil
}

// SyncQueue is a helper for creating a kubernetes controller easily
// It requires a rate limit workqueue , a syncHandler and an optional key function.
// After running the syncQueue, you can call it's Enqueque function to enqueue items.
// SyncQueue will get key from the items by keyFunc, and add the key to the rate limit workqueue.
// The worker will be invoked to call the syncHandler.
type SyncQueue struct {
	// SyncType is the object type in the queue
	SyncType reflect.Type
	// queue is the work queue the worker polls
	Queue workqueue.RateLimitingInterface
	// SyncHandler is called for each item in the queue
	SyncHandler syncHandler
	// KeyFunc is called to get key from obj
	keyFunc keyFunc

	waitGroup sync.WaitGroup

	maxRetries int
}

// NewSyncQueue returns a new SyncQueue, enqueue key of obj using default keyFunc
func NewSyncQueue(syncObject runtime.Object, queue workqueue.RateLimitingInterface, syncHandler syncHandler) *SyncQueue {
	return NewSyncQueueForKeyFunc(syncObject, queue, syncHandler, nil)
}

// NewSyncQueueForKeyFunc returns a new SyncQueue using custom keyFunc
func NewSyncQueueForKeyFunc(syncObject runtime.Object, queue workqueue.RateLimitingInterface, syncHandler syncHandler, keyFunc keyFunc) *SyncQueue {
	sq := &SyncQueue{
		SyncType:    reflect.TypeOf(syncObject),
		Queue:       queue,
		SyncHandler: syncHandler,
		keyFunc:     keyFunc,
		waitGroup:   sync.WaitGroup{},
		maxRetries:  maxRetries,
	}

	if keyFunc == nil {
		sq.keyFunc = sq.defaultKeyFunc
	}

	return sq
}

// SetMaxRetries sets the max retry times of the queue
func (sq *SyncQueue) SetMaxRetries(max int) {
	if max > 0 {
		sq.maxRetries = max
	}
}

// Run starts n workers to sync
func (sq *SyncQueue) Run(workers int, stopCh <-chan struct{}) {
	for i := 0; i < workers; i++ {
		go wait.Until(sq.worker, time.Second, stopCh)
	}
}

// Enqueue wraps queue.Add
func (sq *SyncQueue) Enqueue(obj interface{}) {

	if sq.IsShuttingDown() {
		return
	}

	key, err := sq.keyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for %v %#v: %v", sq.SyncType, obj, err))
		return
	}
	sq.Queue.Add(key)
}

// EnqueueRateLimited wraps queue.AddRateLimited. It adds an item to the workqueue
// after the rate limiter says its ok
func (sq *SyncQueue) EnqueueRateLimited(obj interface{}) {

	if sq.IsShuttingDown() {
		return
	}

	key, err := sq.keyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for %v %#v: %v", sq.SyncType, obj, err))
		return
	}
	sq.Queue.AddRateLimited(key)
}

// EnqueueAfter wraps queue.AddAfter. It adds an item to the workqueue after the indicated duration has passed
func (sq *SyncQueue) EnqueueAfter(obj interface{}, after time.Duration) {

	if sq.IsShuttingDown() {
		return
	}

	key, err := sq.keyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for %v %#v: %v", sq.SyncType, obj, err))
		return
	}
	sq.Queue.AddAfter(key, after)
}

func (sq *SyncQueue) defaultKeyFunc(obj interface{}) (interface{}, error) {
	key, err := KeyFunc(obj)
	if err != nil {
		return "", err
	}
	return key, nil
}

// Worker is a common worker for controllers
// worker runs a worker thread that just dequeues items, processes them, and marks them done.
// It enforces that the syncHandler is never invoked concurrently with the same key.
func (sq *SyncQueue) worker() {
	sq.waitGroup.Add(1)
	defer sq.waitGroup.Done()
	// invoked oncely process any until exhausted
	for sq.processNextWorkItem() {
	}
}

// ProcessNextWorkItem processes next item in queue by syncHandler
func (sq *SyncQueue) processNextWorkItem() bool {
	obj, quit := sq.Queue.Get()
	if quit {
		return false
	}
	defer sq.Queue.Done(obj)

	err := sq.SyncHandler(obj)
	sq.handleSyncError(err, obj)

	return true
}

// HandleSyncError handles error when sync obj error and retry n times
func (sq *SyncQueue) handleSyncError(err error, obj interface{}) {
	if err == nil {
		// no err
		sq.Queue.Forget(obj)
		return
	}

	var key interface{}

	// get short key no matter what the keyfunc is
	key, kerr := KeyFunc(obj)
	if kerr != nil {
		key = obj
	}

	if sq.Queue.NumRequeues(obj) < sq.maxRetries {
		log.Warn("Error syncing object, retry", log.Fields{"type": sq.SyncType, "obj": key, "err": err})
		sq.Queue.AddRateLimited(obj)
		return
	}

	utilruntime.HandleError(err)
	log.Warn("Dropping object out of queue", log.Fields{"type": sq.SyncType, "obj": key, "err": err})
	sq.Queue.Forget(obj)
}

// ShutDown shuts down the work queue and waits for the worker to ACK
func (sq *SyncQueue) ShutDown() {
	// sq shutdown the queue, then worker can't get key from queue
	// processNextWorkItem return false, and then waitGroup -1
	sq.Queue.ShutDown()
	sq.waitGroup.Wait()
}

// IsShuttingDown returns if the method Shutdown was invoked
func (sq *SyncQueue) IsShuttingDown() bool {
	return sq.Queue.ShuttingDown()
}

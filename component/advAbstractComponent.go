// Copyright (c) 2013-2020 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package Component

import (
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"sync"
)

// AdvAbstractComponent are considered newer type of components used by pipelines that instantiate first
// then run with potential modification. Thus, rwmutex are needed
type AdvAbstractComponent struct {
	id     string
	logger *log.CommonLogger

	pipelines      map[string]common.Pipeline
	eventListeners map[string]EventListenersMap
	mutex          sync.RWMutex
}

func NewAdvAbstractComponentWithLogger(id string, logger *log.CommonLogger) *AdvAbstractComponent {
	return &AdvAbstractComponent{
		id:             id,
		pipelines:      make(map[string]common.Pipeline),
		eventListeners: make(map[string]EventListenersMap),
		logger:         logger,
	}
}

func (c *AdvAbstractComponent) Id() string {
	return c.id
}

// Returns nil if not found
func (c *AdvAbstractComponent) GetPipeline(topic string) common.Pipeline {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.pipelines[topic]
}

func (c *AdvAbstractComponent) SpecificAsyncComponentEventListeners(topic string) map[string]common.AsyncComponentEventListener {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	listeners, ok := c.eventListeners[topic]
	if !ok {
		return nil
	}
	listenerMap := make(map[string]common.AsyncComponentEventListener)
	listeners.exportToMap(listenerMap)
	return listenerMap
}

func (c *AdvAbstractComponent) RegisterComponentEventListener(eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	return base.ErrorNotImplemented
}

func (c *AdvAbstractComponent) UnRegisterComponentEventListener(eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	return base.ErrorNotImplemented
}

func (c *AdvAbstractComponent) RegisterSpecificComponentEventListener(topic string, eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	c.mutex.RLock()

	v, ok := c.eventListeners[topic]
	if !ok {
		c.mutex.RUnlock()
		c.mutex.Lock()
		_, ok = c.eventListeners[topic]
		if !ok {
			c.eventListeners[topic] = make(EventListenersMap)
			v = c.eventListeners[topic]
		}
		v.registerListerNoLock(eventType, listener)
		c.mutex.Unlock()
		return nil
	}

	v.registerListerNoLock(eventType, listener)
	c.mutex.RUnlock()
	return nil
}

func (c *AdvAbstractComponent) UnRegisterSpecificComponentEventListener(topic string, eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	c.mutex.RLock()

	v, ok := c.eventListeners[topic]
	if !ok {
		c.mutex.RUnlock()
		return base.ErrorInvalidInput
	}

	err := v.unregisterEventListenerNoLock(eventType, listener)
	if err != nil {
		c.mutex.RUnlock()
		return err
	}

	if len(c.eventListeners[topic]) == 0 {
		c.mutex.RUnlock()
		c.mutex.Lock()
		if len(c.eventListeners[topic]) == 0 {
			delete(c.eventListeners, topic)
		}
		c.mutex.Unlock()
		return nil
	}

	c.mutex.RUnlock()
	return nil
}

func (c *AdvAbstractComponent) RaiseEvent(event *common.Event) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, listener := range c.eventListeners {
		listener.raiseEvent(event)
	}
}

func (c *AdvAbstractComponent) AsyncComponentEventListeners() map[string]common.AsyncComponentEventListener {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	listenersMap := make(map[string]common.AsyncComponentEventListener)
	for _, listener := range c.eventListeners {
		listener.exportToMap(listenersMap)
	}

	return listenersMap
}

func (c *AdvAbstractComponent) Logger() *log.CommonLogger {
	return c.logger
}
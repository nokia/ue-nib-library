/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader_test

import (
	"errors"
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenibreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var someGNb string
var someEvent string
var someEventCategory uenibreader.EventCategory
var someChannel string
var anotherEvent string

func init() {
	someGNb = "somegnb:310-410-b5c67788"
	someEvent = "someEvent"
	someEventCategory = uenibreader.DualConnectivity
	someChannel = someGNb + "_" + someEventCategory.String()
	anotherEvent = "anotherEvent"
}

type eventCbArgs struct {
	gNb           string
	eventCategory uenibreader.EventCategory
	events        []string
}

type eventTracker struct {
	callCount    int
	eventCbCalls []eventCbArgs
}

func (handler *eventTracker) callback(gNb string, eventCategory uenibreader.EventCategory, events []string) {
	handler.callCount = handler.callCount + 1
	handler.eventCbCalls = append(handler.eventCbCalls, eventCbArgs{gNb, eventCategory, events})
}

func (handler *eventTracker) verify(t *testing.T, expectedCallCount int, expectedEvents ...eventCbArgs) {
	assert.Equal(t, expectedCallCount, handler.callCount)
	assert.Equal(t, len(expectedEvents), len(handler.eventCbCalls))
	for i, expectedEvent := range expectedEvents {
		assert.Equal(t, expectedEvent, handler.eventCbCalls[i])
	}
}

func TestSubscribeEventsCanSubscribeDualConnectivityEventCategory(t *testing.T) {
	m, i := setup()
	eventCategory := uenibreader.DualConnectivity
	m.On("SubscribeChannel", mock.AnythingOfType("func(string, ...string)"),
		[]string{someGNb + "_DUAL_CONNECTIVITY"}).Return(nil).Once()

	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{eventCategory},
		func(string, uenibreader.EventCategory, []string) {})
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestSubscribeEventsFailure(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Backend Error")
	m.On("SubscribeChannel", mock.AnythingOfType("func(string, ...string)"),
		[]string{someChannel}).Return(dbError).Once()

	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{someEventCategory},
		func(string, uenibreader.EventCategory, []string) {})
	assert.NotNil(t, err)
	uenibError, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, uenibError.Temporary())
	assert.Contains(t, err.Error(), "database backend error: Some DB Backend Error")
	m.AssertExpectations(t)
}

func TestSubscribeEventsGetOneEvent(t *testing.T) {
	m, i := setup()
	tracker := eventTracker{}
	var storedSdlCallback func(string, ...string)
	m.On("SubscribeChannel", mock.AnythingOfType("func(string, ...string)"),
		[]string{someChannel}).Run(func(args mock.Arguments) {
		storedSdlCallback = args.Get(0).(func(string, ...string))
	}).Return(nil).Once()

	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{someEventCategory}, tracker.callback)
	assert.Nil(t, err)
	storedSdlCallback(someChannel, someEvent)
	tracker.verify(t, 1, eventCbArgs{someGNb, someEventCategory, []string{someEvent}})
	m.AssertExpectations(t)
}

func TestSubscribeEventsUnknownEventCategoryFailure(t *testing.T) {
	m, i := setup()
	tracker := eventTracker{}
	unknownEventCategory := uenibreader.EventCategory(9999999999)
	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{unknownEventCategory}, tracker.callback)
	tracker.verify(t, 0)
	uenibError, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, uenibError.Temporary())
	assert.Contains(t, err.Error(), fmt.Sprintf("validation error: Unknown event category ID: %d", unknownEventCategory))
	m.AssertExpectations(t)
}

func TestSubscribeEventsTwoEventsInOneEvent(t *testing.T) {
	m, i := setup()
	tracker := eventTracker{}
	var storedSdlCallback func(string, ...string)
	m.On("SubscribeChannel", mock.AnythingOfType("func(string, ...string)"), []string{someChannel}).Run(
		func(args mock.Arguments) {
			storedSdlCallback = args.Get(0).(func(string, ...string))
		}).Return(nil).Once()
	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{someEventCategory}, tracker.callback)
	assert.Nil(t, err)
	storedSdlCallback(someChannel, someEvent, anotherEvent)
	tracker.verify(t, 1, eventCbArgs{someGNb, someEventCategory, []string{someEvent, anotherEvent}})
	m.AssertExpectations(t)
}

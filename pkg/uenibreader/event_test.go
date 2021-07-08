/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader_test

import (
	"errors"
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenib"
	"github.com/nokia/ue-nib-library/pkg/uenibreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var someEvNs string
var someGNb string
var someEventCategory uenibreader.EventCategory
var someChannel string
var dcAddEvent string
var expParsedDcAddEvent uenibreader.DcEvent
var dcRemoveEvent string
var expParsedDcRemoveEvent uenibreader.DcEvent
var dcS1ULTunnelEstablishEvent string
var expParsedDcS1ULTunnelEstablishEvent uenibreader.DcEvent
var dcS1ULTunnelReleaseEvent string
var expParsedDcS1ULTunnelReleaseEvent uenibreader.DcEvent
var dcRemoveAllUesEvent string
var expParsedDcRemoveAllUesEvent uenibreader.DcEvent
var dcEmptyS1ULTunnelEstablishEvent string
var expParsedEmptyDcS1ULTunnelEstablishEvent uenibreader.DcEvent

func init() {
	someGNb = "somegnb:310-410-b5c67788"
	someEvNs = "uenib/" + someGNb
	someEventCategory = uenibreader.DualConnectivity
	someChannel = someGNb + "_" + someEventCategory.String()

	dcAddEvent = "somegnb:310-410-b5c67788#100#200_ADD"
	expParsedDcAddEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_ADD,
		UeID: uenib.UeID{
			GNb:         someGNb,
			GNbUeX2ApID: "100",
			ENbUeX2ApID: "200",
		},
	}

	dcRemoveEvent = "somegnb:310-410-b5c67788#100#200_REMOVE"
	expParsedDcRemoveEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_REMOVE,
		UeID: uenib.UeID{
			GNb:         someGNb,
			GNbUeX2ApID: "100",
			ENbUeX2ApID: "200",
		},
	}

	dcS1ULTunnelEstablishEvent = "somegnb:310-410-b5c67788#100#200_10.20.30.40#5000#20.30.40.50#6000_S1UL_TUNNEL_ESTABLISH"
	expParsedDcS1ULTunnelEstablishEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_S1UL_TUNNEL_ESTABLISH,
		UeID: uenib.UeID{
			GNb:         someGNb,
			GNbUeX2ApID: "100",
			ENbUeX2ApID: "200",
		},
		S1ULGtpTunnels: []uenibreader.DcEventTunnel{
			uenibreader.DcEventTunnel{
				Addr: "10.20.30.40",
				Teid: 5000,
			},
			uenibreader.DcEventTunnel{
				Addr: "20.30.40.50",
				Teid: 6000,
			},
		},
	}

	dcS1ULTunnelReleaseEvent = "somegnb:310-410-b5c67788#100#200_20.30.40.50#6000_S1UL_TUNNEL_RELEASE"
	expParsedDcS1ULTunnelReleaseEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_S1UL_TUNNEL_RELEASE,
		UeID: uenib.UeID{
			GNb:         someGNb,
			GNbUeX2ApID: "100",
			ENbUeX2ApID: "200",
		},
		S1ULGtpTunnels: []uenibreader.DcEventTunnel{
			uenibreader.DcEventTunnel{
				Addr: "20.30.40.50",
				Teid: 6000,
			},
		},
	}

	dcRemoveAllUesEvent = "GNB_ALL_UES_REMOVE"
	expParsedDcRemoveAllUesEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_GNB_ALL_UES_REMOVE,
	}

	dcEmptyS1ULTunnelEstablishEvent = "somegnb:310-410-b5c67788#100#200_#_S1UL_TUNNEL_ESTABLISH"
	expParsedEmptyDcS1ULTunnelEstablishEvent = uenibreader.DcEvent{
		EventType: uenibreader.DC_EVENT_S1UL_TUNNEL_ESTABLISH,
		UeID: uenib.UeID{
			GNb:         someGNb,
			GNbUeX2ApID: "100",
			ENbUeX2ApID: "200",
		},
		S1ULGtpTunnels: []uenibreader.DcEventTunnel{
			uenibreader.DcEventTunnel{},
		},
	}
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
	m.On("SubscribeChannel", someEvNs, mock.AnythingOfType("func(string, ...string)"),
		[]string{someGNb + "_DUAL_CONNECTIVITY"}).Return(nil).Once()

	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{eventCategory},
		func(string, uenibreader.EventCategory, []string) {})
	assert.Nil(t, err)
	m.AssertExpectations(t)
}

func TestSubscribeEventsFailure(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Backend Error")
	m.On("SubscribeChannel", someEvNs, mock.AnythingOfType("func(string, ...string)"),
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
	m.On("SubscribeChannel", someEvNs, mock.AnythingOfType("func(string, ...string)"),
		[]string{someChannel}).Run(func(args mock.Arguments) {
		storedSdlCallback = args.Get(1).(func(string, ...string))
	}).Return(nil).Once()

	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{someEventCategory}, tracker.callback)
	assert.Nil(t, err)
	storedSdlCallback(someChannel, dcAddEvent)
	tracker.verify(t, 1, eventCbArgs{someGNb, someEventCategory, []string{dcAddEvent}})
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
	m.On("SubscribeChannel", someEvNs, mock.AnythingOfType("func(string, ...string)"), []string{someChannel}).Run(
		func(args mock.Arguments) {
			storedSdlCallback = args.Get(1).(func(string, ...string))
		}).Return(nil).Once()
	err := i.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{someEventCategory}, tracker.callback)
	assert.Nil(t, err)
	storedSdlCallback(someChannel, dcAddEvent, dcRemoveEvent)
	tracker.verify(t, 1, eventCbArgs{someGNb, someEventCategory, []string{dcAddEvent, dcRemoveEvent}})
	m.AssertExpectations(t)
}

func TestParseDcEventSuccessForAddEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcAddEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedDcAddEvent, retEvt)
}

func TestParseDcEventSuccessForRemoveEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcRemoveEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedDcRemoveEvent, retEvt)
}

func TestParseDcEventSuccessForS1ULTunnelEstablishEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcS1ULTunnelEstablishEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedDcS1ULTunnelEstablishEvent, retEvt)
}

func TestParseDcEventSuccessForS1ULTunnelReleaseEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcS1ULTunnelReleaseEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedDcS1ULTunnelReleaseEvent, retEvt)
}

func TestParseDcEventSuccessForRemoveAllUesEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcRemoveAllUesEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedDcRemoveAllUesEvent, retEvt)
}

func TestParseDcEventPassThroughWithSuccessForUnknownEvent(t *testing.T) {
	var unknownEvent string = "somegnb:310-410-b5c67788#100#200_SOME_UNKNOWN_EVENT"
	retEvt, err := uenibreader.ParseDcEvent(unknownEvent)
	assert.Nil(t, err)
	assert.Equal(t, uenibreader.DC_EVENT_UNKNOWN, retEvt.EventType)
}

func TestParseDcEventReturnsErrorIfNoGnbInEvent(t *testing.T) {
	var illegalEvent string = "100#200_ADD"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventReturnsErrorIfNoGnbInTunnelEvent(t *testing.T) {
	var illegalEvent string = "100#200_20.30.40.50#6000_S1UL_TUNNEL_RELEASE"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventReturnsErrorIfNoTunnelFieldInEvent(t *testing.T) {
	var illegalEvent string = "somegnb:310-410-b5c67788#100#200_S1UL_TUNNEL_ESTABLISH"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventReturnsErrorIfNoTunnelIpInEvent(t *testing.T) {
	var illegalEvent string = "somegnb:310-410-b5c67788#100#200_10.20.30.40#5000#6000_S1UL_TUNNEL_ESTABLISH"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventReturnsErrorIfNoTunnelTeidInEvent(t *testing.T) {
	var illegalEvent string = "somegnb:310-410-b5c67788#100#200_10.20.30.40#5000#20.30.40.50_S1UL_TUNNEL_ESTABLISH"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventReturnsErrorIfTunnelTeidExceedsMaxU32InEvent(t *testing.T) {
	var illegalEvent string = "somegnb:310-410-b5c67788#100#200_10.20.30.40#4294967296_S1UL_TUNNEL_RELEASE"
	_, err := uenibreader.ParseDcEvent(illegalEvent)
	assert.NotNil(t, err)
}

func TestParseDcEventPassThroughWithSuccessIfEmptyTunnelIpAndTeidInEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent(dcEmptyS1ULTunnelEstablishEvent)
	assert.Nil(t, err)
	assert.Equal(t, expParsedEmptyDcS1ULTunnelEstablishEvent, retEvt)
}

func TestParseDcEventPassThroughWithSuccessIfEmptyEvent(t *testing.T) {
	retEvt, err := uenibreader.ParseDcEvent("")
	assert.Nil(t, err)
	assert.Equal(t, uenibreader.DcEvent{}, retEvt)
}

func TestDcEventStringPanicsIfStringMapEntryNotFound(t *testing.T) {
	var evt uenibreader.DcEventType = uenibreader.DC_EVENT_GNB_ALL_UES_REMOVE + 1
	assert.Panics(t, func() { evt.String() },
		"Too big event type didn't cause panic. Check event string map implementation")
}

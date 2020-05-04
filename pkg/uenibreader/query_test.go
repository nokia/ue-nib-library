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
	"testing"
)

var someUeID uenib.UeID
var anotherUeID uenib.UeID
var someErabID uenib.ErabID
var anotherErabID uenib.ErabID
var someDbKeyGNbUeX2ApID string
var someDbKeyENbUeX2ApID string
var someDbKeyPsCellPci string
var someDbKeyPsCellSsbFreq string
var someDbKeyUeStateEvent string
var someDbKeyUeStateCause string
var someDbKeyBearerIDs string
var someDbKeyBearerDrbID string
var someDbKeyBearerS1ULTepAddr string
var someDbKeyBearerS1ULTepTeid string
var someDbKeyBearerArpPL string
var someDbKeyBearerQci string

var anotherDbKeyBearerDrbID string
var anotherDbKeyBearerS1ULTepAddr string
var anotherDbKeyBearerS1ULTepTeid string
var anotherDbKeyBearerArpPL string
var anotherDbKeyBearerQci string

func init() {
	someUeID = uenib.UeID{
		GNb:         "somegnb:310-410-b5c67788",
		ENbUeX2ApID: "100",
	}
	anotherUeID = uenib.UeID{
		GNb:         "somegnb:310-410-b5c67788",
		GNbUeX2ApID: "200",
	}
	someErabID = 1000
	anotherErabID = 2000

	someDbKeyENbUeX2ApID = someUeID.GNb + "," + anotherUeID.GNbUeX2ApID + ",UEMAP_ENBUEX2APID"
	someDbKeyGNbUeX2ApID = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UEMAP_GNBUEX2APID"
	someDbKeyPsCellPci = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UE_PSCELL_PCI"
	someDbKeyPsCellSsbFreq = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UE_PSCELL_FREQ"
	someDbKeyUeStateEvent = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UE_STATE_EVENT"
	someDbKeyUeStateCause = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UE_STATE_CAUSE"
	someDbKeyBearerIDs = someUeID.GNb + "," + someUeID.ENbUeX2ApID + ",UE_ERAB_IDS"

	someDbKeyBearerDrbID = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(someErabID) + ",UE_ERAB_DRB_ID"
	someDbKeyBearerS1ULTepAddr = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(someErabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_ADDR"
	someDbKeyBearerS1ULTepTeid = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(someErabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_TEID"
	someDbKeyBearerArpPL = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(someErabID) + ",UE_ERAB_QOS_ARP_PL"
	someDbKeyBearerQci = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(someErabID) + ",UE_ERAB_QOS_QCI"

	anotherDbKeyBearerDrbID = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(anotherErabID) + ",UE_ERAB_DRB_ID"
	anotherDbKeyBearerS1ULTepAddr = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(anotherErabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_ADDR"
	anotherDbKeyBearerS1ULTepTeid = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(anotherErabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_TEID"
	anotherDbKeyBearerArpPL = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(anotherErabID) + ",UE_ERAB_QOS_ARP_PL"
	anotherDbKeyBearerQci = someUeID.GNb + "," + someUeID.ENbUeX2ApID + "," +
		fmt.Sprint(anotherErabID) + ",UE_ERAB_QOS_QCI"
}

func getTestCellEntry(pci uint32, ssbFreq uint32) *uenib.Cell {
	return &uenib.Cell{
		Pci:     pci,
		SsbFreq: ssbFreq,
	}
}

func getTestStateEntry(event string, cause string) *uenib.UeState {
	return &uenib.UeState{
		Event: event,
		Cause: cause,
	}
}

func getTestStateEntryOnlyEvent(event string) *uenib.UeState {
	return &uenib.UeState{
		Event: event,
	}
}

func getTestErabIDs(erabID1 uenib.ErabID, erabID2 uenib.ErabID) []uenib.ErabID {
	return []uenib.ErabID{erabID1, erabID2}
}

func getTestErabs() []uenib.Bearer {
	return []uenib.Bearer{
		uenib.Bearer{
			ErabID: 1000,
			DrbID:  150,
			ArpPL:  1,
			Qci:    10,
			S1ULGtpTE: uenib.TunnelEndpoint{
				Address: []byte("10.20.30.40"),
				Teid:    []byte("1999"),
			},
		},
		uenib.Bearer{
			ErabID: 2000,
			DrbID:  250,
			ArpPL:  2,
			Qci:    20,
			S1ULGtpTE: uenib.TunnelEndpoint{
				Address: []byte("20.20.30.40"),
				Teid:    []byte("2999"),
			},
		},
	}
}

func expectDbError(t *testing.T, err error, expCause string) {
	assert.NotNil(t, err)
	uenibError, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, uenibError.Temporary())
	assert.Equal(t, true, uenibreader.IsBackendError(err))
	assert.Equal(t, false, uenibreader.IsValueNotFoundFailure(err))
	assert.Equal(t, false, uenibreader.IsValidationError(err))
	assert.Equal(t, false, uenibreader.IsInternalError(err))
	assert.Contains(t, err.Error(), "database backend error: "+expCause)
}

func expectValueNotFoundFailure(t *testing.T, err error, name string) {
	assert.NotNil(t, err)
	uenibError, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, uenibError.Temporary())
	assert.Equal(t, true, uenibreader.IsValueNotFoundFailure(err))
	assert.Equal(t, false, uenibreader.IsBackendError(err))
	assert.Equal(t, false, uenibreader.IsValidationError(err))
	assert.Equal(t, false, uenibreader.IsInternalError(err))
	assert.Contains(t, err.Error(), "value of DB key '"+name+"' not found")
	assert.Contains(t, err.Error(), "not found")
}

func expectValidationError(t *testing.T, err error, name string) {
	assert.NotNil(t, err)
	uenibError, ok := err.(uenibreader.Error)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, uenibError.Temporary())
	assert.Equal(t, true, uenibreader.IsValidationError(err))
	assert.Equal(t, false, uenibreader.IsValueNotFoundFailure(err))
	assert.Equal(t, false, uenibreader.IsBackendError(err))
	assert.Equal(t, false, uenibreader.IsInternalError(err))
	assert.Contains(t, err.Error(), "validation error:")
	assert.Contains(t, err.Error(), name)
}

func getTestTunnelEndpointEntry(addr string, teid string) *uenib.TunnelEndpoint {
	return &uenib.TunnelEndpoint{
		Address: []byte(addr),
		Teid:    []byte(teid),
	}
}

func TestGetMeNbUEX2APIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()

	ret, err := i.GetMeNbUEX2APID(&anotherUeID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(100), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetMeNbUEX2APID(
		&uenib.UeID{
			GNbUeX2ApID: "200",
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Equal(t, uint32(0), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfNoGnbX2ApIDInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetMeNbUEX2APID(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "GNbUeX2ApID")
	assert.Equal(t, uint32(0), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfDbKeyNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{}, nil,
	).Once()

	ret, err := i.GetMeNbUEX2APID(&anotherUeID)

	expectValueNotFoundFailure(t, err, someDbKeyENbUeX2ApID)
	assert.Equal(t, uint32(0), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: nil}, nil,
	).Once()

	ret, err := i.GetMeNbUEX2APID(&anotherUeID)

	expectValueNotFoundFailure(t, err, someDbKeyENbUeX2ApID)
	assert.Equal(t, uint32(0), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetMeNbUEX2APID(&anotherUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetMeNbUEX2APIDReturnsErrorIfValueConvertToUint32Fails(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "IamNotInt"}, nil,
	).Once()

	ret, err := i.GetMeNbUEX2APID(&anotherUeID)

	expectValidationError(t, err, "IamNotInt")
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyGNbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyGNbUeX2ApID: "200"}, nil,
	).Once()

	ret, err := i.GetSgNbUEX2APID(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(200), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetSgNbUEX2APID(
		&uenib.UeID{
			GNbUeX2ApID: "200",
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfNoEnbX2ApIDInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetSgNbUEX2APID(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "ENbUeX2ApID")
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfDbKeyNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyGNbUeX2ApID}).Return(
		map[string]interface{}{}, nil,
	).Once()

	ret, err := i.GetSgNbUEX2APID(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyGNbUeX2ApID)
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyGNbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyGNbUeX2ApID: nil}, nil,
	).Once()

	ret, err := i.GetSgNbUEX2APID(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyGNbUeX2ApID)
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyGNbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetSgNbUEX2APID(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDReturnsErrorIfValueConvertToUint32Fails(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyGNbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyGNbUeX2ApID: "IamNotInt"}, nil,
	).Once()

	ret, err := i.GetSgNbUEX2APID(&someUeID)

	expectValidationError(t, err, "IamNotInt")
	assert.Equal(t, uint32(0), ret)
}

func TestGetPsCellSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyPsCellPci, someDbKeyPsCellSsbFreq}).Return(
		map[string]interface{}{someDbKeyPsCellPci: "10", someDbKeyPsCellSsbFreq: "20"}, nil,
	).Once()

	ret, err := i.GetPsCell(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestCellEntry(10, 20), ret)
}

func TestGetPsCellWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{someDbKeyPsCellPci, someDbKeyPsCellSsbFreq}).Return(
		map[string]interface{}{someDbKeyPsCellPci: "10", someDbKeyPsCellSsbFreq: "20"}, nil,
	).Once()

	ret, err := i.GetPsCell(&anotherUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestCellEntry(10, 20), ret)
}

func TestGetPsCellReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetPsCell(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetPsCellReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetPsCell(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetPsCellReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyPsCellPci, someDbKeyPsCellSsbFreq}).Return(nil, dbError).Once()

	ret, err := i.GetPsCell(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetPsCellWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetPsCell(&anotherUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetPsCellReturnsErrorIfPciDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyPsCellPci, someDbKeyPsCellSsbFreq}).Return(
		map[string]interface{}{someDbKeyPsCellPci: nil, someDbKeyPsCellSsbFreq: "20"}, nil,
	).Once()

	ret, err := i.GetPsCell(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyPsCellPci)
	assert.Nil(t, ret)
}

func TestGetPsCellReturnsErrorIfSsbFreqDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyPsCellPci, someDbKeyPsCellSsbFreq}).Return(
		map[string]interface{}{someDbKeyPsCellPci: "10", someDbKeyPsCellSsbFreq: nil}, nil,
	).Once()

	ret, err := i.GetPsCell(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyPsCellSsbFreq)
	assert.Nil(t, ret)
}

func TestGetStateSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyUeStateEvent, someDbKeyUeStateCause}).Return(
		map[string]interface{}{
			someDbKeyUeStateEvent: "2020-04-30T09:02:39.364571+03:00;SGNB-RECONF-CMPLT",
			someDbKeyUeStateCause: nil,
		}, nil,
	).Once()

	ret, err := i.GetState(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestStateEntryOnlyEvent("2020-04-30T09:02:39.364571+03:00;SGNB-RECONF-CMPLT"), ret)
}

func TestGetStateWithCauseSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyUeStateEvent, someDbKeyUeStateCause}).Return(
		map[string]interface{}{
			someDbKeyUeStateEvent: "2020-04-30T09:02:39.364571+03:00;SGNB-ADD-REQ-REJ",
			someDbKeyUeStateCause: "SGNB-ADD-REQ-REJ;radioNetwork;no_radio_resources_available",
		}, nil,
	).Once()

	ret, err := i.GetState(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestStateEntry(
		"2020-04-30T09:02:39.364571+03:00;SGNB-ADD-REQ-REJ",
		"SGNB-ADD-REQ-REJ;radioNetwork;no_radio_resources_available",
	), ret)
}

func TestGeStateWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{someDbKeyUeStateEvent, someDbKeyUeStateCause}).Return(
		map[string]interface{}{
			someDbKeyUeStateEvent: "STATE-123",
			someDbKeyUeStateCause: nil,
		}, nil,
	).Once()

	ret, err := i.GetState(&anotherUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestStateEntryOnlyEvent("STATE-123"), ret)
}

func TestGetStateReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetState(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetStateReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetState(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetStateReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyUeStateEvent, someDbKeyUeStateCause}).Return(nil, dbError).Once()

	ret, err := i.GetState(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetStateWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetState(&anotherUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetStateReturnsErrorIfEventDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyUeStateEvent, someDbKeyUeStateCause}).Return(
		map[string]interface{}{
			someDbKeyUeStateEvent: nil,
			someDbKeyUeStateCause: "something",
		}, nil,
	).Once()

	ret, err := i.GetState(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyUeStateEvent)
	assert.Nil(t, ret)
}

func TestGetBearerIDsSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000,2000"}, nil,
	).Once()

	ret, err := i.GetBearerIDs(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestErabIDs(1000, 2000), ret)
}

func TestGetBearerIDsWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000,2000"}, nil,
	).Once()

	ret, err := i.GetBearerIDs(&anotherUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestErabIDs(1000, 2000), ret)
}

func TestGetBearerIDsReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetBearerIDs(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetBearerIDsReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetBearerIDs(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetBearerIDsReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyBearerIDs}).Return(nil, dbError).Once()

	ret, err := i.GetBearerIDs(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetBearerIDsWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetBearerIDs(&anotherUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetBearerIDsReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: nil}, nil,
	).Once()

	ret, err := i.GetBearerIDs(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerIDs)
	assert.Nil(t, ret)
}

func TestGetBearerIDsReturnsErrorIfValueConvertToUint32Fails(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "IamNotInt"}, nil,
	).Once()

	ret, err := i.GetBearerIDs(&someUeID)

	expectValidationError(t, err, "IamNotInt")
	assert.Nil(t, ret)
}

func TestGetBearersSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000,2000"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerDrbID,
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
		someDbKeyBearerArpPL,
		someDbKeyBearerQci,
		anotherDbKeyBearerDrbID,
		anotherDbKeyBearerS1ULTepAddr,
		anotherDbKeyBearerS1ULTepTeid,
		anotherDbKeyBearerArpPL,
		anotherDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerDrbID:          "150",
			someDbKeyBearerS1ULTepAddr:    "10.20.30.40",
			someDbKeyBearerS1ULTepTeid:    "1999",
			someDbKeyBearerArpPL:          "1",
			someDbKeyBearerQci:            "10",
			anotherDbKeyBearerDrbID:       "250",
			anotherDbKeyBearerS1ULTepAddr: "20.20.30.40",
			anotherDbKeyBearerS1ULTepTeid: "2999",
			anotherDbKeyBearerArpPL:       "2",
			anotherDbKeyBearerQci:         "20",
		}, nil).Once()

	ret, err := i.GetBearers(&someUeID)

	assert.Nil(t, err)
	assert.Equal(t, getTestErabs(), ret)
}

func TestGetBearersReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetBearers(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetBearers(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")

	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerDrbID,
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
		someDbKeyBearerArpPL,
		someDbKeyBearerQci,
	}).Return(nil, dbError).Once()

	ret, err := i.GetBearers(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfBearerIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyBearerIDs}).Return(nil, dbError).Once()

	ret, err := i.GetBearers(&someUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetBearersWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetBearers(&anotherUeID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerDrbID,
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
		someDbKeyBearerArpPL,
		someDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerDrbID: nil,
		}, nil).Once()

	ret, err := i.GetBearers(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerDrbID)
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfBearerIDDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: nil}, nil,
	).Once()

	ret, err := i.GetBearers(&someUeID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerIDs)
	assert.Nil(t, ret)
}

func TestGetBearersReturnsErrorIfValueConvertToUint32Fails(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyBearerIDs}).Return(
		map[string]interface{}{someDbKeyBearerIDs: "1000"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerDrbID,
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
		someDbKeyBearerArpPL,
		someDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerDrbID: "IamNotInt",
		}, nil).Once()

	ret, err := i.GetBearers(&someUeID)

	expectValidationError(t, err, "IamNotInt")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTESuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: "10.20.30.40",
			someDbKeyBearerS1ULTepTeid: "1999",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, getTestTunnelEndpointEntry("10.20.30.40", "1999"), ret)
}

func TestGetErabS1ULGtpTEWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: "10.20.30.40",
			someDbKeyBearerS1ULTepTeid: "1999",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, getTestTunnelEndpointEntry("10.20.30.40", "1999"), ret)
}

func TestGetErabS1ULGtpTEReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTE(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
		someErabID,
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTE(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
		someErabID,
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
	}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTE(&anotherUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEReturnsErrorIfAddressDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: nil,
			someDbKeyBearerS1ULTepTeid: "1999",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerS1ULTepAddr)
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEReturnsErrorIfTeidDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: "10.20.30.40",
			someDbKeyBearerS1ULTepTeid: nil,
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerS1ULTepTeid)
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEAddrSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: "10.20.30.40",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTEAddr(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, []byte("10.20.30.40"), ret)
}

func TestGetErabS1ULGtpTEAddrWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: "10.20.30.40",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTEAddr(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, []byte("10.20.30.40"), ret)
}

func TestGetErabS1ULGtpTEAddrReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTEAddr(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
		someErabID,
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEAddrReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTEAddr(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
		someErabID,
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEAddrReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
	}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTEAddr(&someUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEAddrWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTEAddr(&anotherUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTEAddrReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepAddr,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepAddr: nil,
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTEAddr(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerS1ULTepAddr)
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTETeidSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepTeid: "1999",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTETeid(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, []byte("1999"), ret)
}

func TestGetErabS1ULGtpTETeidWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepTeid: "1999",
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTETeid(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, []byte("1999"), ret)
}

func TestGetErabS1ULGtpTETeidReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTETeid(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
		someErabID,
	)

	expectValidationError(t, err, "GNb")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTETeidReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabS1ULGtpTETeid(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
		someErabID,
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTETeidReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{
		someDbKeyBearerS1ULTepTeid,
	}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTETeid(&someUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTETeidWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetErabS1ULGtpTETeid(&anotherUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Nil(t, ret)
}

func TestGetErabS1ULGtpTETeidReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerS1ULTepTeid,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerS1ULTepTeid: nil,
		}, nil).Once()

	ret, err := i.GetErabS1ULGtpTETeid(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerS1ULTepTeid)
	assert.Nil(t, ret)
}

func TestGetErabQosArpPLSuccessSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerArpPL,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerArpPL: "1",
		}, nil).Once()

	ret, err := i.GetErabQosArpPL(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(1), ret)
}

func TestGetErabQosArpPLWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerArpPL,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerArpPL: "1",
		}, nil).Once()

	ret, err := i.GetErabQosArpPL(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(1), ret)
}

func TestGetErabQosArpPLReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabQosArpPL(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
		someErabID,
	)

	expectValidationError(t, err, "GNb")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosArpPLReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabQosArpPL(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
		someErabID,
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosArpPLReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{
		someDbKeyBearerArpPL,
	}).Return(nil, dbError).Once()

	ret, err := i.GetErabQosArpPL(&someUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosArpPLWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetErabQosArpPL(&anotherUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosArpPLReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerArpPL,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerArpPL: nil,
		}, nil).Once()

	ret, err := i.GetErabQosArpPL(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerArpPL)
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosQciSuccessSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerQci: "10",
		}, nil).Once()

	ret, err := i.GetErabQosQci(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(10), ret)
}

func TestGetErabQosQciWithGnbX2ApIDSuccess(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(
		map[string]interface{}{someDbKeyENbUeX2ApID: "100"}, nil,
	).Once()
	m.On("Get", []string{
		someDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerQci: "10",
		}, nil).Once()

	ret, err := i.GetErabQosQci(&someUeID, someErabID)

	assert.Nil(t, err)
	assert.Equal(t, uint32(10), ret)
}

func TestGetErabQosQciReturnsErrorIfNoGNbInUeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabQosQci(
		&uenib.UeID{
			ENbUeX2ApID: "100",
		},
		someErabID,
	)

	expectValidationError(t, err, "GNb")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosQciReturnsErrorIfNoUeX2ApID2UeID(t *testing.T) {
	_, i := setup()

	ret, err := i.GetErabQosQci(
		&uenib.UeID{
			GNb: "somegnb:310-410-b5c67788",
		},
		someErabID,
	)

	expectValidationError(t, err, "both UeX2ApIDs")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosQciReturnsErrorIfDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{
		someDbKeyBearerQci,
	}).Return(nil, dbError).Once()

	ret, err := i.GetErabQosQci(&someUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosQciWithGnbX2ApIDReturnsErrorIfEnbX2ApIDDbQueryFails(t *testing.T) {
	m, i := setup()
	dbError := errors.New("Some DB Error")
	m.On("Get", []string{someDbKeyENbUeX2ApID}).Return(nil, dbError).Once()

	ret, err := i.GetErabQosQci(&anotherUeID, someErabID)

	expectDbError(t, err, "Some DB Error")
	assert.Equal(t, uint32(0), ret)
}

func TestGetErabQosQciReturnsErrorIfDbKeyValueNotFound(t *testing.T) {
	m, i := setup()
	m.On("Get", []string{
		someDbKeyBearerQci,
	}).Return(
		map[string]interface{}{
			someDbKeyBearerQci: nil,
		}, nil).Once()

	ret, err := i.GetErabQosQci(&someUeID, someErabID)

	expectValueNotFoundFailure(t, err, someDbKeyBearerQci)
	assert.Equal(t, uint32(0), ret)
}

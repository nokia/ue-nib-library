/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader_test

import (
	"github.com/nokia/ue-nib-library/pkg/uenib"
	"github.com/stretchr/testify/assert"
	"testing"
)

var someUeID uenib.UeID
var someErabID uenib.ErabID

func init() {
	someUeID = uenib.UeID{
		GNb:         "somegnb:310-410-b5c67788",
		GNbUeX2ApID: "200",
		ENbUeX2ApID: "100",
	}
	someErabID = 1000
}

func TestGetMeNbUEX2APIDSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetMeNbUEX2APID(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), ret)
}

func TestGetSgNbUEX2APIDSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetSgNbUEX2APID(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), ret)
}

func TestGetPsCellSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetPsCell(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, uenib.Cell{}, *ret)
}

func TestGetUeStateSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetState(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, uenib.UeState{}, *ret)
}

func TestGetBearersSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetBearers(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, make([]uenib.Bearer, 0), ret)
}

func TestGetBearerIDsSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetBearerIDs(&someUeID)
	assert.Nil(t, err)
	assert.Equal(t, make([]uenib.ErabID, 0), ret)
}

func TestGetS1ULGtpTESuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetErabS1ULGtpTE(&someUeID, someErabID)
	assert.Nil(t, err)
	assert.Equal(t, uenib.TunnelEndpoint{}, *ret)
}

func TestGetS1ULGtpTEAddrSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetErabS1ULGtpTEAddr(&someUeID, someErabID)
	assert.Nil(t, err)
	assert.Equal(t, make([]byte, 0), ret)
}

func TestGetS1ULGtpTETeidSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetErabS1ULGtpTETeid(&someUeID, someErabID)
	assert.Nil(t, err)
	assert.Equal(t, make([]byte, 0), ret)
}

func TestGetQosQciSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetErabQosQci(&someUeID, someErabID)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), ret)
}

func TestGetQosArpPLSuccess(t *testing.T) {
	_, i := setup()
	ret, err := i.GetErabQosArpPL(&someUeID, someErabID)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), ret)
}

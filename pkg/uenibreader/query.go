/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import (
	"github.com/nokia/ue-nib-library/pkg/uenib"
)

//GetMeNbUEX2APID returns UE MeNbUEX2APID.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetMeNbUEX2APID(ueID *uenib.UeID) (uint32, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return uint32(0), err
}

//GetSgNbUEX2APID returns UE SgNbUEX2APID.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetSgNbUEX2APID(ueID *uenib.UeID) (uint32, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return uint32(0), err
}

//GetPsCell returns UE radio resource information container, called as a Primary
//Cell in secondary Node (PSCell).
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetPsCell(ueID *uenib.UeID) (*uenib.Cell, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return &uenib.Cell{}, err
}

//GetState returns UE's last known state in UE-NIB
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetState(ueID *uenib.UeID) (*uenib.UeState, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return &uenib.UeState{}, err
}

//GetBearers returns active bearers (E-RABs) of an UE.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetBearers(ueID *uenib.UeID) ([]uenib.Bearer, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return make([]uenib.Bearer, 0), err
}

//GetBearerIDs returns existing bearer identifiers (E-RAB IDs) of an UE.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetBearerIDs(ueID *uenib.UeID) ([]uenib.ErabID, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return make([]uenib.ErabID, 0), err
}

//GetErabS1ULGtpTE returns UE bearer's S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTE(ueID *uenib.UeID, erabID uenib.ErabID) (*uenib.TunnelEndpoint, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return &uenib.TunnelEndpoint{}, err
}

//GetErabS1ULGtpTEAddr returns UE bearer's transport layer address of an S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTEAddr(ueID *uenib.UeID, erabID uenib.ErabID) ([]byte, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return make([]byte, 0), err
}

//GetErabS1ULGtpTETeid returns UE bearer's GTP tunnel endpoint ID (TEID) of an S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTETeid(ueID *uenib.UeID, erabID uenib.ErabID) ([]byte, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return make([]byte, 0), err
}

//GetErabQosQci returns UE bearer's QoS Class Identifier (QCI).
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabQosQci(ueID *uenib.UeID, erabID uenib.ErabID) (uint32, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return uint32(0), err
}

//GetErabQosArpPL returns UE bearer's QoS Allocation and Retention Priority (ARP) priority level (PL)
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabQosArpPL(ueID *uenib.UeID, erabID uenib.ErabID) (uint32, error) {
	var err error
	//@todo Add value reading from UE-NIB database.
	return uint32(0), err
}

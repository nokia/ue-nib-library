/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import (
	"errors"
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenib"
	"github.com/nokia/ue-nib-library/pkg/uenibreader/internal"
	"strconv"
	"strings"
)

//GetMeNbUEX2APID returns UE MeNbUEX2APID.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetMeNbUEX2APID(ueID *uenib.UeID) (uint32, error) {
	if len(ueID.GNb) == 0 {
		return uint32(0), toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing GNb", ueID.String())))
	}

	if len(ueID.GNbUeX2ApID) == 0 {
		return uint32(0), toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing GNbUeX2ApID", ueID.String())))
	}

	key := internal.DbKeyUeMapGNbToENbUeX2ApID(ueID)

	q, err := reader.newGetQuery(ueID, []string{key})
	if err != nil {
		return uint32(0), err
	}

	return q.getKeyUint32Value(ueID, key)
}

//GetSgNbUEX2APID returns UE SgNbUEX2APID.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetSgNbUEX2APID(ueID *uenib.UeID) (uint32, error) {
	if len(ueID.GNb) == 0 {
		return uint32(0), toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing GNb", ueID.String())))
	}

	if len(ueID.ENbUeX2ApID) == 0 {
		return uint32(0), toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing ENbUeX2ApID", ueID.String())))
	}

	key := internal.DbKeyUeMapENbToGNbUeX2ApID(ueID)

	q, err := reader.newGetQuery(ueID, []string{key})
	if err != nil {
		return uint32(0), err
	}

	return q.getKeyUint32Value(ueID, key)
}

//GetPsCell returns UE radio resource information container, called as a Primary
//Cell in secondary Node (PSCell).
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetPsCell(ueID *uenib.UeID) (*uenib.Cell, error) {
	var q *query
	var retCell uenib.Cell

	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	pciKey := internal.DbKeyPsCellPci(id)
	freqKey := internal.DbKeyPsCellSsbFreq(id)

	if q, err = reader.newGetQuery(ueID, []string{pciKey, freqKey}); err != nil {
		return nil, err
	}

	if retCell.Pci, err = q.getKeyUint32Value(ueID, pciKey); err != nil {
		return nil, err
	}
	if retCell.SsbFreq, err = q.getKeyUint32Value(ueID, freqKey); err != nil {
		return nil, err
	}
	return &retCell, err
}

//GetState returns UE's last known state in UE-NIB and the last GTP Cause code if there has
//been any Cause IEs set in any UE's X2 messages.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetState(ueID *uenib.UeID) (*uenib.UeState, error) {
	var q *query
	var retState uenib.UeState

	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	stateEventKey := internal.DbKeyUeStateEvent(id)
	stateCauseKey := internal.DbKeyUeStateCause(id)

	if q, err = reader.newGetQuery(ueID, []string{stateEventKey, stateCauseKey}); err != nil {
		return nil, err
	}

	if retState.Event, err = q.getKeyStringValue(ueID, stateEventKey); err != nil {
		return nil, err
	}

	//These is no cause value in UE-NIB, if UE's all mobility procedures are done successfully.
	//That's why catch valueNotFoundFailure error and ignore it.
	if retState.Cause, err = q.getKeyStringValue(ueID, stateCauseKey); err != nil {
		if IsValueNotFoundFailure(err) {
			err = nil
		}
	}

	return &retState, err
}

//GetBearerIDs returns existing bearer identifiers (E-RAB IDs) of an UE.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetBearerIDs(ueID *uenib.UeID) ([]uenib.ErabID, error) {
	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	return reader.getErabIDs(id)
}

//GetBearers returns active bearers (E-RABs) of an UE.
//Parameter ueID identifies User equipment (UE).
func (reader *Reader) GetBearers(ueID *uenib.UeID) ([]uenib.Bearer, error) {
	var q *query
	var erabIDKeys []string
	var retBearers []uenib.Bearer

	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	erabIDs, err := reader.getErabIDs(id)
	if err != nil {
		return nil, err
	}

	for _, erabID := range erabIDs {
		erabIDKeys = append(erabIDKeys, internal.GetErabAllDbKeys(id, erabID)...)
	}

	if q, err = reader.newGetQuery(ueID, erabIDKeys); err != nil {
		return nil, err
	}

	for _, erabID := range erabIDs {
		var br uenib.Bearer
		br.ErabID = erabID
		if br.DrbID, err = q.getKeyUint32Value(id, internal.DbKeyErabDrbID(id, erabID)); err != nil {
			return nil, err
		}
		if br.ArpPL, err = q.getKeyUint32Value(id, internal.DbKeyErabQosArpPL(id, erabID)); err != nil {
			return nil, err
		}
		if br.Qci, err = q.getKeyUint32Value(id, internal.DbKeyErabQosQci(id, erabID)); err != nil {
			return nil, err
		}
		if br.S1ULGtpTE.Address, err = q.getKeyByteSliceValue(id, internal.DbKeyErabS1UlGtpTendpAddr(id, erabID)); err != nil {
			return nil, err
		}
		if br.S1ULGtpTE.Teid, err = q.getKeyByteSliceValue(id, internal.DbKeyErabS1UlGtpTendpTeid(id, erabID)); err != nil {
			return nil, err
		}
		retBearers = append(retBearers, br)
	}
	return retBearers, err
}

//GetErabS1ULGtpTE returns UE bearer's S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTE(ueID *uenib.UeID, erabID uenib.ErabID) (*uenib.TunnelEndpoint, error) {
	var q *query
	var retTEp uenib.TunnelEndpoint

	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	addrKey := internal.DbKeyErabS1UlGtpTendpAddr(id, erabID)
	teidKey := internal.DbKeyErabS1UlGtpTendpTeid(id, erabID)

	if q, err = reader.newGetQuery(ueID, []string{addrKey, teidKey}); err != nil {
		return nil, err
	}

	if retTEp.Address, err = q.getKeyByteSliceValue(ueID, addrKey); err != nil {
		return nil, err
	}
	if retTEp.Teid, err = q.getKeyByteSliceValue(ueID, teidKey); err != nil {
		return nil, err
	}
	return &retTEp, err
}

//GetErabS1ULGtpTEAddr returns UE bearer's transport layer address of an S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTEAddr(ueID *uenib.UeID, erabID uenib.ErabID) ([]byte, error) {
	var q *query
	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	addrKey := internal.DbKeyErabS1UlGtpTendpAddr(id, erabID)

	if q, err = reader.newGetQuery(ueID, []string{addrKey}); err != nil {
		return nil, err
	}

	return q.getKeyByteSliceValue(ueID, addrKey)
}

//GetErabS1ULGtpTETeid returns UE bearer's GTP tunnel endpoint ID (TEID) of an S1 uplink GTP tunnel endpoint.
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabS1ULGtpTETeid(ueID *uenib.UeID, erabID uenib.ErabID) ([]byte, error) {
	var q *query
	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return nil, err
	}

	teidKey := internal.DbKeyErabS1UlGtpTendpTeid(id, erabID)

	if q, err = reader.newGetQuery(ueID, []string{teidKey}); err != nil {
		return nil, err
	}

	return q.getKeyByteSliceValue(ueID, teidKey)
}

//GetErabQosArpPL returns UE bearer's QoS Allocation and Retention Priority (ARP) priority level (PL)
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabQosArpPL(ueID *uenib.UeID, erabID uenib.ErabID) (uint32, error) {
	var q *query
	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return uint32(0), err
	}

	arpPLKey := internal.DbKeyErabQosArpPL(id, erabID)

	if q, err = reader.newGetQuery(ueID, []string{arpPLKey}); err != nil {
		return uint32(0), err
	}

	return q.getKeyUint32Value(ueID, arpPLKey)
}

//GetErabQosQci returns UE bearer's QoS Class Identifier (QCI).
//Parameter ueID identifies User equipment (UE).
//Parameter erabID identifies bearer.
func (reader *Reader) GetErabQosQci(ueID *uenib.UeID, erabID uenib.ErabID) (uint32, error) {
	var q *query
	id, err := reader.validateUeIDAndResolveENbX2ApID(ueID)
	if err != nil {
		return uint32(0), err
	}

	qciKey := internal.DbKeyErabQosQci(id, erabID)

	if q, err = reader.newGetQuery(ueID, []string{qciKey}); err != nil {
		return uint32(0), err
	}

	return q.getKeyUint32Value(ueID, qciKey)
}

func (reader *Reader) validateUeIDAndResolveENbX2ApID(ueID *uenib.UeID) (*uenib.UeID, error) {
	var err error
	if err = validateUe(ueID); err != nil {
		return nil, err
	}
	//Make own copy of UeID not to alter the original ueID received in UE-NIB Reader API.
	id := *ueID

	if len(id.ENbUeX2ApID) == 0 {
		key := internal.DbKeyUeMapGNbToENbUeX2ApID(ueID)
		if id.ENbUeX2ApID, err = reader.resolveENbX2ApID(ueID, key); err != nil {
			return nil, err
		}
	}
	return &id, err
}

func (reader *Reader) resolveENbX2ApID(ueID *uenib.UeID, key string) (string, error) {
	q, err := reader.newGetQuery(ueID, []string{key})
	if err != nil {
		return "", err
	}

	return q.getKeyStringValue(ueID, key)
}

func (reader *Reader) getErabIDs(ueID *uenib.UeID) ([]uenib.ErabID, error) {
	var strVal string

	key := internal.DbKeyUeErabIDs(ueID)

	q, err := reader.newGetQuery(ueID, []string{key})
	if err != nil {
		return nil, err
	}

	if strVal, err = q.getKeyStringValue(ueID, key); err != nil {
		return nil, err
	}
	return parseErabIDsStringToErabIDSlice(ueID, strVal)
}

func (reader *Reader) newGetQuery(ueID *uenib.UeID, keys []string) (*query, error) {
	var err error
	q := &query{}
	if q.kvMap, err = reader.db.Get(keys); err != nil {
		return nil, toBackendError(ueID, err)
	}
	return q, err
}

type query struct {
	kvMap map[string]interface{}
}

func (q *query) getKeyStringValue(ueID *uenib.UeID, key string) (string, error) {
	var err error
	if val, ok := q.kvMap[key]; ok {
		if val == nil {
			return "", toValueNotFoundFailure(ueID, key)
		}
		return val.(string), err
	}
	return "", toValueNotFoundFailure(ueID, key)
}

func (q *query) getKeyUint32Value(ueID *uenib.UeID, key string) (uint32, error) {
	strVal, err := q.getKeyStringValue(ueID, key)
	if err != nil {
		return uint32(0), err
	}
	return parseStringToUint32(ueID, strVal)
}

func (q *query) getKeyByteSliceValue(ueID *uenib.UeID, key string) ([]byte, error) {
	strVal, err := q.getKeyStringValue(ueID, key)
	if err != nil {
		return nil, err
	}
	return []byte(strVal), err
}

func toValueNotFoundFailure(ueID *uenib.UeID, key string) *valueNotFoundFailure {
	return &valueNotFoundFailure{ueID: *ueID, name: key, temporary: true}
}

func toBackendError(ueID *uenib.UeID, err error) *backendError {
	return &backendError{ueID: *ueID, err: err.Error(), temporary: true}
}

func toValidationError(ueID *uenib.UeID, err error) *validationError {
	return &validationError{ueID: *ueID, err: err.Error()}
}

func validateUe(ueID *uenib.UeID) error {
	var err error
	if len(ueID.GNb) == 0 {
		err = toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing GNb", ueID.String())))
		return err
	}

	if len(ueID.GNbUeX2ApID) == 0 && len(ueID.ENbUeX2ApID) == 0 {
		err = toValidationError(ueID, errors.New(fmt.Sprintf("%s :: missing both UeX2ApIDs", ueID.String())))
		return err
	}
	return err
}

func parseStringToUint32(ueID *uenib.UeID, str string) (uint32, error) {
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return uint32(0), toValidationError(ueID, err)
	}
	return uint32(val), err
}

func parseErabIDsStringToErabIDSlice(ueID *uenib.UeID, strList string) ([]uenib.ErabID, error) {
	var err error
	var uIntVal uint32
	strVals := strings.Split(strList, ",")
	erabVals := make([]uenib.ErabID, len(strVals))

	for i, s := range strVals {
		if uIntVal, err = parseStringToUint32(ueID, s); err != nil {
			return nil, err
		}
		erabVals[i] = uenib.ErabID(uIntVal)
	}
	return erabVals, err
}

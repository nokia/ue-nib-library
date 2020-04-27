/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

//Package uenib provides UE-NIB public data types.
package uenib

import (
	"fmt"
)

//UEID struct is a holder for a User equipment (UE) identifier in UE-NIB.
//ENb is not in use at the moment, because RIC is unaware of eNBs.
//GNbUeX2ApID and ENbUeX2ApID are optional but either one or both must be set.
type UeID struct {
	GNb         string //Mandatory. Contains GNb RanName form of: <Antenna-Type>:<3 MCC digits>-<3 MNC digits>-<Node ID>
	ENb         string //Not used at the moment, because RIC is unaware of eNBs.
	GNbUeX2ApID string //Optional. Either GNbUeX2ApID or ENbUeX2ApID must be set.
	ENbUeX2ApID string //Optional. Either ENbUeX2ApID or GNbUeX2ApID must be set.
}

//ErabID type is a type used to identify a bearer (E-RAB) of an UE.
type ErabID uint32

//Bearer type is a holder for a User equipment (UE) Bearer level information.
type Bearer struct {
	ErabID    ErabID
	DrbID     uint32
	ArpPL     uint32
	Qci       uint32
	S1ULGtpTE TunnelEndpoint
}

//TunnelEndpoint is a holder for a GTP tunnel endpoint.
type TunnelEndpoint struct {
	Address []byte
	Teid    []byte
}

//Cell type is a holder for a User equipment (UE) Radio resource information.
type Cell struct {
	Pci     uint32 //Physical cell ID
	SsbFreq uint32 //Frequency of the SSB to be used for the serving cell.
}

//UeState is a holder for a UE state.
type UeState struct {
	Event string //A string of a timestamp and UE's last X2 message what UE-NIB has detected.
	Cause string //X2 message's Cause IE value what UE-NIB has lastly detected.
}

//Helper function to print UeID.
func (ueID UeID) String() string {
	return fmt.Sprintf("UeID:[GNb:%s,ENb:%s,GNbUeX2ApID:%s,ENbUeX2ApID:%s]",
		ueID.GNb, ueID.ENb, ueID.GNbUeX2ApID, ueID.ENbUeX2ApID)
}

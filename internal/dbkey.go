/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package internal

import (
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenib"
)

func DbKeyUeMapGNbToENbUeX2ApID(ueID *uenib.UeID) string {
	return ueID.GNbUeX2ApID + ",UEMAP_ENBUEX2APID"
}

func DbKeyUeMapENbToGNbUeX2ApID(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UEMAP_GNBUEX2APID"
}

func DbKeyUeStateEvent(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UE_STATE_EVENT"
}

func DbKeyUeStateCause(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UE_STATE_CAUSE"
}

func DbKeyPsCellPci(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UE_PSCELL_PCI"
}

func DbKeyPsCellSsbFreq(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UE_PSCELL_FREQ"
}

func DbKeyUeErabIDs(ueID *uenib.UeID) string {
	return ueID.ENbUeX2ApID + ",UE_ERAB_IDS"
}

func DbKeyErabDrbID(ueID *uenib.UeID, erabID uenib.ErabID) string {
	return ueID.ENbUeX2ApID + "," + fmt.Sprint(erabID) + ",UE_ERAB_DRB_ID"
}

func DbKeyErabS1UlGtpTendpAddr(ueID *uenib.UeID, erabID uenib.ErabID) string {
	return ueID.ENbUeX2ApID + "," + fmt.Sprint(erabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_ADDR"
}

func DbKeyErabS1UlGtpTendpTeid(ueID *uenib.UeID, erabID uenib.ErabID) string {
	return ueID.ENbUeX2ApID + "," + fmt.Sprint(erabID) + ",UE_ERAB_S1_UL_GTP_TUNNEL_TEID"
}

func DbKeyErabQosArpPL(ueID *uenib.UeID, erabID uenib.ErabID) string {
	return ueID.ENbUeX2ApID + "," + fmt.Sprint(erabID) + ",UE_ERAB_QOS_ARP_PL"
}

func DbKeyErabQosQci(ueID *uenib.UeID, erabID uenib.ErabID) string {
	return ueID.ENbUeX2ApID + "," + fmt.Sprint(erabID) + ",UE_ERAB_QOS_QCI"
}

func GetErabAllDbKeys(ueID *uenib.UeID, erabID uenib.ErabID) []string {
	return []string{
		DbKeyErabDrbID(ueID, erabID),
		DbKeyErabS1UlGtpTendpAddr(ueID, erabID),
		DbKeyErabS1UlGtpTendpTeid(ueID, erabID),
		DbKeyErabQosArpPL(ueID, erabID),
		DbKeyErabQosQci(ueID, erabID),
	}
}

func GetUeNibNs(gNb string) string {
	return "uenib/" + gNb
}

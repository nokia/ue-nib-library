/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import (
	"fmt"
	"github.com/nokia/ue-nib-library/internal"
	"github.com/nokia/ue-nib-library/pkg/uenib"
	"strconv"
	"strings"
)

//EventCategory groups related events. Events can be subscribed per event category.
type EventCategory int

const (
	//DualConnectivity events are triggered after dual connectivity related data has been updated.
	//
	//Following events are possible in this category:
	//    <UE_ID>_ADD
	//        -UE ENDC established.
	//    <UE_ID>_REMOVE
	//        -UE ENDC released.
	//    <UE_ID>_<S1UL_TUN_ENDPOINT>_S1UL_TUNNEL_ESTABLISH
	//        -S1 uplink tunnel established.
	//    <UE_ID>_<S1UL_TUN_ENDPOINT>_S1UL_TUNNEL_RELEASE
	//        -S1 uplink tunnel released.
	//    GNB_ALL_UES_REMOVE
	//        All UEs within the gNB were removed.
	//
	//<UE_ID> identifies a UE in question. It consists of three sub-fields separated by
	//hashtag '#':
	//<GNb>#<GNbUeX2ApID>#<ENbUeX2ApID>. Note that GNb is form of a RanName:
	//<Antenna-Type>:<3 MCC digits>-<3 MNC digits>-<Node ID>.
	//<S1UL_TUN_ENDPOINT> identifies bearer's S1 uplink GTP tunnel endpoint. It consists of two
	//sub-fields separated by hashtag '#':
	//<Transport address>#<GTP TEID>. IP address and TEID are strings. IP address string can
	//contain an IPv4 dotted decimal ("192.0.2.1"), an IPv6 ("2001:db8::68") or a dual address
	//where IPv4 and IPv6 addresses are separated by '+' character.
	//Note that multiple S1 uplink GTP tunnel endpoints can be notified by a single event.
	//Multiple S1 uplink GTP tunnel endpoints are separated by hashtag '#' in an event string.
	DualConnectivity EventCategory = iota
)

//DcEventType defines possible dual connectivity event types.
type DcEventType int

const (
	DC_EVENT_UNKNOWN DcEventType = iota
	DC_EVENT_ADD
	DC_EVENT_REMOVE
	DC_EVENT_S1UL_TUNNEL_ESTABLISH
	DC_EVENT_S1UL_TUNNEL_RELEASE
	DC_EVENT_GNB_ALL_UES_REMOVE
)

//DcEvent defines all the entities what can be parsed from received dual connectivity event.
type DcEvent struct {
	EventType      DcEventType
	UeID           uenib.UeID
	S1ULGtpTunnels []DcEventTunnel
}

//DcEventTunnel defines tunnel endpoint address and tunnel identifier, which can
//be parsed from received event.
//Note that an IP address in the 'Addr' field is in dotted-decimal format:
//IPV4 ("192.0.2.1"), IPv6 ("2001:db8::68"). In dual address case both IPv4 and
//IPv6 addresses are set to 'Addr' field separated by '+' character.
//Field 'Teid' is Tunnel Endpoint ID (TEID) in host byte order.
type DcEventTunnel struct {
	Addr string
	Teid uint32
}

//String returns dual connectivity event type as a string.
func (dcEvt DcEventType) String() string {
	evtStrMap := [...]string{
		"UNKNOWN",
		"_ADD",
		"_REMOVE",
		"_S1UL_TUNNEL_ESTABLISH",
		"_S1UL_TUNNEL_RELEASE",
		"GNB_ALL_UES_REMOVE",
	}
	if dcEvt > DC_EVENT_GNB_ALL_UES_REMOVE {
		panic(fmt.Sprintf("DC event ID %d overflows name string array.\n", dcEvt))
	}
	return evtStrMap[dcEvt]
}

//String returns event category as a string.
func (category EventCategory) String() string {
	categories := [...]string{"DUAL_CONNECTIVITY"}
	if category < DualConnectivity || category > DualConnectivity {
		return "Unknown"
	}
	return categories[category]
}

//EventCallback defines the signature for the event callback function.
//Parameter gNb identifies GNb RanName what is form of: <Antenna-Type>:<3 MCC digits>-<3 MNC digits>-<Node ID>.
//Parameter eventCategory identifies event category.
type EventCallback func(gNb string, eventCategory EventCategory, events []string)

//SubscribeEvents is used to subscribe events from UE-NIB data changes
//(event publishing is done by uenibwriter along with the backend database modification).
//
//Events are subscribed per gNB and event category (at least one of each is required).
//
//Possible events related to each category are listed above in event category constant declaration.
//
//Same callback function may be given for several SubscribeEvents() calls. There might be several
//events combined to one callback function call.
//
//Event delivery protocol is not reliable (a reliable protocol is a protocol which verifies whether
//the delivery of data was successful). Published events should rarely be lost, but it is possible.
//
//In failure case SubscribeEvents() returns an error value indicating an abnormal state.
//In addition to Error() method defined in built-in error interface a function caller can test
//returned error value for a reader.Error with a type assertion and then distinguish temporal errors
//from permanent ones by using Temporary() method. In case of temporal error, the caller of
//SubscribeEvents() may retry the call after a short period of time.
func (reader *Reader) SubscribeEvents(gNbs []string, eventCategories []EventCategory, callback EventCallback) error {
	for gNbIndex := range gNbs {
		for eventCategoriesIndex := range eventCategories {
			if eventCategories[eventCategoriesIndex].String() == "Unknown" {
				return newValidationError("Unknown event category ID: %d", eventCategories[eventCategoriesIndex])
			}
			channel := gNbs[gNbIndex] + "_" + eventCategories[eventCategoriesIndex].String()
			ns := internal.GetUeNibNs(gNbs[gNbIndex])
			err := reader.db.SubscribeChannel(ns, reader.eventCallback(gNbs[gNbIndex], eventCategories[eventCategoriesIndex], callback), channel)
			if err != nil {
				return newBackendError(err.Error())
			}
		}
	}
	return nil
}

func (reader *Reader) eventCallback(gNb string, eventCategory EventCategory, clientCallback EventCallback) func(ch string, ev ...string) {
	return func(ch string, ev ...string) {
		clientCallback(gNb, eventCategory, ev)
	}
}

//ParseDcEvent parses an event string of the dual connectivity category and it returns
//two values: parsing results in a return value of 'DcEvent' type and status of parsing
//in a return value of standard 'error' type. Error status is returned, if parsing has
//been failed for some reason. If parsing has succeeded, parsed values from event string
//are returned inside 'DcEvent' type. Note that success status is also returned, when
//event type in parsed event string is unknown for the parser. In this case event type is
//set to DC_EVENT_UNKNOWN in returned 'DcEvent' type.
func ParseDcEvent(evtStr string) (DcEvent, error) {
	var err error
	var ret DcEvent
	evtFieldStr := parseDcEventType(evtStr, &ret)

	switch ret.EventType {
	case DC_EVENT_S1UL_TUNNEL_ESTABLISH:
		err = parseDcTunnelEvent(evtFieldStr, &ret)
	case DC_EVENT_S1UL_TUNNEL_RELEASE:
		err = parseDcTunnelEvent(evtFieldStr, &ret)
	case DC_EVENT_ADD:
		err = parseUeFromDcEvent(evtFieldStr, &ret)
	case DC_EVENT_REMOVE:
		err = parseUeFromDcEvent(evtFieldStr, &ret)
	}
	return ret, err
}

func parseDcEventType(evtStr string, ret *DcEvent) string {
	if matched := strings.HasSuffix(evtStr, DC_EVENT_GNB_ALL_UES_REMOVE.String()); matched {
		ret.EventType = DC_EVENT_GNB_ALL_UES_REMOVE
		return strings.TrimSuffix(evtStr, ret.EventType.String())
	}
	if matched := strings.HasSuffix(evtStr, DC_EVENT_S1UL_TUNNEL_ESTABLISH.String()); matched {
		ret.EventType = DC_EVENT_S1UL_TUNNEL_ESTABLISH
		return strings.TrimSuffix(evtStr, ret.EventType.String())
	}
	if matched := strings.HasSuffix(evtStr, DC_EVENT_S1UL_TUNNEL_RELEASE.String()); matched {
		ret.EventType = DC_EVENT_S1UL_TUNNEL_RELEASE
		return strings.TrimSuffix(evtStr, ret.EventType.String())
	}
	if matched := strings.HasSuffix(evtStr, DC_EVENT_ADD.String()); matched {
		ret.EventType = DC_EVENT_ADD
		return strings.TrimSuffix(evtStr, ret.EventType.String())
	}
	if matched := strings.HasSuffix(evtStr, DC_EVENT_REMOVE.String()); matched {
		ret.EventType = DC_EVENT_REMOVE
		return strings.TrimSuffix(evtStr, ret.EventType.String())
	}
	return evtStr
}

func parseDcTunnelEvent(evtStr string, ret *DcEvent) error {
	fields := strings.Split(evtStr, "_")
	if cnt := len(fields); cnt != 2 {
		return fmt.Errorf("Event '%s' parse failure: no UE ID or tunnel field in '%s'",
			ret.EventType.String(), evtStr)
	}
	if err := parseUeFromDcEvent(fields[0], ret); err != nil {
		return err
	}
	return parseTunnelsFromDcEvent(fields[1], ret)
}

func parseUeFromDcEvent(ueEvtStr string, ret *DcEvent) error {
	ueFields := strings.Split(ueEvtStr, "#")
	if cnt := len(ueFields); cnt != 3 {
		return fmt.Errorf("Event '%s' parse failure: wrong UE ID fields in '%s'",
			ret.EventType.String(), ueEvtStr)
	}
	ret.UeID.GNb = ueFields[0]
	ret.UeID.GNbUeX2ApID = ueFields[1]
	ret.UeID.ENbUeX2ApID = ueFields[2]
	return nil
}

func parseTunnelsFromDcEvent(tunEvtStr string, ret *DcEvent) error {
	tunFields := strings.Split(tunEvtStr, "#")
	if len(tunFields)%2 != 0 {
		return fmt.Errorf("Event '%s' parse failure: wrong tunnel fields in '%s'",
			ret.EventType.String(), tunEvtStr)
	}
	for i := 0; i < len(tunFields); i = i + 2 {
		var teid uint32
		if tunFields[i+1] != "" {
			u64Teid, err := strconv.ParseUint(tunFields[i+1], 10, 32)
			if err != nil {
				return fmt.Errorf("Event '%s' parse failure: wrong TEID field in '%s' conversion error:'%s'",
					ret.EventType.String(), tunEvtStr, err.Error())
			}
			teid = uint32(u64Teid)
		}
		t := DcEventTunnel{
			Addr: tunFields[i],
			Teid: teid,
		}
		ret.S1ULGtpTunnels = append(ret.S1ULGtpTunnels, t)
	}
	return nil
}

/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import ()

//EventCategory groups related events. Events can be subscribed per event category.
type EventCategory int

const (
	//DualConnectivity events are triggered after dual connectivity related data has been updated.
	//
	//Following events are possible in this category:
	//    <UE_ID>_ESTABLISH
	//        -UE ENDC established.
	//    <UE_ID>_RELEASE
	//        -UE ENDC released.
	//    GNB_ALL_UES_REMOVE
	//        All UEs within the gNB were removed.
	//Where <UE_ID> identifies an UE in question and it consists of three sub-fields:
	//<GNb>#<GNbUeX2ApID>#<ENbUeX2ApID>. Note that GNb is form of a RanName:
	//<Antenna-Type>:<3 MCC digits>-<3 MNC digits>-<Node ID>.
	DualConnectivity EventCategory = iota
)

func (category EventCategory) String() string {
	categories := [...]string{"DUAL_CONNECTIVITY"}
	if category < DualConnectivity || category > DualConnectivity {
		return "Unknown"
	}
	return categories[category]
}

//EventCallback defines the signature for the event callback function.
//Parameter gNb identifies GNb RanNamewhat is form of: <Antenna-Type>:<3 MCC digits>-<3 MNC digits>-<Node ID>.
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
			err := reader.db.SubscribeChannel(reader.eventCallback(gNbs[gNbIndex], eventCategories[eventCategoriesIndex], callback), channel)
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

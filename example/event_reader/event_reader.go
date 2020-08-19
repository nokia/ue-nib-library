/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package main

import (
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenib"
	"github.com/nokia/ue-nib-library/pkg/uenibreader"
	"sync"
	"time"
)

type logEvent struct {
	function string
	lines    []string
}

var someGNb string
var myReader *uenibreader.Reader
var ueDcReaderWaitGroup sync.WaitGroup
var ueDcEventHandlerChannel chan uenibreader.DcEvent
var loggerWaitGroup sync.WaitGroup
var logChannel chan logEvent

func init() {
	myReader = uenibreader.NewReader()
	//A channel is created per GNb 'someGNb'
	someGNb = "somegnb:310-410-b5c67788"
	//someGNb = "gnb-0"
	ueDcEventHandlerChannel = make(chan uenibreader.DcEvent)
	logChannel = make(chan logEvent)
}

func eventHandler(wg *sync.WaitGroup) {
	fmt.Printf("Reader Query Go routine starts\n")
	defer wg.Done()
	for {
		select {
		case evtInfo, ok := <-ueDcEventHandlerChannel:
			if !ok {
				fmt.Printf("Event handler channel closed, exit Reader Query Go routine\n")
				return
			}
			handleEvent(&evtInfo)
		}
	}
}

func handleEvent(evtInfo *uenibreader.DcEvent) {
	log := logEvent{function: evtInfo.EventType.String()}
	log.lines = append(log.lines, fmt.Sprintf("Event = %v", *evtInfo))
	logChannel <- log
	switch evtInfo.EventType {
	case uenibreader.DC_EVENT_S1UL_TUNNEL_ESTABLISH:
		doSomeDbQueries(evtInfo)
	case uenibreader.DC_EVENT_S1UL_TUNNEL_RELEASE:
		doSomeDbQueries(evtInfo)
	case uenibreader.DC_EVENT_ADD:
	case uenibreader.DC_EVENT_REMOVE:
	case uenibreader.DC_EVENT_GNB_ALL_UES_REMOVE:
	}
}

func doSomeDbQueries(evtInfo *uenibreader.DcEvent) {
	doSomeUeDbQueries(evtInfo)
}

func doSomeUeDbQueries(evtInfo *uenibreader.DcEvent) {
	ueID := &evtInfo.UeID
	log := logEvent{function: evtInfo.EventType.String()}
	eNbUeX2ApID, err := myReader.GetMeNbUEX2APID(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetMeNbUEX2APID(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}
	gNbUeX2ApID, err := myReader.GetSgNbUEX2APID(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetSgNbUEX2APID(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}
	ueState, _ := myReader.GetState(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetState(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}
	cell, err := myReader.GetPsCell(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetPsCell(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}
	erabIDs, err := myReader.GetBearerIDs(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetBearerIDs(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}
	erabs, _ := myReader.GetBearers(ueID)
	if err != nil {
		panic(fmt.Sprintf("GetBearers(%s) failed, error: %s\n", ueID.String(), err.Error()))
	}

	log.lines = append(log.lines, fmt.Sprintf("GetMeNbUEX2APID(%s) = %d", ueID.String(), eNbUeX2ApID))
	log.lines = append(log.lines, fmt.Sprintf("GetSgNbUEX2APID(%s) = %d", ueID.String(), gNbUeX2ApID))
	log.lines = append(log.lines, fmt.Sprintf("GetState(%s) = %v", ueID.String(), ueState))
	log.lines = append(log.lines, fmt.Sprintf("GetPsCell(%s) = %v", ueID.String(), cell))
	log.lines = append(log.lines, fmt.Sprintf("GetBearerIDs(%s) = %v", ueID.String(), erabIDs))
	log.lines = append(log.lines, fmt.Sprintf("GetBearers(%s) = %v", ueID.String(), erabs))
	logChannel <- log

	for _, erabID := range erabIDs {
		doSomeUeErabDbQueries(evtInfo, erabID)
	}
}

func doSomeUeErabDbQueries(evtInfo *uenibreader.DcEvent, erabID uenib.ErabID) {
	ueID := &evtInfo.UeID
	log := logEvent{function: evtInfo.EventType.String()}

	te, err := myReader.GetErabS1ULGtpTE(ueID, erabID)
	if err != nil {
		panic(fmt.Sprintf("GetErabS1ULGtpTE(%s, %d) failed, error: %s\n", ueID.String(), erabID, err.Error()))
	}
	teAddr, err := myReader.GetErabS1ULGtpTEAddr(ueID, erabID)
	if err != nil {
		panic(fmt.Sprintf("GetErabS1ULGtpTEAddr(%s, %d) failed, error: %s\n", ueID.String(), erabID, err.Error()))
	}
	teTeid, err := myReader.GetErabS1ULGtpTETeid(ueID, erabID)
	if err != nil {
		panic(fmt.Sprintf("GetErabS1ULGtpTETeid(%s, %d) failed, error: %s\n", ueID.String(), erabID, err.Error()))
	}
	arpPL, err := myReader.GetErabQosArpPL(ueID, erabID)
	if err != nil {
		panic(fmt.Sprintf("GetErabQosArpPL(%s, %d) failed, error: %s\n", ueID.String(), erabID, err.Error()))
	}
	qci, err := myReader.GetErabQosQci(ueID, erabID)
	if err != nil {
		panic(fmt.Sprintf("GetErabQosQci(%s, %d) failed, error: %s\n", ueID.String(), erabID, err.Error()))
	}

	log.lines = append(log.lines, fmt.Sprintf("GetErabS1ULGtpTE(%s, %d) = %v", ueID.String(), erabID, te))
	log.lines = append(log.lines, fmt.Sprintf("GetErabS1ULGtpTEAddr(%s, %d) = %v", ueID.String(), erabID, teAddr))
	log.lines = append(log.lines, fmt.Sprintf("GetErabS1ULGtpTETeid(%s, %d) = %v", ueID.String(), erabID, teTeid))
	log.lines = append(log.lines, fmt.Sprintf("GetErabQosArpPL(%s, %d) = %d", ueID.String(), erabID, arpPL))
	log.lines = append(log.lines, fmt.Sprintf("GetErabQosQos(%s, %d) = %d", ueID.String(), erabID, qci))

	logChannel <- log
}

func logger(wg *sync.WaitGroup, channel chan logEvent) {
	defer wg.Done()
	for event := range channel {
		fmt.Printf("\n%s\n", time.Now().Format(time.StampMicro))
		for _, line := range event.lines {
			fmt.Println(line)
		}
	}
}

func subscribeEvents() {
	err := myReader.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{uenibreader.DualConnectivity},
		func(evGNb string, eventCategory uenibreader.EventCategory, evs []string) {
			for _, ev := range evs {
				switch eventCategory {
				case uenibreader.DualConnectivity:
					evInfo, err := uenibreader.ParseDcEvent(ev)
					if err != nil {
						panic(fmt.Sprintf("Event parsing failed: %s\n", err.Error()))
					}
					if evInfo.EventType != uenibreader.DC_EVENT_UNKNOWN {
						ueDcEventHandlerChannel <- evInfo
					}
				}
			}
		})
	if err != nil {
		panic(fmt.Sprintf("SubscribeEvents failed, error: %s\n", err.Error()))
	}
}

func setup() {
	subscribeEvents()
	loggerWaitGroup.Add(1)
	go logger(&loggerWaitGroup, logChannel)
}

func teardown() {
	close(ueDcEventHandlerChannel)
	if !wait(&ueDcReaderWaitGroup, time.Second) {
		panic("Timeout while waiting reader closing.")
	}
	close(logChannel)
	if !wait(&loggerWaitGroup, time.Second) {
		panic("Timeout while waiting logger closing.")
	}
}

func wait(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}

//@todo implement a fake UE-NIB writer what would trigger event sending.

func main() {
	setup()
	//Give some time for SDL database backend to make event subscriptions effective
	time.Sleep(time.Second)

	ueDcReaderWaitGroup.Add(1)
	go eventHandler(&ueDcReaderWaitGroup)

	//@todo Better closing, now just close the example after 2 seconds.
	time.Sleep(20 * time.Second)
	teardown()
}

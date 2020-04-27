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
	"regexp"
	"strings"
	"sync"
	"time"
)

type ueEventInfo struct {
	ueID             uenib.UeID
	s1ULTunEndpoints []gtpTunnel
}

type gtpTunnel struct {
	address string
	teid    string
}

type logEvent struct {
	function string
	lines    []string
}

var someGNb string
var myReader *uenibreader.Reader
var ueDcReaderWaitGroup sync.WaitGroup
var ueDcEstablishEventChannel chan ueEventInfo
var ueDcReleaseEventChannel chan ueEventInfo
var loggerWaitGroup sync.WaitGroup
var logChannel chan logEvent

func init() {
	myReader = uenibreader.NewReader()
	//A channel is created per GNb 'someGNb'
	someGNb = "somegnb:310-410-b5c67788"
	//someGNb = "someGNbID0"
	ueDcEstablishEventChannel = make(chan ueEventInfo)
	ueDcReleaseEventChannel = make(chan ueEventInfo)
	logChannel = make(chan logEvent)
}

func queryExecutor(wg *sync.WaitGroup) {
	fmt.Printf("Reader Query Go routine starts\n")
	defer wg.Done()
	for {
		select {
		case info, ok := <-ueDcEstablishEventChannel:
			if !ok {
				fmt.Printf("UE DC establish channel closed, exit Reader Query Go routine\n")
				return
			}
			doSomeDbQueries(&info, "Establish")
		case info, ok := <-ueDcReleaseEventChannel:
			if !ok {
				fmt.Printf("UE DC release channel closed, exit Reader Query Go routine\n")
				return
			}
			doSomeDbQueries(&info, "Release")
		}
	}
}

func doSomeDbQueries(info *ueEventInfo, evType string) {
	doSomeUeDbQueries(&info.ueID, evType)
}

func doSomeUeDbQueries(ueID *uenib.UeID, evType string) {
	log := logEvent{function: evType}
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
		doSomeUeErabDbQueries(ueID, erabID, evType)
	}
}

func doSomeUeErabDbQueries(ueID *uenib.UeID, erabID uenib.ErabID, evType string) {
	log := logEvent{function: evType}

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

func eventMatches(event string, pattern string) bool {
	matched, err := regexp.MatchString(pattern, event)
	if err != nil {
		panic(fmt.Sprintf("regexp.MatchString failed, error: %s\n", err.Error()))
	}
	return matched
}

func parseEventInfo(event string) ueEventInfo {
	var retInfo ueEventInfo
	ueAndS1UPTEP := strings.Split(event, "_")
	if cnt := len(ueAndS1UPTEP); cnt != 2 {
		panic(fmt.Sprintf("Wrong number (%d) of UE/S1UL tunnel endpoints in event string: %s\n", cnt, event))
	}
	ue := ueAndS1UPTEP[0]
	s1UPTEP := ueAndS1UPTEP[1]

	retInfo.ueID = parseEventInfoUeId(ue)
	retInfo.s1ULTunEndpoints = parseEventInfoS1ULTEPs(s1UPTEP)
	return retInfo
}

func parseEventInfoUeId(ueStr string) uenib.UeID {
	var retUeID uenib.UeID
	ueFields := strings.Split(ueStr, "#")
	if cnt := len(ueFields); cnt != 3 {
		panic(fmt.Sprintf("Wrong number (%d) of UE fields in event string: %s\n", cnt, ueStr))
	}
	retUeID.GNb = ueFields[0]
	retUeID.GNbUeX2ApID = ueFields[1]
	retUeID.ENbUeX2ApID = ueFields[2]
	return retUeID
}

func parseEventInfoS1ULTEPs(tepsStr string) []gtpTunnel {
	var retTEPs []gtpTunnel
	tepFields := strings.Split(tepsStr, "#")
	if len(tepFields)%2 != 0 {
		panic(fmt.Sprintf("S1UL tunnel endpoints address and teid count (%d) not even. String: %s\n",
			len(tepFields), tepFields))
	}
	for i := 0; i < len(tepFields); i = i + 2 {
		t := gtpTunnel{
			address: tepFields[i],
			teid:    tepFields[i+1],
		}
		retTEPs = append(retTEPs, t)
	}
	return retTEPs
}

func subscribeEvents() {
	err := myReader.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{uenibreader.DualConnectivity},
		func(evGNb string, eventCategory uenibreader.EventCategory, evs []string) {
			for _, ev := range evs {
				switch eventCategory {
				case uenibreader.DualConnectivity:
					if eventMatches(ev, ".*_S1UL_TUNNEL_ESTABLISH") {
						ueDcEstablishEventChannel <- parseEventInfo(strings.TrimSuffix(ev, "_S1UL_TUNNEL_ESTABLISH"))
					}
					if eventMatches(ev, ".*_S1UL_TUNNEL_RELEASE") {
						ueDcReleaseEventChannel <- parseEventInfo(strings.TrimSuffix(ev, "_S1UL_TUNNEL_RELEASE"))
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
	close(ueDcEstablishEventChannel)
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
	go queryExecutor(&ueDcReaderWaitGroup)

	//@todo Better closing, now just close the example after 2 seconds.
	time.Sleep(20 * time.Second)
	teardown()
}

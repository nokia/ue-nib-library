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
	ueID uenib.UeID
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
	someGNb = "somegnb:310-410-b5c67788"
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
	log := logEvent{function: evType}
	cell, err := myReader.GetPsCell(&info.ueID)
	if err != nil {
		panic(fmt.Sprintf("GetPsCell failed, error: %s\n", err.Error()))
	}
	log.lines = append(log.lines, fmt.Sprintf("GetPsCell(%s)", info.ueID.String()))
	log.lines = append(log.lines, fmt.Sprintf("   PsCell:%v", cell))
	logChannel <- log
	//@todo Add here more examples of other UE-NIB read queries.
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
	ueFields := strings.Split(event, "#")
	if cnt := len(ueFields); cnt != 3 {
		panic(fmt.Sprintf("Too many (%d) UE fields in event string: %s\n", cnt, event))
	}
	retInfo.ueID.GNb = ueFields[0]
	retInfo.ueID.GNbUeX2ApID = ueFields[1]
	retInfo.ueID.ENbUeX2ApID = ueFields[2]
	return retInfo
}

func subscribeEvents() {
	err := myReader.SubscribeEvents([]string{someGNb}, []uenibreader.EventCategory{uenibreader.DualConnectivity},
		func(evGNb string, eventCategory uenibreader.EventCategory, evs []string) {
			for _, ev := range evs {
				switch eventCategory {
				case uenibreader.DualConnectivity:
					if eventMatches(ev, ".*_ESTABLISH") {
						ueDcEstablishEventChannel <- parseEventInfo(strings.TrimSuffix(ev, "_ESTABLISH"))
					}
					if eventMatches(ev, ".*_RELEASE") {
						ueDcReleaseEventChannel <- parseEventInfo(strings.TrimSuffix(ev, "_RELEASE"))
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
	time.Sleep(2 * time.Second)
	teardown()
}

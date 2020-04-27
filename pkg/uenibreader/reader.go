/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

//Package uenibreader implements UE-NIB database event subscription and data query functions.
package uenibreader

import (
	"fmt"
	sdl "gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
)

//Reader is used to read UE data from RIC Radio Network Information Base (UE-NIB) database.
//NOTE: Use NewReader() function to create a Reader instance.
type Reader struct {
	db iDbBackend
}

//NewReader creates and initializes a new Reader instance.
func NewReader() *Reader {
	reader := &Reader{}
	if !disableSdlCreationInConstructor {
		reader.setDbBackend(sdl.NewSdlInstance("uenib", sdl.NewDatabase()))
	}
	return reader
}

//Close closes the connection to the database.
//It is recommended to call Close() after Reader is not used any more, otherwise client process may
//have hanging file descriptor open for the socket which was used for the backend database
//connection.
//In failure case Close() returns an error value indicating an abnormal state.
//In addition to Error() method defined in built-in error interface a function caller can test
//returned error value for a reader.Error with a type assertion and then distinguish temporal errors
//from permanent ones by using Temporary() method. In case of temporal error, the caller of
//Close() may retry the call after a short period of time.
func (reader *Reader) Close() error {
	err := reader.db.Close()
	if err != nil {
		return newBackendError(err.Error())
	}
	return err
}

//Variable used for bypassing the real SDL database backend usage in unit tests of UE-NIB.
var disableSdlCreationInConstructor bool

//SDL database backend interface.
type iDbBackend interface {
	Get(keys []string) (map[string]interface{}, error)
	SubscribeChannel(cb func(string, ...string), channels ...string) error
	Close() error
}

func setDisableSdlCreationInConstructor(disabled bool) {
	disableSdlCreationInConstructor = disabled
}

func (reader *Reader) setDbBackend(dbBackend iDbBackend) {
	reader.db = dbBackend
}

func newValidationError(err string, vals ...interface{}) *validationError {
	return &validationError{err: fmt.Sprintf(err, vals...)}
}
func newInternalError(err string) *internalError {
	return &internalError{err: err}
}

func newBackendError(err string) *backendError {
	return &backendError{err: err, temporary: true}
}

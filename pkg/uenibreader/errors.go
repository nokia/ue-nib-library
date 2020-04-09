/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import (
	"fmt"
)

//An Error interface represents a UE-NIB error.
type Error interface {
	error            //Embedded built-in error interface
	Temporary() bool //Returns true if an error is temporary
}

//An validationError is UE-NIB private type to hold classification data of validation type of error.
type validationError struct {
	err string //Error message
}

//An internalError is UE-NIB private type to hold classification data of internal type of error.
type internalError struct {
	err string //Error message
}

//A backendError is UE-NIB private type to hold classification data of database backend error
//type of error. Backend errors are temporal by their nature. UE-NIB API user is adviced to try
//again failed UE-NIB API operation.
type backendError struct {
	err       string //Error message
	temporary bool   //Defines whether the error is temporary or not
}

//Error implements built-in error interface for validationError type.
func (e *validationError) Error() string {
	return fmt.Sprintf("UE-NIB validation error: %s", e.err)
}

//Error implements built-in error interface for internalError type.
func (e *internalError) Error() string {
	return fmt.Sprintf("UE-NIB internal error: %s", e.err)
}

//Error implements built-in error interface for backendError type.
func (e *backendError) Error() string {
	return fmt.Sprintf("UE-NIB database backend error: %s", e.err)
}

//Temporary implements Error interface for backendError type.
//Returns true if backend failure was temporal and it is recommended to re-try failed UE-NIB
//operation.
func (f *backendError) Temporary() bool {
	return f.temporary
}

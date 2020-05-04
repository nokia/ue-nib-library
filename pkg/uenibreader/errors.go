/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

import (
	"fmt"
	"github.com/nokia/ue-nib-library/pkg/uenib"
)

//An Error interface represents a UE-NIB error.
type Error interface {
	error            //Embedded built-in error interface
	Temporary() bool //Returns true if an error is temporary and it is recommended to re-try failed operation
}

//A valueNotFoundFailure is UE-NIB private type for a circumstance when queried value is not found from database.
type valueNotFoundFailure struct {
	ueID      uenib.UeID //Identity of a user in question
	name      string     //Parameter name
	temporary bool       //Defines whether the error is temporary or not
}

//An validationError is UE-NIB private type to hold classification data of validation type of error.
type validationError struct {
	ueID      uenib.UeID //Identity of a user in question
	err       string     //Error message
	temporary bool       //Defines whether the error is temporary or not
}

//An internalError is UE-NIB private type to hold classification data of internal type of error.
type internalError struct {
	ueID      uenib.UeID //Identity of a user in question
	err       string     //Error message
	temporary bool       //Defines whether the error is temporary or not
}

//A backendError is UE-NIB private type to hold classification data of database backend error
//type of error. Backend errors are temporal by their nature. UE-NIB API user is adviced to try
//again failed UE-NIB API operation.
type backendError struct {
	ueID      uenib.UeID //Identity of a user in question
	err       string     //Error message
	temporary bool       //Defines whether the error is temporary or not
}

//Error implements built-in error interface for valueNotFoundFailure type.
func (e *valueNotFoundFailure) Error() string {
	return fmt.Sprintf("UE-NIB %s value of DB key '%s' not found", e.ueID.String(), e.name)
}

//IsValueNotFoundFailure returns true if failure is UE-NIB valueNotFoundFailure type.
func IsValueNotFoundFailure(e interface{}) bool {
	if _, ok := e.(*valueNotFoundFailure); ok {
		return true
	}
	return false
}

//Temporary implements Error interface for valueNotFoundFailure type.
//Returns always true for a valueNotFoundFailure error type. Failure is temporal and hence it is
//recommended to re-try failed UE-NIB operation.
func (e *valueNotFoundFailure) Temporary() bool {
	return e.temporary
}

//Error implements built-in error interface for validationError type.
func (e *validationError) Error() string {
	return fmt.Sprintf("UE-NIB %s validation error: %s", e.ueID.String(), e.err)
}

//IsValidationError returns true if an error is UE-NIB validationError type.
func IsValidationError(e interface{}) bool {
	if _, ok := e.(*validationError); ok {
		return true
	}
	return false
}

//Temporary implements Error interface for validationError type.
//Returns always false for a validationError error type. Error is permanent and hence is not worth
//to re-try failed UE-NIB operation.
func (e *validationError) Temporary() bool {
	return e.temporary
}

//Error implements built-in error interface for internalError type.
func (e *internalError) Error() string {
	return fmt.Sprintf("UE-NIB %s internal error: %s", e.ueID.String(), e.err)
}

//IsInternalError returns true if an error is UE-NIB internalError type.
func IsInternalError(e interface{}) bool {
	if _, ok := e.(*internalError); ok {
		return true
	}
	return false
}

//Temporary implements Error interface for internalError type.
//Returns always false for an internalError error type. Error is permanent and hence is not worth
//to re-try failed UE-NIB operation.
func (e *internalError) Temporary() bool {
	return e.temporary
}

//Error implements built-in error interface for backendError type.
func (e *backendError) Error() string {
	return fmt.Sprintf("UE-NIB %s database backend error: %s", e.ueID.String(), e.err)
}

//IsBackendError returns true if an error is UE-NIB backendError type.
func IsBackendError(e interface{}) bool {
	if _, ok := e.(*backendError); ok {
		return true
	}
	return false
}

//Temporary implements Error interface for backendError type.
//Returns always true for a backendError error type. Error is temporal and hence it is recommended
//to re-try failed UE-NIB operation.
func (e *backendError) Temporary() bool {
	return e.temporary
}

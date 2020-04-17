/*
   Copyright (c) 2020 Nokia.

   Licensed under the BSD 3-Clause Clear License.
   SPDX-License-Identifier: BSD-3-Clause-Clear
*/

package uenibreader

//SetDisableSdlCreationInConstructor exports the private setDisableSdlCreationInConstructor
//function for unit tests.
func SetDisableSdlCreationInConstructor(disabled bool) {
	setDisableSdlCreationInConstructor(disabled)
}

//SetDbBackend exports the private setDbBackend function for unit tests.
//Used to inject mock implementation for database operations.
func (reader *Reader) SetDbBackend(dbBackend iDbBackend) {
	reader.setDbBackend(dbBackend)
}

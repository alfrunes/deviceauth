// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package main

import (
	"time"
)

type Device struct {
	Id          string    `json:"id"`
	TenantToken string    `json:"tenant_token"`
	PubKey      string    `json:"pubkey"`
	IdData      string    `json:"id_data"`
	Status      string    `json:"status"`
	CreatedTs   time.Time `json:"created_ts"`
	UpdatedTs   time.Time `json:"updated_ts"`
}

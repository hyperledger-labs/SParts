package main

// This file contains the restful API support functions.

// Licensing: (Apache-2.0 AND BSD-3-Clause AND BSD-2-Clause)

/*
 * NOTICE:
 * =======
 *  Copyright (c) 2018 Wind River Systems, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software  distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
 * OR CONDITIONS OF ANY KIND, either express or implied.
 */

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	//"log"
	//"path"
	//"path/filepath"
	"net/http"
	//"os/user"
)

const (
	_LEDGER_ADDRESS_KEY    = "node.ledger_address"
	_CONDUCTOR_ADDRESS_KEY = "node.conductor_address"
)

// Ledger Node record
type UUIDRecord struct {
	UUID string `json:"uuid"` // UUID provide w/previous registration
}

func httpGetAPIRequest(net_address string, api_request string) ([]byte, error) {
	response, err := http.Get("http://" + net_address + api_request)
	if err != nil {
		////fmt.Printf("%s\n", err)
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			////fmt.Printf("%s\n", err)
			return nil, err
		}
		return contents, nil
	}
}

type StatusMsg struct {
	Status string `json:"status"`
}

func pingLedger() string {

	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		"/api/sparts/ping")

	var msg StatusMsg
	err = json.Unmarshal(replyAsBytes, &msg)
	if err != nil {
		return ""
	}
	// return status
	return msg.Status
}

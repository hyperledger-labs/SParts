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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	//"log"
	//"path/filepath"
	"net/http"
)

// sendGetRequest handles the sending of  all GET restful api calls
// and assigns the reply contents to 'reply'
func sendGetRequest(serverAddress string, apiCall string, reply interface{}) error {

	if _DEBUG_REST_API_ON {
		fmt.Printf("Using http server api url: %s%s\n", serverAddress, apiCall)
	}

	// create reply record of the same type
	replyRecord := reply

	replyAsBytes, err := httpGetAPIRequest(serverAddress, apiCall)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("http server: '%s' is not accessible", serverAddress)
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(replyAsBytes, &objmap)
	if err != nil {
		reply = nil
		return fmt.Errorf("unable to parse http server reponse")
	}
	// Make sure status field exists
	if _, ok := objmap["status"]; !ok {
		reply = nil
		return fmt.Errorf("'status' field is missing from server's response")
	}
	var resultStatus string
	err = json.Unmarshal(*objmap["status"], &resultStatus)
	if err != nil {
		reply = nil
		return fmt.Errorf("problem accessing 'status' field from '%s' server's response", serverAddress)
	}
	if resultStatus != _SUCCESS {
		var message string
		err = json.Unmarshal(*objmap["message"], &message)
		return fmt.Errorf("received failed response from http server: %s - %s", serverAddress, message)
	}

	// Make sure 'result_type' field exists - to avoid panic
	if _, ok := objmap["result_type"]; !ok {
		reply = nil
		return fmt.Errorf("'result_type' field is missing from service response")
	}
	var resultType string
	err = json.Unmarshal(*objmap["result_type"], &resultType)
	if err != nil {
		reply = nil
		return fmt.Errorf("Reply type is missing from service response")
	}

	// __resultType := "*" + resultType // need to add prefix '*' to type string
	// Need to remove '*' from type - e.g., *ListOf:SupplierRecord -> ListOf:SupplierRecord
	expectedType := strings.Replace(getType(replyRecord), "*", "", 1)
	////&& strings.ToLower(resultType) != strings.ToLower("ArtifactRecord")
	if strings.ToLower(resultType) != strings.ToLower(expectedType) {
		reply = nil
		return fmt.Errorf("api call '%s' type response '%s' is not valid. Expecting: '%s'", apiCall, resultType, expectedType)
	}

	err = json.Unmarshal(*objmap["result"], &replyRecord)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("Service response may not be properly formatted. Expecting: '%s'", resultType)
	}
	reply = replyRecord
	return nil
}

/****************************************
// sendGetRequest handles the sending of  all GET restful api calls
// and assigns the reply contents to 'reply'
func sendGetRequest(apiCall string, reply interface{}) error {

	if _DEBUG_REST_API_ON {
		fmt.Printf("api url %s\n", apiCall)
	}

	// create reply record of the same type
	replyRecord := reply

	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		apiCall)

	//replyAsBytes, err := httpGetAPIRequest("localhost:3075", "/ledger/api/v1/suppliers")
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("ledger is not accessible")
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(replyAsBytes, &objmap)
	if err != nil {
		reply = nil
		return fmt.Errorf("unable to unmarshal ledger reponse")
	}

	// Make sure status field exists
	if _, ok := objmap["status"]; !ok {
		reply = nil
		return fmt.Errorf("ledger response 'status' field is missing")
	}
	var resultStatus string
	err = json.Unmarshal(*objmap["status"], &resultStatus)
	if err != nil {
		reply = nil
		return fmt.Errorf("problem accessing 'status' field")
	}
	if resultStatus != _SUCCESS {
		var message string
		err = json.Unmarshal(*objmap["message"], &message)
		return fmt.Errorf("received failed response from ledger: %s", message)
	}

	// Make sure 'result_type' field exists - to avoid panic
	if _, ok := objmap["result_type"]; !ok {
		reply = nil
		return fmt.Errorf("ledger response 'result_type' field is missing")
	}
	var resultType string
	err = json.Unmarshal(*objmap["result_type"], &resultType)
	if err != nil {
		reply = nil
		return fmt.Errorf("ledger response type is missing")
	}

	// __resultType := "*" + resultType // need to add prefix '*' to type string
	// Need to remove '*' from type - e.g., *ListOf:SupplierRecord -> ListOf:SupplierRecord
	expectedType := strings.Replace(getType(replyRecord), "*", "", 1)
	////&& strings.ToLower(resultType) != strings.ToLower("ArtifactRecord")
	if strings.ToLower(resultType) != strings.ToLower(expectedType) {
		reply = nil
		return fmt.Errorf("ledger response type '%s' is not valid. Expecting: '%s'", resultType, expectedType)
	}

	err = json.Unmarshal(*objmap["result"], &replyRecord)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("Ledger response may not be properly formatted. Expecting: '%s'", resultType)
	}
	reply = replyRecord
	return nil
}

func sendGetRequest2(server string, apiCall string, reply interface{}) error {

	var serverAddress string
	switch server {
	case _LEDGER:
		serverAddress = getLocalConfigValue(_LEDGER_ADDRESS_KEY)
	case _ATLAS:
		serverAddress = getGlobalConfigValue(_ATLAS_ADDRESS_KEY)
	default:
		// server address not known
		return fmt.Errorf("Address for server '%s' is not known", server)
	}

	if _DEBUG_REST_API_ON {
		fmt.Printf("%s server api url: %s%s\n", serverAddress, apiCall)
	}

	// create reply record of the same type
	replyRecord := reply

	replyAsBytes, err := httpGetAPIRequest(serverAddress, apiCall)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("'%s' server is not accessible", server)
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(replyAsBytes, &objmap)
	if err != nil {
		reply = nil
		return fmt.Errorf("unable to unmarshal server '%s' reponse", server)
	}
	// Make sure status field exists
	if _, ok := objmap["status"]; !ok {
		reply = nil
		return fmt.Errorf("'status' field is missing from '%s' server's response", server)
	}
	var resultStatus string
	err = json.Unmarshal(*objmap["status"], &resultStatus)
	if err != nil {
		reply = nil
		return fmt.Errorf("problem accessing 'status' field from '%s' server's response", server)
	}
	if resultStatus != _SUCCESS {
		var message string
		err = json.Unmarshal(*objmap["message"], &message)
		return fmt.Errorf("received failed response from server '%s': %s", server, message)
	}

	// Make sure 'result_type' field exists - to avoid panic
	if _, ok := objmap["result_type"]; !ok {
		reply = nil
		return fmt.Errorf("'result_type' field is missing from '%s' server's response", server)
	}
	var resultType string
	err = json.Unmarshal(*objmap["result_type"], &resultType)
	if err != nil {
		reply = nil
		return fmt.Errorf("server '%s''s reply type is missing from response")
	}

	// __resultType := "*" + resultType // need to add prefix '*' to type string
	// Need to remove '*' from type - e.g., *ListOf:SupplierRecord -> ListOf:SupplierRecord
	expectedType := strings.Replace(getType(replyRecord), "*", "", 1)
	////&& strings.ToLower(resultType) != strings.ToLower("ArtifactRecord")
	if strings.ToLower(resultType) != strings.ToLower(expectedType) {
		reply = nil
		return fmt.Errorf("%s's type response type '%s' is not valid. Expecting: '%s'", server, resultType, expectedType)
	}

	err = json.Unmarshal(*objmap["result"], &replyRecord)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		reply = nil
		return fmt.Errorf("%s's response may not be properly formatted. Expecting: '%s'", server, resultType)
	}
	reply = replyRecord
	return nil
}
************************/

func httpGetAPIRequest(net_address string, api_request string) ([]byte, error) {
	////response, err := http.Get("http://" + net_address + api_request)
	response, err := http.Get(generateURL(net_address, api_request))
	if _DEBUG_REST_API_ON {
		fmt.Println("Requesting URL:", generateURL(net_address, api_request))
	}
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
		if _DEBUG_REST_API_ON {
			fmt.Printf("%s\n", contents)
		}
		return contents, nil
	}
}

// sendPostRequest handles the marshalling up data for all the http POST restful api calls
// It will check for success if replyReord is of type 'ReplyType' and return errpr == nil if
// successful.
func sendPostRequest(apiCall string, requestRecord interface{}, replyRecord interface{}) error {
	// convert request data to bytes
	recordAsBytes, err := json.Marshal(requestRecord)
	if err != nil {
		replyRecord = nil
		return err
	}
	// Make the http post request
	ledgerAddress := getLocalConfigValue(_LEDGER_ADDRESS_KEY)
	//// requestURL := "http://" + ledgerAddress + apiCall
	requestURL := generateURL(ledgerAddress, apiCall)
	if _DEBUG_REST_API_ON {
		fmt.Println("Requesting URL:", requestURL)
	}
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(recordAsBytes))
	if err != nil {
		////fmt.Printf("Error: %s\n", err)
		replyRecord = nil
		return err
	}

	// Read the ledger response
	req.Header.Set("X-Custom-Header", getType(requestRecord))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Printf("Error: %s\n", err)
		replyRecord = nil
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if _DEBUG_REST_API_ON {
		fmt.Println("--------------------------------")
		fmt.Println("Requesting URL:", requestURL)
		fmt.Println()
		fmt.Printf("POST Record: %s\n\n", recordAsBytes)
		fmt.Println("Response Body: ", string(body))
		fmt.Println("--------------------------------")
	}

	switch getType(replyRecord) {
	case "ReplyType":
		var reply ReplyType
		err = json.Unmarshal(body, &reply)
		if err != nil {
			if strings.Contains(err.Error(), "unexpected end of JSON input") {
				replyRecord = nil
				return fmt.Errorf("ledger is not responding at address: %s", ledgerAddress)
			} else {
				replyRecord = nil
				return fmt.Errorf("ledger is not accessible")
			}
		}
		// check return status
		if reply.Status == _SUCCESS {
			replyRecord = reply
			return nil
		} else {
			replyRecord = nil
			return fmt.Errorf("ledger replied with an error message: %s", reply.Message)
		}
	default:
		// unknow reply type. Nothing to do.
		return nil
	}

	// Don't believe execution reaches this point. Complier expects a return statement.
	return nil
}

// getLedgerNodeList retrieves a list of all the ledger nodes from the Atlas look up service.
func getLedgerNodeList() ([]LedgerNodeRecord, error) {

	////list := []LedgerNodeRecord{}

	if getGlobalConfigValue(_ATLAS_ADDRESS_KEY) == "" {
		return nil, fmt.Errorf("error - %s ip address is not set in global configuration", _ATLAS)
	}
	ok, err := pingServer(_ATLAS)
	if !ok {
		return nil, err
	}

	apiStr := fmt.Sprintf("%s%s", _ATLAS_LIST_LEDGER_NODES_API, getLocalConfigValue(_LEDGER_NETWORK_KEY))
	replyAsBytes, err := httpGetAPIRequest(getGlobalConfigValue(_ATLAS_ADDRESS_KEY), apiStr)

	var reply ReplyType
	err = json.Unmarshal(replyAsBytes, &reply)
	if err != nil {
		if _DEBUG_REST_API_ON && _DEBUG_DISPLAY_ON {
			fmt.Println(err)
		}
		return nil, err
	}

	// check return status
	if !(reply.Status == _SUCCESS) {
		if reply.Status == _FAILURE {
			return nil, fmt.Errorf("error - %s server response: %s", _ATLAS, reply.Message)
		}
		return nil, fmt.Errorf("error - %s server response not successful ", _ATLAS)
	}

	type NodeListReply struct {
		Status  string             `json:"status"`
		Message string             `json:"message"`
		Type    string             `json:"result_type"`
		Result  []LedgerNodeRecord `json:"result,omitempty"`
	}

	// We recieved a successul response. No Check result type
	// to mak sure result is []LedgerNodeRecord.
	var nodeList NodeListReply
	// We need to convert type to a language neutral format.
	if getType(nodeList.Result) == reply.Type {
		err = json.Unmarshal(replyAsBytes, &nodeList)
		if err != nil {
			if _DEBUG_REST_API_ON && _DEBUG_DISPLAY_ON {
				fmt.Println(err)
			}
			return nil, err
		}
	}
	return nodeList.Result, nil
}

// pingServer will send a standard restful api "ping" request to
// a sparts service. If a successful response is recieve it returns 'true'
// otherwise it returns 'false'
func pingServer(server string) (bool, error) {

	var replyAsBytes []byte
	var err error
	var ipAddress string
	var serverType string

	switch strings.ToLower(server) {
	case strings.ToLower(_ATLAS):
		ipAddress = getGlobalConfigValue(_ATLAS_ADDRESS_KEY)
		serverType = fmt.Sprintf("%s directory", _ATLAS)
		if ipAddress == "" {
			return false, fmt.Errorf("'%s' address not set in local config file", _ATLAS_ADDRESS_KEY)
		}
		replyAsBytes, err = httpGetAPIRequest(ipAddress, _ATLAS_PING_API)
	case strings.ToLower(_LEDGER):
		ipAddress = getLocalConfigValue(_LEDGER_ADDRESS_KEY)
		serverType = _LEDGER
		if ipAddress == "" {
			return false, fmt.Errorf("'%s' address not set in local config file", _LEDGER_ADDRESS_KEY)
		}
		replyAsBytes, err = httpGetAPIRequest(ipAddress, _LEDGER_PING_API)
	default:
		ipAddress = server
		serverType = "server"
		if ipAddress == "" {
			return false, fmt.Errorf("'%s' address not set in local config file", server)
		}
		replyAsBytes, err = httpGetAPIRequest(ipAddress, _LEDGER_PING_API)
	}

	var reply ReplyType
	err = json.Unmarshal(replyAsBytes, &reply)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			return false, fmt.Errorf("%s is not responding at address: %s", serverType, ipAddress)
		} else {
			return false, fmt.Errorf("ledger is not accessible")
		}
	}
	// check return status
	if reply.Status == _SUCCESS {
		return true, nil
	} else {
		return false, fmt.Errorf("ledger replied with an error message: %s", reply.Message)
	}
}

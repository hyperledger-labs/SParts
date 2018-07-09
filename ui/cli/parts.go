package main

/*
	The functions for the artifact routines can be found in this file.
*/

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
	//"log"
	//"path"
	//"path/filepath"
	"net/http"
	//"os/user"
	"strings"
)

type PartSupplierRecord struct {
	PartUUID     string `json:"part_uuid"`     // Part uuid
	SupplierUUID string `json:"supplier_uuid"` // Suppler uuid
}

type PartArtifactRecord struct {
	PartUUID     string `json:"part_uuid"`     // Part uuid
	ArtifactUUID string `json:"envelope_uuid"` // Suppler uuid
}

func displayParts(partsList []PartItemRecord) {
	if len(partsList) == 0 {
		// empty list
		return
	}
	////fmt.Println("  Parts: ")
	for k := range partsList {
		part, err := getPartInfo(partsList[k].PartUUID)
		if err != nil {
			// error retrieving part
			fmt.Println("Could not retrieve part for uuid=", partsList[k].PartUUID)
			continue // skip to next part.
		}
		fmt.Println()
		fmt.Println("    " + _CYAN_FG + part.Name + _COLOR_END)
		fmt.Print("    ")
		fmt.Print()
		///whiteSpace := createWhiteSpace (len (part.Name))
		///fmt.Printf ("%s%s%s\n", createLine (part.Name), "       ", "------------------------" )
		fmt.Println("-------------------------------------------------")

		fmt.Println("    Name: \t " + part.Name)
		fmt.Println("    Version: \t " + part.Version)
		fmt.Println("    UUID: \t " + part.UUID)
		// Format the descriptions greater
		chuckSize := 60
		for len(part.Description) > chuckSize && part.Description[chuckSize] != ' ' {
			chuckSize++
		}
		chuckSize++
		descriptionStr := strings.Join(chunkString(part.Description, chuckSize), "\n                 ")
		fmt.Println("    Description: " + descriptionStr)
		fmt.Println()
	}

}

// Create part on ledger
func createPart(name string, version string, label string, licensing string,
	description string, url string, checksum string, uuid string) (bool, error) {

	var part PartRecord

	part.Name = name
	part.Version = version
	part.Alias = label
	part.Label = label
	part.Licensing = licensing
	//part.URI = url
	part.Description = description
	part.Checksum = checksum
	part.UUID = uuid

	var requestRecord PartAddRecord
	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	if requestRecord.PrivateKey == "" || requestRecord.PublicKey == "" {
		return false, fmt.Errorf("Private and/or Public key(s) are not set \n Use 'sparts config' to set keys")
	}
	requestRecord.Part = part

	var replyRecord ReplyType
	err := sendPostRequest(_PARTS_API, requestRecord, replyRecord)
	if err != nil {
		return false, err
	}

	return true, nil
}

func createPartSupplierRelationship(part_uuid string, supplier_uuid string) (bool, error) {
	var requestRecord PartToSupplierRecord
	var partSupplierItem PartSupplierPair

	partSupplierItem.PartUUID = part_uuid
	partSupplierItem.SupplierUUID = supplier_uuid

	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	if requestRecord.PrivateKey == "" || requestRecord.PublicKey == "" {
		return false, fmt.Errorf("Private and/or Public key(s) are not set \n Use 'sparts config' to set keys")
	}
	requestRecord.Relation = partSupplierItem

	var replyRecord ReplyType
	err := sendPostRequest(_PARTS_TO_SUPPLIER_API, requestRecord, replyRecord)
	if err != nil {
		return false, err
	}
	return true, nil

	/******
	partSupplierAsBytes, err := json.Marshal(part_supplier_info)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false, false
	}

	//fmt.Println (string(supplierAsBytes))
	requestURL := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/mapping/PartSupplier"
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(partSupplierAsBytes))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false, false
	}
	req.Header.Set("X-Custom-Header", "PartToSupplier")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false, false
	}

	// Read response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false, false
	}
	fmt.Println("PartToSupplier: response Body:", string(body))
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		return false, true
	} else {
		return false, false
	}
	****/
}

func createPartArtifactRelationship(part_uuid string, artifact_uuid string) bool {
	var partArtifactInfo PartArtifactRecord

	partArtifactInfo.PartUUID = part_uuid
	partArtifactInfo.ArtifactUUID = artifact_uuid

	partArtifactAsBytes, err := json.Marshal(partArtifactInfo)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}

	//fmt.Println (string(supplierAsBytes))
	requestURL := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/parts/AddEnvelope"
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(partArtifactAsBytes))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	req.Header.Set("X-Custom-Header", "PartToArtifact")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	fmt.Println("PartToArtifact: response Body:", string(body))
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		return true
	} else {
		return false
	}
}

func getPartInfo(uuid string) (PartRecord, error) {
	var part PartRecord
	////part.Name = ""
	////part.UUID = ""
	//check that uuid is valid.
	if !isValidUUID(uuid) {
		return part, fmt.Errorf("'%s' UUID is not in a valid format", uuid)
	}

	// WORK AROUND - ledger returning wrong format:
	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		_PARTS_API+"/"+uuid)

	err = json.Unmarshal(replyAsBytes, &part)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		if _DEBUG_REST_API_ON {
			fmt.Printf("\n%s\n", replyAsBytes)
		}
		return part, fmt.Errorf("Ledger response may not be properly formatted")
	}

	/*******
	// WORK AROUND - This is what it SHOULD BE
	err := sendGetRequest(_PARTS_API+"/"+uuid, &part)
	if err != nil {
		// error occurred - return err
		return part, err
	}
	*****/
	/*********
		// TODO: do we need to check returned uuid is same?
		if part.UUID != uuid {
			return part, errors.New(fmt.Sprintf("Part not found in ledger with uuid = '%s'", uuid))
		}
		return part, nil
	}
	*****/
	return part, err
}

func getPartList() ([]PartRecord, error) {
	var partList = []PartRecord{}
	err := sendGetRequest(_PARTS_API, &partList)
	return partList, err
}

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
	"errors"
	"fmt"
	"io/ioutil"
	//"log"
	//"path"
	//"path/filepath"
	"net/http"
	//"os/user"
	"strings"
)

type PartRecord struct {
	Name        string `json:"name"`                  // Fullname
	Version     string `json:"version,omitempty"`     // Version if exists.
	Label       string `json:"label,omitempty"`       // 1-5 alphanumeric characters (unique)
	Licensing   string `json:"licensing,omitempty"`   // License expression
	Description string `json:"description,omitempty"` // Part description (1-3 sentences)
	Checksum    string `json:"checksum,omitempty"`    // License expression
	UUID        string `json:"uuid"`                  // UUID provide w/previous registration
	URI         string `json:"src_uri,omitempty"`     //
}

type PartSupplierRecord struct {
	PartUUID     string `json:"part_uuid"`     // Part uuid
	SupplierUUID string `json:"supplier_uuid"` // Suppler uuid
}

type PartArtifactRecord struct {
	PartUUID     string `json:"part_uuid"`     // Part uuid
	ArtifactUUID string `json:"envelope_uuid"` // Suppler uuid
}

func displayParts(partsList []Part) {
	if len(partsList) == 0 {
		// empty list
		return
	}
	fmt.Println("  Parts: ")

	for k := range partsList {
		part, err := getPart(partsList[k].PartId)
		if err != nil {
			// error retrieving part
			fmt.Println("Could not retrieve part for uuid=", partsList[k].PartId)
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

		chuckSize := 70
		urlStr := strings.Join(chunkString(part.URI, chuckSize), "\n                 ")
		fmt.Println("    Url: \t " + urlStr)

		// Format the descriptions greater
		chuckSize = 60
		for len(part.Description) > chuckSize && part.Description[chuckSize] != ' ' {
			chuckSize++
		}
		chuckSize++
		descriptionStr := strings.Join(chunkString(part.Description, chuckSize), "\n                 ")
		fmt.Println("    Description: " + descriptionStr)
	}

}

// Create part on ledger
func createPart(name string, version string, label string, licensing string,
	description string, url string, checksum string, uuid string) string {
	var part PartRecord

	part.Name = name
	part.Version = version
	part.Label = label
	part.Licensing = licensing
	part.URI = url
	part.Description = description
	part.Checksum = checksum
	part.UUID = uuid

	// convert part data to bytes
	partAsBytes, err := json.Marshal(part)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}

	//fmt.Println (string(supplierAsBytes))
	requestURL := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/parts"
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(partAsBytes))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}
	req.Header.Set("X-Custom-Header", "part value")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Create Part: response Body:", string(body))
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		return part.UUID
	} else {
		return ""
	}
}

func createPartSupplierRelationship(part_uuid string, supplier_uuid string) bool {
	var part_supplier_info PartSupplierRecord

	part_supplier_info.PartUUID = part_uuid
	part_supplier_info.SupplierUUID = supplier_uuid

	partSupplierAsBytes, err := json.Marshal(part_supplier_info)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}

	//fmt.Println (string(supplierAsBytes))
	requestURL := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/mapping/PartSupplier"
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(partSupplierAsBytes))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	req.Header.Set("X-Custom-Header", "PartToSupplier")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}

	// Read response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	fmt.Println("PartToSupplier: response Body:", string(body))
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		return true
	} else {
		return false
	}
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

func getPart(uuid string) (PartRecord, error) {
	var part PartRecord

	var err error = nil

	part.Name = ""
	part.UUID = ""

	//check that uuid is valid.
	if !isValidUUID(uuid) {
		err := errors.New(fmt.Sprintf("'%s' UUID is not in a valid format", uuid))
		return part, err
	}
	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		"/api/sparts/ledger/parts/"+uuid)

	err = json.Unmarshal(replyAsBytes, &part)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		return part, errors.New(fmt.Sprintf("Ledger response may not be properly formatted"))
	}

	// Check if supplier exists
	if part.UUID != uuid {
		return part, errors.New(fmt.Sprintf("Part not found in ledger with uuid = '%s'", uuid))
	}
	//fmt.Printf ("Name = %s\t UUID= %s\t Descrip = %s\n", part.Name, part.UUID, part.Description)
	return part, nil

}

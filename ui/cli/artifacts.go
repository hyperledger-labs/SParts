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
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ArtfactItem struct {
	UUID      string `json:"uuid"` // UUID provide w/previous registration
	Timestamp string `json:"timestamp,omitempty"`
}

type ArtifactRecord struct {
	Id           int           `json:"id,omitempty"`           // Database Id
	UUID         string        `json:"uuid"`                   // UUID provide w/previous registration
	Name         string        `json:"filename"`               // File name
	ShortId      string        `json:"short_id,omitempty"`     // 1-5 alphanumeric characters (unique)
	Label        string        `json:"label,omitempty"`        // Display name
	Checksum     string        `json:"checksum"`               // <host_address:port> in  http://<host_address:port>
	URI          string        `json:"uri,omitempty"`          // Universal Resource Identifier
	Path         string        `json:"path,omitempty"`         // Path within Envelope
	OpenChain    string        `json:"openchain,omitempty"`    // True if aritfact was prepared under an OpenChain program
	Type         string        `json:"content_type,omitempty"` // Source, notices, data, spdx, envelope, other
	EnvelopePath string        `json:"local_path,omitempty"`   // Local directory path
	Timestamp    string        `json:"timestamp,omitempty"`    // Timestamp in UTC format
	SubArtifact  []ArtfactItem `json:"sub_artifact,omitempty"` // Timestamp in UTC format
	_verified    bool          `json:"_verified,omitempty"`    // boolean used to compare to artifact lists. Not sent to ledger
}

type Envelope struct {
	Artifacts []ArtifactRecord `json:"artifacts"` // Artifact list
}

func createEnvelopeChecksum(artifactList []ArtifactRecord) string {

	//TODO: sort checksums
	var sha1List string
	for i := 0; i < len(artifactList); i++ {
		// TODO ?? What happens if Checksum is "" ??
		sha1List = sha1List + artifactList[i].Checksum
	}
	hasher := sha1.New()
	hasher.Write([]byte(sha1List))
	//sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	sha := fmt.Sprintf("%x", hasher.Sum(nil))
	return sha
}

func postEnvelopeToledger(artifacts []ArtifactRecord) bool {
	envelope := Envelope{Artifacts: artifacts}

	envelopeAsBytes, err := json.Marshal(envelope)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}

	//fmt.Println (string(supplierAsBytes))
	request_url := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/envelopes"
	req, err := http.NewRequest("POST", request_url, bytes.NewBuffer(envelopeAsBytes))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return false
	}
	req.Header.Set("X-Custom-Header", "CreateEnvelope")
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
	fmt.Println("Create Envelope: response Body:", string(body))
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		return true
	} else {
		return false
	}
}

func getPartArtifacts(part_uuid string) ([]ArtifactRecord, error) {
	//check that uuid is valid.
	if !isValidUUID(part_uuid) {
		return nil, errors.New(fmt.Sprintf("UUID '%s' is not in a valid format", part_uuid))
	}

	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		"/api/sparts/ledger/parts/artifact/"+part_uuid)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		return nil, errors.New(fmt.Sprintf("Ledger may not be accessible."))
	}

	var artifactList []ArtifactRecord
	err = json.Unmarshal(replyAsBytes, &artifactList)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		return nil, errors.New(fmt.Sprintf("Ledger response may not be properly formatted"))
	}

	return artifactList, nil
}

func createEnvelopeFromDirectory(directory string) ([]ArtifactRecord, error) {

	//var anyError error = nil
	artifacts := []ArtifactRecord{}
	var extension string = ""

	// Create envelope artifact record first
	envelope := ArtifactRecord{}
	_, envelope.Name, _, extension = FilenameDirectorySplit(directory)
	envelopeName := envelope.Name // use latter to replace 'name/' with './'
	if extension != ".env" {
		envelope.Name = envelope.Name + ".env"
	}
	envelope.Type = "this"
	envelope.Path = "/"
	envelope.URI = "/"
	envelope.UUID = getUUID()
	envelope.Label = envelope.Name
	envelope.ShortId = envelope.Name
	envelope.OpenChain = "false"
	// add envelope to artifact list
	artifacts = append(artifacts, envelope)
	// Compute envelope index for use later on - typically it should be 0
	envelopeIndex := len(artifacts) - 1

	// Let's traverse the directory collecting up all the files as artifacts.
	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf(" Error: can't access directory: %s\n", directory)
			return err
		}
		if isHidden(path) {
			//fmt.Printf("------ ignored: %s\n", path)
			return nil
		}
		if info.IsDir() {
			// It is a directory
			// Do nothing
			//fmt.Printf("ddddd: visited directory: %s\n", path)
		} else {
			// It is a file
			//create artifact record
			a := ArtifactRecord{}

			a.Path, a.Name, _, a.Type = FilenameDirectorySplit(path)
			a.Path = strings.Replace(a.Path, `\`, `/`, -1)
			// replace envelope name is '.' - e.g.,  env1/dir1/file1 -> ./dir1/file1
			a.Path = strings.Replace(a.Path, envelopeName, ``, 1)
			a.URI = "envelope://" + a.Path
			a.UUID = getUUID()
			//a.UUID = "b17e649e-86f0-4542-639d-0488e7ac0ec9"
			a.Label = a.Name
			a.ShortId = a.Name
			////checksum, _ := getFileSHA1(path)
			a.Checksum, _ = getFileSHA1(path)
			a.OpenChain = "false"
			artifacts = append(artifacts, a)
		}
		return nil
	}) // end of filepath.Walk

	// Compute the envelope checksum which is a function of it's artifact checksums
	artifacts[envelopeIndex].Checksum = createEnvelopeChecksum(artifacts)
	// completed processing directory

	return artifacts, nil
}

func saveEnvelope(artifactsAsJSON string) {
	fmt.Println(" saveEnvelope - Not Implemented")
}

func getArtifactFileType(file string) string {
	_, _, _, extension := FilenameDirectorySplit(file)

	switch strings.ToLower(extension) {
	case ".c":
		return "source"
	default:
		return "tbd"
	}
}

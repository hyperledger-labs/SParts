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

// createEnvelopeChecksum generates the checksum for an envelope.
// The input is a list of artifacts where the first artifact is the
// top level envelope. Other artifacts may also represent an envelope
// but they would be nested envelopes.
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

// postArtifactToLedger adds an artifact to the ledger.
// Return is true if successful, false otherwise.
// error will indicate the error encountered. It does not display
// error messages to the terminal.
func postArtifactToLedger(artifact ArtifactRecord) (bool, error) {
	var replyRecord ReplyType
	var requestRecord ArtifactAddRecord

	// Check uuid
	if !isValidUUID(artifact.UUID) {
		return false, fmt.Errorf("UUID '%s' is not in a valid format", artifact.UUID)
	}

	// TODO: Check for most important fields are filled in

	// Initialize post record.
	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	requestRecord.Artifact = artifact

	// Send artifact post request to ledger
	err := sendPostRequest(_ARTIFACTS_API, requestRecord, replyRecord)
	if err != nil {
		return false, err
	}

	// See if any uris to add
	if len(artifact.URIList) == 0 {
		return true, nil
	}

	const _ARTIFACT_SUCCESS = "Artifact added successfully\n"
	var errorString = _ARTIFACT_SUCCESS
	var uriRequestRecord URIAddRecord
	for _, uri := range artifact.URIList {
		uriRequestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
		uriRequestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
		uriRequestRecord.UUID = artifact.UUID
		uriRequestRecord.URI = uri
		err := sendPostRequest(_ARTIFACTS_URI_API, requestRecord, replyRecord)
		if err != nil {
			errorString += fmt.Sprintf("problem with adding uri %s to artifact\n", uri.Location)
		}
	}

	if errorString == _ARTIFACT_SUCCESS {
		return true, nil
	} else {
		return false, fmt.Errorf("%s", errorString)
	}
}

func postEnvelopeToledger(artifacts []ArtifactRecord) bool {
	envelope := artifacts

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

// getPartArtifacts accepts a part uuid and returns a list of artifact records
// for the part. The func does not display error messages - they are returned to
// the calling routine.
func getPartArtifacts(part_uuid string) ([]ArtifactRecord, error) {
	//check that uuid is valid.
	if !isValidUUID(part_uuid) {
		return nil, fmt.Errorf("UUID '%s' is not in a valid format", part_uuid)
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

// createEnvelopeFromDirectory generates a list of artifacts for the files
// contained within a designate directory. Func returns a list of artifact records
// where the first represents the top level envelope and the remainining artifact
// artifact records represent each of files in the directory (and sub directory)
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
	envelope.ContentType = _ENVELOPE_TYPE
	envelope.UUID = getUUID()
	envelope.Label = envelope.Name
	envelope.Alias = envelope.Name
	envelope.OpenChain = "false"
	envelope.URIList = []URIRecord{} // initalize it the empty list
	// add envelope to artifact list
	artifacts = append(artifacts, envelope)
	//Compute envelope index for use later on - typically it should be 0
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

			a._path, a.Name, _, a.ContentType = FilenameDirectorySplit(path)
			a._path = strings.Replace(a._path, `\`, `/`, -1)
			// replace envelope name is '.' - e.g.,  env1/dir1/file1 -> ./dir1/file1
			a._path = strings.Replace(a._path, envelopeName, ``, 1)
			//a.path =   "envelope://" + a.Path
			a.UUID = getUUID()
			a.Label = a.Name
			a.Alias = a.Name
			////checksum, _ := getFileSHA1(path)
			a.Checksum, _ = getFileSHA1(path)
			a.OpenChain = "false"
			a.URIList = []URIRecord{} // initalize it the empty list
			artifacts = append(artifacts, a)
		}
		return nil
	}) // end of filepath.Walk

	// Compute the envelope checksum which is a function of it's artifact checksums
	/////envelope.Checksum = createEnvelopeChecksum(artifacts)
	artifacts[envelopeIndex].Checksum = createEnvelopeChecksum(artifacts)

	// Create envelope artifact list
	artifactItemList := []ArtifactItem{}
	for i, artifact := range artifacts {
		// Don't include the envelope to the artifact list.
		if i == envelopeIndex ||
			artifacts[envelopeIndex].UUID == artifact.UUID {
			continue
		}
		var item ArtifactItem
		item.UUID = artifact.UUID
		item.Path = artifact._path
		artifactItemList = append(artifactItemList, item)
	}

	artifacts[envelopeIndex].ArtifactList = artifactItemList
	return artifacts, nil
}

// getArtifactFileType determines the artifact type of a file based on the
// files extension. For example, extension ".c" returns type "source".
// Return types are: "binary/audio", binary/executable", "binary/image",
//    "binary/video", "document", "data", "source", "other"
func getArtifactFileType(file string) string {
	_, _, _, extension := FilenameDirectorySplit(file)

	switch strings.ToLower(extension) {
	case ".aac", ".aiff", ".alac", ".flac", ".mp3", ".pcm", ".wav", ".wma":
		return "binary/audio"
	case ".exe", ".jar", ".lib", ".scr", ".so":
		return "binary/executable"
	case ".gif", ".ico", ".jpg", ".jpeg", ".png", ".ttf":
		return "binary/image"
	case ".acc", ".avi", ".flv", ".mov", ".mpg", ".mp2", ".mpeg", ".mpe", ".mpv", ".mp4", ".m4p",
		".oca", ".ogg", ".wmv":
		return "binary/video"
	case ".doc", ".html", ".jnl", ".md", ".pdf", ".ps", ".txt":
		return "document"
	case ".db", ".conf", ".config", ".log":
		return "data"
	case ".asm", ".asp", ".awk", ".bat", ".c", ".class", ".cmd", ".cpp", ".cxx",
		".def", ".dll", ".dpc", ".dpj", ".dtd", ".dump", ".font",
		".h", ".hdl", ".hpp", ".hrc", ".hxx", ".idl", ".inc", ".ini",
		".java", ".js", ".jsp", ".l", ".pl", ".perl", ".pm", ".pmk", ".r", ".rc",
		".res", "rpm", ".s", "sbl", ".sh", ".src", ".tar", ".url", ".y", "yxx":
		return "source"
	case ".bz2", ".gz", ".tgz", ".xz", ".zip":
		return "source"
	default:
		return "other"
	}
}

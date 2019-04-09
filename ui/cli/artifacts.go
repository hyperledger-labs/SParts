package main

/*
	The functions for the artifact routines can be found here.
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
	"strconv"
	"strings"
	"text/tabwriter"
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

func addURIToArtifact(artifactUUID string, uri URIRecord) (bool, error) {
	const _ARTIFACT_SUCCESS = "Artifact added successfully\n"

	var replyRecord ReplyType
	var errorString = _ARTIFACT_SUCCESS
	var uriRequestRecord URIAddRecord
	// Check uuid
	if !isValidUUID(artifactUUID) {
		return false, fmt.Errorf("UUID '%s' is not in a valid format", artifactUUID)
	}

	uriRequestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	uriRequestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	uriRequestRecord.UUID = artifactUUID
	uriRequestRecord.URI = uri
	err := sendPostRequest(_ARTIFACTS_URI_API, uriRequestRecord, replyRecord)
	if err != nil {
		errorString += fmt.Sprintf("problem with adding uri %s to artifact %s\n", uri.Location, artifactUUID)
	}

	if errorString == _ARTIFACT_SUCCESS {
		return true, nil
	} else {
		return false, fmt.Errorf("%s", errorString)
	}
}

// pushArtifactToLedger adds an artifact to the ledger.
// Return is true if successful, false otherwise.
// error will indicate the error encountered. It does not display
// error messages to the terminal.
func pushArtifactToLedger(artifact ArtifactRecord) (bool, error) {
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
	uriRequestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	uriRequestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	uriRequestRecord.UUID = artifact.UUID
	for _, uri := range artifact.URIList {
		uriRequestRecord.URI = uri
		err := sendPostRequest(_ARTIFACTS_URI_API, uriRequestRecord, replyRecord)
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

	request_url := generateURL(getLocalConfigValue(_LEDGER_ADDRESS_KEY), "/api/sparts/ledger/envelopes")
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

func getArtifactFromLedger(artifactUUID string) (ArtifactRecord, error) {
	var artifact = ArtifactRecord{}

	//check that uuid is valid.
	if !isValidUUID(artifactUUID) {
		return artifact, fmt.Errorf("UUID '%s' is not in a valid format", artifactUUID)
	}

	err := sendGetRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY), _ARTIFACTS_API+"/"+artifactUUID, &artifact)

	return artifact, err
}

// getPartArtifacts accepts a part uuid and returns a list of artifact records
// for the part. The func does not display error messages - they are returned to
// the calling routine.
func getEnvelopeArtifactsFromLedger(envelopeUUID string) ([]ArtifactRecord, error) {

	var envelope ArtifactRecord
	var list = []ArtifactRecord{}
	var err error

	//check that uuid is valid.
	if !isValidUUID(envelopeUUID) {
		return list, fmt.Errorf("UUID '%s' is not in a valid format", envelopeUUID)
	}

	envelope, err = getArtifactFromLedger(envelopeUUID)
	if err != nil {
		return list, err
	}

	for _, artifactItem := range envelope.ArtifactList {
		artifactRecord, err := getArtifactFromLedger(artifactItem.UUID)
		if err != nil {
			return list, err
		}

		// This is a temporary workaround until the ledger returns field: 'name' and not 'filename'
		// .Name2 == filename
		if len(artifactRecord.Name) == 0 {
			artifactRecord.Name = artifactRecord.Name2
		}

		list = append(list, artifactRecord)
	}

	return list, nil
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

// createArtifactFromFile creates an artifact record from a file's system info
// The default value for:
// 		ArtifactRecord.OpenChain is false.
// A number of errors may occur:
//		- If 'file' does not exist
//		- If 'file' is a directory
//		- file path can't be obtained
func createArtifactFromFile(file string) (ArtifactRecord, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return ArtifactRecord{}, fmt.Errorf("No such file or directory")
	}
	// Check if it is a directory. Expecting a file only.
	if fileInfo.IsDir() {
		return ArtifactRecord{}, fmt.Errorf("'%s' is a directory. Use --dir flag for directories", file)
	}
	// Extra check to make sure we have a file. IsRegular == true. May not be necessary
	if !fileInfo.Mode().IsRegular() {
		return ArtifactRecord{}, fmt.Errorf("'%s' is not a file. Expecting a file", file)
	}

	var name, label, fullpath string

	fullpath, err = filepath.Abs(file)
	if err != nil {
		return ArtifactRecord{}, fmt.Errorf("Could not obtain full path for '%s'", file)
	}

	checksum, err := getFileSHA1(fullpath)
	if err != nil {
		return ArtifactRecord{}, fmt.Errorf("Error computing SHA1 for file: %s", getAbridgedFilePath(fullpath))
	}

	// Create envelope artifact record
	artifact := ArtifactRecord{}
	_, name, label, _ = FilenameDirectorySplit(fullpath)
	// Clean up string on Windows platform replace '\' with '/'
	fullpath = strings.Replace(fullpath, `\`, `/`, -1)

	////artifactName := artifact.Name // use latter to replace 'name/' with './'

	artifact.UUID = getUUID()
	artifact.Name = name
	artifact.Alias = label
	artifact.Label = name
	artifact.Checksum = checksum
	artifact.ContentType = getArtifactFileType(file)
	artifact.OpenChain = _FALSE
	artifact.URIList = []URIRecord{} // initalize to the empty list
	artifact._contentPath = fullpath
	artifact._onLedger = _FALSE
	artifact._envelopeUUID = _NULL_UUID
	artifact._envelopePath = "/" //default value

	return artifact, nil
}

// createEnvelopeFromDirectory generates a list of artifacts for the files
// contained within a designate directory. Func returns a list of artifact records
// where the first represents the top level envelope and the remainining artifact
// artifact records represent each of files in the directory (and sub directory)
func createEnvelopeFromDirectory(directory string, openchainFlag bool) ([]ArtifactRecord, error) {

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
	envelope.OpenChain = _FALSE
	if openchainFlag {
		envelope.OpenChain = _TRUE
	}
	envelope._onLedger = _FALSE
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

			a._envelopePath, a.Name, _, a.ContentType = FilenameDirectorySplit(path)
			a._envelopePath = strings.Replace(a._envelopePath, `\`, `/`, -1)
			// replace envelope name is '.' - e.g.,  env1/dir1/file1 -> ./dir1/file1
			a._envelopePath = strings.Replace(a._envelopePath, envelopeName, ``, 1)
			a._contentPath = path
			a._contentPath = strings.Replace(a._contentPath, `\`, `/`, -1)
			a.UUID = getUUID()
			a.Label = a.Name
			a.Alias = a.Name
			a._envelopeUUID = envelope.UUID
			////checksum, _ := getFileSHA1(path)
			a.Checksum, _ = getFileSHA1(path)
			a.OpenChain = _FALSE
			if openchainFlag {
				a.OpenChain = _TRUE
			}
			a._onLedger = _FALSE
			////a.OpenChain = "false"
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
		item.Path = artifact._contentPath
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
	file = strings.ToLower(file)
	_, _, _, extension := FilenameDirectorySplit(file)

	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		return "url"
	}

	switch extension {
	case ".aac", ".aiff", ".alac", ".flac", ".mp3", ".pcm", ".wav", ".wma":
		return "binary/audio"
	case ".exe", ".jar", ".lib", ".scr", ".so":
		return "binary/executable"
	case ".gif", ".ico", ".jpg", ".jpeg", ".png", ".ttf":
		return "binary/image"
	case ".acc", ".avi", ".flv", ".mov", ".mpg", ".mp2", ".mpeg", ".mpe", ".mpv", ".mp4", ".m4p",
		".oca", ".ogg", ".wmv":
		return "binary/video"
	case ".doc", ".html", ".jnl", ".md", ".pdf", ".ps", ".rst", ".txt", ".text":
		return "document"
	case ".db", ".conf", ".config", ".log":
		return "data"
	case ".asm", ".asp", ".awk", ".bat", ".c", ".class", ".cmd", ".cpp", ".cxx",
		".def", ".dll", ".dpc", ".dpj", ".dtd", ".dump", ".font", ".go",
		".h", ".hdl", ".hpp", ".hrc", ".hxx", ".idl", ".inc", ".ini",
		".java", ".js", ".jsp", ".l", ".pl", ".perl", ".pm", ".pmk", ".r", ".rc",
		".res", "rpm", ".s", "sbl", ".sh", ".src", ".tar", ".url", ".y", "yxx":
		return "source"
	case ".bz2", ".gz", ".tgz", ".xz", ".zip":
		return "source"
	case ".spdx":
		return "spdx"
	default:
		return "other"
	}
}

// displayStagingTable prints the staging table to the terminal.
func displayStagingTable() {

	fmt.Println()
	ledgerNetwork := getLocalConfigValue(_LEDGER_NETWORK_KEY)
	var color string
	if ledgerNetwork == "" || ledgerNetwork == "tbd" {
		color = _RED_FG
		ledgerNetwork = "tdb"
	} else {
		color = _CYAN_FG
	}
	fmt.Printf("  Network: %s%s%s\n", color, ledgerNetwork, _COLOR_END)

	// See if alias is available for part
	part_uuid := getLocalConfigValue(_PART_KEY)
	partAlias, err := getAliasUsingValue(part_uuid)
	if partAlias != "" && err == nil {
		//trimUUID(part_uuid, 5)
		//part_uuid = partAlias + " [" + part_uuid + "]"
		// part_uuid = partAlias + " " + trimUUID(part_uuid, 5)
		part_uuid = partAlias
	}
	if part_uuid == _NULL_UUID {
		part_uuid = _RED_FG + part_uuid + _COLOR_END
	} else {
		part_uuid = _GREEN_FG + part_uuid + _COLOR_END
	}

	// See if alias is available for envelope
	envelope_uuid := getLocalConfigValue(_ENVELOPE_KEY)
	envelopeUUID := envelope_uuid // We will need the unmodified uuid later in func.
	envelopeAlias, err := getAliasUsingValue(envelope_uuid)
	if envelopeAlias != "" && err == nil {
		//envelope_uuid = envelopeAlias + " [" + envelope_uuid + "]"
		//trimUUID(part_uuid, 5)
		//envelope_uuid = envelopeAlias + " " + trimUUID(envelope_uuid, 5)
		envelope_uuid = envelopeAlias
	}

	if envelope_uuid == _NULL_UUID {
		envelope_uuid = _RED_FG + envelope_uuid + _COLOR_END
	} else {
		envelope_uuid = _GREEN_FG + envelope_uuid + _COLOR_END
	}

	fmt.Println(" |--------------------------- Staging --------------------------------")
	// Decide 'focus' - i.e., whether to display 'part',envelope' or both

	//artifacts are grouped into three catagories:
	// focus parts only = orphan

	focus := getLocalConfigValue(_FOCUS_KEY)
	switch focus {
	case _PART_FOCUS:
		fmt.Printf(" |     %s%s%s : %s%s\n", _CYAN_FG, "part", _COLOR_END, part_uuid)
		envelopeUUID = _NULL_UUID
	case _ENVELOPE_FOCUS:
		fmt.Printf(" | %s%s%s : %s\n", _CYAN_FG, "envelope", _COLOR_END, envelope_uuid)
	case _BOTH_FOCUS:
		fmt.Printf(" |     %s%s%s : %s\n", _CYAN_FG, "part", _COLOR_END, part_uuid)
		fmt.Printf(" | %s%s%s : %s\n", _CYAN_FG, "envelope", _COLOR_END, envelope_uuid)
	case _NO_FOCUS:
		// orphan only = !, set envelopeUUID to null
		envelopeUUID = _NULL_UUID
	default:
	}

	displayArtifacts, err := getEnvelopeArtifactList(envelopeUUID, true)
	if err != nil {
		displayErrorMsg(err.Error())
		return
	}

	fmt.Println(" |--------------------------------------------------------------------")

	if len(displayArtifacts) == 0 {
		// nothing waiting to post
		fmt.Println(" |")
		if envelopeUUID == _NULL_UUID {
			fmt.Println(" | [No atifacts have been placed into the staging area for above PART]")
		} else {
			fmt.Println(" | [No atifacts have been placed into the staging area for the above Envelope]")
		}
		fmt.Println()
		fmt.Printf(" Use '%s add' to stage artifacts prior to posting to ledger\n", filepath.Base(os.Args[0]))
		fmt.Println()
		return // we're done
	}

	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	fmt.Fprintf(w, " \t%s \t%s \t%s\t%s\t %s\n", "  Id", "  Name  ", " Type", "OpCh", " File Path or URI")
	fmt.Fprintf(w, " \t%s\t%s \t%s\t%s\t %s\n", " ----", " ----------", "------", "-----", "---------------------")

	var openChain string
	var id, name, path string
	for i := range displayArtifacts {
		if displayArtifacts[i].ContentType == _ENVELOPE_TYPE {
			continue // skip envelopes
		}
		id = strconv.Itoa(displayArtifacts[i]._ID)
		name = displayArtifacts[i].Name
		if len(name) > 40 {
			name = name[0:39]
		}
		if isPathURL(displayArtifacts[i]._contentPath) {
			// path is a url link
			path = displayArtifacts[i]._contentPath
		} else {
			// It is a file, convert path relative to the .sparts working directory.
			////fmt.Println("Path:", artifacts[i]._contentPath)
			//path = getAbridgedFilePath(displayArtifacts[i]._contentPath)
			path = displayArtifacts[i]._envelopePath
		}
		if displayArtifacts[i].OpenChain == _TRUE {
			openChain = " Y"
		} else {
			openChain = " -"
		}

		// Three states"
		//  + added to staging table
		//  > inserted into envelope
		//  = pushed to ledger
		statusChar := ""
		if displayArtifacts[i]._envelopeUUID == _NULL_UUID {
			// added but not save (inserted)
			statusChar = ">"
		} else if displayArtifacts[i]._envelopeUUID == envelopeUUID && displayArtifacts[i]._onLedger == _FALSE {
			// inserted into envelope, but not pushed to ledge
			statusChar = ">"
			//path = displayArtifacts[i]._envelopePath
		} else if displayArtifacts[i]._envelopeUUID == envelopeUUID && displayArtifacts[i]._onLedger == _TRUE {
			// pushed to ledger
			statusChar = "="
		}
		id = statusChar + id
		if len(id) < 3 {
			id = " " + id
		}
		fmt.Fprintf(w, "\t %s\t %s \t%s\t %s\t %s\n", id, name, displayArtifacts[i].ContentType, openChain, path)
	}

	//fmt.Fprintf(w, "\n")

	w.Flush()
	//fmt.Println("  -----")
	fmt.Println(" |--------------------------------------------------------------------")
	fmt.Printf("   %s  New or updated artifact not yet pushed to ledger\n", _PRE_LEDGER_TOKEN)
	fmt.Printf("   %s  Artifact that has been pushed to the ledger\n", _POST_LEDGER_TOKEN)
	fmt.Printf("   tip: use '%s remove id1 id2 ...' to remove items from the staging area\n", filepath.Base(os.Args[0]))
	fmt.Println()
}

// displayStagingTable prints the staging table to the terminal.
func displayStagingTable2() {

	fmt.Println()
	ledgerNetwork := getLocalConfigValue(_LEDGER_NETWORK_KEY)
	var color string
	if ledgerNetwork == "" || ledgerNetwork == "tbd" {
		color = _RED_FG
		ledgerNetwork = "tdb"
	} else {
		color = _CYAN_FG
	}

	fmt.Printf("|--------------------------- %sStaging%s --------------------------------\n", _CYAN_FG, _COLOR_END)
	fmt.Printf("|  %snetwork%s: %s\n", color, _COLOR_END, ledgerNetwork)

	// See if alias is available for part
	part_uuid := getLocalConfigValue(_PART_KEY)
	partAlias, err := getAliasUsingValue(part_uuid)
	if partAlias != "" && err == nil {
		part_uuid = partAlias
	}
	if part_uuid == _NULL_UUID {
		part_uuid = _RED_FG + part_uuid + _COLOR_END
	} else {
		part_uuid = _GREEN_FG + part_uuid + _COLOR_END
	}

	// See if alias is available for envelope
	envelope_uuid := getLocalConfigValue(_ENVELOPE_KEY)
	envelopeUUID := envelope_uuid // We will need the unmodified uuid later in func.
	envelopeAlias, err := getAliasUsingValue(envelope_uuid)
	if envelopeAlias != "" && err == nil {
		envelope_uuid = envelopeAlias
	}

	if envelope_uuid == _NULL_UUID {
		envelope_uuid = _RED_FG + envelope_uuid + _COLOR_END
	} else {
		envelope_uuid = _GREEN_FG + envelope_uuid + _COLOR_END
	}

	// Decide 'focus' - i.e., whether to display 'part',envelope' or both
	//artifacts are grouped into three catagories:
	// focus parts only = orphan
	focus := getLocalConfigValue(_FOCUS_KEY)
	switch focus {
	case _PART_FOCUS:
		fmt.Printf("|     %s%s%s : %s%s\n", _CYAN_FG, "part", _COLOR_END, part_uuid)
		//// orphan only = !, set envelopeUUID to null
		envelopeUUID = _NULL_UUID
	case _ENVELOPE_FOCUS:
		fmt.Printf("| %s%s%s: %s\n", _CYAN_FG, "envelope", _COLOR_END, envelope_uuid)
		//// orphan + envelope = ! + *  envelopeUUID has correct uuid already.
	case _BOTH_FOCUS:
		fmt.Printf("|     %s%s%s: %s\n", _CYAN_FG, "part", _COLOR_END, part_uuid)
		fmt.Printf("| %s%s%s: %s\n", _CYAN_FG, "envelope", _COLOR_END, envelope_uuid)
		//// orphan + envelope = ! + *  envelopeUUID has correct uuid already.
	case _NO_FOCUS:
		// orphan only = !, set envelopeUUID to null
		envelopeUUID = _NULL_UUID
	default:
	}

	displayArtifacts, err := getEnvelopeArtifactList(envelopeUUID, true)
	if err != nil {
		displayErrorMsg(err.Error())
		return
	}

	fmt.Println("|--------------------------------------------------------------------")

	if len(displayArtifacts) == 0 {
		// nothing waiting to post
		fmt.Println(" |")
		if envelopeUUID == _NULL_UUID {
			fmt.Println(" | [No atifacts have been placed into the staging area for above PART]")
		} else {
			fmt.Println(" | [No atifacts have been placed into the staging area for the above Envelope]")
		}
		fmt.Println()
		fmt.Printf(" Use '%s add' to stage artifacts prior to posting to ledger\n", filepath.Base(os.Args[0]))
		fmt.Println()
		return // we're done
	}

	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)

	header := []string{"Id", "  Name  ", " Type ", "OpeCh", "File Path or URI"}
	PrintRow(w, PaintRowUniformly(CyanText, header))
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator

	var openChain string
	var id, name, envelopePath string
	for i := range displayArtifacts {
		if displayArtifacts[i].ContentType == _ENVELOPE_TYPE {
			continue // skip envelopes
		}
		id = strconv.Itoa(displayArtifacts[i]._ID)
		name = displayArtifacts[i].Name
		if len(name) > 40 {
			name = name[0:39]
		}
		if isPathURL(displayArtifacts[i]._contentPath) {
			// path is a url link
			envelopePath = displayArtifacts[i]._contentPath
		} else {
			// It is a file, convert path relative to the .sparts working directory.
			//path = getAbridgedFilePath(displayArtifacts[i]._contentPath)
			envelopePath = displayArtifacts[i]._envelopePath
		}
		if displayArtifacts[i].OpenChain == _TRUE {
			openChain = " Y"
		} else {
			openChain = " -"
		}
		var colors []Color
		//colors = []Color{DefaultText, DefaultText, DefaultText, DefaultText, DefaultText}
		//fmt.Fprintf(w, "\t %s\t %s \t%s\t %s\t %s\n", id, name, displayArtifacts[i].ContentType, openChain, path)
		switch displayArtifacts[i]._onLedger {
		case "false":
			colors = []Color{DefaultText, YellowText, DefaultText, DefaultText, DefaultText}
		case "true":
			colors = []Color{DefaultText, GreenText, DefaultText, DefaultText, DefaultText}
		default:
			colors = []Color{DefaultText, DefaultText, DefaultText, DefaultText, DefaultText}
		}
		PrintRow(w, PaintRow(colors, []string{id, name, displayArtifacts[i].ContentType, openChain, envelopePath}))
	}

	//fmt.Fprintf(w, "\n")

	w.Flush()
	//fmt.Println("  -----")
	fmt.Println("|--------------------------------------------------------------------")
	//fmt.Println("  <>File paths are relative to the sparts working directory.")
	//fmt.Println()
	// fmt.Println("    * New or updated and NOT assigned to a part or envelope")
	fmt.Printf("   %s<name>%s - new or updated artifact NOT yet pushed to ledger\n", YellowText, Reset)
	fmt.Printf("   %s<name>%s - artifact that has been pushed to the ledger\n", GreenText, Reset)
	fmt.Printf("   tip: use '%s remove id1 id2 ...' to remove items from the staging area\n", filepath.Base(os.Args[0]))
	fmt.Println()
}

func getEnvelopeArtifactList(envelopeUUID string, useUnassigned bool) ([]ArtifactRecord, error) {
	envelopeList := []ArtifactRecord{}
	artifactList, err := getArtifactListDB()
	if err != nil {
		return envelopeList, fmt.Errorf("sparts working database not accessible")
	}
	for _, artifact := range artifactList {
		if artifact.ContentType == _ENVELOPE_TYPE {
			continue // skip - artifact is the envelope.
		}
		envelopeList = append(envelopeList, artifact)
		/****
		if artifact._envelopeUUID == envelopeUUID ||
			(useUnassigned && artifact.UUID == _NULL_UUID) {
			envelopeList = append(envelopeList, artifact)
		}
		****/
	}
	return envelopeList, nil
}

/***************
func getEnvelopeArtifactList(envelopeUUID string, nonEnvelope bool) ([]ArtifactRecord, error) {
	var envelopeList []ArtifactRecord
	//var artifact ArtifactRecord
	artifactList, err := getArtifactListDB()
	////here(2, err)
	////fmt.Printf("XName is: '%s', path is: %s\n", artifactList[1].Name, artifactList[1].._onLedger)
	if err != nil {
		return envelopeList, fmt.Errorf("sparts working database not accessible")
	}
	for _, artifact := range artifactList {
		if artifact.UUID == envelopeUUID {
			continue // skip - artifact is the envelope.
		}
		if artifact._envelopeUUID == envelopeUUID ||
			(nonEnvelope && artifact._envelopeUUID == _NULL_UUID) {
			envelopeList = append(envelopeList, artifact)
			////fmt.Printf("Name is: '%s', path is: %s\n", artifact.Name, artifact._envelopePath)
		}
	}
	return envelopeList, nil
}

****************/

func createArtifactOfEnvelopeRelation(artifactUUID string, envelopeUUID string, path string) error {
	var replyRecord ReplyType
	var requestRecord ArtifactOfEnvelopeRecord

	// Check uuid
	if !isValidUUID(artifactUUID) {
		return fmt.Errorf("UUID '%s' for artifact is not in a valid format", artifactUUID)
	}
	if !isValidUUID(envelopeUUID) {
		return fmt.Errorf("UUID '%s' for envelope is not in a valid format", envelopeUUID)
	}

	// TODO: Check for most important fields are filled in
	// Initialize post record.
	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	requestRecord.Relation.ArtifactUUID = artifactUUID
	requestRecord.Relation.EnvelopeUUID = envelopeUUID
	requestRecord.Relation.Path = path
	// Send artifact post request to ledger
	err := sendPostRequest(_ARTIFACT_OF_ENV_API, requestRecord, replyRecord)
	if err != nil {
		return err
	}

	return nil
}

func createArtifactOfPartRelation(artifactUUID string, partUUID string) error {
	var replyRecord ReplyType
	var requestRecord ArtifactOfPartRecord

	// Check uuid
	if !isValidUUID(artifactUUID) {
		return fmt.Errorf("UUID '%s' for artifact is not in a valid format", artifactUUID)
	}
	if !isValidUUID(partUUID) {
		return fmt.Errorf("UUID '%s' for envelope is not in a valid format", partUUID)
	}

	// TODO: Check for most important fields are filled in
	// Initialize post record.
	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	requestRecord.Relation.ArtifactUUID = artifactUUID
	requestRecord.Relation.PartUUID = partUUID

	// Send artifact post request to ledger
	err := sendPostRequest(_ARTIFACT_OF_PART_API, requestRecord, replyRecord)
	if err != nil {
		return err
	}

	return nil
}

func displayListComparison(artifactList1 []ArtifactRecord, artifactList2 []ArtifactRecord,
	listTitle1 string, listTitle2 string,
	artifactSetName1 string, artifactSetName2 string) error {

	var colors []Color
	type comparisonRecord struct {
		Artifact       ArtifactRecord
		ChecksumsMatch bool
		Name1          string
		Name2          string
		NamesMatch     bool
	}

	comparisonList := []comparisonRecord{}
	var record comparisonRecord

	const equalStr = " = "
	const notEqualStr = " X "
	const noMatchStr = "---------"

	for i := 0; i < len(artifactList1); i++ {
		for k := 0; k < len(artifactList2); k++ {
			// check that it is not the envelope container
			if artifactList1[i].ContentType == _ENVELOPE_TYPE || artifactList2[k].ContentType == _ENVELOPE_TYPE {
				if artifactList1[i].ContentType == _ENVELOPE_TYPE {
					artifactList1[i]._verified = true
				}
				if artifactList2[k].ContentType == _ENVELOPE_TYPE {
					artifactList2[k]._verified = true
				}
				continue
			}
			// See if we have a match (that does not involve one of the main envelopes)
			if artifactList1[i].Checksum == artifactList2[k].Checksum {
				// we have a match. Records both names (they could differ, although rare)
				// and whether there names match.
				record.Artifact = artifactList1[i]
				record.Name1 = artifactList1[i].Name
				record.Name2 = artifactList2[k].Name
				record.NamesMatch = record.Name1 == record.Name2
				record.ChecksumsMatch = true
				comparisonList = append(comparisonList, record)
				// mark artifacted as reviewed
				artifactList1[i]._verified = true
				artifactList2[k]._verified = true
			}
		}
	}

	// Run through the FIRST list to see if any unmatched artifacts.
	for i := 0; i < len(artifactList1); i++ {
		if !artifactList1[i]._verified {
			record.Artifact = artifactList1[i]
			record.Name1 = artifactList1[i].Name
			record.Name2 = noMatchStr
			record.NamesMatch = false
			record.ChecksumsMatch = false
			comparisonList = append(comparisonList, record)
		}
	}

	// Run through the SEOCOND list to see if any unmatched artifacts.
	for k := 0; k < len(artifactList2); k++ {
		if !artifactList2[k]._verified {
			record.Artifact = artifactList2[k]
			record.Name1 = noMatchStr
			record.Name2 = artifactList2[k].Name
			record.NamesMatch = false
			record.ChecksumsMatch = false
			comparisonList = append(comparisonList, record)
		}
	}

	// Display comparison table

	fmt.Println()
	const padding = 0
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	//header := []string{"   Artifacts   ", listTitle1, "", listTitle2, "  Notes  "}
	header := []string{" : ", listTitle1, "", listTitle2, "  Checksum  "}
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator
	PrintRow(w, PaintRowUniformly(CyanText, header))
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator

	// Print artifacts
	for i, record := range comparisonList {
		indexStr := strconv.Itoa(i+1) + ":"
		notes := "  " + trimUUID(record.Artifact.Checksum, 5)
		if record.ChecksumsMatch {
			if record.NamesMatch {
				colors = []Color{DefaultText, DefaultText, BrightGreenText, DefaultText, DefaultText}
			} else {
				//checksums are the same but names are not.
				colors = []Color{DefaultText, YellowText, BrightGreenText, YellowText, YellowText}
				notes = fmt.Sprintf("%s, different names", notes)
			}

			PrintRow(w, PaintRow(colors, []string{
				indexStr,
				record.Name1,
				equalStr,
				record.Name2,
				notes}))
		} else if record.Name1 == noMatchStr {
			colors = []Color{DefaultText, RedText, BrightRedText, DefaultText, DefaultText}
			PrintRow(w, PaintRow(colors, []string{
				indexStr,
				record.Name1,
				notEqualStr,
				record.Name2,
				notes}))
		} else {
			colors = []Color{DefaultText, DefaultText, BrightRedText, RedText, DefaultText}
			PrintRow(w, PaintRow(colors, []string{
				indexStr,
				record.Name1,
				notEqualStr,
				record.Name2,
				notes}))
		}
	}

	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header)))
	w.Flush()
	fmt.Printf("  **%s%s%s List\n", _CYAN_FG, artifactSetName1, _COLOR_END)
	fmt.Printf("  ++%s%s%s List\n", _CYAN_FG, artifactSetName2, _COLOR_END)
	fmt.Println()

	return nil

}

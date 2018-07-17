package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

/*
	The functions for the envelope routines can be found here.
*/

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

// displaySupplierList retrieves the supplier list and
// prints the  list to the terminal.
func displayEnvelopeList() {

	var envelopeList []ArtifactRecord //Envelopes are artifacts
	//envelopeList, err := getEnvelopeListFromDB()
	envelopeList, err := getArtifactListInDBWhere("ContentType", _ENVELOPE_TYPE)
	if checkAndReportError(err) {
		return
	}

	// Let's check if the list of suppliers is empty
	if len(envelopeList) == 0 {
		fmt.Println("  no envelopes exist in the working space.")
		return
	}

	//Sort the list
	//supplierList = sortSupplierList(supplierList)

	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ',
		tabwriter.Debug)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", "   Name  ", " Alias", "  UUID  ")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")

	for k := range envelopeList {
		alias, _ := getAliasUsingValue(envelopeList[k].UUID)
		// format alias field for nil values for short length ones
		if alias == "" {
			alias = "   - "
		} else if len(alias) < 4 {
			alias = "  " + alias
		}

		fmt.Fprintf(w, "\t %s\t %s\t %s\n", envelopeList[k].Name, alias, envelopeList[k].UUID)
	}
	//fmt.Println()
	fmt.Fprintf(w, "\n")
	w.Flush()
}

/********************
func pushEnvelopToLedger(envelope ArtifactRecord) error {

	if envelope._notOnLedger == _FALSE {
		return nil
	}

	var replyRecord ReplyType
	var requestRecord ArtifactAddRecord

	// Check uuid
	if !isValidUUID(envelope.UUID) {
		return fmt.Errorf("UUID '%s' is not in a valid format", envelope.UUID)
	}

	// TODO: Check for most important fields are filled in

	// Initialize post record.
	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	requestRecord.Artifact = envelope

	// Send artifact post request to ledger
	err := sendPostRequest(_ARTIFACTS_API, requestRecord, replyRecord)
	if err != nil {
		return err
	}


	// See if any uris to add
	// TODO: push URI to ledger
	//if len(artifact.URIList) == 0 {
	//	return true, nil
	//}


	return nil
}
********************/

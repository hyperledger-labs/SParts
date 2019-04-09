package main

import "regexp"

/*
	PURPOSE:
	   This is the main code for security related functions. For example,
	   public and private key creation.
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

//"github.com/btcsuite/btcd/btcec"              // License: ISC License
//"github.com/btcsuite/btcd/chaincfg"           // License: ISC License
//"github.com/btcsuite/btcd/chaincfg/chainhash" // License: ISC License
//"github.com/btcsuite/btcutil"                 // License: ISC License

func isPrivateKeyValid(privateKey string) bool {
	//r, err := regexp.Compile("^[5KL][1-9A-HJ-NP-Za-km-z]{50,51}$")
	r, err := regexp.Compile("^[5KL][1-9A-HJ-NP-Za-km-z]{51}$")
	if err != nil {
		return false
	}
	if !r.MatchString(privateKey) || len(privateKey) != 51 {
		return false
	}
	return true
}

// Check if WIF public key is valid.
func isPublicKeyValid(publicKey string) bool {
	if len(publicKey) == 66 {
		return true
	} else {
		return false
	}
}

// getSupplierInfo retirieve a single supplier record from the
// ledger. 'uuid' is the id of the supplier.
// supplier.UUID == "" if an error occurs.
func getPrivatePublicKeys() (KeyPairRecord, error) {
	var keys KeyPairRecord

	err := sendGetRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY), _KEY_PAIR_API, &keys)
	return keys, err
}

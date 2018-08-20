package main

/*
	The functions for the supplier routines can be found in this file.
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

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

func getSupplierList() ([]SupplierRecord, error) {

	var supplierList = []SupplierRecord{}
	err := sendGetRequest(_SUPPLIER_API, &supplierList)
	return supplierList, err
}

// displaySupplierList retrieves the supplier list and
// prints the  list to the terminal.
func displaySupplierList() {

	var supplierList []SupplierRecord
	supplierList, err := getSupplierList()
	if checkAndReportError(err) {
		return
	}

	// Let's check if the list of suppliers is empty
	if len(supplierList) == 0 {
		fmt.Println("  No suppliers are registered with the ledger.")
		return
	}

	//Sort the list
	supplierList = sortSupplierList(supplierList)

	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ',
		tabwriter.Debug)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", "   Name  ", " Alias", "  UUID  ")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")

	for k := range supplierList {
		url := supplierList[k].Url
		if url == "" {
			url = "      "
		}
		alias, _ := getAliasUsingValue(supplierList[k].UUID)
		// format alias field for nil values for short length ones
		if alias == "" {
			alias = "   - "
		} else if len(alias) < 4 {
			alias = "  " + alias
		}

		fmt.Fprintf(w, "\t %s\t %s\t %s\n", supplierList[k].Name, alias, supplierList[k].UUID)
	}
	//fmt.Println()
	fmt.Fprintf(w, "\n")
	w.Flush()
}

// getSupplierInfo retirieve a single supplier record from the
// ledger and prints it. 'uuid' is the id of the supplier.
// supplier.UUID == "" if an error occurs.
func getSupplierInfo(uuid string) (SupplierRecord, error) {

	var supplier SupplierRecord

	supplier.UUID = "" // set in case there is an error later
	if !isValidUUID(uuid) {
		return supplier, fmt.Errorf("Supplier UUID is not in a valid format: %s", uuid)
	}

	// err := getReponse(_SUPPLIER_RECORD_API+uuid, &supplier)
	// return supplier, err

	replyAsBytes, err := httpGetAPIRequest(getLocalConfigValue(_LEDGER_ADDRESS_KEY),
		_SUPPLIER_API+"/"+uuid)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		return supplier, errors.New(fmt.Sprintf("Ledger may not be accessible."))
	}

	// Let's unpack the response as json structure
	err = json.Unmarshal(replyAsBytes, &supplier)
	if err != nil {
		if _DEBUG_DISPLAY_ON {
			displayErrorMsg(err.Error())
		}
		return supplier, errors.New(fmt.Sprintf("Ledger response may not be properly formatted"))
	}

	// Check if supplier exists. UUID should of been set in the json.Unmarshal call
	if supplier.UUID != uuid {
		return supplier, fmt.Errorf("Supplier not found with uuid = '%s'", uuid)
	}
	return supplier, nil
}

// displaySupplier prints the supplier record.
func displaySupplier(uuid string) {
	supplier, err := getSupplierInfo(uuid)
	if err != nil || supplier.UUID == "" {
		displayErrorMsg(err.Error())
		return
	}

	alias, _ := getAliasUsingValue(uuid)
	// format alias field for nil or for short length ones
	if alias == "" {
		alias = "<not defined>"
	} else {
		alias = _ALIAS_TOKEN + _CYAN_FG + alias + _COLOR_END
	}
	fmt.Println("  -----------------------------------------------")
	fmt.Printf("  Name   : %s%s%s\n", _CYAN_FG, supplier.Name, _COLOR_END)
	fmt.Println("  -----------------------------------------------")
	fmt.Println("  Label  :", supplier.Alias)
	fmt.Println("  UUID   :", supplier.UUID)
	fmt.Println("  Alias  :", alias)
	fmt.Println("  URL    :", supplier.Url)

	if len(supplier.Parts) == 0 {
		// Supplier has no parts
		fmt.Println("  Parts  : <none> ")
	} else {
		displayPartsFromLedger(supplier.Parts)
	}
}

func pushSupplierToLedger(supplier SupplierRecord) error {

	var requestRecord SupplierAddRecord

	requestRecord.PrivateKey = getLocalConfigValue(_PRIVATE_KEY)
	requestRecord.PublicKey = getLocalConfigValue(_PUBLIC_KEY)
	if requestRecord.PrivateKey == "" || requestRecord.PublicKey == "" {
		return fmt.Errorf("Private and/or Public key(s) are not set \n Use 'sparts config' to set keys")
	}
	requestRecord.Supplier = supplier
	var replyRecord ReplyType
	err := sendPostRequest(_SUPPLIER_API, requestRecord, replyRecord)
	if err != nil {
		return err
	}

	return nil
}

// createSupplier create a supplier and adds it to the ledger.
// if 'uuid' == "" then it will automatically be generated.
func createSupplier(name string, alias string, uuid string, passwd string, url string) (string, error) {
	var supplierInfo SupplierRecord

	if uuid != "" && !isValidUUID(uuid) {
		return uuid, fmt.Errorf("Supplier UUID is not in a valid format.")
	} else if uuid == "" {
		supplierInfo.UUID = getUUID()
	} else {
		supplierInfo.UUID = uuid
	}

	supplierInfo.Name = name
	supplierInfo.Alias = alias
	supplierInfo.Url = url
	supplierInfo.Parts = []PartItemRecord{}

	var replyRecord ReplyType
	err := sendPostRequest(_SUPPLIER_API, supplierInfo, replyRecord)
	if err != nil {
		return uuid, err
	} else {
		return uuid, nil
	}

	/*********************

	supplierAsBytes, err := json.Marshal(supplier)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "", ""
	}

	//fmt.Println (string(supplierAsBytes))
	requestURL := "http://" + getLocalConfigValue(_LEDGER_ADDRESS_KEY) + "/api/sparts/ledger/suppliers"
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(supplierAsBytes))
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
		return "", ""
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//  {"status":"success"}
	if strings.Contains(string(body), "success") {
		//fmt.Println("response Body:", string(body))
		return "", supplier.UUID
	} else {
		return "", ""
	}
	***************************/
}

//----------------------------------------
// Supplier Slice List Sorting
//-----------------------------------------

// This code sorts a list that is represented as a slice.

// We use Go's sort slice functionn.
// https://golang.org/pkg/sort/#Slice

type By func(p1, p2 *SupplierRecord) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(theList []SupplierRecord) {
	ps := &listSorter{
		theList: theList,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// listSorter joins a By function and a slice of SupplierRecord to be sorted.
type listSorter struct {
	theList []SupplierRecord
	by      func(p1, p2 *SupplierRecord) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *listSorter) Len() int {
	return len(s.theList)
}

// Swap is part of sort.Interface.
func (s *listSorter) Swap(i, j int) {
	s.theList[i], s.theList[j] = s.theList[j], s.theList[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *listSorter) Less(i, j int) bool {
	return s.by(&s.theList[i], &s.theList[j])
}

func sortSupplierList(theList []SupplierRecord) []SupplierRecord {

	name := func(p1, p2 *SupplierRecord) bool {
		return p1.Name < p2.Name
	}
	By(name).Sort(theList)
	return theList
}

package main

/*
	The functions for the network space routines can be found here.
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
	"fmt"
	"os"
	"text/tabwriter"
)

func getNetworkList() ([]NetworkSpaceRecord, error) {

	var networkList = []NetworkSpaceRecord{}
	err := sendGetRequest2(_ATLAS, _NETWORK_LIST_API, &networkList)
	return networkList, err
}

// displayNetworkList retrieves the network list and
// prints the list to the terminal.
func displayNetworkList() {

	var networkList []NetworkSpaceRecord
	networkList, err := getNetworkList()
	if checkAndReportError(err) {
		return
	}

	// Let's check if the list of networks is empty
	if len(networkList) == 0 {
		fmt.Println("  No networks are registered.")
		return
	}

	//Sort the list
	//supplierList = sortSupplierList(supplierList)
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ',
		tabwriter.Debug)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", "   Name  ", " Status", "  Description  ")
	fmt.Fprintf(w, "\t%s\t %s\t %s\n", " ------------------", "-------", "------------------------------------")

	var description string
	var name string
	for k := range networkList {
		// We want to print up to 20 characters of the network name.
		if len(networkList[k].Name) <= 25 {
			name = networkList[k].Name
		} else {
			name = networkList[k].Name[:25] + " ..."
		}

		// We can print up to 40 characters of the descripton.
		// testting lenght: networkList[k].Description = "This is a long description to see how well 40 characters work	"
		if len(networkList[k].Description) <= 40 {
			description = networkList[k].Description
		} else {
			description = networkList[k].Description[:40] + " ..."
		}

		fmt.Fprintf(w, "\t %s\t %s\t %s\n", name, networkList[k].Status, description)
	}
	//fmt.Println()
	fmt.Fprintf(w, "\n")
	w.Flush()
}

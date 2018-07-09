package main

/*
	PURPOSE:
   	This is the main code for the sparts comand line interface (CLI)
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
	"bufio"
	"context"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Global
var ArgsAlias []string // Tracks which command line args used which aliases

func main() {
	// make sure sparts command  is provide
	if len(os.Args[1:]) == 0 {
		fmt.Println("Command not found. Expecting sparts command.")
		// Print help overview
		fmt.Println(_HELP_CONTENT)
		return // Done - exit.
	}

	// Let's make sure there is a valid spart directory in current directory (or parent)
	// If not we need to exit. It is ok if it does not exist for the 'init' command.
	_, err := getSpartsDirectory()
	//If an error detected then there is no valid sparts directory.
	if err != nil && strings.ToLower(os.Args[1]) != "init" {
		fmt.Println("  error: Not in a sparts working directory (or any of the parent directories)")
		os.Exit(_DIR_ACCESS_ERROR)
	}

	// sparts directory exists - let's proceed ...

	// Let's see if command line includes any aliases. If so convert
	// alias to value it represents. Save the alias name originally found in
	// the os.Args[i] variable in ArgsAlias [i] so it is accessible later.
	ArgsAlias = append(ArgsAlias, os.Args[0]) // skip the sparts term
	ArgsAlias = append(ArgsAlias, os.Args[1]) // skip the command
	for i := 2; i < len(os.Args); i++ {
		// Look at first characters of each argument to see it if starts with the alias token
		// _ALIAS_TOKEN is the characters to look for
		ArgsAlias = append(ArgsAlias, "") // initial the alias argument to ""
		tokenSize := len(_ALIAS_TOKEN)
		argSize := len(os.Args[i])
		if argSize > tokenSize && os.Args[i][:tokenSize] == _ALIAS_TOKEN {
			// arg i starts with token characters.
			alias := os.Args[i][tokenSize:]
			// fetch value for alias
			value, err := getAlias(alias)
			if checkAndReportError(err) {
				return // error found in use alias - report and exit
			}
			// We have the value. Save a copy in the ArgsAlias array for later reference.
			ArgsAlias[i] = alias
			os.Args[i] = value
		}
	}

	// arg[1] is the sparts command to execute
	switch strings.ToLower(os.Args[1]) {
	case "about":
		fmt.Println(_ABOUT_SPARTS_HELP)
	case "add":
		addRequest()
	case "alias":
		aliasRequest()
	case "artifact":
		artifactRequest()
	case "config":
		configRequest()
	case "compare":
		compareRequest()
	case "delete":
		deleteRequest()
	case "dir":
		dirRequest()
	case "envelope":
		envelopeRequest()
	case "help", "-help", "-h", "--help":
		helpRequest()
	case "init":
		initRequest()
	case "part":
		partRequest()
	case "ping":
		pingRequest()
	case "remove":
		removeRequest()
	case "seed":
		seedRequest()
	case "status":
		statusRequest()
	case "supplier":
		supplierRequest()
	case "synch":
		synchRequest()
	case "test":
		testRequest()
	case "version":
		versionRequest()
	default:
		fmt.Printf("  '%s' is not a valid %s command.\n", os.Args[1], _TOOL_NAME)
		fmt.Printf("  try: %s --help\n", _TOOL_NAME)
	}
}

// sparts add
func addRequest() {
	// if there are no other arguments display help.
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_ADD_HELP_CONTENT)
		return
	}
	// At least one additional argument. Prepare to iterate through args
	// each representing an artifact
	scanner := bufio.NewScanner(os.Stdin)
	for i := 2; i < len(os.Args); i++ {
		artifact_arg := os.Args[i]

		// Check if "help" is the second argument.
		switch strings.ToLower(artifact_arg) {
		case "--help", "-help", "-h":
			if i == 2 {
				// --help, -h is the second argument. Print the help text
				fmt.Println(_ADD_HELP_CONTENT)
				return // we are done
			}
			continue // if --help, -h comes after second argument then ignore it
		}
		var name, label, checksum, uri, fullpath string
		var err error
		fmt.Println()
		// Generate UUID
		uuid := getUUID()

		// See if a uri or local file
		if strings.Contains(artifact_arg, "uri=") {
			// this is a link (uri) and not a local file
			// uri=http://github.com/abc
			splitStr := strings.Split(artifact_arg, "=")
			uri = splitStr[1]
			name = filepath.Base(uri)
			label = name
			fullpath = _NONE
			checksum = _NONE
			fmt.Println("uri:", uri)
		} else {
			// assume artifact is a local file
			fullpath, err = filepath.Abs(artifact_arg)
			if err != nil {
				displayErrorMsg(fmt.Sprintf("Could not locate full path for %s", artifact_arg))
				continue // process next Arg
			}
			// See if a directory
			if stat, err := os.Stat(fullpath); err == nil && stat.IsDir() {
				messageStr := fmt.Sprintf("%s :is a directory", artifact_arg)
				lineStr := createLine(messageStr)
				fmt.Printf(" %s\n", lineStr)
				fmt.Printf(" %s\n", messageStr)
				fmt.Printf(" %s\n", lineStr)
				fmt.Printf(" %s\n", "Directory will be skipped")
				continue // process next Arg
			}

			if _, err := os.Stat(fullpath); os.IsNotExist(err) {
				// path/to/whatever does not exist
				displayErrorMsg(fmt.Sprintf("%s: No such file or directory", os.Args[i]))
				continue // process next Args
			}

			checksum, err = getFileSHA1(fullpath)
			if err != nil {
				fmt.Println("Error computing SHA1 for", fullpath)
				continue // process next Args
			}

			_, name, label, _ = FilenameDirectorySplit(fullpath)

			// Clean up string on Windows platform replace '\' with '/'
			fullpath = strings.Replace(fullpath, `\`, `/`, -1)

		} // else end - file artifact

		// ------------ For all artifacts -----------------
		titleStr := fmt.Sprintf("For artifact '%s':", artifact_arg)
		lineStr := createLine(titleStr)
		fmt.Printf(" %s\n", titleStr)
		fmt.Printf(" %s\n", lineStr)
		atype := ""
		for atype == "" {

			fmt.Println(" Enter artifact type selection (1)-(6):")
			fmt.Printf("   1) Source\n   2) notices\n   3) envelope\n   4) spdx\n   5) data\n   6) Other\n   7) [skip artifact] \n")

			fmt.Print(" > ")
			scanner.Scan()
			atype = scanner.Text()
			switch strings.ToLower(atype) {
			case "1", "source":
				atype = "source"
			case "2", "notices":
				atype = "notices"
			case "3", "envelope":
				atype = "envelope"
			case "4", "spdx":
				atype = "spdx"
			case "5", "data":
				atype = "data"
			case "6", "other":
				atype = "other"
			case "7", "skip":
				atype = "7"
				break
			default:
				fmt.Printf("  '%s' is not a valid repsonse\n\n", atype)
				atype = ""
			}
		}
		if atype == "7" {
			//skip this artifact
			continue
		}

		// We want to shorten the local path - e.g., /C/Users/mitch/cmd/notices.pdf to ./notices.pdf
		// First make sure we do not have a uri (fullpath == _NONE)
		var path string
		if fullpath == _NONE {
			// We have a uri
			path = uri
			fullpath = uri
		} else {
			path = getAbridgedFilePath(fullpath)
		}
		fmt.Println()
		fmt.Println("|--------------------------------------------------")
		fmt.Printf("| %sArtifact%s: %s%s%s\n", _WHITE_FG, _COLOR_END, _CYAN_FG, path, _COLOR_END)
		const padding = 0
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '.', tabwriter.Debug)
		fmt.Fprintf(w, "\t%s\t%s\n", " -----------", "-------------------------------------")
		fmt.Fprintf(w, "\t %s\t %s\n", "Name", name)
		fmt.Fprintf(w, "\t %s\t %s\n", "UUID", uuid)
		fmt.Fprintf(w, "\t %s\t %s\n", "Label", label)
		fmt.Fprintf(w, "\t %s\t %s\n", "Type", atype)
		fmt.Fprintf(w, "\t %s\t %s\n", "Checksum", checksum)
		////fmt.Fprintf(w, "\t %s\t %s\n", "URI", uri)
		fmt.Fprintf(w, "\t %s\t %s\n", "Full Path", fullpath)
		fmt.Fprintf(w, "\n")
		w.Flush()

		var artifact ArtifactRecord
		artifact.Name = name
		artifact.UUID = uuid
		artifact.Alias = label
		artifact.Label = label
		artifact.Checksum = checksum
		artifact.OpenChain = "false"
		artifact.ContentType = atype
		artifact._path = fullpath

		AddArtifactToDB(artifact)
	}
}

func artifactRequest() {
	// if there are no other arguments display help.
	if len(os.Args[1:]) == 1 {
		fmt.Println(_ARTIFACT_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--add", "-a":

	case "--help", "-help", "-h":
		// Display help
		fmt.Println(_ARTIFACT_HELP_CONTENT)
	default:
		fmt.Printf("'%s' is not a valid argument for the '%s %s' command\n", os.Args[2], filepath.Base(os.Args[0]), os.Args[1])
	}
}

func aliasRequest() {
	// if there are no other arguments display help.
	if len(os.Args[1:]) == 1 {
		fmt.Println(_ALIAS_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--get", "-g":
		value, err := getAlias(os.Args[3])
		if checkAndReportError(err) {
			return
		}
		fmt.Println(value)
	case "--help", "-help", "-h":
		// Display help
		fmt.Println(_ALIAS_HELP_CONTENT)
	case "--list", "-l":
		displayAliases()
	case "--set", "-s", "--create":
		if len(os.Args[3:]) < 2 {
			displayErrorMsg("Expecting additional arguments: --set <alias_name>  <value> ")
			return
		}
		alias := os.Args[3]
		value := os.Args[4]
		r, err := regexp.Compile("^[A-Za-z0-9_][A-Za-z0-9._-]*$")
		if checkAndReportError(err) {
			return
		}

		if len(alias) <= _ALIAS_LENGTH && r.MatchString(alias) {
			// the alias of proper length and is syntactically correct.
			setAlias(alias, value)
		} else {
			// the alias is NOT syntactically correct.
			displayErrorMsg(fmt.Sprintf("'%s' is not syntactically correct", os.Args[3]))
			fmt.Println("Aliases must:")
			fmt.Printf("  be %d characters or less in length; and\n", _ALIAS_LENGTH)
			fmt.Println("  begin with an alphanumeric or '_' character;")
			fmt.Println("  followed by a combination of alphanumeric, '_', or '.' characters.")
		}
	default:
		fmt.Printf("'%s' is not a valid argument for the '%s %s' command\n", os.Args[2], filepath.Base(os.Args[0]), os.Args[1])
	}
}

// config command
func configRequest() {
	// if there are no other arguments then display help.
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_CONFIG_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--alias", "-a":
		if len(os.Args[3:]) < 2 {
			fmt.Println("Error: Expecting additional arguments.")
			return
		}
		switch strings.ToLower(os.Args[3]) {
		case "--set", "-s":
			if len(os.Args[4:]) < 2 {
				fmt.Println("Error: Expecting additional argument: <value>")
				return
			}
			setAlias(os.Args[4], os.Args[5])
		case "--get", "-g":
			value, err := getAlias(os.Args[4])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(value)
			}
		case "--list", "-l":
			fmt.Println("Lists all aliases - Not implemented yet.")
		default:
			fmt.Printf("'%s' is not a valid argument for %s %s %s\n", os.Args[3], filepath.Base(os.Args[0]), os.Args[1], os.Args[2])
		}
	case "--help", "-help", "-h":
		// Display help
		fmt.Println(_CONFIG_HELP_CONTENT)
	case "--list", "-l":
		fmt.Println()
		fmt.Println(" -------------")
		fmt.Println(" --- Local ---")
		fmt.Println(" -------------")
		displayLocalConfigData()
		fmt.Println(" --------------")
		fmt.Println(" --- Global ---")
		fmt.Println(" --------------")
		displayGlobalConfigData()
	case "--local":
		if len(os.Args[2:]) == 1 {
			// only --local option, no third option
			// i.e., config --local
			// This is non-valid option
			fmt.Println(_CONFIG_HELP_CONTENT)
			return
		}
		// Have at least arguments.
		switch os.Args[3] {
		case "--list", "-l":
			displayLocalConfigData()
		default:
			if len(os.Args[1:]) == 4 {
				// Request to assign a field (os.Args[3]) a new value (os.Args[4]).
				// Illegal fields will be detect in the setLocalConfigValue() routine.
				setLocalConfigValue(os.Args[3], os.Args[4])
			} else {
				// need 4 arguments but num = 3 or > 4.
				// display help
				fmt.Println(_CONFIG_HELP_CONTENT)
			}
		}
	case "--global":
		if len(os.Args[2:]) == 1 {
			// only have --global option, expected a third argument
			//   i.e., sparts config --global  <- this is not valid option
			fmt.Println(_CONFIG_HELP_CONTENT)
			return
		}
		// We have at least three arguments.
		switch os.Args[3] {
		case "--list":
			displayGlobalConfigData()
		case _ATLAS_ADDRESS_KEY,
			_USER_NAME_KEY,
			_USER_EMAIL_KEY:
			setGlobalConfigValue(os.Args[3], os.Args[4])
		default:
			fmt.Printf("  '%s' is not a validate global configuration value\n", os.Args[3])
		}
	default:
		fmt.Printf("  '%s' is not a valid config option.\n", os.Args[2])
	}
}

func compareRequest() {
	var artifactList1, artifactList2 []ArtifactRecord
	var artifactName1, artifactName2 string = "111", "222"
	var listTitle1, listTitle2 string
	var repoCount int
	var err error

	// At least one additional argument.
	//scanner := bufio.NewScanner(os.Stdin)
	repoCount = 0
	for i := 2; i < len(os.Args) && repoCount < 2; i++ {
		// Check if "help" is the second argument.
		artifact_arg := os.Args[i]
		switch strings.ToLower(artifact_arg) {
		case "--help", "-help", "-h":
			fmt.Println(_COMPARE_HELP_CONTENT)
			return
		case "--dir":
			////fmt.Println (len (os.Args[:i]))
			if len(os.Args[i:]) == 1 {
				displayErrorMsg(fmt.Sprintf("Missing next argument: a directory was expected for argument %d", i))
				return // we are done. exit.
			} else {
				directory := os.Args[i+1]
				i++
				if !isDirectory(directory) {
					displayErrorMsg(fmt.Sprintf("Argument %d: '%s' is not a directory", i, directory))
					return // we are done. exit.
				}
				switch repoCount {
				case 0:
					artifactList1, _ = createEnvelopeFromDirectory(directory)
					artifactName1 = directory
					listTitle1 = " Directory**"
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					artifactList2, _ = createEnvelopeFromDirectory(directory)
					artifactName2 = directory
					listTitle2 = " Directory++"
					repoCount++ // we can accept up to two repositories (director and/or part)
				}
				continue
			}
		case "--part":
			if len(os.Args[i:]) == 1 {
				displayErrorMsg(fmt.Sprintf("Missing next argument: a directory was expected for argument %d", i))
				return // we are done. exit.
			} else {
				part_uuid := os.Args[i+1]
				i++
				if !isValidUUID(part_uuid) {
					displayErrorMsg(fmt.Sprintf("'%s' is not valid uuid", part_uuid))
					return // we are done. exit.
				}
				switch repoCount {
				case 0:
					artifactList1, err = getPartArtifacts(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactName1, _ = getAliasUsingValue(part_uuid)
					if len(artifactName1) < 1 {
						artifactName1 = trimUUID(part_uuid, 5)
					}
					listTitle1 = "   Ledger**"
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					artifactList2, err = getPartArtifacts(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactName2, _ = getAliasUsingValue(part_uuid)
					if len(artifactName2) < 1 {
						artifactName2 = trimUUID(part_uuid, 5)
					}
					listTitle2 = "   Ledger++"
					repoCount++ // we can accept up to two repositories (director and/or part)
				}
				continue
			}
		default:
			displayErrorMsg(fmt.Sprintf("'%s' is not a valid argument.\n", os.Args[i]))
			return // we are done. exit.
		} // switch strings.ToLower(artifact_arg)
	} // for i :=

	if repoCount < 2 { // make sure we have two repositories to compare.

		displayErrorMsg(fmt.Sprintf("Missing two repositories to compare. Try: %s %s --help",
			filepath.Base(os.Args[0]), os.Args[1]))
		return
	}
	// check if any artifacts to display
	if len(artifactList1) == 0 {
		fmt.Printf("No artifacts are contained within repo %s\n", artifactName1)
		return
	}
	// check if any artifacts to display.
	if len(artifactList2) == 0 {
		fmt.Printf("No artifacts are contained within repo %s\n", artifactName2)
		return
	}
	// Display comparison table
	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	const noMatchStr = "   -"
	fmt.Println()
	fmt.Println(" 	Comparing ...")
	fmt.Fprintf(w, " \t %s \t%s \t%s \t\n", "----------------------", " -----------", " -----------")
	//fmt.Fprintf(w, " \t %s \t%s\t%s\t\n", "       Artifacts", " Directory*", "  Ledger*")
	fmt.Fprintf(w, " \t %s\t%s \t%s\t\n", "       Artifacts", listTitle1, listTitle2)
	fmt.Fprintf(w, " \t %s\t%s \t%s \t\n", "----------------------", " -----------", " -----------")
	for i := 0; i < len(artifactList1); i++ {
		for k := 0; k < len(artifactList2); k++ {
			// check that it is not the envelope container
			if artifactList1[i].ContentType == "envelope" || artifactList2[k].ContentType == "envelope" {
				if artifactList1[i].ContentType == _ENVELOPE_TYPE {
					artifactList1[i]._verified = true
				}
				if artifactList2[k].ContentType == _ENVELOPE_TYPE {
					artifactList2[k]._verified = true
				}
				continue
			}
			// See if we have a match (that is not the main envelope)
			if artifactList1[i].Checksum == artifactList2[k].Checksum {
				// we have a match
				fmt.Fprintf(w, " \t %s \t  %s \t  %s\t\n", artifactList1[i].Name, trimUUID(artifactList1[i].Checksum, 5), trimUUID(artifactList2[k].Checksum, 5))
				////fmt.Fprintf(w, "\t  %s\t %s\t %s\t %s\n", id, namea, artifacts[i].Type, path)
				artifactList1[i]._verified = true
				artifactList2[k]._verified = true
			}
		}
	}
	// Now run through the first list to see if any unverified.
	for i := 0; i < len(artifactList1); i++ {
		if !artifactList1[i]._verified {
			fmt.Fprintf(w, " \t %s \t  %s \t  %s\t\n", artifactList1[i].Name, trimUUID(artifactList1[i].Checksum, 5), noMatchStr)
		}
	}

	// Now run through the second list to see if any unmatched.
	for k := 0; k < len(artifactList2); k++ {
		if !artifactList2[k]._verified {
			////id_2 := part_list_2[k].Checksum
			////id_2 = id_2[:5]
			fmt.Fprintf(w, " \t %s \t  %s \t  %s\t\n", artifactList2[k].Name, noMatchStr, trimUUID(artifactList2[k].Checksum, 5))
		}
	}
	// Write out comparison table.
	fmt.Fprintf(w, " \t %s \t%s\t%s\t\n", "----------------------", " -----------", " -----------")
	w.Flush()
	fmt.Printf("   **%s%s%s List\n", _CYAN_FG, artifactName1, _COLOR_END)
	fmt.Printf("   ++%s%s%s List\n", _CYAN_FG, artifactName2, _COLOR_END)
	fmt.Println()
}

// sparts delete ...
func deleteRequest() {
	// the only possible option is --help
	if len(os.Args[2:]) >= 1 {
		// Display help
		fmt.Println(_DELETE_HELP_CONTENT)
		return
	}
	// no arguments - which is expected for delete
	// Proceed to delete
	// Read from standard input.
	scanner := bufio.NewScanner(os.Stdin)
	answer := ""

	fmt.Println("  Are you sure you want to DELLETE ALL 'sparts' data for the working directory (y/n)?")
	fmt.Print("  > ")
	scanner.Scan()
	answer = strings.ToLower(scanner.Text())
	if answer == "y" || answer == "yes" {
		// They said yes - let's double check
		fmt.Println("  Really? Are you sure? - just double checking (y/n)")
		fmt.Print("  > ")
		scanner.Scan()
		answer = strings.ToLower(scanner.Text())
		if answer == "y" || answer == "yes" {
			fmt.Println("  deleting sparts workspace data.....")
			spartsDir, err := getSpartsDirectory()
			if err != nil {
				fmt.Println(err)
				return
			}
			// remove the sparts directory
			err = os.RemoveAll(spartsDir)
			if err == nil {
				// Succefully delete director
				fmt.Printf("The %s directoty has been successfully deleted.\n", filepath.Base(os.Args[0]))
				return
			} else {
				// Directory deletion failed.
				// Try to guest based on error
				errorMsg := err.Error()
				if _DEBUG_DISPLAY_ON {
					displayErrorMsg(errorMsg)
				}
				if strings.Contains(errorMsg, _LOCAL_DB_FILE) {
					displayErrorMsg(fmt.Sprintf("Can't delete %s directory because '%s' is being used by another process",
						_SPARTS_DIRECTORY, _LOCAL_DB_FILE))
				} else if strings.Contains(errorMsg, _LOCAL_CONFIG_FILE) {
					displayErrorMsg(fmt.Sprintf("Can't delete %s directory because '%s' is being used by another process",
						_SPARTS_DIRECTORY, _LOCAL_CONFIG_FILE))
				} else {
					displayErrorMsg(errorMsg)
				}
				fmt.Println("delete request has been CANCELLED.")
				return
			}
		}
	}
	// They decided not to delete the work space
	fmt.Println("  delete request has been CANCELLED.")
}

// sparts dir command
func dirRequest() {
	// the only possible option is --help
	// Anything else we will also display the help contents
	if len(os.Args[2:]) >= 1 {
		// Display help
		fmt.Println(_DIRECTORY_HELP_CONTENT)
		return
	}
	// proceed to obtain and display directory
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("  sparts directory:")
	fmt.Println("     ", spartsDir)

	globalConfigFile, err := getGlobalConfigFile()
	//// Clean up string on Windows replace '\' with '/'
	////globalConfigFile = strings.Replace(globalConfigFile, `\`, `/`, -1)
	if globalConfigFile != "" {
		fmt.Println("  sparts global config file:")
		fmt.Println("     ", globalConfigFile)
	} else {
		fmt.Println("  sparts global config file does not exist.")
	}
}

// sparts envelope
func envelopeRequest() {
	if len(os.Args[1:]) == 1 {
		// No 'part' arguments
		// Display help
		fmt.Println(_ENVELOPE_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--list", "-l":
		if len(os.Args[2:]) > 1 && os.Args[3] == "--all" {
			fmt.Println("  - Not Implemented Yet - list all parts in network")
		} else {
			// displaySupplierParts()
			fmt.Println("  - Not Implemented Yet - list my (supplier) parts")
		}
	case "--help", "-help", "help", "-h":
		// Display help
		fmt.Println(_ENVELOPE_HELP_CONTENT)
	case "--create":
		if len(os.Args[2:]) <= 1 {
			fmt.Println("  Error - Expecting a directory argument.")
			fmt.Printf("  try: %s %s --help\n", filepath.Base(os.Args[0]), os.Args[1])
			return
		}
		// Have at least 1 directory to create envelope from
		// For each directory
		for i := 3; i < len(os.Args); i++ {
			directory := os.Args[i]

			// Check if "help" is the second argument.
			switch strings.ToLower(directory) {
			case "--help", "-help", "-h":
				if i == 2 {
					// --help, -h is the second argument. Print the help text
					fmt.Println(_ENVELOPE_HELP_CONTENT)
					return // we are done
				}
				continue // if --help, -h comes after second argument then ignore it
			}

			if !isDirectory(directory) {
				fmt.Printf(" %s: is not a directory.\n", directory)
				continue // Not a directory - proceed to next directory
			}

			fileCount := getFileCount(directory)
			// Make number of files are within expectation.
			if fileCount > _MAX_FILE_WARNING_COUNT {
				fmt.Printf(" %s directory contains %d files\n", directory, fileCount)
				fmt.Printf(" Do you want to proceed (y/n)?")
				fmt.Print(" > ")
				// Read from standard input.
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				answer := strings.ToLower(scanner.Text())
				if answer == "y" || answer == "yes" {
					fmt.Println(" proceeding ...")
				} else {
					// User does not want to continue. Too many files. Go process next directory
					continue
				}
			}

			artifactList, _ := createEnvelopeFromDirectory(directory)
			artifactsAsJSON, err := createJSONFormat(artifactList)

			if err != nil {
				fmt.Println("Unable to create evelope JSON format. Canceling envelope creation")
				continue
			}
			fmt.Println(artifactsAsJSON)

			partUUID := getLocalConfigValue(_PART_UUID_KEY)
			alias, _ := getAliasUsingValue(partUUID)
			if alias == "" {
				alias = partUUID
			}

			if !getkeyboardYesNoReponse(fmt.Sprintf("The envelope contents is listed above. Do you want to post it \n to the Ledger for part '%s%s%s')?", _GREEN_FG, alias, _COLOR_END)) {
				// They do not wish to proceed.
				continue
			}

			// Returns true if successful
			postEnvelopeToledger(artifactList)

			// Now create relationship between parts and artifacts
			for i := 1; i < len(artifactList); i++ {
				createPartArtifactRelationship(partUUID, artifactList[i].UUID)
				fmt.Println(i)
			}

			//saveEnvelope (artifactsAsJSON)
			// TODO Create TOC.json

			// completed processing directory
			fmt.Println()

		} // for loop - see if another directory
	default:
		//fmt.Println (isHidden(os.Args[2]))
		fmt.Printf("  '%s' is not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Printf("  try: %s %s --help\n", filepath.Base(os.Args[0]), os.Args[1])

	} // end of switch
}

// Dipslays the general help
func helpRequest() {
	// Print help overview
	fmt.Println(_HELP_CONTENT)
}

func initRequest() {
	var spartsDirectory string
	var localConfigFile string
	//var globalConfigFile string

	//see if sub directory specified
	if len(os.Args[2:]) > 1 {
		// Too many arguments specified
		fmt.Println(_INIT_HELP_CONTENT)
		return
	} else if len(os.Args[2:]) == 1 {
		// Additional argument specificied -> sub-diretory
		spartsDirectory = getDirectory() + "/" + os.Args[2] + "/" + _SPARTS_DIRECTORY
		// spartsConfigFile =  getDirectory() +  "/" + os.Args[2] + "/.sparts/config"
		//spartsConfigFile =  getDirectory() +  "/" + os.Args[2] + LOCAL_CONFIG_FILE
	} else {
		// local directory (which is most common)
		spartsDirectory = getDirectory() + "/" + _SPARTS_DIRECTORY
	}

	localConfigFile = spartsDirectory + "/" + _LOCAL_CONFIG_FILE
	//fmt.Println ("spartsDirectory", spartsDirectory)
	//fmt.Println ("spartsConfigFile", spartsConfigFile)

	if isSpartsDirectory(spartsDirectory) {
		// .sparts directory alreadt exists
		// Re-intialize - not sure what that entails yet - just print simple message for now
		fmt.Printf("Reinitialized existing %s working directory:\n", _TOOL_NAME)
		fmt.Println("  ", spartsDirectory)
	} else {
		// Create new directory
		fmt.Printf("Initialized empty %s working directory in:\n", _TOOL_NAME)
		fmt.Println("  ", spartsDirectory)
		createDirectory(spartsDirectory)
		f, err := os.Create(localConfigFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		_, err = f.WriteString(_LOCAL_CONFIG_FILE_CONTENT)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// The following will create the global config file if it does not exist.
	globalConfigFile, err := getGlobalConfigFile()
	if err != nil && globalConfigFile != "" {
		fmt.Println("Created global config file:", globalConfigFile)
	}
	//fmt.Println( usr.HomeDir )

	// Initialize the database
	initializeDB()
}

func partRequest() {
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_PART_HELP_CONTENT)
		return
	}
	switch os.Args[2] {
	case "--create":
		// Let's first see if the supplier uuid is set in the config file.
		// TODO: allow user to enter UUID if not set in local config.
		supplierUUID := getLocalConfigValue("supplier_uuid")

		if supplierUUID == "" {
			fmt.Println("Supplier UUID not set in local the config file.")
			return
		}
		// TODO: Check if UUID is correcrly formatted.

		// Read from standard input.
		scanner := bufio.NewScanner(os.Stdin)
		name, version, label, license, description, url, uuid, checksum := "", "", "", "", "", "", "", ""

		fmt.Println("  -------------------------------------------------------------------------")
		fmt.Println("  required field (*)")

		for name == "" {
			fmt.Print("  name(*): ")
			scanner.Scan()
			name = scanner.Text()
		}

		fmt.Print("  version: ")
		scanner.Scan()
		version = scanner.Text()

		fmt.Print("  label (nickname): ")
		scanner.Scan()
		label = scanner.Text()

		fmt.Print("  licensing: ")
		scanner.Scan()
		license = scanner.Text()

		fmt.Print("  description: ")
		scanner.Scan()
		description = scanner.Text()

		h := sha1.New()
		h.Write([]byte("Some String"))
		bs := h.Sum(nil)
		checksum = fmt.Sprintf("%x", bs)

		fmt.Print("  url: ")
		scanner.Scan()
		url = scanner.Text()

		fmt.Print("  UUID (auto generated if blank): ")
		scanner.Scan()
		uuid = scanner.Text()

		if uuid == "" {
			uuid = getUUID()
			//uuid ="f1c2d3..."
		}

		fmt.Println("  \n  Do you want to proceed to create a part with the following values (y/n)?")
		fmt.Println("  -------------------------------------------------------------------------")
		fmt.Println("\tname      	= " + name)
		fmt.Println("\tversion   	= " + version)
		fmt.Println("\tlabel     	= " + label)
		fmt.Println("\tlicensing 	= " + license)
		fmt.Println("\turl       	= " + url)
		fmt.Println("\tchecksum  	= " + checksum)
		fmt.Println("\tuuid      	= " + uuid)
		fmt.Println("\tdescription 	= " + description)

		fmt.Print("(y/n)> ")
		scanner.Scan()
		confirmation := strings.ToLower(scanner.Text())
		if confirmation == "y" || confirmation == "yes" {
			fmt.Println("submitting part info ....")
			ok, err := createPart(name, version, label, license, description, url, checksum, uuid)
			//result := ""

			if ok == false || err != nil {
				if checkAndReportError(err) {
					return
				}
			} else {
				// Part was successfully registered with the legder
				// Create relationship between part and supplier.
				ok, err := createPartSupplierRelationship(uuid, supplierUUID)
				if ok == true {
					fmt.Println("Create part on ledger was SUCCESSFUL.")
					if getkeyboardYesNoReponse("Would you like to create an alias for this part (y/n)?") {
						alias := getkeyboardReponse("Enter the alias you would like to use?")
						err := setAlias(alias, uuid)
						if err != nil {
							fmt.Printf("Error: Unable to create alias '%s'.\n", alias)
							return
						}
						fmt.Printf("You can Use expression '%s%s' at the command line to reference the part\n", _ALIAS_TOKEN, alias)
						return
					}
					return
				} else {
					fmt.Println("Create part-to-suppiler relationship ledger api call FAILED.")
					if _DEBUG_DISPLAY_ON {
						fmt.Println(err)
					}
					return
				}
			}
		} else {
			fmt.Println("Part creation request has been cancelled.")
		}

	case "--get":
		var part PartRecord
		var partID string
		var err error

		// Any additional arguments ?
		if len(os.Args[3:]) == 0 {
			// no other arguments (e.g., no uuid). Assume local config supplier uuid
			// e.g., sparts part --get
			partID = getLocalConfigValue(_PART_UUID_KEY)
			if partID == _NULL_PART {
				fmt.Println("Part uuid is not assigned in local config file. Try:")
				fmt.Printf("  %s part --get uuid=<uuid>\n", filepath.Base(os.Args[0]))
				fmt.Println("or")
				fmt.Printf("  %s part --set uuid=<uuid>\n", filepath.Base(os.Args[0]))
				fmt.Printf("  %s part --get\n", filepath.Base(os.Args[0]))
				return // we are done. exit func.
			}
		} else { // forth argument exists, should be uuid=<uuid>. Let's check the format
			// next argument is a uuid or alias for a uuid.
			partID = os.Args[3]
			if !isValidUUID(partID) {
				if strings.ToLower(partID) == strings.ToLower(_NULL_PART) {
					fmt.Printf("  '%s' is not acceptable value here\n", partID)
					return // we are done. exit func.
				} else {
					fmt.Printf("  '%s' is not a properly formatted UUID\n", partID)
					return // we are done. exit func.
				}
			}
			// uuid syntax is not properly formated

		} // else
		// part_id holds a properly formated uuid.
		part, err = getPartInfo(partID)
		if err != nil {
			displayErrorMsg(err.Error())
			return // we are done. exit func.
		}

		// No errors. Proceed to print.
		alias, err := getAliasUsingValue(part.UUID)
		if err != nil {
			alias = ""
		}
		// Display part info
		////fmt.Println()

		//fmt.Printf ("  Part Info: %s%s%s\n",_GREEN_FG, alias, _COLOR_END)
		fmt.Println("  -------------------------------------------------------------")
		fmt.Printf("  Part Name   : %s%s%s\n", _GREEN_FG, part.Name, _COLOR_END)
		fmt.Println("  -------------------------------------------------------------")
		if alias != "" {
			//fmt.Printf("  Alias (id=) : %s%s%s\n", _GREEN_FG, alias, _COLOR_END)
			fmt.Printf("  Alias (id=) : %s\n", alias)
		}
		fmt.Println("  Version     :", part.Version)
		fmt.Println("  Label       :", part.Label)
		fmt.Println("  License     :", part.Licensing)
		fmt.Println("  Checksum    :", part.Checksum)
		fmt.Println("  UUID        :", part.UUID)
		//fmt.Println("  URI         :", part.URI)
		fmt.Println("  Description : " + formatDisplayString(part.Description, 60))
		fmt.Println("  -------------------------------------------------------------")

	case "--help", "-help", "help", "-h":
		// Display help
		fmt.Println(_PART_HELP_CONTENT)
	case "--list", "-l":
		if len(os.Args[2:]) > 1 && os.Args[3] == "--all" {
			////if len(os.Args[3:]) == 0 {
			partsList, err := getPartList()
			if checkAndReportError(err) {
				return
			}
			/////displayParts(partList)
			partItemList := []PartItemRecord{}
			var partItem PartItemRecord
			for k := range partsList {
				partItem.PartUUID = partsList[k].UUID
				partItemList = append(partItemList, partItem)
			}
			displayParts(partItemList)
		} else {
			// displaySupplierParts()
			fmt.Println("  - Not Implemented Yet - list my (supplier) parts")
			fmt.Printf("  - try: '%s %s --list --all' to list all the parts registered with the ledger\n", filepath.Base(os.Args[0]), os.Args[1])
		}
	case "--set":
		var uuid string
		if len(os.Args[3:]) >= 1 {
			uuid = strings.ToLower(os.Args[3])
			if isValidUUID(uuid) || strings.ToLower(uuid) == strings.ToLower(_NULL_PART) {
				// we are good.
				setLocalConfigValue(_PART_UUID_KEY, uuid)
				return // we are done.
			}
		}
		// the uuid is not valid
		fmt.Printf("'%s' is not a valid uuid\n", uuid)
		return
		/*******
		var idStr []string
		uuidStrValid := true // Intially assume true. Set to false once we learn not true.
		// See if next argument uuid=xxx exists
		if len(os.Args[3:]) >= 1 && strings.Contains(strings.ToLower(os.Args[3]), "uuid=") {
			idStr = strings.Split(os.Args[3], "=")
			//if isValidUUID(idStr[0]) {
			if len(idStr) == 2 {
				if isValidUUID(idStr[1]) || strings.ToLower(idStr[1]) == strings.ToLower(_NULL_PART) {
					// we are good. idStr[1] holds the uuid.
					setLocalConfigValue(_PART_UUID_KEY, strings.ToLower(idStr[1]))
					return // we are done.
				} else {
					// uuid is not valid format.
					uuidStrValid = false
				}
			}
		}
		// If we get this far then --set argument is not valid.
		if uuidStrValid {
			fmt.Printf("Format is not valid. Expecting: \n")
			fmt.Printf("   %s --set uuid=<uuid>  \n", filepath.Base(os.Args[0]))
			fmt.Printf("   %s --set uuid=%s    .\n",
				filepath.Base(os.Args[0]), strings.ToLower(_NULL_PART))
		} else {
			fmt.Printf("UUID '%s' is not properly formatted\n", idStr[1])
		}
		****/
	default:
		fmt.Printf("%s: not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Println(_PART_HELP_CONTENT)
	}
}

func pingRequest() {
	//ok, err := pingServer(_ATLAS)
	ok, err := pingServer(_LEDGER)
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		fmt.Println("ping was successful")
	} else {
		fmt.Println("ping failed")
	}
}

func removeRequest() {
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_REMOVE_HELP_CONTENT)
		return
	}
	// At least one argument.
	switch strings.ToLower(os.Args[2]) {
	case "-h", "--help":
		// Display help
		fmt.Println(_REMOVE_HELP_CONTENT)
	default:
		firsttime := true
		for i := 2; i <= len(os.Args[1:]); i++ {
			if _, err := strconv.Atoi(os.Args[i]); err != nil {
				fmt.Printf("  %s is not an integer (id)\n", os.Args[i])
				if firsttime {
					firsttime = false
					fmt.Println("  Usage: sparts remove <id> ...")
					fmt.Println("   e.g.,  sparts remove 4")
					fmt.Println("  Obtain id from: 'sparts status'")
					fmt.Println()
				}
				continue
			}
			// get the artifact record for specified id
			artifact := getArtifactFromDB("id", os.Args[i])
			// Delete aritfact record for id. Give error if id not found.
			if !deleteArtifactFromDB(artifact) {
				fmt.Println("  artifact not found for id =", os.Args[i])
			}
		}
	}
}

//sparts status
func statusRequest() {
	var id, name, path string

	artifacts, err := getArtifactListDB()
	if err != nil {
		fmt.Println("  fatal: sparts working database not accessible.")
		os.Exit(_DB_ACCESSS_ERROR) // exit program
	}

	fmt.Println()
	supplychain := getLocalConfigValue("supply_chain")
	fmt.Println("Network: ", supplychain)
	////fmt.Println("Staged, waiting for commit:")

	part_uuid := getLocalConfigValue(_PART_UUID_KEY)
	partAlias, err := getAliasUsingValue(part_uuid)
	if partAlias != "" && err == nil {
		part_uuid = partAlias + " = " + part_uuid
	}

	if part_uuid == _NULL_PART {
		part_uuid = _RED_FG + part_uuid + _COLOR_END
	} else {
		part_uuid = _GREEN_FG + part_uuid + _COLOR_END
	}

	if len(artifacts) == 0 {
		// nothing waiting to commit
		fmt.Println("nothing to commit (use 'spart add' to stage artifact for commit)")
		return // we're done
	}
	// At least one items pending a commit
	fmt.Println("Staged artifacts to be committed:")
	fmt.Println(" (use 'spart remove id1 id2 ...' to unstage)")
	fmt.Println()
	fmt.Println(" |---------------------------------------------------------------")
	fmt.Println(" | ")
	/******
	partAlias, err := getAliasUsingValue(part_uuid)
	fmt.Println("partAlias", partAlias)
	if partAlias != "" && err == nil {
		part_uuid = partAlias
	}
	*****/
	fmt.Println(" | part :", part_uuid)
	fmt.Println(" | ")
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ',
		tabwriter.Debug)
	fmt.Fprintf(w, "\t%s\t%s\t %s\t %s \n", " ----", " ----------", "------", "------------------------")
	fmt.Fprintf(w, "\t%s\t%s\t %s\t %s\n", "  Id", "  Name  ", " Type", "  File Path* or URI")
	fmt.Fprintf(w, "\t%s\t%s\t %s\t %s\n", " ----", " ----------", "------", "------------------------")

	for i := range artifacts {
		id = strconv.Itoa(artifacts[i]._ID)
		name = artifacts[i].Name
		if isPathURL(artifacts[i]._path) {
			// path is a url link
			path = artifacts[i]._path
		} else {
			// It is a file, convert path relative to the .sparts working directory.
			////fmt.Println("Path:", artifacts[i]._path)
			path = getAbridgedFilePath(artifacts[i]._path)
		}
		fmt.Fprintf(w, "\t  %s\t %s\t %s\t %s\n", id, name, artifacts[i].ContentType, path)
	}

	//fmt.Fprintf(w, "\n")

	w.Flush()
	fmt.Println("   ---")
	fmt.Println(" * File paths are relative to the sparts working directory.")
}

func supplierRequest() {
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_SUPPLIER_HELP_CONTENT)
		return
	}
	//fmt.Println ("num", len(os.Args[1:]))
	switch os.Args[2] {
	case "--list", "-l":
		displaySupplierList()
	case "--help", "-help", "help", "-h":
		// Display help
		fmt.Println(_SUPPLIER_HELP_CONTENT)
	case "--get":
		if len(os.Args[3:]) == 0 {
			// no other arguments (e.g., no uuid). Assume local config supplier uuid
			displaySupplier(getLocalConfigValue("supplier_uuid"))
		} else {
			// next argument should be uuid.
			// TODO - check if uuid syntax is correct.
			displaySupplier(os.Args[3])
		}
	case "--create", "-c":
		// for each additional arg
		name, short_id, url := " ", " ", " "
		for i := 1; i <= len(os.Args[3:]); i++ {
			//fmt.Println(os.Args[2+i])
			arg := strings.Split(os.Args[2+i], "=")
			if len(arg) != 2 {
				fmt.Println("  Error with arguments. Expecting: name=value")
				fmt.Println("  e.g., name=[...] short_id=[...] url=[...]")
				return
			}
			switch arg[0] {
			case "name":
				name = arg[1]
			case "short_id":
				short_id = arg[1]
			case "url":
				url = arg[1]
			default:
				fmt.Printf("  error - '%s'is not a valid argument for %s\n", os.Args[2+i], os.Args[2])
				fmt.Printf("  Expecting name=[...] short_id=[...] and url=[...]\n")
				return
			}
		}
		if name == " " {
			fmt.Printf("  error - expecting at least argument name=[supplier name] \n")
		} else {
			// send request to ledger to create new supplier
			uuid := createSupplier(name, short_id, "", "abc123", url)
			if uuid == "" {
				fmt.Printf("  Not able to create new supplier: '%s'\n", name)
			} else {
				fmt.Printf("  new supplier '%s' has uuid = %s\n", name, uuid)
			}
			//fmt.Printf ("  do it name=%s, short_id=%s, url=%s\n", name, short_id, url)
		}
	default:
		fmt.Printf("%s: not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Println(_SUPPLIER_HELP_CONTENT)
	}
}

func synchRequest() {

	fmt.Println("  Network: ", getLocalConfigValue(_SUPPLY_CHAIN_NETWORK_KEY))
	ok, err := pingServer(_LEDGER)
	//ok, err := pingServer(_ATLAS)
	//ok := false
	if ok {
		// things are good.
		fmt.Println("  Current ledger node is ACTIVE at address:", getLocalConfigValue(_LEDGER_ADDRESS_KEY))
		return // done.
	}

	// could not successful access the current ledger node. Proceed to check for other nodes.
	fmt.Println("  Default ledger node is NOT ACTIVE:", getLocalConfigValue(_LEDGER_ADDRESS_KEY))
	fmt.Println("  Searching for new ledger node .....")

	// Obtain current list of available ledger nodes from look up directory (atlas)
	nodeList, err := getLedgerNodeList()
	if err != nil {
		fmt.Println(err)
		// Suggest a fix for certain circumstances
		if strings.Contains(err.Error(), "does not exist") {
			fmt.Printf("  You may need to set or update local config value: '%s'\n", _SUPPLY_CHAIN_NETWORK_KEY)
		}
		return
	}
	// Check if list is empty
	if len(nodeList) == 0 {
		fmt.Printf("  The network '%s' has no ledger nodes registered\n", getLocalConfigValue(_SUPPLY_CHAIN_NETWORK_KEY))
		return
	}
	newNodeFound := false
	for _, node := range nodeList {
		if newNodeFound {
			break // from for loop.
		}
		//fmt.Println("  Checking node:", node.APIURL)
		ok, err := pingServer(node.APIURL)
		if err != nil {
			//fmt.Println("   ", err)
		}
		if ok {
			newNodeFound = true
			setLocalConfigValue(_LEDGER_ADDRESS_KEY, node.APIURL)
			fmt.Printf("  Found ACTIVE ledger node at: '%s'\n", node.APIURL)
			fmt.Println("  UPDATING default ledger node in config to:", node.APIURL)
		}
	}
	if newNodeFound == false {
		fmt.Printf("  Not able to locate an active ledger node using the %s directory\n", _ATLAS)
	}
}

// Temporary Used to initally test new commands. This will be removed.
func lfsRequest() {
	//fmt.Println(" Not implemented")

	for i := 1; i < 30; i++ {
		fmt.Println(getUUID())
	}
	return

	file := os.Args[2]
	_, _, _, fileExtension := FilenameDirectorySplit(file)
	sha1, _ := getFileSHA1(file)
	fi, err := os.Stat(file)
	if checkAndReportError(err) {
		return
	}
	fileSize := strconv.FormatInt(fi.Size(), 10)
	fileRepoPath := "_content/" + sha1 + "." + fileSize + fileExtension

	//First we need to set up the authenticate token with the github server.
	context := context.Background()
	// get token from: https://github.com/settings/tokens/new
	// and you need to enter it in the configuration file
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "e841f93ab39d5711e5ae48d917c82b041f0a63de"})
	tokenClient := oauth2.NewClient(context, tokenService)
	githubClient := github.NewClient(tokenClient)

	///fileContent := []byte("This is the content of my file\nand the 2nd line of it")
	fileBytes, err := ioutil.ReadFile(file)

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String("This is my commit message"),
		Content:   fileBytes,
		Branch:    github.String("master"),
		Committer: &github.CommitAuthor{Name: github.String("Mark"), Email: github.String("user@example.com")},
	}
	//_, _, err := githubClient.Repositories.CreateFile(context, "g-snoop", "zephyr-content", "content/README.md", opts)
	_, _, err = githubClient.Repositories.CreateFile(context, "g-snoop", "zephyr-content", fileRepoPath, opts)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// sparts version
func versionRequest() {
	// if there are no other arguments display the version for sparts.
	if len(os.Args[1:]) == 1 {
		// Display version
		fmt.Printf("  %s version: %s\n", filepath.Base(os.Args[0]), _VERSION)
		return
	}
	// there is an additional argument
	switch os.Args[2] {
	case "--help", "-help", "-h":
		// Display help
		fmt.Println(_VERSION_HELP_CONTENT)
	case "--all", "-a":
		fmt.Printf("%s version: %s   data model: %s\n", filepath.Base(os.Args[0]), _VERSION, _DB_Model)
	default:
		fmt.Printf("  '%s' is not a valid version option. Try --help\n", os.Args[2])
	}

}

func seedRequest() {

	fmt.Println(" Not implemented")

}

func testRequest() {

}

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
	"net/url"
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
		fmt.Printf("  tip: use '%s init' to create %s working directory\n", filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
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
	case "focus":
		focusRequest()
	case "help", "-help", "-h", "--help":
		helpRequest()
	case "init":
		initRequest()
	case "network":
		networkRequest()
	case "org":
		orgRequest()
	case "part":
		partRequest()
	case "ping":
		pingRequest()
	case "push":
		pushRequest()
	case "quick":
		quickRequest()
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
	case "tips":
		tipsRequest()
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
	/////scanner := bufio.NewScanner(os.Stdin)
	openChain := _FALSE
	var artifact ArtifactRecord
	var err error
	for i := 2; i < len(os.Args); i++ {
		artifactArg := os.Args[i]

		// Check if "help" is the second argument.
		switch strings.ToLower(artifactArg) {
		case "--help", "-help", "-h":
			if i == 2 {
				// --help, -h is the second argument. Print the help text
				fmt.Println(_ADD_HELP_CONTENT)
				return // we are done
			}
			continue // if --help, -h comes after second argument then ignore it
		case "--url", "-u":
			// We only allow one link argument per line.
			// if number of remaining arguments is 2 OR
			// if --openchain present AND remaining arguments therefore must be 3
			if len(os.Args[2:]) == 2 || (openChain == _TRUE && len(os.Args[2:]) == 3) {
				theURL := os.Args[3]
				// process link
				_, err := url.Parse(theURL)
				if err != nil {
					fmt.Printf("	error:  %s%s%s - invalid link syntax\n", _RED_FG, theURL, _COLOR_END)
					return
				}
				artifact = ArtifactRecord{}
				_, name, _, _ := FilenameDirectorySplit(theURL)
				if strings.HasPrefix(theURL, "https") {
					artifact.Name = "https://.../" + name
				} else {
					artifact.Name = "http://.../" + name
				}
				artifact.UUID = getUUID()
				artifact.Alias = name
				artifact.Label = name
				artifact.ContentType = "url"
				artifact.Checksum, _ = getStringSHA1(theURL)
				artifact.URIList = []URIRecord{} // initalize to the empty list
				artifact._contentPath = theURL
				if openChain == _TRUE {
					artifact.OpenChain = _TRUE
				} else {
					artifact.OpenChain = _FALSE
				}
				artifact._envelopeUUID = _NULL_UUID
				artifact._onLedger = _FALSE
				// Add to database
				err = addArtifactToDB(artifact)
				if err != nil {
					fmt.Printf("	error:  %s%s%s\n", _RED_FG, theURL, _COLOR_END)
				} else {
					fmt.Printf("	adding: %s%s%s\n", _GREEN_FG, theURL, _COLOR_END)
				}
				return
			}
			// else more than one argument. Set flag and proceed to process all the files
			openChain = _TRUE
			continue
		case "--dir":
			var directory string
			// advance argument index to directory name
			i++
			if len(os.Args[i]) > 0 {
				directory = os.Args[i]
			}
			// Make sure it is a directory
			if !isDirectory(directory) {
				displayErrorMsg(fmt.Sprintf("'%s' is not a directory", directory))
				return // exit
			}
			// obtain envelope artifact list from directory (the last argument)
			artifactList, err := createEnvelopeFromDirectory(directory, openChain == _TRUE)
			if len(artifactList) == 0 || err != nil {
				displayErrorMsg(fmt.Sprintf("cannot generate envelope from directory: '%s': %s", directory, err.Error()))
				return // exit
			}
			for i, artifact := range artifactList {
				// Assign the artifact envelope UUID for non-envelope artifacts.
				if i == 0 {
					continue // skip first (top directory record)
				}
				artifact._envelopeUUID = _NULL_UUID
				// Add to database
				artifactName := artifact._envelopePath + artifact.Name
				err := addArtifactToDB(artifact)
				if err != nil {
					fmt.Printf("	error:  %s%s%s\n", _RED_FG, artifactName, _COLOR_END)
				} else {
					fmt.Printf("	adding: %s%s%s\n", _GREEN_FG, artifactName, _COLOR_END)
				}
			}
			continue
		case "--openchain", "-oc":
			// If --openchain is the only argument, we need to
			// assign all the artifacts in staging area to be OpenChain.
			// If there additional argumnets just set the openchain flag to true/yes

			/************
			// First see if this is the only argument
			//if len(os.Args[2:]) == 2 && os.Args[3] == "-all" {
			if len(os.Args[2:]) == 1 {
				fmt.Println("Only OpenChain")
				// if more than 5 non-openchain ask - do you want to sent them all ?
				// set all
				artifacts, err := getArtifactListDB()
				if err != nil {
					displayErrorMsg("sparts working database not accessible.")
					return
				}
				for _, artifact := range artifacts {
					fmt.Println(artifact.Name, artifact.OpenChain)
					if artifact.OpenChain == _FALSE {
						artifact.OpenChain = _TRUE
						addArtifactToDB(artifact)
					}
				}
				return
			}
			****************/
			// else more than one argument. Set flag and proceed to process all the files
			openChain = _TRUE
			continue
		default:
			fmt.Println("artifactArg", artifactArg)
			artifact, err = createArtifactFromFile(artifactArg)
			if err != nil {
				fmt.Printf("	error:  %s%s%s - %s\n", _RED_FG, artifactArg, _COLOR_END, err)
				continue // go process next artifact
			}
			if openChain == _TRUE {
				artifact.OpenChain = _TRUE
			}
			// Add to database
			err = addArtifactToDB(artifact)
			if err != nil {
				fmt.Printf("	error:  %s%s%s\n", _RED_FG, artifactArg, _COLOR_END)
			} else {
				fmt.Printf("	adding: %s%s%s\n", _GREEN_FG, artifactArg, _COLOR_END)
			}
			// Go process next artifact (if any remaining)
		}
	} // For loop
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

// handle: 'sparts alias' command
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

// handle: 'sparts config' command
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
		displayLocalConfigData()
		fmt.Println()
		displayGlobalConfigData()
		fmt.Println()
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
			fmt.Println()
			displayLocalConfigData()
			fmt.Println()
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
			fmt.Println()
			displayGlobalConfigData()
			fmt.Println()
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

// handle: 'sparts compare' command
func compareRequest() {
	var artifactList1, artifactList2 []ArtifactRecord
	var artifactSetName1, artifactSetName2 string = "111", "222"
	var listTitle1, listTitle2 string
	var repoCount int
	var err error

	if len(os.Args[1:]) == 1 {
		fmt.Println(_COMPARE_HELP_CONTENT)
		return
	}

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
					artifactList1, _ = createEnvelopeFromDirectory(directory, false)
					artifactSetName1 = directory
					listTitle1 = "  Directory**  "
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					artifactList2, _ = createEnvelopeFromDirectory(directory, false)
					artifactSetName2 = directory
					listTitle2 = "  Directory++  "
					repoCount++ // we can accept up to two repositories (director and/or part)
				}
				continue
			}
		case "--env":
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
					//artifactList1, err = getPartArtifacts(part_uuid)
					artifactList1, err = getEnvelopeArtifactsFromLedger(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactSetName1, _ = getAliasUsingValue(part_uuid)
					if len(artifactSetName1) < 1 {
						artifactSetName1 = trimUUID(part_uuid, 5)
					}
					listTitle1 = "  Ledger**  "
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					// artifactList2, err = getPartArtifacts(part_uuid)
					artifactList2, err = getEnvelopeArtifactsFromLedger(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactSetName2, _ = getAliasUsingValue(part_uuid)
					if len(artifactSetName2) < 1 {
						artifactSetName2 = trimUUID(part_uuid, 5)
					}
					listTitle2 = "  Ledger++  "
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
		fmt.Printf("No artifacts are contained within repo %s\n", artifactSetName1)
		return
	}
	// check if any artifacts to display.
	if len(artifactList2) == 0 {
		fmt.Printf("No artifacts are contained within repo %s\n", artifactSetName2)
		return
	}

	err = displayListComparison(artifactList1, artifactList2, listTitle1, listTitle2, artifactSetName1, artifactSetName2)
}

// handle: 'sparts compare' command
func compareRequest2() {
	var artifactList1, artifactList2 []ArtifactRecord
	var artifactName1, artifactName2 string = "111", "222"
	var listTitle1, listTitle2 string
	var repoCount int
	var err error

	if len(os.Args[1:]) == 1 {
		fmt.Println(_COMPARE_HELP_CONTENT)
		return
	}

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
					artifactList1, _ = createEnvelopeFromDirectory(directory, false)
					artifactName1 = directory
					listTitle1 = "Directory**"
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					artifactList2, _ = createEnvelopeFromDirectory(directory, false)
					artifactName2 = directory
					listTitle2 = "Directory++"
					repoCount++ // we can accept up to two repositories (director and/or part)
				}
				continue
			}
		case "--env":
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
					//artifactList1, err = getPartArtifacts(part_uuid)
					artifactList1, err = getEnvelopeArtifactsFromLedger(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactName1, _ = getAliasUsingValue(part_uuid)
					if len(artifactName1) < 1 {
						artifactName1 = trimUUID(part_uuid, 5)
					}
					listTitle1 = "  Ledger**  "
					repoCount++ // we can accept up to two repositories (director and/or part)
				case 1:
					// artifactList2, err = getPartArtifacts(part_uuid)
					artifactList2, err = getEnvelopeArtifactsFromLedger(part_uuid)
					if err != nil {
						displayErrorMsg(err.Error())
						return // we are done. exit.
					}
					artifactName2, _ = getAliasUsingValue(part_uuid)
					if len(artifactName2) < 1 {
						artifactName2 = trimUUID(part_uuid, 5)
					}
					listTitle2 = "  Ledger++  "
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
	const equalStr = "="
	const notEqualStr = "X"
	const noMatchStr = "     -     "

	/*******************************************
	const padding = 0
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	// writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Println()

	//fmt.Println(" 	Comparing ...")
	fmt.Fprintf(w, " \t %s\t%s\t%s\t%s \t\n", "----------------------", " -----------", "  ", " -----------")
	//fmt.Fprintf(w, " \t %s \t%s\t%s\t\n", "       Artifacts", " Directory*", "  Ledger*")
	fmt.Fprintf(w, " \t %s\t%s\t%s\t%s\t\n", "       Artifacts", listTitle1, "  ", listTitle2)
	//fmt.Fprintf(w, " \t %s\t%s %s %s \t\n", "       Artifacts", listTitle1, "  ", listTitle2)
	fmt.Fprintf(w, " \t %s\t%s\t%s\t%s \t\n", "----------------------", " -----------", "  ", " -----------")
	//fmt.Fprintf(w, " \t %s\t%s %s %s \t\n", "----------------------", " -----------", "  ", " -----------")
	*********************************************/
	fmt.Println()
	const padding = 0
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	header := []string{"   Artifacts   ", listTitle1, "", listTitle2}
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator
	PrintRow(w, PaintRowUniformly(CyanText, header))
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator

	var colors []Color

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
				colors = []Color{DefaultText, DefaultText, GreenText, DefaultText}
				////fmt.Fprintf(w, "\t  %s\t %s\t %s\t %s\n", id, namea, artifacts[i].Type, path)
				//fmt.Fprintf(w, " \t %s \t  %s \t  %s \t  %s\t\n", artifactList1[i].Name, trimUUID(artifactList1[i].Checksum, 5), equalStr, trimUUID(artifactList2[k].Checksum, 5))
				PrintRow(w, PaintRow(colors, []string{
					artifactList1[i].Name,
					"  " + trimUUID(artifactList1[i].Checksum, 5),
					equalStr,
					"  " + trimUUID(artifactList2[k].Checksum, 5)}))
				////fmt.Fprintf(w, " \t %s \t  %s  %s  %s\t\n", artifactList1[i].Name, trimUUID(artifactList1[i].Checksum, 5), equalStr, trimUUID(artifactList2[k].Checksum, 5))
				artifactList1[i]._verified = true
				artifactList2[k]._verified = true
			}
		}
	}
	// Now run through the first list to see if any unverified.
	for i := 0; i < len(artifactList1); i++ {
		if !artifactList1[i]._verified {
			colors = []Color{DefaultText, DefaultText, RedText, RedText}
			PrintRow(w, PaintRow(colors, []string{
				artifactList1[i].Name,
				"  " + trimUUID(artifactList1[i].Checksum, 5),
				notEqualStr,
				noMatchStr}))
			//fmt.Fprintf(w, " \t %s \t  %s \t  %s \t  %s\t\n", artifactList1[i].Name, trimUUID(artifactList1[i].Checksum, 5), notEqualStr, noMatchStr)
		}
	}

	// Now run through the second list to see if any unmatched.
	for k := 0; k < len(artifactList2); k++ {
		if !artifactList2[k]._verified {
			////id_2 := part_list_2[k].Checksum
			////id_2 = id_2[:5]
			//fmt.Fprintf(w, " \t %s \t  %s \t  %s \t  %s\t\n", artifactList2[k].Name2, noMatchStr, notEqualStr, trimUUID(artifactList2[k].Checksum, 5))
			//fmt.Fprintf(w, " \t %s \t  %s \t  %s \t  %s\t\n", artifactList2[k].Name2, noMatchStr, notEqualStr, trimUUID(artifactList2[k].Checksum, 5))
			colors = []Color{DefaultText, RedText, RedText, DefaultText}
			PrintRow(w, PaintRow(colors, []string{
				artifactList2[k].Name2,
				noMatchStr,
				notEqualStr,
				"  " + trimUUID(artifactList2[k].Checksum, 5)}))
		}
	}
	// Write out comparison table.
	//fmt.Fprintf(w, " \t %s \t%s\t%s\t\n", "----------------------", " -----------", " -----------")
	//fmt.Fprintf(w, " \t %s\t%s\t%s\t%s \t\n", "----------------------", " -----------", "  ", " -----------")
	//fmt.Fprintf(w, " \t %s\t%s\t%s\t%s \t\n", "----------------------", " -----------", "  ", " -----------")
	PrintRow(w, PaintRowUniformly(DefaultText, AnonymizeRow(header)))
	w.Flush()
	fmt.Printf("  **%s%s%s List\n", _CYAN_FG, artifactName1, _COLOR_END)
	fmt.Printf("  ++%s%s%s List\n", _CYAN_FG, artifactName2, _COLOR_END)
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

// handle: 'sparts dir' command
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

func focusRequest() {
	if len(os.Args) == 2 {
		// No arguments. Expecting: part, envelope, both, none. Display help
		fmt.Println(_FOCUS_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--both":
		setLocalConfigValue(_FOCUS_KEY, _BOTH_FOCUS)
	case "--envelope":
		setLocalConfigValue(_FOCUS_KEY, _ENVELOPE_FOCUS)
		if len(os.Args) == 4 && isValidUUID(os.Args[3]) {
			//TODO: check that uuid is in fact an envelop in the db.
			setLocalConfigValue(_ENVELOPE_KEY, os.Args[3])
		}
	case "--help", "-h", "help":
		fmt.Println(_FOCUS_HELP_CONTENT)
	case "--none":
		setLocalConfigValue(_FOCUS_KEY, _NO_FOCUS)
	case "--part", "-p":
		setLocalConfigValue(_FOCUS_KEY, _PART_FOCUS)
		if len(os.Args) == 4 && isValidUUID(os.Args[3]) {
			//TODO: check that uuid is in fact an envelop in the db.
			setLocalConfigValue(_ENVELOPE_KEY, os.Args[3])
		}
	}
}

// handle: 'sparts envelope' command
func envelopeRequest() {
	if len(os.Args[1:]) == 1 {
		// No 'part' arguments
		// Display help
		fmt.Println(_ENVELOPE_HELP_CONTENT)
		return
	}
	switch strings.ToLower(os.Args[2]) {
	case "--list", "-l":
		displayEnvelopeList()
	case "--help", "help", "-h":
		// Display help
		fmt.Println(_ENVELOPE_HELP_CONTENT)
	case "--set", "-s":
		if len(os.Args) == 4 {
			uuid := os.Args[3]
			if isValidUUID(uuid) {
				//TODO: check if valid envelope uuid.
				err := setLocalConfigValue(_ENVELOPE_KEY, uuid)
				if err != nil {
					displayErrorMsg(err.Error())
					return
				}
				return
			} else {
				displayErrorMsg(fmt.Sprintf("uuid '%s' is not in a valid format.", uuid))
				return
			}
		} else {
			displayErrorMsg("Incorrect number of arguments (expecting four). See sparts --envelope --help")
			return
		}
	case "--create":
		var err error
		var envelope ArtifactRecord
		// TODO: check that no extra Args are provided,
		// last argument should be envelope name
		// TODO: check syntax of envelope name
		lastArg := os.Args[len(os.Args)-1] // make more readable - create lastArg variable.
		envelope.Name = lastArg
		envelope.UUID = getUUID()
		envelope.Alias = envelope.Name
		envelope.Label = envelope.Name
		envelope.ContentType = _ENVELOPE_TYPE
		envelope.Checksum, _ = getStringSHA1(envelope.Name)
		envelope.ArtifactList = []ArtifactItem{}
		envelope.URIList = []URIRecord{}
		envelope._onLedger = _FALSE
		envelope._envelopeUUID = _NULL_UUID
		envelope._envelopePath = "/"
		// check if openchain flag is present
		envelope.OpenChain = _FALSE
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == strings.ToLower("--openchain") {
				envelope.OpenChain = _TRUE
			}
		}
		// Add to database
		envelopeName := envelope._envelopePath + envelope.Name
		err = addArtifactToDB(envelope)
		if err != nil {
			fmt.Printf("	error:  %s%s%s\n", _RED_FG, envelopeName, _COLOR_END)
			return // exit
		} else {
			fmt.Printf("	creating: %s%s%s\n", _GREEN_FG, envelopeName, _COLOR_END)
		}
		setLocalConfigValue(_ENVELOPE_KEY, envelope.UUID)
		setAlias(envelope.Name, envelope.UUID)
		return // exit

		/*************
		alias := ""
		notDone := true
		for notDone {
			if getkeyboardYesNoReponse("Would you like to create an alias for the envelope (y/n)?") {
				alias = getkeyboardReponse("Enter the alias you would like to use?")
				if alias == "" {
					continue
				}
				err := setAlias(alias, envelope.UUID)
				if err != nil {
					fmt.Printf("Error: Unable to create alias '%s'.\n", alias)
				} else {
					fmt.Printf("Use expression '%s%s' to reference this envelope\n", _ALIAS_TOKEN, alias)
				}
				notDone = false
			} else {
				notDone = false
			}
		}
		*******************/

		/***********
		r, err := regexp.Compile("^[A-Za-z0-9_][A-Za-z0-9._-]*$")
		name := ""
		for name == "" {
			name = getkeyboardReponse("What is the envelope name (q=quit)?")
			if strings.ToLower(name) == "q" {
				return
			}
			if len(name) > 20 || !r.MatchString(name) {
				fmt.Printf("invalid syntax '%s'. Please try again", name)
				name = ""
			}
		}
		***********************/
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

// handle: 'sparts init' command
func initRequest() {
	var spartsDirectory string
	var localConfigFile string
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
	if _SEED_FUNCTION_ON {
		seedRequest()
	}
}

// handle: 'sparts network' command
func networkRequest() {
	if len(os.Args[1:]) == 1 {
		//// Display help
		////fmt.Println(_NETWORK_HELP_CONTENT)
		networkName := getLocalConfigValue(_LEDGER_NETWORK_KEY)
		if len(networkName) == 0 {
			fmt.Println("The network has not be assigned")
		}
		fmt.Println(networkName)
		return
	}
	switch os.Args[2] {
	case "--get", "-g":
		networkName := getLocalConfigValue(_LEDGER_NETWORK_KEY)
		if len(networkName) == 0 {
			fmt.Println(" The network has not be assigned")
		}
		fmt.Println(networkName)
	case "--help", "-help", "help", "-h":
		// Display help
		fmt.Println(_NETWORK_HELP_CONTENT)
	case "--list", "-l":
		displayNetworkList()
	default:
		fmt.Printf("%s: not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Println(_NETWORK_HELP_CONTENT)
	}
}

// handle: 'sparts org' command
func orgRequest() {
	if len(os.Args[1:]) == 1 {
		// Display help
		fmt.Println(_SUPPLIER_HELP_CONTENT)
		return
	}
	//fmt.Println ("num", len(os.Args[1:]))
	switch os.Args[2] {
	case "--list", "-l":
		displaySupplierList2()
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
		name, alias, url := " ", " ", " "
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
			case "alias":
				alias = arg[1]
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
			uuid, err := createSupplier(name, alias, "", "abc123", url)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), strings.ToLower("A connection attempt failed")) {
					// expecting something like error msg:
					//   Post http://35.166.246.146:818/ledger/api/v1/suppliers: dial tcp 35.166.246.146:818:
					//   connectex: A connection attempt failed because the connected party did not properly respond after
					//   a period of time, or established connection failed because connected host has failed to respond.
					//displayErrorMsg(err.Error())
					displayErrorMsg(fmt.Sprintf("ledger node not responding. Might try '%s synch' to locate active ledger node.", filepath.Base(os.Args[0])))
				}

				////fmt.Printf("  Not able to create new supplier: '%s'\n", name)
				return
			} else {
				fmt.Printf("  new supplier '%s' has uuid = %s\n", name, uuid)
			}
		}
	default:
		fmt.Printf("%s: not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Println(_SUPPLIER_HELP_CONTENT)
	}
}

// handle: 'sparts part' command
func partRequest() {
	var part PartRecord

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
		if !isValidUUID(supplierUUID) {
			fmt.Println("Supplier UUID not properly set in local the config file.")
			return
		}
		// Read from standard input.
		scanner := bufio.NewScanner(os.Stdin)
		//name, version, alias, license, description, url, uuid, checksum := "", "", "", "", "", "", "", ""

		fmt.Println("  -------------------------------------------------------------------------")
		fmt.Println("  required field (*)")
		part.Name = ""
		for part.Name == "" {
			fmt.Print("  name(*): ")
			scanner.Scan()
			part.Name = scanner.Text()
		}

		fmt.Print("  version: ")
		scanner.Scan()
		part.Version = scanner.Text()

		fmt.Print("  alias (nickname): ")
		scanner.Scan()
		part.Alias = scanner.Text()
		if len(part.Alias) > 0 {

		}

		fmt.Print("  licensing: ")
		scanner.Scan()
		part.Licensing = scanner.Text()

		fmt.Print("  description: ")
		scanner.Scan()
		part.Description = scanner.Text()

		h := sha1.New()
		h.Write([]byte("Some String"))
		bs := h.Sum(nil)
		part.Checksum = fmt.Sprintf("%x", bs)

		// generate part uuid
		part.UUID = getUUID()
		// Get part root UUID and store in Label for now.
		part.Label = "root:" + getUUID()

		fmt.Println("  \n  Do you want to proceed to create a part with the following values (y/n)?")
		fmt.Println("  -------------------------------------------------------------------------")
		fmt.Println("\tname      	= " + part.Name)
		fmt.Println("\tversion   	= " + part.Version)
		fmt.Println("\talias     	= " + part.Alias)
		fmt.Println("\tlicensing 	= " + part.Licensing)
		fmt.Println("\tchecksum  	= " + part.Checksum)
		fmt.Println("\tuuid      	= " + part.UUID)
		fmt.Println("\tdescription 	= " + part.Description)

		fmt.Print("(y/n)> ")
		scanner.Scan()
		confirmation := strings.ToLower(scanner.Text())
		if confirmation == "y" || confirmation == "yes" {
			// Add part to db
			err := addPartToDB(part)
			if err != nil {
				displayErrorMsg(err.Error())
				return
			}
			// Set default part in local config file
			setLocalConfigValue(_PART_KEY, part.UUID)

			fmt.Printf("%s%s%s%s\n", _INDENT_STR, _CYAN_FG, "submitting part to ledger ....", _COLOR_END)
			err = pushPartToLedger(part)
			if err != nil {
				displayErrorMsg(err.Error())
				return
			}
			// Part was successfully registered with the legder and added to the db.
			// Now create th ledger relationship between part and supplier.
			ok, err := createPartSupplierRelationship(part.UUID, supplierUUID)
			if ok == true {
				////fmt.Println("Part creation on ledger was SUCCESSFUL.")
				fmt.Printf("	pushed part: %s%s%s\n", _GREEN_FG, part.Name, _COLOR_END)
				if getkeyboardYesNoReponse("Would you like to create an alias for this part (y/n)?") {
					alias := getkeyboardReponse("Enter the alias you would like to use?")
					err := setAlias(alias, part.UUID)
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
			partID = getLocalConfigValue(_PART_KEY)
			if partID == _NULL_UUID {
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
				if strings.ToLower(partID) == strings.ToLower(_NULL_UUID) {
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
		if len(os.Args) == 3 {
			// get supplier part list from db (local cache)
			displayPartList()
			return
		}

		if len(os.Args[2:]) > 1 && (os.Args[3] == "--all" || os.Args[3] == "-a") {
			////if len(os.Args[3:]) == 0 {
			partsList, err := getPartListFromLedger()
			if checkAndReportError(err) {
				return
			}
			partItemList := []PartItemRecord{}
			var partItem PartItemRecord
			for k := range partsList {
				partItem.PartUUID = partsList[k].UUID
				partItemList = append(partItemList, partItem)
			}
			displayPartsFromLedger(partItemList)
		} else {
			// displaySupplierParts()
			fmt.Println("  - Not Implemented Yet - list my (supplier) parts")
			fmt.Printf("  - try: '%s %s --list --all' to list all the parts registered with the ledger\n", filepath.Base(os.Args[0]), os.Args[1])
		}
	case "--set":
		var uuid string
		if len(os.Args[3:]) >= 1 {
			uuid = strings.ToLower(os.Args[3])
			if isValidUUID(uuid) || strings.ToLower(uuid) == strings.ToLower(_NULL_UUID) {
				// we are good.
				setLocalConfigValue(_PART_KEY, uuid)
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
				if isValidUUID(idStr[1]) || strings.ToLower(idStr[1]) == strings.ToLower(_NULL_UUID) {
					// we are good. idStr[1] holds the uuid.
					setLocalConfigValue(_PART_KEY, strings.ToLower(idStr[1]))
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
				filepath.Base(os.Args[0]), strings.ToLower(_NULL_UUID))
		} else {
			fmt.Printf("UUID '%s' is not properly formatted\n", idStr[1])
		}
		****/
	default:
		fmt.Printf("%s: not a valid argument for %s\n", os.Args[2], os.Args[1])
		fmt.Println(_PART_HELP_CONTENT)
	}
}

// handle: 'sparts ping' command
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

func pushRequest() {
	//var partRecord PartRecord
	//var envelopeRecord ArtifactRecord
	var envelopeUUID string
	var partUUID string
	var err error

	if len(os.Args) == 3 && (strings.ToLower(os.Args[2]) == "--help" || strings.ToLower(os.Args[2]) == "-h") {
		fmt.Println(_PUSH_HELP_CONTENT)
		return
	}
	if len(os.Args) != 4 {
		displayErrorMsg(fmt.Sprintf("wrong number of arguments specified. Try '%s %s --help'", filepath.Base(os.Args[0]), os.Args[1]))
		return
	}
	if strings.ToLower(os.Args[2]) == "--help" || strings.ToLower(os.Args[2]) == "-h" {
		fmt.Println(_PUSH_HELP_CONTENT)
		return
	}
	if !(strings.ToLower(os.Args[2]) == "envelope" || strings.ToLower(os.Args[3]) == "ledger") {
		displayErrorMsg(fmt.Sprintf("Argument '%s' is not valid. See '%s %s --help'", os.Args[2], filepath.Base(os.Args[0]), os.Args[1]))
		return
	}

	// Ok, let's push the staging areas to the ledger

	// Check part and envelope are assigned as default.
	partUUID = getLocalConfigValue(_PART_KEY)
	if partUUID == _NULL_UUID || !isValidUUID(partUUID) {
		displayErrorMsg(fmt.Sprintf("Default PART uuid has NOT been assigned in local config file. See '%s %s --help'", filepath.Base(os.Args[0]), os.Args[1]))
		return
	}
	/*****
	// grab part record from local cache (db)
	partList, err := getPartListFromDBWhere("UUID", partUUID)
	// I expect one (and only one record).
	if len(partList) == 0 || err != nil {
		// no part found
		displayErrorMsg("Cannot obtain default part uuid from local config file")
		return // exit
	}
	// again, we expect partList to contain one (and only one record)
	//partRecord = partList[0]
	*****/

	// Get envelope record
	if envelopeUUID = getLocalConfigValue(_ENVELOPE_KEY); envelopeUUID == _NULL_UUID {
		//displayErrorMsg(fmt.Sprintf("Default ENVELOPE uuid has NOT been assigned in local config file. See '%s %s --help'", filepath.Base(os.Args[0]), os.Args[1]))
		//return
		if !getkeyboardYesNoReponse("No envelope specified. Do you want to push to the part root envelope (y/n)?") {
			fmt.Println("push request has been cancelled")
			return
		}
		/*****
		// will push to the part root envelope.
		//
		// grab root envelope id from part record
		partRootEnvelopeUUID = getRootEnvelope(partRecord)
		// make sure part root envelope id exists.
		if !isValidUUID (partRootEnvelopeUUID) {
			// Report error but user does not need to understand about the part root envelope
			displayErrorMsg("Cannot obtain part root envelope record.")
			if _DEBUG_DISPLAY_ON {
				// print for easier debgging
				displayErrorMsg("trouble locating the root envelope")
			}
			return // exit
		}
		envelopeUUID = partRootEnvelopeUUID
		*****/
	}

	/***************
		// grab part record from local cache (db)
		partList, err := getPartListFromDBWhere("UUID", partUUID)
		// I expect one (and only one record).
		if len(partList) = 0 || err != nil {
			// no part found
			displayErrorMsg("Cannot obtain default part uuid from local config file")
			return // exit
		}
		// again, we expect partList to contain one (and only one record)
		partRecord := partList[0]
		// grab root envelope id from part record
		partRootEnvelopeUUID = getRootEnvelope()
		// make sure part root envelope id exists.
		if !isValidUUID (partRootEnvelopeUUID) {
			// Report error but user does not need to understand about the part root envelope
			displayErrorMsg("Cannot obtain part record.")
			if _DEBUG_DISPLAY_ON {
				// print for easier debgging
				displayErrorMsg("trouble locating the root envelope")
			}
			return // exit
		}
	*****************/

	// Get list of artifacts.
	var displayArtifacts []ArtifactRecord

	// getArfitactFromDB where uuid == envelopeID

	// Push envelope if not on ledger already.
	var envelope ArtifactRecord
	list, err := getArtifactListInDBWhere("UUID", envelopeUUID)
	if err != nil {
		displayErrorMsg(err.Error())
		return
	}
	// TODO: We are getting a list but show modify call to return a single record
	if len(list) > 0 {
		envelope = list[0]
	}

	//fmt.Printf("Isss: '%s'  '%d'  '%s'\n", envelope.UUID, len(list), envelopeUUID)
	//return
	envelope._envelopePath = "/"

	if envelope._onLedger == _FALSE {
		ok, err := pushArtifactToLedger(envelope)
		if !ok || err != nil {
			// Error occurred
			fmt.Printf("	error pushing:  %s%s%s\n", _RED_FG, envelope.Name, _COLOR_END)
			fmt.Printf("  Aborting the 'push' request.\n")
			return
		}
		// So far so good. Update the envelope status in db.
		envelope._onLedger = _TRUE
		id := strconv.Itoa(envelope._ID)
		err = updateArtifactInDB("_onLedger", _TRUE, id)
		if err != nil {
			if _DEBUG_DISPLAY_ON {
				fmt.Printf("DB error updating db for: %s - \n", envelope.Name, err)
			}
		}
		/*****
		partUUID := getLocalConfigValue(_PART_KEY)
		if !isValidUUID(partUUID) {
			displayErrorMsg("Default part UUID not properly set in local config file.")
			return
		}
		****/
		partAlias, err := getAliasUsingValue(partUUID)
		if partAlias == "" || err != nil {
			partAlias = partUUID
		}
		envelopeAlias, err := getAliasUsingValue(envelope.UUID)
		if envelopeAlias == "" || err != nil {
			envelopeAlias = envelope.UUID
		}
		err = createArtifactOfPartRelation(envelope.UUID, partUUID)
		/****
		if err != nil {
			fmt.Printf("	error pushing:  %s%s/%s%s\n", _RED_FG, partAlias, envelope.UUID, _COLOR_END)
		} else {
			// Report success.
			fmt.Printf("	pushing relation: %s%s/%s%s\n", _GREEN_FG, partAlias, envelope.UUID, _COLOR_END)
		}
		*****/
	}
	// Made sure the Envelope is on the ledger
	// Now push the artifacts.Start by getting the artifact list
	displayArtifacts, err = getEnvelopeArtifactList(envelopeUUID, true)
	if err != nil {
		displayErrorMsg(err.Error())
		return
	}
	for _, artifact := range displayArtifacts {
		//if artifact.ContentType == _ENVELOPE_TYPE || artifact.UUID == envelopeUUID {
		if artifact.ContentType == _ENVELOPE_TYPE {
			continue // it's the envelope. skip
		}
		if artifact._onLedger == _FALSE {
			fmt.Printf("%s%s%s%s\n", _INDENT_STR, _CYAN_FG, "submitting artifact to ledger ....", _COLOR_END)
			ok, err := pushArtifactToLedger(artifact)
			if !ok || err != nil {
				// Error occurred
				fmt.Printf("	error pushing:  %s%s%s\n", _RED_FG, artifact.Name, _COLOR_END)
				////fmt.Sprintf("error pushing artifact '%s. See '%s %s --help'", filepath.Base(os.Args[0]), os.Args[1]))
				continue // go get next artifact
			} else {
				// So far so good. Update the artifact status in db.
				artifact._onLedger = _TRUE
				id := strconv.Itoa(artifact._ID)
				err = updateArtifactInDB("_onLedger", _TRUE, id)
				if err != nil && _DEBUG_DISPLAY_ON {
					fmt.Printf("DB error updating db onLedger status for: %s - \n", artifact.Name, err)
				}
				////artifact._envelopeUUID = envelope.UUID
				// Assign artifact the envelope id
				err = updateArtifactInDB("_envelopeUUID", envelope.UUID, id)
				if err != nil && _DEBUG_DISPLAY_ON {
					fmt.Printf("DB error updating db envelope UUID for: %s - \n", artifact.Name, err)
				}

				// Create relationship between envelope and artifact on ledger
				////fmt.Println(" '%s' '%s' '%s'", artifact.UUID, envelopeUUID)
				err = createArtifactOfEnvelopeRelation(artifact.UUID, envelopeUUID, artifact._envelopePath)
				if err != nil {
					fmt.Printf("	error pushing:  %s%s%s\n", _RED_FG, artifact.Name, _COLOR_END)
				} else {
					// Report success.
					fmt.Printf("	pushed artifact: %s%s%s\n", _GREEN_FG, artifact.Name, _COLOR_END)
				}
			}
		}
	}
}

// handle: 'sparts remove' command
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
	case "--all":
		artifacts, err := getArtifactListDB()
		if err != nil {
			fmt.Println("  fatal: sparts working database not accessible.")
			os.Exit(_DB_ACCESSS_ERROR) // exit program
		}
		for _, artifact := range artifacts {
			// Delete aritfact record
			deleteArtifactFromDB(artifact)
		}

	default:
		////r, err := regexp.Compile("^[1-9]-[A-Za-z0-9._-]*$")
		//firsttime := true
		for i := 2; i <= len(os.Args[1:]); i++ {
			////if r.MatchString(os.Args[i])
			id := os.Args[i]
			if _, err := strconv.Atoi(id); err != nil {
				fmt.Printf("	error:  %s%s%s - id has invalid syntax. Expecting an integer.\n", _RED_FG, id, _COLOR_END)
				continue
			}
			/**************
			if _, err := strconv.Atoi(os.Args[i]); err != nil {
				fmt.Printf("  %s is not a valid integer (id)\n", os.Args[i])
				if firsttime {
					firsttime = false
					fmt.Println("  Usage: sparts remove -all|<id>+. Obtain id using: 'sparts status'")
					fmt.Println("   e.g.,  sparts remove 4")
					fmt.Println("  Obtain id from: 'sparts status'")
					fmt.Println()
				}
				continue
			}
			************************/
			// get the artifact record for specified id
			artifact, err := getArtifactFromDB("id", id)
			if err != nil {
				fmt.Printf("	error:   %sid=%s%s - is not valid. Try 'sparts status'.\n", _RED_FG, id, _COLOR_END)
				continue
			}
			// Delete aritfact record for id. Give error if id not found.
			if !deleteArtifactFromDB(artifact) {
				/////fmt.Println("  artifact not found for id =", os.Args[i])
				fmt.Printf("	error:   %s%s%s - artifact record not found\n", _RED_FG, id, _COLOR_END)
			} else {
				fmt.Printf("	removed: %s%s: %s %s\n", _GREEN_FG, id, artifact.Name, _COLOR_END)
			}
		}
	}
}

// handle: 'sparts status' command
func statusRequest() {

	// At least one argument.

	if len(os.Args) == 2 {
		//  'sparts status'
		// ---------------------------
		// Display Staging Area Table
		//----------------------------
		displayStagingTable2()
		return
	}

	// one or more additional arguments.
	switch strings.ToLower(os.Args[2]) {
	case "-h", "--help":
		// Display help
		fmt.Println(_STATUS_HELP_CONTENT)
	case "-v", "--view":
		numArguments := len(os.Args)
		if numArguments == 3 {
			displayErrorMsg("The '--view' option is expecting another argument. Try '--help' for more details")
			return
		}
		// sparts status --view id ...
		//   0		1	  2     3
		//                      ^
		for i := 3; i < numArguments; i++ {
			id := os.Args[i] // argument is database id
			artifact, err := getArtifactFromDB("id", id)
			if err != nil {
				fmt.Printf("Artifact with id=%s does not exist\n", id)
				continue
			}
			var path string
			if isPathURL(artifact._contentPath) {
				// path is a url
				path = artifact._contentPath
			} else {
				// We have a file path (and not a url). Obtain abridged version
				path = getAbridgedFilePath(artifact._contentPath)
			}
			path = artifact._contentPath

			fmt.Println()
			fmt.Println("|--------------------------------------------------------")
			fmt.Printf("| %sArtifact%s: %s%s%s\n", _WHITE_FG, _COLOR_END, _CYAN_FG, artifact.Name, _COLOR_END)
			const padding = 0
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '.', tabwriter.Debug)
			//w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '.', 0)
			fmt.Fprintf(w, "\t%s\t%s\n", " -----------", "-----------------------------------------")
			//fmt.Fprintf(w, "\t %s\t %s\n", "Name", artifact.Name)
			fmt.Fprintf(w, "\t %s\t %s\n", "UUID", artifact.UUID)
			fmt.Fprintf(w, "\t %s\t %s\n", "Alias", artifact.Alias)
			fmt.Fprintf(w, "\t %s\t %s\n", "Label", artifact.Label)
			fmt.Fprintf(w, "\t %s\t %s\n", "Type", artifact.ContentType)
			fmt.Fprintf(w, "\t %s\t %s\n", "Checksum", artifact.Checksum)
			fmt.Fprintf(w, "\t %s\t %s\n", "OpenChain", artifact.OpenChain)
			fmt.Fprintf(w, "\t %s\t %s\n", "Content Path", path)

			fmt.Fprintf(w, "\t %s\t %s\n", "On Ledger", artifact._onLedger)

			if len(artifact._envelopePath) > 0 {
				fmt.Fprintf(w, "\t %s\t %s\n", "Envelope Path", artifact._envelopePath)
			}
			//fmt.Fprintf(w, "\n")
			w.Flush()
			fmt.Println("|--------------------------------------------------------")
			///////fmt.Println("is ...", artifact._onLedger)
		}
	default:
		fmt.Printf("  error - '%s'is not a valid argument for %s\n", os.Args[2], filepath.Base(os.Args[0]))
		return
	} // end of switch
}

// handle: 'sparts supplier' command - which has been discontinued.
func supplierRequest() {
	fmt.Print("  'supplier' command is no longer supported. Use 'org' in its place.\n")
}

// handle: 'sparts synch' command
func synchRequest() {

	fmt.Println("  Network: ", getLocalConfigValue(_LEDGER_NETWORK_KEY))
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
	fmt.Println("  Searching for a new primary ledger node .....")

	// Obtain current list of available ledger nodes from look up directory (atlas)
	nodeList, err := getLedgerNodeList()
	if err != nil {
		fmt.Println(" ", err)
		// Suggest a fix for certain circumstances
		if strings.Contains(err.Error(), "does not exist") {
			fmt.Printf("  You may need to set or update local config variable: '%s'\n", _LEDGER_NETWORK_KEY)
			fmt.Printf("  To view local and global variables try: %s config --list\n", filepath.Base(os.Args[0]))
		}
		return
	}
	// Check if list is empty
	if len(nodeList) == 0 {
		fmt.Printf("  The network '%s' has no ledger nodes registered\n", getLocalConfigValue(_LEDGER_NETWORK_KEY))
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
		fmt.Printf("  Not able to locate an active ledger node referring the %s directory\n", _ATLAS)
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

// handle: 'sparts version' command
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

// // handle: 'sparts seed' command
func seedRequest() {
	var err error
	if !_SEED_FUNCTION_ON {
		fmt.Println(" Seed function NOT implemented")
		return
	}

	//os.Args[0] = "sparts"

	// seed
	//os.Args = append(os.Args, "test2") // 2
	//os.Args = append(os.Args, "test3") // 3
	//os.Args = append(os.Args, "test4") // 4

	//setLocalConfigValue(_LEDGER_NETWORK_KEY, "sparts-test-network")
	setLocalConfigValue(_LEDGER_NETWORK_KEY, "zephyr-parts-network")
	setLocalConfigValue(_LEDGER_ADDRESS_KEY, "35.197.7.42:818")
	//setLocalConfigValue(_LEDGER_ADDRESS_KEY, "147.11.176.111:818")

	//setLocalConfigValue(_PRIVATE_KEY, "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq")
	//setLocalConfigValue(_PUBLIC_KEY, "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515")

	setLocalConfigValue(_PRIVATE_KEY, "147b72b747a643136d313962eb3c774b972eebb8f47e33a494ffcd542f8f22b8")
	setLocalConfigValue(_PUBLIC_KEY, "03241be9afb64bc15844c2e0f319ee75c41509b927230e04c02f55fc07a78bc014")

	synchRequest()
	/////setLocalConfigValue(_PART_KEY, "zephyr-parts-network")
	if err = setLocalConfigValue(_SUPPLIER_KEY, "3568f20a-8faa-430e-7c65-e9fce9aa155d"); err != nil {
		fmt.Println("Seeding", err)
	}
	//if err = setLocalConfigValue(_PART_KEY, "fd6462e4-9560-4c7f-614c-a87f8ff792b8"); err != nil {

	/*************
	if err = setLocalConfigValue(_PART_KEY, "tbd"); err != nil {
		fmt.Println("Seeding", err)
	}
	**************/

	//setAlias("p1", "fd6462e4-9560-4c7f-614c-a87f8ff792b8")

	ok, err := pingServer(_LEDGER)
	if !ok {
		fmt.Println("Can't access ledger. Can't complete the 'seed' request")
		return
	}

	// list, err := getSupplierList()
	_, err = getSupplierList()
	if err == nil /* && len(list) == 0 */ {
		supplier := OrganizationRecord{}
		supplier.Name = "Wind River"
		supplier.Alias = "WR"
		supplier.Type = "supplier"
		//supplier.Description = "Wind River, leading RTOS supplier"
		supplier.Description = "Wind"

		supplier.UUID = "2567f20a-8faa-430e-7c65-e9fce9aa155d"
		supplier.Url = "http://www.windriver.com"
		supplier.Parts = []PartItemRecord{}
		err = pushSupplierToLedger(supplier)
		if err != nil {
			displayErrorMsg("encountered problem adding Wind River supplier to ledger")
		} else {
			setAlias("wr", supplier.UUID)
		}

		supplier = OrganizationRecord{}
		supplier.Name = "Zephyr Project"
		supplier.Alias = "Zephyr"
		supplier.Type = "supplier"
		supplier.Description = "Zephyr-project-part-network"
		supplier.UUID = "7234f20a-85bc-121a-39ac-2c5ce9dc167a"
		supplier.Url = "http://www.zephyrproject.org"
		supplier.Parts = []PartItemRecord{}
		err = pushSupplierToLedger(supplier)
		if err != nil {
			displayErrorMsg("encountered problem adding Zephyr Project supplier to ledger")
		} else {
			setAlias("zephyr", supplier.UUID)
		}

		supplier = OrganizationRecord{}
		supplier.Name = "Intel Corp"
		supplier.Alias = "Intel"
		supplier.Type = "supplier"
		supplier.Description = "Intel, leading hardware supplier"
		supplier.UUID = "1f54f20a-85bc-9e1a-81d1-611ce9d2b122"
		supplier.Url = "http://www.intel.com"
		supplier.Parts = []PartItemRecord{}
		err = pushSupplierToLedger(supplier)
		if err != nil {
			displayErrorMsg("encountered problem adding 'Intel' supplier to ledger")
		} else {
			setAlias("intel", supplier.UUID)
		}

	}

}

// handle: 'sparts tips' command
func tipsRequest() {
	fmt.Println(_TIPS_CONTENT)
}

func quickRequest() {
	fmt.Println("	1)	Add URI to Artifact")

	choice := getkeyboardReponse("Select: ?")

	switch choice {
	case "1":
		var uri URIRecord

		uuid := getkeyboardReponse("Artifact uuid?")
		uri.Version = getkeyboardReponse("URI Version?")
		uri.Checksum = getkeyboardReponse("URI Checksum?")
		uri.Size = getkeyboardReponse("URI Size (bytes)?")
		uri.ContentType = getkeyboardReponse("Content Type (.pdf)?")
		uri.URIType = getkeyboardReponse("URI Type (http)?")
		uri.Location = getkeyboardReponse("Location?")

		addURIToArtifact(uuid, uri)
	}
}

// Used for special testing
func testRequest() {

	/****
	keys, err := getPrivatePublicKeys()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("success")
		fmt.Println("  Public Key:", keys.PublicKey)
		fmt.Println("  Private Key:", keys.PrivateKey)
	}
	*****/

	// keys:
	// Private: bf6bb6df3afdbe2cdda1ce4e92d4fbda46a49586832c2dd09900981bfdd37f2b
	// Public: 030c9148861b4ee085118bc44a235d961b56dbfd4c01f0d8d6391e923fe04889e9

	err := registerUser("user007", "mark@windriver.com", "member", "allow", "030c9148861b4ee085118bc44a235d961b56dbfd4c01f0d8d6391e923fe04889e9")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("success")
	}

	/*********
		artifact, err := getArtifactFromLedger("d2538468-9245-446c-4b6b-90068f2d8713")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(len(artifact.ArtifactList))
	***********/

	/***********
		envelopeUUID := os.Args[2]

		partUUID := getLocalConfigValue(_PART_KEY)
		if !isValidUUID(partUUID) {
			displayErrorMsg("Default part UUID not properly set in local config file.")
			return
		}
		alias, err := getAlias(partUUID)
		if alias == "" || err != nil {
			alias = partUUID
		}
		err = createArtifactOfPartRelation(envelopeUUID, partUUID)
		if err != nil {
			fmt.Printf("	error pushing:  %s%s%s\n", _RED_FG, alias, _COLOR_END)
		} else {
			// Report success.
			fmt.Printf("	pushing: %s%s%s\n", _GREEN_FG, alias, _COLOR_END)
		}
	*******************/

	/****
	var part PartRecord
	part.Label = os.Args[2]
	fmt.Println(getRootEnvelope(part))
	****/

	/******

	r, _ := regexp.Compile(`^[1-9]+\-[0-9]+$`)
	s, _ := regexp.Compile(`^[1-9][0-9]*$`)
	var idList []int
	for i := 2; i < len(os.Args); i++ {
		if r.MatchString(os.Args[i]) || s.MatchString(os.Args[i]) {
			if r.MatchString(os.Args[i]) {
				list := strings.Split(os.Args[i], "-")
				first, _ := strconv.Atoi(list[0])
				last, _ := strconv.Atoi(list[1])
				for i := first; i <= last; i++ {
					////fmt.Print(i, " ")
					idList = append(idList, i)
				}

			} else {
				id, err := strconv.Atoi(os.Args[i])
				if err == nil {
					idList = append(idList, id)
				}
			}

		}
	}
	fmt.Println("len=", len(idList))
	for i := 0; i < len(idList); i++ {
		fmt.Print(idList[i], " ")
	}
	fmt.Println()
	************/
}

package main

// This module is where the sparts config  routines are found.
// There is support for the "sparts config" and "sparts alias" command

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
	//"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
)

// Config file format
// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type localConfig struct {
	LookUp string `json:"lookup"`
	Node   struct {
		LedgerAddress    string `json:"ledger_address"`
		ConductorAddress string `json:"conductor_address"`
	} `json:"node"`
	PartUUID     string `json:"part_uuid"`
	PrivateKey   string `json:"private_key"`
	PublicKey    string `json:"public_key"`
	SupplyChain  string `json:"supply_chain"`
	SupplierUUID string `json:"supplier_uuid"`
}

// The default local config file contents can be found in build_config.go
//_LOCAL_CONFIG_FILE_CONTENT = ...

// local config file field names
const (
	_CONDUCTOR_ADDRESS_KEY    = "node.conductor_address"
	_LEDGER_ADDRESS_KEY       = "node.ledger_address"
	_PART_UUID_KEY            = "part_uuid"
	_SUPPLY_CHAIN_NETWORK_KEY = "supply_chain"
	_PRIVATE_KEY              = "private_key"
	_PUBLIC_KEY               = "public_key"
)

// global Config file format
type globalConfig struct {
	Atlas string `json:"atlas"`
	User  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

// The default gola config file contents can be found in build_config.go
//_GLOBAL_CONFIG_FILE_CONTENT = ...

// global config file field names
const (
	_ATLAS_ADDRESS_KEY = "atlas"
	_USER_NAME_KEY     = "user.name"
	_USER_EMAIL_KEY    = "user.email"
)

type AlisaRecord struct {
	Alias string `json:"name"`
	Value string `json:"value"`
}

//----------------------------------------------
// 			 Local Config routines
//----------------------------------------------

// getLocalConfigValue returns the value of a local config value.
// for key
func getLocalConfigValue(key string) string {
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	configFile := spartsDir + "/" + _LOCAL_CONFIG_FILE
	configData, err := readLocalConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	switch strings.ToLower(key) {
	case "node.ledger_address":
		return configData.Node.LedgerAddress

	case "node.conductor_address":
		return configData.Node.ConductorAddress
	case "part_uuid":
		return configData.PartUUID
	case _PRIVATE_KEY:
		return configData.PrivateKey
	case _PUBLIC_KEY:
		return configData.PublicKey
	case "supply_chain":
		return configData.SupplyChain
	case "supplier_uuid":
		return configData.SupplierUUID
	case "lookup":
		return configData.LookUp
	}
	return ""
}

// setLocalConfigValue assigns the value of a local config variable.
// The variable identified by key is set to newValue
func setLocalConfigValue(key string, newValue string) {
	// Read config file
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}
	configFile := spartsDir + "/" + _LOCAL_CONFIG_FILE
	configData, err := readLocalConfig(configFile)

	if err != nil {
		// Error reading config file
		fmt.Println(err)
		return
	}
	switch strings.ToLower(key) {
	case "node.ledger_address":
		configData.Node.LedgerAddress = newValue
	case "node.conductor_address":
		configData.Node.ConductorAddress = newValue
	case "part_uuid":
		if isValidUUID(newValue) || strings.ToLower(newValue) == strings.ToLower(_NULL_PART) {
			configData.PartUUID = newValue
		} else {
			fmt.Println("  UUID syntax is not valid.")
			fmt.Println("  The updated value was not saved.")
			return // done due to error.
		}
	case _PRIVATE_KEY:
		configData.PrivateKey = newValue
	case _PUBLIC_KEY:
		configData.PublicKey = newValue
	case "supplier_uuid":
		if isValidUUID(newValue) {
			configData.SupplierUUID = newValue
		} else {
			fmt.Println("  UUID syntax is not valid.")
			fmt.Println("  The updated value was not saved.")
			return // done due to error.
		}
	case "lookup":
		configData.LookUp = newValue
	default:
		fmt.Printf("  '%s' is not a validate local configuration value\n", key)
	}
	// Save updated config values
	d, err := yaml.Marshal(&configData)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configFile, d, 0644)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))
}

// readLocalConfig reads in the contents of the local config file
func readLocalConfig(configPath string) (configData localConfig, err error) {
	y, err := ioutil.ReadFile(configPath)
	if err != nil {
		return configData, err
	}

	err = yaml.Unmarshal(y, &configData)
	if err != nil {
		////log.Printf("yaml.Unmarshal err   #%v ", err)
		return configData, fmt.Errorf("yaml.Unmarshal err   #%v ", err)
	}
	return configData, nil
}

// displayLocalConfigData displays the local config
func displayLocalConfigData() {
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := readLocalConfig(spartsDir + "/" + _LOCAL_CONFIG_FILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	d, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	//fmt.Println("node.supply_chain=", data.SupplyChain)
	fmt.Printf("%s\n", string(d))
}

//----------------------------------------------
// 			Global Config File routines
//----------------------------------------------

func getGlobalConfigValue(value string) string {
	////spartsDir, err := getSpartsDirectory()

	configFile, err := getGlobalConfigFile()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	configData, err := readGlobalConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	switch strings.ToLower(value) {
	case _ATLAS_ADDRESS_KEY:
		return configData.Atlas

	case _USER_NAME_KEY:
		return configData.User.Name

	case _USER_EMAIL_KEY:
		return configData.User.Email
	}
	return ""
}

func getGlobalConfigFile() (string, error) {
	// Get home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error locating home directory")
		fmt.Println(err)
		return "", err
	}

	globalConfigFile := usr.HomeDir + "/" + _GLOBAL_CONFIG_FILE
	// Check if global config file exists.
	if _, err := os.Stat(globalConfigFile); os.IsNotExist(err) {
		// global config file does not exist. Create one.
		fmt.Println("File does not exist:", globalConfigFile)
		fmt.Println("Creating config file")
		globalFile, err := os.Create(globalConfigFile)
		if err != nil {
			fmt.Println("Error Creating global config file: ", err)
			return "", err
		}
		defer globalFile.Close()
		_, err = globalFile.WriteString(_GLOBAL_CONFIG_FILE_CONTENT)
		if err != nil {
			fmt.Println("Error - writing to global config file: ", err)
			return "", err
		}
	}
	// Clean up string on Windows replace '\' with '/'. Does nothing on Linux.
	globalConfigFile = strings.Replace(globalConfigFile, `\`, `/`, -1)
	return globalConfigFile, nil
}

func setGlobalConfigValue(field string, newValue string) {
	// Read global configation file
	globalConfigFile, _ := getGlobalConfigFile()
	configData, err := readGlobalConfig(globalConfigFile)
	if err != nil {
		// Error reading config file
		fmt.Println("Error reading global configuation file")
		fmt.Println(err)
		return
	}
	switch strings.ToLower(field) {
	case _ATLAS_ADDRESS_KEY:
		configData.Atlas = newValue
	case "user.name":
		configData.User.Name = newValue
	case "user.email":
		configData.User.Email = newValue
	}
	// Save updated config file
	d, err := yaml.Marshal(&configData)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(globalConfigFile, d, 0644)
	if err != nil {
		panic(err)
	}
}

func readGlobalConfig(configPath string) (configData globalConfig, err error) {
	y, err := ioutil.ReadFile(configPath)
	if err != nil {
		//log.Printf ("config file ReadFile err   #%v ", err)
		return configData, err
	}

	err = yaml.Unmarshal(y, &configData)
	if err != nil {
		log.Printf("yaml.Unmarshal err   #%v ", err)
		return configData, err
	}
	return configData, nil
}

func displayGlobalConfigData() {
	// Get home directory and read it in
	globalCoinfigFile, err := getGlobalConfigFile()
	data, err := readGlobalConfig(globalCoinfigFile)
	if err != nil {
		fmt.Println("Error reading the global configuration file")
		fmt.Println(err)
		return
	}

	d, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Printf("%s\n", string(d))
}

//----------------------------------------------
// 			 Alias routines
//----------------------------------------------

func setAlias(alias string, value string) error {

	addAliasValueToDB(alias, value)

	return nil
}

func getAlias(alias string) (string, error) {
	aliasValue, err := getAliasValueFromDB(alias)
	if err != nil {
		return "", err
	}
	return aliasValue, nil
}

func getAliasUsingValue(value string) (string, error) {
	alias, err := getAliasUsingValueFromDB(value)
	if err != nil {
		return "", err
	} else {
		return alias, nil
	}
}

func displayAliases() {

	aliasList, err := getAlisaListFromDB()
	if checkAndReportError(err) {
		return // error, we are done
	}

	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '.',
		tabwriter.Debug)
	fmt.Fprintf(w, " \t %s \t %s\n", "--------------", "  --------------------------")
	fmt.Fprintf(w, " \t %s \t %s\n", "     Alias*    	", "        Value")
	fmt.Fprintf(w, " \t %s \t %s\n", "--------------", "  --------------------------")
	for i := 0; i < len(aliasList); i++ {
		fmt.Fprintf(w, " \t %s \t %s\n", aliasList[i].Alias, aliasList[i].Value)
	}
	// print to screen
	w.Flush()

	fmt.Printf("   ------\n")
	fmt.Printf("   *Aliases are case sensitive\n")
	fmt.Printf("   *Use 'id=' prefix to apply (e.g., sparts supplier --get id=mycompany)\n")
	fmt.Println()

}

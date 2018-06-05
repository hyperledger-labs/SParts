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
	SupplyChain  string `json:"supply_chain"`
	SupplierUUID string `json:"supplier_uuid"`
	PartUUID     string `json:"part_uuid"`
	LookUp       string `json:"lookup"`
	Node         struct {
		LedgerAddress    string `json:"ledger_address"`
		ConductorAddress string `json:"conductor_address"`
	} `json:"node"`
}

const _GLOBAL_CONFIG_FILE_CONTENT = `user:
    name:
    email:
`

// global Config file format
type globalConfig struct {
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

type AlisaRecord struct {
	Alias string `json:"name"`
	Value string `json:"value"`
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
	//fmt.Println( usr.HomeDir )
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

	return globalConfigFile, nil
}

func readLocalConfig(configPath string) (configData localConfig, err error) {
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

func getLocalConfigValue(value string) string {
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
	switch strings.ToLower(value) {
	case "node.ledger_address":
		return configData.Node.LedgerAddress

	case "node.conductor_address":
		return configData.Node.ConductorAddress
	case "part_uuid":
		return configData.PartUUID

	case "supply_chain":
		return configData.SupplyChain

	case "supplier_uuid":
		return configData.SupplierUUID

	case "lookup":
		return configData.LookUp
	}
	return ""
}

func setLocalConfigValue(field string, newValue string) {

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
	switch strings.ToLower(field) {
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
	case "supply_chain":
		configData.SupplyChain = newValue
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
		fmt.Printf("  '%s' is not a validate local configuration value\n", field)
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

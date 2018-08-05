package main

/*
 * This module is where the sparts config  routines are found.
 * It includes support for the "sparts config" and "sparts alias" commands
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
	//"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/ghodss/yaml"
)

// Config file format
// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type localConfig struct {
	AutoSynch     string `json:"auto_synch"`
	EnvelopeUUID  string `json:"envelope_uuid"`
	Focus         string `json:"focus"`
	LedgerAddress string `json:"ledger_address"`
	LedgerNetwork string `json:"ledger_network"`
	PartUUID      string `json:"part_uuid"`
	PrivateKey    string `json:"private_key"`
	PublicKey     string `json:"public_key"`
	SupplierUUID  string `json:"supplier_uuid"`
}

// The default local config file contents can be found in build_config.go
//_LOCAL_CONFIG_FILE_CONTENT = ...

// local config file field names
const (
	//_CONDUCTOR_ADDRESS_KEY = "node.conductor_address"
	_AUTO_SYNCH_KEY     = "auto_synch"
	_ENVELOPE_KEY       = "envelope_uuid"
	_FOCUS_KEY          = "focus"
	_LEDGER_ADDRESS_KEY = "node.ledger_address"
	_LEDGER_NETWORK_KEY = "ledger_network"
	_SUPPLIER_KEY       = "supplier_uuid"
	_PART_KEY           = "part_uuid"
	_PRIVATE_KEY        = "private_key"
	_PUBLIC_KEY         = "public_key"
)

// global Config file format

/***********
type globalConfig struct {
	Atlas string `json:"atlas_address"`
	User  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}
************/
type globalConfig struct {
	Atlas     string `json:"atlas_address"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

// The default gola config file contents can be found in build_config.go
//_GLOBAL_CONFIG_FILE_CONTENT = ...

// global config file field names
const (
	_ATLAS_ADDRESS_KEY = "atlas_address"
	_USER_NAME_KEY     = "user_name"
	_USER_EMAIL_KEY    = "user_email"
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
	case _LEDGER_ADDRESS_KEY:
		return configData.LedgerAddress
	case _ENVELOPE_KEY:
		return configData.EnvelopeUUID
	case _FOCUS_KEY:
		return configData.Focus
	case _PART_KEY:
		return configData.PartUUID
	case _PRIVATE_KEY:
		return configData.PrivateKey
	case _PUBLIC_KEY:
		return configData.PublicKey
	case _LEDGER_NETWORK_KEY:
		return configData.LedgerNetwork
	case _SUPPLIER_KEY:
		return configData.SupplierUUID
	case _AUTO_SYNCH_KEY:
		return configData.AutoSynch
	}
	return ""
}

// setLocalConfigValue assigns the value of a local config variable.
// The variable identified by key is set to newValue
func setLocalConfigValue(key string, newValue string) error {
	// Read config file
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		//fmt.Println(err)
		return err
	}
	configFile := spartsDir + "/" + _LOCAL_CONFIG_FILE
	configData, err := readLocalConfig(configFile)

	if err != nil {
		// Error reading config file
		fmt.Println(err)
		return err
	}
	switch strings.ToLower(key) {
	case _LEDGER_ADDRESS_KEY:
		configData.LedgerAddress = newValue
	//case "node.conductor_address":
	//	configData.Node.ConductorAddress = newValue
	case _LEDGER_NETWORK_KEY:
		configData.LedgerNetwork = newValue
	case _ENVELOPE_KEY:
		configData.EnvelopeUUID = newValue
	case _FOCUS_KEY:
		configData.Focus = newValue
	case _PART_KEY:
		if isValidUUID(newValue) || strings.ToLower(newValue) == strings.ToLower(_NULL_UUID) {
			configData.PartUUID = newValue
		} else {
			return fmt.Errorf("UUID syntax is not valid.")
		}
	case _PRIVATE_KEY:
		configData.PrivateKey = newValue
	case _PUBLIC_KEY:
		configData.PublicKey = newValue
	case _SUPPLIER_KEY:
		if isValidUUID(newValue) {
			configData.SupplierUUID = newValue
		} else {
			return fmt.Errorf("UUID syntax is not valid.")
		}
	case _AUTO_SYNCH_KEY:
		configData.AutoSynch = newValue
	default:
		return fmt.Errorf("  '%s' is not a validate local configuration value\n", key)
	}

	// Save updated config values
	d, err := yaml.Marshal(&configData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, d, 0644)
	if err != nil {
		return err
	}
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))

	return nil
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

	s := reflect.ValueOf(&data).Elem()
	typeOfT := s.Type()
	fmt.Println("|-----------------------------------------------------")
	fmt.Println("|	Local Configuration")
	fmt.Println("|-----------------------------------------------------")
	////fmt.Printf("| %sArtifact%s: %s%s%s\n", _WHITE_FG, _COLOR_END, _CYAN_FG, artifact.Name, _COLOR_END)
	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	////fmt.Fprintf(w, "\t%s\t%s\n", " --------------", "-----------------------------------------")
	//fmt.Fprintf(w, "\t %s\t %s\n", "Name", artifact.Name)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		//fmt.Fprintf(w, "\t %s\t %v\n", typeOfT.Field(i).Name, f.Interface())
		fmt.Fprintf(w, "\t %s\t %v\n", typeOfT.Field(i).Tag.Get("json"), f.Interface())
	}
	w.Flush()
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
		return configData.UserName

	case _USER_EMAIL_KEY:
		return configData.UserEmail
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
		configData.UserName = newValue
	case "user.email":
		configData.UserEmail = newValue
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

// readGlobalConfig reads the global config file and returns a globalConfig record
func readGlobalConfig(configPath string) (configData globalConfig, err error) {
	y, err := ioutil.ReadFile(configPath)
	if err != nil {
		//log.Printf ("config file ReadFile err   #%v ", err)
		return configData, err
	}

	err = yaml.Unmarshal(y, &configData)
	if err != nil {
		return configData, fmt.Errorf("yaml.Unmarshal err   #%v ", err)
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

	////d, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	//fmt.Printf("%s\n", string(d))

	s := reflect.ValueOf(&data).Elem()
	typeOfT := s.Type()

	fmt.Println("|-----------------------------------------------------")
	fmt.Println("|	Global Configuration")
	fmt.Println("|-----------------------------------------------------")
	const padding = 0
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	//fmt.Fprintf(w, "\t%s\t%s\n", " --------------", "-----------------------------------------")
	//fmt.Fprintf(w, "\t %s\t %s\n", "Name", artifact.Name)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		////value := f.Interface()
		////fmt.Println(f.Type())
		////fmt.Println(len(f.Interface()))
		////fmt.Println(getType(f.Interface()))
		////fmt.Println(val.Type().Field(i).Tag.Get("json")
		//fmt.Fprintf(w, "\t %s\t %v\n", typeOfT.Field(i).Name, f.Interface())
		fmt.Fprintf(w, "\t %s\t %v\n", typeOfT.Field(i).Tag.Get("json"), f.Interface())
	}
	w.Flush()
}

//----------------------------------------------
// 			 Alias routines
//----------------------------------------------

// setAlias records the value for 'alias' in the database
func setAlias(alias string, value string) error {

	addAliasValueToDB(alias, value)
	return nil
}

// getAlias retrieves the value for 'alias'
func getAlias(alias string) (string, error) {
	aliasValue, err := getAliasValueFromDB(alias)
	if err != nil {
		return "", err
	}
	return aliasValue, nil
}

// getAliasUsingValue returns the alias that has 'value'.
// If more than one alias exist, it will return one selected at random.
func getAliasUsingValue(value string) (string, error) {
	alias, err := getAliasUsingValueFromDB(value)
	if err != nil {
		return "", err
	} else {
		return alias, nil
	}
}

// displayAliases displays a list of all the alias on the terminal
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

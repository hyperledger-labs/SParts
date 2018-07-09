package main

/*
	Purpose: All the database access routines are found in this file.
*/

// Licensing: Apache-2.0
/*
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

/**** https://github.com/mattn/go-sqlite3
 Copyright (c) 2014 Yasuhiro Matsumoto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
****/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// File state globals
var theDB *sql.DB

const CONFIG_RECORD = "Config_Info"

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Open the database.
func openDB() {

	var err error
	// Using SQLite
	spartsDir, err := getSpartsDirectory()
	if err != nil {
		fmt.Println(err)
		fmt.Println("  fatal: sparts working directory not accessible.")
		os.Exit(_DIR_ACCESS_ERROR)
	}
	theDB, err = sql.Open("sqlite3", spartsDir+"/"+_LOCAL_DB_FILE)
	if err != nil {
		fmt.Println("  fatal: sparts working directory database is not accessible.")
		os.Exit(_DB_ACCESSS_ERROR)
	}
}

// Initialize the database
func initializeDB() {
	createDBTables()
}

// Create database tables
func createDBTables() {

	// TODO: create and store database model version

	openDB()
	defer theDB.Close()

	// Create Aritfacts  Table - a list of Artifacts
	sql_cmd := `
	CREATE TABLE IF NOT EXISTS Artifacts (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		UUID TEXT,
		Name TEXT,
		Alias TEXT,	
		Label TEXT,
		Checksum TEXT,
		Path TEXT,
		OpenChain TEXT,
		ContentType TEXT,
		InsertedDatetime DATETIME
	);
	`
	_, err := theDB.Exec(sql_cmd)
	if err != nil {
		fmt.Println(err)
		fmt.Println("  fatal: sparts database table creation not feasible. ")
		os.Exit(3)
	}

	// Create the Aliases Table - a list of Aliases
	sql_cmd = `
	CREATE TABLE IF NOT EXISTS  Aliases (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Alias TEXT,
		Value TEXT,
		InsertedDatetime DATETIME
	);
	`
	_, err = theDB.Exec(sql_cmd)
	if err != nil {
		fmt.Println(err)
		fmt.Println("  fatal: Could not create sparts database table: Aliases")
		os.Exit(3)
	}

	// Set Alias field to be unique in the Aliases table. If Insert detects existing Alias
	// it will replace existing with new record.
	sql_cmd = `CREATE UNIQUE INDEX idx_Apps_Alias 
				ON Aliases (Alias);`

	_, err = theDB.Exec(sql_cmd)
	if err != nil {
		fmt.Println(err)
	}
}

// Insert Aftifact record into the DB
func AddArtifactToDB(record ArtifactRecord) {

	// TODO: return succes/failure status
	openDB()
	defer theDB.Close()

	sql_additem := `
	INSERT OR REPLACE INTO Artifacts (
		UUID,
		Name,
		Alias,
		Label,
		Checksum,
		Path,
		ContentType,
		OpenChain,
		InsertedDatetime
		) values(?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	_, err2 := stmt.Exec(record.UUID, record.Name, record.Alias, record.Label, record.Checksum,
		record._path, record.ContentType, record.OpenChain)
	if err2 != nil {
		fmt.Println(err2)
	}
	//_, err = res.LastInsertId()
}

// Insert Supply Chain network Application record into the DB
func deleteArtifactFromDB(record ArtifactRecord) bool {

	// TODO: return succes/failure status
	openDB()
	defer theDB.Close()

	//stmt, err := theDB.Prepare(fmt.Sprintf("DELETE FROM Artifacts where id=%s", record.Id))
	stmt, err := theDB.Prepare("DELETE FROM Artifacts WHERE id=?")
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err2 := stmt.Exec(record._ID)
	if err2 != nil {
		fmt.Println(err2)
		return false
	}

	//_, err = res.LastInsertId()

	return true
}

func getArtifactListDB() ([]ArtifactRecord, error) {

	var list []ArtifactRecord
	var record ArtifactRecord

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT ID, UUID, Name, Alias, Label, Checksum, Path, OpenChain, ContentType FROM Artifacts")

	if err != nil {
		//fmt.Println("error:", err)
		return list, err // return empty list
	}

	for rows.Next() {
		err = rows.Scan(&record._ID, &record.UUID, &record.Name, &record.Alias, &record.Label, &record.Checksum,
			&record._path, &record.OpenChain, &record.ContentType)
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		list = append(list, record)
	}
	rows.Close() //good habit to close
	return list, err
}

func getArtifactFromDB(field string, value string) ArtifactRecord {

	var record ArtifactRecord
	var query_str string

	openDB()
	defer theDB.Close()

	switch strings.ToLower(field) {
	case "id":
		query_str = fmt.Sprintf("SELECT ID, UUID, Name, Alias, Label, Checksum, Path, OpenChain, ContentType FROM Artifacts WHERE ID=%s", value)

	case "checksum":
		//
		fmt.Println()
	default:
		fmt.Println()
	}

	////stmt, err := theDB.Prepare
	rows, err := theDB.Query(query_str)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&record._ID, &record.UUID, &record.Name, &record.Alias, &record.Label,
			&record.Checksum, &record._path, &record.OpenChain, &record.ContentType)
		checkErr(err)
	}
	rows.Close() //good habit to close

	return record
}

// Insert alias value into the DB
func addAliasValueToDB(alias string, value string) {

	// TODO: return succes/failure status
	openDB()
	defer theDB.Close()

	sql_additem := `
	INSERT OR REPLACE INTO Aliases (
		Alias,
		Value,
		InsertedDatetime
		) values(?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err2 := stmt.Exec(alias, value)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func getAliasValueFromDB(alias string) (string, error) {
	var value string = ""
	/////var timestamp =""
	openDB()
	defer theDB.Close()

	//rows, err := theDB.Query("SELECT Value FROM Aliases WHERE Alias=?", alias)
	rows, err := theDB.Query(fmt.Sprintf("SELECT Value FROM Aliases WHERE Alias='%s'", alias))
	if err != nil {
		return "", err
	}

	count := 0
	for rows.Next() {
		count++
		err = rows.Scan(&value)
		if err != nil {
			return "", err
		}
	}
	rows.Close()
	if count == 0 {
		// alias not found in db.
		return "", errors.New(fmt.Sprintf("alias '%s' does not exist. Aliases are case sensitive", alias))
	} else {
		return value, nil
	}
}

func getAliasUsingValueFromDB(value string) (string, error) {
	var alias string = ""
	openDB()
	defer theDB.Close()

	rows, err := theDB.Query(fmt.Sprintf("SELECT Alias FROM Aliases WHERE Value='%s' ORDER BY InsertedDatetime ASC", value))
	if err != nil {
		return "", err
	}

	count := 0
	for rows.Next() {
		count++
		err = rows.Scan(&alias)
		if err != nil {
			return "", err
		}
	}
	rows.Close()
	if count == 0 {
		// alias not found in db.
		return "", errors.New(fmt.Sprintf("Alias for value: '%s' does not exist", value))
	} else {
		return alias, nil
	}
}

func getAlisaListFromDB() ([]AlisaRecord, error) {

	var list []AlisaRecord
	var record AlisaRecord

	openDB()
	defer theDB.Close()
	rows, err := theDB.Query("SELECT Alias, Value FROM Aliases ORDER BY Alias")

	if err != nil {
		return list, err // return empty list + erro
	}
	for rows.Next() {
		err = rows.Scan(&record.Alias, &record.Value)
		if err != nil {
			return nil, err // return empty list + error
		}
		list = append(list, record)
	}
	rows.Close()
	return list, err
}

func dumpDBTable(table_name string) (string, error) {

	openDB()
	defer theDB.Close()
	// Prepare statement to get the native types.
	stmt, err := theDB.Prepare(fmt.Sprintf("SELECT * FROM %s", table_name))

	if err != nil {
		return "", err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}

		tableData = append(tableData, entry)
	}

	jsonData99, err := json.Marshal(tableData)

	return string(jsonData99), nil
}

/******************************************************************************************************/
/******************************************************************************************************/
/******************************************************************************************************/

type Supplier_struct struct {
	Id         string `json:"id,omitempty"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	SKU_Symbol string `json:"sku_symbol"` // WR, INTEL, SAMSG
	Short_Id   string `json:"short_id"`   // WR-2e72b7
	Passwd     string `json:"passwd"`
	Type       string `json:"type"` // commercial, non-profit, individual
	Url        string `json:"url"`
}

type ledgerNode struct {
	Id          string `json:"id"`          // Primary key Id
	Name        string `json:"name"`        // Fullname
	ShortId     string `json:"short_id"`    //	1-5 alphanumeric characters (unique)
	IPAddress   string `json:"ip_address"`  // IP address - e.g., 147.11.153.122
	Port        int    `json:"port"`        // Port e.g., 5000
	UUID        string `json:"uuid"`        // 	UUID provide w/previous registration
	Label       string `json:"label"`       // 1-5 words display description
	Description string `json:"description"` // 2-3 sentence description
	Available   int    `json:"available"`   // 0 or 1 int (boolean)
}

func dumpTable2(table string) (string, error) {

	//rows, err := db.Query(sqlString)

	openDB()
	defer theDB.Close()
	////sqlString = fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := theDB.Query(fmt.Sprintf("SELECT * FROM %s", table))

	checkErr(err)

	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

// Add supplier record to the database
func AddSupplierToDB(s Supplier_struct) {

	openDB()
	defer theDB.Close()
	sql_additem := `
	INSERT OR REPLACE INTO Suppliers (
		UUID, 
		Name, 
		SKU_Symbol, 
		Alias, 
		Passwd, 
		Type, 
		Url, 
		InsertedDatetime
		) values(?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	stmt, err := theDB.Prepare(sql_additem)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	res, err2 := stmt.Exec(s.UUID, s.Name, s.SKU_Symbol, s.Short_Id, s.Passwd, s.Type, s.Url)
	if err2 != nil {
		panic(err2)
	}

	// for debugging - could probably be removed.
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error:")
		fmt.Println("last Id inserted", id)
		fmt.Println("Inserted")
		fmt.Println(s)
	}
}

// Get Supplier record from database.
func GetSupplier(db *sql.DB) []Supplier_struct {
	sql_readall := `
	SELECT Id, UUID, Name, SKU_Symbol, Short_Id, Passwd, Type, Url FROM Suppliers
	ORDER BY datetime(InsertedDatetime) DESC
	`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Supplier_struct
	for rows.Next() {
		item := Supplier_struct{}
		err2 := rows.Scan(&item.Id, &item.UUID, &item.Name, &item.SKU_Symbol, &item.Short_Id, &item.Passwd, &item.Type, &item.Url)
		if err2 != nil {
			panic(err2)
		}
		result = append(result, item)
	}
	return result
}

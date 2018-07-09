package main

// Common utilties are defined in this file.

/*
LICENSE NOTICE:
===============
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

/************** https://github.com/nu7hatch/gouuid *************
LICENSE NOTICE:
===============

/**
Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*****/

/**
Functions:
func FilenameDirectorySplit(full_file_path string) (string, string, string, string)

**/

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/nu7hatch/gouuid"
)

// Check error, and if true print error message.
// This is a commonly performed task that can be performed in one line.
func checkAndReportError(err error) bool {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return true
	} else {
		return false
	}
}

// This is a simple debug print command.
func here(i int, err error) {
	fmt.Printf("HERE: %d %s\n", i, err)
}

// This function exists to make error formating/display constient throughout the code
func displayErrorMsg(msg string) {
	fmt.Printf("error: %s\n", msg)
}

// Obtain abridge version of the file's path
func getAbridgedFilePath(fullpath string) (abridedPath string) {

	// For example: /C/Users/mitch/Documents/gospace/src/cmd/notices.pdf
	// 			To: ./notices.pdf
	// All relative to the sparts working directory

	// If string is empty nothing to do.
	if len(fullpath) == 0 {
		return "" // Path is empty - return empty
	}

	systemDir, _ := getSpartsDirectory() // C:/Users/mitch/Documents/gospace/src/cmd/.sparts
	path := fullpath
	pathDirectoryLength := len(path)
	count := 0
	//fmt.Println ("count: ", count)
	for len(path) == pathDirectoryLength && len(systemDir) > 3 {
		systemDir = filepath.Dir(systemDir)
		systemDir = strings.Replace(systemDir, `\`, `/`, -1)
		count++
		path = strings.TrimPrefix(path, systemDir)
	}
	if count == 1 {
		path = "." + path
	} else {
		// for count == 2
		path = strings.TrimPrefix(path, "/")
		// for count > 2
		for i := 1; i < count; i++ {
			path = "../" + path
		}
	}
	return path
}

//Create a line the length of the input string.
func createLine(str string) string {
	var buffer bytes.Buffer
	for l := 0; l < len(str); l++ {
		// --------------------------------
		buffer.WriteString("-")
	}
	return buffer.String()
}

// Create a white space string
func createWhiteSpace(size int) string {
	var buffer bytes.Buffer
	if size <= 0 {
		return ""
	}
	for l := 0; l < size; l++ {
		buffer.WriteString(" ")
	}
	return buffer.String()
}

// Check that the syntax of a UUID is correct.
func isValidUUID(uuid_str string) bool {
	_, err := uuid.ParseHex(uuid_str)
	return err == nil
}

// Get UUID - Unique Universal Identifier
func getUUID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	return u4.String()
}

// Checks if name starts with a "."
// Does not determine if file or directory is valid.
func isHidden(filename string) bool {
	dirPath, name, _, _ := FilenameDirectorySplit(filename)
	// if filename starts with "." or  directory path OR a parent directory THEN true
	if name[0:1] == "." ||
		(len(dirPath) > 0 && dirPath[0:1] == ".") ||
		strings.Contains(filename, "/.") {
		return true
	} else {
		return false
	}
}

func isPathURL(path string) bool {
	path = strings.ToLower(path)
	if strings.Contains(path, "http") {
		return true
	}
	return false
}

// Obtain the file count for a directory (and sub directory)
func getFileCount(path string) int {

	// find . -type f | wc -l  // recursively counts # files from a given directory
	count := 0

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//fmt.Printf(" Error: can't access directory: %s\n", directory)
			// don't count
			return nil
		}
		// ignore if directory or hidden file - files or directories that start with . or have a /. in
		if info.IsDir() || isHidden(path) {
			// It is a directory or a hidden file/directory
			//fmt.Printf("visited directory: %s\n", path)
			// don't count
		} else {
			count++
		}
		return nil
	}) // end of filepath.Walk

	return count
}

// Compute the SHA1 of a file.
func getFileSHA1(file string) (string, error) {

	hasher := sha1.New()
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(hasher, f); err != nil {
		fmt.Println(err)
		return "", err
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))
	return checksum, nil
}

// Break up a long string into chunkSize pieces.
func chunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}

// Parses a a full  file path name
// RETURNS: directory path, filename, file base name (name w/o ext), file extension
// Example: Path: "./d1/d2/d3/my_code.go"
//			Returns: "./d1/d2/d3/", "my_data.db", "my_data", ".db"
func FilenameDirectorySplit(full_file_path string) (dir_path string, base_name string, filename string, file_extension string) {

	filename = filepath.Base(full_file_path)
	file_extension = filepath.Ext(full_file_path)
	base_name = filename[:len(filename)-len(file_extension)]
	dir_path = full_file_path[:(len(full_file_path) - len(filename))]

	return dir_path, filename, base_name, file_extension
}

// Pretty print (format) the json reply.
func createJSONFormat(data interface{}) (string, error) {

	// We want to pretty print the json reply. We need to wrap:
	//    json.NewEncoder(http_reply).Encode(reply)
	// with the following code:

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "   ") // tells how much to indent "  " spaces.
	err := encoder.Encode(data)

	if err != nil {
		return "", err
	} else {
		return buffer.String(), nil
	}
}

// Check if a directory is the main working directory.
func isTopDirectory(directory string) bool {
	spartsDirectory := directory + "/" + _SPARTS_DIRECTORY
	if _, err := os.Stat(spartsDirectory); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Check if directory is the spart's data directory.
func isSpartsDirectory(directory string) bool {
	// check the
	name := filepath.Base(directory)
	if name != _SPARTS_DIRECTORY {
		return false
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Get the current directory.
func getCurrentDirectory() (directory string, err error) {

	directory, err = filepath.Abs(".")
	if err != nil {
		return "", err
	}
	// Clean up string on Windows platform replace '\' with '/'
	directory = strings.Replace(directory, `\`, `/`, -1)
	return directory, nil
}

//Accept terminal input for simple Yes/No reponses.
func getkeyboardYesNoReponse(msg string) bool {
	fmt.Printf("%s\n", msg)
	fmt.Print("[y/n] > ")
	// Read from standard input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.ToLower(scanner.Text())
	if answer == "y" || answer == "yes" {
		// They chose Yes.
		return true
	} else {
		// They chose No.
		return false
	}
}

//Accept terminal input for an open ended string reponse.
func getkeyboardReponse(msg string) string {
	fmt.Printf("%s\n", msg)
	fmt.Print("> ")
	// Read from standard input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// Find the parent sparts directory starting from a specified directory.
func getSpartsDirectory() (directory string, err error) {

	dir := getDirectory()
	previousDir := dir
	end_of_directory_path := false
	spartsDir := dir + "/" + _SPARTS_DIRECTORY

	for !isSpartsDirectory(spartsDir) {
		// Have not found the directory - go up one level
		previousDir = dir
		dir = filepath.Dir(dir)
		if !isDirectory(dir) || previousDir == dir {
			// either dir is no longer a directory or dir did not change
			// We are done walking back the path.
			end_of_directory_path = true
			break
		} else {
			// keep looking
			spartsDir = dir + "/" + _SPARTS_DIRECTORY
			//fmt.Println ("try next level:", dir)
		}
	}

	if end_of_directory_path {
		// did not find a sparts working directory
		return "", errors.New("Did not find valid " + _TOOL_NAME + " directory in path.")
	} else {
		// found directory
		// Clean up string on Windows replace '\' with '/'
		spartsDir = strings.Replace(spartsDir, `\`, `/`, -1)
		return spartsDir, nil
	}
}

func getType(object interface{}) string {

	// types have a "main." prefix. We need to remove it
	goType := strings.Replace(reflect.TypeOf(object).String(), "main.", "", 1)

	if strings.Contains(goType, "[]") {
		return strings.Replace(goType, "[]", "ListOf:", 1)
	} else {
		return goType
	}
}

// format a string for display such that no line is more than 'checkSize' long.
func formatDisplayString(DisplayStr string, chuckSize int) string {

	for len(DisplayStr) > chuckSize && DisplayStr[chuckSize] != ' ' {
		chuckSize++
	}
	chuckSize++
	//fmt.Println (chuckSize)
	//fmt.Println (chunkString(part.Description, chuckSize))
	return strings.Join(chunkString(DisplayStr, chuckSize), "\n                 ")
}

// create a directory
func createDirectory(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			panic(err)
		}
	}
}

/****
func createFile (file string) {
    f, err := os.Create(file)
    check(err)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    err := f.Write("This is a YAML file")
    if err != nil {
        panic(err)
    }
}
*****/

// Obtian the path to the sparts local config file
func getLocalConfigFile() string {
	var file string
	file = "./.sparts/config"
	fmt.Println("config file is:", file)
	return file
}

// Get path to current directory
func getDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return (dir)
}

// Determine if path provided is a directory.
func isDirectory(directory string) bool {
	candidate, err := os.Stat(directory)
	if err != nil {
		return false
	}

	//mode := fi.Mode()
	if candidate.Mode().IsDir() {
		return true
	} else {
		return false
	}
}

// Perform AES data encryption
func aesEncrypt(key []byte, message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	return
}

// Perform AES data decryption
func aesDecrypt(key []byte, securemess string) (decodedmess string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess = string(cipherText)
	return
}

// Compare to time values.
func compareTime(a, b time.Time) bool {
	secondsDiff := a.Unix() - b.Unix()

	if secondsDiff != 0 {
		return secondsDiff < 0
	}

	nanoDiff := a.Nanosecond() - b.Nanosecond()

	return nanoDiff <= 0
}

// True if A happened before B
func compareTimestamps(a, b string) (bool, error) {
	var ret bool
	layout := "2006-01-02 15:04:05.999999"

	timeA, err := time.Parse(layout, a)
	if err != nil {
		return ret, err
	} else if timeA.Unix() < 0 {
		return ret, fmt.Errorf("Time Overflow")
	}

	timeB, err := time.Parse(layout, b)
	if err != nil {
		return ret, err
	} else if timeB.Unix() < 0 {
		return ret, fmt.Errorf("Time Overflow")
	}

	ret = compareTime(timeA, timeB)

	return ret, nil
}

// reduce the lengh of a UUID to just the last 'szie' characters.
// Return as [<num characters]. Example: "[2f56e]" for size = 5.
func trimUUID(uuid string, size int) string {
	substring := uuid[len(uuid)-size : len(uuid)]
	return "[" + substring + "]"
}

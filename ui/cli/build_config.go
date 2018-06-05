package main

/*
	This file contains build configuration parameters for the sparts cli.
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

const (
	_VERSION            = "0.8"
	_DB_Model           = "0.8" // sqlite db data model
	_LOCAL_DB_FILE      = "status.0.8.db"
	_GLOBAL_CONFIG_FILE = ".spartsconfig"
	_LOCAL_CONFIG_FILE  = "config.yml"

	_SPARTS_DIRECTORY = ".sparts"
	_ENCRYPT_ID_STR   = "mY cRypTO String is nOT fIShy"
	_ENCRYPT_KEY      = "43a670c8-8f12-42"
	_TOOL_NAME        = "sparts"
)

//Runtime options
const (
	_DEBUG_DISPLAY_ON = false
)

const _LOCAL_CONFIG_FILE_CONTENT = `supply_chain: /zephyr
supplier_uuid: dde3e600-cd79-4ee5-464e-e74e1ce764bb
part_uuid: TBD
look_up: false
node:
    ledger_address: 147.11.176.111:818
    conductor_address: 10.37.133.106:811
`

// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
const (
	_CYAN_FG   = "\x1b[36;1m"
	_GREEN_FG  = "\x1b[32;1m"
	_RED_FG    = "\x1b[31;1m"
	_YELLOW_FG = "\x1b[36;1m"
	_WHITE_FG  = "\x1b[39;1m"
	_COLOR_END = "\x1b[0m"
)

// Error codes
const (
	_INITIALIZE_ERROR        = 1
	_DIR_ACCESS_ERROR        = 2
	_DB_ACCESSS_ERROR        = 3
	_DB_TABLE_CREATION_ERROR = 4
)

// Common value strings
const (
	_NONE         = "NONE"
	_NULL_PART    = "tbd"
	_ALIAS_TOKEN  = "id="
	_ALISA_LENGTH = 15
)

// Limits
const (
	_MAX_FILE_WARNING_COUNT = 49 // Give warning If an envelope is created with greater then max value # files
)

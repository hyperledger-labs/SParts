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

const _LOCAL_CONFIG_FILE_CONTENT = `look_up: false
part_uuid: tbd
node:
  ledger_address:
  conductor_address:
public_key:
private_key:
supply_chain: zephyr
supplier_uuid:
`

const _GLOBAL_CONFIG_FILE_CONTENT = `atlas: localhost:3075
user:
	email:
	name:
`

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
	_DEBUG_DISPLAY_ON  = false
	_DEBUG_REST_API_ON = false
)

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
	_ALIAS_LENGTH  = 15
	_ALIAS_TOKEN   = "id="
	_ENVELOPE_TYPE = "envelope"
	_LEDGER        = "ledger"
	_ATLAS         = "atlas"
	_NONE          = "NONE"
	_NULL_PART     = "tbd"
)

// API constants
const (
	_LISTOF_LEDGER_NODE          = "ListOf:LedgerNodeRecord"
	_ATLAS_LIST_LEDGER_NODES_API = "/atlas/api/v1/ledgerlist/"
)

// Limits
const (
	_MAX_FILE_WARNING_COUNT = 49 // Give warning If an envelope is created with greater then max value # files
)

// rest_api.go variables
const (
	// Atlas directory look up
	_ATLAS_PING_API = "/atlas/api/v1/ping"

	// Ledger
	_ARTIFACTS_API         = "/ledger/api/v1/artifacts"
	_LEDGER_PING_API       = "/ledger/api/v1/ping"
	_PARTS_API             = "/ledger/api/v1/parts"
	_PARTS_TO_SUPPLIER_API = "/ledger/api/v1/parts/supplier"
	_SUPPLIERS_API         = "/ledger/api/v1/suppliers"
	_ARTIFACTS_URI_API     = "/ledger/api/v1/artifacts/uri"
)

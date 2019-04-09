package main

/*
	This file contains compiler build configuration parameters for the sparts cli.
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

//Runtime debug flags
const (
	_DEBUG_DISPLAY_ON  = false
	_DEBUG_REST_API_ON = false
	_SEED_FUNCTION_ON  = true
)

// This is the default settings fot the local configuration
const _LOCAL_CONFIG_FILE_CONTENT = `auto_synch: false
envelope_uuid: <tbd>
focus: BOTH
ledger_address:
part_uuid: <tbd>
private_key:
public_key:
ledger_network: sparts-test-network
supplier_uuid:
`

// This is the default settings fot the global configuration
const _GLOBAL_CONFIG_FILE_CONTENT = `
atlas_address: https://spartshub.org
user_email
user_name
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

// Const values to display color to a terminal screen
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
	_INITIALIZE_ERROR = 1
	_DIR_ACCESS_ERROR = 2
	_DB_ACCESSS_ERROR = 3
)

// Misc  values
const (
	_ALIAS_LENGTH      = 15
	_ALIAS_TOKEN       = "id="
	_INDENT_STR        = "   "
	_ENVELOPE_TYPE     = "envelope"
	_LEDGER            = "ledger"
	_ATLAS             = "atlas"
	_NONE              = "NONE"
	_ENVELOPE_FOCUS    = "ENVELOPE"
	_PART_FOCUS        = "PART"
	_BOTH_FOCUS        = "BOTH"
	_NO_FOCUS          = "NONE"
	_TRUE              = "true"
	_FALSE             = "false"
	_PRE_LEDGER_TOKEN  = ">"
	_POST_LEDGER_TOKEN = "="
	_NULL_UUID         = "<tbd>"
	_ROOT_ENV          = "root"
	_ROOT_TOKEN        = "root:"
)

// Limits
const (
	_MAX_FILE_WARNING_COUNT = 49 // Give warning If an envelope is created with greater then max value # files
)

// rest_api.go variables
const (
	// Atlas (spartshub.org) directory look up
	_ATLAS_PING_API              = "/atlas/api/v1/ping"
	_ATLAS_LIST_LEDGER_NODES_API = "/atlas/api/v1/network_node_list/"
	_NETWORK_LIST_API            = "/atlas/api/v1/network_space"

	// Ledger
	_ARTIFACTS_API        = "/ledger/api/v1/artifacts"
	_ARTIFACTS_URI_API    = "/ledger/api/v1/artifacts/uri"
	_ARTIFACT_OF_ENV_API  = "/ledger/api/v1/envelope/artifact"
	_ARTIFACT_OF_PART_API = "/ledger/api/v1/artifacts/part"
	_KEY_PAIR_API         = "/ledger/api/v1/keys"
	_PARTS_API            = "/ledger/api/v1/parts"
	_LEDGER_PING_API      = "/ledger/api/v1/ping"
	//_PARTS_TO_SUPPLIER_API = "/ledger/api/v1/parts/supplier"
	_PARTS_TO_SUPPLIER_API = "/ledger/api/v1/parts/org"
	// _SUPPLIER_API          = "/ledger/api/v1/suppliers"
	_ORGS_API          = "/ledger/api/v1/orgs"
	_REGISTER_USER_API = "/ledger/api/v1/registeruser"
)

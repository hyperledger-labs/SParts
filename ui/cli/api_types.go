package main

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
	_SUCCESS     = "success"
	_FAILURE     = "failed"
	_SUCCESS_MSG = "Ok"
)

type ArtifactRecord struct {
	UUID         string         `json:"uuid"`
	Name         string         `json:"name"`
	Alias        string         `json:"short_id,omitempty"`
	Label        string         `json:"label,omitempty"` // Display name
	Checksum     string         `json:"checksum"`
	OpenChain    string         `json:"openchain,omitempty"`
	ContentType  string         `json:"content_type,omitempty"`
	Timestamp    string         `json:"timestamp,omitempty"`
	ArtifactList []ArtifactItem `json:"artifact_list,omitempty"`
	URIList      []URIRecord    `json:"uri_list, omitempty"`
	// Internal use only
	_ID       int
	_path     string
	_verified bool
}

type ArtifactOfEnvelopeRecord struct {
	PublicKey  string               `json:"public_key"`
	PrivateKey string               `json:"private_key"`
	Relation   ArtifactEnvelopePair `json:"relation"`
}

type ArtifactEnvelopePair struct {
	ArtifactUUID string `json:"artifact_uuid"`
	EnvelopeUUID string `json:"envelope_uuid"`
}

type ArtifactOfPart struct {
	PublicKey  string           `json:"public_key"`
	PrivateKey string           `json:"private_key"`
	Relation   ArtifactPartPair `json:"relation"`
}

type ArtifactPartPair struct {
	ArtifactUUID string `json:"artifact_uuid"`
	PartUUID     string `json:"part_uuid"`
}

type EmptyRecord struct {
}

var _EMPTY_RECORD EmptyRecord

type EnvelopeArtifactRecord struct {
	PrivateKey string               `json:"private_key"`
	public_key string               `json:"public_key"`
	relation   EnvelopeArtifactPair `json:"relation"`
}
type EnvelopeArtifactPair struct {
	EnvelopeUUID string `json:"envelope_uuid"`
	ArtifactUUID string `json:"artifact_uuid"`
}

type ArtifactAddRecord struct {
	PublicKey  string         `json:"public_key"`
	PrivateKey string         `json:"private_key"`
	Artifact   ArtifactRecord `json:"artifact"`
}

type ArtifactItem struct {
	UUID string `json:"uuid"` // Artifact Universal Unique IDentifier
	Path string `json:"path"` // Path of artifact within the envelope
}

type PartAddRecord struct {
	PublicKey  string     `json:"public_key"`
	PrivateKey string     `json:"private_key"`
	Part       PartRecord `json:"part"`
}

type PartToSupplierRecord struct {
	PublicKey  string           `json:"public_key"`
	PrivateKey string           `json:"private_key"`
	Relation   PartSupplierPair `json:"relation"`
}

type PartSupplierPair struct {
	PartUUID     string `json:"part_uuid"`
	SupplierUUID string `json:"supplier_uuid"`
}

type PartRecord struct {
	UUID         string           `json:"uuid"`                  // Unique identifier
	Name         string           `json:"name"`                  // Fullname
	Version      string           `json:"version,omitempty"`     // Version if exists.
	Label        string           `json:"label,omitempty"`       // Display name
	Alias        string           `json:"alias,omitempty"`       // 1-15 alphanumeric characters (unique)
	Licensing    string           `json:"licensing,omitempty"`   // License expression
	Description  string           `json:"description,omitempty"` // Part description (1-3 sentences)
	Checksum     string           `json:"checksum,omitempty"`    // License expression
	ArtifactList []ArtifactItem   `json:"artifacts,omitempty"`
	Suppliers    []SupplierRecord `json:"suppliers,omitempty"`
	Categories   []CategoryRecord `json:"categories,omitempty"`
	// URIList      []URIRecord    `json:"uri_list,omitempty"`
}

type PartItemRecord struct {
	PartUUID string `json:"part_id"` // Part uuid
}

/******
type PartSupplierRelateRecord struct {
	PublicKey  string     `json:"public_key"`
	PrivateKey string     `json:"private_key"`
	relation   PartRecord `json:"part"`
}
*******/

type ReplyType struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Type    string      `json:"result_type"`
	Result  interface{} `json:"result,omitempty"`
}

type CategoryRecord struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"` // Part description (1-3 sentences)
}

type UserRecord struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Authorized string `json:"authorized"`
	PublicKey  string `json:"public_key"`
}

type URIRecord struct {
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	ContentType string `json:"content_type"`   // text, envelope, binary, archive
	Size        string `json:"size,omitempty"` // size in bytes
	URIType     string `json:"uri_type"`       // e.g., http, ipfs
	Location    string `json:"location"`       // actual link
}

type URIAddRecord struct {
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
	UUID       string    `json:"uuid"`
	URI        URIRecord `json:"uri"`
}

type SupplierRecord struct {
	UUID  string `json:"uuid"`            // UUID provide w/previous registration
	Name  string `json:"name"`            // Fullname
	Alias string `json:"alias,omitempty"` // 1-15 alphanumeric characters
	Url   string `json:"url,omitempty"`   // 2-3 sentence description
	Parts []PartItemRecord
}

// ================================================
// 			Atlas API Types
//=================================================

type LedgerNodeRecord struct {
	UUID        string `json:"uuid"`                  // UUID
	Name        string `json:"name"`                  // Fullname
	NetworkName string `json:"network_name"`          // Network Space name
	Alias       string `json:"alias,omitempty"`       // 1-15 alphanumberic alias
	APIURL      string `json:"api_url"`               // e.g., http://147.52.17.33:5000
	PublicKey   string `json:"public_key,omitempty"`  // Public key to verify authorization
	Description string `json:"description,omitempty"` // 2-3 sentence description
	Status      string `json:"status,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

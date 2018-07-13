# Ledger Node API

## Overview

The API for the SParts ledger is presented here. The ledger API calls are defined in part I of this document. The record types (objects) past between the ledger and client application are defined in part II of this document. Types include supplier, part, category as so forth. 

## I) Ledger API Calls

#### +Ping Request

------

Send request to see if the ledger is currently available.

```
GET /ledger/api/v1/ping
```

Example of a successful response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "EmptyRecord",
	result: 	{}
}
```

Since there is not data to return the record type **EmptyRecord** is specified in the results field.  **EmptyRecord** is defined in part II. If the ledger is not available then no response will be received.



#### +Artifact Record

```
GET /ledger/api/v1/artifacts/{uuid}
```

(This call use to be: /api/ledger/envelopes/{uuid})

An artifact represents an item of evidence. Typically an artifact is a single document (e.g., notice file, source code archive, bill of materials). An envelope is a special instance of an artifact which represents a collection of artifacts potentially including  other envelopes. For single artifacts the artifact_list field will be empty. For an envelope it will contain a list of zero of more artifacts and the content_type field will be set to "envelope". The uri_list field is a list because copies of the artifact could exist in multiple locations.

Response form:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ArtifactRecord",
	result: 	{
		name: "...",
		uuid: "...",
		filename: "...",
		checksum: "...",
		content_type: "...",  // envelope, notice, source, spdx, doc, other
		alias: "...",
		label: "...",
		openchain: "...",
		timestamp: "..."
		artifact_list: [...]   /* used for envelopes but not for singular artifact */
		uri_list: [ {
          			version: "...",
							  alias: "...",
							  checksum: "...",
							  size:	"..."
							  content_type: "...",  // http, ipfs, ...
							  location: "https://...."
						   }
						 ]
				}
}
```

Example of a <u>single</u> artifact response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ArtifactRecord",
	result: {
		name: "Zephyr 1.12 Notice File",
		uuid: "26559ed4-6868-488d-a5a7-3e81714beb00",
		filename: "Zephyr-1.12-Notices.txt",
		checksum: "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
		content_type: "notices",
		alias: "zephyr-notices-1.12",
		label: "Zephyr Notices 1.12",
		openchain: "True",
		timestamp: "2018-06-18 00:30:12.498167"
		artifact_list: []   /* not used for singular artifact */
		uri_list: [ {
  			version: "1.0",
			   alias: "zephyr-notices-1.12",
			   checksum: "Zephyr Notices 1.12",
			   size:	"235120"
			   content_type: "http",
			   location: "https://...."
				   }
				 ]
}
```

Example of an <u>envelope</u> response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ArtifactRecord",
	result: {    name: "Zephyr 1.12 Envelope",
				uuid: "9b602058-c73f-4f02-9237-b71a2760fc15",
				filename: "Zephyr-1.12-envelope.zip",
				checksum: "a1e2486417f4cd7fc670bf5facd5870af9c1e3a5",
				content_type: "envelope",
				alias: "zephyr-notices-1.12",
				label: "Zephyr Notices 1.12",
				openchain: "True",
				timestamp: "2018-06-18 00:30:12.498167"
				artifact_list: [
                    		{ uuid: "731ef148-5f81-11e8-9c2d-fa7ae01bbebc",
                               path: "/spdx"},
						   { uuid: "f2cef148-5f81-11e8-8f51-fa7ae01bb93b",
						     path: "/notices"}
				] 
				uri_list: [ {
                    			version: "1.0",
							   alias: "zephyr-envelope-1.12",
							   checksum: "f67d3213907a52012a4367d8ad4f093b65abc016",
							   size:	"235120"
							   content_type: "http",
							   location: "https://...."
						   }
						 ]
}
```

Note that the envelope record utilizes the artifact_list field where a single artifact does not. 



#### Artifact Add*

```
POST /ledger/api/v1/artifacts
```

Use the **ArtifactRecord** (the artifact_list  and uri_list fields are not used in this post). The request must be performed by a user with Roles: admin or supplier.

| Field                | Type   | Description                                                  |
| -------------------- | ------ | ------------------------------------------------------------ |
| uuid                 | string | unique identifier                                            |
| name                 | string | file or envelope name                                        |
| alias (was short_id) | string | alias for typing                                             |
| label                | string | // Display name                                              |
| checksum             | string | artifact checksum                                            |
| openchain            | string | true/false If prepared under an OpenChain comforting program |
| content_type         | string | envelope, notices, spdx, source, ...                         |

An example single artifact request:

```
{
   private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   artifact: {	
   			  uuid: "7709ca8d-01f4-4de2-69ed-16b7ebae704a",
    		  name: "Zephypr 1.12 SPDX file",
			  alias: "zephypr_1.12",
			  label: "Zephypr 1.12 SPDX file",
			  checksum: "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
			  openchain: true,
			  content_type: "spdx"
		}
}
```

Example curl Request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", "artifact": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c","name": "Zephypr 1.12 SPDX file", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.12", "label": "Zephypr 1.12 SPDX file", "openchain": "true", "content_type": "spdx"} }' http://147.11.176.111:818/ledger/api/v1/artifacts
```



**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the add.
- One or more of the required fields UUID, checksum are missing.
- The UUID is not in a valid format. 



#### Artifact URI Add*

```
POST /ledger/api/v1/artifacts/uri
```

The request must be performed by a user with Roles: Admin or Supplier.

| Field        | Type   | Description                               |
| ------------ | ------ | ----------------------------------------- |
| version      | string | name of use                               |
| checksum     | string | artifact checksum                         |
| content_type | string | type (e.g., text, binary, archive, other) |
| size         | int    | file size in bytes                        |
| uri_type     | string | (e.g., http, ipfs, ...)                   |
| location     | string | link, path                                |

An example request:

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   uuid: "bcb083a1-89c7-4bd2-a568-8450350e8195",
   uri:  { version: "1.0",
		   checksum: "f67d3213907a52012a4367d8ad4f093b65abc016",
		   size:	"235120"
		   content_type: ".pdf",
		   uri_type: "http",
		   location: 	  
		      "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016"
		}
}
```

The uri field is of type **URIRecord**.

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", 
"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c", "uri": {"version": "1.0", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "size": "235120", "content_type": ".pdf", "uri_type": "http", "location": "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016" } }' http://147.11.176.111:818/ledger/api/v1/artifacts/uri
```



#### Artifact Of Envelope Relation*

```
POST /ledger/api/v1/envelope/artifact
```

The user is identified by the public key. The following input record is: **EnvelopeArtifactRecord**

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   relation: { 
   				artifact_uuid: "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
   				envelope_uuid: "dec6b86a-f794-43d6-64bc-ca4146548048"
		  	}
}
```

Example curl request

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "relation": {"artifact_uuid": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "envelope_uuid": "dec6b86a-f794-43d6-64bc-ca4146548048"} }' http://localhost:3075/ledger/api/v1/envelope/artifact
```

**To do:**

- Ledger check for envelope to be an envelope or send back error
- Need to use uuid for user. - are public keys unique??

POST /ledger/api/v1/relation/part_artifact
POST /ledger/api/v1/relation/category_part



#### Artifact of Part Relation*

```
POST /ledger/api/v1/artifacts/part
```

 **ArtifactOfPartRecord** 

Example request:

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   relation: { 
   				part_uuid: "cb00a696-14d0-447a-6096-69be4c5d93a5",
   				artifact_uuid: "2745a756-eed8-4093-683f-a1f6b56f7249"
		  	}
}
```

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "relation": {"part_uuid": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "artifact_uuid": "dec6b86a-f794-43d6-64bc-ca4146548048"} }' http://localhost:3075/ledger/api/v1/artifacts/part
```



#### Supplier Record

```
GET /ledger/api/v1/suppliers/{uuid}
```

Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "SupplierRecord",
	result: 	{    name: "...",
					uuid: "...",
					alias: "...",
					url: "...",
					parts: [...]
				}
}
```

Example Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "SupplierRecord",
	result: {	name: "Wind River Systems",
			 	uuid: "dde3e600-cd79-4ee5-464e-e74e1ce764bb",
				alias: "WindRiver",
				parts: [
					{ part_id: "731ef148-5f81-11e8-9c2d-fa7ae01bbebc" },
					{ part_id: "f2cef148-5f81-11e8-8f51-fa7ae01bb93b" },
					{ part_id: "ee3b9d57-4c98-4d6f-5ecb-a54c97a7cda2" }
				],
				url: "http://www.windriver.com"
			}
}
```

See section II for the json definition of **OrganizationRecord**. If an organization is not found the following error response will be received:

```
{	status: 	"failed",
	message: 	"Organization record not found",
	result_type: "EmptyRecord",
	result: 	{}
}
```



#### Supplier List

Obtain a list of organizations. 

```
GET /ledger/api/v1/suppliers
```

Example Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ListOf:SupplierRecord",
	result: [ {	name: "Wind River Systems",
			 	uuid: "dde3e600-cd79-4ee5-464e-e74e1ce764bb",
				alias: "WindRiver",
				parts: [
					{ part_uuid: "731ef148-5f81-11e8-9c2d-fa7ae01bbebc" },
					{ part_uuid: "f2cef148-5f81-11e8-8f51-fa7ae01bb93b" },
					{ part_uuid: "ee3b9d57-4c98-4d6f-5ecb-a54c97a7cda2" }
				],
				url: "http://www.windriver.com"
			},
			{	name: "Intel Corporation",
			 	uuid: ""2584a6ce-16a7-44c0-7e53-21969d1e026b",
				alias: "Intel",
				parts: [
					{ part_uuid: "6584a6ce-16a7-44c0-7e53-21969d1e026b" },
				],
				url: "http://www.intel.com"
			}
	]
}
```

See section II for the json definition of **OrganizationRecord**. If there are not organizations registered then the empty list will be received:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ListOf.SupplierRecord",
	result: [ ]
}
```



#### Supplier Add*

```
POST /ledger/api/v1/suppliers
```

The request must be performed by a user with Role:Admin access.

| Name  | Type    | Description                                           |
| ----- | ------- | ----------------------------------------------------- |
| uuid  | string  | universal unique identifier                           |
| name  | string  | name of organization                                  |
| alias | strings | quick identifier 1-15 alphanumeric characters website |
| url   | string  | website url                                           |

Example Request:

```
{
	uuid: "dde3e600-cd79-4ee5-464e-e74e1ce764bb",
	name: "WindRiver",
	alias: "supplier-76bb1",
	url: "http://www.windriver.com"
}
```

Example curl command:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"name": "Wind River", "uuid": "3568f20a-8faa-430e-7c65-e9fce9aa155d", "alias": "wr", "url": "http://www.windriver.com"}' http://localhost:3075/ledger/api/v1/orgs
```



#### Part Get

To obtain a part record use:

```
GET /ledger/api/v1/parts/{uuid}
```

Response Record

```
{	status: 	"success",
	message: 	"OK",
	result_type: "PartRecord",
	result: {    uuid: "...",
	  			name: "...",
	  			version: "...",
	  			label: "..."
	  			alias: "..."
	  			checksum: "...",
	  			label: "...",
	  			licensing: "..."
	  			description: "..."
	  			artifacts: [...]
	  			suppliers: [...]
	  			categories: [...]
	}
}
```

Example response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "PartRecord",
	result: {	 uuid: "f1eae70d-86ba-4440-583a-28127e447f83",
				name: "Zephyr Runtime 1.10",
				version: "1.10",
				label: "zephyr 1.10"
				alias: "z1
				description: "Zephyr runtime for the ACX 11 board support bpackage"
				licensing: "Apache-2.0",
				checksum: "d9be5fcf820e88b217a760f7869959af49898dbe",
				suppliers: [
						      { supplier_id: "3568f20a-8faa-430e-7c65-e9fce9aa155d" }
						   ],
				artifacts: [ ],
				categories: [ ],
	}
```

Error Messages include:

- Part uuid not found



#### Part List

```
GET /ledger/api/v1/parts
```

Returns a list of **PartRecords**. Example response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ListOf.PartRecord",
	result: []
}
```



#### Part Add*

```
POST /ledger/api/v1/parts
```

Example:

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   part: { 
			uuid: "731ef148-5f81-11e8-9c2d-fa7ae01bbebc",
	  		name: "Zephyr Runtime 1.10",
	  		label: "Zephyr 1.10"
	  		alias: "zephyr_1.10"
	  		version: "1.10",
	  		checksum: "1e673213907a52012a4367d8ad4f093b65abc222",
	  		label: "Test Part 1.0",
	  		licensing: "MIT"
	  		description: "Zephyr is a small real-time operating supporting multiple architectures"
		 }
}
```

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "part": {"uuid": "731ef148-5f81-11e8-9c2d-fa7ae01bbebc","name": "Zephypr 1.12 SPDX file", "version": "1.12", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.10", "label": "Zephyr 1.10", "licensing": "MIT", "description": "Zephyr is a small real-time operating supporting multiple architectures" } }' http://localhost:3075/ledger/api/v1/parts
```



#### Part Of Supplier Relation*

```
POST /ledger/api/v1/parts/supplier
```

 PartOfSupplierRecord 

Example request:

```
{  private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   relation: { 
   				supplier_uuid: "cb00a696-14d0-447a-6096-69be4c5d93a5",
   				part_uuid: "2745a756-eed8-4093-683f-a1f6b56f7249"
		  	}
}
```

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "relation": {"part_uuid": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "supplier_uuid": "dec6b86a-f794-43d6-64bc-ca4146548048"} }' http://localhost:3075/ledger/api/v1/parts/supplier
```



#### Get Public/Private Key Pair

------

```
GET /ledger/api/v1/keys
```

Example Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "PrivatePublicKeyRecord",
	result: {
			  private_key: "5K6Q2kMHaMUrRjvSb4EPXQEJi1nAy3uXAYhqYvq2qNiLEGuFuVS",
			  public_key: "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"
	}
}
```



#### User Add 

------

```
POST /ledger/api/v1/users
```

Here is an example of how to register a new user. 

The request must be performed by a user with Role:Admin.

| Name       | Type   | Description                           |
| ---------- | ------ | ------------------------------------- |
| name       | string | name of user                          |
| email      | string | The email address of user             |
| role       | string | specific role value [admin, supplier] |
| authorized | string | specific access value [allow, deny]   |
| public_key | string | public key value                      |

Example Request:

```
{	name: 	"John Doe",
	email: 	"john.doe@windriver.com",
	role: "admin",
	authorized: "allow",
	public_key: "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"
}
	
```

Example response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "EmptyRecord",
	result: {}
}
```

Example curl request:

```
 curl -i -H "Content-Type: application/json" -X POST -d  '{"name": "John Doe", "email": "john.doe@windriver.com", "role": "admin", "authorized": "allow", "public_key": "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"}' http://localhost:3075/ledger/api/v1/users

```



#### User Get

```
GET /ledger/api/v1/users/{public_key}
```



```
{	name: 	"Sameer Ahmed",
	email: 	"sameer.ahmed@windriver.com",
	organization: "Wind River"
	public_key: "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"
}
```



#### GET Category

------

```
GET /ledger/api/v1/categories/{uuid}
```

 Example Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "CategoryRecord",
	result: {	
		name: "OS",
		description: "Operating System",
		uuid: "43903f02-00fd-43a3-bdaa-befe4a2fcd7e"
	}
}
```

An example error response where the category uuid does not exist:

```
{	status: 	"failed",
	message: 	"Category record not found",
	result_type: "EmptyRecord",
	result: 	{}
}
```



#### List Categories

------

Obtain a list of the categories.

```
GET /ledger/api/v1/categories
```

Example Response:

```
{	status: 	"success",
	message: 	"OK",
	result_type: "ListOf.CategoryRecord",
	result: [ {	name: "operating system",
			    description: "Operating System",
				uuid: "43903f02-00fd-43a3-bdaa-befe4a2fcd7e"
			 },
			 {	name: "libraries",
			    description: "Operating System",
				uuid: "12d7f0c2-00fd-43a3-bdaa-befe4a2fcd7e"
			 }
	]
}
```



## II) API Types

#### ArtifactRecord

```
{ 
	UUID         string         `json:"uuid"`
    Name         string         `json:"name"`     
	Alias        string         `json:"short_id,omitempty"` 
	Label        string         `json:"label,omitempty"`  // Display name
	Checksum     string         `json:"checksum"`             
	OpenChain    string         `json:"openchain,omitempty"`  
	ContentType  string         `json:"content_type,omitempty"`  
	Timestamp    string         `json:"timestamp,omitempty"`     
	ArtifactList ListOf.ArtifactItem `json:"artifact_list,omitempty"` 
	URIList      ListOf.URIRecord    `json:"uri_list, omitempty"`     
}


```



#### ArtifactItem

```
{
	UUID string `json:"uuid"` // Artifact Universal Unique IDentifier
	Path string `json:"path"` // Path of artifact within the envelope
}
```



#### EmptyRecord

```
{ }
```

#### EnvelopeArtifactRecord

```

```







#### SupplierRecord

```
{
	UUID    string `json:"uuid"`               // UUID provide w/previous registration
	Name    string `json:"name"`               // Fullname
	Alias   string `json:"alias,omitempty"`    // 1-15 alphanumeric characters
	Url     string `json:"url,omitempty"`      // 2-3 sentence description
	Parts   ListOf.PartItemRecord
}
```



#### PartItemRecord 

```
{
	PartUUID string `json:"part_uuid"` // Part uuid
}
```



#### PartRecord

```
{
	Name        string `json:"name"`                  // Fullname
	Version     string `json:"version,omitempty"`     // Version if exists.
	Alias       string `json:"label,omitempty"`       // 1-15 alphanumeric characters
	Licensing   string `json:"licensing,omitempty"`   // License expression
	Description string `json:"description,omitempty"` // Part description (1-3 sentences)
	Checksum    string `json:"checksum,omitempty"`    // License expression
	UUID        string `json:"uuid"`                  // UUID provide w/previous registration
	URIList     []URIRecord `json:"uri_list,omitempty"`     //
}
```



#### PrivatePublicKeyRecord

```
{
	PrivateKey string `json:"private_key"` // Private key
	PublicKey  string `json:"public_key"`  // Pubkic key
}
```



#### URIRecord

```
{
	Version     string `json:"version"`
	Checksum    string `json:"checksum"`
	ContentType string `json:"content_type"`   // text, envelope, binary, archive
	Size        string `json:"size,omitempty"` // size in bytes
	URIType     string `json:"uri_type"`       // e.g., http, ipfs
	Location    string `json:"location"`       // actual link
}
```



#### 

<style>
	table {
	  font-family: arial, sans-serif;
	  border-collapse: collapse;
	  width: 100%;
	}
	
	td, th {
	  border: 1px solid #dddddd;
	  text-align: left;
	  padding: 8px;
	}
	
	tr:nth-child(even) {
	  background-color: #dddddd;
	}
</style>

# Ledger Node API


## I) Overview

The API for the SParts ledger is presented here. The ledger API calls are defined in part I of this document. The record types (objects) past between the ledger and client application are defined in part II of this document. Types include supplier, part, category as so forth.


## II) Ledger API Calls


### Ping Request
----

Send request to see if the ledger is currently available.

```
GET /ledger/api/v1.1/ping
```

Example of a successful response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "EmptyRecord",
	result: 	{}
}
```

Since there is no data to return the record type **EmptyRecord** is specified in the results field.  **EmptyRecord** is defined in part II of this document. If the ledger is not available then no response will be received.


### Artifact Family
----

An artifact represents an item of evidence. Typically, an artifact is a single document (e.g., notice file, source code archive, bill of materials). An envelope is a special instance of an artifact which represents a collection of artifacts potentially including  other envelopes. For single artifacts the artifact_list field will be empty. For an envelope it will contain a list of zero of more artifacts and the content_type field will be set to "envelope". The uri_list field is a list because copies of the artifact could exist in multiple locations.


#### Artifact Create

```
POST /ledger/api/v1.1/artifacts
```

Allows the user to create an artifact into the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
    <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>file or envelope name</td>
        </tr>
        <tr>
            <td>alias (was short_id)</td>
            <td>string</td>
            <td>alias for typing</td>
        </tr>
        <tr>
            <td>label</td>
            <td>string</td>
            <td>// Display name</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>artifact checksum</td>
        </tr>
        <tr>
            <td>openchain</td>
            <td>string</td>
            <td>true/false If prepared under an OpenChain comforting program</td>
        </tr>
        <tr>
            <td>content_type</td>
            <td>string</td>
            <td>envelope, notices, spdx, source, ...</td>
        </tr>
    </tbody>
</table>
<br>

Example of single artifact request:

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

Example of curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", "artifact": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c","name": "Zephypr 1.12 SPDX file", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.12", "label": "Zephypr 1.12 SPDX file", "openchain": "true", "content_type": "spdx"} }' http://147.11.176.111:818/ledger/api/v1.1/artifacts
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the create.
- One or more of the required fields are missing.
- The UUID is not in a valid format. 
- The UUID is not unique to the artifact.


#### Artifact Amend

```
POST /ledger/api/v1.1/artifacts/amend
```

Allows the user to amend an artifact in the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
    <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>file or envelope name</td>
        </tr>
        <tr>
            <td>alias (was short_id)</td>
            <td>string</td>
            <td>alias for typing</td>
        </tr>
        <tr>
            <td>label</td>
            <td>string</td>
            <td>// Display name</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>artifact checksum</td>
        </tr>
        <tr>
            <td>openchain</td>
            <td>string</td>
            <td>true/false If prepared under an OpenChain comforting program</td>
        </tr>
        <tr>
            <td>content_type</td>
            <td>string</td>
            <td>envelope, notices, spdx, source, ...</td>
        </tr>
    </tbody>
</table>
<br>

Example of single artifact request:

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

	or

{
   private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   artifact: {	
		uuid: "7709ca8d-01f4-4de2-69ed-16b7ebae704a",
	  	name: "Zephypr 1.12 SPDX file",
	  	alias: "zephypr_1.12"
	}
}
```

(Note: The payload is valid as long as "private_key", "public_key", "artifact" and "uuid" are present.)

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", "artifact": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c","name": "Zephypr 1.12 SPDX file", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.12", "label": "Zephypr 1.12 SPDX file", "openchain": "true", "content_type": "spdx"} }' http://147.11.176.111:818/ledger/api/v1.1/artifacts/amend

	or

curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K92SiHianMJRtqRiMaQ6xwzuYz7xaFRa2C8ruBQT6edSBg87Kq", "public_key" : "02be88bd24003b714a731566e45d24bf68f89ede629ae6f0aa5ce33baddc2a0515", "artifact": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae705c","name": "Zephypr 1.12 SPDX file", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.12"} }' http://147.11.176.111:818/ledger/api/v1.1/artifacts/amend
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the create.
- No fields were amended.
- The UUID is not in a valid format. 
- The UUID does not exist.


#### Artifact List

```
GET /ledger/api/v1.1/artifacts
```

Allows the user to obtain a list of artifact from the sPart ledger. This is a public function.

Example of list artifact response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "ArtifactRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

If there are no artifacts registered, then an empty list will be returned as shown:

```
{   
    status:     "success",
    message:    "OK",
    result_type: "ArtifactRecord",
    result:     []
}
```

**Potential Errors**:

- Data cannot be deserialized due to encoding error.


#### Artifact Retrieve

```
GET /ledger/api/v1.1/artifacts/{uuid}
```

Allows the user to obtain the artifact data associating with the uuid from the sPart ledger. This is a public function.

Example of a <u>single</u> artifact response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "ArtifactRecord",
	result: {
		name: "Zephyr 1.12 Notice File",
		uuid: "26559ed4-6868-488d-a5a7-3e81714beb00",
		checksum: "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
		content_type: "notices",
		alias: "zephyr-notices-1.12",
		label: "Zephyr Notices 1.12",
		openchain: "True",
		timestamp: "2018-06-18 00:30:12.498167"
		artifact_list: []   /* not used for singular artifact */
		uri_list: [ 
			{
				version: "1.0",
		   		alias: "zephyr-notices-1.12",
			   	checksum: "Zephyr Notices 1.12",
			   	size:	"235120"
			   	content_type: "http",
			   	location: "https://...."
			}
		]
	}
}
```

Example of an <u>envelope</u> response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "ArtifactRecord",
	result: {   
		name: "Zephyr 1.12 Envelope",
		uuid: "9b602058-c73f-4f02-9237-b71a2760fc15",
		checksum: "a1e2486417f4cd7fc670bf5facd5870af9c1e3a5",
		content_type: "envelope",
		alias: "zephyr-notices-1.12",
		label: "Zephyr Notices 1.12",
		openchain: "True",
		timestamp: "2018-06-18 00:30:12.498167"
		artifact_list: [
    		{
    			uuid: "731ef148-5f81-11e8-9c2d-fa7ae01bbebc",
            	path: "/spdx"
    		},
			{
				uuid: "f2cef148-5f81-11e8-8f51-fa7ae01bb93b",
			    path: "/notices"
			}
		] 
		uri_list: [ 
			{
    			version: "1.0",
				alias: "zephyr-envelope-1.12",
				checksum: "f67d3213907a52012a4367d8ad4f093b65abc016",
				size:	"235120"
				content_type: "http",
				location: "https://...."
			}
		]
	}
}
```

Note that the envelope record utilizes the artifact_list field where a single artifact does not.

**Potential Errors**:

- The UUID does not exist.


#### Artifact Retrieve History

```
GET /ledger/api/v1.1/artifacts/history/{uuid}
```

Allows the user to obtain the historical data of an artifact associating with the uuid from the sPart ledger. This is a public function.

Example of historical artifact response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "ArtifactRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

**Potential Errors**:

- The UUID does not exist.


#### Artifact Retrieve Range

```
GET /ledger/api/v1.1/artifacts/{uuid}/date/{yyyymmdd}
```

Allows the user to obtain the artifact data associating with the uuid from the sPart ledger on the specified date. This is a public function and it returns the state which is most relevant to the given date.

Example of range artifact response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "ArtifactRecord",
    result: {
        ...
    }
}
```

**Potential Errors**:

- The UUID does not exist.


### Category Family
----

A category represents the category in which the part is assocaited. For instance, a part can be operating system, security sofware, and etc.


#### Category Create

```
POST /ledger/api/v1.1/categories
```

Allows the user to create a category into the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of category</td>
        </tr>
    </tbody>
</table>
<br>

Example of single category request:

```
{
    "private_key": "4761b2ef44d595c98022dd3a59da5cc135f7331193eaca31531eff9e1a122d73",
    "public_key": "02161c3bdef1135f21e0018f906d92d6fd790799d53e73a9474787ed9a99a30510",
    "category": {
    	"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a",
    	"name": "OS",
    	"description": "operating_system"
	}
}
```

Example of curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg","public_key": "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a","category": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a","name": "OS","description":"operating_system"}}' http://147.11.176.111:818/ledger/api/v1.1/categories
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID is not in a valid format.
- The UUID is not unique to the category.


#### Category Amend

```
POST /ledger/api/v1.1/category/amend
```

Allows the user to amend a category in the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of category</td>
        </tr>
    </tbody>
</table>
<br>

Example of single category request:

```
{
    "private_key": "4761b2ef44d595c98022dd3a59da5cc135f7331193eaca31531eff9e1a122d73",
    "public_key": "02161c3bdef1135f21e0018f906d92d6fd790799d53e73a9474787ed9a99a30510",
    "category": {
    	"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a",
    	"name": "OS",
    	"description": "operating_system"
	}
}

	or

{
    "private_key": "4761b2ef44d595c98022dd3a59da5cc135f7331193eaca31531eff9e1a122d73",
    "public_key": "02161c3bdef1135f21e0018f906d92d6fd790799d53e73a9474787ed9a99a30510",
    "category": {
    	"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a",
    	"name": "iOS"
	}
}	
```

(Note: The payload is valid as long as "private_key", "public_key", "category" and "uuid" are present.)

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg","public_key": "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a","category": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a","name": "OS","description":"operating_system"}}' http://147.11.176.111:818/ledger/api/v1.1/categories/amend

	or

curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg","public_key": "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a","category": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae703a","name": "iOS"}}' http://147.11.176.111:818/ledger/api/v1.1/categories/amend
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the create.
- No fields were amended.
- The UUID is not in a valid format. 
- The UUID does not exist.


#### Category List

```
GET /ledger/api/v1.1/categories
```

Allows the user to obtain a list of category from the sPart ledger. This is a public function.

Example of list category response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "CategoryRecord",
	result: [ 
		{	
			name: "operating system",
			description: "Operating System",
			uuid: "43903f02-00fd-43a3-bdaa-befe4a2fcd7e"
		},
		{	
			name: "libraries",
			description: "Operating System",
			uuid: "12d7f0c2-00fd-43a3-bdaa-befe4a2fcd7e"
		}
	]
}
```

If there are no categories registered, then an empty list is returned as shown:

```
{   
    status:     "success",
    message:    "OK",
    result_type: "CategoryRecord",
    result:     []
}
```

**Potential Errors**:

- Data cannot be deserialized due to encoding error.


#### Category Retrieve

```
GET /ledger/api/v1.1/categories/{uuid}
```

Allows the user to obtain the category data associating with the uuid from the sPart ledger. This is a public function.

Example of a <u>single</u> category response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "CategoryRecord",
	result: {	
		name: "OS",
		description: "Operating System",
		uuid: "43903f02-00fd-43a3-bdaa-befe4a2fcd7e"
	}
}
```

**Potential Errors**:

- The UUID does not exist.


#### Category Retrieve History

```
GET /ledger/api/v1.1/categories/history/{uuid}
```

Allows the user to obtain the historical data of a category associating with the uuid from the sPart ledger. This is a public function.

Example of historical category response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "CategoryRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

**Potential Errors**:

- The UUID does not exist.


#### Category Retrieve Range

```
GET /ledger/api/v1.1/categories/{uuid}/date/{yyyymmdd}
```

Allows the user to obtain the category data associating with the uuid from the sPart ledger on the specified date. This is a public function and it returns the state which is most relevant to the given date.

Example of range category response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "CategoryRecord",
    result: {
        ...
    }
}
```

**Potential Errors**:

- The UUID does not exist.


### Organization Family
----

A organization can represent a company, foundation, project or individual whom are associated with parts or artifacts. For instance, "Mac" is owned by "Apple Inc", which is an organization.


#### Organization Create

```
POST /ledger/api/v1.1/orgs
```

Allows the user to create an organization into the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>alias</td>
            <td>string</td>
            <td>alias of the organization</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>type</td>
            <td>string</td>
            <td>the type associated with organization (ie. software, hardware, ...)</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of category</td>
        </tr>
        <tr>
            <td>url</td>
            <td>string</td>
            <td>the url to website associated with the organization</td>
        </tr>
    </tbody>
</table>
<br>

Example of single organization request:

```
{
	"private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
	"public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
	"organization" : {
		"uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155e",
		"alias" : "oracle",
		"name" : "Oracle",
		"type" : "software",
		"description" : "java",
		"url" : "http://www.oracle"
	}
}
```

Example of curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg", "public_key" : "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a","supplier":{"name": "Oracle", "uuid": "3568f20a-8faa-430e-7c65-e9fce9aa155e", "alias": "oracle", "url": "http://www.oracle"}}' http://147.11.176.111:818/ledger/api/v1.1/orgs
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID is not in a valid format.
- The UUID is not unique to the category.


#### Organization Amend

```
POST /ledger/api/v1.1/organization/amend
```

Allows the user to amend an organization in the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>alias</td>
            <td>string</td>
            <td>alias of the organization</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>type</td>
            <td>string</td>
            <td>the type associated with organization (ie. software, hardware, ...)</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of category</td>
        </tr>
        <tr>
            <td>url</td>
            <td>string</td>
            <td>the url to website associated with the organization</td>
        </tr>
    </tbody>
</table>
<br>

Example of single organization request:

```
{
	"private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
	"public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
	"organization" : {
		"uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155e",
		"alias" : "oracle",
		"name" : "Oracle",
		"type" : "software",
		"description" : "java",
		"url" : "http://www.oracle"
	}
}

	or

{
	"private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
	"public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
	"organization" : {
		"uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155e",
		"alias" : "oracle",
		"name" : "Oracle",
		"url" : "http://www.oracle"
	}
}
```

(Note: The payload is valid as long as "private_key", "public_key", "organization" and "uuid" are present.)

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg", "public_key" : "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a", "organization":{"name": "Oracle", "uuid": "3568f20a-8faa-430e-7c65-e9fce9aa155e", "alias": "oracle", "type" : "software", "description" : "java","url": "http://www.oracle"}}' http://147.11.176.111:818/ledger/api/v1.1/orgs/amend

	or

curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg", "public_key" : "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a", "organization": {"name": "Oracle", "uuid": "3568f20a-8faa-430e-7c65-e9fce9aa155e", "alias": "oracle", "url": "http://www.oracle"}}' http://147.11.176.111:818/ledger/api/v1.1/orgs/amend
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the create.
- No fields were amended.
- The UUID is not in a valid format. 
- The UUID does not exist.


#### Organization List

```
GET /ledger/api/v1.1/orgs
```

Allows the user to obtain a list of organization from the sPart ledger. This is a public function.

Example of list organization response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "OrganizationRecord",
	result: [ 
		{	
			name: "Tesla, Inc.",
			uuid: "31e3e600-cd79-4ee5-464e-e74e1ce763cc",
			alias: "Tesla",
			type: "customer",
			description: "Company specializing in electric vehicles and lithium-ion battery energy storage."
			url: "http://www.tesla.com"
		},
		{	
			name: "General Motors Corporation",
			uuid: ""2584a6ce-16a7-44c0-7e53-21969d1e026b",
			alias: "GM",
			type: "customer"
			description: "United States automotive manufacturer."
			url: "http://www.gm.com"
		},
		{	
			name: "Wind River Systems",
			uuid: ""3568f20a-8faa-430e-7c65-e9fce9aa155d",
			alias: "WindRiver",
			type: "supplier"
			description: "United States automotive manufacturer."
			url: "http://www.windriver.com"
		}
	]
}
```

If there are no organizations registered, then an empty list will be returned as shown: 

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "OrganizationRecord",
	result: 	[]
}
```

**Potential Errors**:

- Data cannot be deserialized due to encoding error.


#### Organization Retrieve

```
GET /ledger/api/v1.1/orgs/{uuid}
```

Allows the user to obtain the organization data associating with the uuid from the sPart ledger. This is a public function.

Example of <u>single</u> organization response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "OrganizationRecord",
	result: {
		name: "...",
		uuid: "...",
		alias: "...",
		type: "..."
		description: "...",
		url: "..."
	}
}
```

**Potential Errors**:

- The UUID does not exist.


#### Organization Retrieve History

```
GET /ledger/api/v1.1/orgs/history/{uuid}
```

Allows the user to obtain the historical data of a part associating with the uuid from the sPart ledger. This is a public function.

Example of historical organization response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "OrganizationRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

**Potential Errors**:

- The UUID does not exist.


#### Organization Retrieve Range

```
GET /ledger/api/v1.1/orgs/{uuid}/date/{yyyymmdd}
```

Allows the user to obtain the organization data associating with the uuid from the sPart ledger on the specified date. This is a public function and it returns the state which is most relevant to the given date.

Example of range organization response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "OrganizationRecord",
    result: {
        ...
    }
}
```

**Potential Errors**:

- The UUID does not exist.


### Part Family
----

A part in sPart ledger primarily represents file. That being said, part can be anything and it is the smallest matter in the sPart ledger to how atoms are considered the building blocks of matters.


#### Part Create

```
POST /ledger/api/v1.1/parts
```

Allows the user to create a part into the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>label</td>
            <td>string</td>
            <td>label of the part</td>
        </tr>
        <tr>
            <td>alias</td>
            <td>string</td>
            <td>alias of the part</td>
        </tr>
        <tr>
            <td>version</td>
            <td>string</td>
            <td>version of the part</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>checksum of the part</td>
        </tr>
        <tr>
            <td>licensing</td>
            <td>string</td>
            <td>licensing of the part</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of the part</td>
        </tr>
    </tbody>
</table>
<br>

Example of single part request:

```
{
	"private_key": "74d39e1cd7c0ca98fbe860547569b7b54a23ab445f88a42cdc4ac2a041e14c77", 
	"public_key" : "03b7bc00503dc13596e1f3bc9216c20fbd4bb168a345d44a94bb8af0d42f1137f9",
	"part": {
		"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae701a",
		"name": "Zephypr 1.12 SPDX file",
		"version": "1.12",
		"checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
		"alias": "zephypr_1.10",
		"label": "Zephyr 1.10",
		"licensing": "MIT",
		"description": "Zephyr is a small real-time operating supporting multiple architectures"
	}
}
```

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "part": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae701a", "name": "Zephypr 1.12 SPDX file", "version": "1.12", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.10", "label": "Zephyr 1.10", "licensing": "MIT", "description": "Zephyr is a small real-time operating supporting multiple architectures" } }' http://localhost:3075/ledger/api/v1.1/parts
```
**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID is not in a valid format.
- The UUID is not unique to the part.


#### Part Amend

```
POST /ledger/api/v1.1/parts/amend
```

Allows the user to amend a part in the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>uuid</td>
            <td>string</td>
            <td>unique identifier</td>
        </tr>
        <tr>
            <td>name</td>
            <td>string</td>
            <td>category name</td>
        </tr>
        <tr>
            <td>label</td>
            <td>string</td>
            <td>label of the part</td>
        </tr>
        <tr>
            <td>alias</td>
            <td>string</td>
            <td>alias of the part</td>
        </tr>
        <tr>
            <td>version</td>
            <td>string</td>
            <td>version of the part</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>checksum of the part</td>
        </tr>
        <tr>
            <td>licensing</td>
            <td>string</td>
            <td>licensing of the part</td>
        </tr>
        <tr>
            <td>description</td>
            <td>string</td>
            <td>description of the part</td>
        </tr>
    </tbody>
</table>
<br>

Example of single part request:

```
{
	"private_key": "74d39e1cd7c0ca98fbe860547569b7b54a23ab445f88a42cdc4ac2a041e14c77", 
	"public_key" : "03b7bc00503dc13596e1f3bc9216c20fbd4bb168a345d44a94bb8af0d42f1137f9",
	"part": {
		"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae701a",
		"name": "Zephypr 1.12 SPDX file",
		"version": "1.12",
		"checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
		"alias": "zephypr_1.10",
		"label": "Zephyr 1.10",
		"licensing": "MIT",
		"description": "Zephyr is a small real-time operating supporting multiple architectures"
	}
}
```

(Note: The payload is valid as long as "private_key", "public_key", "part" and "uuid" are present.)

Example curl request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "part": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae701a", "name": "Zephypr 1.12 SPDX file", "version": "1.12", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "alias": "zephypr_1.10", "label": "Zephyr 1.10", "licensing": "MIT", "description": "Zephyr is a small real-time operating supporting multiple architectures" } }' http://localhost:3075/ledger/api/v1.1/parts

	or

curl -i -H "Content-Type: application/json" -X POST -d  '{"private_key": "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt", "public_key" : "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9", "part": {"uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae701a", "name": "Zephypr 1.12 SPDX file", "version": "1.12", "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1", "description": "Zephyr is a small real-time operating supporting multiple architectures" } }' http://localhost:3075/ledger/api/v1.1/parts
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the create.
- No fields were amended.
- The UUID is not in a valid format. 
- The UUID does not exist.


#### Part List

```
GET /ledger/api/v1.1/parts
```

Allows the user to obtain a list of part from the sPart ledger. This is a public function.

Example of list part response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "PartRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

If there are not parts registered, then an empty list is returned as shown:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "PartRecord",
	result: []
}
```

**Potential Errors**:

- Data cannot be deserialized due to encoding error.


#### Part Retrieve

```
GET /ledger/api/v1.1/parts/{uuid}
```

Allows the user to obtina the part data associating with the uuid from the sPart ledger. This is a public function.

Example of a <u>single</u> part response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "PartRecord",
	result: {	
		...
	}
}
```

**Potential Errors**:

- The UUID does not exist.


#### Part Retrieve History

```
GET /ledger/api/v1.1/parts/history/{uuid}
```

Allows the user to obtain the historical data of a part associating with the uuid from the sPart ledger. This is a public function.

Example of historical part response:

```
{
    status:     "success",
    message:    "OK",
    result_type: "PartRecord",
    result: [
        {
            ...
        },
        {
            ...
        },
        {
            ...
        }
    ]
}
```

**Potential Errors**:

- The UUID does not exist.


#### Part Retrieve Range

```
GET /ledger/api/v1.1/parts/{uuid}/date/{yyyymmdd}
```

Allows the user to obtain the part data associating with the uuid from the sPart ledger on the specified date. This is a public function and it returns the state which is most relevant to the given date.

Example of range part response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "PartRecord",
	result: {	
		...
	}
}
```

**Potential Errors**:

- The UUID does not exist.


### Relational Call
----

Relational calls allow the user to establish relations between different families. For instance, same part can be supplied by two different organizations.


#### Add Sub-Artifact to Artifact

```
POST /ledger/api/v1.1/artifacts/artifact
```

Allows the user to establish relationship between artifacts. In other words, creating a sub-folder in a folder. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the primary artifact (root)</td>
        </tr>
        <tr>
            <td>sub_artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the secondary artifact (child)</td>
        </tr>
        <tr>
            <td>path</td>
            <td>string</td>
            <td>the patch in which the secondary artifact resides</td>
        </tr>
    </tbody>
</table>
<br>

Example of add sub-artifact request:

```
{
   "private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg",
   "public_key": "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a",
   "relation": {
      "artifact_uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae704e",
      "sub_artifact_uuid": "/ledger/api/v1/artifacts",
      "path": "/notices"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of either artifact does not exists in the ledger.


#### Sever Sub-Artifact to Artifact

```
POST /ledger/api/v1.1/artifacts/artifact/delete
```

Allows the user to sever relationship between artifacts. In other words, removing a sub-folder in a folder. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the primary artifact (root)</td>
        </tr>
        <tr>
            <td>sub_artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the secondary artifact (child)</td>
        </tr>
        <tr>
            <td>path</td>
            <td>string</td>
            <td>the patch in which the secondary artifact resides</td>
        </tr>
    </tbody>
</table>
<br>

Example of sever sub-artifact request:

```
{
   "private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg",
   "public_key": "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a",
   "relation": {
      "artifact_uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae704e",
      "sub_artifact_uuid": "/ledger/api/v1/artifacts",
      "path": "/notices"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of either artifact does not exists in the ledger.
- Identical "path" is required to sever the relationship.


#### Add URI to Artifact

```
POST /ledger/api/v1.1/artifacts/uri
```

Allows the user to establish relationship between artifact and its uri. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>version</td>
            <td>string</td>
            <td>version of the uri</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>checksum of the uri</td>
        </tr>
        <tr>
            <td>size</td>
            <td>string</td>
            <td>size of the uri</td>
        </tr>
        <tr>
            <td>content_type</td>
            <td>string</td>
            <td>content type of the uri</td>
        </tr>
        <tr>
            <td>uri_type</td>
            <td>string</td>
            <td>uri type of the uri</td>
        </tr>
        <tr>
            <td>location</td>
            <td>string</td>
            <td>location of the uri</td>
        </tr>
    </tbody>
</table>
<br>

Example of add uri request:

```
{
   "private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg",
   "public_key" : "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a",
   "uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae704e",
   "uri": {
      "version": "1.0",
      "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
      "size": "235120",
      "content_type": ".pdf",
      "uri_type": "http",
      "location": "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016"
   }
}
```

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of artifact does not exists in the ledger.


#### Sever URI to Artifact

```
POST /ledger/api/v1.1/artifacts/uri/delete
```

Allows the user to sever relationship between artifact and its uri. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
        <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>version</td>
            <td>string</td>
            <td>version of the uri</td>
        </tr>
        <tr>
            <td>checksum</td>
            <td>string</td>
            <td>checksum of the uri</td>
        </tr>
        <tr>
            <td>size</td>
            <td>string</td>
            <td>size of the uri</td>
        </tr>
        <tr>
            <td>content_type</td>
            <td>string</td>
            <td>content type of the uri</td>
        </tr>
        <tr>
            <td>uri_type</td>
            <td>string</td>
            <td>uri type of the uri</td>
        </tr>
        <tr>
            <td>location</td>
            <td>string</td>
            <td>location of the uri</td>
        </tr>
    </tbody>
</table>
<br>

Example of sever uri request:

```
{
   "private_key": "5HvGd1pTeL6vECR1Whk6Hfk6rXuEtvug3g69GyL2LdnPiz8AJMg",
   "public_key" : "03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a",
   "uuid": "7709ca8d-01f4-4de2-69ed-16b7ebae704e",
   "uri": {
      "version": "1.0",
      "checksum": "f855d41c49e80b9d6f2a13148e5eb838607e92f1",
      "size": "235120",
      "content_type": ".pdf",
      "uri_type": "http",
      "location": "https://github.com/zephyrstorage/_content/master/f67d3213907a52012a4367d8abc016"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of artifact does not exists in the ledger.
- All the fields insdie "uri" must be identical to sever the relationship.


#### Add Part to Organization and Organization to Part

```
POST /ledger/api/v1.1/parts/orgs
```

Allows the user to establish relationship between part and organization. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>organization_uuid</td>
            <td>string</td>
            <td>unique identifier of the organization</td>
         </tr>
   </tbody>
</table>
<br>

Example of add organization to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "organization_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of either one does not exists in the ledger.


#### Sever Part to Organization and Organization to Part

```
POST /ledger/api/v1.1/parts/orgs/delete
```

Allows the user to sever relationship between part and organization. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>organization_uuid</td>
            <td>string</td>
            <td>unique identifier of the organization</td>
         </tr>
   </tbody>
</table>
<br>

Example of sever organization to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "organization_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of part one does not exists in the ledger.
- The UUID of organization does not exists in the part.


#### Add Artifact to Part

```
POST /ledger/api/v1.1/artifacts/part
```

Allows the user to establish relationship between part and artifact. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the artifact</td>
         </tr>
   </tbody>
</table>
<br>

Example of add artifact to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "artifact_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of either one does not exists in the ledger.


#### Sever Artifact to Part

```
POST /ledger/api/v1.1/artifacts/part/delete
```

Allows the user to sever relationship between part and artifact. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>artifact_uuid</td>
            <td>string</td>
            <td>unique identifier of the artifact</td>
         </tr>
   </tbody>
</table>
<br>

Example of sever artifact to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "artifact_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of part one does not exists in the ledger.
- The UUID of artifact does not exists in the part.


#### Add Category to Part

```
POST /ledger/api/v1.1/parts/categories
```

Allows the user to establish relationship between part and category. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>category_uuid</td>
            <td>string</td>
            <td>unique identifier of the category</td>
         </tr>
   </tbody>
</table>
<br>

Example of add category to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "category_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of either one does not exists in the ledger.


#### Sever Category to Part

```
POST /ledger/api/v1.1/parts/categories/delete
```

Allows the user to sever relationship between part and category. The request must be performed by a user with roles: "admin" or "member".


<table>
   <thead>
         <tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
         </tr>
   </thead>
   <tbody>
         <tr>
            <td>part_uuid</td>
            <td>string</td>
            <td>unique identifier of the part</td>
         </tr>
         <tr>
            <td>category_uuid</td>
            <td>string</td>
            <td>unique identifier of the category</td>
         </tr>
   </tbody>
</table>
<br>

Example of sever category to part request:

```
{
   "private_key" :  "338c2bb1b6985771a42144181af15d5d1d36a2ff047ce2217eadf784aade9d3f",
   "public_key" :  "02ba744c3f3ab596fc1af38bd0e3e2426703817257a8fda736187c8a6bb098c464",
   "relation" : {
      "part_uuid" : "7709ca8d-01f4-4de2-69ed-16b7ebae701a", 
      "category_uuid" : "3568f20a-8faa-430e-7c65-e9fce9aa155d"
   }
}
```

**Potential Errors**

- The requesting user does not have the appropriate access credentails to perform the create.
- One or more of the required fields are missing.
- The UUID of part one does not exists in the ledger.
- The UUID of category does not exists in the part.


### User Family
----

A user represents a client who is registered either by the "admin" or "member" into the sPart ledger. Typically, a user is someone who is associated to an organization but not need to be.


#### Get Key Pair

```
GET /ledger/api/v1/keys
```

Allows the user to obtain a new set of private and public keys which can be used later to register user into the sPart ledger. This is a public function.

Example Response:

```
{	
	status: 	"success",
	message: 	"OK",
	result_type: "PrivatePublicKeyRecord",
	result: {
		private_key: "5K6Q2kMHaMUrRjvSb4EPXQEJi1nAy3uXAYhqYvq2qNiLEGuFuVS",
		public_key: "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"
	}
}
```


#### Register User

```
POST /ledger/api/v1/registeruser
```

Allows the user to register other users into the sPart ledger. The request must be performed by a user with roles: "admin" or "member".


<table>
	<thead>
		<tr>
            <th>Field</th>
            <th>Type</th>
            <th>Description</th>
        </tr>
   	</thead>
   	<tbody>
        <tr>
            <td>user_name</td>
            <td>string</td>
            <td>user name</td>
        </tr>
        <tr>
            <td>email_address</td>
            <td>string</td>
            <td>user email</td>
        </tr>
        <tr>
            <td>role</td>
            <td>string</td>
            <td>the role (e.g., "member")</td>
        </tr>
        <tr>
            <td>authorized</td>
            <td>string</td>
            <td>the authorization (e.g., "allow")</td>
        </tr>
        <tr>
            <td>public_key</td>
            <td>string</td>
            <td>the user's public key</td>
        </tr>
    </tbody>
</table>
<br>

Example of single user request:

```
{
   private_key: "5K9ft3F4CDHMdGbeUZSyt77b1TJavfR7CAEgDZ7nXbdno1aynbt",
   public_key: "034408551a7b24b917103ccfafb402195713cd2e5dcdc588e7dc537f07b195bcf9",
   user: {	
		user_name: "7709ca8d-01f4-4de2-69ed-16b7ebae704a",
		email_address: "John.Doe@intel.com",
		role: "member",
		authorized: "allow",
   		public_key: 	
   			"03ef24753779355b4841dcef68a28044d1bc41b508b75bf8455b8518a5a61da50a"
	}
}
```

Example curl Request:

```
curl -i -H "Content-Type: application/json" -X POST -d  '{"name": "John Doe", "email": "john.doe@windriver.com", "role": "admin", "authorized": "allow", "public_key": "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"}' http://192.168.44.1:818/ledger/api/v1/registeruser
```

**Potential Errors**:

- The requesting user does not have the appropriate access credentials to perform the add.
- One or more of the required fields are missing.
- The UUID is not in a valid format.


#### User List (NO TYET IMPLEMENTED)

```
GET /ledger/api/v1/users
```

Allows the user to obtain a list of user from the sPart ledger. This is a public function

Example of list user response:

```
NOT YET IMPLEMENTED
```

If there are no users registered, then an empty list will be returned as shown:

```
NOT YET IMPLEMENTED
```

**Potential Errors**:

- Data cannot be deserialized due to encoding error.


#### User Retrieve (NOT YET IMPLEMENTED)

```
GET /ledger/api/v1/users/{public_key}
```

Allows the user to obtain the user data associating with the public key from the sPart ledger. This is a public function.

Example of a <u>single</u> user response:

```
{
	status: 	"success",
	message: 	"OK",
	result_type: "UserRecord",
	result: {	
		name: "Sameer Ahmed",
		email: "sameer.ahmed@windriver.com",
		organization: "Wind River"
		public_key: "0315d60b8dd9a90c55f2f7643270bc46d20798c1e5a38a30c9cb839882398d537f"
	}
}
```

**Potential Errors**:

- The public key does not exist.


## III) API Types

The explanation of the fields in the record associated with different families in the sPart ledger.

### Artifact (UD)
----

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


### Category (UD)
----


### Organization (UD)
----

#### OrganizationRecord

```
{
	UUID    string `json:"uuid"`               // universal unique identifier
	Name    string `json:"name"`               // Fullname
	Alias   string `json:"alias"`    // 1-15 alphanumeric characters
	Url     string `json:"url"`      // 2-3 sentence description
	Parts   ListOf:PartItemRecord
}
```


### Part (UD)
----

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

#### PartItemRecord 

```
{
	PartUUID string `json:"part_uuid"` // Part uuid
}
```


### User (UD)
----

#### UserRecord

```
{
	Name       string `json:"user_name"`
	Email      string `json:"email_address"`
	Role       string `json:"role"`
	Authorized string `json:"authorized"`
	PublicKey  string `json:"user_public_key"`
}
```

#### UserRegisterRecord

```
{
	User       UserRecord `json:"user"`
	PrivateKey string     `json:"private_key"`
	PublicKey  string     `json:"public_key"`
}
```

#### PrivatePublicKeyRecord

```
{
	PrivateKey string `json:"private_key"` // Private key
	PublicKey  string `json:"public_key"`  // Pubkic key
}
```


### Other (UD)
----

#### EmptyRecord

```
{ }
```

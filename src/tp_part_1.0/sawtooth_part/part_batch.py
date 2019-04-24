# Copyright 2016 Intel Corporation
# Copyright 2017 Wind River
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------------------
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
import hashlib
import base64
import time
import requests
import yaml
import datetime
import json
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader
from sawtooth_sdk.protobuf.transaction_pb2 import Transaction
from sawtooth_sdk.protobuf.batch_pb2 import BatchList
from sawtooth_sdk.protobuf.batch_pb2 import BatchHeader
from sawtooth_sdk.protobuf.batch_pb2 import Batch
from sawtooth_part.exceptions import PartException
################################################################################
#                            GLOBAL FUNCTIONS                                  #
################################################################################
def _sha512(data):
    """
    Creates the string of sha512 hashed to the passed in data.
    
    Args:
        data (bytes): The data to be hashed
    
    Returns:
        type: str
        The sha512 hashed data in string of hex values.
        
    """
    return hashlib.sha512(data).hexdigest()
################################################################################
#                                  CLASS                                       #
################################################################################
class PartBatch:
    """
    Class for creating batch of the Transaction Family : Part
    
    Attributes:
        base_url (str): The base url of the transaction family
    
    """
    
    def __init__(self, base_url):
        """
        Constructs the PartBatch object.
        
        Args:
            base_url (str): The base url of the transaction family
        
        """
        self._base_url = base_url
################################################################################
#                            PUBLIC FUNCTIONS                                  #
################################################################################    
    def create_part(self, pt_id, pt_name, checksum, version, alias, licensing,
                label, description, private_key, public_key):
        """
        Constructs the batch payload for the "create" command.
        
        Args:
            pt_id (str): The uuid of the part
            pt_name (str): The name of the part
            checksum (str): The checksum of the part
            version (str): The version of the part
            alias (str): The alias of the part
            licensing (str): The licensing of the part
            label (str): The label of the part
            description (str): The description of the part
            private_key (str): The private key of the user
            public_key (str): The public key of the user
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the 
            "create" command.
                
                or
            
            type: None
            None object if the UUID already exists in the ledger.
            
        """
        # Checking if the uuid is unique
        address = self._get_address(pt_id)
        response_bytes = self._send_request(
                        "state/{}".format(address), pt_id=pt_id, creation=True
                    )
        if response_bytes != None:
            return None
        
        # Creating the batch object to be returned
        cur = self._get_block_num()
        return self.create_part_transaction(pt_id, pt_name, checksum, version, 
                    alias, licensing, label, description, "create", private_key,
                    public_key, [], [], [], "0", cur, 
                    str(datetime.datetime.utcnow()))
  
    def amend_part(self, pt_id, pt_name, checksum, version, alias, licensing,
                label, description, private_key, public_key):
        """
        Constructs the batch payload for the "amend" command.
        
        Args:
            pt_id (str): The uuid of the part
            pt_name (str): The name of the part
            checksum (str): The checksum of the part
            version (str): The version of the part
            alias (str): The alias of the part
            licensing (str): The licensing of the part
            label (str): The label of the part
            description (str): The description of the part
            private_key (str): The private key of the user
            public_key (str): The public key of the user
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the
            "amend" command.
            
                or
            
            type: None
            None object if the UUID does not exist in the ledger.
            
                or 
            
            type: list pertaining None
            List containing only None object if no member was amended.
            
        """
        # Checking if the uuid exists
        response_bytes = self.retrieve_part(pt_id)
        if response_bytes != None:
            
            # Loading the data to perform checks
            jresponse = json.loads(response_bytes.decode())
            
            # Checking if params are "null"; if yes, replace with prior values
            if pt_name == "null":
                pt_name = jresponse["name"]
            if checksum == "null":
                checksum = jresponse["checksum"]
            if version == "null":
                version = jresponse["version"]
            if alias == "null":
                alias = jresponse["alias"]
            if licensing == "null":
                licensing = jresponse["licensing"]
            if label == "null":
                label = jresponse["label"]
            if description == "null":
                description = jresponse["description"]
            
            # Checking if any of the params were changed; if not, return [None]
            if (jresponse["name"]        == pt_name      and
                jresponse["checksum"]    == checksum     and
                jresponse["version"]     == version      and
                jresponse["alias"]       == alias        and
                jresponse["licensing"]   == licensing    and
                jresponse["label"]       == label        and
                jresponse["description"]    == description):
                return [None]
            else:
                # Creating the batch object to be returned
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, pt_name, checksum, 
                            version, alias, licensing, label, description, 
                            "amend", private_key, public_key, 
                            jresponse["artifact_list"], 
                            jresponse["category_list"],
                            jresponse["organization_list"], 
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
        
        return None
    
    def add_artifact(self, pt_id, artifact_id, private_key, public_key, 
                        del_flag=False):
        """
        Constructs the batch payload for the "AddArtifact" command.
        
        Args:
            pt_id (str): The uuid of the part
            artifact_id (str): The uuid of the artifact
            private_key (str): The private key of the user
            public_key (str): The public key of the user
            del_flag (bool): The flag for "--delete" option (default False)
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the
            "AddArtifact" command.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
                or
            
            type: list pertaining None and str
            List containing None object and error message:
                * If "--delete"
                    > If "artifact_list" is empty
                    > If "artifact_id" is not in "artifact_list"
                * If "artifact_id" is in "artifact_list"
        
        """
        # Checking if "--delete" is invoked
        if del_flag:
            
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if len(jresponse["artifact_list"]) == 0:
                    return  [
                                None,
                                "No {} to remove from this {}." \
                                    .format("Artifact", "Part")
                            ]
                        
                if artifact_id not in jresponse["artifact_list"]:
                    return  [
                                None,
                                "No such {} in this {}." \
                                    .format("Artifact", "Part")
                            ]
                
                # Removing the "artifact_id" from "artifact_list" and creating
                # the batch object to be returned 
                jresponse["artifact_list"].remove(artifact_id)
                
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddArtifact", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None
        else:
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Checking if artifact to be added exists
                if self._validate_artifact_id(artifact_id) == None:
                    return  [
                                None,
                                "ArtifactException : UUID does not exist."
                            ]
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if artifact_id not in jresponse["artifact_list"]:
                    jresponse["artifact_list"].append(artifact_id)
                else:
                    return  [
                                None,
                                "Duplicate Artifact UUID in the Part."
                            ]
                
                # Creating a batch object to be returned along with 
                # updated "artifact_list"
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddArtifact", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None
        
    def add_category(self, pt_id, category_id, private_key, public_key,
                        del_flag=False):
        """
        Constructs the batch payload for the "AddCategory" command.
        
        Args:
            pt_id (str): The uuid of the part
            category_id (str): The uuid of the category
            private_key (str): The private key of the user
            public_key (str): The public key of the user
            del_flag (bool): The flag for "--delete" option (default False)
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the
            "AddCategory" command.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
                or
            
            type: list pertaining None and str
            List containing None object and error message:
                * If "--delete"
                    > If "category_list" is empty
                    > If "category_id" is not in "category_list"
                * If "category_id" is in "category_list"
        
        """
        # Checking if "--delete" is invoked
        if del_flag:
            
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if len(jresponse["category_list"]) == 0:
                    return  [
                                None,
                                "No {} to remove from this {}." \
                                    .format("Category", "Part")
                            ]
                
                if category_id not in jresponse["category_list"]:
                    return  [
                                None,
                                "No such {} in this {}." \
                                    .format("Category", "Part")
                            ]
                
                # Removing the "category_id" from "category_list" and creating
                # the batch object to be returned        
                jresponse["category_list"].remove(category_id)
                
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddCategory", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None
        else:
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Checking if category to be added exists
                if self._validate_category_id(category_id) == None:
                    return  [
                                None,
                                "CategoryException : UUID does not exist."
                            ]
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if category_id not in jresponse["category_list"]:
                    jresponse["category_list"].append(category_id)
                else:
                    return  [
                                None,
                                "Duplicate Category UUID in the Part."
                            ]
                
                # Creating a batch object to be returned along with 
                # updated "category_list"
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddCategory", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None
   
    def add_organization(self, pt_id, organization_id, private_key, public_key,
                        del_flag=False):
        """
        Constructs the batch payload for the "AddOrganization" command.
        
        Args:
            pt_id (str): The uuid of the part
            organization_id (str): the uuid of the organization
            private_key (str): The private key of the user
            public_key (str): The public key of the user
            del_flag (bool): The flag for "--delete" option (default False)
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the
            "AddOrganization" command.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
                or
            
            type: list pertaining None and str
            List containing None object and error message:
                * If "--delete"
                    > If "organization_list" is empty
                    > If "organization_id" is not in "organization_list"
                * If "organization_id" is in "organization_list"
            
        """
        # Checking if "--delete" is invoked
        if del_flag:
            
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if len(jresponse["organization_list"]) == 0:
                    return  [
                                None, 
                                "No {} to remove from this {}." \
                                    .format("Organization", "Part")
                            ]
                
                if organization_id not in jresponse["organization_list"]:
                    return  [
                                None,
                                "No such {} in this {}." \
                                    .format("Organization", "Part")
                            ]   
                
                # Removing the "organization_id" from "organization_list"
                # and creating the batch object to be returned
                jresponse["organization_list"].remove(organization_id)
                
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddOrganization", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None
        else:
            # Loading the data to perform checks
            response_bytes = self.retrieve_part(pt_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Checking if organization to be added exists
                if self._validate_organization_id(organization_id) == None:
                    return  [
                                None,
                                "OrganizationException : UUID does not exist."
                            ]
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                if organization_id not in jresponse["organization_list"]:
                    jresponse["organization_list"].append(organization_id)
                else:
                    return  [
                                None,
                                "Duplicate Organization UUID in the Part."
                            ]
                
                # Creating a batch object to be returned along with 
                # updated "organization_list"
                cur = self._get_block_num()
                return self.create_part_transaction(pt_id, jresponse["name"], 
                            jresponse["checksum"], jresponse["version"], 
                            jresponse["alias"], jresponse["licensing"], 
                            jresponse["label"], jresponse["description"], 
                            "AddOrganization", private_key, public_key, 
                            jresponse["artifact_list"],
                            jresponse["category_list"], 
                            jresponse["organization_list"],
                            jresponse["cur_block"], cur, 
                            str(datetime.datetime.utcnow()))
            
            return None

    def list_part(self):
        """
        Fetches the data from ledger and constructs the list of part.
        
        Returns:
            type: list of dict
            List of JSON (Python dict) associated with 
            the Transaction Family : Part.
            
                or
            
            type: None
            None object if deserialization of the data failed.
        
        """
        part_prefix = self._get_prefix()

        result = self._send_request(
            "state?address={}".format(part_prefix)
        )

        try:
            encoded_entries = yaml.safe_load(result)["data"]

            return [
                json.loads(base64.b64decode(entry["data"]).decode()) for entry \
                    in encoded_entries
            ]

        except BaseException:
            return None
    
    def retrieve_part(self, pt_id, all_flag=False, range_flag=None):
        """
        Fetches the data associating with UUID from ledger.
        
        Args:
            pt_id (str): The uuid of the part
            all_flag (bool): The flag for "--all" option (default False)
            range_flag (list of int):
                The flag for "--range" option (default None)
            
        Returns:
            type: bytes
            Bytes containing the data associated to the UUID.
            
                or
            
            type: list of dict
            List of JSON (Python dict) associated with the UUID.
                * If "--all" or "--range" are not default values
            
                or
            
            type: None
            None object if decoding failed.
        
        """
        # Checking if "--all" is invoked
        if all_flag:
            
            # Loading and instatiating to perform checks
            retVal = []
        
            response = self.retrieve_part(pt_id).decode()
            response = json.loads(response)
            
            # Checking if "--range" is invoked and performing checks to append
            if range_flag != None:
                curTime = int(response["timestamp"].split()[0].replace("-", ""))
                if (curTime <= int(range_flag[1]) and 
                        curTime >= int(range_flag[0])):
                    retVal.append(response)
            else:
                retVal.append(response)
            
            # While not "create" perform checks to append to list
            while str(response["prev_block"]) != "0":
                
                response = json.loads(self._get_payload_(
                                int(response["prev_block"])).decode())
                
                timestamp       = response["timestamp"] 
                
                del response["action"]
                
                if range_flag != None:
                    curTime = int(timestamp.split()[0].replace("-", ""))
                    if curTime < int(range_flag[0]):
                        break
                    elif curTime <= int(range_flag[1]):
                        retVal.append(response)
                else:
                    retVal.append(response)
            
            # Returning the list of JSON
            return retVal
        else:
            address = self._get_address(pt_id)
    
            result = self._send_request(
                            "state/{}".format(address), pt_id=pt_id
                        )
            try:
                return base64.b64decode(yaml.safe_load(result)["data"])
    
            except BaseException:
                return None
################################################################################
#                            PRIVATE FUNCTIONS                                 #
################################################################################
    def _get_prefix(self):
        """
        Constructs and returns a string of SHA512 hashed data of
        the Transaction Family : Part.
        
        Returns:
            type: str
            The first 6 characters of the SHA512 hashed data of
            the Transaction Family : Part.
            
        """
        return _sha512("pt".encode("utf-8"))[0:6]
    
    def _get_address(self, pt_id):
        """
        Constructs and returns a string of unique hashed data for the
        passed in UUID.
        
        Args:
            pt_id (str): The uuid of the part
        
        Returns:
            type: str
            The address-to-be, which associates the uuid and the prefix.
            
        """
        part_prefix = self._get_prefix()
        address = _sha512(pt_id.encode("utf-8"))[0:64]
        return part_prefix + address
    
    def _get_block_num(self):
        """
        Fetches the current block ID of the ledger.
        
        Returns:
            type: str
            The current block ID of the ledger as a string.
            
                or
            
            type: None
            None object if there is no block in the ledger.
        
        """
        part_prefix = self._get_prefix()
        result = self._send_request(
            "blocks?={}".format(part_prefix)
        )
        
        if result != None or result != "":
            result = json.loads(result)
            return str(len(result["data"]))
        return None
    
    def _get_payload_(self, blocknum):
        """
        Fetches the payload associated with the given block ID in the ledger.
        
        Args:
            blocknum (int): The block ID of the previous state of the UUID
        
        Returns:
            type: bytes
            The payload on the given block ID in bytes.
            
            type: None
            None object if there is no block in the ledger.
            
        """
        part_prefix = self._get_prefix()
        result = self._send_request(
            "blocks?={}".format(part_prefix)
        )
        
        if result != None or result != "":
            result = json.loads(result)
            payload = result["data"][-(blocknum + 1)]["batches"][0]\
                        ["transactions"][0]["payload"]
            
            return base64.b64decode(payload)
        return None
    
    def _validate_artifact_id(self, artifact_id):
        """
        Validates if the artifact UUID exists in the ledger.
        
        Args:
            artifact_id (str): The uuid of the artifact
        
        Returns:
            type: str
            Data as string if UUID exist in the ledger.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
        """
        artifact_prefix = _sha512("artifact".encode("utf-8"))[0:6]
        address = _sha512(artifact_id.encode("utf-8"))[0:64]
        address = artifact_prefix + address
        return self._send_request("state/{}".format(address))
    
    def _validate_category_id(self, category_id):
        """
        Validates if the category UUID exists in the ledger.
        
        Args:
            category_id (str): The uuid of the category
        
        Returns:
            type: str
            Data as string if UUID exist in the ledger.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
        
        """
        category_prefix = _sha512("category".encode("utf-8"))[0:6]
        address = _sha512(category_id.encode("utf-8"))[0:64]
        address = category_prefix + address
        return self._send_request("state/{}".format(address))
    
    def _validate_organization_id(self, organization_id):
        """
        Validates if the organization UUID exists in the ledger.
        
        Args:
            organization_id (str): The uuid of the organization
        
        Returns:
            type: str
            Data as string if UUID exist in the ledger.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
        """
        category_prefix = _sha512("organization".encode("utf-8"))[0:6]
        address = _sha512(organization_id.encode("utf-8"))[0:64]
        address = category_prefix + address
        return self._send_request("state/{}".format(address))
    
    def _send_request(self, suffix, data=None, content_type=None,
                        pt_id=None, creation=False):
        """
        Performs RESTful API call on the given params.
        
        Args:
            suffix (str): The suffix of the url in query
            data (str): The data to be sent in POST request (default None)
            content_type (str): The data type (default None)
            pt_id (str): The uuid of the part (default None)
            creation (bool): The flag for "create" command (default False)
        
        Returns:
            type: str
            Data associated with suffix as a string.
            
                or
                
            type: None
            None object if any exception occurs or "404" was raised during
            "create" command.
            
        Raises:
            PartException:
                * If "404" was raised for the request
                * If status was "sucessful"
            
        """
        # Building the URL
        if self._base_url.startswith("http://"):
            url = "{}/{}".format(self._base_url, suffix)
        else:
            url = "http://{}/{}".format(self._base_url, suffix)

        headers = {}
        if content_type is not None:
            headers["Content-Type"] = content_type
        
        # Performing appropriate RESTful API
        try:
            if data is not None:
                result = requests.post(url, headers=headers, data=data)
            else:
                result = requests.get(url, headers=headers)

            if result.status_code == 404:
                if creation:
                    return None
                raise PartException("No part found: {}".format(pt_id))

            elif not result.ok:
                raise PartException("Error {}: {}".format(
                    result.status_code, result.reason))

        except BaseException as err:
            print(err)
            return None
        
        # Returning the data as string
        return result.text
   
    def create_part_transaction(self, pt_id, pt_name, checksum, version, alias, 
                        licensing, label, description, action, private_key, 
                        public_key, artifact_id, category_id, organization_id,
                        prev, cur, timestamp):
        """
        Constructs the Batch to be posted and sent the request to be posted on
        the ledger.
        
        Args:
            pt_id (str): The uuid of the part
            pt_name (str): The name of the part
            checksum (str): The checksum of the part
            version (str): The version of the part
            alias (str): The alias of the part
            licensing (str): The licensing of the part
            label (str): The label of the part
            description (str): The description of the part
            private_key (str): The private key of the user
            public_key (str): The public key of the user
            artifact_id (list of str):
                The list of the artifact uuid associated with the part
            category_id (list of str):
                The list of the category uuid associated with the part
            organization_id (list str):
                The list of the organization uuid associated with the part
            prev (str): The previous block id of the transaction (default "0")
            cur (str): the current block id of the transaction
            timestamp (str): The UTC time for when the transaction was submitted
            action (str): The action performed
            
        Returns:
            type: str
            Data associated with suffix as a string.
            
                or
            
            type: None
            None object if _send_request failed.
            
        """
        # Constructing Batch to be sent and stored
        self._public_key = public_key
        self._private_key = private_key

        payload = {
            "uuid"              : str(pt_id),
            "name"              : str(pt_name),
            "checksum"          : str(checksum),
            "version"           : str(version),
            "alias"             : str(alias),
            "licensing"         : str(licensing),
            "label"             : str(label),
            "description"       : str(description),
            "action"            : str(action),
            "prev_block"        : str(prev),
            "cur_block"         : str(cur),
            "timestamp"         : str(timestamp),
            "artifact_list"     : artifact_id,
            "category_list"     : category_id,
            "organization_list" : organization_id
        }
        payload = json.dumps(payload).encode()
        address = self._get_address(pt_id)

        header = TransactionHeader(
            signer_public_key = self._public_key,
            family_name = "pt",
            family_version = "1.0",
            inputs = [address],
            outputs = [address],
            dependencies = [],
            payload_sha512 = _sha512(payload),
            batcher_public_key = self._public_key,
            nonce = time.time().hex().encode()
        ).SerializeToString()
        
        signature = CryptoFactory(create_context("secp256k1")) \
            .new_signer(Secp256k1PrivateKey.from_hex(self._private_key)) \
            .sign(header)

        transaction = Transaction(
            header = header,
            payload = payload,
            header_signature = signature
        )
        
        # Creating batch list
        batch_list = self._create_batch_list([transaction])
        
        return self._send_request(
            "batches", batch_list.SerializeToString(),
            "application/octet-stream"
        )

    def _create_batch_list(self, transactions):
        """
        Helps create a batch list to be transmitted to the ledger.
        
        Args:
            transactions (list of Transaction): List containing transaction IDs
        
        Returns:
            type: BatchList
            BatchList object where each batch in the list are constructed in
            the function. 
            
        """
        transaction_signatures = [t.header_signature for t in transactions]

        header = BatchHeader(
            signer_public_key=self._public_key,
            transaction_ids=transaction_signatures
        ).SerializeToString()

        signature = CryptoFactory(create_context("secp256k1")) \
            .new_signer(Secp256k1PrivateKey.from_hex(self._private_key)) \
            .sign(header)

        batch = Batch(
            header=header,
            transactions=transactions,
            header_signature=signature
        )
        return BatchList(batches=[batch])
################################################################################
#                                                                              #
################################################################################

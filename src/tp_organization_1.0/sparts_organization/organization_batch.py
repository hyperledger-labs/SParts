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
from sparts_organization.exceptions import OrganizationException
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
class OrganizationBatch:
    """
    Class for creating batch of the Transaction Family : Organization
    
    Attributes:
        base_url (str): The base url of the transaction family
    
    """
    
    def __init__(self, base_url):
        """
        Constructs the OrganizationBatch object.
        
        Args:
            base_url (str): The base url of the transaction family
        
        """
        self._base_url = base_url
################################################################################
#                            PUBLIC FUNCTIONS                                  #
################################################################################
    def create_organization(self, org_id, org_alias, org_name, org_type,
                description, org_url, private_key, public_key):
        """
        Constructs the batch payload for the "create" command.
        
        Args:
            org_id (str): The uuid of the organization
            org_alias (str): The alias of the organization
            org_name (str): The name of the organization
            org_type (str): The type of the organization
            description (str): The description of the organization
            org_url (str): The url of the organization
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
        address = self._get_address(org_id)
        response_bytes = self._send_request(
                    "state/{}".format(address), org_id=org_id, creation=True
                )
        if response_bytes != None:
            return None
        
        # Creating the batch object to be returned
        cur = self._get_block_num()
        return self.create_organization_transaction(org_id, org_alias, org_name, 
                    org_type, description, org_url, "create", private_key, 
                    public_key, "0", cur, str(datetime.datetime.utcnow()), "")
    
    def amend_organization(self, org_id, org_alias, org_name, org_type,
                description, org_url, private_key, public_key):
        """
        Constructs the batch payload for the "amend" command.
        
        Args:
            org_id (str): The uuid of the organization
            org_alias (str): The alias of the organization
            org_name (str): The name of the organization
            org_type (str): The type of the organization
            description (str): The description of the organization
            org_url (str): The url of the organization
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
        response_bytes = self.retrieve_organization(org_id)
        if response_bytes != None:
            
            # Loading the data to perform checks
            jresponse = json.loads(response_bytes.decode())
            
            # Checking if params are "null"; if yes, replace with prior values
            if org_alias == "null":
                org_alias = jresponse["alias"]
            if org_name == "null":
                org_name = jresponse["name"]
            if  org_type == "null":
                org_type = jresponse["type"]
            if description == "null":
                description = jresponse["description"]
            if org_url == "null":
                org_url = jresponse["url"]
            
            # Checking if any of the params were changed; if not, return [None]
            if (jresponse["alias"] == org_alias and
                jresponse["name"] == org_name and
                jresponse["type"] == org_type and
                jresponse["description"] == description and
                jresponse["url"] == org_url):
                return [None]
            else:
                # Creating the batch object to be returned
                cur = self._get_block_num()
                return self.create_organization_transaction(org_id, org_alias, 
                            org_name, org_type, description, org_url, "amend", 
                            private_key, public_key, jresponse["cur_block"], 
                            cur, str(datetime.datetime.utcnow()), 
                            jresponse["pt_list"])
                            
        return None

    def list_organization(self):
        """
        Fetches the data from ledger and constructs the list of organization.
        
        Returns:
            type: list of dict
            List of JSON (Python dict) associated with 
            the Transaction Family : Organization.
            
                or
            
            type: None
            None object if deserialization of the data failed.
        
        """
        organization_prefix = self._get_prefix()

        result = self._send_request(
            "state?address={}".format(organization_prefix)
        )

        try:
            encoded_entries = yaml.safe_load(result)["data"]

            return [
                json.loads(base64.b64decode(entry["data"]).decode()) for entry \
                    in encoded_entries
            ]

        except BaseException:
            return None

    def retrieve_organization(self, org_id, all_flag=False, range_flag=None):
        """
        Fetches the data associating with UUID from ledger.
        
        Args:
            org_id (str): The uuid of the organization
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
            
            response = self.retrieve_organization(org_id).decode()
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
            address = self._get_address(org_id)
    
            result = self._send_request("state/{}".format(address), 
                        org_id=org_id)
            
            try:
                return base64.b64decode(yaml.safe_load(result)["data"])
    
            except BaseException:
                return None
    
    def add_part(self, org_id, pt_id, private_key, public_key, del_flag=False):
        """
        Constructs the batch payload for the "AddPart" command.
        
        Args:
            org_id (str): The uuid of the organization
            pt_id (str): the uuid of the part
            private_key (str): The private key of the user
            public_key (str): The public key of the user
            del_flag (bool): The flag for "--delete" option (default False)
        
        Returns:
            type: Batch
            The batch object which pertains all the data associating with the
            "AddPart" command.
            
                or
            
            type: None
            None object if UUID does not exist in the ledger.
            
                or
            
            type: list pertaining None and str
            List containing None object and error message:
                * If "--delete"
                    > If "pt_list" is empty
                    > If "pt_id" is not in "pt_list"
                * If "pt_id" is in "pt_list"
            
        """
        # Checking if "--delete" is invoked
        if del_flag:
            
            # Loading the data to perform checks
            response_bytes = self.retrieve_organization(org_id)
            
            # Checking if uuid exists
            if response_bytes != None:
                
                # Loading the state of uuid
                jresponse = json.loads(response_bytes.decode())
                
                # Removing the "pt_id" from "pt_list" and creating the batch
                # object to be returned
                jresponse["pt_list"].remove(pt_id)
                
                cur = self._get_block_num()
                return self.create_organization_transaction(org_id, 
                            jresponse["alias"], 
                            jresponse["name"], 
                            jresponse["type"], 
                            jresponse["description"],
                            jresponse["url"], 
                            "AddPart", private_key, public_key, 
                            jresponse["cur_block"], cur,
                            str(datetime.datetime.utcnow()),
                            jresponse["pt_list"])
                            
            return None
        else:
            # Loading the data to perform checks
            response_bytes = self.retrieve_organization(org_id)
             
            # Checking if uuid exists
            if response_bytes != None:
                
                # Loading the state of uuid to perform checks
                jresponse = json.loads(response_bytes.decode())
                
                # Creating a batch object to be returned along with 
                # updated "pt_list"
                jresponse["pt_list"].append(pt_id)
                
                cur = self._get_block_num()
                return self.create_organization_transaction(org_id, 
                            jresponse["alias"], 
                            jresponse["name"], 
                            jresponse["type"], 
                            jresponse["description"],
                            jresponse["url"], 
                            "AddPart", private_key, public_key, 
                            jresponse["cur_block"], cur,
                            str(datetime.datetime.utcnow()),
                            jresponse["pt_list"])
                            
            return None
################################################################################
#                            PRIVATE FUNCTIONS                                 #
################################################################################
    def _get_prefix(self):
        """
        Constructs and returns a string of SHA512 hashed data of
        the Transaction Family : Organization.
        
        Returns:
            type: str
            The first 6 characters of the SHA512 hashed data of
            the Transaction Family : Organization.
            
        """
        return _sha512("organization".encode("utf-8"))[0:6]

    def _get_address(self, org_id):
        """
        Constructs and returns a string of unique hashed data for the
        passed in UUID.
        
        Args:
            org_id (str): The uuid of the organization
        
        Returns:
            type: str
            The address-to-be, which associates the uuid and the prefix.
            
        """
        organization_prefix = self._get_prefix()
        address = _sha512(org_id.encode("utf-8"))[0:64]
        return organization_prefix + address
    
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
        organization_prefix = self._get_prefix()

        result = self._send_request(
            "blocks?={}".format(organization_prefix)
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
        organization_prefix = self._get_prefix()

        result = self._send_request(
            "blocks?={}".format(organization_prefix)
        )
        
        if result != None or result != "":
            result = json.loads(result)
            payload = result["data"][-(blocknum + 1)]["batches"][0]\
                        ["transactions"][0]["payload"]
            
            return base64.b64decode(payload)
        return None
    
    def _send_request(self, suffix, data=None, content_type=None,
                        org_id=None, creation=False):
        """
        Performs RESTful API call on the given params.
        
        Args:
            suffix (str): The suffix of the url in query
            data (str): The data to be sent in POST request (default None)
            content_type (str): The data type (default None)
            org_id (str): The uuid of the organization (default None)
            creation (bool): The flag for "create" command (default False)
        
        Returns:
            type: str
            Data associated with suffix as a string.
            
                or
                
            type: None
            None object if any exception occurs or "404" was raised during
            "create" command.
            
        Raises:
            OrganizationException:
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
                raise OrganizationException("No such organization: {}".format(id))

            elif not result.ok:
                raise OrganizationException("Error {}: {}".format(
                    result.status_code, result.reason))

        except BaseException as err:
            print(err)
            return None
        
        # Returning the data as string
        return result.text

    def create_organization_transaction(self, org_id, org_alias, org_name, 
                org_type, description, org_url, action, private_key, public_key, 
                prev, cur, timestamp, pt_id):
        """
        Constructs the Batch to be posted and sent the request to be posted on
        the ledger.
        
        Args:
            org_id (str): The uuid of the organization
            org_alias (str): The name of the organization
            org_name (str): The checksum of the organization
            org_type (str): The version of the organization
            description (str): The description of the organization
            org_url (str): The url of the organization
            pt_id (str):
                The uuid of the part to be added to "pt_list" in organization
            private_key (str): The private key of the user
            public_key (str): The public key of the user
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
        
        payload  = {
            "uuid"          : str(org_id),
            "alias"         : str(org_alias),
            "name"          : str(org_name),
            "type"          : str(org_type),
            "description"   : str(description),
            "url"           : str(org_url),
            "action"        : str(action),
            "prev_block"    : str(prev),
            "cur_block"     : str(cur),
            "timestamp"     : str(timestamp),
            "pt_list"       : pt_id
        }
        payload = json.dumps(payload).encode()
        address = self._get_address(org_id)

        header = TransactionHeader(
            signer_public_key = self._public_key,
            family_name = "organization",
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
            signer_public_key = self._public_key,
            transaction_ids = transaction_signatures
        ).SerializeToString()

        signature = CryptoFactory(create_context("secp256k1")) \
            .new_signer(Secp256k1PrivateKey.from_hex(self._private_key)) \
            .sign(header)

        batch = Batch(
            header = header,
            transactions = transactions,
            header_signature = signature
        )
        return BatchList(batches=[batch])
################################################################################
#                                                                              #
################################################################################

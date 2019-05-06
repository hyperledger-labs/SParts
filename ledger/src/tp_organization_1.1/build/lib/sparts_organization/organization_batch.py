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
from base64 import b64encode
import time
import requests
import yaml
import datetime
import json
# import sawtooth_signing.secp256k1_signer as signing
#
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
#

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
    return hashlib.sha512(data).hexdigest()

################################################################################
#                                  CLASS                                       #
################################################################################
class OrganizationBatch:
   
    def __init__(self, base_url):
        self._base_url = base_url
################################################################################
#                            PUBLIC FUNCTIONS                                  #
################################################################################
    def create_organization(self, org_id, org_alias, org_name, org_type,
                description, org_url, private_key, public_key):
        address = self._get_address(org_id)
        response_bytes = self._send_request("state/{}".format(address), 
                    org_id=org_id, creation=True)            
        
        if response_bytes != None:
            return None
        
        cur = self._get_block_num()
        return self.create_organization_transaction(org_id, org_alias, org_name, 
                    org_type, description, org_url, "create", private_key, 
                    public_key, "0", cur, str(datetime.datetime.utcnow()), "")
    
    def amend_organization(self, org_id, org_alias, org_name, org_type,
                description, org_url, private_key, public_key):
        response_bytes = self.retrieve_organization(org_id)
        
        if response_bytes != None:
            
            jresponse = json.loads(response_bytes.decode())
            
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
            
            if (jresponse["alias"] == org_alias and
                jresponse["name"] == org_name and
                jresponse["type"] == org_type and
                jresponse["description"] == description and
                jresponse["url"] == org_url):
                return [None]
            else:
                cur = self._get_block_num()
                return self.create_organization_transaction(org_id, org_alias, 
                            org_name, org_type, description, org_url, "amend", 
                            private_key, public_key, jresponse["cur_block"], 
                            cur, str(datetime.datetime.utcnow()), 
                            jresponse["pt_list"])
                            
        return None

    def list_organization(self):
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
        if all_flag:
            
            retVal = []
            
            response = self.retrieve_organization(org_id).decode()
            response = json.loads(response)
            
            if range_flag != None:
                curTime = int(response["timestamp"].split()[0].replace("-", ""))
                if (curTime <= int(range_flag[1]) and 
                        curTime >= int(range_flag[0])):
                    retVal.append(response)
            else:
                retVal.append(response)
                
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
        if del_flag:
            response_bytes = self.retrieve_organization(org_id)
            
            if response_bytes != None:
                
                jresponse = json.loads(response_bytes.decode())
                
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
            response_bytes = self.retrieve_organization(org_id)
             
            if response_bytes != None:
                
                jresponse = json.loads(response_bytes.decode())
                
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
        return _sha512("organization".encode("utf-8"))[0:6]

    def _get_address(self, org_id):
        organization_prefix = self._get_prefix()
        address = _sha512(org_id.encode("utf-8"))[0:64]
        return organization_prefix + address
    
    def _get_block_num(self):
        organization_prefix = self._get_prefix()

        result = self._send_request(
            "blocks?={}".format(organization_prefix)
        )
        
        if result != None or result != "":
            result = json.loads(result)
            return str(len(result["data"]))
        return None
    
    def _get_payload_(self, blocknum):
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
        
        if self._base_url.startswith("http://"):
            url = "{}/{}".format(self._base_url, suffix)
        else:
            url = "http://{}/{}".format(self._base_url, suffix)

        headers = {}
        if content_type is not None:
            headers["Content-Type"] = content_type

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
        
        return result.text

    def create_organization_transaction(self, org_id, org_alias, org_name, 
                org_type, description, org_url, action, private_key, public_key, 
                prev, cur, timestamp, pt_id):
        
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

        # Construct the address
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

        batch_list = self._create_batch_list([transaction])
        
        return self._send_request(
            "batches", batch_list.SerializeToString(),
            "application/octet-stream"
        )

    def _create_batch_list(self, transactions):
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

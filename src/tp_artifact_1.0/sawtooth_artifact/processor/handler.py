# Copyright 2016 Intel Corporation
#
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

import hashlib
import logging
import json
from datetime import datetime
from collections import OrderedDict
# from sawtooth_sdk.processor.state import StateEntry
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
# from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader
from sawtooth_sdk.processor.handler import TransactionHandler


LOGGER = logging.getLogger(__name__)


class ArtifactTransactionHandler:
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return 'artifact'

    @property
    def family_versions(self):
        return ['1.0']

    @property
    def encodings(self):
        return ['csv-utf8']

    @property
    def namespaces(self):
        return [self._namespace_prefix]

    def apply(self, transaction, context):


        # header = TransactionHeader()
        # header.ParseFromString(transaction.header)

        try:
            # The payload is csv utf-8 encoded string
            artifact_id,alias,artifact_name,artifact_type,artifact_checksum,label,openchain,timestamp,action, sub_artifact_id,version,artifact_checksum_uri,content_type,size,uri_type,location,path = transaction.payload.decode().split(",")
        except ValueError:
            raise InvalidTransaction("Invalid payload serialization")

        validate_transaction(artifact_id,action)
               
        data_address = make_artifact_address(self._namespace_prefix,artifact_id)
        
        if  artifact_id  == "":
            raise InvalidTransaction("Artifact Data is required")

        if action == "":
            raise InvalidTransaction("Action is required")
          
        # state_entries = state_store.get([data_address])
        state_entries = context.get_state([data_address])

       
        if len(state_entries) != 0:
            try:
                   
                    stored_artifact_id, stored_artifact_str = \
                    state_entries[0].data.decode().split(",",1)
                             
                    stored_artifact = json.loads(stored_artifact_str)
            except ValueError:
                raise InternalError("Failed to deserialize data.")
 
        else:
            stored_artifact_id = stored_artifact = None
            
        # 3. Validate the artifact data
        if action == "create" and stored_artifact_id is not None:
            raise InvalidTransaction("Invalid Action-artifact already exists.")

        elif action == "AddArtifact":
            if stored_artifact_id is None:
                raise InvalidTransaction(
                    "Invalid Action-Add Artifact requires an existing artifact."
                )
        
        elif action == "AddURI":
            if stored_artifact_id is None:
                raise InvalidTransaction(
                    "Invalid Action-Add URI requires an existing artifact."
                ) 
     
        if action == "create":
            artifact = create_artifact(artifact_id,alias,artifact_name,artifact_type,artifact_checksum,label,openchain,timestamp)
            stored_artifact_id = artifact_id
            stored_artifact = artifact
           
        if action == "AddArtifact":
            if sub_artifact_id not in stored_artifact_str:
                artifact = add_artifact(sub_artifact_id,stored_artifact,path)
                stored_artifact = artifact
                
                
        if action == "AddURI":
            artifact = add_URI(stored_artifact,version,artifact_checksum_uri,content_type,size,uri_type,location)
            stored_artifact = artifact
            
            
      
        stored_artifact_str = json.dumps(stored_artifact)
        data=",".join([stored_artifact_id,stored_artifact_str]).encode()
        addresses = context.set_state({data_address:data})
       
        # addresses = state_store.set([
        #     StateEntry(
        #         address=data_address,
        #         data=",".join([stored_artifact_id, stored_artifact_str]).encode()
                
        #     )
        # ])
       
        return addresses
        
        
def add_artifact(uuid,parent_artifact,path):
    
    artifact_list = parent_artifact['artifact_list']
    artifact_dic = {'artifact_id': uuid,'path':path}
    artifact_list.append(artifact_dic)
    parent_artifact['artifact_list'] = artifact_list  
    return parent_artifact     

def add_URI(artifact,version,artifact_checksum,content_type,size,uri_type,location):
   
    URI_list = artifact['uri_list']
    URI_dic = {'version': version,'checksum': artifact_checksum,'content_type': content_type,'size':size,'uri_type':uri_type,'location':location}      
    URI_list.append(URI_dic)
    artifact['uri_list'] = URI_list
    return artifact

    

def create_artifact(uuid,alias,art_name,art_type,art_checksum,label,openchain,timestamp):
    artifact = {'artifact_id': uuid,'alias':alias,'artifact_name': art_name,'artifact_type': art_type,'artifact_checksum': art_checksum,'label':label,'openchain': openchain,'timestamp':timestamp,'artifact_list':[],'uri_list':[]}
    return artifact 


def validate_transaction( artifact_id, action):
    if not artifact_id:
        raise InvalidTransaction('Artifact ID is required')
    
    if not action:
        raise InvalidTransaction('Action is required')

    if action not in ("AddArtifact", "create","AddURI"):
        raise InvalidTransaction('Invalid action: {}'.format(action))

def make_artifact_address(namespace_prefix, artifact_id):
    return namespace_prefix + \
        hashlib.sha512(artifact_id.encode('utf-8')).hexdigest()[:64]

def _display(msg):
    n = msg.count("\n")

    if n > 0:
        msg = msg.split("\n")
        length = max(len(line) for line in msg)
    else:
        length = len(msg)
        msg = [msg]

    LOGGER.debug("+" + (length + 2) * "-" + "+")
    for line in msg:
        LOGGER.debug("+ " + line.center(length) + " +")
    LOGGER.debug("+" + (length + 2) * "-" + "+")

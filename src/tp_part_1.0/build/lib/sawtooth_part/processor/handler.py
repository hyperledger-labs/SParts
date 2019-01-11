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

import hashlib
import logging
import json
# from sawtooth_sdk.processor.state import StateEntry
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
# from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader
from sawtooth_sdk.processor.handler import TransactionHandler


LOGGER = logging.getLogger(__name__)


class PartTransactionHandler:
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return 'pt'

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
            pt_id,pt_name,checksum,version,alias,licensing,label,description,action,artifact_id,category_id,supplier_id = transaction.payload.decode().split(",")
        except ValueError:
            raise InvalidTransaction("Invalid payload serialization")
        
        validate_transaction( pt_id,action)  
             
        data_address = make_part_address(self._namespace_prefix,pt_id)
        
        # Retrieve the data from state storage  
        # state_entries = state_store.get([data_address])
        state_entries = context.get_state([data_address])
     
        if len(state_entries) != 0:
            try:
                   
                    stored_pt_id, stored_pt_str = \
                    state_entries[0].data.decode().split(",",1)
                             
                    stored_pt = json.loads(stored_pt_str)
            except ValueError:
                raise InternalError("Failed to deserialize data.")
            

        else:
            stored_pt_id = stored_pt = None
      
        
        if action == "create" and stored_pt_id is not None:
            raise InvalidTransaction("Invalid part already exists.")

        elif action == "AddArtifact" or action == "AddSupplier" or action == "AddCategory":
            if stored_pt_id is None:
                raise InvalidTransaction(
                    "Invalid the operation requires an existing part."
                )
           
        if action == "create":
            pt = create_part(pt_id,pt_name,checksum,version,alias,licensing,label,description)
            stored_pt_id = pt_id
            stored_pt = pt
            _display("Created a part.")
          
         
        if action == "AddArtifact":
            if artifact_id not in stored_pt_str:
                pt = add_artifact(artifact_id,stored_pt)
                stored_pt = pt
                
                
            
        if action == "AddSupplier":
            if supplier_id not in stored_pt_str:
                pt = add_supplier(supplier_id,stored_pt)
                stored_pt = pt
        
        if action == "AddCategory":
            if category_id not in stored_pt_str:
                pt = add_category(category_id,stored_pt)
                stored_pt = pt
        
         
        # 6. Put data back in state storage
        stored_pt_str = json.dumps(stored_pt)

        data=",".join([stored_pt_id,stored_pt_str]).encode()
        addresses = context.set_state({data_address:data})

        
        # addresses = state_store.set([
        #     StateEntry(
        #         address=data_address,
        #         data=",".join([stored_pt_id, stored_pt_str]).encode()
        #     )
        # ])
        return addresses

        
        
def add_artifact(uuid,parent_pt):
    
    pt_list = parent_pt['artifacts']
    pt_dic = {'artifact_id': uuid}
    pt_list.append(pt_dic)
    parent_pt['artifacts'] = pt_list
    return parent_pt  


def add_supplier(uuid,parent_pt):
    
    pt_list = parent_pt['suppliers']
    pt_dic = {'supplier_id': uuid}
    pt_list.append(pt_dic)
    parent_pt['suppliers'] = pt_list
    return parent_pt     

def add_category(uuid,parent_pt):
    
    pt_list = parent_pt['categories']
    pt_dic = {'category_id': uuid}
    pt_list.append(pt_dic)
    parent_pt['categories'] = pt_list
    return parent_pt        


def create_part(pt_id,pt_name,checksum,version,alias,licensing,label,description):
    ptD = {'pt_id': pt_id,'pt_name': pt_name,'checksum': checksum,'version': version,'alias':alias,'licensing':licensing,'label':label,'description':description,'artifacts':[],'suppliers':[],'categories':[]}
    return ptD 


def validate_transaction( pt_id,action):
    if not pt_id:
        raise InvalidTransaction('Part ID is required')
 
    if not action:
        raise InvalidTransaction('Action is required')

    if action not in ("AddArtifact", "create","AddCategory","AddSupplier","list-part","retrieve"):
        raise InvalidTransaction('Invalid action: {}'.format(action))


def make_part_address(namespace_prefix, part_id):
    return namespace_prefix + \
        hashlib.sha512(part_id.encode('utf-8')).hexdigest()[:64]


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

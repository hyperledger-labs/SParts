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
from collections import OrderedDict
from sawtooth_sdk.processor.state import StateEntry
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader

LOGGER = logging.getLogger(__name__)


class CategoryTransactionHandler:
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return 'category'

    @property
    def family_versions(self):
        return ['1.0']

    @property
    def encodings(self):
        return ['csv-utf8']

    @property
    def namespaces(self):
        return [self._namespace_prefix]

    def apply(self, transaction, state_store):

        header = TransactionHeader()
        header.ParseFromString(transaction.header)

        try:
            # The payload is csv utf-8 encoded string
            category_id,category_name,description,action = transaction.payload.decode().split(",")
        except ValueError:
            raise InvalidTransaction("Invalid payload")

        validate_transaction( category_id,category_name,description,action)
               
        data_address = create_category_address(self._namespace_prefix,category_id)
          
        state_entries = state_store.get([data_address])
        # Retrieve data from state storage
        if len(state_entries) != 0:
            try:
                   
                    stored_category_id, stored_category_str = \
                    state_entries[0].data.decode().split(",",1)                    
                    stored_category = json.loads(stored_category_str)
            except ValueError:
                raise InternalError("Failed to deserialize data.")
            
        else:
            stored_category_id = stored_category = None
            
        # Validate category data
        if action == "create" and stored_category_id is not None:
            raise InvalidTransaction("Invalid Action-category already exists.")
               
           
        if action == "create":
            category = create_category_payload(category_id,category_name,description)
            stored_category_id = category_id
            stored_category = category
            _display("Created a category.")
                
        # Insert data back
        stored_supp_str = json.dumps(stored_category)
        addresses = state_store.set([
            StateEntry(
                address=data_address,
                data=",".join([stored_category_id, stored_supp_str]).encode()
            )
        ])
        return addresses
        
        
def create_category_payload(category_id,category_name,description):
    categoryP = {'category_id': category_id,'category_name': category_name,'description': description}
    return categoryP 


def validate_transaction( category_id,category_name,description,action):
    if not category_id:
        raise InvalidTransaction('Category ID is required')
    
    if not action:
        raise InvalidTransaction('Action is required')

    if action not in ('create','list-category','retrieve'):
        raise InvalidTransaction('Invalid action: {}'.format(action))

    
def create_category_address(namespace_prefix, category_id):
    return namespace_prefix + \
        hashlib.sha512(category_id.encode('utf-8')).hexdigest()[:64]


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

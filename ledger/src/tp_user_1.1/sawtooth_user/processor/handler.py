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
#from sawtooth_sdk.processor.state import StateEntry
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
#from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader
from sawtooth_sdk.processor.handler import TransactionHandler

LOGGER = logging.getLogger(__name__)


class UserTransactionHandler:
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return 'user'

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
            user_public_key,user_name,email_address,authorized,role,action = transaction.payload.decode().split(",")
        except ValueError:
            raise InvalidTransaction("Invalid payload")

        validate_transaction(user_public_key,user_name,email_address,authorized,role,action)
               
        data_address = create_user_address(self._namespace_prefix,user_public_key)
          
        # state_entries = state_store.get([data_address])
        state_entries = context.get_state([data_address])
        # Retrieve data from state storage
        if len(state_entries) != 0:
            try:
                   
                    stored_user_id, stored_user_str = \
                    state_entries[0].data.decode().split(",",1)                    
                    stored_user = json.loads(stored_user_str)
            except ValueError:
                raise InternalError("Failed to deserialize data.")
            
        else:
            stored_user_id = stored_user = None
            
        # Validate user data
        if action == "register" and stored_user_id is not None:
            raise InvalidTransaction("Invalid Action-user already exists.")
               
           
        if action == "register":
            user = create_user_payload(user_public_key,user_name,email_address,authorized,role)
            stored_user_id = user_public_key
            stored_user = user
            _display("Register a user.")
                
        # Insert data back
        stored_userAccount_str = json.dumps(stored_user)
        data=",".join([stored_user_id,stored_userAccount_str]).encode()
        addresses = context.set_state({data_address:data})
        # addresses = state_store.set([
        #     StateEntry(
        #         address=data_address,
        #         data=",".join([stored_user_id, stored_userAccount_str]).encode()
        #     )
        # ])
        return addresses
        
        
def create_user_payload(user_public_key,user_name,email_address,authorized,role):
    userP = {'public_key': user_public_key,'user_name': user_name,'email_address': email_address,'authorized':authorized,'role': role}
    return userP 


def validate_transaction(user_public_key,user_name,email_address,authorized,role,action):
    if not user_public_key:
        raise InvalidTransaction('User Public Key is required')
    
    if not action:
        raise InvalidTransaction('Action is required')

    if action not in ('register','list-user','retrieve'):
        raise InvalidTransaction('Invalid action: {}'.format(action))

    
def create_user_address(namespace_prefix, user_public_key):
    return namespace_prefix + \
        hashlib.sha512(user_public_key.encode('utf-8')).hexdigest()[:64]


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

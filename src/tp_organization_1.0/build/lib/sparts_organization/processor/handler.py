
# Copyright 2017 Wind River Systems
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
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
import hashlib
import logging
import json
from collections import OrderedDict
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
from sawtooth_sdk.processor.handler import TransactionHandler

LOGGER = logging.getLogger(__name__)
################################################################################
#                               HANDLER OBJ                                    #
################################################################################
class OrganizationTransactionHandler:
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return 'organization'

    @property
    def family_versions(self):
        return ['1.0']

    @property
    def encodings(self):
        return ['csv-utf8']

    @property
    def namespaces(self):
        return [self._namespace_prefix]
################################################################################
#                                 FUNCTIONS                                    #
################################################################################
    def apply(self, transaction, context):
        
        try:
            payload = json.loads(transaction.payload.decode())
            org_id      = payload["uuid"]
            org_alias   = payload["alias"]
            org_name    = payload["name"]
            org_type    = payload["type"]
            description = payload["description"]
            org_url     = payload["url"]
            action      = payload["action"]
            prev        = payload["prev_block"]
            cur         = payload["cur_block"]
            timestamp   = payload["timestamp"]
            pt_id       = payload["pt_list"]
            
        except ValueError:
            raise InvalidTransaction("Invalid payload serialization")

        validate_transaction(org_id, action)
               
        data_address = make_organization_address(self._namespace_prefix, org_id)
        
        state_entries = context.get_state([data_address])
        
        if len(state_entries) != 0:
            try:

                stored_organization = json.loads(state_entries[0].data.decode())
                stored_organization_id = stored_organization["uuid"]
                
            except ValueError:
                raise InternalError("Failed to deserialize data.")
            
        else:
            stored_organization_id = stored_organization = None
            
        if action == "create" and stored_organization_id is not None:
            raise InvalidTransaction(
                        "Invalid Action-organization already exists."
                    )
    
        elif action == "create":
            organization = create_organization(org_id, org_alias, org_name, 
                                org_type, description, org_url, prev, cur, 
                                timestamp)
            _display("Created an organization.")
        elif action == "amend" and stored_organization_id is not None:
            organization = create_organization(org_id, org_alias, org_name, 
                                org_type, description, org_url, prev, cur, 
                                timestamp, pt_id)
            _display("Amended an organization.")
        elif action == "AddPart" and stored_organization_id is not None:
            organization = create_organization(org_id, org_alias, org_name, 
                                org_type, description, org_url, prev, cur, 
                                timestamp, pt_id)
            
        data = json.dumps(organization).encode()
        addresses = context.set_state({data_address:data})
        
        return addresses
  

def create_organization(org_id, org_alias, org_name, org_type, description, 
                        org_url, prev, cur, timestamp, pt_id=[]):
    
    return {
                "uuid"          : org_id, 
                "alias"         : org_alias, 
                "name"          : org_name, 
                "type"          : org_type, 
                "description"   : description, 
                "url"           : org_url,
                "prev_block"    : prev, 
                "cur_block"     : cur,
                "timestamp"     : timestamp,
                "pt_list"       : pt_id
            } 

def validate_transaction(org_id, action):
    if not org_id:
        raise InvalidTransaction('Organization ID is required') 
    if not action:
        raise InvalidTransaction('Action is required')
    if action not in ('create', "amend", "AddPart"):
        raise InvalidTransaction('Invalid action: {}'.format(action))

def make_organization_address(namespace_prefix, org_id):
    return namespace_prefix + \
        hashlib.sha512(org_id.encode('utf-8')).hexdigest()[:64]

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
################################################################################
#                                                                              #
################################################################################

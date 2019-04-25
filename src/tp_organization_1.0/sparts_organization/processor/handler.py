# Copyright 2016 Intel Corporation
# Copyright 2017 Wind River Systems
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
#                         LIBRARIES & DEPENDENCIES                             #
################################################################################
import hashlib
import logging
import json
from sawtooth_sdk.processor.exceptions import InvalidTransaction
from sawtooth_sdk.processor.exceptions import InternalError
from sawtooth_sdk.processor.handler import TransactionHandler

LOGGER = logging.getLogger(__name__)
################################################################################
#                               HANDLER OBJ                                    #
################################################################################
class OrganizationTransactionHandler:
    """
    Class for handling the Transaction Family : Organization
    
    Attributes:
        namespace_prefix (str): The namespace prefix of the transaction family
        
    """
    
    def __init__(self, namespace_prefix):
        """
        Constructs the OrganizationTransactionHandler object.
        
        Args:
            namespace_prefix (str):
                The namepsace prefix of the transaction family
                
        """
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        """
        type: str
        Returns the family name of the handler object.
        
        """
        return "organization"

    @property
    def family_versions(self):
        """
        type: list of str
        Returns the family version of the handler object.
        
        """
        return ["1.0"]

    @property
    def encodings(self):
        """
        type: list of str
        Returns the encoding scheme used for the data for the handler object.
        
        """
        return ["csv-utf8"]

    @property
    def namespaces(self):
        """
        type: list of str
        Returns the namespaces associating with the handler object.
        
        """
        return [self._namespace_prefix]
################################################################################
#                                 FUNCTIONS                                    #
################################################################################
    def apply(self, transaction, context):
        """
        Applys the payload from transaction onto the state storage.
        
        Args:
            transaction (Transaction): The transaction pertaining the payload
            context (State): The current state of the ledger
            
        Returns:
            type: State
            The new state of the ledger, which includes the data from the
            transaction, is returned to be stored on the state storage.
        
        Raises:
            InvalidTransaction:
                * If deserialization for payload from transaction failed
                * If "create" was called on non-unique uuid
                * If "amend" was called on non-existing uuid
                * If "Add..." were called on non-existing uuid
                * If invalid operation was called
            InternalError:
                * If deserialization of State.data failed
            
        """
        
        # Parsing required fields from transaction payload
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
        
        # Soft sanity check and loading required data
        validate_transaction(org_id, action)
        data_address = make_organization_address(self._namespace_prefix, org_id)
        state_entries = context.get_state([data_address])
        
        # Hard sanity check before creating final payload for the state storage
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
        
        # Adding the final payload to the state storage   
        data = json.dumps(organization).encode()
        addresses = context.set_state({data_address:data})
        
        return addresses
  

def create_organization(org_id, org_alias, org_name, org_type, description, 
                        org_url, prev, cur, timestamp, pt_id=[]):
    """
    Constructs the payload to be stored in the state storage.
    
    Args:
        org_id (str): The uuid of the organization
        org_alias (str): The alias of the organization
        org_name (str): The name of the organization
        org_type (str): The type of the organization
        description (str): The description of the organization
        org_url (str): The url of the organization
        prev (str): The previous block id of the transaction (default "0")
        cur (str): the current block id of the transaction
        timestamp (str): The UTC time for when the transaction was submitted
        pt_id (list of str):
            The list of the part uuid associated with the organization
            (default [])
        
    Returns:
        type: dict
        The dictionary pertaining all the param is created and returned to
        be stored on the state storage.
    
    """
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
    """
    Performs soft sanity check in order to improve runtime by eliminating the
    obvious exception errors.
    
    Args:
        org_id (str): The uuid of the organization
        action (str): The command to be performed
    
    Raises:
        InvalidTransaction:
            If the uuid or the action are not passed in or the 
            action is not a valid action.
    
    """
    if not org_id:
        raise InvalidTransaction("Organization ID is required") 
    if not action:
        raise InvalidTransaction("Action is required")
    if action not in ("create", "amend", "AddPart"):
        raise InvalidTransaction("Invalid action: {}".format(action))

def make_organization_address(namespace_prefix, org_id):
    """
    Creates an organization address which will be used to recover the associated
    UUID if the part already exists in the state storage; or, used as a key to
    store the new data into the state storage.
    
    Args:
        namespace_prefix (str):
            The prefix associating with the transaction family
        org_id (str): The uuid of the organization
        
    Returns:
        type: str
        The address-to-be, which associates the uuid and the namespace prefix.
    
    """
    return namespace_prefix + \
        hashlib.sha512(org_id.encode("utf-8")).hexdigest()[:64]

def _display(msg):
    """
    Logs the message to the debug logger.
    
    Args:
        msg (str): The message that is to be logged into the debug logger
    
    """
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

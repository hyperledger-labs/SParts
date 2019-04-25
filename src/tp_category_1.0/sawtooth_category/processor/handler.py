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
class CategoryTransactionHandler(TransactionHandler):
    """
    Class for handling the Transaction Family : Category
    
    Attributes:
        namespace_prefix (str): The namespace prefix of the transaction family
        
    """
    
    def __init__(self, namespace_prefix):
        """
        Constructs the CategoryTransactionHandler object.
        
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
        return "category"

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
                * If invalid operation was called
            InternalError:
                * If deserialization of State.data failed
            
        """
        
        # Parsing required fields from transaction payload
        try:
            
            payload = json.loads(transaction.payload.decode())
            category_id     = payload["uuid"]
            category_name   = payload["name"]
            description     = payload["description"]
            action          = payload["action"]
            prev            = payload["prev_block"]
            cur             = payload["cur_block"]
            timestamp       = payload["timestamp"]
            
        except ValueError:
            raise InvalidTransaction("Invalid payload")
        
        # Soft sanity check and loading required data
        validate_transaction(category_id, action)
        data_address = create_category_address(self._namespace_prefix, 
                                                category_id)
        state_entries = context.get_state([data_address])
        
        # Hard sanity check before creating final payload for the state storage
        if len(state_entries) != 0:
            try:

                stored_category = json.loads(state_entries[0].data.decode())
                stored_category_id = stored_category["uuid"]
                
            except ValueError:
                raise InternalError("Failed to deserialize data.")

        else:
            stored_category_id = stored_category = None

        if action == "create" and stored_category_id is not None:
            raise InvalidTransaction("Invalid Action-category already exists.")
        elif action == "create":
            category = create_category_payload(category_id, category_name, 
                                    description, prev, cur, timestamp)
            _display("Created a category.")
        elif action == "amend" and stored_category_id is not None:
            category = create_category_payload(category_id, category_name, 
                                    description, prev , cur, timestamp)
            _display("Amended a category.")
        else:
            raise InvalidTransaction(
                "Invalid operation."
            )
        
        # Adding the final payload to the state storage
        data = json.dumps(category).encode()
        addresses = context.set_state({data_address:data})
     
        return addresses
################################################################################
#                             HELPER FUNCTIONS                                 #
################################################################################
def create_category_payload(category_id, category_name, description, 
                            prev, cur, timestamp):
    """
    Constructs the payload to be stored in the state storage.
    
    Args:
        category_id (str): The uuid of the category
        category_name (str): The name of the category
        description (str): The description of the category
        prev (str): The previous block id of the transaction (default "0")
        cur (str): the current block id of the transaction
        timestamp (str): The UTC time for when the transaction was submitted
        
    Returns:
        type: dict
        The dictionary pertaining all the param is created and returned to
        be stored on the state storage.
    
    """
    return {
                "uuid"          : category_id,
                "name"          : category_name,
                "description"   : description,
                "prev_block"    : prev,
                "cur_block"     : cur,
                "timestamp"     : timestamp
            }


def validate_transaction(category_id, action):
    """
    Performs soft sanity check in order to improve runtime by eliminating the
    obvious exception errors.
    
    Args:
        category_id (str): The uuid of the category
        action (str): The command to be performed
    
    Raises:
        InvalidTransaction:
            If the uuid or the action are not passed in or the 
            action is not a valid action.
    
    """
    if not category_id:
        raise InvalidTransaction("Category ID is required")
    if not action:
        raise InvalidTransaction("Action is required")
    if action not in ("create", "list-category", "retrieve", "amend"):
        raise InvalidTransaction("Invalid action: {}".format(action))


def create_category_address(namespace_prefix, category_id):
    """
    Creates a category address which will be used to recover the associated UUID
    if the category already exists in the state storage; or, used as a key to
    store the new data into the state storage.
    
    Args:
        namespace_prefix (str):
            The prefix associating with the transaction family
        part_id (str): The uuid of the category
        
    Returns:
        type: str
        The address-to-be, which associates the uuid and the namespace prefix.
    
    """
    return namespace_prefix + \
        hashlib.sha512(category_id.encode("utf-8")).hexdigest()[:64]


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

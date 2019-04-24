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
class ArtifactTransactionHandler:
    """
    Class for handling the Transaction Family : Artifact
    
    Attributes:
        namespace_prefix (str): The namespace prefix of the transaction family
        
    """
    
    def __init__(self, namespace_prefix):
        """
        Constructs the ArtifactTransactionHandler object.
        
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
        return "artifact"

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
            artifact_id             = payload["uuid"]
            artifact_alias          = payload["alias"]
            artifact_name           = payload["name"]
            artifact_type           = payload["content_type"]
            artifact_checksum       = payload["checksum"]
            artifact_label          = payload["label"]
            artifact_openchain      = payload["openchain"]
            action                  = payload["action"]
            prev                    = payload["prev_block"]
            cur                     = payload["cur_block"]
            timestamp               = payload["timestamp"]
            artifact_list           = payload["artifact_list"]
            uri_list                = payload["uri_list"]
            
        except ValueError:
            raise InvalidTransaction("Invalid payload serialization")
        
        # Soft sanity check and loading required data
        validate_transaction(artifact_id, action)
        data_address = make_artifact_address(self._namespace_prefix, 
                                                    artifact_id)
        state_entries = context.get_state([data_address])
        
        # Hard sanity check before creating final payload for the state storage
        if len(state_entries) != 0:
            try:
                
                stored_artifact = json.loads(state_entries[0].data.decode())
                stored_artifact_id = stored_artifact["uuid"]
                             
            except ValueError:
                raise InternalError("Failed to deserialize data.")
 
        else:
            stored_artifact_id = stored_artifact = None
            
        if action == "create" and stored_artifact_id is not None:
            raise InvalidTransaction("Invalid Action-artifact already exists.")
        
        elif action == "create":
            artifact = create_artifact(artifact_id, artifact_alias, 
                            artifact_name, artifact_type, artifact_checksum, 
                            artifact_label, artifact_openchain, 
                            prev, cur, timestamp)
        elif action == "amend" and stored_artifact_id is not None:
            artifact = create_artifact(artifact_id, artifact_alias, 
                            artifact_name, artifact_type, artifact_checksum, 
                            artifact_label, artifact_openchain,
                            prev, cur, timestamp, artifact_list, uri_list)
        elif action == "AddArtifact" or action == "AddURI":
            if stored_artifact_id is None:
                raise InvalidTransaction(
                    "Invalid Action-requires an existing artifact."
                )
            artifact = create_artifact(artifact_id, artifact_alias, 
                            artifact_name, artifact_type, artifact_checksum, 
                            artifact_label, artifact_openchain, 
                            prev, cur, timestamp, 
                            artifact_list, uri_list)
        
        # Adding the final payload to the state storage    
        data = json.dumps(artifact).encode()
        addresses = context.set_state({data_address:data})
       
        return addresses
################################################################################
#                             HELPER FUNCTIONS                                 #
################################################################################
def create_artifact(artifact_id, artifact_alias, artifact_name, artifact_type, 
                    artifact_checksum, artifact_label, artifact_openchain, 
                    prev, cur, timestamp, artifact_list=[], uri_list=[]):
    """
    Constructs the payload to be stored in the state storage.
    
    Args:
        artifact_uuid (str): The uuid of the artifact
        artifact_alias (str): The alias of the artifact
        artifact_name (str): The name of the artifact
        artifact_type (str): The type of the artifact
        artifact_checksum (str): The checksum of the artifact
        artifact_label (str): The label of the artifact
        artifact_openchain (str): The openchain of the artifact
        prev (str): The previous block id of the transaction (default "0")
        cur (str): the current block id of the transaction
        timestamp (str): The UTC time for when the transaction was submitted
        artifact_list (list of dict):
            The list of the artifact uuid associated with the artifact
            (default [])
        uri_list (list of dict):
            The list of the uri associated with the artifact (default [])
        
    Returns:
        type: dict
        The dictionary pertaining all the param is created and returned to
        be stored on the state storage.
    
    """
    return {    
                "uuid"          : artifact_id,
                "alias"         : artifact_alias,
                "name"          : artifact_name,
                "content_type"  : artifact_type,
                "checksum"      : artifact_checksum,
                "label"         : artifact_label,
                "openchain"     : artifact_openchain,
                "prev_block"    : prev, 
                "cur_block"     : cur,
                "timestamp"     : timestamp,
                "artifact_list" : artifact_list,
                "uri_list"      : uri_list
            }

def validate_transaction(artifact_id, action):
    """
    Performs soft sanity check in order to improve runtime by eliminating the
    obvious exception errors.
    
    Args:
        artifact_id (str): The uuid of the artifact
        action (str): The command to be performed
    
    Raises:
        InvalidTransaction:
            If the uuid or the action are not passed in or the 
            action is not a valid action.
    
    """
    if not artifact_id:
        raise InvalidTransaction("Artifact ID is required")
    if not action:
        raise InvalidTransaction("Action is required")
    if action not in ("AddArtifact", "create", "AddURI", "amend"):
        raise InvalidTransaction("Invalid action: {}".format(action))

def make_artifact_address(namespace_prefix, artifact_id):
    """
    Creates an artifact address which will be used to recover the associated
    UUID if the artifact already exists in the state storage; or, used as a key to
    store the new data into the state storage.
    
    Args:
        namespace_prefix (str):
            The prefix associating with the transaction family
        artifact_id (str): The uuid of the artifact
        
    Returns:
        type: str
        The address-to-be, which associates the uuid and the namespace prefix.
    
    """
    return namespace_prefix + \
        hashlib.sha512(artifact_id.encode("utf-8")).hexdigest()[:64]

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

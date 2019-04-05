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
#                                STATE OBJ                                     #
################################################################################
class StateEntry:
    def __init__(self, address, data):
        self.address = address
        self.data = data
################################################################################
#                               HANDLER OBJ                                    #
################################################################################
class CategoryTransactionHandler(TransactionHandler):
    def __init__(self, namespace_prefix):
        self._namespace_prefix = namespace_prefix

    @property
    def family_name(self):
        return "category"

    @property
    def family_versions(self):
        return ["1.0"]

    @property
    def encodings(self):
        return ["csv-utf8"]

    @property
    def namespaces(self):
        return [self._namespace_prefix]

    def apply(self, transaction, context):

        try:
            # The payload is csv utf-8 encoded string
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

        validate_transaction(category_id, category_name, description, action)

        data_address = create_category_address(self._namespace_prefix, 
                                                category_id)

        state_entries = context.get_state([data_address])
        
        # Retrieve data from state storage
        if len(state_entries) != 0:
            try:

                stored_category = json.loads(state_entries[0].data.decode())
                stored_category_id = stored_category["uuid"]
                
            except ValueError:
                raise InternalError("Failed to deserialize data.")

        else:
            stored_category_id = stored_category = None

        # Validate category data
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
        
        data = json.dumps(category).encode()
        addresses = context.set_state({data_address:data})
     
        return addresses


def create_category_payload(category_id, category_name, description, 
                            prev, cur, timestamp):
    return {
                "uuid"          : category_id,
                "name"          : category_name,
                "description"   : description,
                "prev_block"    : prev,
                "cur_block"     : cur,
                "timestamp"     : timestamp
            }


def validate_transaction(category_id, category_name, description, action):
    if not category_id:
        raise InvalidTransaction("Category ID is required")

    if not action:
        raise InvalidTransaction("Action is required")

    if action not in ("create", "list-category", "retrieve", "amend"):
        raise InvalidTransaction("Invalid action: {}".format(action))


def create_category_address(namespace_prefix, category_id):
    return namespace_prefix + \
        hashlib.sha512(category_id.encode("utf-8")).hexdigest()[:64]


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

# Copyright 2016 Intel Corporation
# Copyright 2017 Wind River
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
import base64
from base64 import b64encode
import time
import requests
import yaml

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

from sawtooth_user.exceptions import UserException


def _sha512(data):
    return hashlib.sha512(data).hexdigest()


class UserBatch:
    def __init__(self, base_url):

        self._base_url = base_url

    def register_user(self, user_public_key,user_name,email_address,authorized,role,adprivatekey,adpublickey):
        return self.send_user_transactions(user_public_key,user_name,email_address,authorized,role, "register",adprivatekey,adpublickey)
    
    
    #Remove this function
    def list_user(self):
        user_prefix = self._get_prefix()

        result = self._send_request(
            "state?address={}".format(user_prefix)   
        )
        try:
            encoded_entries = yaml.safe_load(result)["data"]

            return [
                base64.b64decode(entry["data"]) for entry in encoded_entries
            ]

        except BaseException:
            return None

    def retreive_user(self,  user_public_key):
        address = self._get_address( user_public_key)

        result = self._send_request("state/{}".format(address),  user_public_key= user_public_key,
                                   )
        try:
            return base64.b64decode(yaml.safe_load(result)["data"])

        except BaseException:
            return None

    def _get_prefix(self):
        return _sha512('user'.encode('utf-8'))[0:6]

    def _get_address(self,  user_public_key):
        user_prefix = self._get_prefix()
        address = _sha512( user_public_key.encode('utf-8'))[0:64]
        return user_prefix + address

    def _send_request(
            self, suffix, data=None,
            content_type=None,  user_public_key=None):
        if self._base_url.startswith("http://"):
            url = "{}/{}".format(self._base_url, suffix)
        else:
            url = "http://{}/{}".format(self._base_url, suffix)

        headers = {}
       
        if content_type is not None:
            headers['Content-Type'] = content_type

        try:
            if data is not None:
                result = requests.post(url, headers=headers, data=data)
            else:
                result = requests.get(url, headers=headers)

            if result.status_code == 404:
                raise UserException("No such user: {}".format(user_public_key))

            elif not result.ok:
                raise UserException("Error {}: {}".format(
                    result.status_code, result.reason))

        except BaseException as err:
            raise UserException(err)

        return result.text

    def send_user_transactions(self, user_public_key,user_name,email_address,authorized,role, action,adprivatekey,adpublickey
                     ):
        
        self._public_key = adpublickey
        self._private_key = adprivatekey
        payload = ",".join([user_public_key,user_name,email_address,authorized,role, action]).encode()

        # Form the address
        address = self._get_address(user_public_key)

        header = TransactionHeader(
            signer_public_key= self._public_key,
            family_name="user",
            family_version="1.0",
            inputs=[address],
            outputs=[address],
            dependencies=[],
            # payload_encoding="csv-utf8",
            payload_sha512=_sha512(payload),
            batcher_public_key= self._public_key,
            nonce=time.time().hex().encode()
        ).SerializeToString()

        # signature = signing.sign(header, self._private_key)

        signature = CryptoFactory(create_context('secp256k1')) \
            .new_signer(Secp256k1PrivateKey.from_hex(self._private_key)).sign(header)

        transaction = Transaction(
            header=header,
            payload=payload,
            header_signature=signature
        )

        batch_list = self._create_batch_list([transaction])
        
        return self._send_request(
            "batches", batch_list.SerializeToString(),
            'application/octet-stream',
          
        )

    def _create_batch_list(self, transactions):
        transaction_signatures = [t.header_signature for t in transactions]

        header = BatchHeader(
            signer_public_key=self._public_key,
            transaction_ids=transaction_signatures
        ).SerializeToString()

        # signature = signing.sign(header, self._private_key)

        signature = CryptoFactory(create_context('secp256k1')) \
            .new_signer(Secp256k1PrivateKey.from_hex(self._private_key)).sign(header)
        batch = Batch(
            header=header,
            transactions=transactions,
            header_signature=signature
        )
        return BatchList(batches=[batch])

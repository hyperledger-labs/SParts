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
import base64
from base64 import b64encode
import time
import requests
import yaml

import sawtooth_signing.secp256k1_signer as signing

from sawtooth_sdk.protobuf.transaction_pb2 import TransactionHeader
from sawtooth_sdk.protobuf.transaction_pb2 import Transaction
from sawtooth_sdk.protobuf.batch_pb2 import BatchList
from sawtooth_sdk.protobuf.batch_pb2 import BatchHeader
from sawtooth_sdk.protobuf.batch_pb2 import Batch

from sawtooth_supplier.exceptions import SupplierException


def _sha512(data):
    return hashlib.sha512(data).hexdigest()


class SupplierBatch:
   
    def __init__(self, base_url, keyfile):

        self._base_url = base_url

        try:
            with open(keyfile) as fd:
                self._private_key = fd.read().strip()
                fd.close()
        except:
            raise IOError("Failed to read keys.")

        self._public_key = signing.generate_pubkey(self._private_key)

    def create(self,supplier_id,short_id,supplier_name,passwd,supplier_url, auth_user=None, auth_password=None):
        return self.create_supplier_transaction(supplier_id,short_id,supplier_name,passwd,supplier_url, "create",
                                 auth_user=auth_user,
                                 auth_password=auth_password)

  
    def add_part(self,supplier_id,part_id):
        return self.create_supplier_transaction(supplier_id,"","","","","AddPart",part_id)

    def list_supplier(self, auth_user=None, auth_password=None):
        supplier_prefix = self._get_prefix()

        result = self._send_request(
            "state?address={}".format(supplier_prefix),
            auth_user=auth_user,
            auth_password=auth_password
        )

        try:
            encoded_entries = yaml.safe_load(result)["data"]

            return [
                base64.b64decode(entry["data"]) for entry in encoded_entries
            ]

        except BaseException:
            return None

    def retrieve_supplier(self, supplier_id, auth_user=None, auth_password=None):
        address = self._get_address(supplier_id)

        result = self._send_request("state/{}".format(address), supplier_id=supplier_id,
                                    auth_user=auth_user,
                                    auth_password=auth_password)
        try:
            return base64.b64decode(yaml.safe_load(result)["data"])

        except BaseException:
            return None

 

    def _get_prefix(self):
        return _sha512('supplier'.encode('utf-8'))[0:6]

    def _get_address(self, supplier_id):
        supplier_prefix = self._get_prefix()
        address = _sha512(supplier_id.encode('utf-8'))[0:64]
        return supplier_prefix + address

    def _send_request(
            self, suffix, data=None,
            content_type=None, supplier_id=None, auth_user=None, auth_password=None):
        if self._base_url.startswith("http://"):
            url = "{}/{}".format(self._base_url, suffix)
        else:
            url = "http://{}/{}".format(self._base_url, suffix)

        headers = {}
        if auth_user is not None:
            auth_string = "{}:{}".format(auth_user, auth_password)
            b64_string = b64encode(auth_string.encode()).decode()
            auth_header = 'Basic {}'.format(b64_string)
            headers['Authorization'] = auth_header

        if content_type is not None:
            headers['Content-Type'] = content_type

        try:
            if data is not None:
                result = requests.post(url, headers=headers, data=data)
            else:
                result = requests.get(url, headers=headers)

            if result.status_code == 404:
                raise SupplierException("No such supplier: {}".format(supplier_id))

            elif not result.ok:
                raise SupplierException("Error {}: {}".format(
                    result.status_code, result.reason))

        except BaseException as err:
            raise SupplierException(err)

        return result.text

    def create_supplier_transaction(self, supplier_id,short_id="",supplier_name="",passwd="",supplier_url="", action="",part_id="",
                     auth_user=None, auth_password=None):
        
        payload = ",".join([supplier_id,str(short_id),str(supplier_name),str(passwd),str(supplier_url), action,str(part_id)]).encode()

        # Construct the address
        address = self._get_address(supplier_id)

        header = TransactionHeader(
            signer_pubkey=self._public_key,
            family_name="supplier",
            family_version="1.0",
            inputs=[address],
            outputs=[address],
            dependencies=[],
            payload_encoding="csv-utf8",
            payload_sha512=_sha512(payload),
            batcher_pubkey=self._public_key,
            nonce=time.time().hex().encode()
        ).SerializeToString()

        signature = signing.sign(header, self._private_key)

        transaction = Transaction(
            header=header,
            payload=payload,
            header_signature=signature
        )

        batch_list = self._create_batch_list([transaction])
        batch_id = batch_list.batches[0].header_signature
        
        return self._send_request(
            "batches", batch_list.SerializeToString(),
            'application/octet-stream',
            auth_user=auth_user,
            auth_password=auth_password
        )

    def _create_batch_list(self, transactions):
        transaction_signatures = [t.header_signature for t in transactions]

        header = BatchHeader(
            signer_pubkey=self._public_key,
            transaction_ids=transaction_signatures
        ).SerializeToString()

        signature = signing.sign(header, self._private_key)

        batch = Batch(
            header=header,
            transactions=transactions,
            header_signature=signature
        )
        return BatchList(batches=[batch])

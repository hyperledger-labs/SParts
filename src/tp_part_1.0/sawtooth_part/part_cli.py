# Copyright 2017 Intel Corporation
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

from __future__ import print_function

import argparse
import configparser
import getpass
import logging
import os
import traceback
import sys
import shutil
import pkg_resources
import json
import re
import requests

from colorlog import ColoredFormatter

# import sawtooth_signing.secp256k1_signer as signing

#
import binascii
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
#

from sawtooth_part.part_batch import PartBatch
from sawtooth_part.exceptions import PartException


DISTRIBUTION_NAME = 'sawtooth-part'


def create_console_handler(verbose_level):
    clog = logging.StreamHandler()
    formatter = ColoredFormatter(
        "%(log_color)s[%(asctime)s %(levelname)-8s%(module)s]%(reset)s "
        "%(white)s%(message)s",
        datefmt="%H:%M:%S",
        reset=True,
        log_colors={
            'DEBUG': 'cyan',
            'INFO': 'green',
            'WARNING': 'yellow',
            'ERROR': 'red',
            'CRITICAL': 'red',
        })

    clog.setFormatter(formatter)

    if verbose_level == 0:
        clog.setLevel(logging.WARN)
    elif verbose_level == 1:
        clog.setLevel(logging.INFO)
    else:
        clog.setLevel(logging.DEBUG)

    return clog


def setup_loggers(verbose_level):
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    logger.addHandler(create_console_handler(verbose_level))


def add_create_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('create', parents=[parent_parser])

    parser.add_argument(
        'pt_id',
        type=str,
        help='an identifier for the part')
    
    parser.add_argument(
        'pt_name',
        type=str,
        help='provide part name')
    
    parser.add_argument(
        'checksum',
        type=str,
        help='Provide checksum')
    
    parser.add_argument(
        'version',
        type=str,
        help='provide version for the part')
    
    parser.add_argument(
        'alias',
        type=str,
        help='provide alias')
    
    parser.add_argument(
        'licensing',
        type=str,
        help='provide licensing')
    
    parser.add_argument(
        'label',
        type=str,
        help='provide label')
    
    parser.add_argument(
        'description',
        type=str,
        help='provide description')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')


    parser.add_argument(
        '--disable-client-validation',
        action='store_true',
        default=False,
        help='disable client validation')



def add_list_part_parser(subparsers, parent_parser):
    subparsers.add_parser('list-part', parents=[parent_parser])


def add_retrieve_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('retrieve', parents=[parent_parser])

    parser.add_argument(
        'pt_id',
        type=str,
        help='part identifier')
    


def add_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddArtifact', parents=[parent_parser])
    
    parser.add_argument(
        'pt_id',
        type=str,
        help='part identifier')

    parser.add_argument(
        'artifact_id',
        type=str,
        help='the UUID identifier for artifact')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')


def create_parent_parser(prog_name):
    parent_parser = argparse.ArgumentParser(prog=prog_name, add_help=False)
    parent_parser.add_argument(
        '-v', '--verbose',
        action='count',
        help='enable more verbose output')

    try:
        version = pkg_resources.get_distribution(DISTRIBUTION_NAME).version
    except pkg_resources.DistributionNotFound:
        version = 'UNKNOWN'

    parent_parser.add_argument(
        '-V', '--version',
        action='version',
        version=(DISTRIBUTION_NAME + ' (Hyperledger Sawtooth) version {}')
        .format(version),
        help='print version information')

    return parent_parser


def create_parser(prog_name):
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        parents=[parent_parser],
        formatter_class=argparse.RawDescriptionHelpFormatter)

    subparsers = parser.add_subparsers(title='subcommands', dest='command')

    add_create_parser(subparsers, parent_parser)
   
    add_list_part_parser(subparsers, parent_parser)
    add_retrieve_parser(subparsers, parent_parser)
    add_artifact_parser(subparsers, parent_parser)
    add_supplier_parser(subparsers,parent_parser)
    add_category_parser(subparsers,parent_parser)

    return parser


def add_category_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddCategory', parents=[parent_parser])
    
    parser.add_argument(
        'pt_id',
        type=str,
        help='the identifier for the part')

    parser.add_argument(
        'category_id',
        type=str,
        help='the identifier for Supplier')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')

# Provide the UUID of the parent artifact and the UUID of the supplier
def add_supplier_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddSupplier', parents=[parent_parser])
    
    parser.add_argument(
        'pt_id',
        type=str,
        help='the identifier for the part')

    parser.add_argument(
        'supplier_id',
        type=str,
        help='the identifier for Supplier')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')
    
def do_list_part(args, config):
    b_url = config.get('DEFAULT', 'url')
   
    client = PartBatch(base_url=b_url)
    result = client.list_part()

    if result is not None:
        result = refine_output(str(result))
        output = ret_msg("success","OK","ListOf:PartRecord",result)
        print(output)
    else:
        raise PartException("Could not retrieve part listing.")


def do_retrieve(args, config):
    
    pt_id = args.pt_id

    b_url = config.get('DEFAULT', 'url')
   
    client = PartBatch(base_url=b_url)

    result = client.retrieve_part(pt_id).decode()

    if result is not None:
        result = filter_output(str(result))
        output = ret_msg("success","OK","PartRecord",result)
        print(result)
     
    else:
        raise PartException("Part not found: {}".format(pt_id))



def filter_output(inputstr):
    
    ptlist = inputstr.split(',',1)
    ptstr = ptlist[1]
    jsonstr = ptstr.replace('pt_id','uuid').replace('pt_name','name')
    data = json.loads(jsonstr)
    jsonstr = json.dumps(data)
    return jsonstr


def amend_part_fields(inputstr):
    output = inputstr.replace("\\","").replace('pt_id','uuid').replace('pt_name','name')
    return output

        
def refine_output(inputstr):
    inputstr = inputstr[1:-1]
    outputstr = inputstr.replace('b\'','').replace('}\'','}')  
    outputstr = outputstr[:-1]
    slist = outputstr.split("},")
    supplierlist = []
    for line in slist:
        record = "{"+line.split(",{",1)[-1]+"}"
        supplierlist.append(record)
    joutput = str(supplierlist)
    joutput = joutput.replace("'{","{").replace("}'","}").replace(", { {",", {")
    joutput = amend_part_fields(joutput)
    
    if joutput == "[{}]":
                    joutput = "[]"
    return joutput


def do_create(args, config): 
    pt_id = args.pt_id
    pt_name = args.pt_name
    checksum = args.checksum
    version = args.version
    alias = args.alias
    licensing = args.licensing
    label = args.label
    description = args.description
    private_key = args.private_key
    public_key = args.public_key


    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"]=[{"role":"admin"},{"role":"member"}]
    payload = json.dumps(key)
       
    headers = {'content-type': 'application/json'}

    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth",data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get('status')and statusinfo.get('message'):
            
        status = statusinfo['status']
        message = statusinfo['message']
            
        if status == 'success' and message == 'authorized':
            b_url = config.get('DEFAULT', 'url')
            client = PartBatch(base_url=b_url)
            response = client.create(
            pt_id,pt_name,checksum,version,alias,licensing,label,description,private_key,public_key
            )
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
   



def add_Category(args, config):
    pt_id = args.pt_id
    category_id = args.category_id
    private_key = args.private_key
    public_key = args.public_key

    # #
    # context = create_context('secp256k1')
    # private_key = context.new_random_private_key()
    # public_key = context.get_public_key(private_key)
    # #
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"]=[{"role":"admin"},{"role":"member"}]
    payload = json.dumps(key)
       
    headers = {'content-type': 'application/json'}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth",data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get('status')and statusinfo.get('message'):
            
        status = statusinfo['status']
        message = statusinfo['message']
            
        if status == 'success' and message == 'authorized':
            b_url = config.get('DEFAULT', 'url')
            client = PartBatch(base_url=b_url)
            response = client.add_category(pt_id,category_id,private_key,public_key)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
  

def do_add_artifact(args, config):
    pt_id = args.pt_id
    artifact_id = args.artifact_id
    private_key = args.private_key
    public_key = args.public_key

    # #
    # context = create_context('secp256k1')
    # private_key = context.new_random_private_key()
    # public_key = context.get_public_key(private_key)
    # #
    
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"]=[{"role":"admin"},{"role":"member"}]
    payload = json.dumps(key)
       
    headers = {'content-type': 'application/json'}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth",data=json.dumps(key),headers=headers)
    
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
    if statusinfo.get('status')and statusinfo.get('message'):
            
        status = statusinfo['status']
        message = statusinfo['message']
            
        if status == 'success' and message == 'authorized':
            b_url = config.get('DEFAULT', 'url')
            client = PartBatch(base_url=b_url)
            response = client.add_artifact(pt_id,artifact_id,private_key,public_key)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
   

# add the relationship between parent artifact and supplier
def add_Supplier(args, config):
    pt_id = args.pt_id
    supplier_id = args.supplier_id
    private_key = args.private_key
    public_key = args.public_key

    # #
    # context = create_context('secp256k1')
    # private_key = context.new_random_private_key()
    # public_key = context.get_public_key(private_key)
    # #
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"]=[{"role":"admin"},{"role":"member"}]
    payload = json.dumps(key)
       
    headers = {'content-type': 'application/json'}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth",data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get('status')and statusinfo.get('message'):
            
        status = statusinfo['status']
        message = statusinfo['message']
            
        if status == 'success' and message == 'authorized':
            b_url = config.get('DEFAULT', 'url')
            client = PartBatch(base_url=b_url)
            response = client.add_supplier(pt_id,supplier_id,private_key,public_key)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
  
def print_msg(response):
    
    if "batch_statuses?id" in response:
        print(ret_msg("success","OK","EmptyRecord","{}"))
    else:
        print(ret_msg("failed","Exception raised","EmptyRecord","{}"))
        

def load_config():
   
    config = configparser.ConfigParser()
    config.set('DEFAULT', 'url', 'http://127.0.0.1:8008')
    return config


def main(prog_name=os.path.basename(sys.argv[0]), args=None):
    if args is None:
        args = sys.argv[1:]
    parser = create_parser(prog_name)
    args = parser.parse_args(args)

    if args.verbose is None:
        verbose_level = 0
    else:
        verbose_level = args.verbose

    setup_loggers(verbose_level=verbose_level)

    config = load_config()

    if args.command == 'create':
        do_create(args, config)
    
    elif args.command == 'list-part':
        do_list_part(args, config)
    elif args.command == 'retrieve':
        do_retrieve(args, config)
    elif args.command == 'AddArtifact':
        do_add_artifact(args, config)     
    elif args.command == 'AddSupplier':
        add_Supplier(args, config)     
    elif args.command == 'AddCategory':
        add_Category(args, config)          
    else:
        raise PartException("invalid command: {}".format(args.command))

def removekey(d,key):
    r = dict(d)
    del r[key]
    return r

def ret_msg(status,message,result_type,result):
    msgJSON = "{}"
    key = json.loads(msgJSON)
    key["status"] = status
    key["message"] = message
    key["result_type"] = result_type
    key["result"] = json.loads(result)
    
    msgJSON = json.dumps(key)
    return msgJSON

def main_wrapper():
    try:
        main()
    except PartException as err:
        errmsg = str(err)
        if '404' in errmsg:
            exp = ret_msg("failed","404 Not Found","EmptyRecord","{}")
            print(PartException(exp))
           
        else:
            exp = ret_msg("failed",errmsg,"EmptyRecord","{}")
            print(PartException()) 
        sys.exit(1)
    except KeyboardInterrupt:
        pass
    except SystemExit as err:
        raise err
    except BaseException as err:
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)

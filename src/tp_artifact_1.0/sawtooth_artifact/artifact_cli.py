# Copyright 2017 Intel Corporation
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
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
#

from datetime import datetime
from sawtooth_artifact.artifact_batch import ArtifactBatch
from sawtooth_artifact.exceptions import ArtifactException


DISTRIBUTION_NAME = 'sawtooth-artifact'


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
        'artifact_id',
        type=str,
        help='an identifier for the artifact')
    
    parser.add_argument(
        'alias',
        type=str,
        help='an short identifier for the artifact')
    
    parser.add_argument(
        'artifact_name',
        type=str,
        help='Provide artifact name')
    
    parser.add_argument(
        'artifact_type',
        type=str,
        help='provide artifact type')
    
    parser.add_argument(
        'artifact_checksum',
        type=str,
        help='provide artifact checksum')
        
    parser.add_argument(
        'label',
        type=str,
        help='provide artifact label')
    
    parser.add_argument(
        'openchain',
        type=str,
        help='provide artifact Open Chain status')
    
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


def add_list_artifact_parser(subparsers, parent_parser):
    subparsers.add_parser('list-artifact', parents=[parent_parser])


def add_retrieve_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('retrieve', parents=[parent_parser])

    parser.add_argument(
        'artifact_id',
        type=str,
        help='the identifier for the artifact')


def add_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddArtifact', parents=[parent_parser])
    
    parser.add_argument(
        'artifact_id',
        type=str,
        help='the identifier for the artifact')

    parser.add_argument(
        'sub_artifact_id',
        type=str,
        help='the UUID identifier for sub artifact')
    
    parser.add_argument(
        'path',
        type=str,
        help='path of the artifact')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')
    
    
def add_uri_to_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddURI', parents=[parent_parser])
    
    parser.add_argument(
        'artifact_id',
        type=str,
        help='Provide identifier for the artifact')

    parser.add_argument(
        'version',
        type=str,
        help='Provide version')
    
    parser.add_argument(
        'checksum',
        type=str,
        help='Provide Artifact Checksum')
     
    parser.add_argument(
        'content_type',
        type=str,
        help='Provide type')
      
    parser.add_argument(
        'size',
        type=str,
        help='Provide Artifact size')
    
    parser.add_argument(
        'uri_type',
        type=str,
        help='Provide URI type')
    
    parser.add_argument(
        'location',
        type=str,
        help='Provide link/path for the artifact')
    
    parser.add_argument(
        'private_key',
        type=str,
        help='Provide User Private Key')
    
    parser.add_argument(
        'public_key',
        type=str,
        help='Provide User Public Key')
    
    


def ret_msg(status,message,result_type,result):
    msgJSON = "{}"
    key = json.loads(msgJSON)
    key["status"] = status
    key["message"] = message
    key["result_type"] = result_type
    key["result"] = json.loads(result)
    msgJSON = json.dumps(key)
    return msgJSON

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
    add_list_artifact_parser(subparsers, parent_parser)
    add_retrieve_artifact_parser(subparsers, parent_parser)
    add_artifact_parser(subparsers, parent_parser)
    add_uri_to_artifact_parser(subparsers, parent_parser)
    
    return parser

def do_list_artifact(args, config):
    b_url = config.get('DEFAULT', 'url')

    client = ArtifactBatch(base_url=b_url)

    result = client.list_artifact()
    
    if (result is not None):
        result = refine_output_artifact(str(result))
        result = amend_artifact_fields(result)
        output = ret_msg("success","OK","ListOf:ArtifactRecord",result)
        print(output)
    else:     
        raise ArtifactException("Could not retrieve artifact listing")

def removekey(d,key):
    r = dict(d)
    del r[key]
    return r


def do_retrieve_artifact(args, config):
    artifact_id = args.artifact_id

    b_url = config.get('DEFAULT', 'url')
  
    client = ArtifactBatch(base_url=b_url)

    result = client.retrieve_artifact(artifact_id).decode()

    if result is not None:
        output = ''
        result = filter_output(result)
        output = ret_msg("success","OK","ArtifactRecord",result)
        print(output)
    else:
        raise ArtifactException("Artifact not found {}".format(artifact_id))
    

def do_create(args, config):
    
    artifact_id = args.artifact_id
    alias = args.alias
    artifact_name = args.artifact_name
    artifact_type = args.artifact_type
    artifact_checksum = args.artifact_checksum
    
    label = args.label
    openchain = args.openchain
    private_key = args.private_key
    public_key = args.public_key
    timestampUTC = datetime.utcnow()
    timestamp = str(timestampUTC) 

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
            client = ArtifactBatch(base_url=b_url)
            response = client.create(private_key,public_key,artifact_id,alias,artifact_name,artifact_type,artifact_checksum,label,openchain,timestamp)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
   
 
 
def do_add_uri_to_artifact(args, config):
    artifact_id = args.artifact_id
    version = args.version
    checksum = args.checksum
    content_type = args.content_type
    size = args.size
    uri_type = args.uri_type
    location = args.location 
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
            client = ArtifactBatch(base_url=b_url)
            response = client.add_uri(private_key,public_key,artifact_id,version,checksum,content_type,size,uri_type,location)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)   

def do_add_sub_artifact(args, config):
    artifact_id = args.artifact_id
    sub_artifact_id = args.sub_artifact_id
    path = args.path
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
            client = ArtifactBatch(base_url=b_url)
            response = client.add_artifact(private_key,public_key,artifact_id,sub_artifact_id,path)
            print_msg(response)
        else:
            print(output)
    else:
        print(output)
   
    

def filter_output(result):
    
    mylist = result.split(',',1)
    newstr = mylist[1]
    jsonStr = newstr.replace('artifact_id','uuid').replace('artifact_name','filename').replace('artifact_checksum','checksum').replace('artifact_type','content_type')
    data = json.loads(jsonStr)
    jsonStr = json.dumps(data)
    return jsonStr


def refine_output_artifact(inputstr):
    inputstr = inputstr[1:-1]
    output = inputstr
    output = "["+output+"]"  
    output = re.sub('b\'.*?,{',"{",output)
    output = output.replace("}']","}]")
    output = output.replace("}', {","},{")
                               
    return output

def amend_artifact_fields(inputstr):
        output = inputstr.replace("\\","").replace("artifact_name","filename").replace("artifact_id","uuid").replace("artifact_checksum","checksum").replace("artifact_type","content_type")
        return output

        
def load_config():

    config = configparser.ConfigParser()
    config.set('DEFAULT', 'url', 'http://127.0.0.1:8008')
    return config


def print_msg(response):
    
    if "batch_statuses?id" in response:
        print(ret_msg("success","OK","EmptyRecord","{}"))
    else:
        print(ret_msg("failed","Exception raised","EmptyRecord","{}"))


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
    elif args.command == 'list-artifact':
        do_list_artifact(args, config)
    elif args.command == 'retrieve':
        do_retrieve_artifact(args, config)
    elif args.command == 'AddArtifact':
        do_add_sub_artifact(args, config)  
    elif args.command == 'AddURI':
        do_add_uri_to_artifact(args,config)
    else:    
        raise ArtifactException("invalid command {}".format(args.command))
    

def main_wrapper():
    try:
        main()
    except ArtifactException as err:
        errmsg = str(err)
        if '404' in errmsg:
            exp = ret_msg("failed","404 Not Found","EmptyRecord","{}")
            print(ArtifactException(exp))
           
        else:
            exp = ret_msg("failed",errmsg,"EmptyRecord","{}")
            print(ArtifactException()) 
        sys.exit(1)
    except KeyboardInterrupt:
        pass
    except SystemExit as err:
        raise err
    except BaseException as err:
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)

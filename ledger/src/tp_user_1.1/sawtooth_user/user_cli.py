# Copyright 2017 Intel Corporation
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
import requests
import shlex
import subprocess


from colorlog import ColoredFormatter

# import sawtooth_signing.secp256k1_signer as signing

#
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
#

from sawtooth_user.user_batch import UserBatch
from sawtooth_user.exceptions import UserException


DISTRIBUTION_NAME = 'sawtooth-user'

def add_register_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('register', parents=[parent_parser])

    parser.add_argument(
        'user_public_key',
        type=str,
        help='User Public Key')
    
    parser.add_argument(
        'user_name',
        type=str,
        help='Provide user name')
    
    
    parser.add_argument(
        'email_address',
        type=str,
        help='Provide user email address')
    
    parser.add_argument(
        'authorized',
        type=str,
        help='User Authorization')
    
    parser.add_argument(
        'role',
        type=str,
        help='User Role')
       
    parser.add_argument(
        'ad_private_key',
        type=str,
        help='Provide Admin Private Key')
    
    parser.add_argument(
        'ad_public_key',
        type=str,
        help='Provide Admin Public Key')
    
    
    parser.add_argument(
        '--disable-client-validation',
        action='store_true',
        default=False,
        help='disable client validation')
    
    
def add_register_init_user_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('register_init', parents=[parent_parser])

    parser.add_argument(
        'user_public_key',
        type=str,
        help='User Public Key')
    
    parser.add_argument(
        'user_name',
        type=str,
        help='Provide user name')
    
    
    parser.add_argument(
        'email_address',
        type=str,
        help='Provide user email address')
    
    parser.add_argument(
        'authorized',
        type=str,
        help='User Authorization')
    
    parser.add_argument(
        'role',
        type=str,
        help='User Role')
     
    
    parser.add_argument(
        '--disable-client-validation',
        action='store_true',
        default=False,
        help='disable client validation')
    



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


#remove function
def add_list_user_parser(subparsers, parent_parser):
    subparsers.add_parser('list-user', parents=[parent_parser])
    
    
def add_retrieve_user_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('retrieve', parents=[parent_parser])

    parser.add_argument(
        'user_public_key',
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

    add_register_parser(subparsers, parent_parser)
    add_list_user_parser(subparsers, parent_parser)
    add_retrieve_user_parser(subparsers, parent_parser)
    add_register_init_user_parser(subparsers,parent_parser)
        
    return parser
 

#remove function
def do_list_user(args, config):
    b_url = config.get('DEFAULT', 'url')

    client = UserBatch(base_url=b_url)
    user_list = client.list_user()

    if user_list is not None:
        output = refine_output(str(user_list))
        print (output)
    else:
        raise UserException("Could not retrieve user list.")


def do_retrieve_user(args, config):
    user_public_key = args.user_public_key
    b_url = config.get('DEFAULT', 'url')
    client = UserBatch(base_url=b_url)

    data = client.retreive_user(user_public_key)

    if data is not None:
        output = filter_output(str(data))
        print (output)
    else:
        raise UserException("User not found: {}".format(user_public_key))



def do_register_init(args, config):
    
    # priv_key = signing.generate_privkey()
    # pub_key = signing.generate_pubkey(priv_key)

    user_public_key = args.user_public_key
    user_name = args.user_name
    email_address = args.email_address
    authorized = args.authorized
    role = args.role

    #
    context = create_context('secp256k1')
    priv_key = context.new_random_private_key()
    pub_key = context.get_public_key(priv_key)
    #
    priv_key = priv_key.as_hex()
    pub_key = pub_key.as_hex()
    
    cmd = "user list-user"

    cmd = shlex.split(cmd)
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
    process.wait()
    output = ''
    for line in process.stdout:
        output += line.decode("utf-8").strip()

    if output == "[]" and role == "admin" and len(user_public_key) == 66:
        
        b_url = config.get('DEFAULT', 'url')            
        client = UserBatch(base_url=b_url)
        response = client.register_user(user_public_key,user_name,email_address,authorized,role,priv_key,pub_key)
        print_msg(response)        
    else:
        print(ret_access_denied__msg('Invalid operation.'))
        

def do_register_user(args, config):
    user_public_key = args.user_public_key
    user_name = args.user_name
    email_address = args.email_address
    authorized = args.authorized
    role = args.role
    ad_private_key = args.ad_private_key
    ad_public_key = args.ad_public_key

    # #
    # context = create_context('secp256k1')
    # user_private_key = context.new_random_private_key()
    # user_public_key = context.get_public_key(user_private_key)
    # #
    
    if len(user_public_key) == 66:
        payload = "{}"
        key = json.loads(payload)
        key["publickey"] = ad_public_key
        key["privatekey"] = ad_private_key
        key["allowedrole"]=[{"role":"admin"}]
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
                
                client = UserBatch(base_url=b_url)
                response = client.register_user(
                       user_public_key,user_name,email_address,authorized,role,ad_private_key,ad_public_key
                       )
                print_msg(response)
            else:
                print(output)
            
        else:
            print(output)
    else:
        print(ret_access_denied__msg('Invalid key'))
        
    
def filter_output(result):    
    uslist = result.split(',',1)
    output = uslist[1]
    output = output.replace("\\","")
    output = output[:-1]
    return output

def refine_output(inputstr):
    outputstr=inputstr.replace('b\'','').replace('}\'','}').replace("}]","")
    ulist = outputstr.split("},")
    userlist = []
    for line in ulist:
        record = "{"+line.split(",{",1)[-1]+"}"
        userlist.append(record)
    joutput = str(userlist)
    joutput = joutput.replace("'{","{").replace("}'","}").replace(", { {",", {")
    joutput = amend_fields(joutput)
    if joutput == "[{[]}]":
        joutput = "[]"
    return joutput

def ret_access_denied__msg(message):
    expJson = "{}"
    key = json.loads(expJson)
    key["status"] = "failed"
    key["message"] = message
    expJson = json.dumps(key)
    return expJson 

def amend_fields(inputstr):
    output = inputstr.replace("\\","")
    return output

def load_config():
    config = configparser.ConfigParser()
    config.set('DEFAULT', 'url', 'http://127.0.0.1:8008')
    return config

def print_msg(response):
    if "batch_statuses?id" in response:
        print ("{\"status\":\"success\"}")
    else:
        print ("{\"status\":\"exception\"}")
        
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

    if args.command == 'register':
        do_register_user(args, config)
    elif args.command == 'list-user':
        do_list_user(args, config)
    elif args.command == 'retrieve':
        do_retrieve_user(args, config)
    elif args.command == 'register_init':
        do_register_init(args,config)
   
    else:
        raise UserException("invalid command: {}".format(args.command))


def main_wrapper():
    try:
        main()
    except UserException as err:
        newstr = str(err)
        if '404' in newstr:
            print("{\"status\":\"404 Not Found\"}")
        else:
            error_message = "{\"status\":\"failed\",\"message\":\""
            closing_str = "\"}"
            print (error_message+newstr+closing_str)
            
        sys.exit(1)
    except KeyboardInterrupt:
        pass
    except SystemExit as err:
        raise err
    except BaseException as err:
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)

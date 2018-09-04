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
import json
import os
import traceback
import sys
import shutil
import pkg_resources
import re

from colorlog import ColoredFormatter

import sawtooth_signing.secp256k1_signer as signing

from sawtooth_supplier.supplier_batch import SupplierBatch
from sawtooth_supplier.exceptions import SupplierException


DISTRIBUTION_NAME = 'sawtooth-supplier'


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
        'supplier_id',
        type=str,
        help='an identifier for the supplier')
    
    parser.add_argument(
        'short_id',
        type=str,
        help='an short identifier for the supplier')
    
    parser.add_argument(
        'supplier_name',
        type=str,
        help='Provide supplier name')
    
    parser.add_argument(
        'passwd',
        type=str,
        help='provide hashed password')
    
    parser.add_argument(
        'supplier_url',
        type=str,
        help='provide URL')

    parser.add_argument(
        '--disable-client-validation',
        action='store_true',
        default=False,
        help='disable client validation')


def add_init_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('init', parents=[parent_parser])

    parser.add_argument(
        '--username',
        type=str,
        help='the name of the user')

    parser.add_argument(
        '--url',
        type=str,
        help='the url of the REST API')



def add_list_supplier_parser(subparsers, parent_parser):
    subparsers.add_parser('list-supplier', parents=[parent_parser])


def add_retrieve_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('retrieve', parents=[parent_parser])

    parser.add_argument(
        'supplier_id',
        type=str,
        help='an identifier for the supplier')
    
    

def add_part_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('AddPart', parents=[parent_parser])
    
    parser.add_argument(
        'supplier_id',
        type=str,
        help='the identifier for the supplier')

    parser.add_argument(
        'part_id',
        type=str,
        help='the UUID identifier for Part')
    



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

    parent_parser.add_argument(
        '--auth-user',
        type=str,
        help='username for authentication if REST API is using Basic Auth')

    parent_parser.add_argument(
        '--auth-password',
        type=str,
        help='password for authentication if REST API is using Basic Auth')

    return parent_parser


def create_parser(prog_name):
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        parents=[parent_parser],
        formatter_class=argparse.RawDescriptionHelpFormatter)

    subparsers = parser.add_subparsers(title='subcommands', dest='command')

    add_create_parser(subparsers, parent_parser)
    add_init_parser(subparsers, parent_parser)   
    add_list_supplier_parser(subparsers, parent_parser)
    add_retrieve_parser(subparsers, parent_parser)
    add_part_parser(subparsers, parent_parser)

    return parser


def do_init(args, config):
    username = config.get('DEFAULT', 'username')
    if args.username is not None:
        username = args.username

    url = config.get('DEFAULT', 'url')
    if args.url is not None:
        url = args.url

    config.set('DEFAULT', 'username', username)
    print("set username: {}".format(username))
    config.set('DEFAULT', 'url', url)
    print("set url: {}".format(url))

    save_config(config)

    priv_filename = config.get('DEFAULT', 'key_file')
    if priv_filename.endswith(".priv"):
        addr_filename = priv_filename[0:-len(".priv")] + ".addr"
    else:
        addr_filename = priv_filename + ".addr"

    if not os.path.exists(priv_filename):
        try:
            if not os.path.exists(os.path.dirname(priv_filename)):
                os.makedirs(os.path.dirname(priv_filename))

            privkey = signing.generate_privkey()
            pubkey = signing.generate_pubkey(privkey)
            addr = signing.generate_identifier(pubkey)

            with open(priv_filename, "w") as priv_fd:
                print("writing file: {}".format(priv_filename))
                priv_fd.write(privkey)
                priv_fd.write("\n")

            with open(addr_filename, "w") as addr_fd:
                print("writing file: {}".format(addr_filename))
                addr_fd.write(addr)
                addr_fd.write("\n")
        except IOError as ioe:
            raise SupplierException("IOError: {}".format(str(ioe)))




def do_list_supplier(args, config):
    url = config.get('DEFAULT', 'url')
    key_file = config.get('DEFAULT', 'key_file')
    auth_user, auth_password = _get_auth_info(args)

    client = SupplierBatch(base_url=url, keyfile=key_file)

    result = client.list_supplier(auth_user=auth_user,
                                 auth_password=auth_password)

    if result is not None:
        out = refine_output_supplier(str(result))
        output = refine_output(out)
        print(output)
    else:
        raise SupplierException("Could not retrieve supplier listing.")

		
def refine_output_supplier(inputstr):
    inputstr = inputstr[1:-1]
    output = re.sub(r'\[.*?\]', '',inputstr)
    output = "["+output+"]"  
    return output

def amend_supplier_fields(inputstr):
        output = inputstr.replace("\\","").replace('supplier_id','uuid').replace('supplier_name','name').replace('supplier_url','url')
        return output

        
def refine_output(inputstr):
    
                subpartstr = "\"parts\": ,"
                outputstr=inputstr.replace(subpartstr,"").replace('b\'','').replace('}\'','}').replace(", \"parts\": ","")
                outputstr=outputstr.replace('b\'','').replace('}\'','}')
                slist = outputstr.split("},")
                supplierlist = []
                for line in slist:
                        record = "{"+line.split(",{",1)[-1]+"}"
                        supplierlist.append(record)
                joutput = str(supplierlist)
                joutput = joutput.replace("'{","{").replace("}'","}").replace(", { {",", {").replace("}]}]","}]")
                joutput = amend_supplier_fields(joutput)
                return joutput

def do_retrieve(args, config):
    supplier_id = args.supplier_id

    url = config.get('DEFAULT', 'url')
    key_file = config.get('DEFAULT', 'key_file')
    auth_user, auth_password = _get_auth_info(args)

    client = SupplierBatch(base_url=url, keyfile=key_file)

    result = client.retrieve_supplier(supplier_id, auth_user=auth_user, auth_password=auth_password).decode()

    if result is not None:
        result = filter_output(result)
        print(result)

    else:
        raise SupplierException("Supplier not found: {}".format(supplier_id))

def removekey(d,key):
    r = dict(d)
    del r[key]
    return r

def print_msg(response):
    if "batch_status?id" in response:
        print ("{\"status\":\"success\"}")
    else:
        print ("{\"status\":\"exception\"}")

def filter_output(result):
    
    supplierlist = result.split(',',1)
    suppstr = supplierlist[1]
    jsonStr = suppstr.replace('supplier_id','uuid').replace('supplier_name','name').replace('supplier_url','url')
    data = json.loads(jsonStr)
    jsonStr = json.dumps(data)
    return jsonStr


def do_create(args, config):
    supplier_id = args.supplier_id
    short_id = args.short_id
    supplier_name = args.supplier_name
    passwd = args.passwd
    supplier_url = args.supplier_url

    url = config.get('DEFAULT', 'url')
    key_file = config.get('DEFAULT', 'key_file')
    auth_user, auth_password = _get_auth_info(args)

    client = SupplierBatch(base_url=url, keyfile=key_file)

   
    response = client.create(
            supplier_id,short_id,supplier_name,passwd,supplier_url,
            auth_user=auth_user,
            auth_password=auth_password)
    
    print_msg(response)



def load_config():
    home = os.path.expanduser("~")
    real_user = getpass.getuser()

    config_file = os.path.join(home, ".sawtooth", "supplier.cfg")
    key_dir = os.path.join(home, ".sawtooth", "keys")

    config = configparser.ConfigParser()
    config.set('DEFAULT', 'url', 'http://127.0.0.1:8080')
    config.set('DEFAULT', 'key_dir', key_dir)
    config.set('DEFAULT', 'key_file', '%(key_dir)s/%(username)s.priv')
    config.set('DEFAULT', 'username', real_user)
    if os.path.exists(config_file):
        config.read(config_file)

    return config


def save_config(config):
    home = os.path.expanduser("~")

    config_file = os.path.join(home, ".sawtooth", "supplier.cfg")
    if not os.path.exists(os.path.dirname(config_file)):
        os.makedirs(os.path.dirname(config_file))

    with open("{}.new".format(config_file), "w") as fd:
        config.write(fd)
    if os.name == 'nt':
        if os.path.exists(config_file):
            os.remove(config_file)
    os.rename("{}.new".format(config_file), config_file)


def _get_auth_info(args):
    auth_user = args.auth_user
    auth_password = args.auth_password
    if auth_user is not None and auth_password is None:
        auth_password = getpass.getpass(prompt="Auth Password: ")

    return auth_user, auth_password


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
    elif args.command == 'init':
        do_init(args, config)
    elif args.command == 'list-supplier':
        do_list_supplier(args, config)
    elif args.command == 'retrieve':
        do_retrieve(args, config)
    elif args.command == 'AddPart':
        do_addpart(args, config) 

    else:
        raise SupplierException("invalid command: {}".format(args.command))

def do_addpart(args, config):
    supplier_id = args.supplier_id
    part_id = args.part_id
   
    url = config.get('DEFAULT', 'url')
    key_file = config.get('DEFAULT', 'key_file')

    client = SupplierBatch(base_url=url,
                      keyfile=key_file)
    response = client.add_part(supplier_id,part_id)
    print_msg(response)


def main_wrapper():
    try:
        main()
    except SupplierException as err:
        newstr = str(err)
        if '404' in newstr:
            print("{\"status\":\"404 Not Found\"}")
        else:
            error_message = "{\"error\":\"failed\",\"error_message\":\""
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

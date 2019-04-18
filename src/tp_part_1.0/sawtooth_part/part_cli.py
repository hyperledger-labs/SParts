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
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
from __future__ import print_function

import argparse
import configparser
import logging
import os
import traceback
import sys
import pkg_resources
import json
import requests
from colorlog import ColoredFormatter
from sawtooth_part.part_batch import PartBatch
from sawtooth_part.exceptions import PartException

DISTRIBUTION_NAME = "sawtooth-part"
################################################################################
def create_console_handler(verbose_level):
    """
    """
    
    clog = logging.StreamHandler()
    formatter = ColoredFormatter(
        "%(log_color)s[%(asctime)s %(levelname)-8s%(module)s]%(reset)s "
        "%(white)s%(message)s",
        datefmt="%H:%M:%S",
        reset=True,
        log_colors={
            "DEBUG": "cyan",
            "INFO": "green",
            "WARNING": "yellow",
            "ERROR": "red",
            "CRITICAL": "red",
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
    """
    """
    
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    logger.addHandler(create_console_handler(verbose_level))
################################################################################
#                                   OBJ                                        #
################################################################################
def add_create_parser(subparsers, parent_parser):
    """
    """
    
    parser = subparsers.add_parser("create", parents=[parent_parser])

    parser.add_argument(
        "pt_id",
        type=str,
        help="an identifier for the part")
    
    parser.add_argument(
        "pt_name",
        type=str,
        help="provide part name")
    
    parser.add_argument(
        "checksum",
        type=str,
        help="Provide checksum")
    
    parser.add_argument(
        "version",
        type=str,
        help="provide version for the part")
    
    parser.add_argument(
        "alias",
        type=str,
        help="provide alias")
    
    parser.add_argument(
        "licensing",
        type=str,
        help="provide licensing")
    
    parser.add_argument(
        "label",
        type=str,
        help="provide label")
    
    parser.add_argument(
        "description",
        type=str,
        help="provide description")
    
    parser.add_argument(
        "private_key",
        type=str,
        help="Provide User Private Key")
    
    parser.add_argument(
        "public_key",
        type=str,
        help="Provide User Public Key")


    parser.add_argument(
        "--disable-client-validation",
        action="store_true",
        default=False,
        help="disable client validation")

def add_list_part_parser(subparsers, parent_parser):
    """
    """
    
    subparsers.add_parser("list-part", parents=[parent_parser])

def add_retrieve_parser(subparsers, parent_parser):
    """
    """
    parser = subparsers.add_parser("retrieve", parents=[parent_parser])

    parser.add_argument(
        "pt_id",
        type=str,
        help="part identifier")
    
    parser.add_argument(
        "-a", "--all",
        action="store_true",
        default=False,
        help="show history of uuid")
        
    parser.add_argument(
        "--range",
        nargs=2,
        metavar=("START", "END"),
        default=None,
        help="show history of uuid within the range; FORMAT : yyyymmdd")

def add_amend_parser(subparsers, parent_parser):
    """
    """
    
    parser = subparsers.add_parser("amend", parents=[parent_parser])

    parser.add_argument(
        "pt_id",
        type=str,
        help="an identifier for the part")
    
    parser.add_argument(
        "pt_name",
        type=str,
        help="provide part name")
    
    parser.add_argument(
        "checksum",
        type=str,
        help="Provide checksum")
    
    parser.add_argument(
        "version",
        type=str,
        help="provide version for the part")
    
    parser.add_argument(
        "alias",
        type=str,
        help="provide alias")
    
    parser.add_argument(
        "licensing",
        type=str,
        help="provide licensing")
    
    parser.add_argument(
        "label",
        type=str,
        help="provide label")
    
    parser.add_argument(
        "description",
        type=str,
        help="provide description")
    
    parser.add_argument(
        "private_key",
        type=str,
        help="Provide User Private Key")
    
    parser.add_argument(
        "public_key",
        type=str,
        help="Provide User Public Key")


    parser.add_argument(
        "--disable-client-validation",
        action="store_true",
        default=False,
        help="disable client validation")
   
def add_artifact_parser(subparsers, parent_parser):
    """
    """
    
    parser = subparsers.add_parser("AddArtifact", parents=[parent_parser])
    
    parser.add_argument(
        "pt_id",
        type=str,
        help="part identifier")

    parser.add_argument(
        "artifact_id",
        type=str,
        help="the UUID identifier for artifact")
    
    parser.add_argument(
        "private_key",
        type=str,
        help="Provide User Private Key")
    
    parser.add_argument(
        "public_key",
        type=str,
        help="Provide User Public Key")
        
    parser.add_argument(
        "-D", "--delete",
        action="store_true",
        default=False,
        help="removes the artifact")

def add_category_parser(subparsers, parent_parser):
    """
    """
    
    parser = subparsers.add_parser("AddCategory", parents=[parent_parser])
    
    parser.add_argument(
        "pt_id",
        type=str,
        help="the identifier for the part")

    parser.add_argument(
        "category_id",
        type=str,
        help="the identifier for Category")
    
    parser.add_argument(
        "private_key",
        type=str,
        help="Provide User Private Key")
    
    parser.add_argument(
        "public_key",
        type=str,
        help="Provide User Public Key")
        
    parser.add_argument(
        "-D", "--delete",
        action="store_true",
        default=False,
        help="removes the category")

# Provide the UUID of the parent artifact and the UUID of the organization
def add_organization_parser(subparsers, parent_parser):
    """
    """
    
    parser = subparsers.add_parser("AddOrganization", parents=[parent_parser])
    
    parser.add_argument(
        "pt_id",
        type=str,
        help="the identifier for the part")

    parser.add_argument(
        "organization_id",
        type=str,
        help="the identifier for Organization")
    
    parser.add_argument(
        "private_key",
        type=str,
        help="Provide User Private Key")
    
    parser.add_argument(
        "public_key",
        type=str,
        help="Provide User Public Key")
        
    parser.add_argument(
        "-D", "--delete",
        action="store_true",
        default=False,
        help="removes the organization")
################################################################################
#                                   CREATE                                     #
################################################################################
def create_parent_parser(prog_name):
    """
    """
    
    parent_parser = argparse.ArgumentParser(prog=prog_name, add_help=False)
    parent_parser.add_argument(
        "-v", "--verbose",
        action="count",
        help="enable more verbose output")

    try:
        version = pkg_resources.get_distribution(DISTRIBUTION_NAME).version
    except pkg_resources.DistributionNotFound:
        version = "UNKNOWN"

    parent_parser.add_argument(
        "-V", "--version",
        action="version",
        version=(DISTRIBUTION_NAME + " (Hyperledger Sawtooth) version {}")
        .format(version),
        help="print version information")

    return parent_parser

def create_parser(prog_name):
    """
    """
    
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        parents=[parent_parser],
        formatter_class=argparse.RawDescriptionHelpFormatter)

    subparsers = parser.add_subparsers(title="subcommands", dest="command")

    add_create_parser(subparsers, parent_parser)
   
    add_list_part_parser(subparsers, parent_parser)
    add_retrieve_parser(subparsers, parent_parser)
    add_amend_parser(subparsers, parent_parser)
    
    add_artifact_parser(subparsers, parent_parser)
    add_organization_parser(subparsers,parent_parser)
    add_category_parser(subparsers,parent_parser)
    
    return parser
################################################################################
#                               FUNCTIONS                                      #
################################################################################    
def do_list_part(args, config):
    """
    """
    
    b_url = config.get("DEFAULT", "url")
    client = PartBatch(base_url=b_url)
    result = client.list_part()

    if result is not None:
        result.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(result)
        
        output = ret_msg("success", "OK", "ListOf:PartRecord", result)
        
        print(output)
    else:
        raise PartException("Could not retrieve part listing.")

def do_retrieve(args, config):
    """
    """
    
    all_flag = args.all
    range_flag = args.range
    
    pt_id = args.pt_id
    
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = PartBatch(base_url=b_url)
    data = client.retrieve_part(pt_id, all_flag, range_flag)
    
    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "PartRecord", data.decode())
        else:
            output = ret_msg("success", "OK", "PartRecord", data)
            
        print(output)
    else:
        raise PartException("Part not found: {}".format(pt_id))

def do_create_part(args, config):
    """
    """
    
    pt_id       = args.pt_id
    pt_name     = args.pt_name
    checksum    = args.checksum
    version     = args.version
    alias       = args.alias
    licensing   = args.licensing
    label       = args.label
    description = args.description
    private_key = args.private_key
    public_key  = args.public_key

    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}

    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.create_part(pt_id, pt_name, checksum, version, alias, 
                            licensing, label, description, private_key, 
                            public_key)
            
            print_msg(response, "create")
        else:
            print(output)
    else:
        print(output)

def do_amend_part(args, config):
    """
    """
    
    pt_id       = args.pt_id
    pt_name     = args.pt_name
    checksum    = args.checksum
    version     = args.version
    alias       = args.alias
    licensing   = args.licensing
    label       = args.label
    description = args.description
    private_key = args.private_key
    public_key  = args.public_key

    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}

    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.amend_part(pt_id, pt_name, checksum, version, alias, 
                            licensing, label, description, private_key, 
                            public_key)
            
            print_msg(response, "amend")
        else:
            print(output)
    else:
        print(output)
   
def do_add_artifact(args, config):
    """
    """
    
    deleteArt   = args.delete
    
    pt_id       = args.pt_id
    artifact_id = args.artifact_id
    private_key = args.private_key
    public_key  = args.public_key
    
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key), headers=headers)
    
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_artifact(
                                pt_id, artifact_id, private_key, public_key,
                                deleteArt
                            )
                            
            print_msg(response, "AddArtifact")
        else:
            print(output)
    else:
        print(output)
        
def do_add_category(args, config):
    """
    """
    
    deleteCat   = args.delete
    
    pt_id       = args.pt_id
    category_id = args.category_id
    private_key = args.private_key
    public_key  = args.public_key
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_category(
                                pt_id, category_id, private_key, public_key,
                                deleteCat
                            )
                            
            print_msg(response, "AddCategory")
        else:
            print(output)
    else:
        print(output)      

# add the relationship between parent artifact and organization
def do_add_organization(args, config):
    """
    """
    
    del_flag   = args.delete
    
    pt_id           = args.pt_id
    organization_id = args.organization_id
    private_key     = args.private_key
    public_key      = args.public_key
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_organization(
                                pt_id, organization_id, private_key, public_key,
                                del_flag
                            )
                            
            print_msg(response, "AddOrganization")
        else:
            print(output)
    else:
        print(output)
################################################################################
#                                  PRINT                                       #
################################################################################ 
def load_config():
    """
    """
    
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    return config

def print_msg(response, cmd=None):
    """
    """
    
    try:
        if type(response) is list and response[0] == None:
            if len(response) > 1:
                raise PartException(
                        "PartException : {}".format(response[1])
                    )
            raise PartException(
                        "PartException : No change."
                    )
        
        if response == None:
            if cmd == "create":
                raise PartException("PartException : Duplicate UUID.")
                
            elif (cmd == "amend" or cmd == "AddOrganization" or 
                    cmd == "AddArtifact" or cmd == "AddCategory"):
                raise PartException(
                            "PartException : UUID does not exist."
                        )
                
            raise PartException("Exception raised.")
        elif "batch_statuses?id" in response:
            print(ret_msg("success", "OK", "PartRecord", "{}"))
            return ret_msg("success", "OK", "PartRecord", "{}")
        else:
            raise PartException("Exception raised.")
    except BaseException as err:
        output = ret_msg(
                            "failed",
                            str(err),
                            "PartRecord", "{}"
                        )
        print(output)
        return output
        
def ret_msg(status, message, result_type, result):
    """
    """
    
    msgJSON = "{}"
    key = json.loads(msgJSON)
    key["status"] = status
    key["message"] = message
    key["result_type"] = result_type
    key["result"] = result if type(result) is list else json.loads(result)
   
    msgJSON = json.dumps(key)
    return msgJSON
################################################################################
#                                   MAIN                                       #
################################################################################
def main(prog_name=os.path.basename(sys.argv[0]), args=None):
    """
    """
    
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

    if args.command == "create":
        do_create_part(args, config)
    elif args.command == "list-part":
        do_list_part(args, config)
    elif args.command == "retrieve":
        do_retrieve(args, config)
    elif args.command == "amend":
        do_amend_part(args, config)
    elif args.command == "AddArtifact":
        do_add_artifact(args, config)     
    elif args.command == "AddOrganization":
        do_add_organization(args, config)     
    elif args.command == "AddCategory":
        do_add_category(args, config)          
    else:
        raise PartException("invalid command: {}".format(args.command))

def main_wrapper():
    """
    """
    
    try:
        main()
    except PartException as err:
        errmsg = str(err)
        if "404" in errmsg:
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
################################################################################
#                                 API                                          #
################################################################################
def api_do_create_part(args, config):
    """
    """
    
    param_check = _payload_check_(args, creation=True)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    pt_id       = args["part"]["uuid"]
    pt_name     = args["part"]["name"]
    checksum    = args["part"]["checksum"]
    version     = args["part"]["version"]
    alias       = args["part"]["alias"]
    licensing   = args["part"]["licensing"]
    label       = args["part"]["label"]
    description = args["part"]["description"]
    private_key = args["private_key"]
    public_key  = args["public_key"]

    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}

    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.create_part(pt_id, pt_name, checksum, version,
                            alias, licensing, label, description, private_key, 
                            public_key)
            
            return print_msg(response, "create")
        else:
            return output
    else:
        return output
        
def api_do_amend_part(args, config):
    """
    """
    
    param_check = _payload_check_(args)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    pt_id       = args["part"]["uuid"]
    
    pt_name     = _null_cast(args["part"], "name")
    checksum    = _null_cast(args["part"], "checksum")
    version     = _null_cast(args["part"], "version")
    alias       = _null_cast(args["part"], "alias")
    licensing   = _null_cast(args["part"], "licensing")
    label       = _null_cast(args["part"], "label")
    description = _null_cast(args["part"], "description")
    
    private_key = args["private_key"]
    public_key  = args["public_key"]

    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}

    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.amend_part(pt_id, pt_name, checksum, version,
                            alias, licensing, label, description, private_key, 
                            public_key)
            
            return print_msg(response, "amend")
        else:
            return output
    else:
        return output
        
def api_do_list_part(config):
    """
    """
    
    b_url = config.get("DEFAULT", "url")
    client = PartBatch(base_url=b_url)
    result = client.list_part()

    if result is not None:
        result.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(result)
        
        output = ret_msg("success", "OK", "ListOf:PartRecord", result)
        
        return output
    else:
        return ret_msg(
                    "failed", 
                    "PartException : Could not retrieve part listing.", 
                    "PartRecord", "{}"
                )

def api_do_retrieve_part(pt_id, config, all_flag=False, range_flag=None):
    """
    """
    
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = PartBatch(base_url=b_url)
    data = client.retrieve_part(pt_id, all_flag, range_flag)
    
    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "PartRecord", data.decode())
        else:
            output = ret_msg("success", "OK", "PartRecord", data)
            
        return output
    else:
        return ret_msg(
                    "failed",
                    "PartException : UUID {} does not exist." \
                    .format(pt_id),
                    "PartRecord", "{}"
                )
                
def api_do_add_organization(args, config, del_flag=False):
    """
    """
    
    param_check = _payload_check_(args, cmd="AddOrganization")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    pt_id           = args["relation"]["part_uuid"]
    organization_id = args["relation"]["organization_uuid"]
    private_key     = args["private_key"]
    public_key      = args["public_key"]
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key), headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_organization(
                                pt_id, organization_id, private_key, public_key,
                                del_flag
                            )
            
            return print_msg(response, "AddOrganization")
        else:
            return output
    else:
        return output
        
def api_do_add_category(args, config, del_flag=False):
    """
    """
    
    param_check = _payload_check_(args, cmd="AddCategory")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    pt_id       = args["relation"]["part_uuid"]
    category_id = args["relation"]["category_uuid"]
    private_key = args["private_key"]
    public_key  = args["public_key"]
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_category(
                                pt_id, category_id, private_key, public_key,
                                del_flag
                            )
                            
            return print_msg(response, "AddCategory")
        else:
            return output
    else:
        return output
        
def api_do_add_artifact(args, config, del_flag=False):
    """
    """
    
    param_check = _payload_check_(args, cmd="AddArtifact")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    pt_id       = args["relation"]["part_uuid"]
    artifact_id = args["relation"]["artifact_uuid"]
    private_key = args["private_key"]
    public_key  = args["public_key"]
   
    payload = "{}"
    key = json.loads(payload)
    key["publickey"] = public_key
    key["privatekey"] = private_key
    key["allowedrole"] = [{"role" : "admin"}, {"role" : "member"}]
    payload = json.dumps(key)
       
    headers = {"content-type" : "application/json"}
    response = requests.post("http://127.0.0.1:818/api/sparts/ledger/auth", 
                    data=json.dumps(key), headers=headers)
    
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = PartBatch(base_url=b_url)
            response = client.add_artifact(
                                pt_id, artifact_id, private_key, public_key,
                                del_flag
                            )
                            
            return print_msg(response, "AddArtifact")
        else:
            return output
    else:
        return output
################################################################################
#                           API PRIVATE FUNCTIONS                              #
################################################################################
def _payload_check_(args, creation=False, cmd=None):
    """
    """
    
    if cmd != None:
        if cmd == "AddOrganization":
            if "relation" not in args:
                return [True, "Relation missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "part_uuid" not in args["relation"]:
                return [True, "Part UUID missing."]
            elif "organization_uuid" not in args["relation"]:
                return [True, "Organization UUID missing."]
            else:
                return [False]
        elif cmd == "AddCategory":
            if "relation" not in args:
                return [True, "Relation missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "part_uuid" not in args["relation"]:
                return [True, "Part UUID missing."]
            elif "category_uuid" not in args["relation"]:
                return [True, "Category UUID missing."]
            else:
                return [False]
        elif cmd == "AddArtifact":
            if "relation" not in args:
                return [True, "Relation missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "part_uuid" not in args["relation"]:
                return [True, "Part UUID missing."]
            elif "artifact_uuid" not in args["relation"]:
                return [True, "Artifact UUID missing."]
            else:
                return [False]
        else:
                return [False]
    else:  
        if creation:
            if "part" not in args:
                return [True, "Part missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["part"]:
                return [True, "UUID missing."]
            elif "name" not in args["part"]:
                return [True, "Name missing."]
            elif "checksum" not in args["part"]:
                return [True, "Checksum missing."]
            elif "version" not in args["part"]:
                return [True, "Version missing."]
            elif "alias" not in args["part"]:
                return [True, "Alias missing."]
            elif "licensing" not in args["part"]:
                return [True, "Licensing missing."]
            elif "label" not in args["part"]:
                return [True, "Label missing."]    
            elif "description" not in args["part"]:
                return [True, "Description missing."]
            else:
                return [False]
        else:
            if "part" not in args:
                return [True, "Part missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["part"]:
                return [True, "UUID missing."]
            else:
                return [False]

def _null_cast(dic, key):
    """
    """
    
    if key not in dic:
        return "null"
    return dic[key]
################################################################################
#                                                                              #
################################################################################

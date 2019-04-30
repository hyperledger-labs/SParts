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
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
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

from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey

from sawtooth_artifact.artifact_batch import ArtifactBatch
from sawtooth_artifact.exceptions import ArtifactException


DISTRIBUTION_NAME = "sawtooth-artifact"
################################################################################
def create_console_handler(verbose_level):
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
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    logger.addHandler(create_console_handler(verbose_level))
################################################################################
#                                   OBJ                                        #
################################################################################
def add_create_parser(subparsers, parent_parser):
    parser = subparsers.add_parser("create", parents=[parent_parser])

    parser.add_argument(
        "artifact_id",
        type=str,
        help="an identifier for the artifact")
    
    parser.add_argument(
        "alias",
        type=str,
        help="an short identifier for the artifact")
    
    parser.add_argument(
        "artifact_name",
        type=str,
        help="Provide artifact name")
    
    parser.add_argument(
        "artifact_type",
        type=str,
        help="provide artifact type")
    
    parser.add_argument(
        "artifact_checksum",
        type=str,
        help="provide artifact checksum")
        
    parser.add_argument(
        "label",
        type=str,
        help="provide artifact label")
    
    parser.add_argument(
        "openchain",
        type=str,
        help="provide artifact Open Chain status")
    
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

def add_list_artifact_parser(subparsers, parent_parser):
    subparsers.add_parser("list-artifact", parents=[parent_parser])

def add_retrieve_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser("retrieve", parents=[parent_parser])

    parser.add_argument(
        "artifact_id",
        type=str,
        help="the identifier for the artifact")
    
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
    parser = subparsers.add_parser("amend", parents=[parent_parser])

    parser.add_argument(
        "artifact_id",
        type=str,
        help="an identifier for the artifact")
    
    parser.add_argument(
        "alias",
        type=str,
        help="an short identifier for the artifact")
    
    parser.add_argument(
        "artifact_name",
        type=str,
        help="Provide artifact name")
    
    parser.add_argument(
        "artifact_type",
        type=str,
        help="provide artifact type")
    
    parser.add_argument(
        "artifact_checksum",
        type=str,
        help="provide artifact checksum")
        
    parser.add_argument(
        "label",
        type=str,
        help="provide artifact label")
    
    parser.add_argument(
        "openchain",
        type=str,
        help="provide artifact Open Chain status")
    
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
    parser = subparsers.add_parser("AddArtifact", parents=[parent_parser])
    
    parser.add_argument(
        "artifact_id",
        type=str,
        help="the identifier for the artifact")

    parser.add_argument(
        "sub_artifact_id",
        type=str,
        help="the UUID identifier for sub artifact")
    
    parser.add_argument(
        "path",
        type=str,
        help="path of the artifact")
    
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
        help="removes the sub artifact")
    
def add_uri_to_artifact_parser(subparsers, parent_parser):
    parser = subparsers.add_parser("AddURI", parents=[parent_parser])
    
    parser.add_argument(
        "artifact_id",
        type=str,
        help="Provide identifier for the artifact")

    parser.add_argument(
        "version",
        type=str,
        help="Provide version")
    
    parser.add_argument(
        "checksum",
        type=str,
        help="Provide Artifact Checksum")
     
    parser.add_argument(
        "content_type",
        type=str,
        help="Provide type")
      
    parser.add_argument(
        "size",
        type=str,
        help="Provide Artifact size")
    
    parser.add_argument(
        "uri_type",
        type=str,
        help="Provide URI type")
    
    parser.add_argument(
        "location",
        type=str,
        help="Provide link/path for the artifact")
    
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
        help="removes the URI")
################################################################################
#                                   CREATE                                     #
################################################################################    
def create_parent_parser(prog_name):
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
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        parents=[parent_parser],
        formatter_class=argparse.RawDescriptionHelpFormatter)

    subparsers = parser.add_subparsers(title="subcommands", dest="command")

    add_create_parser(subparsers, parent_parser)
    add_list_artifact_parser(subparsers, parent_parser)
    add_retrieve_artifact_parser(subparsers, parent_parser)
    add_amend_parser(subparsers, parent_parser)
    
    add_artifact_parser(subparsers, parent_parser)
    add_uri_to_artifact_parser(subparsers, parent_parser)
    
    return parser
################################################################################
#                               FUNCTIONS                                      #
################################################################################
def do_list_artifact(args, config):
    b_url = config.get("DEFAULT", "url")
    client = ArtifactBatch(base_url=b_url)
    result = client.list_artifact()
    
    if result is not None:
        result.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(result)
        
        output = ret_msg("success", "OK", "ListOf:ArtifactRecord", result)
        
        print(output)
    else:     
        raise ArtifactException("Could not retrieve artifact listing.")

def do_retrieve_artifact(args, config):
    all_flag    = args.all
    range_flag  = args.range
    
    artifact_id = args.artifact_id
    
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = ArtifactBatch(base_url=b_url)
    data = client.retrieve_artifact(artifact_id, all_flag, range_flag)

    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "ArtifactRecord", data.decode())
        else:
            output = ret_msg("success", "OK", "ArtifactRecord", data)
        
        print(output)
    else:
        raise ArtifactException("Artifact not found {}".format(artifact_id))

def do_create_artifact(args, config):
    artifact_id         = args.artifact_id
    artifact_alias      = args.alias
    artifact_name       = args.artifact_name
    artifact_type       = args.artifact_type
    artifact_checksum   = args.artifact_checksum
    artifact_label      = args.label
    artifact_openchain  = args.openchain
    private_key         = args.private_key
    public_key          = args.public_key 

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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.create_artifact(
                            private_key, public_key, artifact_id,
                            artifact_alias, artifact_name, artifact_type, 
                            artifact_checksum, artifact_label,
                            artifact_openchain
                        )
                        
            print_msg(response, "create")
        else:
            print(output)
    else:
        print(output)

def do_amend_artifact(args, config):
    artifact_id         = args.artifact_id
    artifact_alias      = args.alias
    artifact_name       = args.artifact_name
    artifact_type       = args.artifact_type
    artifact_checksum   = args.artifact_checksum
    artifact_label      = args.label
    artifact_openchain  = args.openchain
    private_key         = args.private_key
    public_key          = args.public_key 

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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.amend_artifact(
                            private_key, public_key, artifact_id, 
                            artifact_alias, artifact_name, artifact_type, 
                            artifact_checksum, artifact_label, 
                            artifact_openchain
                        )
                        
            print_msg(response, "amend")
        else:
            print(output)
    else:
        print(output)   

def do_add_sub_artifact(args, config):
    deleteSub       = args.delete
    
    artifact_id     = args.artifact_id
    sub_artifact_id = args.sub_artifact_id
    path            = args.path
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
                    data=json.dumps(key), headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.add_artifact(
                            private_key, public_key, artifact_id,
                            sub_artifact_id, path, deleteSub
                        )
                            
            print_msg(response, "AddArtifact")
        else:
            print(output)
    else:
        print(output)
    
def do_add_uri_to_artifact(args, config):
    deleteURI       = args.delete
    
    artifact_id     = args.artifact_id
    version         = args.version
    checksum        = args.checksum
    content_type    = args.content_type
    size            = args.size
    uri_type        = args.uri_type
    location        = args.location 
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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.add_uri(private_key, public_key, artifact_id, 
                                version, checksum, content_type, size, uri_type,
                                location, deleteURI)
                                
            print_msg(response, "AddURI")
        else:
            print(output)
    else:
        print(output)   
################################################################################
#                                  PRINT                                       #
################################################################################
def load_config():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    return config

def print_msg(response, cmd=None):
    try:
        if type(response) is list and response[0] == None:
            if len(response) > 1:
                raise ArtifactException(
                        "ArtifactException : {}".format(response[1])
                    )
            raise ArtifactException(
                        "ArtifactException : No change."
                    )
        
        if response == None:
            if cmd == "create":
                raise ArtifactException("ArtifactException : Duplicate UUID.")
                
            elif cmd == "amend" or cmd == "AddArtifact" or cmd == "AddURI":
                raise ArtifactException(
                            "ArtifactException : UUID does not exist."
                        )
                
            raise ArtifactException("Exception raised.")
        elif "batch_statuses?id" in response:
            print(ret_msg("success", "OK", "ArtifactRecord", "{}"))
            return ret_msg("success", "OK", "ArtifactRecord", "{}")
        else:
            raise ArtifactException("Exception raised.")
    except BaseException as err:
        output = ret_msg(
                            "failed",
                            str(err),
                            "ArtifactRecord", "{}"
                        )
        print(output)
        return output

def ret_msg(status, message, result_type, result):
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
        do_create_artifact(args, config) 
    elif args.command == "list-artifact":
        do_list_artifact(args, config)
    elif args.command == "retrieve":
        do_retrieve_artifact(args, config)
    elif args.command == "amend":
        do_amend_artifact(args, config)    
    elif args.command == "AddArtifact":
        do_add_sub_artifact(args, config)  
    elif args.command == "AddURI":
        do_add_uri_to_artifact(args,config)
    else:    
        raise ArtifactException("invalid command {}".format(args.command))
    
def main_wrapper():
    try:
        main()
    except ArtifactException as err:
        errmsg = str(err)
        if "404" in errmsg:
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
################################################################################
#                                 API                                          #
################################################################################
def api_do_create_artifact(args, config):
    param_check = _payload_check_(args, creation=True)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    artifact_id         = args["artifact"]["uuid"]
    artifact_alias      = args["artifact"]["alias"]
    artifact_name       = args["artifact"]["name"]
    artifact_type       = args["artifact"]["content_type"]
    artifact_checksum   = args["artifact"]["checksum"]
    artifact_label      = args["artifact"]["label"]
    artifact_openchain  = args["artifact"]["openchain"]
    private_key         = args["private_key"]
    public_key          = args["public_key"] 

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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.create_artifact(private_key, public_key, artifact_id, 
                                artifact_alias, artifact_name, artifact_type, 
                                artifact_checksum, artifact_label, 
                                artifact_openchain)
            
            return print_msg(response, "create")
        else:
            return output
    else:
        return output

def api_do_amend_artifact(args, config):
    param_check = _payload_check_(args)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    artifact_id         = args["artifact"]["uuid"]
    
    artifact_alias      = _null_cast(args["artifact"], "alias")
    artifact_name       = _null_cast(args["artifact"], "name")
    artifact_type       = _null_cast(args["artifact"], "content_type")
    artifact_checksum   = _null_cast(args["artifact"], "checksum")
    artifact_label      = _null_cast(args["artifact"], "label")
    artifact_openchain  = _null_cast(args["artifact"], "openchain")
    
    private_key         = args["private_key"]
    public_key          = args["public_key"] 

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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.amend_artifact(private_key, public_key, artifact_id, 
                                artifact_alias, artifact_name, artifact_type, 
                                artifact_checksum, artifact_label, 
                                artifact_openchain)
            
            return print_msg(response, "amend")
        else:
            return output
    else:
        return output

def api_do_list_artifact(config):
    b_url = config.get("DEFAULT", "url")
    client = ArtifactBatch(base_url=b_url)
    artifact_list = client.list_artifact()
    
    if artifact_list is not None:
        artifact_list.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(artifact_list)
        
        output = ret_msg("success", "OK", "ListOf:ArtifactRecord", result)
        
        return output
    else:
        return ret_msg(
                    "failed", 
                    "ArtifactException : Could not retrieve artifact listing.", 
                    "ArtifactRecord", "{}"
                )
                
def api_do_retrieve_artifact(
        artifact_id, config, all_flag=False, range_flag=None
    ):
    
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = ArtifactBatch(base_url=b_url)
    data = client.retrieve_artifact(artifact_id, all_flag, range_flag)

    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "ArtifactRecord", data.decode())
        else:
            output = ret_msg("success", "OK", "ArtifactRecord", data)
        
        return output
    else:
        return ret_msg(
                    "failed",
                    "ArtifactException : UUID {} does not exist." \
                    .format(artifact_id), 
                    "ArtifactRecord", "{}"
                )
                
def api_do_add_sub_artifact(args, config, del_flag=False):
    param_check = _payload_check_(args, cmd="AddArtifact")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    artifact_id     = args["relation"]["artifact_uuid"]
    sub_artifact_id = args["relation"]["sub_artifact_uuid"]
    path            = args["relation"]["path"]
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
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.add_artifact(
                            private_key, public_key, artifact_id,
                            sub_artifact_id, path, del_flag
                        )
                            
            return print_msg(response, "AddArtifact")
        else:
            return output
    else:
        return output
        
def api_do_add_uri_to_artifact(args, config, del_flag=False):
    param_check = _payload_check_(args, cmd="AddURI")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    artifact_id     = args["uuid"]
    
    version         = args["uri"]["version"]
    checksum        = args["uri"]["checksum"]
    content_type    = args["uri"]["content_type"]
    size            = args["uri"]["size"]
    uri_type        = args["uri"]["uri_type"]
    location        = args["uri"]["location"]
    
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
                    data=json.dumps(key),headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status") and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = ArtifactBatch(base_url=b_url)
            response = client.add_uri(
                            private_key, public_key, artifact_id, version,
                            checksum, content_type, size, uri_type, location,
                            del_flag
                        )
                                
            return print_msg(response, "AddURI")
        else:
            return output
    else:
        return output
################################################################################
#                           API PRIVATE FUNCTIONS                              #
################################################################################
def _payload_check_(args, creation=False, cmd=None):
    if cmd != None:
        if cmd == "AddArtifact":
            if "relation" not in args:
                return [True, "Relation missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "artifact_uuid" not in args["relation"]:
                return [True, "Artifact UUID missing."]
            elif "sub_artifact_uuid" not in args["relation"]:
                return [True, "Sub-Artifact UUID missing."]
            else:
                return [False]
        elif cmd == "AddURI":
            if "uuid" not in args:
                return [True, "UUID missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uri" not in args:
                return [True, "URI missing."]
            elif "version" not in args["uri"]:
                return [True, "URI : Version missing."]
            elif "checksum" not in args["uri"]:
                return [True, "URI : Checksum missing."]
            elif "size" not in args["uri"]:
                return [True, "URI : Size missing."]
            elif "content_type" not in args["uri"]:
                return [True, "URI : Content-type missing."]
            elif "uri_type" not in args["uri"]:
                return [True, "URI : URI-type missing."]
            elif "location" not in args["uri"]:
                return [True, "URI : Location missing."]
            else:
                return [False]
        else:
            return [False]
    else:
        if creation:
            if "artifact" not in args:
                return [True, "Artifact missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["artifact"]:
                return [True, "UUID missing."]
            elif "name" not in args["artifact"]:
                return [True, "Name missing."]
            elif "checksum" not in args["artifact"]:
                return [True, "Checksum missing."]
            elif "alias" not in args["artifact"]:
                return [True, "Alias missing."]
            elif "label" not in args["artifact"]:
                return [True, "Label missing."]
            elif "openchain" not in args["artifact"]:
                return [True, "Openchain missing."]
            elif "content_type" not in args["artifact"]:
                return [True, "Content_type missing."]    
            else:
                return [False]
        else:
            if "artifact" not in args:
                return [True, "Artifact missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["artifact"]:
                return [True, "UUID missing."]
            else:
                return [False]

def _null_cast(dic, key):
    if key not in dic:
        return "null"
    return dic[key]
################################################################################
#                                                                              #
################################################################################

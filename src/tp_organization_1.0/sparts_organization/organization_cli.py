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
from sparts_organization.organization_batch import OrganizationBatch
from sparts_organization.exceptions import OrganizationException

DISTRIBUTION_NAME = "sparts-organization-family"
################################################################################
def create_console_handler(verbose_level):
    """
    Helpes create a console handler for the Transaction Family : Organization.
    
    Returns:
        type: logging
        Logging object which contains the console handler config.
    
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
    Sets up logger for the Transaction Family : Organization
    
    Args:
        verbose_level (int): Verbose level of the logged message
        
    """
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    logger.addHandler(create_console_handler(verbose_level))
################################################################################
#                                   OBJ                                        #
################################################################################
def add_create_parser(subparsers, parent_parser):
    """
    Bash "create" subcommand handler for the Transaction Family : Organization
    
    Args:
        subparsers (ArgumentParser): Subcommand parser
        parent_parser (ArgumentParser):
            ArgumentParser object containing all the parameters
    
    """
    parser = subparsers.add_parser("create", parents=[parent_parser])

    parser.add_argument(
        "org_id",
        type=str,
        help="an identifier for the organization")
    
    parser.add_argument(
        "alias",
        type=str,
        help="Alias for the organization")
    
    parser.add_argument(
        "name",
        type=str,
        help="Provide organization name")
    
    parser.add_argument(
        "type",
        type=str,
        help="type of organization")
    
    parser.add_argument(
        "description",
        type=str,
        help="description ")

    parser.add_argument(
        "url",
        type=str,
        help="provide URL")
    
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

def add_list_organization_parser(subparsers, parent_parser):
    """
    Bash "list" subcommand handler for the Transaction Family : Organization
    
    Args:
        subparsers (ArgumentParser): Subcommand parser
        parent_parser (ArgumentParser):
            ArgumentParser object containing all the parameters
    
    """
    subparsers.add_parser("list-organization", parents=[parent_parser])

def add_retrieve_parser(subparsers, parent_parser):
    """
    Bash "retrieve" subcommand handler for the Transaction Family : Organization
    
    Args:
        subparsers (ArgumentParser): Subcommand parser
        parent_parser (ArgumentParser):
            ArgumentParser object containing all the parameters
    
    """
    parser = subparsers.add_parser("retrieve", parents=[parent_parser])

    parser.add_argument(
        "org_id",
        type=str,
        help="an identifier for the organization")
        
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
    Bash "amend" subcommand handler for the Transaction Family : Organization
    
    Args:
        subparsers (ArgumentParser): Subcommand parser
        parent_parser (ArgumentParser):
            ArgumentParser object containing all the parameters
            
    """
    parser = subparsers.add_parser("amend", parents=[parent_parser])
    
    parser.add_argument(
        "org_id",
        type=str,
        help="an identifier for the organization")
    
    parser.add_argument(
        "alias",
        type=str,
        help="Alias for the organization")
    
    parser.add_argument(
        "name",
        type=str,
        help="Provide organization name")
    
    parser.add_argument(
        "type",
        type=str,
        help="type of organization")
    
    parser.add_argument(
        "description",
        type=str,
        help="description ")

    parser.add_argument(
        "url",
        type=str,
        help="provide URL")
    
    
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

def add_part_parser(subparsers, parent_parser):
    """
    Bash "AddPart" subcommand handler for the Transaction Family : Organization
    
    Args:
        subparsers (ArgumentParser): Subcommand parser
        parent_parser (ArgumentParser):
            ArgumentParser object containing all the parameters
    
    """
    parser = subparsers.add_parser("AddPart", parents=[parent_parser])
    
    parser.add_argument(
        "org_id",
        type=str,
        help="the identifier for the organization")

    parser.add_argument(
        "pt_id",
        type=str,
        help="the identifier for Part")
    
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
        help="removes the pt")
################################################################################
#                                   CREATE                                     #
################################################################################
def create_parent_parser(prog_name):
    """
    Instantiates the ArgumentParser for the program.
    
    Args:
        prog_name (str): Name of the Transaction Family
    
    Returns:
        type: ArgumentParser
        ArgumentParser object with the basic configurations to perform a method
        for the program.
    
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
    Creates the ArgumentParser object which parses the bash input and stored
    the required parameters to perfrom the command on the
    Transaction Family : Organization
    
    Args:
        prog_name (str): Name of the Transaction Family
        
    Returns:
        type: ArgumentParser
        ArgumentParser object with all the required parameters stored to
        perform a method for the program.
    
    """
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        parents=[parent_parser],
        formatter_class=argparse.RawDescriptionHelpFormatter)

    subparsers = parser.add_subparsers(title="subcommands", dest="command")

    add_create_parser(subparsers, parent_parser)
    add_amend_parser(subparsers, parent_parser)
    
    add_list_organization_parser(subparsers, parent_parser)
    add_retrieve_parser(subparsers, parent_parser)
    
    add_part_parser(subparsers, parent_parser)
    
    return parser
################################################################################
#                               FUNCTIONS                                      #
################################################################################
def do_list_organization(args, config):
    """
    Lists out all the state associating with the UUIDs in the
    Transaction Family : Organization
    
    Args:
        config (ConfigParser): ConfigParser which contains the default url
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
    
    Raises:
        OrganizationException:
            * If failed to retrieve the list
            
    """
    b_url   = config.get("DEFAULT", "url")
    client  = OrganizationBatch(base_url=b_url)
    result  = client.list_organization()

    if result is not None:
        result.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(result)
        
        output = ret_msg("success", "OK", "ListOf:OrganizationRecord", result)
        
        print(output)
    else:
        raise OrganizationException("Could not retrieve organization listing.")

def do_retrieve_organization(args, config):
    """
    Retrieves the state associating with the UUID in the
    Transaction Family : Organization
    
    Args:
        args (ArgumentParser):
            ArgumentParser object containing required parameters
        config (ConfigParser): ConfigParser which contains the default url
        
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
    
    Raises:
        OrganizationException:
            * If failed to retrieve the uuid
    
    """
    all_flag    = args.all
    range_flag  = args.range
    
    org_id      = args.org_id
    
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = OrganizationBatch(base_url=b_url)
    data = client.retrieve_organization(org_id, all_flag, range_flag)
    
    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "OrganizationRecord", 
                        data.decode())
        else:
            output = ret_msg("success", "OK", "OrganizationRecord", data)
        
        print(output)
    else:
        raise OrganizationException("Organization not found: {}".format(org_id))

def do_create_organization(args, config):
    """
    Creates the state associating with the UUID in the
    Transaction Family : Organization
    
    Args:
        args (ArgumentParser):
            ArgumentParser object containing required parameters
        config (ConfigParser): ConfigParser which contains the default url
        
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
    
    """
    org_id      = args.org_id
    org_alias   = args.alias
    org_name    = args.name
    org_type    = args.type
    description = args.description
    org_url     = args.url
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
                    data=json.dumps(key), headers=headers)
    
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = OrganizationBatch(base_url=b_url)
            response = client.create_organization(org_id, org_alias, org_name,
                            org_type, description, org_url, private_key,
                            public_key
                        )
            print_msg(response, "create")
        else:
            print(output)
    else:
        print(output)

def do_amend_organization(args, config):
    """
    Amends the state associating with the UUID in the
    Transaction Family : Organization
    
    Args:
        args (ArgumentParser):
            ArgumentParser object containing required parameters
        config (ConfigParser): ConfigParser which contains the default url
        
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
    
    """
    org_id      = args.org_id
    org_alias   = args.alias
    org_name    = args.name
    org_type    = args.type
    description = args.description
    org_url     = args.url
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
                    data=json.dumps(key), headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)

    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = OrganizationBatch(base_url=b_url)
            response = client.amend_organization(org_id, org_alias, org_name,
                            org_type, description, org_url, private_key,
                            public_key
                        )
                        
            print_msg(response, "amend")
        else:
            print(output)
    else:
        print(output)

def do_addpart(args, config):
    """
    Establishes relationship between Organization and Part in the state
    associating with the UUID of the Transaction Family : Organization
    
    Args:
        args (ArgumentParser):
            ArgumentParser object containing required parameters
        config (ConfigParser): ConfigParser which contains the default url
        
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    """
    deletePart  = args.delete
    
    org_id      = args.org_id
    pt_id       = args.pt_id
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
                data=json.dumps(key), headers=headers)
    output = response.content.decode("utf-8").strip()
    statusinfo = json.loads(output)
       
    if statusinfo.get("status")and statusinfo.get("message"):
            
        status = statusinfo["status"]
        message = statusinfo["message"]
            
        if status == "success" and message == "authorized":
            b_url = config.get("DEFAULT", "url")
            client = OrganizationBatch(base_url=b_url)
            response = client.add_part(org_id, pt_id, private_key, 
                                public_key, deletePart)
                                
            print_msg(response, "AddPart")
        else:
            print(output)
    else:
        print(output)
################################################################################
#                                  PRINT                                       #
################################################################################   
def load_config():
    """
    Helps construct ConfigParser object pertaining default url for
    the sawtooth validator.
    
    Returns:
        type: ConfigParser
        ConfigParser object with default url.
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    return config

def print_msg(response, cmd=None):
    """
    Helps create the return message for the terminal or the web-browser.
    
    Args:
        response (None or list containing None and str):
            Contains the data for the function to construct return message
        cmd (None or str): The subcommand which was performed
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure. 
    
    Raises:
        OrganizationException:
            * If response is None
            * If response is unknown
            * If response is a list with None
    
    """
    try:
        if type(response) is list and response[0] == None:
            raise OrganizationException(
                        "OrganizationException : No change."
                    )
        
        if response == None:
            if cmd == "create":
                raise OrganizationException(
                            "OrganizationException : Duplicate UUID."
                        )
                
            elif cmd == "amend" or cmd == "AddPart":
                raise OrganizationException(
                            "OrganizationException : UUID does not exist."
                        )
                
            raise OrganizationException("Exception raised.")
        elif "batch_statuses?id" in response:
            print(ret_msg("success", "OK", "OrganizationRecord", "{}"))
            return ret_msg("success", "OK", "OrganizationRecord", "{}")
        else:
            raise OrganizationException("Exception raised.")
    except BaseException as err:
        output = ret_msg(
                            "failed",
                            str(err),
                            "OrganizationRecord", "{}"
                        )
        print(output)
        return output
    
def ret_msg(status, message, result_type, result):
    """
    Helps create the message to be returned.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
    
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
        do_create_organization(args, config)
    elif args.command == "list-organization":
        do_list_organization(args, config)
    elif args.command == "retrieve":
        do_retrieve_organization(args, config)
    elif args.command == "AddPart":
        do_addpart(args, config)
    elif args.command == "amend":
        do_amend_organization(args, config)
    else:
        raise OrganizationException("invalid command: {}".format(args.command))

def main_wrapper():
    try:
        main()
    except OrganizationException as err:
        errmsg = str(err)
        if "404" in errmsg:
            exp = ret_msg("failed","404 Not Found","EmptyRecord","{}")
            print(OrganizationException(exp))
           
        else:
            exp = ret_msg("failed",errmsg,"EmptyRecord","{}")
            print(OrganizationException()) 
        sys.exit(1)
    except KeyboardInterrupt:
        pass
    except SystemExit as err:
        raise err
    except BaseException as err:
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)
################################################################################
#                                   API                                        #
################################################################################
def api_do_create_organization(args, config):
    """
    API version of "do_create_organization" function.
    """
    param_check = _payload_check_(args, creation=True)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    org_id      = args["organization"]["uuid"]
    org_alias   = args["organization"]["alias"]
    org_name    = args["organization"]["name"]
    org_type    = args["organization"]["type"]
    description = args["organization"]["description"]
    org_url     = args["organization"]["url"]
    private_key = args["private_key"]
    public_key  = args["public_key"]

    payload             = "{}"
    key                 = json.loads(payload)
    key["publickey"]    = public_key
    key["privatekey"]   = private_key
    key["allowedrole"]  = [{"role" : "admin"}, {"role" : "member"}]
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
            client = OrganizationBatch(base_url=b_url)
            response = client.create_organization(org_id, org_alias, org_name,
                            org_type, description, org_url, private_key,
                            public_key
                        )
            
            return print_msg(response, "create")
        else:
            return output
    else:
        return output

def api_do_amend_organization(args, config):
    """
    API version of "do_amend_organization" function.
    """
    param_check = _payload_check_(args)
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    org_id      = args["organization"]["uuid"]
    
    org_alias   = _null_cast(args["organization"], "alias")
    org_name    = _null_cast(args["organization"], "name")
    org_type    = _null_cast(args["organization"], "type")
    description = _null_cast(args["organization"], "description")
    org_url     = _null_cast(args["organization"], "url")
    
    private_key = args["private_key"]
    public_key  = args["public_key"]

    payload             = "{}"
    key                 = json.loads(payload)
    key["publickey"]    = public_key
    key["privatekey"]   = private_key
    key["allowedrole"]  = [{"role" : "admin"}, {"role" : "member"}]
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
            client = OrganizationBatch(base_url=b_url)
            response = client.amend_organization(org_id, org_alias, org_name,
                            org_type, description, org_url, private_key,
                            public_key
                        )
            
            return print_msg(response, "amend")
        else:
            return output
    else:
        return output

def api_do_list_organization(config):
    """
    API version of "do_list_organization" function.
    """
    b_url   = config.get("DEFAULT", "url")
    client  = OrganizationBatch(base_url=b_url)
    organization_list  = client.list_organization()

    if organization_list is not None:
        organization_list.sort(key=lambda x:x["timestamp"], reverse=True)
        result = json.dumps(organization_list)
        
        output = ret_msg("success", "OK", "ListOf:OrganizationRecord", result)
        
        return output
    else:
        return ret_msg(
                    "failed", 
                    "{} : Could not retrieve {}.".format(
                        "OrganizationException", "organization listing"
                    ),
                    "OrganizationRecord", "{}"
                )

def api_do_retrieve_organization(org_id, config, all_flag=False,
                                    range_flag=None):
    """
    API version of "do_retrieve_organization" function.
    """
    if range_flag != None:
        all_flag = True
    
    b_url = config.get("DEFAULT", "url")
    client = OrganizationBatch(base_url=b_url)
    data = client.retrieve_organization(org_id, all_flag, range_flag)
    
    if data is not None:
        
        if all_flag == False:
            output = ret_msg("success", "OK", "OrganizationRecord", 
                        data.decode())
        else:
            output = ret_msg("success", "OK", "OrganizationRecord", data)
        
        return output
    else:
        return ret_msg(
                    "failed",
                    "OrganizationException : UUID {} does not exist." \
                    .format(org_id),
                    "OrganizationRecord", "{}"
                )
                
def api_do_addpart(args, config, del_flag=False):
    """
    API version of "do_addpart" function.
    """
    param_check = _payload_check_(args, cmd="AddPart")
    
    if param_check[0]:
        return ret_msg("failed", param_check[1], "EmptyRecord", "{}")
    
    org_id      = args["relation"]["organization_uuid"]
    pt_id       = args["relation"]["part_uuid"]
    private_key = args["private_key"]
    public_key  = args["public_key"]
    
    payload             = "{}"
    key                 = json.loads(payload)
    key["publickey"]    = public_key
    key["privatekey"]   = private_key
    key["allowedrole"]  = [{"role" : "admin"}, {"role" : "member"}]
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
            client = OrganizationBatch(base_url=b_url)
            response = client.add_part(org_id, pt_id, private_key, 
                                public_key, del_flag)
                                
            return print_msg(response, "AddPart")
        else:
            return output
    else:
        return output
################################################################################
#                           API PRIVATE FUNCTIONS                              #
################################################################################
def _payload_check_(args, creation=False, cmd=None):
    """
    Checks payload for correct JSON format for a given command.
    
    Args:
        args (dict): Pass in payload
        creation (bool): True if "create", false otherwise
        cmd (None or str): str if "Add...", None otherwise
    
    Returns:
        type: list containing bool or bool and str
        List with False or list with True and error message. False stands for
        do not terminate the process.
        
    """
    if cmd != None:
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
    else:
        if creation:
            if "organization" not in args:
                return [True, "Organization missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["organization"]:
                return [True, "UUID missing."]
            elif "alias" not in args["organization"]:
                return [True, "Alias missing."]
            elif "name" not in args["organization"]:
                return [True, "Name missing."]
            elif "type" not in args["organization"]:
                return [True, "Type missing."]
            elif "description" not in args["organization"]:
                return [True, "Description missing."]
            elif "url" not in args["organization"]:
                return [True, "URL missing."]
            else:
                return [False]
        else:
            if "organization" not in args:
                return [True, "Organization missing."]
            elif "private_key" not in args:
                return [True, "Private-Key missing."]
            elif "public_key" not in args:
                return [True, "Public-Key missing."]
            elif "uuid" not in args["organization"]:
                return [True, "UUID missing."]
            else:
                return [False]

def _null_cast(dic, key):
    """
    Allows the user to load value, given key from the dictionary. If the key
    is not found, return "null".
    
    Args:
        dic (dict): Dictionary in look for (key, value) pair
        key (str): Key to look search in the dictionary
        
    Returns:
        type: str
        Either "null" string or previous data stored in the field.
    
    """
    if key not in dic:
        return "null"
    return dic[key]
################################################################################
#                                                                              #
################################################################################

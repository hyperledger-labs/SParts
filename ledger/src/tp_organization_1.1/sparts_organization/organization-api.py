# Copyright 2017 Wind River Systems
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
#                         LIBRARIES & DEPENDENCIES                             #
################################################################################
from flask import Flask, jsonify, make_response, request, json
import organization_cli
import configparser
################################################################################
#                                FLASK APP                                     #
################################################################################
app = Flask(__name__)

# PING
@app.route("/tp/organization/ping", methods=["GET"])
def get_ping_result():
    """
    Allows the client side API call to "ping" the port localhost:851 to ensure
    that the port is up and running.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    """
    output = ret_msg("success", "OK", "EmptyRecord", "Organization")
    return output 

# CREATE
@app.route("/tp/organization", methods=["POST"])
def create_organization():
    """
    Allows the client side API call to "create" the organization given a correct
    JSON formatted payload.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
        
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
            
        output = organization_cli \
                    .api_do_create_organization(request.json, config)    
        
        return output
    except Exception as e:
        return e

# AMEND
@app.route("/tp/organization/amend", methods=["POST"])
def amend_organization():
    """
    Allows the client side API call to "amend" the organization given a correct
    JSON formatted payload.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
            
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = organization_cli \
                    .api_do_amend_organization(request.json, config)    
        
        return output
    except Exception as e:
        return e

# LIST
@app.route("/tp/organization", methods=["GET"])
def list_organization():
    """
    Allows the client side API call to "list" the organization.
    
    Returns:
        type: str
        String representing JSON object which contains the result of
        the "organization list-organization" if the call was a success; else,
        JSON object which contains error message.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = organization_cli.api_do_list_organization(config)
        
        return output
    except Exception as e:
        return e

# RETRIEVE MOST RECENT BY UUID
@app.route("/tp/organization/<string:organization_id>", methods=["GET"])
def retrieve_organization(organization_id):
    """
    Allows the client side API call to "retrieve" the organization.
    
    Returns:
        type: str
        String representing JSON object which contains the result of
        the "organization retrieve {uuid}" if the call was a success; else,
        JSON object which contains error message.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = organization_cli.api_do_retrieve_organization(
                    organization_id, config
                )
        
        return output
    except Exception as e:
        return e

# RETRIEVE HISTORY OF UUID
@app.route("/tp/organization/history/<string:organization_id>", methods=["GET"])
def retrieve_organization_history(organization_id):
    """
    Allows the client side API call to "retrieve" the organization and display
    its history up to its creation block.
    
    Returns:
        type: str
        String representing JSON object which contains the result of
        the "organization retrieve --all {uuid}" if the call was a success;
        else, JSON object which contains error message.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = organization_cli.api_do_retrieve_organization(
                        organization_id, config, all_flag=True
                    )
        return output
    except Exception as e:
        return e

# RETRIEVE UUID ON CERTAIN DATE     
@app.route(
    "/tp/organization/<string:organization_id>/date/<string:START>",
    methods=["GET"]
)
def retrieve_organization_history_date(organization_id, START):
    """
    Allows the client side API call to "retrieve" the organization and display
    its history for the specified date.
    
    Returns:
        type: str
        String representing JSON object which contains the result of
        the "organization retrieve --range START END {uuid}" if the call was
        a success; else, JSON object which contains error message.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = organization_cli.api_do_retrieve_organization(
                        organization_id, config, range_flag=["0", START]
                    )
        return output
    except Exception as e:
        return e
   
# ADD PART
@app.route("/tp/organization/addpart", methods=["POST"])
def add_part_organization():
    """
    Allows the client side API call to "AddPart" to the organization given
    a correct JSON formatted payload.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = organization_cli.api_do_addpart(request.json, config)    
        
        return output
    except Exception as e:
        return e
        
# ADD PART --DELETE
@app.route("/tp/organization/addpart/delete", methods=["POST"])
def add_part_organization_delete():
    """
    Allows the client side API call to "AddPart --delete" to the organization
    given a correct JSON formatted payload.
    
    Returns:
        type: str
        String representing JSON object which allows the client to know that
        the call was either a success or a failure.
        
    Raises:
        Exception:
            * If the request does not contain JSON payload
    
    """
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = organization_cli.api_do_addpart(request.json, config, True)    
        
        return output
    except Exception as e:
        return e
################################################################################
#                                  PRINT                                       #
################################################################################
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
    key["result"] = result
    msgJSON = json.dumps(key)
    return msgJSON
################################################################################
#                                   MAIN                                       #
################################################################################
if __name__ == "__main__":
    app.run(host="0.0.0.0", port="851")
################################################################################
#                                                                              #
################################################################################


# copyright 2017 Wind River Systems
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

from flask import Flask, jsonify, make_response, request, json
import organization_cli
import configparser
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
app = Flask(__name__)

# PING
@app.route("/tp/organization/ping", methods=["GET"])
def get_ping_result():
    
    output = ret_msg("success","OK","EmptyRecord","Organization")
    return output 

# CREATE
@app.route("/tp/organization", methods=["POST"])
def create_organization():
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
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = organization_cli.api_do_retrieve_organization(
                        organization_id, config, range_flag=[START, START]
                    )
        return output
    except Exception as e:
        return e
   
# ADD PART
@app.route("/tp/organization/addpart", methods=["POST"])
def add_part_organization():
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
#                                   TEST                                       #
################################################################################
@app.route("/tp/test", methods=["POST"])
def testing_():
    try:
        return json.dumps(request.json)
    except Exception as e:
        return e

@app.route("/tp/test", methods=["GET"])
def testing_get():
    try:
        return "phyo test get was called successfully"
    except Exception as e:
        return e
################################################################################
#                                  PRINT                                       #
################################################################################
def ret_msg(status, message, result_type, result):
    msgJSON = "{}"
    key = json.loads(msgJSON)
    key["status"] = status
    key["message"] = message
    key["result_type"] = result_type
    key["result"] = result
    msgJSON = json.dumps(key)
    return msgJSON
    
# @app.errorhandler(500)
# def custom500(message):
#     return "yolo"
################################################################################
#                                   MAIN                                       #
################################################################################
if __name__ == "__main__":
    app.run(host="0.0.0.0", port="851")
################################################################################
#                                                                              #
################################################################################

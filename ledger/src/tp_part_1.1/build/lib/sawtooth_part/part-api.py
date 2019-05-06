
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
import part_cli
import configparser
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
app = Flask(__name__)

# PING
@app.route("/tp/part/ping", methods=["GET"])
def get_ping_result():
    
    output = ret_msg("success","OK","EmptyRecord","Part")
    return output 

# CREATE
@app.route("/tp/part", methods=["POST"])
def create_part():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
            
        output = part_cli.api_do_create_part(request.json, config)    
        
        return output
    except Exception as e:
        return e

# AMEND
@app.route("/tp/part/amend", methods=["POST"])
def amend_part():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_amend_part(request.json, config)    
        
        return output
    except Exception as e:
        return e

# LIST
@app.route("/tp/part", methods=["GET"])
def list_part():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = part_cli.api_do_list_part(config)
        
        return output
    except Exception as e:
        return e

# RETRIEVE MOST RECENT BY UUID
@app.route("/tp/part/<string:part_id>", methods=["GET"])
def retrieve_part(part_id):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = part_cli.api_do_retrieve_part(
                    part_id, config
                )
        
        return output
    except Exception as e:
        return e

# RETRIEVE HISTORY OF UUID
@app.route("/tp/part/history/<string:part_id>", methods=["GET"])
def retrieve_part_history(part_id):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = part_cli.api_do_retrieve_part(
                        part_id, config, all_flag=True
                    )
        return output
    except Exception as e:
        return e

# RETRIEVE UUID ON CERTAIN DATE     
@app.route(
    "/tp/part/<string:part_id>/date/<string:START>",
    methods=["GET"]
)
def retrieve_part_history_date(part_id, START):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = part_cli.api_do_retrieve_part(
                        part_id, config, range_flag=[START, START]
                    )
        return output
    except Exception as e:
        return e
   
# ADD ORGANIZATION
@app.route("/tp/part/addorganization", methods=["POST"])
def add_part_organization():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_organization(request.json, config)    
        
        return output
    except Exception as e:
        return e

# ADD ORGANIZATION --DELETE
@app.route("/tp/part/addorganization/delete", methods=["POST"])
def add_part_organization_delete():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_organization(request.json, config, True)    
        
        return output
    except Exception as e:
        return e
        
# ADD CATEGORY
@app.route("/tp/part/addcategory", methods=["POST"])
def add_part_category():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_category(request.json, config)    
        
        return output
    except Exception as e:
        return e
        
# ADD CATEGORY --DELETE
@app.route("/tp/part/addcategory/delete", methods=["POST"])
def add_part_category_delete():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_category(request.json, config, True)    
        
        return output
    except Exception as e:
        return e
        
# ADD ARTIFACT
@app.route("/tp/part/addartifact", methods=["POST"])
def add_part_artifact():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_artifact(request.json, config)    
        
        return output
    except Exception as e:
        return e
        
# ADD ARTIFACT --DELETE
@app.route("/tp/part/addartifact/delete", methods=["POST"])
def add_part_artifact_delete():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = part_cli.api_do_add_artifact(request.json, config, True)    
        
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
################################################################################
#                                   MAIN                                       #
################################################################################
if __name__ == "__main__":
    app.run(host="0.0.0.0", port="852")
################################################################################
#                                                                              #
################################################################################

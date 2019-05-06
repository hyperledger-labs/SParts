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

from flask import Flask, jsonify, make_response, request, json
import artifact_cli
import configparser
################################################################################
#                               LIBS & DEPS                                    #
################################################################################
app = Flask(__name__)

# PING
@app.route("/tp/artifact/ping", methods=["GET"])
def get_ping_result():
    
    output = ret_msg("success", "OK", "EmptyRecord", "Artifact")
    return output 

# CREATE
@app.route("/tp/artifact", methods=["POST"])
def create_artifact():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
            
        output = artifact_cli.api_do_create_artifact(request.json, config)  
        
        return output
    except Exception as e:
        return e

# AMEND
@app.route("/tp/artifact/amend", methods=["POST"])
def amend_artifact():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = artifact_cli.api_do_amend_artifact(request.json, config)    
        
        return output
    except Exception as e:
        return e

# LIST
@app.route("/tp/artifact", methods=["GET"])
def list_artifact():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = artifact_cli.api_do_list_artifact(config)
        
        return output
    except Exception as e:
        return e

# RETRIEVE MOST RECENT BY UUID
@app.route("/tp/artifact/<string:artifact_id>", methods=["GET"])
def retrieve_artifact(artifact_id):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = artifact_cli.api_do_retrieve_artifact(artifact_id, config)
        
        return output
    except Exception as e:
        return e

# RETRIEVE HISTORY OF UUID
@app.route("/tp/artifact/history/<string:artifact_id>", methods=["GET"])
def retrieve_artifact_history(artifact_id):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = artifact_cli.api_do_retrieve_artifact(
                        artifact_id, config, all_flag=True
                    )
        return output
    except Exception as e:
        return e

# RETRIEVE UUID ON CERTAIN DATE     
@app.route(
    "/tp/artifact/<string:artifact_id>/date/<string:START>",
    methods=["GET"]
)
def retrieve_artifact_history_date(artifact_id, START):
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        output = artifact_cli.api_do_retrieve_artifact(
                        artifact_id, config, range_flag=[START, START]
                    )
        return output
    except Exception as e:
        return e
        
# ADD ARTIFACT
@app.route("/tp/artifact/addartifact", methods=["POST"])
def add_artifact_artifact():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = artifact_cli.api_do_add_sub_artifact(request.json, config)    
        
        return output
    except Exception as e:
        return e
        
# ADD ARTIFACT --DELETE
@app.route("/tp/artifact/addartifact/delete", methods=["POST"])
def add_artifact_artifact_delete():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = artifact_cli.api_do_add_sub_artifact(
                        request.json, config, True
                    )    
        
        return output
    except Exception as e:
        return e
        
# ADD URI
@app.route("/tp/artifact/adduri", methods=["POST"])
def add_artifact_uri():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = artifact_cli.api_do_add_uri_to_artifact(request.json, config)    
        
        return output
    except Exception as e:
        return e
        
# ADD URI --DELETE
@app.route("/tp/artifact/adduri/delete", methods=["POST"])
def add_artifact_uri_delete():
    config = configparser.ConfigParser()
    config.set("DEFAULT", "url", "http://127.0.0.1:8008")
    
    try:
        if not request.json:
            return "Expecting JSON Object."
        
        output = artifact_cli.api_do_add_uri_to_artifact(
                        request.json, config, True
                    )    
        
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
    app.run(host="0.0.0.0", port="853")
################################################################################
#                                                                              #
################################################################################


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

#!flask/bin/python
import os
import subprocess, shlex, re
from flask import Flask, jsonify, make_response, request, json
import requests
# import sawtooth_signing.secp256k1_signer as signing
#
from sawtooth_signing import create_context
from sawtooth_signing import CryptoFactory
from sawtooth_signing import ParseError
from sawtooth_signing.secp256k1 import Secp256k1PublicKey
from sawtooth_signing.secp256k1 import Secp256k1PrivateKey
#
import base64

import uuid
from random import randint
app = Flask(__name__)

@app.route("/ledger/api/v1/ping", methods=["GET"])
def get_ping_result():
    
    output = ret_msg("success","OK","EmptyRecord","{}")
    return output    
################################################################################
#                                  ARTIFACT                                    #
################################################################################
def create_artifact_cmd(artifact_id, alias, filename, content_type, checksum,
                            label, openchain, private_key, public_key):
    return "artifact create {} {} {} {} {} {} {} {} {}".format(
                str_qt(artifact_id), str_qt(alias), str_qt(filename),
                str_qt(content_type), str_qt(checksum), str_qt(label),
                str_qt(openchain), str_qt(private_key), str_qt(public_key)
            )
    # return "artifact create {} {} {} {} {} {} {} {} {}".format(
    #             artifact_id, alias, filename, content_type, checksum, label,
    #             openchain, private_key, public_key
    #         )
    # cmd = "artifact create " + str_qt(artifact_id) + " " + str_qt(alias) + " " + str_qt(filename)+ " " + str_qt(content_type) + " " + str_qt(checksum) + " " + str_qt(label) + " " + str_qt(openchain) + " "+ str_qt(private_key) + " "+ str_qt(public_key)
    # return cmd

@app.route("/ledger/api/v1/artifacts", methods=["POST"])
def create_artifact():
    try:
        if (not request.json or 
            "private_key" not in request.json or
            "public_key" not in request.json or 
            "artifact" not in request.json):
            return ret_exception_msg("Invalid Input")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        artifact_id = request.json["artifact"]["uuid"]
        filename = request.json["artifact"]["name"]
        alias = request.json["artifact"]["alias"]
        checksum = request.json["artifact"]["checksum"]
        content_type = request.json["artifact"]["content_type"]       
        label = request.json["artifact"]["label"]
        openchain = request.json["artifact"]["openchain"]
          
        cmd = create_artifact_cmd(artifact_id, alias, filename, content_type,
                            checksum, label, openchain, private_key, public_key) 
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output 
    except Exception as e:
        exp =  ret_exception_msg(e)
        return exp

def amend_artifact_cmd(artifact_id, alias, filename, content_type, checksum,
                            label, openchain, private_key, public_key):
    return "artifact amend {} {} {} {} {} {} {} {} {}".format(
                str_qt(artifact_id), str_qt(alias), str_qt(filename),
                str_qt(content_type), str_qt(checksum), str_qt(label),
                str_qt(openchain), str_qt(private_key), str_qt(public_key)
            )
    # return "artifact amend {} {} {} {} {} {} {} {} {}".format(
    #             artifact_id, alias, filename, content_type, checksum, label,
    #             openchain, private_key, public_key
    #         )
    # cmd = "artifact amend " + str_qt(artifact_id) + " " + str_qt(alias) + " " + str_qt(filename)+ " " + str_qt(content_type) + " " + str_qt(checksum) + " " + str_qt(label) + " " + str_qt(openchain) + " "+ str_qt(private_key) + " "+ str_qt(public_key)
    # return cmd

@app.route("/ledger/api/v1/artifacts/amend", methods=["POST"])
def amend_artifact():
    try:
        if (not request.json or 
            "private_key" not in request.json or
            "public_key" not in request.json or
            "artifact" not in request.json):
            return ret_exception_msg("Invalid Input")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        artifact_id = request.json["artifact"]["uuid"]
        
        filename = nullCast(request.json["artifact"], "name")
        alias = nullCast(request.json["artifact"], "alias") 
        checksum = nullCast(request.json["artifact"], "checksum") 
        content_type = nullCast(request.json["artifact"], "content_type")      
        label = nullCast(request.json["artifact"], "label")
        openchain = nullCast(request.json["artifact"], "openchain")
          
        cmd = amend_artifact_cmd(artifact_id, alias, filename, content_type,
                            checksum, label, openchain, private_key, public_key) 
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output 
    except Exception as e:
        exp =  ret_exception_msg(e)
        return exp

@app.route("/ledger/api/v1/artifacts", methods=["GET"])
def get_artifacts():
    try:
        cmd = "artifact list-artifact"
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = "" 
        for line in process.stdout:
            output += line.decode("UTF-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

@app.route("/ledger/api/v1/artifacts/<string:artifact_id>", methods=["GET"])
def get_artifact(artifact_id):
    try:
        cmd = "artifact retrieve " + artifact_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e)
        return exp

def add_artifact_to_envelope_cmd(envelope_uuid, artifact_uuid, private_key,
                                    public_key, path):
    return "artifact AddArtifact {} {} {} {} {}".format(
                str_qt(envelope_uuid), str_qt(artifact_uuid), str_qt(path),
                str_qt(private_key), str_qt(public_key)
            )
    # return "artifact AddArtifact {} {} {} {} {}".format(
    #             envelope_uuid, artifact_uuid, path, private_key, public_key
    #         )
    # cmd = "artifact AddArtifact " + str_qt(envelope_uuid) + " " + str_qt(artifact_uuid) + " " + str_qt(path) + " " + str_qt(private_key) + " " + str_qt(public_key)
    # return cmd

@app.route("/ledger/api/v1/envelope/artifact", methods=["POST"])
def add_artifact_to_envelope():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json or
            "relation" not in request.json):
            return ret_exception_msg("Invalid Input")
        output = "" 
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
         
        artifact_uuid = request.json["relation"]["artifact_uuid"]
        envelope_uuid = request.json["relation"]["envelope_uuid"]
        path = request.json["relation"]["path"]
        
        cmd = add_artifact_to_envelope_cmd(envelope_uuid, artifact_uuid,
                                            private_key, public_key, path)
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output 
    except Exception as e:
        exp =  ret_exception_msg(e)
        return exp

def add_uri_cmd(artifact_id, version, checksum, content_type, size, uri_type,
                location, private_key, public_key):
    return "artifact AddURI {} {} {} {} {} {} {} {} {}".format(
                str_qt(artifact_id), str_qt(version), str_qt(checksum),
                str_qt(content_type), size, str_qt(uri_type), str_qt(location),
                str_qt(private_key), str_qt(public_key)
            )
    # return "artifact AddURI {} {} {} {} {} {} {} {} {}".format(
    #             artifact_id, version, checksum, content_type, size, uri_type,
    #             location, private_key, public_key
    #         )
    # cmd = "artifact AddURI " + str_qt(artifact_id) + " " + str_qt(version) + " " + str_qt(checksum) + " " + str_qt(content_type) + " " + size +  " " +str_qt(uri_type)+" "+ str_qt(location) + " " + str_qt(private_key) + " "+ str_qt(public_key)
    # return cmd

# Create record for the artifact  
@app.route("/ledger/api/v1/artifacts/uri", methods=["POST"])
def add_uri_to_artifact():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json or
            "uri" not in request.json or
            "uuid" not in request.json):
            return ret_exception_msg("Invalid Input")
        output = "" 
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
         
        artifact_id = request.json["uuid"]
        version = request.json["uri"]["version"]
        checksum = request.json["uri"]["checksum"]
        content_type = request.json["uri"]["content_type"] 
        size = request.json["uri"]["size"]
        uri_type = request.json["uri"]["uri_type"]
        location = request.json["uri"]["location"]
              
        cmd = add_uri_cmd(artifact_id, version, checksum, content_type, size,
                            uri_type, location, private_key, public_key) 
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output 
    except Exception as e:
        exp =  ret_exception_msg(e)
        return exp
################################################################################
#                                     USER                                     #
################################################################################    
@app.route("/ledger/api/v1/registeruser", methods=["POST"])
def register_user():
    try:
        if (not request.json or 
            "private_key" not in request.json or
            "public_key" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
            
        user_public_key = request.json["user"]["user_public_key"]
        name = request.json["user"]["user_name"]
        name = format_str(name) 
        email_address = request.json["user"]["email_address"]
        email_address = format_str(email_address)
        authorized = request.json["user"]["authorized"]
        role = request.json["user"]["role"]     
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        # cmd = "user register " + user_public_key + " " + str(name) + " " + str(email_address) + " " + authorized +" " + role+ " " + private_key + " "+ public_key
        cmd = "user register {} {} {} {} {} {} {}".format(
                    user_public_key, str(name), str(email_address), authorized,
                    role, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp
################################################################################
#                                   CATEGORY                                   #
################################################################################    
# Create record for category
@app.route("/ledger/api/v1/categories", methods=["POST"])
def create_category():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return "Invalid JSON"
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        uuid = request.json["category"]["uuid"]
        
        name = request.json["category"]["name"]
        
        description = request.json["category"]["description"]
        
        cmd = "category create {} {} {} {} {}".format(
                    uuid, name, description, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Amend category record
@app.route("/ledger/api/v1/categories/amend", methods=["POST"])
def amend_category():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return "Invalid JSON"
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        uuid = request.json["category"]["uuid"]
        
        name = nullCast(request.json["category"], "name")
        
        description = nullCast(request.json["category"], "description")
        
        cmd = "category amend {} {} {} {} {}".format(
                    uuid, name, description, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves categories list
@app.route("/ledger/api/v1/categories", methods=["GET"])
def get_categories():
    try:
        cmd = "category list-category"
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("UTF-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp
        
# Retrieves most recent category record by category id   
@app.route("/ledger/api/v1/categories/<string:category_id>", methods=["GET"])
def get_category(category_id):
    try:
        cmd = "category retrieve " + category_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves historical category record by category id   
@app.route(
    "/ledger/api/v1/categories/history/<string:category_id>",
    methods=["GET"]
)
def get_uuid_category_history(category_id):
    try:
        cmd = "category retrieve --all " + category_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves historical category record on certain date by category id
@app.route(
    "/ledger/api/v1/categories/<string:category_id>/date/<string:START>",
    methods=["GET"]
)
def get_uuid_category_day(category_id, START):
    try:
        cmd = "category retrieve --range {} {} {}".format(
                    START, START, category_id
                )
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieve most appropriate category record on certain date by category id
# TODO : --limit needs to be implemented first
################################################################################
#                                 ORGANIZATION                                 #
################################################################################
@app.route("/ledger/api/v1/orgs", methods=["POST"])
def create_organization():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"] 
        
        uuid = request.json["organization"]["uuid"]
        
        alias = request.json["organization"]["alias"]
        
        name = request.json["organization"]["name"]
         
        _type = request.json["organization"]["type"]
        
        description = request.json["organization"]["description"]
        
        url = request.json["organization"]["url"]

        cmd = "organization create {} {} {} {} {} {} {} {}".format(
                    uuid, alias, name, _type, description, url,
                    private_key, public_key
                )
        cmd = shlex.split(cmd)
    
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e)
        return exp

@app.route("/ledger/api/v1/orgs/amend", methods=["POST"])
def amend_organization():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        uuid = request.json["organization"]["uuid"]
        
        alias = nullCast(request.json["organization"], "alias")
        
        name = nullCast(request.json["organization"], "name")
        
        org_type = nullCast(request.json["organization"], "type")
        
        description = nullCast(request.json["organization"], "description")
        
        url = nullCast(request.json["organization"], "url")
        
        cmd = "organization amend {} {} {} {} {} {} {} {}".format(
                    uuid, alias, name, org_type, description, url,
                    private_key, public_key
                )
        cmd = shlex.split(cmd)
    
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e)
        return exp

# Retrieves list of all organizations from the ledger
@app.route("/ledger/api/v1/orgs", methods=["GET"])
def get_organizations():
    try:
        cmd = "organization list-organization"
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
                output += line.decode("utf-8").strip()
                output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves organization record by organization id
@app.route("/ledger/api/v1/orgs/<string:org_id>", methods=["GET"])
def get_organization(org_id):
    try:
        cmd = "organization retrieve " + org_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves historical organization record by organization id
@app.route("/ledger/api/v1/orgs/history/<string:org_id>", methods=["GET"])
def get_uuid_organization_history(org_id):
    try:
        cmd = "organization retrieve --all " + org_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves historical organization record on certain date by organization id
@app.route("/ledger/api/v1/orgs/<string:org_id>/date/<string:START>", methods=["GET"])
def get_uuid_organization_day(org_id, START):
    try:
        cmd = "organization retrieve --range {} {} {}".format(
                    START, START, org_id
                )
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Establishes relationship between part and organization 
def add_part_to_organization(uuid, part_uuid, private_key, public_key):
    try:
        cmd = "organization AddPart {} {} {} {}".format(
                    uuid, part_uuid, private_key, public_key
                )
        cmd = shlex.split(cmd)
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
                output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp
################################################################################
#                                    PART                                      #
################################################################################
@app.route("/ledger/api/v1/parts", methods=["POST"])
def create_part():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        uuid = request.json["part"]["uuid"]
        name = request.json["part"]["name"]
        
        checksum = request.json["part"]["checksum"]
        
        version = request.json["part"]["version"]
        
        alias = request.json["part"]["alias"]
        
        licensing = request.json["part"]["licensing"]
        
        label = request.json["part"]["label"]
        
        description = request.json["part"]["description"]
        
        cmd = "pt create {} {} {} {} {} {} {} {} {} {}".format(
                    uuid, name, checksum, version, alias, licensing, label,
                    description, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
            exp = ret_exception_msg(e) 
            return exp

@app.route("/ledger/api/v1/parts/amend", methods=["POST"])
def amend_part():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
        
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        uuid = request.json["part"]["uuid"]
        
        name = nullCast(request.json["part"], "name")
        
        checksum = nullCast(request.json["part"], "checksum")
        
        version = nullCast(request.json["part"], "version")
        
        alias = nullCast(request.json["part"], "alias")
        
        licensing = nullCast(request.json["part"], "licensing")
        
        label = nullCast(request.json["part"], "label")
        
        description = nullCast(request.json["part"], "description")
        
        cmd = "pt amend {} {} {} {} {} {} {} {} {} {}".format(
                    uuid, name, checksum, version, alias, licensing, label,
                    description, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
            exp = ret_exception_msg(e) 
            return exp

# Retrieves list of all parts
@app.route("/ledger/api/v1/parts", methods=["GET"])
def get_parts():
    try:
        cmd = "pt list-part"
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Retrieves part record by part id
@app.route("/ledger/api/v1/parts/<string:part_id>", methods=["GET"])
def get_part(part_id):
    try:    
        cmd = "pt retrieve " + part_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e: 
        exp = ret_exception_msg(e)
        return exp

# Retrieves historical part record by part id
@app.route("/ledger/api/v1/parts/history/<string:part_id>", methods=["GET"])
def get_uuid_part_history(part_id):
    try:    
        cmd = "pt retrieve --all " + part_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e: 
        exp = ret_exception_msg(e)
        return exp

# Retrieves historical part record on certain date by part id
@app.route(
    "/ledger/api/v1/parts/<string:part_id>/date/<string:START>",
    methods=["GET"]
)
def get_uuid_part_day(part_id, START):
    try:    
        cmd = "pt retrieve --range {} {} {}".format(START, START, part_id)
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        output = ""
        for line in process.stdout:
            output += line.decode("utf-8").strip()
            output = refine_output(output)
        return output
    except Exception as e: 
        exp = ret_exception_msg(e)
        return exp

# Establishes relationship between envelopes and part
@app.route("/ledger/api/v1/artifacts/part", methods=["POST"])
def add_artifact_to_part():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json or
            "relation" not in request.json):
            return  ret_exception_msg("Invalid JSON")
        output = ""
        
        uuid = request.json["relation"]["part_uuid"]
        envelope_uuid = request.json["relation"]["artifact_uuid"]
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        # cmd = "pt AddArtifact " + uuid + sp + envelope_uuid + sp + private_key + sp + public_key
        cmd = "pt AddArtifact {} {} {} {}".format(
                    uuid, envelope_uuid, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        for line in process.stdout:
            output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Establishes relationship between organization and part
@app.route("/ledger/api/v1/parts/orgs", methods=["POST"])
def add_organization_to_part():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json or
            "relation" not in request.json):
            return ret_exception_msg("Invalid JSON")
        output = ""
        
        uuid = request.json["relation"]["part_uuid"]
        organization_uuid = request.json["relation"]["organization_uuid"]
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        # cmd = "pt AddSupplier " + uuid + " " + organization_uuid + " "+ private_key + " "+ public_key
        cmd = "pt AddSupplier {} {} {} {}".format(
                    uuid, organization_uuid, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        
        add_part_to_organization(organization_uuid, uuid, private_key, 
                                    public_key)

        for line in process.stdout:
                output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

# Establishes relationship between category and part
@app.route("/api/sparts/ledger/parts/AddCategory", methods=["POST"])
def add_category_to_part():
    try:
        if (not request.json or
            "private_key" not in request.json or
            "public_key" not in request.json):
            return "Invalid JSON"
        output = ""
        
        uuid = request.json["add_category"]["part_uuid"]
        category_uuid = request.json["add_category"]["category_uuid"]
        public_key = request.json["public_key"]
        private_key = request.json["private_key"]
        
        # cmd = "pt AddCategory " + uuid + " " + category_uuid + " " + private_key + " "+ public_key
        cmd = "pt AddCategory {} {} {} {}".format(
                    uuid, category_uuid, private_key, public_key
                )
        cmd = shlex.split(cmd)
        
        process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
        process.wait()
        
        for line in process.stdout:
                output += line.decode("utf-8").strip()
        return output
    except Exception as e:
        exp = ret_exception_msg(e) 
        return exp

@app.route(
    "/api/sparts/ledger/envelopes/searchbychecksum/<string:checksum_id>",
    methods=["GET"]
)
def artifact_verify_checksum(checksum_id):
    try:
        artifactlist = get_envelopes()
        jdata = json.loads(artifactlist)
    
        output = ""
        for i in jdata:
            if i["checksum"] == checksum_id:
                output = json.dumps(i) 
        return output
    except Exception as e:
        exp = ret_exception_msg(e)
        return exp

@app.route(
    "/api/sparts/ledger/parts/artifact/<string:part_id>",
    methods=["GET"]
)
def get_part_artifact(part_id):
    try:
        cmd = "pt retrieve " + part_id
        process = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE)
        process.wait()
        partsOut = ""
        output = ""
        for line in process.stdout:
            partsOut += line.decode("utf-8").strip()
            partsOut = refine_output(partsOut)
    
        artifactListJSON = "[]"
        jArtifactsData = json.loads(artifactListJSON)
        data = json.loads(partsOut)
        for element in data["envelopes"]:
            artifactD = get_envelope(element["envelope_id"])
            artJSONObj = json.loads(artifactD)
            del artJSONObj["sub_artifact"]
            jArtifactsData.append(artJSONObj)
                
        output = json.dumps(jArtifactsData)
        return output
    except Exception as e:
        exp = ret_exception_msg(e)
        return exp
################################################################################
#                             SAWTOOTH VERSION                                 #
################################################################################
@app.route("/ledger/api/v1/parts/sawtooth/version", methods=["GET"])
def get_sawtooth_version():
    output = "{\"name\":\"Hyperledger Sawtooth\",\"version\":\"0.8.8\"}"
    return output
################################################################################
#                            PRIVATE FUNCTIONS                                 #
################################################################################
def ret_success_auth_msg(input):
    if "success" in input:
        msgJSON = "{}"
        key = json.loads(msgJSON)
        key["status"] = "success"
        key["message"] = "authorized"
        msgJSON = json.dumps(key)
        return msgJSON 
    else:
        return input

@app.route("/api/sparts/ledger/auth", methods=["POST"])
def sparts_auth():
    try:
        if (not request.json or
            "privatekey" not in request.json or
            "publickey" not in request.json or
            "allowedrole" not in request.json):
            return ret_exception_msg("Invalid JSON")
        output = ""
        
        uuid = get_uuid()
        private_key = request.json["privatekey"]
        public_key = request.json["publickey"]

        if len(private_key) == 64 and len(public_key) == 66:
            
            signature = get_signature(uuid, private_key, "wif")
            verify = verify_signature(uuid, signature, public_key)
            print("Verify Sign: " + str(verify)) 
            if str(verify) == "True":
                cmd = "user retrieve " + public_key
                cmd = shlex.split(cmd)
                process = subprocess.Popen(cmd, stdout=subprocess.PIPE)
                process.wait()
                ojson = ""
                for line in process.stdout:
                    ojson += line.decode("utf-8").strip()
               
                userinfo = json.loads(ojson)
                
                if userinfo.get("user_name"):
                    authorized = userinfo["authorized"]
                    
                    assignedrole = userinfo["role"]
                    if str(authorized) == "allow":
                        for i in request.json["allowedrole"]:
                            if str(i["role"]) == str(assignedrole):          
                                output = ret_success_auth_msg("success")
                                return output 
                            
                        output = ret_access_denied__msg(
                                        "Unauthorized, access is denied"
                                    )
                    elif str(authorized) == "deny":
                        output = ret_access_denied__msg(
                                        "Unauthorized, access is denied"
                                    )
                elif userinfo.get("message"):
                    output = ret_access_denied__msg(
                                    "Unauthorized, access is denied"
                                )    
            else:
                output = ret_access_denied__msg(
                                "Invalid keys, access is denied"
                            )
        else:
            output = ret_access_denied__msg(
                            "Invalid keys, access is denied"
                        )
        return output
    
    except Exception as e:
        exp = ret_access_denied__msg(str(e))
        return exp

@app.route("/ledger/api/v1/keys", methods=["GET"]) 
def get_keys():
    context = create_context("secp256k1")
    privkey = context.new_random_private_key()
    pubkey = context.get_public_key(privkey)
    userKeyJSON = "{}"
    keys = json.loads(userKeyJSON)
    keys["public_key"] = pubkey.as_hex()
    keys["private_key"] = privkey.as_hex()
    return ret_msg("success", "OK", "Keys", keys)
    
def ret_msg(status, message, result_type, result):
    msgJSON = "{}"
    key = json.loads(msgJSON)
    key["status"] = status
    key["message"] = message
    key["result_type"] = result_type
    key["result"] = result
    msgJSON = json.dumps(key)
    return msgJSON
    
def ret_exception_msg(message):
    msg = str(message)
    expJson = "{}"
    key = json.loads(expJson)
    key["status"] = "failed"
    key["message"] = msg
    key["result_type"] = "EmptyRecord"
    key["result"] = "{}"
    expJson = json.dumps(key)
    return expJson 

@app.errorhandler(401)
def ret_access_denied__msg(message):
    expJson = "{}"
    key = json.loads(expJson)
    key["status"] = "failed"
    key["message"] = message
    expJson = json.dumps(key)
    return expJson

def get_signature(message, private_key, privkey_format="wif"):
    context = create_context("secp256k1")
    factory = CryptoFactory(context)
    
    privkey = Secp256k1PrivateKey.from_hex(private_key)  
    signer = factory.new_signer(privkey)
    signature = signer.sign(message.encode())    
    return signature

def verify_signature(message, signature, public_key):
    try:
        context = create_context("secp256k1")
        pubkey = Secp256k1PublicKey.from_hex(public_key)
        result = context.verify(signature,message.encode(),pubkey)
        return result 
    except Exception:
        return False

def ret_auth_msg(status, message, auth, role):
    expJson = "{}"
    key = json.loads(expJson)
    key["status"] = status
    key["message"] = message
    key["authorized"] = auth;
    key["role"] = role
    expJson = json.dumps(key)
    return expJson 

def str_qt(input):
    output = "'" + input + "'"
    return output
    
def get_uuid():
    return str(uuid.uuid4())

def not_found():
    status = "{\"error\":\"Not found\"}"
    return status

def format_str(inputstr):
    output = "\"{}\"".format(inputstr)
    return output 

def refine_output(inputstr):
    ins = inputstr
    if "##" in ins:
        ins = ins.replace("##", ",")
    return ins
    
def nullCast(dic, key):
    if key not in dic:
        return "null"
    else:
        return dic[key]
################################################################################
#                            API to API ARTIFACT                               #
################################################################################
################################################################################
#                            API to API CATEGORY                               #
################################################################################
# API : CATEGORY CREATE 
@app.route("/phyo/api/create/category", methods=["POST"])
def api_create_category():
    
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:850/tp/category", 
                    data=json.dumps(request.json), headers=headers)
    output = response.content.decode("utf-8")
    
    return output

# API : CATEGORY AMEND
@app.route("/phyo/api/amend/category", methods=["POST"])
def api_amend_category():
    
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:850/tp/category/amend", 
                    data=json.dumps(request.json), headers=headers)
    output = response.content.decode("utf-8")
    
    return output

# API : CATEGORY LIST CATEGORY
@app.route("/phyo/api/list/category", methods=["GET"])
def api_list_category():
    response = requests.get("http://127.0.0.1:850/tp/category")
    output = response.content.decode("utf-8").strip()
    
    return output

# API : CATEGORY RETRIEVE {UUID}
@app.route("/phyo/api/retrieve/category/<string:category_id>", methods=["GET"])
def api_retrieve_category(category_id):
    response = requests.get(
                    "http://127.0.0.1:850/tp/category/{}".format(category_id)
                )
    output = response.content.decode("utf-8").strip()
    
    return output
    
# API : CATEGORY RETRIEVE --ALL {UUID}
@app.route(
    "/phyo/api/retrieve/category/history/<string:category_id>",
    methods=["GET"]
)
def api_retrieve_category_history(category_id):
    response = requests.get(
                    "http://127.0.0.1:850/tp/category/history/{}" \
                    .format(category_id)
                )
    output = response.content.decode("utf-8").strip()
    
    return output

# API : CATEGORY RETRIEVE --RANGE START END {UUID}
@app.route(
    "/phyo/api/retrieve/category/<string:category_id>/date/<string:START>",
    methods=["GET"]
)
def api_category_history_date(category_id, START):
    response = requests.get(
                    "http://127.0.0.1:850/tp/category/{}/date/{}" \
                    .format(category_id, START)
                )
    output = response.content.decode("utf-8").strip()
    
    return output
################################################################################
#                          API to API ORGANIZATION                             #
################################################################################
# API : ORGANIZATION CREATE 
@app.route("/phyo/api/create/organization", methods=["POST"])
def api_create_organization():
    
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:851/tp/organization", 
                    data=json.dumps(request.json), headers=headers)
    output = response.content.decode("utf-8")
    
    return output

# API : ORGANIZATION AMEND
@app.route("/phyo/api/amend/organization", methods=["POST"])
def api_amend_organization():
    
    headers = {"content-type": "application/json"}
    response = requests.post("http://127.0.0.1:851/tp/organization/amend", 
                    data=json.dumps(request.json), headers=headers)
    output = response.content.decode("utf-8")
    
    return output

# API : ORGANIZATION LIST ORGANIZATION
@app.route("/phyo/api/list/organization", methods=["GET"])
def api_list_organization():
    response = requests.get("http://127.0.0.1:851/tp/organization")
    output = response.content.decode("utf-8").strip()
    
    return output

# API : ORGANIZATION RETRIEVE {UUID}
@app.route("/phyo/api/retrieve/organization/<string:org_id>", methods=["GET"])
def api_retrieve_organization(org_id):
    response = requests.get(
                    "http://127.0.0.1:851/tp/organization/{}" \
                        .format(org_id)
                )
    output = response.content.decode("utf-8").strip()
    
    return output
    
# API : ORGANIZATION RETRIEVE --ALL {UUID}
@app.route(
    "/phyo/api/retrieve/organization/history/<string:org_id>",
    methods=["GET"]
)
def api_retrieve_organization_history(org_id):
    response = requests.get(
                    "http://127.0.0.1:851/tp/organization/history/{}" \
                    .format(org_id)
                )
    output = response.content.decode("utf-8").strip()
    
    return output

# API : ORGANIZATION RETRIEVE --RANGE START END {UUID}
@app.route(
    "/phyo/api/retrieve/organization/<string:org_id>/date/<string:START>",
    methods=["GET"]
)
def api_organization_history_date(org_id, START):
    response = requests.get(
                    "http://127.0.0.1:851/tp/organization/{}/date/{}" \
                    .format(org_id, START)
                )
    output = response.content.decode("utf-8").strip()
    
    return output
################################################################################
#                              API to API PART                                 #
################################################################################
################################################################################
#                            API to API RELATION                               #
################################################################################
################################################################################
#                                   MAIN                                       #
################################################################################
if __name__ == "__main__":
    app.run(host="0.0.0.0", port="818")
################################################################################
#                                                                              #
################################################################################

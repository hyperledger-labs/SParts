#!/bin/bash

 # NOTICE:
 # =======
 #  Copyright (c) 2018 Wind River Systems, Inc.
 #
 # Licensed under the Apache License, Version 2.0 (the "License");
 # you may not use this file except in compliance with the License.
 # You may obtain a copy of the License at:
 #       http://www.apache.org/licenses/LICENSE-2.0
 #
 # Unless required by applicable law or agreed to in writing, software  distributed
 # under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
 # OR CONDITIONS OF ANY KIND, either express or implied.
 #

# Set the bin directory to someplace within your $PATH (otherwise build in local directory)
BIN_DIR="."
go build -o $BIN_DIR/sparts
cp sparts spa

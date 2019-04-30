# ------------------------------------------------------------------------------
# Copyright 2017 Intel Corporation
# Copyright 2018 Wind River Systems
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

import os
import subprocess

from setuptools import setup, find_packages

setup(name='sparts-organization-family',
      version='1.0',
      description='Sparts Organization Family',
      author='Sameer Ahmed, Wind River System',
      url='https://github.com/hyperledger-labs/SParts',
      packages=find_packages(),
      install_requires=[
          'aiohttp',
          'colorlog',
          'protobuf',
          'sawtooth-sdk',
          'sawtooth-signing',
          'PyYAML',
          ],
      entry_points={
          'console_scripts': [
              'organization = sparts_organization.organization_cli:main_wrapper',
              'tp_organization = sparts_organization.processor.main:main',
          ]
      })

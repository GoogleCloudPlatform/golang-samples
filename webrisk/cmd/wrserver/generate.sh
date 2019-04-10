#!/bin/bash
# Copyright 2016 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
set -e

# This script builds the generated Go code for the static web files.
# The statik tool must be installed. The recommended version is:
#
#	github.com/rakyll/statik: 2940084503a48359b41de178874e862c5bc3efe8
for TOOL in statik; do
	command -v $TOOL >/dev/null 2>&1 || { echo "Could not locate $TOOL. Aborting." >&2; exit 1; }
done

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

statik -src public -dest .

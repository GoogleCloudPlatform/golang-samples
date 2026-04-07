/*
 Copyright 2025 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

import fs from "node:fs";
import path from "node:path";
import setupVars from "../setup-vars.js";

const project_id = process.env.PROJECT_ID;
if (!project_id) {
  console.error(
    "Please set the PROJECT_ID environment variable to your Google Cloud project."
  );
  process.exit(1);
}

const core = {
  exportVariable: (_key, _value) => null,
};

const setupFile = process.argv[2];
if (!setupFile) {
  console.error("Please provide the path to a setup file.");
  process.exit(1);
}
const data = fs.readFileSync(path.join("..", "..", setupFile), "utf8");
const setup = JSON.parse(data);

setupVars({ project_id, core, setup });

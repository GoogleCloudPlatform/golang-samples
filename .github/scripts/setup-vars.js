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

export default function setupVars(
  { projectId, core, setup, serviceAccount, idToken },
  runId = null
) {
  // Define automatic variables plus custom variables.
  const vars = {
    PROJECT_ID: projectId,
    RUN_ID: runId || uniqueId(),
    SERVICE_ACCOUNT: serviceAccount,
    ...(setup.env || {}),
  };

  // Apply variable interpolation.
  const env = Object.fromEntries(
    Object.keys(vars).map((key) => [key, substituteVars(vars[key], vars)])
  );

  // Export environment variables.
  console.log("env:");
  for (const key in env) {
    const value = env[key];
    console.log(`  ${key}: ${value}`);
    core.exportVariable(key, value);
  }

  // Show exported secrets, for logging purposes.
  // TODO: We might want to fetch the secrets here and export them directly.
  //       https://cloud.google.com/secret-manager/docs/create-secret-quickstart#secretmanager-quickstart-nodejs
  console.log("secrets:");
  for (const key in setup.secrets || {}) {
    // This is the Google Cloud Secret Manager secret ID.
    // NOT the secret value, so it's ok to show.
    console.log(`  ${key}: ${setup.secrets[key]}`);
  }

  // Set global secret for the Service Account identity token
  // Use in place of 'gcloud auth print-identity-token' or auth.getIdTokenClient
  // usage: curl -H 'Bearer: $ID_TOKEN' https://
  core.exportVariable("ID_TOKEN", idToken);
  core.setSecret(idToken);
  // For logging, show the source of the ID_TOKEN
  console.log(`  ID_TOKEN: steps.auth.outputs.id_token (from GitHub Action)`);

  // Return env and secrets to use for further steps.
  return {
    env: env,
    // Transform secrets into the format needed for the GHA secret manager step.
    secrets: Object.keys(setup.secrets || {})
      .map((key) => `${key}:${setup.secrets[key]}`)
      .join("\n"),
  };
}

export function substituteVars(value, env) {
  for (const key in env) {
    let re = new RegExp(`\\$(${key}\\b|\\{\\s*${key}\\s*\\})`, "g");
    value = value.replaceAll(re, env[key]);
  }
  return value;
}

export function uniqueId(length = 6) {
  const min = 2 ** 32;
  const max = 2 ** 64;
  return Math.floor(Math.random() * max + min)
    .toString(36)
    .slice(0, length);
}

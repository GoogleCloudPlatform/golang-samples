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

import { deepStrictEqual } from "assert";
import setupVars from "./setup-vars.js";
import { substituteVars, uniqueId } from "./setup-vars.js";

const projectId = "my-test-project";
const serviceAccount = "my-sa@my-project.iam.gserviceaccount.com";
const core = {
  exportVariable: (_key, _value) => null,
  setSecret: (_key) => null,
};

const autovars = {
  PROJECT_ID: projectId,
  RUN_ID: "run-id",
  SERVICE_ACCOUNT: serviceAccount,
};

describe("setupVars", () => {
  describe("env", () => {
    it("empty", () => {
      const setup = {};
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = autovars;
      deepStrictEqual(vars.env, expected);
    });

    it("zero vars", () => {
      const setup = { env: {} };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = autovars;
      deepStrictEqual(vars.env, expected);
    });

    it("one var", () => {
      const setup = { env: { A: "x" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = { ...autovars, A: "x" };
      deepStrictEqual(vars.env, expected);
    });

    it("three vars", () => {
      const setup = { env: { A: "x", B: "y", C: "z" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = { ...autovars, A: "x", B: "y", C: "z" };
      deepStrictEqual(vars.env, expected);
    });

    it("should override automatic variables", () => {
      const setup = {
        env: { PROJECT_ID: "custom-value", SERVICE_ACCOUNT: "baz@foo.com" },
      };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = {
        PROJECT_ID: "custom-value",
        RUN_ID: "run-id",
        SERVICE_ACCOUNT: "baz@foo.com",
      };
      deepStrictEqual(vars.env, expected);
    });

    it("should interpolate variables", () => {
      const setup = { env: { A: "x", B: "y", C: "$A/${B}" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = { ...autovars, A: "x", B: "y", C: "x/y" };
      deepStrictEqual(vars.env, expected);
    });

    it("should not interpolate secrets", () => {
      const setup = {
        env: { C: "$x/$y" },
        secrets: { A: "x", B: "y" },
      };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = { ...autovars, C: "$x/$y" };
      deepStrictEqual(vars.env, expected);
    });
  });

  describe("secrets", () => {
    it("zero secrets", () => {
      const setup = { secrets: {} };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      deepStrictEqual(vars.secrets, "");
    });

    it("one secret", () => {
      const setup = { secrets: { A: "x" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = "A:x";
      deepStrictEqual(vars.secrets, expected);
    });

    it("three secrets", () => {
      const setup = { secrets: { A: "x", B: "y", C: "z" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = "A:x\nB:y\nC:z";
      deepStrictEqual(vars.secrets, expected);
    });

    it("should not interpolate variables", () => {
      const setup = {
        env: { A: "x", B: "y" },
        secrets: { C: "$A/$B" },
      };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = "C:$A/$B";
      deepStrictEqual(vars.secrets, expected);
    });

    it("should not interpolate secrets", () => {
      const setup = { secrets: { A: "x", B: "y", C: "$A/$B" } };
      const vars = setupVars(
        { projectId, core, setup, serviceAccount },
        "run-id"
      );
      const expected = "A:x\nB:y\nC:$A/$B";
      deepStrictEqual(vars.secrets, expected);
    });
  });
});

describe("substituteVars", () => {
  it("should interpolate $VAR", () => {
    const got = substituteVars("$A-$B", { A: "x", B: "y" });
    const expected = "x-y";
    deepStrictEqual(got, expected);
  });

  it("should interpolate ${VAR}", () => {
    const got = substituteVars("${A}-${B}", { A: "x", B: "y" });
    const expected = "x-y";
    deepStrictEqual(got, expected);
  });

  it("should interpolate ${ VAR }", () => {
    const got = substituteVars("${ A }-${ \tB\t }", { A: "x", B: "y" });
    const expected = "x-y";
    deepStrictEqual(got, expected);
  });

  it("should not interpolate on non-word boundary", () => {
    const got = substituteVars("$Ab", { A: "x" });
    const expected = "$Ab";
    deepStrictEqual(got, expected);
  });
});

describe("uniqueId", () => {
  it("should match length", () => {
    const n = 6;
    deepStrictEqual(uniqueId(n).length, n);
  });
});

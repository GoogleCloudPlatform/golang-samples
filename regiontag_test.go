// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package samples

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestRegionTags(t *testing.T) {
	type RegionTag struct {
		tag  string
		path string
	}
	errs := []RegionTag{}

	cmd := exec.Command("go", "list", "./...")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		path := scanner.Text()

		goDoc := exec.Command("go", "doc", path)
		grep := exec.Command("grep", "START")
		grep.Stdin, _ = goDoc.StdoutPipe()

		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)
		grep.Stdout = writer

		grep.Start()
		goDoc.Run()
		grep.Wait()

		writer.Flush()
		s := buf.String()
		if len(s) > 0 {
			errs = append(errs, RegionTag{s, path})
		}
	}

	if len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("%#v", e.path)
		}
	}

}

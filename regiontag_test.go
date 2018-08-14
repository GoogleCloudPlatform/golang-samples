// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package samples

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

func listPackages() <-chan string {
	c := make(chan string)
	go func() {
		cmd := exec.Command("go", "list", "./...")
		out, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		close(c)
	}()
	return c
}

type RegionTag struct {
	file string
	line string
}

func findRegionTags(pkgs <-chan string) <-chan RegionTag {
	c := make(chan RegionTag)
	go func() {
		for p := range pkgs {
			goDoc := exec.Command("go", "doc", p)
			out, err := goDoc.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}

			if err = goDoc.Start(); err != nil {
				log.Fatal(err)
			}

			scanner := bufio.NewScanner(out)

			// filter affected lines only
			for scanner.Scan() {
				text := scanner.Text()
				if !strings.Contains(text, "[START") {
					continue
				}
				c <- RegionTag{file: p, line: text}
			}
		}
		close(c)
	}()
	return c
}

func merge(cs ...<-chan RegionTag) <-chan RegionTag {
	var wg sync.WaitGroup
	out := make(chan RegionTag)

	output := func(c <-chan RegionTag) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func TestRegionTags(t *testing.T) {

	in := listPackages()
	const workers = 4

	var cs [workers]<-chan RegionTag
	for i := 0; i < workers; i++ {
		cs[i] = findRegionTags(in)
	}

	for tag := range merge(cs[:]...) {
		t.Errorf("\nFile: %v\nLine: %#v", tag.file, tag.line)
	}
}

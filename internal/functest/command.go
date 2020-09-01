package functest

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// runCommand executes shell scripts.
// For logging purposes, it assumes commands take the form: '/bin/sh -c "my-command"'
func runCommand(cmd *exec.Cmd, dir string) ([]byte, error) {
	cmd.Dir = dir
	log.Printf("Executing: %s [dir: %s]\n", strings.Join(cmd.Args[2:], " "), cmd.Dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("----| Error Output |----")
		os.Stderr.Write(out)
		return []byte{}, fmt.Errorf("Command.CombinedOutput: %w", err)
	}

	return bytes.TrimSpace(out), err
}

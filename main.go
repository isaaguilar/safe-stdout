package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type safeStdout struct {
	secrets []string
}

func NewSafeStdout(secrets []string) safeStdout {
	return safeStdout{
		secrets: secrets,
	}
}

func (s safeStdout) Write(b []byte) (int, error) {
	var line bytes.Buffer
	i := 0
	for {
		if i == len(b) {
			break
		}
		fmt.Fprintf(&line, "%s", string(b[i]))
		i++
	}

	safeString := line.String()
	for _, secret := range s.secrets {
		safeString = strings.ReplaceAll(safeString, secret, "*****")
	}

	fmt.Printf("%s", safeString)
	return i, nil
}

func main() {
	cmd := exec.Command("/bin/bash", "scripts/test.sh")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	writer := NewSafeStdout([]string{"something", "sleep", "seconds", "be"})

	cmd.Start()
	io.Copy(writer, stdout)
	cmd.Wait()
}

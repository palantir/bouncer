// Copyright 2017 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bouncer

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func splitCommandString(fullCommand string) (string, []string) {
	commandArr := strings.Split(fullCommand, " ")
	command := commandArr[0]
	var args []string
	if len(commandArr) > 1 {
		args = commandArr[1:]
	}
	return command, args
}

func bufferResults(cmd *exec.Cmd, r io.Reader, inputType string) {
	logger := log.WithFields(log.Fields{
		"Output Source Cmd": cmd.Args,
		"Output Source":     inputType,
	})
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if inputType == "stdout" {
			logger.Info(scanner.Text())
		} else {
			logger.Warn(scanner.Text())
		}
	}
	// TODO investigate how to ensure we never get a "file already closed" error here, where the cmd
	// has exited and closed while this function is still trying to run.  Maybe we just need to test
	// for that specific scanner error?  Something else?
	// err := scanner.Err()
	// if err != nil {
	// 	log.Error(errors.Wrapf(err, "error scanning %s", inputType))
	// }
}

func getCmd(command string, args []string) (*exec.Cmd, error) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "error attaching stdout")
	}
	go bufferResults(cmd, stdout, "stdout")

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "error attaching stderr")
	}
	go bufferResults(cmd, stderr, "stderr")

	return cmd, nil
}

func (r *BaseRunner) executeExternalCommand(fullCommand string) error {
	tmout := r.opts.ItemTimeout
	command, args := splitCommandString(fullCommand)
	log.Infof("Executing pre-terminate command '%s' with args '%s'", command, args)
	r.resetTimeout()
	r.noopCheck()

	cmd, err := getCmd(command, args)
	if err != nil {
		return errors.Wrap(err, "error initializing cmd")
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "error starting command")
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(tmout):
		err = cmd.Process.Kill()
		if err != nil {
			return errors.Wrapf(err, "error killing process after timeout of %s", tmout)
		}
		return errors.Errorf("process killed as timeout of %s reached", tmout)
	case err = <-done:
		return errors.Wrap(err, "error running process")
	}
}

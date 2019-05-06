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
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// DesiredASG contains pieces of the ASG as they _should_ be, but at any given time, since we twiddle
// the desired capacity, may not _actually_ be.
type DesiredASG struct {
	AsgName         string
	DesiredCapacity int64
	// PreTerminateCmd is the external process that needs to be run before terminating an instance in this ASG
	PreTerminateCmd *string
}

// extractDesiredASG takes in a separator-separated string of asgname and desired capacity, and returns a DesiredASG pointer
func extractDesiredASG(asgItem string, separator string, defaultDesired *int64, preTerminateCmd *string) (*DesiredASG, error) {
	asgItems := strings.Split(asgItem, separator)
	var desiredCapacity int64
	var err error

	if len(asgItems) > 2 || (defaultDesired == nil && len(asgItems) == 1) {
		return nil, errors.Errorf("Error parsing '%s'.  Must be in format '%s%s%s'", asgItem, "ASG-NAME", separator, "1")
	} else if len(asgItems) == 2 {
		desiredCapacity, err = strconv.ParseInt(asgItems[1], 10, 64)
		if err != nil {
			return nil, errors.Errorf("Error parsing '%s' from ASG Item '%s' as int64", asgItems[1], asgItem)
		}
	} else {
		desiredCapacity = *defaultDesired
	}

	curASG := DesiredASG{
		AsgName:         asgItems[0],
		DesiredCapacity: desiredCapacity,
		PreTerminateCmd: preTerminateCmd,
	}

	return &curASG, nil
}

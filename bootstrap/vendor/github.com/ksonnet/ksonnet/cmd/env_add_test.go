// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"testing"

	"github.com/ksonnet/ksonnet/actions"
)

func Test_envAddCmd(t *testing.T) {
	cases := []cmdTestCase{
		{
			name:   "in general",
			args:   []string{"env", "add", "prod", "--server", "http://example.com"},
			action: actionEnvAdd,
			expected: map[string]interface{}{
				actions.OptionApp:      ka,
				actions.OptionEnvName:  "prod",
				actions.OptionModule:   "default",
				actions.OptionOverride: false,
				actions.OptionServer:   "http://example.com",
				actions.OptionSpecFlag: "version:v1.7.0",
			},
		},
	}

	runTestCmd(t, cases)
}

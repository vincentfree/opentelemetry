// Copyright 2024 Vincent Free
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package providerconfighttp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEndpoint(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input string
		error error
	}{
		"without protocol": {
			input: "randomurl",
			error: fmt.Errorf("invalid endpoint: randomurl"),
		},
		"with http protocol": {
			input: "http://localhost:1234",
			error: nil,
		},
		"with https protocol": {
			input: "https://localhost:12",
			error: nil,
		},
		"port bigger than max value": {
			input: "https://localhost:65536",
			error: fmt.Errorf("invalid port value: 65536"),
		},
	}
	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			// Act
			err := validateEndpoint(testData.input)

			// Assert
			if testData.error == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, testData.error, err)
			}
		})
	}
}

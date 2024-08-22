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

package providerconfig

import "slices"

type Execution uint8

//sync async

const (
	Sync Execution = iota + 1
	Async
)

var (
	executionTypes = []Execution{Sync, Async}
)

func (e Execution) String() string {
	switch e {
	case Sync:
		return "Sync"
	case Async:
		return "Async"
	default:
		return "undefined"
	}
}

func (e Execution) IsValid() bool {
	return slices.Contains(executionTypes, e)
}

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

type Protocol uint8

var (
	protocols = []Protocol{Grpc, Http}
)

const (
	Grpc Protocol = iota + 1
	Http
)

func (p Protocol) String() string {
	switch p {
	case Grpc:
		return "grpc"
	case Http:
		return "http"
	default:
		return "undefined"
	}
}

func (p Protocol) Port() int {
	switch p {
	case Grpc:
		return grpcPort
	case Http:
		return httpPort
	default:
		return 0
	}
}

func (p Protocol) IsValid() bool {
	return slices.Contains(protocols, p)
}

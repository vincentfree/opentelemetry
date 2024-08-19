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

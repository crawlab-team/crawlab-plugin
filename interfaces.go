package plugin

import "github.com/crawlab-team/crawlab-core/interfaces"

type Plugin interface {
	Register()
	GetClient() interfaces.GrpcClient
}

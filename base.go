package plugin

import "github.com/crawlab-team/crawlab-core/interfaces"

type BasePlugin struct {
	c interfaces.GrpcClient
}

func (p *BasePlugin) Init() {
}

func (p *BasePlugin) GetClient() interfaces.GrpcClient {
	return p.c
}

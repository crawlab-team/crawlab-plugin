package plugin

import "github.com/crawlab-team/crawlab-core/interfaces"

type BasePlugin struct {
	c interfaces.GrpcClient
}

func (p *BasePlugin) Register() {
}

func (p *BasePlugin) GetClient() interfaces.GrpcClient {
	return p.c
}

func NewPlugin() (p Plugin) {
	return &BasePlugin{}
}

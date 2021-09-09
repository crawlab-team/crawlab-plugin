package basic

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type Plugin struct {
}

func (p *Plugin) GetClient() interfaces.GrpcClient {
	panic("implement me")
}

func (p *Plugin) Init() {
	panic("implement me")
}

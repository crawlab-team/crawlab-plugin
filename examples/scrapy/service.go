package scrapy

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	plugin "github.com/crawlab-team/crawlab-plugin"
)

type Service struct {
	c interfaces.GrpcClient
}

func (svc *Service) Init() (err error) {
	return nil
}

func (svc *Service) Start() (err error) {
	return nil
}

func (svc *Service) GetClient() interfaces.GrpcClient {
	panic("implement me")
}

func NewScrapyPluginService() plugin.Plugin {
	return nil
}

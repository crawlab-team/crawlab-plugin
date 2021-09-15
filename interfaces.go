package plugin

import grpc "github.com/crawlab-team/crawlab-grpc"

type Plugin interface {
	Init() error
	Start() error
	Stop() error
}

type EventServiceInterface interface {
	Subscribe() (err error)
	GetStream() (stream grpc.PluginService_SubscribeClient)
}

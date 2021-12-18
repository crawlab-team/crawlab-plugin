package plugin

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
)

type EventService struct {
	internal *Internal
	stream   grpc.PluginService_SubscribeClient
}

func (svc *EventService) Subscribe() (err error) {
	log.Infof("subscribe events")

	// request request data
	data, err := json.Marshal(entity.GrpcEventServiceMessage{
		Type: constants.GrpcEventServiceTypeRegister,
	})
	if err != nil {
		return trace.TraceError(err)
	}

	// register request
	req := &grpc.PluginRequest{
		Name:    svc.internal.p.Name,
		NodeKey: viper.GetString("node.key"),
		Data:    data,
	}

	// register
	_, err = svc.internal.GetGrpcClient().GetPluginClient().Register(context.Background(), req)
	if err != nil {
		return trace.TraceError(err)
	}

	// stream
	svc.stream = svc.internal.GetGrpcClient().GetStream()

	return nil
}

func (svc *EventService) GetStream() (stream grpc.PluginService_SubscribeClient) {
	return svc.stream
}

func NewEventService(internal *Internal) (svc EventServiceInterface) {
	svc = &EventService{
		internal: internal,
	}
	return svc
}

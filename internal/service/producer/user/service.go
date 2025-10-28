package producer

import (
	"github.com/WithSoull/UserServer/internal/service"
	"github.com/WithSoull/platform_common/pkg/kafka"
	"google.golang.org/protobuf/encoding/protojson"
)

type userProducerService struct {
	userCreatedProducer kafka.Producer
	userDeletedProducer kafka.Producer
}

func NewService(
	userCreatedProducer kafka.Producer,
	userDeletedProducer kafka.Producer,
) service.UserProducerService {
	return &userProducerService{
		userCreatedProducer: userCreatedProducer,
		userDeletedProducer: userDeletedProducer,
	}
}

func PJ() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
		Indent:          "",
	}
}

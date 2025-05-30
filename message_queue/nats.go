package message_queue

import (
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	natsService *NatsService = nil
)

type NatsService struct {
	conn      *nats.Conn
	Url       string
	logger    runtime.Logger
	marshaler *protojson.MarshalOptions
}

func InitNatsService(logger runtime.Logger, natsUrl string, marshaler *protojson.MarshalOptions) {
	natsService = &NatsService{
		conn:      nil,
		Url:       natsUrl,
		logger:    logger,
		marshaler: marshaler,
	}
	natsService.Connect()
}

func GetNatsService() *NatsService {
	return natsService
}

func (conn *NatsService) Connect() {
	var err error
	conn.conn, err = nats.Connect(conn.Url)
	if err != nil {
		conn.logger.Error("Cannot connect to nats server %v", err)
	} else {
		conn.logger.Error("Connect to nats server success")
	}
}

func (conn *NatsService) Publish(topic string, data proto.Message) {
	dataByte, err := conn.marshaler.Marshal(data)
	if err != nil {
		conn.logger.Error("Publish topic Marshal data error %v", err)
		return
	}
	err = conn.conn.Publish(topic, dataByte)
	if err != nil {
		conn.logger.Error("Publish topic error %v", err)
	} else {
		conn.logger.Info("Publish topic success %v", topic)
	}
}

func (conn *NatsService) RegisterAllSubject() {
	for topic, _ := range messageHandler {
		_, err := conn.conn.Subscribe(topic, func(msg *nats.Msg) {
			processMessage(msg.Subject, msg.Data)
		})
		if err != nil {
			conn.logger.Error("Subscribe nats topic error %v", err)
		}
	}
}

func (conn *NatsService) Disconnect() {
	conn.conn.Close()
}

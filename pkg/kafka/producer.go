package kafka

import (
	"4microservice/api_gateway/api/handlers/models"
	"4microservice/api_gateway/config"
	pbc "4microservice/api_gateway/genproto/post_service"
	pb "4microservice/api_gateway/genproto/user_service"
	"4microservice/api_gateway/pkg/logger"
	"context"
	"encoding/json"
	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Produce struct {
	Log     logger.Logger
	User    *kafka.Writer
	Product *kafka.Writer
}

func NewProducerInit(conf config.Config, log logger.Logger) *Produce {
	return &Produce{
		Log: log,
		User: &kafka.Writer{
			Addr:                   kafka.TCP(conf.KafkaAddress),
			Topic:                  conf.KafkaTopic,
			AllowAutoTopicCreation: true,
		},
		Product: &kafka.Writer{
			Addr:                   kafka.TCP(conf.KafkaAddress),
			Topic:                  conf.KafkaTopic,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Produce) ProduceUser(ctx context.Context, key string, proto models.User) error {
	event := pb.User{
		Username:  proto.Username,
		Email:     proto.Email,
		Password:  proto.Password,
		FirstName: proto.FirstName,
		LastName:  proto.LastName,
		Bio:       proto.Bio,
		Website:   proto.Website,
	}
	byteData, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: byteData,
	}
	return p.User.WriteMessages(ctx, message)
}

func (p *Produce) ProducePost(ctx context.Context, key string, proto models.Post) error {
	event := pbc.Post{
		Id:       proto.Id,
		UserId:   proto.UserId,
		Content:  proto.Content,
		ImageUrl: proto.ImageUrl,
		Title:    proto.Title,
		Likes:    proto.Likes,
		Dislikes: proto.Dislikes,
		Views:    proto.Views,
	}

	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	return p.Product.WriteMessages(ctx, message)
}

func (p *Produce) Close() {
	if err := p.Product.Close(); err != nil {
		p.Log.Error("error while closing product create", zap.Error(err))
	}
	if err := p.User.Close(); err != nil {
		p.Log.Error("error while closing user create", zap.Error(err))
	}
}

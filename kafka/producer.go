package kafka

import (
	"4microservice/api_gateway/api/handlers/models"
	"context"
)

type ProduceMessages interface {
	ProducePost(ctx context.Context, key string, proto models.Post) error
	ProduceUser(ctx context.Context, key string, proto models.User) error
	Close()
}

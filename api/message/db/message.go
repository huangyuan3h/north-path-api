package db

import (
	"crypto/rand"
	"time"

	"api.north-path.site/message/types"
	db "api.north-path.site/utils/dynamodb"
	"github.com/oklog/ulid"
)

const tableName = "message"

type Message struct {
	types.Message
	client *db.Client
}

type UserMethod interface {
	CreateNew(message *types.Message) error
}

func New() UserMethod {
	client := db.New(tableName)

	return Message{client: &client}
}

func (m Message) CreateNew(message *types.Message) error {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	message.Id = id.String()
	message.CreatedDate = time.Now().Format(time.RFC3339)

	return m.client.CreateOrUpdate(message)
}

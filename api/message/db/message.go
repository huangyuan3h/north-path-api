package db

import (
	"crypto/rand"
	"time"

	"north-path.it-t.xyz/message/types"
	db "north-path.it-t.xyz/utils/dynamodb"
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

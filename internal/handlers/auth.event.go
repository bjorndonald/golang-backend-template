package handlers

import (
	"encoding/json"
	"fmt"

	// "log"

	"github.com/IBM/sarama"
	"github.com/bjorndonald/lasgcce/internal/bootstrap"
	"github.com/bjorndonald/lasgcce/internal/models"
)

type EventHandler struct {
	Deps *bootstrap.AppDependencies
}

func (h *EventHandler) ProcessSignup(msg *sarama.ConsumerMessage) error {
	// log.Printf("Processing message: topic=%s, partition=%d, offset=%d, key=%s, value=%s",
	// 	msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	user := &models.User{}
	if err := json.Unmarshal(msg.Value, user); err != nil {
		return err
	}

	fmt.Println("new user registration event received", user)

	// maybe send email notification

	return nil
}

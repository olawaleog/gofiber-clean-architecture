package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"gopkg.in/gomail.v2"
	"strconv"
)

// EmailConsumerService handles processing of queued email messages
type EmailConsumerService interface {
	// StartConsumer starts consuming email messages from the queue
	StartConsumer() error

	// ProcessEmail processes a single email message
	ProcessEmail(emailData []byte) error
}

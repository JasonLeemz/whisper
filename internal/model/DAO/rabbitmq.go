package dao

import (
	"github.com/streadway/amqp"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/mq"
)

func NewMQDao() *MQDao {
	return &MQDao{
		Conn: mq.Conn,
	}
}

type MQDao struct {
	Conn *amqp.Connection
}

type MQ interface {
	ProduceMessage(exchange, routingKey string, message []byte)
}

// ProduceMessage exchange: whisper, queneName: whisper.*
func (d *MQDao) ProduceMessage(exchange, routingKey string, message []byte) {
	ctx := context.NewContext()
	ch, err := d.Conn.Channel()
	if err != nil {
		log.Logger.Error(ctx, err)
		return
	}

	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		log.Logger.Error(ctx, "Failed to publish a message: %s", err.Error())
	} else {
		log.Logger.Info(ctx, routingKey, "ok")
	}
}

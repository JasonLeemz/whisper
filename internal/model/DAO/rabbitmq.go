package dao

import (
	"github.com/streadway/amqp"
	"sync"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/mq"
)

var (
	mqDao     *MQDao
	onceMQDao sync.Once
)

func NewMQDao() *MQDao {
	onceMQDao.Do(func() {
		mqDao = &MQDao{
			Conn: mq.Conn,
		}
	})
	return mqDao
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
	//ch.QueueDeclare(routingKey)
	err = ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // 托管。如果收到消息但是路由失败，false会丢弃，true会处理
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

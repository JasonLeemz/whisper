package mq

import (
	"github.com/streadway/amqp"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

type SearchBoxSvc struct {
}

func (s *SearchBoxSvc) ConsumeSearchBoxMsg() error {

	ctx := context.NewContext()
	var msgs <-chan amqp.Delivery

	mqd := dao.NewMQDao()
	ch, err := mqd.Conn.Channel()
	if err != nil {
		log.Logger.Error(ctx, err)
		return err
	}
	msgs, err = ch.Consume(
		QueneNameSearchKey, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)

	if err != nil {
		log.Logger.Error(ctx, "Failed to register a consumer")
		return err
	}

	go func() {
		for d := range msgs {
			log.Logger.Infof("Received a message: %s", d.Body)
			TopNZIncrBy(ctx, redis.KeyHotSearchSearchBox, d.Body)
		}
	}()

	log.Logger.Info(ctx, "Waiting for SearchKey messages...")
	return nil
}

// ConsumeSearchBoxMsgCommand 命令封装成对象
type ConsumeSearchBoxMsgCommand struct {
	*SearchBoxSvc
}

func NewConsumeSearchBoxMsgCMD(svc *SearchBoxSvc) *ConsumeSearchBoxMsgCommand {
	return &ConsumeSearchBoxMsgCommand{SearchBoxSvc: svc}
}

func (cmd ConsumeSearchBoxMsgCommand) Exec(args ...interface{}) error {
	return cmd.ConsumeSearchBoxMsg()
}

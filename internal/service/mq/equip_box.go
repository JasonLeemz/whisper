package mq

import (
	"github.com/streadway/amqp"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

type EquipBoxSvc struct {
}

func (s *EquipBoxSvc) ConsumeEquipBoxMsg() error {

	ctx := context.NewContext()
	var msgs <-chan amqp.Delivery

	mqd := dao.NewMQDao()
	ch, err := mqd.Conn.Channel()
	if err != nil {
		log.Logger.Error(ctx, err)
		return err
	}
	msgs, err = ch.Consume(
		QueneNameEquipBox, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	if err != nil {
		log.Logger.Error(ctx, "Failed to register a consumer")
		return err
	}

	go func() {
		for d := range msgs {
			log.Logger.Infof("Received a message: %s", d.Body)
			TopNZIncrBy(ctx, redis.KeyHotSearchEquipBox, d.Body)
		}
	}()

	log.Logger.Info(ctx, "Waiting for EquipBox messages...")
	return nil
}

// ConsumeEquipBoxMsgCommand 命令封装成对象
type ConsumeEquipBoxMsgCommand struct {
	*EquipBoxSvc
}

func NewConsumeEquipBoxMsgCMD(svc *EquipBoxSvc) *ConsumeEquipBoxMsgCommand {
	return &ConsumeEquipBoxMsgCommand{EquipBoxSvc: svc}
}

func (cmd ConsumeEquipBoxMsgCommand) Exec(args ...interface{}) error {
	return cmd.ConsumeEquipBoxMsg()
}

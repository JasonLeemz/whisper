package mq

import (
	"github.com/streadway/amqp"
	dao "whisper/internal/model/DAO"
	"whisper/pkg/context"
	"whisper/pkg/log"
	"whisper/pkg/redis"
)

func ConsumeSearchKeyMsg() error {
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

func ConsumeEquipBoxMsg() error {
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

func TopNZIncrBy(ctx *context.Context, key string, data []byte) {
	redis.RDB.ZIncrBy(ctx, key, 1, string(data))
}

type HandlerFunc func() error

var ConsumerFunc = []HandlerFunc{
	ConsumeSearchKeyMsg,
	ConsumeEquipBoxMsg,
}

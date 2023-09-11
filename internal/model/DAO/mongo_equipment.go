package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"whisper/internal/model"
	"whisper/pkg/context"
	mongo2 "whisper/pkg/mongo"
)

type MongoEquipment interface {
	Add(ctx *context.Context, ei []*model.EquipIntro) error
	Find(ctx *context.Context, cond map[string]interface{}) ([]*model.EquipIntro, error)
	Delete(ctx *context.Context, cond map[string]interface{}) error
}

type MongoEquipmentDAO struct {
	client     *mongo.Client
	db         *mongo.Database
	collection string
}

func (d *MongoEquipmentDAO) Find(ctx *context.Context, filter bson.M) ([]*model.EquipIntro, error) {

	opts := options.Find().SetSort(bson.D{{"price", -1}})
	cursor, err := d.db.Collection(d.collection).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var results []*model.EquipIntro
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (d *MongoEquipmentDAO) Add(ctx *context.Context, ei []*model.EquipIntro) error {
	// 创建要插入的文档
	data := make([]interface{}, 0, len(ei))
	for _, intro := range ei {
		bsonData, err := bson.Marshal(intro)
		if err != nil {
			return err
		}
		data = append(data, bsonData)
	}
	// 插入文档
	_, err := d.db.Collection(d.collection).InsertMany(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (d *MongoEquipmentDAO) Delete(ctx *context.Context, cond map[string]interface{}) error {
	_, err := d.db.Collection(d.collection).DeleteMany(ctx, cond)
	return err
}

var (
	mgDao  *MongoEquipmentDAO
	mgOnce *sync.Once
)

func NewMongoEquipmentDAO() *MongoEquipmentDAO {
	mgOnce.Do(func() {
		mgDao = &MongoEquipmentDAO{
			client:     mongo2.Client,
			db:         mongo2.Database,
			collection: "equipment",
		}
	})
	return mgDao
}

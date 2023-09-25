package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func InitMilvusClient(ctx context.Context, milvusAddr string) (client.Client, error) {

	c, err := client.NewClient(ctx, client.Config{
		Address: milvusAddr,
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

func CreateCollectionIfNotExist(
	c client.Client, ctx context.Context, collectionName string, schema *entity.Schema, resetCollection bool) (bool, error) {

	// first, lets check the collection exists
	collExists, err := c.HasCollection(ctx, collectionName)
	if err != nil {
		return false, err
	}

	if collExists {
		log.Println("Collection ", collectionName, " already exist.")
		if resetCollection {
			log.Println("Dropping ", collectionName)

			err = c.DropCollection(ctx, collectionName)
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	err = c.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func createDbIfNotExist(c client.Client, ctx context.Context, dbName string) (bool, error) {
	dbList, err := c.ListDatabases(ctx)
	if err != nil {
		return false, err
	}

	for _, db := range dbList {
		if db.Name == dbName {
			return true, nil
		}
	}

	err = c.CreateDatabase(ctx, dbName)
	if err != nil {
		return false, err
	}

	return true, nil
}

func InsertIntoCollection(c client.Client, ctx context.Context, collectionName string, columns []entity.Column) (bool, error) {

	if exist, err := c.HasCollection(ctx, collectionName); err != nil || !exist {
		return false, errors.New("No such collection : " + collectionName)
	}

	// insert into default partition
	_, err := c.Insert(ctx, collectionName, "", columns...)
	if err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}

	log.Println("insert completed")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	err = c.Flush(ctx, collectionName, false)
	if err != nil {
		log.Fatal("failed to flush collection:", err.Error())
	}

	log.Println("flush completed")
	return true, nil
}

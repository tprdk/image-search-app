package main

import (
	"context"
	"log"
	"time"

	"vectordbdemo/storage"

	"github.com/cilium/ebpf/features"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/spf13/viper"
)

const collectionName = `plip_768_collection`
const dbDefaultTimeout = time.Second * 5
const insertTimeout = time.Second * 60

const vectorDim = 768

func main() {
	viper.SetConfigFile("./common/envs/.env")
	viper.ReadInConfig()

	// add env variables as needed
	port := viper.Get("SERVER_PORT").(string)
	dbUrl := viper.Get("DB_URL").(string)

	log.Println("Milvus db url : ", dbUrl)
	log.Println("Server port : ", port)

	dbDefaultCtx := context.Background()
	dbDefaultCtx, cancelConnection := context.WithTimeout(dbDefaultCtx, dbDefaultTimeout)
	defer cancelConnection()

	log.Println("Connection milvus db.")
	mc, err := storage.InitMilvusClient(dbDefaultCtx, dbUrl)

	if err != nil {
		log.Fatal("Db connection could not initialized: ", err.Error())
	}
	defer mc.Close()

	log.Println("Creating schema.")
	schema := storage.GetDbSchema(collectionName, vectorDim)
	exist, err := storage.CreateCollectionIfNotExist(mc, dbDefaultCtx, collectionName, schema, false)

	if err != nil {
		log.Fatal("Failed to check collection exists: ", err.Error())
	}

	log.Println("Collection exist : ", exist)

	insertCtx := context.Background()
	insertCtx, cancelInsert := context.WithTimeout(insertCtx, insertTimeout)
	defer cancelInsert()
	

	featureVectors := LoadFeatureVectors()
	columns := []entity.Column{
		entity.NewColumnFloatVector("pkey", vectorDim, featureVectors)
	}

	storage.InsertIntoCollection()
}

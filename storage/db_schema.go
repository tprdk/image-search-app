package storage

import "github.com/milvus-io/milvus-sdk-go/v2/entity"

func GetDbSchema(collectionName string, vectorDim int) *entity.Schema {
	return &entity.Schema{
		CollectionName: collectionName,
		Description:    "this is the face images collection",
		AutoID:         true,
		Fields: []*entity.Field{
			// currently primary key field is compulsory, and only int64 is allowd
			{
				Name:       "pkey",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			// also the vector field is needed
			{
				Name:     "feature",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{ // the vector dim may changed def method in release
					entity.TypeParamDim: string(vectorDim),
				},
			},
		},
	}
}

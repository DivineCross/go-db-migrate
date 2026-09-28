package driver

import (
	"go.mongodb.org/mongo-driver/bson"
)

func unmarshalExtJson[T any](json string) (T, error) {
	var result T
	err := bson.UnmarshalExtJSON([]byte(json), true, &result)
	return result, err
}

func getField(document bson.D, key string) (any, bool) {
	for _, field := range document {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}

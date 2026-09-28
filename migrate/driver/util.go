package driver

import (
	"migrate/util"

	"go.mongodb.org/mongo-driver/bson"
)

func unmarshalExtJson[T any](json string) T {
	var result T
	util.Throw(bson.UnmarshalExtJSON([]byte(json), true, &result))
	return result
}

func getField(document bson.D, key string) (any, bool) {
	for _, field := range document {
		if field.Key == key {
			return field.Value, true
		}
	}
	return nil, false
}

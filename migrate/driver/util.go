package driver

import "go.mongodb.org/mongo-driver/bson"

func UnmarshalExtJson[T any](json string) T {
	var result T
	throw(bson.UnmarshalExtJSON([]byte(json), true, &result))
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

func throw(err error) {
	if err != nil {
		panic(err)
	}
}

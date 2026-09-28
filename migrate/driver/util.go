package driver

import "go.mongodb.org/mongo-driver/bson"

func UnmarshalExtJson[T any](json string) T {
	var result T
	throw(bson.UnmarshalExtJSON([]byte(json), true, &result))
	return result
}

func throw(err error) {
	if err != nil {
		panic(err)
	}
}

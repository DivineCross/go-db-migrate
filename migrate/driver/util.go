package driver

import "go.mongodb.org/mongo-driver/bson"

func unmarshalExtJson(json string) bson.D {
	var result bson.D
	throw(bson.UnmarshalExtJSON([]byte(json), true, &result))
	return result
}

func throw(err error) {
	if err != nil {
		panic(err)
	}
}

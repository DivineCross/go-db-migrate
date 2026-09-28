package script

import db "migrate/driver"

func init() {
	Register(s002Up, s002Down)
}

func s002Up() {
	db.EnsureSchema("Code", `{
		"bsonType": "object",
		"required": [
			"category",
			"code",
			"text",
			"sort",
			"disabled"
		],
		"properties": {
			"category": {"bsonType": "string"},
			"code": {"bsonType": "string"},
			"text": {"bsonType": "string"},
			"sort": {"bsonType": "int"},
			"disabled": {"bsonType": "bool"}
		}
	}`)
	db.EnsureSchema("EventDef", `{
		"bsonType": "object",
		"required": [
			"eventId",
			"name",
			"level",
			"shouldRecover"
		],
		"properties": {
			"eventId": {"bsonType": "string"},
			"name": {"bsonType": "string"},
			"level": {"bsonType": "int"},
			"shouldRecover": {"bsonType": "bool"},
			"message": {"bsonType": "string"},
			"recoverMessage": {"bsonType": "string"}
		}
	}`)
	db.EnsureSchema("Pricing", `{
		"bsonType": "object",
		"required": [
			"pricingId",
			"startAt",
			"productId",
			"litePrice",
			"plusPrice",
			"proPrice"
		],
		"properties": {
			"pricingId": {"bsonType": "string"},
			"startAt": {"bsonType": "date"},
			"productId": {"bsonType": "string"},
			"litePrice": {"bsonType": "double"},
			"plusPrice": {"bsonType": "double"},
			"proPrice": {"bsonType": "double"}
		}
	}`)
}

func s002Down() {
	db.DropSchema("Pricing")
	db.DropSchema("EventDef")
	db.DropSchema("Code")
}

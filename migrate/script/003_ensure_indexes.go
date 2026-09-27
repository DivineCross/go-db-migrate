package script

import db "migrate/driver"

func init() {
	Register(s003Up, s003Down)
}

func s003Up() {
	db.EnsureIndex("Code", `{
		"key": {"category": 1},
		"name": "category_1"
	}`)
	db.EnsureIndex("EventDef", `{
		"key": {"eventId": 1},
		"name": "eventId_1",
		"unique": true
	}`)
	db.EnsureIndex("Pricing", `{
		"key": {"pricingId": 1},
		"name": "pricingId_1"
	}`)
}

func s003Down() {
}

package script

import db "migrate/driver"

func init() {
	Register(s003Up, s003Down)
}

func s003Up() {
	db.EnsureIndex("Code", `{
		"key": {"category": {"$numberInt": "1"}}
	}`)
	db.EnsureIndex("EventDef", `{
		"key": {"eventId": {"$numberInt": "1"}},
		"unique": true
	}`)
	db.EnsureIndex("Pricing", `{
		"key": {"pricingId": {"$numberInt": "1"}}
	}`)
}

func s003Down() {
}

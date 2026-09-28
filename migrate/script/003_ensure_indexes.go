package script

func init() {
	Register(s003Up, s003Down)
}

func s003Up() {
	db.EnsureIndex("Code", `{
		"key": {"category": {"$numberInt": "1"}},
		"name": "category_1"
	}`)
	db.EnsureIndex("EventDef", `{
		"key": {"eventId": {"$numberInt": "1"}},
		"name": "eventId_1",
		"unique": true
	}`)
	db.EnsureIndex("Pricing", `{
		"key": {"pricingId": {"$numberInt": "1"}},
		"name": "pricingId_1"
	}`)
}

func s003Down() {
	db.DropIndex("Pricing", "pricingId_1")
	db.DropIndex("EventDef", "eventId_1")
	db.DropIndex("Code", "category_1")
}

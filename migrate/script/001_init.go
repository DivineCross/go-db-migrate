package script

import db "migrate/driver"

func init() {
	Register(s001Up, s001Down)
}

func s001Up() {
	db.EnsureColl("Code")
	db.EnsureColl("EventDef")
	db.EnsureColl("Pricing")
}

func s001Down() {
	db.DropColl("Pricing")
	db.DropColl("EventDef")
	db.DropColl("Code")
}

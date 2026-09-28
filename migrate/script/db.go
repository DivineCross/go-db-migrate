package script

import (
	"migrate/driver"
	"migrate/util"
)

var db scriptDb

type scriptDb struct{}

func (scriptDb) EnsureColl(name string) {
	util.Throw(driver.EnsureColl(name))
}

func (scriptDb) DropColl(name string) {
	util.Throw(driver.DropColl(name))
}

func (scriptDb) EnsureSchema(name, schemaJson string) {
	util.Throw(driver.EnsureSchema(name, schemaJson))
}

func (scriptDb) DropSchema(name string) {
	util.Throw(driver.DropSchema(name))
}

func (scriptDb) EnsureIndex(name, indexJson string) {
	util.Throw(driver.EnsureIndex(name, indexJson))
}

func (scriptDb) DropIndex(name, indexName string) {
	util.Throw(driver.DropIndex(name, indexName))
}

func (scriptDb) InsertWhenNotMatched(name, documentsJson string, fields ...string) {
	util.Throw(driver.InsertWhenNotMatched(name, documentsJson, fields...))
}

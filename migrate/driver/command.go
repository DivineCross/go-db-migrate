package driver

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"migrate/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	errCodeNamespaceExists = 48
)

func EnsureColl(name string) {
	err := runCommand(bson.D{{Key: "create", Value: name}})
	if err == nil {
		return
	}

	var commandErr mongo.CommandError
	if errors.As(err, &commandErr) && commandErr.Code == errCodeNamespaceExists {
		return
	}

	util.Throw(fmt.Errorf("create collection %s: %w", name, err))
}

func DropColl(name string) {
	if err := runCommand(bson.D{{Key: "drop", Value: name}}); err != nil {
		util.Throw(fmt.Errorf("drop collection %s: %w", name, err))
	}
}

func EnsureSchema(name, schemaJson string) {
	schema := unmarshalExtJson[bson.D](schemaJson)

	cmd := bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{{Key: "$jsonSchema", Value: schema}}},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}
	if err := runCommand(cmd); err != nil {
		util.Throw(fmt.Errorf("create validator %s: %w", name, err))
	}
}

func DropSchema(name string) {
	cmd := bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{}},
	}
	if err := runCommand(cmd); err != nil {
		util.Throw(fmt.Errorf("drop validator %s: %w", name, err))
	}
}

func EnsureIndex(name, indexJson string) {
	index := unmarshalExtJson[bson.D](indexJson)
	cmd := bson.D{
		{Key: "createIndexes", Value: name},
		{Key: "indexes", Value: bson.A{index}},
	}
	if err := runCommand(cmd); err != nil {
		util.Throw(fmt.Errorf("create index on %s: %w", name, err))
	}
}

func DropIndex(name, indexName string) {
	err := withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).Indexes().DropOne(cmdCtx, indexName)
		return err
	})
	if err != nil {
		util.Throw(fmt.Errorf("drop index %s on %s: %w", indexName, name, err))
	}
}

func InsertWhenNotMatched(name string, documentsJson string, fields ...string) {
	documents := unmarshalExtJson[[]bson.D](documentsJson)
	if len(fields) == 0 {
		util.Throw(fmt.Errorf("insert into %s requires match fields", name))
	}
	if len(documents) == 0 {
		return
	}

	oriDocs := FindAll[bson.D](name)
	var insDocs []any
	for _, doc := range documents {
		matched := false
		for _, oriDoc := range oriDocs {
			matched = true
			for _, field := range fields {
				value, exists := getField(doc, field)
				oriValue, oriExists := getField(oriDoc, field)
				if exists != oriExists || !reflect.DeepEqual(value, oriValue) {
					matched = false
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			insDocs = append(insDocs, doc)
		}
	}

	Inserts(name, insDocs)
}

func Inserts(name string, documents []any) {
	if len(documents) == 0 {
		return
	}
	err := withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).InsertMany(cmdCtx, documents)
		return err
	})
	if err != nil {
		util.Throw(fmt.Errorf("insert documents into %s: %w", name, err))
	}
}

func InsertOne(name string, document any) error {
	return withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).InsertOne(cmdCtx, document)
		return err
	})
}

func UpdateOne(name string, filter, update bson.D) (bool, error) {
	var count int64
	err := withCommandContext(func(cmdCtx context.Context) error {
		result, err := ctx.db.Collection(name).UpdateOne(cmdCtx, filter, update)
		if err == nil {
			count = result.MatchedCount
		}
		return err
	})
	return count == 1, err
}

func FindAll[T any](name string) []T {
	var documents []T
	err := withCommandContext(func(cmdCtx context.Context) error {
		cursor, err := ctx.db.Collection(name).Find(cmdCtx, bson.D{})
		if err != nil {
			return err
		}
		return cursor.All(cmdCtx, &documents)
	})
	if err != nil {
		util.Throw(fmt.Errorf("find all documents in %s: %w", name, err))
	}
	return documents
}

func RunCommand(commandJson string) {
	command := unmarshalExtJson[bson.D](commandJson)

	if err := runCommand(command); err != nil {
		util.Throw(fmt.Errorf("run command %s: %w", commandJson, err))
	}
}

func hasIndex(collName, indexName string) (bool, error) {
	var indexes []*mongo.IndexSpecification
	err := withCommandContext(func(cmdCtx context.Context) (err error) {
		indexes, err = ctx.db.Collection(collName).Indexes().ListSpecifications(cmdCtx)
		return err
	})
	if err != nil {
		return false, err
	}

	for _, index := range indexes {
		if index.Name == indexName {
			return true, nil
		}
	}
	return false, nil
}

func runCommand(cmd bson.D) error {
	return runCommandResult(cmd, nil)
}

func runCommandResult(cmd bson.D, result any) error {
	return withCommandContext(func(cmdCtx context.Context) error {
		response := ctx.db.RunCommand(cmdCtx, cmd)
		if result != nil {
			return response.Decode(result)
		}
		return response.Err()
	})
}

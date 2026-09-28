package driver

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

	throw(fmt.Errorf("create collection %s: %w", name, err))
}

func DropColl(name string) {
	if err := runCommand(bson.D{{Key: "drop", Value: name}}); err != nil {
		throw(fmt.Errorf("drop collection %s: %w", name, err))
	}
}

func EnsureSchema(name, schemaJson string) {
	schema := UnmarshalExtJson[bson.D](schemaJson)

	cmd := bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{{Key: "$jsonSchema", Value: schema}}},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}
	if err := runCommand(cmd); err != nil {
		throw(fmt.Errorf("create validator %s: %w", name, err))
	}
}

func DropSchema(name string) {
	cmd := bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{}},
	}
	if err := runCommand(cmd); err != nil {
		throw(fmt.Errorf("drop validator %s: %w", name, err))
	}
}

func EnsureIndex(name, indexJson string) {
	index := UnmarshalExtJson[bson.D](indexJson)
	cmd := bson.D{
		{Key: "createIndexes", Value: name},
		{Key: "indexes", Value: bson.A{index}},
	}
	if err := runCommand(cmd); err != nil {
		throw(fmt.Errorf("create index on %s: %w", name, err))
	}
}

func DropIndex(name, indexName string) {
	err := withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).Indexes().DropOne(cmdCtx, indexName)
		return err
	})
	if err != nil {
		throw(fmt.Errorf("drop index %s on %s: %w", indexName, name, err))
	}
}

func InsertWhenNotMatched(name string, documents []bson.M, fields ...string) {
	if len(fields) == 0 {
		throw(fmt.Errorf("insert into %s requires match fields", name))
	}
	if len(documents) == 0 {
		return
	}

	oriDocs := FindAll(name)
	var insDocs []any
	for _, doc := range documents {
		matched := false
		for _, oriDoc := range oriDocs {
			matched = true
			for _, field := range fields {
				value, exists := doc[field]
				oriValue, oriExists := oriDoc[field]
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
		throw(fmt.Errorf("insert documents into %s: %w", name, err))
	}
}

func FindAll(name string) []bson.M {
	var documents []bson.M
	err := withCommandContext(func(cmdCtx context.Context) error {
		cursor, err := ctx.db.Collection(name).Find(cmdCtx, bson.D{})
		if err != nil {
			return err
		}
		return cursor.All(cmdCtx, &documents)
	})
	if err != nil {
		throw(fmt.Errorf("find all documents in %s: %w", name, err))
	}
	return documents
}

func RunCommand(commandJson string) {
	command := UnmarshalExtJson[bson.D](commandJson)

	if err := runCommand(command); err != nil {
		throw(fmt.Errorf("run command %s: %w", commandJson, err))
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

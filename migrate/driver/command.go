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

func EnsureColl(name string) error {
	err := runCommand(bson.D{{Key: "create", Value: name}})
	if err == nil {
		return nil
	}

	var commandErr mongo.CommandError
	if errors.As(err, &commandErr) && commandErr.Code == errCodeNamespaceExists {
		return nil
	}
	return err
}

func DropColl(name string) error {
	return runCommand(bson.D{{Key: "drop", Value: name}})
}

func EnsureSchema(name, schemaJson string) error {
	schema, err := unmarshalExtJson[bson.D](schemaJson)
	if err != nil {
		return err
	}

	return runCommand(bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{{Key: "$jsonSchema", Value: schema}}},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	})
}

func DropSchema(name string) error {
	return runCommand(bson.D{
		{Key: "collMod", Value: name},
		{Key: "validator", Value: bson.D{}},
	})
}

func EnsureIndex(name, indexJson string) error {
	index, err := unmarshalExtJson[bson.D](indexJson)
	if err != nil {
		return err
	}

	return runCommand(bson.D{
		{Key: "createIndexes", Value: name},
		{Key: "indexes", Value: bson.A{index}},
	})
}

func DropIndex(name, indexName string) error {
	return withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).Indexes().DropOne(cmdCtx, indexName)
		return err
	})
}

func InsertWhenNotMatched(name string, documentsJson string, fields ...string) error {
	documents, err := unmarshalExtJson[[]bson.D](documentsJson)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return fmt.Errorf("insert into %s requires match fields", name)
	}
	if len(documents) == 0 {
		return nil
	}

	oriDocs, err := FindAll[bson.D](name)
	if err != nil {
		return err
	}

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

	return InsertMany(name, insDocs)
}

func InsertMany(name string, documents []any) error {
	if len(documents) == 0 {
		return nil
	}

	return withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).InsertMany(cmdCtx, documents)
		return err
	})
}

func InsertOne(name string, document any) error {
	return withCommandContext(func(cmdCtx context.Context) error {
		_, err := ctx.db.Collection(name).InsertOne(cmdCtx, document)
		return err
	})
}

func UpdateOne(name string, filter, update bson.D) (bool, error) {
	var matched bool
	err := withCommandContext(func(cmdCtx context.Context) error {
		result, updateErr := ctx.db.Collection(name).UpdateOne(cmdCtx, filter, update)
		if updateErr == nil {
			matched = result.MatchedCount == 1
		}
		return updateErr
	})
	return matched, err
}

func FindAll[T any](name string) ([]T, error) {
	var documents []T
	err := withCommandContext(func(cmdCtx context.Context) error {
		cursor, findErr := ctx.db.Collection(name).Find(cmdCtx, bson.D{})
		if findErr != nil {
			return findErr
		}
		return cursor.All(cmdCtx, &documents)
	})

	if err != nil {
		return nil, err
	}
	return documents, nil
}

func runCommand(cmd bson.D) error {
	return withCommandContext(func(cmdCtx context.Context) error {
		return ctx.db.RunCommand(cmdCtx, cmd).Err()
	})
}

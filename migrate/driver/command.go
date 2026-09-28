package driver

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func CreateSchema(name, schemaJson string) {
	schema := unmarshalExtJson(schemaJson)

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
	index := unmarshalExtJson(indexJson)
	var indexName string
	for _, field := range index {
		if field.Key == "name" {
			indexName, _ = field.Value.(string)
		}
	}
	if indexName == "" {
		throw(fmt.Errorf("index name must be a non-empty string"))
	}

	exists, err := hasIndex(name, indexName)
	if err != nil {
		throw(fmt.Errorf("check index %s on %s: %w", indexName, name, err))
	}
	if exists {
		return
	}

	cmd := bson.D{
		{Key: "createIndexes", Value: name},
		{Key: "indexes", Value: bson.A{index}},
	}
	if err := runCommand(cmd); err != nil {
		throw(fmt.Errorf("create index %s on %s: %w", indexName, name, err))
	}
}

func RunCommand(commandJson string) {
	command := unmarshalExtJson(commandJson)

	if err := runCommand(command); err != nil {
		throw(fmt.Errorf("run command %s: %w", commandJson, err))
	}
}

func hasIndex(collName, indexName string) (bool, error) {
	var indexes []struct {
		Name string `bson:"name"`
	}
	cmd := bson.D{{Key: "listIndexes", Value: collName}}
	if err := runCommandCursor(cmd, &indexes); err != nil {
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

func runCommandCursor(cmd bson.D, result any) error {
	return withCommandContext(func(cmdCtx context.Context) (err error) {
		cursor, err := ctx.db.RunCommandCursor(cmdCtx, cmd)
		if err != nil {
			return err
		}

		timeout := 10 * time.Second
		defer func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			err = errors.Join(err, cursor.Close(closeCtx))
		}()
		return cursor.All(cmdCtx, result)
	})
}

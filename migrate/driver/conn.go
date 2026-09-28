package driver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Close() error {
	timeout := 10 * time.Second

	mu.Lock()
	defer mu.Unlock()
	if ctx == nil {
		return nil
	}
	commandCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := ctx.client.Disconnect(commandCtx)
	ctx = nil
	return err
}

func withCommandContext(action func(context.Context) error) error {
	timeout := 1 * time.Minute

	mu.Lock()
	defer mu.Unlock()
	if err := connect(); err != nil {
		return err
	}
	commandCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return action(commandCtx)
}

func connect() error {
	if ctx != nil {
		return nil
	}

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	uri := os.Getenv("MONGODB_URI")
	database := os.Getenv("MONGODB_DATABASE")
	if uri == "" || database == "" {
		return fmt.Errorf("MONGODB_URI and MONGODB_DATABASE are required")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx = &mongoCtx{client: client, db: client.Database(database)}
	return nil
}

var (
	mu  sync.Mutex
	ctx *mongoCtx
)

type mongoCtx struct {
	client *mongo.Client
	db     *mongo.Database
}

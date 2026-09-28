package runner

import (
	"fmt"
	"time"

	db "migrate/driver"
	"migrate/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const migrationCollName = "__migration__"

type migrationRecord struct {
	Id          primitive.ObjectID `bson:"_id"`
	StartAt     time.Time          `bson:"startAt"`
	EndAt       *time.Time         `bson:"endAt"`
	FromVersion string             `bson:"fromVersion"`
	ToVersion   string             `bson:"toVersion"`
	Error       *string            `bson:"error"`
}

func getCurrentVersion() (*string, error) {
	documents, err := db.FindAll[migrationRecord](migrationCollName)
	if err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		version := "0.0.0"
		return &version, nil
	}

	for _, current := range documents {
		if current.StartAt.IsZero() {
			return nil, fmt.Errorf("migration record has invalid startAt")
		}
	}
	sorted := util.ToSorted(documents, func(a, b migrationRecord) int {
		return b.StartAt.Compare(a.StartAt)
	})

	doc := &sorted[0]
	if doc.EndAt == nil || doc.Error != nil {
		return nil, nil
	}

	if doc.ToVersion == "" {
		return nil, fmt.Errorf("latest migration record has no toVersion")
	}
	return &doc.ToVersion, nil
}

func startMigration(fromVersion, toVersion string) (primitive.ObjectID, error) {
	id := primitive.NewObjectID()

	doc := migrationRecord{
		Id:          id,
		StartAt:     time.Now().UTC(),
		FromVersion: fromVersion,
		ToVersion:   toVersion,
	}
	return id, db.InsertOne(migrationCollName, doc)
}

func endMigration(id primitive.ObjectID, stepErr error) error {
	var errorMessage any
	if stepErr != nil {
		errorMessage = stepErr.Error()
	}
	matched, err := db.UpdateOne(migrationCollName,
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "endAt", Value: time.Now().UTC()},
			{Key: "error", Value: errorMessage},
		}}},
	)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("migration record %s was not found", id.Hex())
	}
	return nil
}

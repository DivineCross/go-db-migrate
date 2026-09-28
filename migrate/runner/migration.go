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

func getCurrentVersion() (version *string, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			version = nil
			if cause, ok := recovered.(error); ok {
				err = fmt.Errorf("get current migration version: %w", cause)
			} else {
				err = fmt.Errorf("get current migration version: %v", recovered)
			}
		}
	}()
	documents := db.FindAll[migrationRecord](migrationCollName)

	if len(documents) == 0 {
		version := "0.0.0"
		return &version, nil
	}

	for _, current := range documents {
		if current.StartAt.IsZero() || current.Id.IsZero() {
			return nil, fmt.Errorf("migration record has invalid startAt or _id")
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
	if err := db.InsertOne(migrationCollName, doc); err != nil {
		return id, fmt.Errorf("record migration start %s -> %s: %w", fromVersion, toVersion, err)
	}
	return id, nil
}

func endMigration(id primitive.ObjectID, stepErr error) error {
	var errorMessage any
	if stepErr != nil {
		errorMessage = stepErr.Error()
	}
	updated, err := db.UpdateOne(migrationCollName,
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "endAt", Value: time.Now().UTC()},
			{Key: "error", Value: errorMessage},
		}}},
	)
	if err != nil {
		return fmt.Errorf("record migration end %s: %w", id.Hex(), err)
	}
	if !updated {
		return fmt.Errorf("migration record %s was not found", id.Hex())
	}
	return nil
}

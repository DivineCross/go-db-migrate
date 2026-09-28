package runner

import (
	"errors"
	"fmt"

	db "migrate/driver"
	"migrate/script"
	"migrate/util"
)

func Run(fromVersion, toVersion string) (err error) {
	fromSeq, exists := script.GetSeq(fromVersion)
	if !exists {
		return fmt.Errorf("migration version %s is not registered", fromVersion)
	}
	toSeq, exists := script.GetSeq(toVersion)
	if !exists {
		return fmt.Errorf("migration version %s is not registered", toVersion)
	}

	isUp := fromSeq <= toSeq
	dir := util.Ternary(isUp, 1, -1)

	var steps []script.Step
	for seq := fromSeq; seq != toSeq; seq += dir {
		stepSeq := util.Ternary(isUp, seq+dir, seq)
		step, exists := script.GetStep(stepSeq)
		if !exists {
			return fmt.Errorf("migration step %d is not registered", stepSeq)
		}
		steps = append(steps, step)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	curVersion, err := getCurrentVersion()
	if err != nil {
		return err
	}
	if curVersion == nil {
		return fmt.Errorf("database version is unknown")
	}
	if *curVersion != fromVersion {
		return fmt.Errorf("%s does not match database version %s", fromVersion, *curVersion)
	}
	if len(steps) == 0 {
		return nil
	}

	recordId, startErr := startMigration(fromVersion, toVersion)
	if startErr != nil {
		return startErr
	}
	for _, step := range steps {
		dirText := util.Ternary(isUp, "up", "down")
		fn := util.Ternary(isUp, step.Up, step.Down)

		fmt.Printf("Running migration %03d %s (%s)\n", step.Seq, dirText, step.Name)
		stepErr := fn()
		if stepErr != nil {
			endErr := endMigration(recordId, stepErr)
			return errors.Join(
				fmt.Errorf("migration %d %s (%s): %w", step.Seq, dirText, step.Name, stepErr),
				endErr,
			)
		}
	}
	return endMigration(recordId, nil)
}

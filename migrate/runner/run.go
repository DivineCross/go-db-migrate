package runner

import (
	"errors"
	"fmt"

	db "migrate/driver"
	"migrate/script"
)

func Run(from, to int) (err error) {
	if from < 0 || to < 0 {
		return fmt.Errorf("step sequences must be non-negative")
	}

	isUp := from <= to
	dir := ternary(isUp, 1, -1)

	var steps []script.Step
	for seq := from; seq != to; seq += dir {
		stepSeq := ternary(isUp, seq+dir, seq)
		step, exists := script.Get(stepSeq)
		if !exists {
			return fmt.Errorf("migration step %d is not registered", stepSeq)
		}
		steps = append(steps, step)
	}
	if len(steps) == 0 {
		return nil
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()
	for _, step := range steps {
		dirText := ternary(isUp, "up", "down")
		fn := ternary(isUp, step.Up, step.Down)
		fmt.Printf("Running migration %03d %s (%s)\n", step.Seq, dirText, step.Name)
		if err := fn(); err != nil {
			return fmt.Errorf("migration %d %s (%s): %w", step.Seq, dirText, step.Name, err)
		}
	}
	return nil
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

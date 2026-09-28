package script

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

type Step struct {
	Seq  int
	Name string
	Up   func() error
	Down func() error
}

func Register(up, down func()) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("cannot determine migration filename")
	}
	if err := register(file, up, down); err != nil {
		panic(err)
	}
}

func Get(seq int) (Step, bool) {
	step, exists := registry[seq]
	return step, exists
}

func register(file string, up, down func()) error {
	name := filepath.Base(file)
	parts := filenameRegexp.FindStringSubmatch(name)
	if parts == nil {
		return fmt.Errorf("invalid migration filename: %s", name)
	}
	seq, err := strconv.Atoi(parts[1])
	if err != nil || seq <= 0 {
		return fmt.Errorf("invalid migration sequence: %s", parts[1])
	}
	if up == nil || down == nil {
		return fmt.Errorf("migration step %d requires up and down functions", seq)
	}
	if _, exists := registry[seq]; exists {
		return fmt.Errorf("migration step %d is already registered", seq)
	}

	registry[seq] = Step{
		Seq:  seq,
		Name: parts[2],
		Up:   catch(up),
		Down: catch(down),
	}
	return nil
}

func catch(action func()) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				if cause, ok := r.(error); ok {
					err = cause
				} else {
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		action()
		return nil
	}
}

var (
	registry       = map[int]Step{}
	filenameRegexp = regexp.MustCompile(`^([0-9]{1,6})_(.+)\.go$`)
)

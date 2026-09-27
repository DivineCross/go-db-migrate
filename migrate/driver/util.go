package driver

type Failure struct {
	Cause error
}

func (f *Failure) Error() string {
	return f.Cause.Error()
}

func (f *Failure) Unwrap() error {
	return f.Cause
}

func throw(err error) {
	if err != nil {
		panic(&Failure{Cause: err})
	}
}

package tapedeck

type TagError struct {
	Path string
	Err  error
}

func (e *TagError) Error() string {
	return e.Path + ": " + e.Err.Error()
}

func newTagError(path string, err error) error {
	if err == nil {
		return nil
	}

	return &TagError{path, err}
}

package apperror

type NotFoundErr struct {
	Resource string
	Detail   string
}

func (e *NotFoundErr) Error() string {
	return e.Resource + " " + e.Detail + ": not found"
}

func NewNotFoundErr(resource, detail string) *NotFoundErr {
	return &NotFoundErr{
		Resource: resource,
		Detail:   detail,
	}
}

type AlreadyExistsErr struct {
	Resource string
	Detail   string
}

func (e *AlreadyExistsErr) Error() string {
	return e.Resource + " " + e.Detail + ": already exists"
}

func NewAlreadyExistsErr(resource, detail string) *AlreadyExistsErr {
	return &AlreadyExistsErr{
		Resource: resource,
		Detail:   detail,
	}
}

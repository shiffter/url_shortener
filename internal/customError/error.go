package customError

import "errors"

var (
	ErrBadRequest   = errors.New("Bad request params")
	UrlAlreadyExist = errors.New("url already exist")
	InternalErr     = errors.New("internal error")
	NotFound        = errors.New("no existing short url")
)

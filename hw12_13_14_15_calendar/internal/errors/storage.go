package internalerrors

import (
	"errors"
)

var (
	ErrDocumentNotFound           = errors.New("document not found")
	ErrDocumentOperationForbidden = errors.New("document operation forbidden")
)

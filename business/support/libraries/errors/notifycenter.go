package errors

import (
	"net/http"
)

var (
	SQSEmpty = NewError(http.StatusOK, "SQS_EMPTY", "The SQS query is empty")
)

package handlers

import (
	"time"

	"github.com/vo1dFl0w/qa-service/internal/transport/utils"
)

var (
	ctxDelay = time.Second * 5
)

func validateMethod(expMethod string, gotMethod string) error {
	if gotMethod != expMethod {
		return utils.ErrMethodNotAllowed
	}

	return nil
}
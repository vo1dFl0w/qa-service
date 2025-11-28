package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/utils"
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

func getIdFromURL(r *http.Request) (int, error) {
	idStr := r.PathValue("id")

	if idStr == "" {
		return 0, fmt.Errorf("empty id")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("cannot convert id from url")
	}

	return id, nil
}
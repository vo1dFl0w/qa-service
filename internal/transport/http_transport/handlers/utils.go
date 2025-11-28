package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ctxDelay = time.Second * 5
)

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
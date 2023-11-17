// cmd/api/helpers.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) storeInRedis(prefix string, hash string, userID uuid.UUID, expiration time.Duration) error {
	ctx := context.Background()
	err := app.redisClient.Set(
		ctx,
		fmt.Sprintf("%s%s", prefix, userID),
		hash,
		expiration,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil, app.config.debug)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

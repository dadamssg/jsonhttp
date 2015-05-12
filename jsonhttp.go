package jsonhttp

import (
	"net/http"

	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type ResponseError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Errors []ResponseError `json:"errors"`
}

func MapOrSendError(w http.ResponseWriter, r *http.Request, data interface{}) error {
	if err := MapJSON(r, data); err != nil {
		SendError(w, 400, err)
		return err
	}

	return nil
}

func SendError(w http.ResponseWriter, statusCode int, err error) {
	errors := []ResponseError{ResponseError{Code: statusCode, Error: err.Error()}}
	SendErrors(w, errors)
}

func SendErrors(w http.ResponseWriter, errors []ResponseError) {
	code := 500
	if len(errors) > 0 {
		code = errors[0].Code
	}
	SendJSON(w, code, ErrorResponse{Errors: errors})
}

func MapJSON(r *http.Request, subject interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return errors.New("Could not read data.")
	}
	if err := r.Body.Close(); err != nil {
		return errors.New("Could not finish reading data.")
	}

	if err := json.Unmarshal(body, &subject); err != nil {
		return errors.New("Invalid json.")
	}

	return nil
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

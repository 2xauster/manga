package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type httpHandler func(w http.ResponseWriter, r *http.Request) Response

type Response struct {
	Status int `json:"status"`
	StatusText string `json:"status_text"`
	Error error `json:"error,omitempty"`
	D interface{} `json:"d"`
	Time time.Time `json:"time"`
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/health", handleRequest(HandleGetHealth))
}

func handleRequest(f httpHandler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		prefix := "[" + path + "] " + "[" + method + "] "

		res := f(w, r)

		if res.Error != nil {
			log.Println(prefix, "[ERROR]", res.Error.Error())
		}
		if err := writeJSON(res, w); err != nil {
			log.Println(prefix, "[ERROR]", "Could not write json", err)
			return
		}
		log.Println(prefix, "[INFO]", "Accessed")
	}
}

func newResponse(status int, d interface{}, err error) Response {
	return Response{
		Status: status,
		StatusText: http.StatusText(status),
		D: d,
		Error: err,
		Time: time.Now(),
	}
}

func writeJSON(r Response, w http.ResponseWriter) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return json.NewEncoder(w).Encode(r)
}

func readJSON(r *http.Request, target interface{}) (err error) {
	if err = json.NewDecoder(r.Body).Decode(target); err != nil {
		return err 
	}

	if err = validate.Struct(target); err != nil {
		return err
	}
	return nil
}
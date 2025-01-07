package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)


type httpHandler func(w http.ResponseWriter, r *http.Request) Response

type Response struct {
	Status int `json:"status"`
	StatusText string `json:"status_text"`
	Error string `json:"error,omitempty"`
	D interface{} `json:"d"`
	Time time.Time `json:"time"`
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/health", handleRequest(HandleGetHealth))
	router.Handle("/manga/latest", handleRequest(HandleGetLatestManga))
	router.Handle("/manga/get/{id}", handleRequest(HandleGetManga))
	router.Handle("/manga/search", handleRequest(HandlePostSearch))
	
	router.Handle("/chapter/panels/{mangaID}/{chapterID}", handleRequest(HandleGetChapterPanels))
}

func handleRequest(f httpHandler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method
		prefix := "[" + path + "] " + "[" + method + "] "
		vars := mux.Vars(r)

		res := f(w, r)

		if res.Error != "" {
			log.Println("[ERROR]", prefix, "(", vars, ") ", res.Error)
		}
		if err := writeJSON(res, w); err != nil {
			log.Println(prefix, "[ERROR]", "Could not write json", err)
			return
		}
		log.Println("[INFO]", prefix, "Accessed")
	}
}

func newResponse(status int, d interface{}, err error) Response {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Status: status,
		StatusText: http.StatusText(status),
		D: d,
		Error: errMsg,
		Time: time.Now(),
	}
}

func writeJSON(r Response, w http.ResponseWriter) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	return json.NewEncoder(w).Encode(r)
}

func readJSON(r *http.Request, target interface{}) (err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) 
	if err = json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return 	
}
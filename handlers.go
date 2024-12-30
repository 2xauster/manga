package main

import "net/http"

func HandleGetHealth(w http.ResponseWriter, r *http.Request) Response {
	return newResponse(http.StatusOK, "+", nil)
}
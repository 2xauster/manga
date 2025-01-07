package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/2xauster/manga/scraper"
	"github.com/gorilla/mux"
)

type PostSearchMangaRequest struct {
	Query string `json:"query"`
}

func HandleGetHealth(w http.ResponseWriter, r *http.Request) Response {
	return newResponse(http.StatusOK, "+", nil)
}


func HandleGetLatestManga(w http.ResponseWriter, r *http.Request) Response {
	data, err := scraper.ScrapeLatest()
	if err != nil {
		return newResponse(http.StatusInternalServerError, nil, err)
	}
	return newResponse(http.StatusOK, data, nil)
}

func HandleGetManga(w http.ResponseWriter, r *http.Request) Response {
	id := mux.Vars(r)["id"]
	if id == "" {
		return newResponse(http.StatusBadRequest, nil, errors.New("missing required parameter: id"))
	}

	data, err := scraper.ScrapeManga(id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, scraper.ErrNotFound) {
			log.Println("nope")
			return newResponse(http.StatusBadRequest, nil, err)
		}
		return newResponse(http.StatusInternalServerError, nil, err)
	}
	return newResponse(http.StatusOK, data, nil)
}

func HandlePostSearch(w http.ResponseWriter, r *http.Request) Response {
	if r.Method != "POST" {
		return newResponse(http.StatusBadRequest, nil, errors.New("invalid method"))
	}
	var body PostSearchMangaRequest
	if err := readJSON(r, &body); err != nil {
		return newResponse(http.StatusBadRequest, nil, err)
	} 

	res, err := scraper.SearchManga(body.Query)
	if err != nil {
		return newResponse(http.StatusBadRequest, nil, err)
	}
	return newResponse(http.StatusOK, res, nil)
}

func HandleGetChapterPanels(w http.ResponseWriter, r *http.Request) Response {
	mangaID, chapterID := mux.Vars(r)["mangaID"], mux.Vars(r)["chapterID"]

	if mangaID == "" || chapterID == "" {
		return newResponse(http.StatusBadRequest, nil, errors.New("missing params"))
	}
	data, err := scraper.ScrapeChapterPanels(mangaID, chapterID)

	if err != nil {
		if errors.Is(err, scraper.ErrNotFound) {
			return newResponse(http.StatusBadRequest, nil, err)
		}
		return newResponse(http.StatusInternalServerError, nil, err)
	}
	return newResponse(http.StatusOK, data,nil)
}
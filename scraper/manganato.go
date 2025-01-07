package scraper

import (
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func createCollector() (*colly.Collector, chan error) {
	errChan := make(chan error, 1)
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		log.Println(r.StatusCode)
		if r.StatusCode == 404 {
			errChan <- ErrNotFound
			return
		}
		errChan <- err
	})
	return c, errChan
}

// Scraping Functions
func ScrapeLatest() (mangas []LatestManga, err error) {
	c, errCh := createCollector()

	c.OnHTML(LatestMangaSelector, func(e *colly.HTMLElement) { 
		link := e.Attr("href")
		title := e.Attr("title")
		img := e.ChildAttr("img.img-loading", "src")

		if link != "" && title != "" && img != "" {
			id := extractIDFromLink(link)
			mangas = append(mangas, LatestManga{
				ID:        id,
				Title:     title,
				Link:      link,
				Thumbnail: img,
			})
		}
	})

	if err = c.Visit(endpoint("/")); err != nil {
		return nil, err
	}

	select {
	case scrapeErr := <-errCh:
		return nil, scrapeErr
	default:
	}
	return mangas, nil
}

func ScrapeManga(id string) (manga Manga, err error) {
	c, errCh := createCollector()

	c.OnHTML(MangaInfoSelector, func(e *colly.HTMLElement) {
		imageSelector := "div.story-info-left span.info-image img.img-loading"

		manga.Thumbnail = e.ChildAttr(imageSelector, "src")
		manga.Title = e.ChildAttr(imageSelector, "title")
		manga.Meta = extractMangaInfoFromTable(e)
	})

	c.OnHTML(MangaDescSelector, func (e *colly.HTMLElement) {
		desc := strings.TrimPrefix(e.Text, "\n        Description :\n        ")
		desc = strings.ReplaceAll(desc, "  ", "")
		if desc != "" {
			manga.Description = desc
		}
	})

	c.Visit(CHAPMANGANATO_URL + "/" + id)

	select {
	case scrapeErr := <- errCh:
		return manga, scrapeErr
	default:
	}
	close(errCh)
	return manga, err
}

func SearchManga(query string) (results []SearchResult, err error) {
	c, errCh := createCollector()

	c.OnHTML(SearchResultSelector, func(e *colly.HTMLElement) {
		e.ForEach("div.search-story-item", func(_ int, h *colly.HTMLElement) {
			link, _ := h.DOM.Find("a.item-img").Attr("href")
			title, _ := h.DOM.Find("a.item-img").Attr("title")
			imgSrc, _ := h.DOM.Find("a.item-img img.img-loading").Attr("src")

			if link != "" {
				results = append(results, SearchResult{
					ID:        extractIDFromLink(link),
					Title:     title,
					Thumbnail: imgSrc,
				})
			}
		})
	})

	
	c.Visit(endpoint("/search/story/" + query))

	select {
	case scrapeErr := <- errCh:
		return nil, scrapeErr
	default:
	}
	close(errCh)

	return results, nil
}

func ScrapeChapterPanels(mangaID, chapterID string) (chapter Chapter, err error) {
	var panels []string
	var title string
	c, errCh := createCollector()

	c.OnHTML(ChapterPanelSelector, func(h *colly.HTMLElement) {
		if src := h.Attr("src"); src != "" {
			panels = append(panels, src)
		}
	})

	c.OnHTML(ChapterTitleSelector, func(h *colly.HTMLElement) {
		title = h.Text
	})

	c.Visit(CHAPMANGANATO_URL + "/" + mangaID + "/" + chapterID)

	chNumber, err := extractNumberFromID(chapterID)
	if err != nil {
		return chapter, err
	}

	select {
	case scrapeErr := <- errCh:
		return chapter, scrapeErr
	default:
	}

	chapter = Chapter{
		ID:        chapterID,
		Number:    chNumber,
		Title:     title,
		PanelURLs: panels,
	}
	return chapter, nil
	}

func extractMangaInfoFromTable(e *colly.HTMLElement) []Metadata {
	var meta []Metadata
	e.ForEach(MangaTableSelector, func(_ int, h *colly.HTMLElement) {
		replacer := strings.NewReplacer(" ", "", ":","")
		label := replacer.Replace(h.ChildText("td.table-label"))
		value := replacer.Replace(h.ChildText("td.table-value"))
		meta = append(meta, Metadata{Label: label, Value: value})
	})
	return meta
}
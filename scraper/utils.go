package scraper

import (
	"fmt"
	"strconv"
	"strings"
)

// Utility Functions
func endpoint(path string) string {
	return MANGANATO_URL + path
}

func extractIDFromLink(link string) string {
	parts := strings.Split(link, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

func extractNumberFromID(chapterID string) (int, error) {
	parts := strings.Split(chapterID, "-")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid chapter ID: %s", chapterID)
	}
	return strconv.Atoi(parts[1])
}
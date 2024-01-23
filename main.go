package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Pilfer/ultimate-guitar-scraper/pkg/ultimateguitar"
)

func tabId(url string) (int64, error) {
	invalidErr := errors.New("invalid UG URL")

	re := regexp.MustCompile("\\d+$")
	match := re.FindString(url)
	if match == "" {
		return 0, invalidErr
	}

	id, err := strconv.Atoi(match)
	if err != nil {
		return 0, invalidErr
	}

	return int64(id), nil
}

type TabOutput struct {
	SongName   string
	ArtistName string
	Url        string
	TabOut     string
}

func fetchTab(url string) (TabOutput, error) {
	tabOutput := TabOutput{}
	id, err := tabId(url)

	if err != nil {
		return tabOutput, err
	}

	s := ultimateguitar.New()
	tab, err := s.GetTabByID(id)
	if err != nil {
		return tabOutput, err
	}

	// Remove the syntax delimiters.
	tabOut := strings.ReplaceAll(tab.Content, "[tab]", "")
	tabOut = strings.ReplaceAll(tabOut, "[/tab]", "")
	tabOut = strings.ReplaceAll(tabOut, "[ch]", "")
	tabOut = strings.ReplaceAll(tabOut, "[/ch]", "")

	tabOutput.SongName = tab.SongName
	tabOutput.ArtistName = tab.ArtistName
	tabOutput.TabOut = tabOut
	tabOutput.Url = tab.URLWeb

	return tabOutput, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "https://tabs.ultimate-guitar.com/tab/misc-traditional/the-parting-glass-chords-1147884"
	tab, err := fetchTab(url)
	if err != nil {
		http.Error(w, "Error fetching tab", 500)
	}

	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, tab)
}

func main() {
	http.HandleFunc("/", handler)
	port := ":8080"
	fmt.Printf("Running on http://localhost%v\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

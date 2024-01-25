package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Pilfer/ultimate-guitar-scraper/pkg/ultimateguitar"
)

var tmpl = template.Must(template.ParseFiles("template.html"))

func tabId(url string) (int64, error) {
	invalidErr := errors.New("Invalid UG URL")

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
	URL        string
	TabOut     string
	Error      string
}

func fetchTab(id int64) (TabOutput, error) {
	tabOutput := TabOutput{}

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
	tabOutput.URL = tab.URLWeb

	return tabOutput, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := TabOutput{}

	path := r.URL.Path[len("/"):]
	if path == "" {
		tmpl.Execute(w, data)
		return
	}

	id, err := tabId(path)
	if err != nil {
		data.Error = err.Error()
		tmpl.Execute(w, data)
		return
	}

	data, err = fetchTab(id)
	if err != nil || data.TabOut == "" {
		data.Error = "Error fetching tab"
	}

	tmpl.Execute(w, data)
}

func main() {
	portInt := flag.Int("port", 3000, "port to run on")
	flag.Parse()
	port := fmt.Sprintf(":%v", *portInt)

	http.Handle("/ugs/static/", http.StripPrefix("/ugs/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/ugs", handler)

	fmt.Printf("Running on http://localhost%v/ugs\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

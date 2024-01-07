package main

import (
	"errors"
	"fmt"
	"log"
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

	// Remove the syntax delimiters as a proof of concept
	tabOut := strings.ReplaceAll(tab.Content, "[tab]", "")
	tabOut = strings.ReplaceAll(tabOut, "[/tab]", "")
	tabOut = strings.ReplaceAll(tabOut, "[ch]", "")
	tabOut = strings.ReplaceAll(tabOut, "[/ch]", "")

	tabOutput.SongName = tab.SongName
	tabOutput.ArtistName = tab.ArtistName
	tabOutput.TabOut = tabOut

	return tabOutput, nil
}

func main() {
	url := "https://tabs.ultimate-guitar.com/tab/misc-traditional/the-parting-glass-chords-1147884"
	tab, err := fetchTab(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", tab)
}

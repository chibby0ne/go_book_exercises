package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
    FirstIssueNumber = 1
	LatestIssueNumber     = 2244
	UrlFormatString = "https://xkcd.com/%v/info.0.json"
)

type Comic struct {
	Title      string
	Transcript string
	Url        string
}

type Comics []Comic

func downloadAllXkcd(filename string) error {
	var comics Comics
	for i := FirstIssueNumber; i <= LatestIssueNumber; i++ {
		if i == 404 { // xkcd being the comic it is, 404 doesn't exist :)
			continue
		}
		fmt.Printf("\rDownloading %v", i)
		url := fmt.Sprintf(UrlFormatString, strconv.FormatInt(int64(i), 10))
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			continue
		}
		defer resp.Body.Close()
		var comic Comic
		if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
			return err
		}
		comic.Url = url
		comics = append(comics, comic)
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&comics); err != nil {
		log.Fatal(err)
	}
	return nil
}

func findComic(match, filename string) (Comic, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var comics Comics
	if err := json.NewDecoder(file).Decode(&comics); err != nil {
		return Comic{}, err
	}
	for _, comic := range comics {
		if strings.Contains(comic.Title, match) {
			return comic, nil
		}
	}
	return Comic{}, nil
}

func main() {
	// General flag
	filenameFlag := flag.String("index", "index.json", "filename of the index")

	// Subcommands
	downloadSubcommand := flag.NewFlagSet("download", flag.ExitOnError)
	querySubcommand := flag.NewFlagSet("query", flag.ExitOnError)

	// query subcommand flags
	matchFlag := querySubcommand.String("match", "api", "word to find a comic match")

	if len(os.Args[1:]) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "download":
		downloadSubcommand.Parse(os.Args[2:])
	case "query":
		querySubcommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if downloadSubcommand.Parsed() {
		if err := downloadAllXkcd(*filenameFlag); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\nSuccessfully downloaded all comics and saved them into %v index file\n", *filenameFlag)
	} else {
		fmt.Printf("Searching comic with match string: %q in its title\n", *matchFlag)
		comic, err := findComic(*matchFlag, *filenameFlag)
		if err != nil {
			log.Fatal(err)
		}
		var emptyComic Comic
		if comic == emptyComic {
			fmt.Printf("No comic found with match string: %q\n", *matchFlag)
			os.Exit(1)
		}
		fmt.Printf("Title: %v\nURL: %v\nTranscript: %v\n", comic.Title, comic.Url, comic.Transcript)
	}
}

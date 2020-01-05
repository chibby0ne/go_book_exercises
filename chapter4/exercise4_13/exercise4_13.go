// The JSON-based websiervce of the Open Movie Database let's you search
// https://omdbapi.com/ for a movie by name and downlaod ists poster image.
// Write a tool poster that downloads the psoter image for a movie named on the
// command line.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	baseUrl = `https://omdbapi.com/?apikey=c672af33&t=%s`
)

type OmdbAPIResponse struct {
	Poster string `json:"Poster"`
}

func downloadPoster(movie *string) (string, error) {
	movieLowerCase := strings.ToLower(*movie)
	movieQuery := url.QueryEscape(movieLowerCase)
	url := fmt.Sprintf(baseUrl, movieQuery)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var response OmdbAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	imageResp, err := http.Get(response.Poster)
	if err != nil {
		return "", err
	}
	defer imageResp.Body.Close()
	extension := response.Poster[strings.LastIndex(response.Poster, "."):]
	movieNoSpaces := strings.ReplaceAll(movieLowerCase, " ", "_")
	filename := fmt.Sprintf("%s%s", movieNoSpaces, extension)
	poster, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(poster, imageResp.Body)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func main() {
	nameFlag := flag.String("name", "Jaws", "Name of the movie")
	flag.Parse()
	fmt.Printf("Downloading poster for %s\n", *nameFlag)
	filename, err := downloadPoster(nameFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\rDownloaded poster for movie %q to file: %q\n", *nameFlag, filename)
}

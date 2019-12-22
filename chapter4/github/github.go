package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	IssuesURL  = "https://api.github.com/search/issues"
	Regex      = `page=([0-9]+)>`
	PageNumber = 5
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func queryUrl(terms []string, page int64) (*http.Response, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q + "&page=" + strconv.FormatInt(page, 10))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	return resp, nil
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	resp, err := queryUrl(terms, 1)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	link := resp.Header["Link"][0] // Contains the link to the next page and last page
	reg := regexp.MustCompile(Regex)
	submatches := reg.FindAllStringSubmatch(link, -1)
	lastPage := submatches[1][1] // Get only the last page
	last, err := strconv.ParseInt(lastPage, 10, 64)
	if err != nil {
		return &result, nil
	}
	if last > PageNumber {
		last = PageNumber
	}
	for i := int64(2); i < last; i++ {
		response, err := queryUrl(terms, i)
		if err != nil {
			return &result, nil
		}
		var res IssuesSearchResult
		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			response.Body.Close()
			return &result, nil
		}
		response.Body.Close()
		result.Items = append(result.Items, res.Items...)
	}
	return &result, nil
}

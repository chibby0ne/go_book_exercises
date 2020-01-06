// Create a web server that queries GitHub once and then allows navigation of
// the list of bug reports, milestones and users.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var issueListTempl = template.Must(template.New("issuelist").Funcs(template.FuncMap{"daysAgo": daysAgo, "makeUrl": makeUrl}).Parse(`
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>Created</th>
  <th>User</th>
  <th>Title</th>
  <th>Milestone Title</th>
  <th>Milestone Description</th>
</tr>
{{range .Issues}}
<tr>
  <td><a href='{{.Number | makeUrl }}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td>{{.CreatedAt | daysAgo}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.Number | makeUrl }}'>{{.Title}}</a></td>
  {{if .Milestone }}
  <td><a href='{{.Milestone.HTMLURL}}'>{{.Milestone.Title}}</a></td>
  {{end}}
</tr>
{{end}}
</table>
`))

var issueTempl = template.Must(template.New("issue").Funcs(template.FuncMap{"daysAgo": daysAgo}).Parse(`
<h1>{{.Title}}</h1>
<dl>
	<dt>User</dt>
	<dd><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></dd>

	<dt>State</dt>
	<dd>{{.State}}</dd>

	<dt>Created</dt>
	<dd>{{.CreatedAt | daysAgo}}</dd>

    <dt>Issue Number</dt>
	<dd><a href='{{.HTMLURL}}'>{{.Number}}</a></dd>

    <dt>Milestone</dt>
    <dd><a href='{{.Milestone.HTMLURL}}'>{{.Milestone.Title}}</a></dd>

</dl>
<p>{{.Body}}</p>
`))

const (
	IssuesURL        = "https://api.github.com/repos/%s/%s/issues"
	Regex            = `<(https:\/\/api.github.com\/repositories\/[0-9]+/issues\?)page=([0-9]+)>`
	PageNumber       = 5
	UseMaxPageNumber = false
	Usage            = `
Usage: exercise4_14 ORGANIZATION REPOSITORY
       exercise4_14 -index INDEX

Serves all the issues of the given Gitub's organization/repository on a http server

When used with -index it will skip downloading all the issues from Github. This
should be a json file, containing an issueCache struct serialized.

The Github API limits unathenticated request to 60 per hour, that means that a
maximum of 60 pages will be reaches before the program is rate limited and
starts hosting the issues.
`
)

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Milestone struct {
	Title       string
	HTMLURL     string `json:"html_url"`
	Description string
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     template.HTML
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      template.HTML
	Milestone *Milestone
}

type IssueCache struct {
	Issues         []Issue
	IssuesByNumber map[int]Issue
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func makeUrl(issueNum int) string {
	return fmt.Sprintf("http://localhost:8080/%s", strconv.Itoa(issueNum))
}

func queryUrl(url string) (*http.Response, error) {
	fmt.Printf("\rReading page: %q", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	return resp, nil
}

func SearchIssues(organization, repository string) (*[]Issue, error) {
	resp, err := queryUrl(fmt.Sprintf(IssuesURL, organization, repository))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result []Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	reg := regexp.MustCompile(Regex)
	linkValue, ok := resp.Header["Link"]
	// If there's only a single page of result then we can return from function immediately
	if !ok {
		return &result, nil
	}
	link := linkValue[0] // Contains the link to the next page and last page
	submatches := reg.FindAllStringSubmatch(link, -1)
	if len(submatches) < 2 {
		return &result, nil
	}
	baseUrl, lastPage := submatches[1][1], submatches[1][2]
	last, err := strconv.ParseInt(lastPage, 10, 64)
	if err != nil {
		return &result, nil
	}
	for i := int64(2); i < last; i++ {
		nextPage := baseUrl + fmt.Sprintf("page=%d", i)
		response, err := queryUrl(nextPage)
		if err != nil {
			break
		}
		var res []Issue
		if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
			response.Body.Close()
			break
		}
		response.Body.Close()
		result = append(result, res...)
	}
	fmt.Println()
	return &result, nil
}

func NewIssueCache(issues *[]Issue) *IssueCache {
	var issueCache IssueCache
	issueCache.Issues = make([]Issue, len(*issues))
	copy(issueCache.Issues, *issues)
	issueCache.IssuesByNumber = make(map[int]Issue)
	for _, issue := range *issues {
		issueCache.IssuesByNumber[issue.Number] = issue
	}
	return &issueCache
}

func (ic *IssueCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprintf(os.Stdout, "Accessing issue list\n")
		issueListTempl.Execute(w, ic)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Use a issue number as path from root (/) URL\n")
		return
	}
	issueNum, err := strconv.Atoi(parts[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Use a issue number as path from root (/) URL\n")
		return
	}
	issue, ok := ic.IssuesByNumber[issueNum]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, fmt.Sprintf("Issue with number %d doesn't exist\n", issueNum))
		return
	}
	fmt.Fprintf(os.Stdout, "Accessing issue: #%d\n", issueNum)
	issueTempl.Execute(w, issue)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, Usage)
	}
	index := flag.String("index", "", "Backup with all the issues")
	flag.Parse()
	fmt.Println(flag.Args())
	if len(flag.Args()) == 0 && *index == "" || len(flag.Args()) != 2 && *index == "" || len(flag.Args()) == 2 && *index != "" {
		flag.Usage()
		os.Exit(1)
	}
	var issueCache *IssueCache
	if *index != "" {
		file, err := os.Open(*index)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Loading from file %q\n", *index)
		if err := json.NewDecoder(file).Decode(&issueCache); err != nil {
			log.Fatal(err)
		}
	} else {
		org, repo := strings.ToLower(flag.Args()[0]), strings.ToLower(flag.Args()[1])
		result, err := SearchIssues(org, repo)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("\rFetched all the issues for %s/%s\n", org, repo)
		issueCache = NewIssueCache(result)
		slice := [...]string{org, repo}
		file, err := os.Create(strings.Join(slice[:], "_") + ".json")
		if err != nil {
			log.Fatal(err)
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "   ")
		if err := encoder.Encode(&issueCache); err != nil {
			log.Fatal(err)
		}
		log.Printf("Stored the cache in file: %q\n", file.Name())
	}
	fmt.Println("Listening in localhost:8080...")
	http.Handle("/", issueCache)
	http.ListenAndServe(":8080", nil)
}

// Use the html/template package to replace printTracks with a functions taht
// displays the tracks as an HTML table. Use the solution to the previous
// exercise to arrange that each click on a column head makes an HTTP request
// to sort the table.
package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

type output int

const (
	lt output = iota
	eq
	gt
)

type TrackComparator func(a, b *Track) output

type MultiTierSort struct {
	tracks []*Track
	comps  []TrackComparator
}

func LessTitle(a, b *Track) output {
	switch {
	case a.Title < b.Title:
		return lt
	case a.Title > b.Title:
		return gt
	default:
		return eq
	}
}

func LessArtist(a, b *Track) output {
	switch {
	case a.Artist < b.Artist:
		return lt
	case a.Artist > b.Artist:
		return gt
	default:
		return eq
	}
}

func LessAlbum(a, b *Track) output {
	switch {
	case a.Album < b.Album:
		return lt
	case a.Album > b.Album:
		return gt
	default:
		return eq
	}
}

func LessYear(a, b *Track) output {
	switch {
	case a.Year < b.Year:
		return lt
	case a.Year > b.Year:
		return gt
	default:
		return eq
	}
}

func LessLength(a, b *Track) output {
	switch {
	case a.Length < b.Length:
		return lt
	case a.Length > b.Length:
		return gt
	default:
		return eq
	}
}

func (m *MultiTierSort) Len() int      { return len(m.tracks) }
func (m *MultiTierSort) Swap(i, j int) { *m.tracks[i], *m.tracks[j] = *m.tracks[j], *m.tracks[i] }
func (m *MultiTierSort) Less(i, j int) bool {
	for _, comp := range m.comps {
		out := comp(m.tracks[i], m.tracks[j])
		switch out {
		case lt:
			return true
		case gt:
			return false
		default:
			continue
		}
	}
	return false
}

var tracks = []*Track{
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("3m39s")},
	{"Ready to Go", "Martin Solveig", "Smash", 2011, length("3m27s")},
	{"Go Ahead 2", "Alicia Keys", "As I Am", 2007, length("3m38s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist }
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byAlbum []*Track

func (x byAlbum) Len() int           { return len(x) }
func (x byAlbum) Less(i, j int) bool { return x[i].Album < x[j].Album }
func (x byAlbum) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byLength []*Track

func (x byLength) Len() int           { return len(x) }
func (x byLength) Less(i, j int) bool { return x[i].Length < x[j].Length }
func (x byLength) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func printTracks(w http.ResponseWriter, req *http.Request) {
	switch req.FormValue("sort") {
	case "Title":
		sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessTitle}})
	case "Artist":
		sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessArtist}})
	case "Album":
		sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessAlbum}})
	case "Year":
		sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessYear}})
	case "Length":
		sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessLength}})
	}
	err := htmlTemplate.Execute(w, tracks)
	if err != nil {
		log.Printf("error executing the template: %s", err)
	}
}

var htmlTemplate = template.Must(template.New("htmlTemplate").Parse(`
<html>
<style>
table,
thead
tbody
tr {
    border: 1px solid #333;
}
</style>
<body>
<table>
    <thead>
        <tr>
            <td><a href="?sort=Title">Title</a></td>
            <td><a href="?sort=Artist">Artist</a></td>
            <td><a href="?sort=Album">Album</a></td>
            <td><a href="?sort=Year">Year</a></td>
            <td><a href="?sort=Length">Length</a></td>
        </tr>
    </thead>
{{range .}}
    <tbody>
        <tr>
            <td>{{.Title}}</td>
            <td>{{.Artist}}</td>
            <td>{{.Album}}</td>
            <td>{{.Year}}</td>
            <td>{{.Length}}</td>
        </tr>
    </tbody>
{{end}}
</table>
</body>
</html>
`))

func main() {
	log.Printf("HTTP server listening in :8080")
	http.HandleFunc("/", printTracks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

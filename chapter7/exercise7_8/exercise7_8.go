// Many GUIs provide a table widget with a stateful multi-tier sort: the
// primary sort key is the most recently clicked column head, the secondary
// sort key is the second-most recently clicked column head, and so on. Define
// an implementation of sort.Interface for use by such a table. Compare that
// approach with repeated sorting using sort.Stable

package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
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

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length", "Address")
	fmt.Fprintf(tw, format, "-----", "-----", "-----", "-----", "------", "--------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length, &t)
	}
	tw.Flush()
}

func copyTracks(dst, src *[]*Track) {
	for _, v := range *src {
		copyV := *v
		*dst = append(*dst, &copyV)
	}
}

func main() {
	var oldTracks []*Track
	copyTracks(&oldTracks, &tracks)
	fmt.Printf("Old tracks\n")
	printTracks(oldTracks)

	fmt.Printf("\n\nUnsorted\n")
	printTracks(tracks)

	fmt.Printf("\n\nSorted by Length\n")
	sort.Sort(&MultiTierSort{tracks, []TrackComparator{LessLength, LessAlbum}})
	printTracks(tracks)

	fmt.Printf("\n\nSorted by Length using sort stable\n")
	tracks = nil
	copyTracks(&tracks, &oldTracks)
	sort.Stable(byLength(tracks))
	sort.Stable(byAlbum(oldTracks))
	printTracks(tracks)

	fmt.Printf("\n\nSorted by Length using sort unstable\n")
	tracks = nil
	copyTracks(&tracks, &oldTracks)
	sort.Sort(byLength(tracks))
	sort.Sort(byAlbum(tracks))
	printTracks(tracks)

}

// Change the handler for /list ot print its output as an HTML table, not text
// You may find the html/template package useful.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

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
            <td>Item</td>
            <td>Price</td>
        </tr>
    </thead>
{{ range $item, $price := . }}
    <tbody>
        <tr>
            <td>{{$item}}</td>
            <td>{{$price}}</td>
        </tr>
    </tbody>
{{end}}
</table>
</body>
</html>
`))

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	sync.Mutex
	mapping map[string]dollars
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	htmlTemplate.Execute(w, db.mapping)
	// for item, price := range db.mapping {
	// 	fmt.Fprintf(w, "%s: %s\n", item, price)
	// }
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	price, ok := db.mapping[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "price cannot be parsed into a float: %q\n", price)
		return
	}
	_, ok := db.mapping[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	db.mapping[item] = dollars(priceFloat)
	fmt.Fprintf(w, "Updated item: %s to price: %s\n", item, dollars(priceFloat))
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "price not a float or int type: %q\n", price)
		return
	}
	db.mapping[item] = dollars(priceFloat)
	fmt.Fprintf(w, "Created item: %s with price: %s\n", item, dollars(priceFloat))
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	item := req.URL.Query().Get("item")
	_, ok := db.mapping[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	delete(db.mapping, item)
	fmt.Fprintf(w, "Deleted entry %s from database\n", item)
}

func main() {
	db := database{mapping: map[string]dollars{"shoes": 50, "socks": 5}}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

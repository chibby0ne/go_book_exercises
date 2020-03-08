// Write a web-based calculator program

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var index = template.Must(template.New("index").Parse(`
<html>
<body>
<h1>Simple Calculator</h1>
<p>Supports +, -, *, /, ^, pow, sin operations in the expression</p>
<form action="/calculate" method="post" class="form-example">
  <div class="form-example">
    <label for="expression">Enter an expression:</label>
    <input type="text" name="expression" id="expression" required>
  </div>
  <div class="form-example">
    <label for="environment">Enter the environment:</label>
    <input type="text" name="environment" id="environment" required>
  </div>
  <div class="form-example">
    <button>Calculate</button>
    <output name="result" for="expression environment">{{.Result}}</output>
  </div>
</form>
</body>
</html>
`))

func parseValuesToEnv(values string) (Env, error) {
	environment := make(Env)
	pairs := strings.Split(values, ",")
	for _, pair := range pairs {
		val := strings.Split(pair, "=")
		if len(val) != 2 {
			return nil, fmt.Errorf("Values of environment are not correctly formatted. Please use the following format: \"a\" = 1, \"b\" = 2")
		}
		value := strings.TrimSpace(val[1])
		variable := strings.TrimSpace(val[0])
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("Value %v cannot be converted to float", value)
		}
		environment[Var(variable)] = valueFloat
	}
	return environment, nil
}

func main() {
	http.HandleFunc("/calculate", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		exprString := r.FormValue("expression")
		expr, err := Parse(exprString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		envString := r.FormValue("environment")
		env, err := parseValuesToEnv(envString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := expr.Eval(env)
		if err = index.Execute(w, struct{ Result float64 }{Result: result}); err != nil {
			log.Fatal(err)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := index.Execute(w, struct{ Result float64 }{Result: 0}); err != nil {
			log.Fatal(err)
		}
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

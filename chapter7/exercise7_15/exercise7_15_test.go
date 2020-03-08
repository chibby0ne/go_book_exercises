package main

import (
	"fmt"
	"testing"
)

func TestMinimumValue(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"x ^ 3", Env{"x": 2}, "8"},
		{"4 ^ 3", Env{"x": 3}, "64"},
	}
	var prevExpr string
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n", test.expr, test.env, got, test.want)
		}
		extExpr, ok := expr.(extendedExpr)
		if !ok {
			t.Errorf("not able to type assert to extended expression")
			continue
		}
		minVal := extExpr.MinimumValue(test.env)
		fmt.Printf("\tExpr: %v, has MinimumValue of: %v\n", extExpr, minVal)
	}
}

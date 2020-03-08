package exercise7_14

import (
	"fmt"
	"math"
	"testing"
)

func TestEvalString(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
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
		// With expr.String() as input to Parse
		new_expr, err := Parse(expr.String())
		if err != nil {
			t.Error(err) // parse error
			continue
		}
		new_got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s (using Expr's String method as input to the Parse)\n", test.env, new_got)
		if new_got != got {
			t.Errorf("Expr after String() = %s, with Env %v = %q, old expr %q\n", new_expr, test.env, new_got, expr)
		}

	}
}

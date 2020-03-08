package exercise7_14

import (
	"fmt"
	"math"
	"strings"
)

type extendedExpr struct {
	op   rune
	x, y Expr
}

func (b extendedExpr) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	case '^':
		return math.Pow(b.x.Eval(env), b.y.Eval(env))
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (b extendedExpr) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/^", b.op) {
		return fmt.Errorf("unexpected binary op: %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b extendedExpr) String() string { return fmt.Sprintf("%v %c %v", b.x, b.op, b.y) }

// MinimumValue returns the minimum value of the two operands
func (b extendedExpr) MinimumValue(env Env) float64 {
	if b.x.Eval(env) < b.y.Eval(env) {
		return b.x.Eval(env)
	} else {
		return b.y.Eval(env)
	}
}

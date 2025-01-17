// Define a new concrete type that satisfies the Expr interface and provides a
// new operation such as computing the minimum value of its operands. Since the
// Parse function does not create instances of this new type, to use it you
// will need to construct a syntax tree directly (or extend the parser)

package main

import (
	"fmt"
	"math"
	"strings"
)

// An Expr is an arithmetic expression
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64

	// Check reports errors in this Expr and adds its Vars to the set
	Check(vars map[Var]bool) error

	// Pretty print syntax tree
	String() string
}

// A Var identifies a variable
type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) String() string { return string(v) }

// A literal is anumeric constant, e.g., 3.141
type literal float64

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (l literal) String() string { return fmt.Sprintf("%f", l) }

// A unary represents a unary operator expression, e.g., -x.
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported uanry operator: %q", u.op))
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op: %q", u.op)
	}
	return u.x.Check(vars)
}

func (u unary) String() string { return fmt.Sprintf("%c%v", u.op, u.x) }

// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
	op   rune // one of '+', '-', '*', '/'
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op: %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b binary) String() string { return fmt.Sprintf("%v %c %v", b.x, b.op, b.y) }

// A call represents a function call expression, e.g., sin(x).
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%v(", c.fn)
	fmt.Fprintf(&builder, "%v", c.args[0])
	for i := 1; i < len(c.args); i++ {
		fmt.Fprintf(&builder, ", %v", c.args[i])
	}
	fmt.Fprint(&builder, ")")
	return builder.String()
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

// To evualuate an expression containing variables, we need an environment that maps variables to values
type Env map[Var]float64

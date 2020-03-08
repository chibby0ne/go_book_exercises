// Write a program that reads a single expression from the standard input,
// prompts the user to provide values for any variables, then evaluates the
// expression in the resulting environment. Handle all errors gracefully.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	for {
		fmt.Println("Please provide a single expression")
		reader := bufio.NewReader(os.Stdin)
		exprString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		exprString = strings.TrimSpace(exprString)

		fmt.Println("Provide values for any variables in the form: x = 1, y = 2, etc...")
		valuesString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		valuesString = strings.TrimSpace(valuesString)
		env, err := parseValuesToEnv(valuesString)
		if err != nil {
			log.Fatal(err)
		}
		expr, err := Parse(exprString)
		if err != nil {
			log.Fatal(err)
		}
		result := expr.Eval(env)
		fmt.Println("Result is:", result)
	}
}

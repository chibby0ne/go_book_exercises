// Construct a tool that reports the set of all packages in the workspace that
// transitively depend on the packages specified by the arguments. Hint: you
// will need to run go list twice, once for the initial packages and once for
// all packages. You may want to parse its JSON output using the encoding/json
// package
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("You need to give at least one argument")
	}

	// Get the import paths of the packages given as arguments just in case the
	// arguments are not passed in resolvable unambiguous way
	args := []string{"list", "-f", "'{{ .ImportPath }}'"}
	for _, arg := range os.Args[1:] {
		args = append(args, arg)
	}
	claCommand := exec.Command("go", args...)
	claOutput, err := claCommand.Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not run command: %v:  %v", claCommand.String(), err))
	}
	inputPackages := strings.Fields(string(claOutput))

	// Make a map so that we can quickly check if the the one of this packages
	// is a transitive dependency
	packagesFromInput := make(map[string]bool)
	for _, inputPackage := range inputPackages {
		packagesFromInput[inputPackage] = true
	}

	// Get all packages in workspace, this is done through the use of the "..."
	// pattern. We need to report the set of packages and we could use the
	// package name we are using the ImportPath since it is unique while
	// package name is not, ex: math/rand  -> package rand, crypt/rand ->
	// package rand, their ImportPath is unique but their package name is not
	allPackagesCommand := exec.Command("go", []string{"list", "-f", `'{{ .ImportPath }} {{ join .Deps " " }}'`, "..."}...)
	allOutput, err := allPackagesCommand.Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not run command: %v, %v", allPackagesCommand.String(), err))
	}

	var dependantPackages []string
	scanner := bufio.NewScanner(bytes.NewReader(allOutput))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		importPath := fields[0]
		dependencies := fields[1:]
		for _, dependency := range dependencies {
			if packagesFromInput[dependency] {
				dependantPackages = append(dependantPackages, importPath)
			}
		}
	}

	for _, pkg := range dependantPackages {
		fmt.Println(pkg)
	}
}

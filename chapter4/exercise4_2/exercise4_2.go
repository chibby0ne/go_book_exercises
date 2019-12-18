package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	SHA256 = "256"
	SHA384 = "384"
	SHA512 = "512"
)

func generateHash(s string, hashSize string) string {
	var hash string
	switch hashSize {
	case SHA256:
		temp := sha256.Sum256([]byte(s))
		hash = fmt.Sprintf("%x", temp)
	case SHA384:
		temp := sha512.Sum384([]byte(s))
		hash = fmt.Sprintf("%x", temp)
	case SHA512:
		temp := sha512.Sum512([]byte(s))
		hash = fmt.Sprintf("%x", temp)
	}
	return hash
}

func handleFlags() (*string, error) {
	size := flag.String("size", SHA256, "Select SHA size: 256, 384 or 512")
	flag.Parse()
	if *size != SHA256 && *size != SHA384 && *size != SHA512 {
		return nil, fmt.Errorf("Wrong argument type: %v\n", *size)
	}
	return size, nil
}

func getInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	return strings.TrimRight(s, "\n"), err
}

func main() {
	size, err := handleFlags()
	if err != nil {
		log.Fatal(err)
	}
	input, err := getInput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(generateHash(input, *size))
}

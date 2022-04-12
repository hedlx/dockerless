package util

import (
	"log"
	"os"
	"strconv"
)

func GetStrVar(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatal(name, " is missing")
	}

	return v
}

func GetIntVar(name string) int {
	s := GetStrVar(name)
	v, err := strconv.Atoi(s)

	if err != nil {
		panic(err)
	}

	return v
}

package main

import "os"

type Environment struct {
	Pwd string `json:"pwd"`
}

func environment() *Environment {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return &Environment{Pwd: wd}
}

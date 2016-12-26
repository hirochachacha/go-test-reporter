package main

type LineCounts struct {
	Total   int `json:"total"`
	Covered int `json:"covered"`
	Missed  int `json:"missed"`
}

package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/cover"
)

type FileCoverage struct {
	Name            string      `json:"name"`
	BlobID          string      `json:"blob_id"`
	Coverage        string      `json:"coverage"`
	CoveredPercent  float64     `json:"covered_percent"`
	CoveredStrength float64     `json:"covered_strength"`
	LineCounts      *LineCounts `json:"line_counts"`
}

func fileCoverage(p *cover.Profile) *FileCoverage {
	path, err := findFile(p.FileName)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	h := sha1.New()
	h.Write([]byte("blob "))
	h.Write([]byte(strconv.Itoa(len(body))))
	h.Write([]byte{0})
	h.Write(body)
	blobID := hex.EncodeToString(h.Sum(nil))

	coverage := make([]interface{}, 1+bytes.Count(body, []byte{'\n'}))

	for _, block := range p.Blocks {
		for i := block.StartLine; i <= block.EndLine; i++ {
			count, _ := coverage[i-1].(int)
			coverage[i-1] = count + block.Count
		}
	}

	jcov, err := json.Marshal(coverage)
	if err != nil {
		panic(err)
	}

	var totalCount int

	l := &LineCounts{
		Total: len(coverage),
	}

	for _, count := range coverage {
		if count == nil || count.(int) == 0 {
			l.Missed++
		} else {
			l.Covered++
			totalCount += count.(int)
		}
	}

	var percent float64
	if l.Covered != 0 {
		percent = float64(l.Covered) / float64(l.Total)
	}

	var strength float64
	if l.Covered != 0 {
		strength = float64(totalCount) / float64(l.Covered)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rel, err := filepath.Rel(wd, path)
	if err != nil {
		panic(err)
	}

	return &FileCoverage{
		Name:            filepath.ToSlash(rel),
		BlobID:          blobID,
		Coverage:        string(jcov),
		LineCounts:      l,
		CoveredStrength: strength,
		CoveredPercent:  percent,
	}
}

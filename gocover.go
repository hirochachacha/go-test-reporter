// Original: github.com/goveralls/gocover.go
// License: https://mattn.mit-license.org/2016

package main

// Much of the core of this is copied from go's cover tool itself.

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The rest is written by Dustin Sallings

import (
	"fmt"
	"go/build"
	"log"
	"path/filepath"

	"golang.org/x/tools/cover"
)

func findFile(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}
	return filepath.Join(pkg.Dir, file), nil
}

// mergeProfs merges profiles for same target packages.
// It assumes each profiles have same sorted FileName and Blocks.
func mergeProfs(pfss [][]*cover.Profile) []*cover.Profile {
	// skip empty profiles ([no test files])
	for i := 0; i < len(pfss); i++ {
		if len(pfss[i]) > 0 {
			pfss = pfss[i:]
			break
		}
	}
	if len(pfss) == 0 {
		return nil
	} else if len(pfss) == 1 {
		return pfss[0]
	}
	head, rest := pfss[0], pfss[1:]
	ret := make([]*cover.Profile, 0, len(head))
	for i, profile := range head {
		for _, ps := range rest {
			if len(ps) == 0 {
				// no test files
				continue
			} else if len(ps) < i+1 {
				log.Fatal("Profile length is different")
			}
			if ps[i].FileName != profile.FileName {
				log.Fatal("Profile FileName is different")
			}
			profile.Blocks = mergeProfBlocks(profile.Blocks, ps[i].Blocks)
		}
		ret = append(ret, profile)
	}
	return ret
}

func mergeProfBlocks(as, bs []cover.ProfileBlock) []cover.ProfileBlock {
	if len(as) != len(bs) {
		log.Fatal("Two block length should be same")
	}
	// cover.ProfileBlock genereated by cover.ParseProfiles() is sorted by
	// StartLine and StartCol, so we can use index.
	ret := make([]cover.ProfileBlock, 0, len(as))
	for i, a := range as {
		b := bs[i]
		if a.StartLine != b.StartLine || a.StartCol != b.StartCol {
			log.Fatal("Blocks are not sorted")
		}
		a.Count += b.Count
		ret = append(ret, a)
	}
	return ret
}

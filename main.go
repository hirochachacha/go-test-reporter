package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/tools/cover"
)

var (
	env       map[string]string
	token     *string
	testFlags *string
)

func init() {
	env = make(map[string]string, 0)
	for _, e := range os.Environ() {
		if i := strings.IndexByte(e, '='); i != -1 {
			key := e[:i]
			val := e[i+1:]

			env[key] = val
		}
	}

	token = flag.String("token", env["CODECLIMATE_REPO_TOKEN"], "Code Climate repo token")
	testFlags = flag.String("testflags", "", "extra flags for go test")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go-test-reporter [coverprofile]")
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	args := flag.Args()

	var profile string
	var runAt int64
	var ps []*cover.Profile

	if len(args) == 0 {
		ps = getProfiles()
		runAt = time.Now().Unix()
	} else {
		profile = args[0]
		st, err := os.Stat(profile)
		if err != nil {
			panic(err)
		}
		runAt = st.ModTime().Unix()
		ps, err = cover.ParseProfiles(profile)
		if err != nil {
			panic(err)
		}
	}

	if len(ps) == 0 {
		return
	}

	l := new(LineCounts)

	var percent float64
	var strength float64

	sourceFiles := make([]*FileCoverage, len(ps))

	for i, p := range ps {
		fc := fileCoverage(p)

		percent += fc.CoveredPercent
		strength += fc.CoveredStrength

		l.Total += fc.LineCounts.Total
		l.Covered += fc.LineCounts.Covered
		l.Missed += fc.LineCounts.Missed

		sourceFiles[i] = fc
	}

	percent = round(percent / float64(len(sourceFiles)))
	strength = round(strength / float64(len(sourceFiles)))

	for _, fc := range sourceFiles {
		fc.CoveredPercent = round(fc.CoveredPercent)
		fc.CoveredStrength = round(fc.CoveredStrength)
	}

	api := &API{
		RepoToken:       *token,
		RunAt:           runAt,
		SourceFiles:     sourceFiles,
		CoveredPercent:  percent,
		CoveredStrength: strength,
		LineCounts:      l,
		Partial:         false,
		Git:             git(),
		CIService:       ci(),
		Environment:     environment(),
	}

	api.Post()
}

func getProfiles() []*cover.Profile {
	cmd := exec.Command("go", "list", "./...")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	pkgs := bytes.Split(output, []byte{'\n'})

	var tmp [][]*cover.Profile

	for _, pkg := range pkgs {
		pkg = bytes.TrimSpace(pkg)
		if len(pkg) != 0 && !bytes.Contains(pkg, []byte("/vendor/")) {
			if ps := getPackageProfiles(string(pkg)); len(ps) != 0 {
				tmp = append(tmp, ps)
			}
		}
	}

	return mergeProfs(tmp)
}

func getPackageProfiles(pkg string) []*cover.Profile {
	f, err := ioutil.TempFile("", "go-test-reporter")
	if err != nil {
		panic(err)
	}
	f.Close()

	defer os.Remove(f.Name())

	args := []string{"test", pkg, "-cover", "-coverpkg", "./...", "-coverprofile", f.Name()}
	if testFlags != nil {
		args = append(args, strings.Fields(*testFlags)...)
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	ps, err := cover.ParseProfiles(f.Name())
	if err != nil {
		panic(err)
	}

	return ps
}

func round(f float64) float64 {
	return float64(int(f*100+math.Copysign(0.5, f))) / 100
}

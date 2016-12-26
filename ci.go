package main

type Travis struct {
	Name            string `json:"name"`
	Branch          string `json:"branch"`
	BuildIdentifier string `json:"build_identifier"`
	PullRequest     string `json:"pull_request"`
}

type CircleCI struct {
	Name            string `json:"name"`
	Branch          string `json:"branch"`
	BuildIdentifier string `json:"build_identifier"`
	CommitSHA       string `json:"commit_sha"`
}

func ci() interface{} {
	switch {
	case env["TRAVIS"] != "":
		return &Travis{
			Name:            "travis-ci",
			Branch:          env["TRAVIS_BRANCH"],
			BuildIdentifier: env["TRAVIS_JOB_ID"],
			PullRequest:     env["TRAVIS_PULL_REQUEST"],
		}
	case env["CIRCLECI"] != "":
		return &CircleCI{
			Name:            "circleci",
			Branch:          env["CIRCLE_BRANCH"],
			BuildIdentifier: env["CIRCLE_BUILD_NUM"],
			CommitSHA:       env["CIRCLE_SHA1"],
		}
	}

	return nil
}

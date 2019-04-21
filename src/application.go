package main

import "os"

type applicationJSON struct {
	GitHash string `json:"git_hash"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

func generateApplicationJSON() applicationJSON {
	return applicationJSON{
		os.Getenv("GIT_HASH"),
		os.Getenv("VERSION"),
		os.Getenv("BUILD"),
	}
}

type applicationStatusJSON struct {
	Status string `json:"status"`
}

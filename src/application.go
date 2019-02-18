package main

type applicationJSON struct {
	Code    int    `json:"code"`
	GitHash string `json:"git_hash"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

func generateApplicationJSON applicationJSON {
	return applicationJSON{
		os.Getenv("CODE"),
		os.Getenv("GIT_HASH"),
		os.Getenv("VERSION"),
		os.Getenv("BUILD"),
	}
}

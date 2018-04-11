package main

type applicationJSON struct {
	Code    int    `json:"code"`
	Text    string `json:"text"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

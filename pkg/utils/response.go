package utils

import (
	"encoding/json"
	"regexp"
)

func ParseResponse(reviewResponse string) (*ResponseJSON, error) {

	var re = regexp.MustCompile(`(?ms)<final_json>(.+)<\/final_json>`)

	matches := re.FindAllStringSubmatch(reviewResponse, -1)

	if matches == nil {
		return nil, nil
	}

	var response ResponseJSON
	err := json.Unmarshal([]byte(matches[0][1]), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type ResponseJSON struct {
	Summary string   `json:"overarching_review_summary"`
	Details []Detail `json:"details"`
}

type Detail struct {
	CodeBlock string          `json:"code_block"`
	Comments  []ReviewComment `json:"comments"`
}

type ReviewComment struct {
	Type    string `json:"type"`
	Kind    string `json:"kind"`
	Comment string `json:"comment"`
}

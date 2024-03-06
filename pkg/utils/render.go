package utils

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func RenderReview(resp *ResponseJSON) {
	fmt.Println()
	if resp == nil {
		return
	}
	for _, detail := range resp.Details {
		fmt.Println(string(markdown.Render(fmt.Sprintf("```\n%s\n```", detail.CodeBlock), 80, 6)))
		for _, comment := range detail.Comments {
			com := fmt.Sprintf("*%s; %s* : %s", comment.Kind, comment.Type, comment.Comment)
			fmt.Println(string(markdown.Render(com, 80, 6)))
		}
	}
}

package prompts

import "fmt"

func ConstructDiffMessage(filename string, fromFileContents []byte, toFileContents []byte) string {
	return fmt.Sprintf(`
	Please help review the following file '%s'

	File Version in destination branch
	----
	%s
	----

	My Change that I want to merge
	----
	%s
	----
	`, filename, string(toFileContents), string(fromFileContents))
}

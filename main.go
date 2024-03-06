package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"peer_reviewer/pkg/llm"
	"peer_reviewer/pkg/prompts"
	"peer_reviewer/pkg/utils"

	"github.com/go-git/go-git/v5"
)

func main() {
	var sourceBranch string
	var destinationBranch string
	var projectDirectoryPath string
	var llmProvider string
	var model string
	var temperature float64

	flag.StringVar(&sourceBranch, "source", "", "Source branch from which the PR/MR is raised")
	flag.StringVar(&destinationBranch, "dest", "", "Destination branch into which the merge is requested")
	flag.StringVar(&projectDirectoryPath, "project-dir", "", "Local project directory path")
	flag.StringVar(&llmProvider, "provider", "anthropic", "The LLM to use")
	flag.StringVar(&model, "model", "claude-2.1", "The model to use")
	flag.Float64Var(&temperature, "temperature", 0.1, "The temperature for the model")
	flag.Parse()

	if sourceBranch == "" || destinationBranch == "" || projectDirectoryPath == "" {
		fmt.Println("Please make sure to key in values for source & destination branches and the project directory root path")
		os.Exit(-1)
	}

	provider, err := llm.GetLLM(llmProvider)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Open repository
	r, err := git.PlainOpen(projectDirectoryPath)
	if err != nil {
		fmt.Printf("Unable to open repository %s\n", err)
		os.Exit(-1)
	}

	latestSourceBranchCommit, err := utils.GetLatestCommit(r, sourceBranch)
	if err != nil {
		fmt.Printf("Unable to fetch the latest commit to the source branch %s\n", err)
		os.Exit(-1)
	}

	latestDestinationBranchCommit, err := utils.GetLatestCommit(r, destinationBranch)
	if err != nil {
		fmt.Printf("Unable to fetch the latest commit to the destination branch %s\n", err)
		os.Exit(-1)
	}

	patch, err := latestSourceBranchCommit.Patch(latestDestinationBranchCommit)
	if err != nil {
		fmt.Printf("Unable to generate a patch %s\n", err)
		os.Exit(-1)
	}

	ctx := context.Background()
	for _, fp := range patch.FilePatches() {
		from, to := fp.Files()
		fromFileContents, err := utils.GetContentsOfFile(r, from.Hash())
		if err != nil {
			fmt.Printf("Unexpected error %s\n", err)
			os.Exit(-1)
		}

		toFileContents, err := utils.GetContentsOfFile(r, to.Hash())
		if err != nil {
			fmt.Printf("Unexpected error %s\n", err)
			os.Exit(-1)
		}

		fmt.Printf("Reviewing file %s...\n", from.Path())
		completion, err := provider.GetCompletions(ctx, model, temperature, []llm.Message{
			{
				Role:    "user",
				Content: prompts.ConstructDiffMessage(from.Path(), fromFileContents, toFileContents),
			},
		}, prompts.SystemMessage)
		if err != nil {
			fmt.Printf("Error from LLM %s\n", err)
			os.Exit(-1)
		}

		response, err := utils.ParseResponse(completion.Content)
		if err != nil {
			fmt.Printf("Could not parse response from LLM %s\n", err)
			os.Exit(-1)
		}

		// fmt.Printf("Response is %+v\n", response)
		utils.RenderReview(response)
	}

}

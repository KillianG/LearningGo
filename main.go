package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/v42/github"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

func getInput(phrase string) string {
	fmt.Print(phrase)
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	res, err := reader.ReadString('\n')
	println()
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return ""
	}
	res = strings.TrimSuffix(res, "\n")
	return res
}

func main() {
	fmt.Print("GitHub Token: ")
	byteToken, _ := terminal.ReadPassword(int(syscall.Stdin))
	println()
	usertoCheck := getInput("Username to check PRs: ")
	repo := getInput("Repository to check PRs: ")
	token := string(byteToken)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	nbPr := 100
	var openPrsList []*github.PullRequest
	page := 0
	fmt.Print("Loading")
	for nbPr == 100 {
		opts := github.PullRequestListOptions{
			State:       "all",
			Head:        "",
			Base:        "",
			Sort:        "",
			Direction:   "",
			ListOptions: github.ListOptions{Page: page, PerPage: 100},
		}
		page += 1
		openPrs, _, err := client.PullRequests.List(ctx, "scality", repo, &opts)
		openPrsList = append(openPrsList, openPrs...)
		if err != nil {
			fmt.Printf("\nerror: %v\n", err)
			return
		}
		fmt.Print(".")
		nbPr = len(openPrs)
	}
	println()

	users := make(map[string]int)
	var totalTime time.Duration
	for _, pr := range openPrsList {
		users[pr.User.GetLogin()] += 1
		created := pr.GetCreatedAt()
		closed := pr.GetClosedAt()
		diff := closed.Sub(created)
		if *pr.State == "closed" {
			totalTime += diff
		}
	}
	averageTimeToClosePr := totalTime.Minutes() / float64(len(openPrsList)) / 60
	fmt.Printf("Usually a PR takes %.2f hours to close in %s\n", averageTimeToClosePr, repo)
	fmt.Printf("User %s openned %d PR in %s\n", usertoCheck, users[usertoCheck], repo)
}

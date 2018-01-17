package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	defaultBaseURL = "https://api.github.com"
)

type argsType struct {
	owner              string
	repo               string
	perPage            int
	insecureSkipVerify bool
	baseURL            string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	ctx := context.Background()

	client, err := newGithubClient(ctx, args.baseURL, args.insecureSkipVerify)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	opts := &github.PullRequestListOptions{
		State:     "all",
		Sort:      "updated",
		Direction: "desc",
	}
	opts.PerPage = args.perPage

	stats, err := collectStats(ctx, client, args.owner, args.repo, opts)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Printf("user\tdone\treviewed\n")
	for _, w := range stats {
		fmt.Printf("%s\t%d\t%d\n", w.user, w.sentPullRequests, w.reviewedPullRequests)
	}
}

func parseArgs() (*argsType, error) {
	args := &argsType{}
	flag.StringVar(&args.owner, "owner", "", "owner of repo")
	flag.StringVar(&args.repo, "repo", "", "repo name")
	flag.IntVar(&args.perPage, "per-page", 10, "count of pull requests to scan")
	flag.BoolVar(&args.insecureSkipVerify, "insecure-skip-verify", false, "skip verification of cert")
	flag.StringVar(&args.baseURL, "base-url", defaultBaseURL, "custom GitHub base URL if you use GitHub Enterprise")
	flag.Parse()

	if args.owner == "" {
		return nil, fmt.Errorf("owner cannnot be empty")
	}

	if args.repo == "" {
		return nil, fmt.Errorf("repo cannot be empty")
	}

	if args.perPage <= 0 {
		return nil, fmt.Errorf("per-page should be positive")
	}

	return args, nil
}

func collectStats(ctx context.Context, client *github.Client, owner, repo string, listOpts *github.PullRequestListOptions) ([]*workloadStat, error) {
	prs, _, err := client.PullRequests.List(ctx, owner, repo, listOpts)
	if err != nil {
		return nil, err
	}

	statsByUser := make(map[string]*workloadStat)

	for _, pr := range prs {
		if assignee := pr.GetAssignee(); assignee != nil {
			if _, ok := statsByUser[assignee.GetLogin()]; !ok {
				statsByUser[assignee.GetLogin()] = &workloadStat{
					user: assignee.GetLogin(),
				}
			}
			statsByUser[assignee.GetLogin()].sentPullRequests++
		}

		reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, pr.GetNumber(), nil)
		if err != nil {
			continue
		}

		for _, review := range reviews {
			if review.GetState() == "COMMENTED" {
				continue
			}
			if _, ok := statsByUser[review.User.GetLogin()]; !ok {
				statsByUser[review.User.GetLogin()] = &workloadStat{
					user: review.User.GetLogin(),
				}
			}
			statsByUser[review.User.GetLogin()].reviewedPullRequests++
		}
	}

	stats := make([]*workloadStat, 0, len(statsByUser))
	for _, w := range statsByUser {
		stats = append(stats, w)
	}

	return stats, nil
}

func newGithubClient(ctx context.Context, baseURL string, insecureSkipVerify bool) (*github.Client, error) {
	ghToken, err := getGithubToken()
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: ghToken,
		},
	)
	tc := oauth2.NewClient(ctx, ts)
	if tct, ok := tc.Transport.(*oauth2.Transport); ok && insecureSkipVerify {
		tct.Base = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if baseURL == defaultBaseURL {
		client := github.NewClient(tc)
		return client, nil
	}

	client, err := github.NewEnterpriseClient(baseURL, baseURL, tc)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getGithubToken() (string, error) {
	token := os.Getenv("GITHUB_API_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_API_TOKEN must be provided")
	}
	return token, nil
}

type workloadStat struct {
	user                 string
	sentPullRequests     int
	reviewedPullRequests int
}

func (w *workloadStat) ratio() float32 {
	return float32(w.reviewedPullRequests) / float32(w.sentPullRequests)
}

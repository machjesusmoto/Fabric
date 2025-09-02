package internal

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/danielmiessler/fabric/cmd/generate_changelog/internal/cache"
	"github.com/danielmiessler/fabric/cmd/generate_changelog/internal/config"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

type ReleaseManager struct {
	cache       *cache.Cache
	githubToken string
	owner       string
	repo        string
}

// getGitHubInfo extracts owner and repo from git remote origin URL
func getGitHubInfo() (owner, repo string, err error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	url := strings.TrimSpace(string(output))

	// Handle both SSH and HTTPS URLs
	// SSH: git@github.com:owner/repo.git
	// HTTPS: https://github.com/owner/repo.git
	var re *regexp.Regexp
	if strings.HasPrefix(url, "git@") {
		re = regexp.MustCompile(`git@github\.com:([^/]+)/([^/.]+)(?:\.git)?`)
	} else {
		re = regexp.MustCompile(`https://github\.com/([^/]+)/([^/.]+)(?:\.git)?`)
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	return matches[1], matches[2], nil
}

func NewReleaseManager(cfg *config.Config) (*ReleaseManager, error) {
	cache, err := cache.New(cfg.CacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	owner, repo, err := getGitHubInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub repository info: %w", err)
	}

	return &ReleaseManager{
		cache:       cache,
		githubToken: cfg.GitHubToken,
		owner:       owner,
		repo:        repo,
	}, nil
}

func (rm *ReleaseManager) Close() error {
	return rm.cache.Close()
}

func (rm *ReleaseManager) UpdateReleaseDescription(version string) error {
	versions, err := rm.cache.GetVersions()
	if err != nil {
		return fmt.Errorf("failed to get versions from cache: %w", err)
	}

	versionData, exists := versions[version]
	if !exists {
		return fmt.Errorf("version %s not found in versions table", version)
	}

	if versionData.AISummary == "" {
		return fmt.Errorf("ai_summary is empty for version %s", version)
	}

	releaseBody := fmt.Sprintf("## Changes\n\n%s", versionData.AISummary)

	ctx := context.Background()
	var client *github.Client

	if rm.githubToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: rm.githubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	// Check if current repository is a fork by getting repo details
	repo, _, err := client.Repositories.Get(ctx, rm.owner, rm.repo)
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	// If repository is a fork, try updating the upstream (parent) repository first
	if repo.Parent != nil {
		parentOwner := repo.Parent.Owner.GetLogin()
		parentRepo := repo.Parent.GetName()

		fmt.Printf("Repository is a fork of %s/%s, attempting to update upstream release...\n", parentOwner, parentRepo)

		err := rm.updateReleaseForRepo(ctx, client, parentOwner, parentRepo, version, releaseBody)
		if err == nil {
			fmt.Printf("Successfully updated release description for %s in upstream repository %s/%s\n", version, parentOwner, parentRepo)
			return nil
		}

		fmt.Printf("Failed to update upstream repository: %v\nFalling back to current repository...\n", err)
	}

	// Update current repository (either not a fork or upstream update failed)
	err = rm.updateReleaseForRepo(ctx, client, rm.owner, rm.repo, version, releaseBody)
	if err != nil {
		return fmt.Errorf("failed to update release description for version %s in repository %s/%s: %w", version, rm.owner, rm.repo, err)
	}

	fmt.Printf("Successfully updated release description for %s in repository %s/%s\n", version, rm.owner, rm.repo)
	return nil
}

func (rm *ReleaseManager) updateReleaseForRepo(ctx context.Context, client *github.Client, owner, repo, version, releaseBody string) error {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, version)
	if err != nil {
		return fmt.Errorf("failed to get release for version %s: %w", version, err)
	}

	release.Body = &releaseBody
	_, _, err = client.Repositories.EditRelease(ctx, owner, repo, *release.ID, release)
	if err != nil {
		return fmt.Errorf("failed to update release description for version %s: %w", version, err)
	}

	return nil
}

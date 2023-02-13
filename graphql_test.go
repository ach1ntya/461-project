package main

import (
	"testing"
)

func TestGitHubGraphQL(t *testing.T) {
	repoName := "461-project"
	owner := "ach1ntya"
	expectedIssueCount := 5
	expectedReleaseCount := 3
	expectedStarCount := 10
	expectedLicense := "MIT"

	issueCount, releaseCount, starCount, license := gitHubGraphQL(repoName, owner)

	if issueCount != expectedIssueCount {
		t.Errorf("Expected issue count to be %d, but got %d", expectedIssueCount, issueCount)
	}

	if releaseCount != expectedReleaseCount {
		t.Errorf("Expected release count to be %d, but got %d", expectedReleaseCount, releaseCount)
	}

	if starCount != expectedStarCount {
		t.Errorf("Expected star count to be %d, but got %d", expectedStarCount, starCount)
	}

	if license != expectedLicense {
		t.Errorf("Expected license to be %s, but got %s", expectedLicense, license)
	}
}

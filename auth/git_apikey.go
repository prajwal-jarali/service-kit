package auth

import (
	"context"
	"fmt"
)

type GitApiKey struct {
	token string
}

func NewGitApiKey(ghpToken string) (*GitApiKey, error) {
	if ghpToken == "" {
		return nil, fmt.Errorf("github token cannot be empty")
	}

	return &GitApiKey{
		token: "Bearer " + ghpToken,
	}, nil
}

func (g *GitApiKey) GetToken(ctx context.Context) (string, error) {
	return g.token, nil
}

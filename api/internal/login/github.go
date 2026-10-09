package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Token struct {
	Access  string    `json:"a"`
	Refresh string    `json:"r,omitempty"`
	Expires time.Time `json:"e,omitzero"`
}

type githubApp struct {
	clientID     string
	clientSecret string
	redirectURL  string
	authorizeURL string
	tokenURL     string
	client       *http.Client
	now          func() time.Time
}

func newGitHubApp(clientID, clientSecret, redirectURL string) *githubApp {
	return &githubApp{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		authorizeURL: "https://github.com/login/oauth/authorize",
		tokenURL:     "https://github.com/login/oauth/access_token",
		client:       &http.Client{Timeout: 30 * time.Second},
		now:          time.Now,
	}
}

func (g *githubApp) authorize(state, challenge string) string {
	query := url.Values{
		"client_id":             {g.clientID},
		"redirect_uri":          {g.redirectURL},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	return g.authorizeURL + "?" + query.Encode()
}

func (g *githubApp) exchange(ctx context.Context, code, verifier string) (Token, error) {
	return g.token(ctx, url.Values{
		"code":          {code},
		"redirect_uri":  {g.redirectURL},
		"code_verifier": {verifier},
	})
}

func (g *githubApp) refresh(ctx context.Context, refreshToken string) (Token, error) {
	return g.token(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	})
}

func (g *githubApp) token(ctx context.Context, form url.Values) (Token, error) {
	form.Set("client_id", g.clientID)
	form.Set("client_secret", g.clientSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return Token{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := g.client.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer func() { _ = res.Body.Close() }()
	var body struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
		Description  string `json:"error_description"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return Token{}, fmt.Errorf("github token response: %w", err)
	}
	if body.Error != "" {
		return Token{}, fmt.Errorf("github refused the token request: %s: %s", body.Error, body.Description)
	}
	if body.AccessToken == "" {
		return Token{}, errors.New("github sent no access token")
	}
	token := Token{Access: body.AccessToken, Refresh: body.RefreshToken}
	if body.ExpiresIn > 0 {
		token.Expires = g.now().Add(time.Duration(body.ExpiresIn) * time.Second)
	}
	return token, nil
}

// Copyright (c) 2026 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

package imagefactory

import (
	"context"
	"fmt"
	"net/url"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/siderolabs/omni/internal/pkg/config"
)

// TokenSource fetches and caches an Auth0 client-credentials token for the
// Image Factory Enterprise service.
//
// It wraps golang.org/x/oauth2's ReuseTokenSource so caching and expiry-based
// refresh are handled automatically. No background goroutine is required.
type TokenSource struct {
	source oauth2.TokenSource
	logger *zap.Logger
}

// NewTokenSource creates a TokenSource for the Omni M2M credentials.
// It performs an initial token fetch to verify the credentials; an error is
// returned if the fetch fails.
func NewTokenSource(cfg config.RegistriesImageFactoryAuth0, logger *zap.Logger) (*TokenSource, error) {
	return newTokenSource(cfg.ClientID, cfg.ClientSecret, cfg.Domain, cfg.Audience, cfg.OrgID, logger)
}

// NewNodeTokenSource creates a TokenSource using the node M2M credentials from cfg.
// Node tokens carry the same org_id as Omni's token but have narrower scopes
// (Image Factory API only, no Management API). They are intended for injection
// into Talos machine configs; configure a longer token lifetime in Auth0 (e.g.
// 30 days) to reduce machine config push frequency.
func NewNodeTokenSource(cfg config.RegistriesImageFactoryAuth0, logger *zap.Logger) (*TokenSource, error) {
	if cfg.NodeClientID == nil || cfg.NodeClientSecret == nil {
		return nil, fmt.Errorf("node client ID and secret are required for node token source")
	}

	return newTokenSource(*cfg.NodeClientID, *cfg.NodeClientSecret, cfg.Domain, cfg.Audience, cfg.OrgID, logger)
}

func newTokenSource(clientID, clientSecret, domain, audience, orgID string, logger *zap.Logger) (*TokenSource, error) {
	ccfg := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://" + domain + "/oauth/token",
		EndpointParams: url.Values{
			"audience":     {audience},
			"organization": {orgID},
		},
	}

	// context.Background() is intentional: token fetches must outlive any
	// individual request context.
	source := oauth2.ReuseTokenSource(nil, ccfg.TokenSource(context.Background()))

	// Verify credentials at startup so misconfiguration surfaces early.
	if _, err := source.Token(); err != nil {
		return nil, fmt.Errorf("initial Auth0 token fetch failed: %w", err)
	}

	return &TokenSource{source: source, logger: logger}, nil
}

// Token returns the current access token string.
// It may block briefly if the cached token is expired and a new one must be fetched.
// Returns "" on error (the error is logged).
func (ts *TokenSource) Token() string {
	tok, err := ts.source.Token()
	if err != nil {
		ts.logger.Error("failed to fetch Auth0 access token for image factory", zap.Error(err))
		return ""
	}

	return tok.AccessToken
}

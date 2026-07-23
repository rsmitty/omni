// Copyright (c) 2026 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

// Package imagefactoryauth provides helpers to build a Talos RegistryAuthConfig
// document for authenticating against the configured Omni image factory.
package imagefactoryauth

import (
	"fmt"
	"net/url"

	"github.com/siderolabs/talos/pkg/machinery/config/types/cri"

	omnicfg "github.com/siderolabs/omni/internal/pkg/config"
)

// TokenSource provides the current access token for the image factory.
// A nil TokenSource means Auth0 is not configured; username/password auth is used instead.
type TokenSource interface {
	Token() string
}

// BuildDoc returns a RegistryAuthConfig doc for the image factory based on the
// configured credentials. Returns nil if no credentials are configured.
//
// When ts is non-nil, the current Bearer token is used as the password (with a
// placeholder username, which image-factory ignores). When ts is nil, the static
// username and password from registries config are used.
//
//nolint:nilnil
func BuildDoc(registries omnicfg.Registries, ts TokenSource) (*cri.RegistryAuthConfigV1Alpha1, error) {
	var username, password string

	if ts != nil {
		token := ts.Token()
		if token == "" {
			// Token not yet available; skip injecting registry auth for now.
			return nil, nil
		}

		username = "omni"
		password = token
	} else {
		username = registries.GetImageFactoryUsername()
		password = registries.GetImageFactoryPassword()

		if username == "" || password == "" {
			return nil, nil
		}
	}

	u, err := url.Parse(registries.GetImageFactoryBaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse image factory base URL: %w", err)
	}

	doc := cri.NewRegistryAuthConfigV1Alpha1(u.Host)
	doc.RegistryUsername = username
	doc.RegistryPassword = password

	return doc, nil
}

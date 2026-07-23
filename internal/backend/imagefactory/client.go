// Copyright (c) 2026 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

package imagefactory

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/siderolabs/image-factory/pkg/client"
	"github.com/siderolabs/image-factory/pkg/schematic"
)

// requestTimeout caps every call to the image factory. Without it, the factory's HTTP client has
// no internal timeout, so a stuck connection would pin a controller reconcile slot indefinitely.
// 30 minutes as the scan report request can be particularly slow.
const requestTimeout = 30 * time.Minute

// serverSnifferTransport wraps an http.RoundTripper and records whether the image factory
// identifies itself as an Enterprise instance via the Server response header.
// It captures the header from the first successful response so no extra requests are needed.
type serverSnifferTransport struct {
	wrapped      http.RoundTripper
	detected     atomic.Bool
	isEnterprise atomic.Bool
}

// bearerTransport injects an Authorization: Bearer header on every request using
// the current token from a TokenSource.
type bearerTransport struct {
	wrapped http.RoundTripper
	ts      *TokenSource
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if token := t.ts.Token(); token != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return t.wrapped.RoundTrip(req)
}

func (t *serverSnifferTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.wrapped.RoundTrip(req)
	if err == nil && t.detected.CompareAndSwap(false, true) {
		t.isEnterprise.Store(strings.Contains(resp.Header.Get("Server"), "Enterprise"))
	}

	return resp, err
}

// Client is the image factory client.
type Client struct {
	*client.Client

	sniffer *serverSnifferTransport
	host    string
}

// NewClient creates a new image factory client with optional Basic Auth credentials.
func NewClient(imageFactoryBaseURL, username, password string) (*Client, error) {
	var opts []client.Option

	if username != "" && password != "" {
		opts = append(opts, client.WithBasicAuth(username, password))
	}

	return newClient(imageFactoryBaseURL, nil, opts...)
}

// NewClientWithTokenSource creates a new image factory client that authenticates
// using Bearer tokens from the provided TokenSource.
func NewClientWithTokenSource(imageFactoryBaseURL string, ts *TokenSource) (*Client, error) {
	return newClient(imageFactoryBaseURL, ts)
}

func newClient(imageFactoryBaseURL string, ts *TokenSource, extraOpts ...client.Option) (*Client, error) {
	sniffer := &serverSnifferTransport{wrapped: http.DefaultTransport}

	var transport http.RoundTripper = sniffer
	if ts != nil {
		transport = &bearerTransport{wrapped: sniffer, ts: ts}
	}

	opts := append([]client.Option{
		client.WithClient(http.Client{Transport: transport, Timeout: requestTimeout}),
	}, extraOpts...)

	factoryClient, err := client.New(imageFactoryBaseURL, opts...)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(imageFactoryBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image factory base URL %q: %w", imageFactoryBaseURL, err)
	}

	return &Client{
		Client:  factoryClient,
		host:    baseURL.Host,
		sniffer: sniffer,
	}, nil
}

// Host returns the host of the image factory client.
func (cli *Client) Host() string {
	if cli == nil {
		return ""
	}

	return cli.host
}

// CachedIsEnterprise reports whether the connected image factory is an Enterprise instance.
// The value is detected from the Server response header of the first successful HTTP response
// and cached; it returns false until at least one response has been received.
func (cli *Client) CachedIsEnterprise() bool {
	return cli.sniffer.isEnterprise.Load()
}

// EnsureSchematic uploads the given schematic to the image factory and returns its ID
// along with the normalized schematic as the factory persisted it.
//
// The factory deduplicates by content: if the same schematic was uploaded before, it returns
// the existing ID without creating a new one.
func (cli *Client) EnsureSchematic(ctx context.Context, inputSchematic schematic.Schematic) (string, *schematic.Schematic, error) {
	id, data, err := cli.SchematicCreate(ctx, inputSchematic)
	if err != nil {
		return "", nil, fmt.Errorf("failed to ensure schematic: %w", err)
	}

	return id, data, nil
}

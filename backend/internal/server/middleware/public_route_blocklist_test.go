package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestLoadPublicRouteBlocklistDefault(t *testing.T) {
	t.Setenv("PUBLIC_ROUTE_BLOCKLIST_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
	t.Setenv("DATA_DIR", "")

	list, err := LoadPublicRouteBlocklist()
	require.NoError(t, err)
	require.True(t, list.Enabled)
	require.True(t, list.UsingDefault)
	require.False(t, list.Matches("/register"))
	require.False(t, list.Matches("/email-verify"))
	require.True(t, list.Matches("/api/v1/settings/public"))
	require.True(t, list.Matches("/api/v1/auth/oauth/github/start"))
	require.False(t, list.Matches("/login"))
}

func TestLoadPublicRouteBlocklistFromFile(t *testing.T) {
	path := writeBlocklistFile(t, `
enabled: false
rules:
  - match: exact
    path: /hidden
  - match: exact
    path: /hidden
  - match: prefix
    path: /private/
`)
	t.Setenv("PUBLIC_ROUTE_BLOCKLIST_FILE", path)
	t.Setenv("DATA_DIR", "")

	list, err := LoadPublicRouteBlocklist()
	require.NoError(t, err)
	require.False(t, list.Enabled)
	require.False(t, list.UsingDefault)
	require.Equal(t, path, list.Source)
	require.Len(t, list.Rules, 2)
	require.False(t, list.Matches("/hidden"))
}

func TestLoadPublicRouteBlocklistRejectsInvalidRules(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "bad_match",
			body: `
enabled: true
rules:
  - match: glob
    path: /hidden
`,
		},
		{
			name: "bad_path",
			body: `
enabled: true
rules:
  - match: exact
    path: hidden
`,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := writeBlocklistFile(t, tc.body)
			t.Setenv("PUBLIC_ROUTE_BLOCKLIST_FILE", path)
			t.Setenv("DATA_DIR", "")

			_, err := LoadPublicRouteBlocklist()
			require.Error(t, err)
		})
	}
}

func TestPublicRouteBlocklistMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	list, err := normalizePublicRouteBlocklist(publicRouteBlocklistFile{
		Enabled: boolPtr(true),
		Rules: []PublicRouteBlocklistRule{
			{Match: "exact", Path: "/hidden"},
			{Match: "prefix", Path: "/private/"},
		},
	}, "test", false)
	require.NoError(t, err)

	router := gin.New()
	router.Use(PublicRouteBlocklistMiddleware(list))
	router.Any("/*path", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	tests := []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{path: "/hidden?x=1", wantStatus: http.StatusNotFound, wantBody: "404 Not Found"},
		{path: "/private/value", wantStatus: http.StatusNotFound, wantBody: "404 Not Found"},
		{path: "/setup", wantStatus: http.StatusNotFound, wantBody: "404 Not Found"},
		{path: "/setup/status", wantStatus: http.StatusNotFound, wantBody: "404 Not Found"},
		{path: "/login", wantStatus: http.StatusOK, wantBody: "ok"},
		{path: "/v1/usage", wantStatus: http.StatusOK, wantBody: "ok"},
		{path: "/assets/app.js", wantStatus: http.StatusOK, wantBody: "ok"},
		{path: "/static/app/res/app.js", wantStatus: http.StatusOK, wantBody: "ok"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			router.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
			require.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestPublicRouteBlocklistMiddlewareDisabledStillHidesSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	list, err := normalizePublicRouteBlocklist(publicRouteBlocklistFile{
		Enabled: boolPtr(false),
		Rules: []PublicRouteBlocklistRule{
			{Match: "exact", Path: "/hidden"},
		},
	}, "test", false)
	require.NoError(t, err)

	router := gin.New()
	router.Use(PublicRouteBlocklistMiddleware(list))
	router.Any("/*path", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/hidden", nil))
	require.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/setup/status", nil))
	require.Equal(t, http.StatusNotFound, w.Code)
}

func writeBlocklistFile(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "public-route-blocklist.yaml")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}

func boolPtr(v bool) *bool {
	return &v
}

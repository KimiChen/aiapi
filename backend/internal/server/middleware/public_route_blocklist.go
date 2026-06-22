package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const publicRouteBlocklistFilename = "public-route-blocklist.yaml"

type PublicRouteBlocklist struct {
	Enabled      bool
	Rules        []PublicRouteBlocklistRule
	Source       string
	UsingDefault bool
}

type PublicRouteBlocklistRule struct {
	Match string `yaml:"match"`
	Path  string `yaml:"path"`
}

type publicRouteBlocklistFile struct {
	Enabled *bool                      `yaml:"enabled"`
	Rules   []PublicRouteBlocklistRule `yaml:"rules"`
}

func LoadPublicRouteBlocklist() (*PublicRouteBlocklist, error) {
	candidates := publicRouteBlocklistCandidatePaths()
	for _, path := range candidates {
		raw, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read public route blocklist %q: %w", path, err)
		}

		var cfg publicRouteBlocklistFile
		if err := yaml.Unmarshal(raw, &cfg); err != nil {
			return nil, fmt.Errorf("parse public route blocklist %q: %w", path, err)
		}
		list, err := normalizePublicRouteBlocklist(cfg, path, false)
		if err != nil {
			return nil, err
		}
		return list, nil
	}

	list, err := normalizePublicRouteBlocklist(defaultPublicRouteBlocklistFile(), "embedded defaults", true)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func PublicRouteBlocklistMiddleware(list *PublicRouteBlocklist) gin.HandlerFunc {
	if list == nil {
		list, _ = normalizePublicRouteBlocklist(defaultPublicRouteBlocklistFile(), "embedded defaults", true)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "" {
			path = "/"
		}

		if isCompletedSetupPath(path) {
			writePlainNotFound(c)
			return
		}

		if !list.Enabled || shouldBypassPublicRouteBlocklist(path) {
			c.Next()
			return
		}

		if list.Matches(path) {
			writePlainNotFound(c)
			return
		}

		c.Next()
	}
}

func (l *PublicRouteBlocklist) Matches(path string) bool {
	if l == nil || !l.Enabled {
		return false
	}
	for _, rule := range l.Rules {
		switch rule.Match {
		case "exact":
			if path == rule.Path {
				return true
			}
		case "prefix":
			if strings.HasPrefix(path, rule.Path) {
				return true
			}
		}
	}
	return false
}

func LogPublicRouteBlocklist(list *PublicRouteBlocklist) {
	if list == nil {
		return
	}
	sourceKind := "file"
	if list.UsingDefault {
		sourceKind = "default"
	}
	log.Printf("Public route blocklist loaded: source=%s path=%q enabled=%v effective_rules=%d",
		sourceKind, list.Source, list.Enabled, len(list.Rules))
}

func publicRouteBlocklistCandidatePaths() []string {
	var paths []string
	if explicit := strings.TrimSpace(os.Getenv("PUBLIC_ROUTE_BLOCKLIST_FILE")); explicit != "" {
		paths = append(paths, explicit)
	}
	if dataDir := strings.TrimSpace(os.Getenv("DATA_DIR")); dataDir != "" {
		paths = append(paths, filepath.Join(dataDir, publicRouteBlocklistFilename))
	}
	paths = append(paths,
		filepath.Join("data", publicRouteBlocklistFilename),
		filepath.Join("/app/data", publicRouteBlocklistFilename),
		publicRouteBlocklistFilename,
		filepath.Join("config", publicRouteBlocklistFilename),
		filepath.Join("/etc/sub2api", publicRouteBlocklistFilename),
	)
	return dedupeStringSlice(paths)
}

func normalizePublicRouteBlocklist(cfg publicRouteBlocklistFile, source string, usingDefault bool) (*PublicRouteBlocklist, error) {
	enabled := true
	if cfg.Enabled != nil {
		enabled = *cfg.Enabled
	}

	rules := make([]PublicRouteBlocklistRule, 0, len(cfg.Rules))
	seen := make(map[string]struct{}, len(cfg.Rules))
	for i, rule := range cfg.Rules {
		match := strings.ToLower(strings.TrimSpace(rule.Match))
		path := strings.TrimSpace(rule.Path)
		if match != "exact" && match != "prefix" {
			return nil, fmt.Errorf("invalid public route blocklist rule %d in %q: unsupported match %q", i+1, source, rule.Match)
		}
		if !strings.HasPrefix(path, "/") {
			return nil, fmt.Errorf("invalid public route blocklist rule %d in %q: path must start with /", i+1, source)
		}
		key := match + "\x00" + path
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		rules = append(rules, PublicRouteBlocklistRule{
			Match: match,
			Path:  path,
		})
	}

	return &PublicRouteBlocklist{
		Enabled:      enabled,
		Rules:        rules,
		Source:       source,
		UsingDefault: usingDefault,
	}, nil
}

func defaultPublicRouteBlocklistFile() publicRouteBlocklistFile {
	enabled := true
	return publicRouteBlocklistFile{
		Enabled: &enabled,
		Rules: []PublicRouteBlocklistRule{
			{Match: "prefix", Path: "/auth/"},
			{Match: "exact", Path: "/forgot-password"},
			{Match: "exact", Path: "/reset-password"},
			{Match: "exact", Path: "/key-usage"},
			{Match: "prefix", Path: "/legal/"},
			{Match: "prefix", Path: "/payment/"},
			{Match: "exact", Path: "/api/event_logging/batch"},
			{Match: "exact", Path: "/api/v1/settings/public"},
			{Match: "exact", Path: "/api/v1/auth/send-verify-code"},
			{Match: "exact", Path: "/api/v1/auth/validate-promo-code"},
			{Match: "exact", Path: "/api/v1/auth/validate-invitation-code"},
			{Match: "exact", Path: "/api/v1/auth/forgot-password"},
			{Match: "exact", Path: "/api/v1/auth/reset-password"},
			{Match: "prefix", Path: "/api/v1/auth/oauth/"},
			{Match: "prefix", Path: "/api/v1/payment/public/"},
			{Match: "prefix", Path: "/api/v1/payment/webhook/"},
			{Match: "exact", Path: "/api/v1/settings/email-unsubscribe"},
			{Match: "prefix", Path: "/api/v1/pages/"},
		},
	}
}

func shouldBypassPublicRouteBlocklist(path string) bool {
	switch path {
	case "/login",
		"/register",
		"/email-verify",
		"/status",
		"/favicon.ico",
		"/logo.png",
		"/responses",
		"/images",
		"/v1",
		"/v1beta",
		"/user/register",
		"/user/login",
		"/user/login/2fa",
		"/user/refresh",
		"/user/logout":
		return true
	}

	for _, prefix := range []string{
		"/static/app/",
		"/assets/",
		"/responses/",
		"/images/",
		"/v1/",
		"/v1beta/",
		"/backend-api/",
		"/antigravity/",
	} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func isCompletedSetupPath(path string) bool {
	return path == "/setup" || strings.HasPrefix(path, "/setup/")
}

func writePlainNotFound(c *gin.Context) {
	c.Data(http.StatusNotFound, "text/plain; charset=utf-8", []byte("404 Not Found"))
	c.Abort()
}

func dedupeStringSlice(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

//go:build unit

package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type settingPublicRepoStub struct {
	values map[string]string
}

func (s *settingPublicRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *settingPublicRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
}

func (s *settingPublicRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *settingPublicRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			out[key] = value
		}
	}
	return out, nil
}

func (s *settingPublicRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
}

func (s *settingPublicRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
}

func (s *settingPublicRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

func TestSettingService_GetPublicSettingsForInjection_IncludesClientEndpointFields(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{
		values: map[string]string{
			SettingKeySiteName:                             "企业数据中台",
			SettingKeySiteSubtitle:                         "统一数据目录、治理与服务编排入口",
			SettingKeyRegistrationEnabled:                  "false",
			SettingKeyPromoCodeEnabled:                     "true",
			SettingKeyGoogleOAuthEnabled:                   "false",
			SettingKeyBackendModeEnabled:                   "false",
			SettingKeyAPIBaseURL:                           "https://api-a.example.test;https://api-b.example.test",
			SettingKeyTableDefaultPageSize:                 "20",
			SettingKeyTablePageSizeOptions:                 "[10,20,50,100]",
			SettingKeyCustomMenuItems:                      "[]",
			SettingKeyCustomEndpoints:                      `[{"name":"HK","endpoint":"https://hk.example.test","description":"Hong Kong"}]`,
			SettingKeyChannelMonitorEnabled:                "false",
			SettingKeyChannelMonitorDefaultIntervalSeconds: "60",
			SettingKeyAllowUserViewErrorRequests:           "false",
		},
	}, &config.Config{})

	payload, err := svc.GetPublicSettingsForInjection(context.Background())
	require.NoError(t, err)

	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var out map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(raw, &out))

	require.JSONEq(t, `{"site_name":"企业数据中台","site_subtitle":"统一数据目录、治理与服务编排入口"}`, string(raw))
	require.NotContains(t, out, "registration_enabled")
	require.NotContains(t, out, "promo_code_enabled")
	require.NotContains(t, out, "google_oauth_enabled")
	require.NotContains(t, out, "backend_mode_enabled")
	require.NotContains(t, out, "api_base_url")
	require.NotContains(t, out, "table_default_page_size")
	require.NotContains(t, out, "table_page_size_options")
	require.NotContains(t, out, "custom_menu_items")
	require.NotContains(t, out, "custom_endpoints")
	require.NotContains(t, out, "channel_monitor_enabled")
	require.NotContains(t, out, "allow_user_view_error_requests")
}

func TestSettingService_GetPublicSettingsForInjection_IncludesEnabledLoginFeatures(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{
		values: map[string]string{
			SettingKeyRegistrationEnabled:              "true",
			SettingKeyEmailVerifyEnabled:               "true",
			SettingKeyRegistrationEmailSuffixWhitelist: `["@example.com"]`,
			SettingKeyPromoCodeEnabled:                 "true",
			SettingKeyInvitationCodeEnabled:            "true",
			SettingKeyTurnstileEnabled:                 "true",
			SettingKeyTurnstileSiteKey:                 "site-key",
			SettingKeyOIDCConnectEnabled:               "true",
			SettingKeyOIDCConnectProviderName:          "CorpID",
		},
	}, &config.Config{})

	payload, err := svc.GetPublicSettingsForInjection(context.Background())
	require.NoError(t, err)

	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	require.Equal(t, true, out["registration_enabled"])
	require.Equal(t, true, out["email_verify_enabled"])
	require.Equal(t, []any{"@example.com"}, out["registration_email_suffix_whitelist"])
	require.Equal(t, true, out["promo_code_enabled"])
	require.Equal(t, true, out["invitation_code_enabled"])
	require.Equal(t, true, out["turnstile_enabled"])
	require.Equal(t, "site-key", out["turnstile_site_key"])
	require.NotContains(t, out, "oidc_oauth_enabled")
	require.NotContains(t, out, "oidc_oauth_provider_name")
}

func TestSettingService_GetPublicSettings_ExposesRegistrationEmailSuffixWhitelist(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyRegistrationEnabled:              "true",
			SettingKeyEmailVerifyEnabled:               "true",
			SettingKeyRegistrationEmailSuffixWhitelist: `["@EXAMPLE.com"," @foo.bar ","*.EDU.CN","@invalid_domain",""]`,
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"@example.com", "@foo.bar", "*.edu.cn"}, settings.RegistrationEmailSuffixWhitelist)
}

func TestSettingService_GetPublicSettings_ExposesTablePreferences(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyTableDefaultPageSize: "50",
			SettingKeyTablePageSizeOptions: "[20,50,100]",
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, 50, settings.TableDefaultPageSize)
	require.Equal(t, []int{20, 50, 100}, settings.TablePageSizeOptions)
}

func TestSettingService_GetPublicSettings_ExposesForceEmailOnThirdPartySignup(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyForceEmailOnThirdPartySignup: "true",
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.True(t, settings.ForceEmailOnThirdPartySignup)
}

func TestSettingService_GetPublicSettings_ExposesAllowUserViewErrorRequests(t *testing.T) {
	repo := &settingPublicRepoStub{
		values: map[string]string{
			SettingKeyAllowUserViewErrorRequests: "true",
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.True(t, settings.AllowUserViewErrorRequests)
}

func TestSettingService_GetPublicSettings_ExposesWeChatOAuthModeCapabilities(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{
		values: map[string]string{
			SettingKeyWeChatConnectEnabled:             "true",
			SettingKeyWeChatConnectAppID:               "wx-mp-app",
			SettingKeyWeChatConnectAppSecret:           "wx-mp-secret",
			SettingKeyWeChatConnectMode:                "mp",
			SettingKeyWeChatConnectScopes:              "snsapi_base",
			SettingKeyWeChatConnectOpenEnabled:         "true",
			SettingKeyWeChatConnectMPEnabled:           "true",
			SettingKeyWeChatConnectRedirectURL:         "https://api.example.com/api/v1/auth/oauth/wechat/callback",
			SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
		},
	}, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.True(t, settings.WeChatOAuthEnabled)
	require.True(t, settings.WeChatOAuthOpenEnabled)
	require.True(t, settings.WeChatOAuthMPEnabled)
}

func TestSettingService_GetPublicSettings_DoesNotExposeMobileOnlyWeChatAsWebOAuthAvailable(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{
		values: map[string]string{
			SettingKeyWeChatConnectEnabled:             "true",
			SettingKeyWeChatConnectMobileEnabled:       "true",
			SettingKeyWeChatConnectMode:                "mobile",
			SettingKeyWeChatConnectMobileAppID:         "wx-mobile-app",
			SettingKeyWeChatConnectMobileAppSecret:     "wx-mobile-secret",
			SettingKeyWeChatConnectFrontendRedirectURL: "/auth/wechat/callback",
		},
	}, &config.Config{})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.False(t, settings.WeChatOAuthEnabled)
	require.False(t, settings.WeChatOAuthOpenEnabled)
	require.False(t, settings.WeChatOAuthMPEnabled)
	require.True(t, settings.WeChatOAuthMobileEnabled)
}

func TestSettingService_GetPublicSettings_FallsBackToConfigForWeChatOAuthCapabilities(t *testing.T) {
	svc := NewSettingService(&settingPublicRepoStub{values: map[string]string{}}, &config.Config{
		WeChat: config.WeChatConnectConfig{
			Enabled:             true,
			OpenEnabled:         true,
			OpenAppID:           "wx-open-config",
			OpenAppSecret:       "wx-open-secret",
			FrontendRedirectURL: "/auth/wechat/config-callback",
		},
	})

	settings, err := svc.GetPublicSettings(context.Background())
	require.NoError(t, err)
	require.True(t, settings.WeChatOAuthEnabled)
	require.True(t, settings.WeChatOAuthOpenEnabled)
	require.False(t, settings.WeChatOAuthMPEnabled)
	require.False(t, settings.WeChatOAuthMobileEnabled)
}

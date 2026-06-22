//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type updateServiceCacheStub struct {
	data string
}

func (s *updateServiceCacheStub) GetUpdateInfo(context.Context) (string, error) {
	if s.data == "" {
		return "", errors.New("cache miss")
	}
	return s.data, nil
}

func (s *updateServiceCacheStub) SetUpdateInfo(_ context.Context, data string, _ time.Duration) error {
	s.data = data
	return nil
}

type updateServiceGitHubClientStub struct {
	release *GitHubRelease
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(context.Context, string) (*GitHubRelease, error) {
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) DownloadFile(context.Context, string, string, int64) error {
	panic("DownloadFile should not be called when no update is available")
}

func (s *updateServiceGitHubClientStub) FetchChecksumFile(context.Context, string) ([]byte, error) {
	panic("FetchChecksumFile should not be called when no update is available")
}

func TestUpdateServicePerformUpdateNoUpdateReturnsSentinel(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.132",
				Name:    "v0.1.132",
			},
		},
		"0.1.132",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoUpdateAvailable))
	require.ErrorIs(t, err, ErrNoUpdateAvailable)
}

func TestUpdateServiceForkVersionUsesUpstreamBaseForUpdateCheck(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
		},
		"0.1.138.kim",
		"release",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.Equal(t, "0.1.138.kim", info.CurrentVersion)
	require.Equal(t, "0.1.138", info.UpstreamCurrentVersion)
	require.Equal(t, "0.1.138", info.LatestVersion)
	require.True(t, info.ForkBuild)
	require.False(t, info.HasUpdate)
}

func TestUpdateServiceForkVersionReportsNewUpstreamVersion(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.139",
				Name:    "v0.1.139",
			},
		},
		"0.1.138.kim",
		"release",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.Equal(t, "0.1.138", info.UpstreamCurrentVersion)
	require.Equal(t, "0.1.139", info.LatestVersion)
	require.True(t, info.ForkBuild)
	require.True(t, info.HasUpdate)
}

func TestUpdateServicePerformUpdateForkBuildRequiresManualUpdate(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.139",
				Name:    "v0.1.139",
			},
		},
		"0.1.138.kim",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.ErrorIs(t, err, ErrManualUpdateRequired)
}

func TestUpstreamVersionBaseSupportsForkSuffixes(t *testing.T) {
	require.Equal(t, "0.1.138", upstreamVersionBase("0.1.138.kim"))
	require.Equal(t, "0.1.138", upstreamVersionBase("0.1.138-kim"))
	require.Equal(t, "0.1.138", upstreamVersionBase("v0.1.138+kim"))
	require.True(t, isForkVersion("0.1.138.kim"))
	require.False(t, isForkVersion("v0.1.138"))
}

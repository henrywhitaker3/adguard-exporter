package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromEnvAllowsBlankCredentials(t *testing.T) {
	t.Setenv("ADGUARD_SERVERS", "http://adguard.test:3000")
	t.Setenv("ADGUARD_USERNAMES", "")
	t.Setenv("ADGUARD_PASSWORDS", "")

	configuration, err := FromEnv()

	require.NoError(t, err)
	require.Len(t, configuration.Configs, 1)
	assert.Equal(t, Config{
		Url:      "http://adguard.test:3000",
		Username: "",
		Password: "",
	}, configuration.Configs[0])
}

func TestFromEnvKeepsBlankCredentialsAlignedWithTheirServer(t *testing.T) {
	t.Setenv("ADGUARD_SERVERS", "http://first.test:3000,http://second.test:3000")
	t.Setenv("ADGUARD_USERNAMES", ",second-user")
	t.Setenv("ADGUARD_PASSWORDS", ",second-password")

	configuration, err := FromEnv()

	require.NoError(t, err)
	assert.Equal(t, []Config{
		{Url: "http://first.test:3000"},
		{Url: "http://second.test:3000", Username: "second-user", Password: "second-password"},
	}, configuration.Configs)
}

package flags

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvConfig_Defaults(t *testing.T) {
	cmd := new(cobra.Command)
	SetDefaults()
	RegisterDockerFlags(cmd)

	err := EnvConfig(cmd)
	require.NoError(t, err)

	assert.Equal(t, "unix:///var/run/docker.sock", os.Getenv("DOCKER_HOST"))
	assert.Equal(t, "", os.Getenv("DOCKER_TLS_VERIFY"))
	// Re-enable this test when we've moved to github actions.
	// assert.Equal(t, DockerAPIMinVersion, os.Getenv("DOCKER_API_VERSION"))
}

func TestEnvConfig_Custom(t *testing.T) {
	cmd := new(cobra.Command)
	SetDefaults()
	RegisterDockerFlags(cmd)

	err := cmd.ParseFlags([]string{"--host", "some-custom-docker-host", "--tlsverify", "--api-version", "1.99"})
	require.NoError(t, err)

	err = EnvConfig(cmd)
	require.NoError(t, err)

	assert.Equal(t, "some-custom-docker-host", os.Getenv("DOCKER_HOST"))
	assert.Equal(t, "1", os.Getenv("DOCKER_TLS_VERIFY"))
	// Re-enable this test when we've moved to github actions.
	// assert.Equal(t, "1.99", os.Getenv("DOCKER_API_VERSION"))
}

func TestGetSecretsFromFilesWithString(t *testing.T) {
	value := "supersecretstring"

	err := os.Setenv("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD", value)
	require.NoError(t, err)

	testGetSecretsFromFiles(t, "notification-email-server-password", value)
}

func TestGetSecretsFromFilesWithFile(t *testing.T) {
	value := "megasecretstring"

	// Create the temporary file which will contain a secret.
	file, err := ioutil.TempFile(os.TempDir(), "watchtower-")
	require.NoError(t, err)
	defer os.Remove(file.Name()) // Make sure to remove the temporary file later.

	// Write the secret to the temporary file.
	secret := []byte(value)
	_, err = file.Write(secret)
	require.NoError(t, err)

	err = os.Setenv("WATCHTOWER_NOTIFICATION_EMAIL_SERVER_PASSWORD", file.Name())
	require.NoError(t, err)

	testGetSecretsFromFiles(t, "notification-email-server-password", value)
}

func TestGetSecretsArrayFromFilesWithFile(t *testing.T) {
	expected := []string{"first line", "second line"}

	// Create the temporary file which will contain a secret.
	file, err := ioutil.TempFile(os.TempDir(), "watchtower-")
	require.NoError(t, err)
	defer os.Remove(file.Name()) // Make sure to remove the temporary file later.

	// Write the secret to the temporary file.
	secret := []byte(strings.Join(expected, "\n"))
	_, err = file.Write(secret)
	require.NoError(t, err)

	err = os.Setenv("WATCHTOWER_NOTIFICATION_URL", file.Name())
	require.NoError(t, err)

	cmd := new(cobra.Command)
	SetDefaults()
	RegisterNotificationFlags(cmd)
	GetSecretsFromFiles(cmd)
	actual, err := cmd.PersistentFlags().GetStringArray("notification-url")
	require.NoError(t, err)

	require.Equal(t, 2, len(actual))
	require.Equal(t, expected[0], actual[0])
	require.Equal(t, expected[1], actual[1])
}

func testGetSecretsFromFiles(t *testing.T, flagName string, expected string) {
	cmd := new(cobra.Command)
	SetDefaults()
	RegisterNotificationFlags(cmd)
	GetSecretsFromFiles(cmd)
	value, err := cmd.PersistentFlags().GetString(flagName)
	require.NoError(t, err)

	assert.Equal(t, expected, value)
}

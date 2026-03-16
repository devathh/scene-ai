package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()

	require.NotNil(t, loader)
	require.NotNil(t, loader.validator)
}

func TestLoader_Load_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yml")

	t.Setenv("APP_ENV", "local")
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DBNAME", "testdb")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "pass")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_USERNAME", "redis_user")
	t.Setenv("REDIS_PASSWORD", "redis_pass")

	content := `app:
  name: test-app
  version: 1.0.0
  env: ${APP_ENV}
server:
  addr: localhost:8080
  tls:
    enable: false
    server-cert-path: ./cert.crt
    server-key-path: ./key.key
  read-timeout: 5s
  write-timeout: 5s
  idle-timeout: 1m
persistence:
  postgres:
    host: ${POSTGRES_HOST}
    port: ${POSTGRES_PORT}
    sslmode: disable
    dbname: ${POSTGRES_DBNAME}
    auth:
      user: ${POSTGRES_USER}
      password: ${POSTGRES_PASSWORD}
    conn:
      max-idles: 5
      max-opens: 10
      max-idle-time: 5m
      max-lifetime: 10m
cache:
  refresh-ttl: 10s
  user-ttl: 30m
  redis:
    host: ${REDIS_HOST}
    port: ${REDIS_PORT}
    db: 1
    auth:
      username: ${REDIS_USERNAME}
      password: ${REDIS_PASSWORD}
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	loader := NewLoader()

	cfg, err := loader.Load(configPath)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "test-app", cfg.App.Name)
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "local", cfg.App.Env)
	assert.Equal(t, "localhost:8080", cfg.Server.Addr)
	assert.Equal(t, 5*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, "localhost", cfg.Persistence.Postgres.Host)
	assert.Equal(t, 5432, cfg.Persistence.Postgres.Port)

	assert.Equal(t, 10*time.Second, cfg.Cache.RefreshTTL)
	assert.Equal(t, 30*time.Minute, cfg.Cache.UserTTL)
	assert.Equal(t, 6379, cfg.Cache.Redis.Port)
	assert.Equal(t, 1, cfg.Cache.Redis.DB)
}

func TestLoader_Load_FileNotFound(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load("/non/existent/path/config.yml")

	require.Error(t, err)
	assert.NotContains(t, err.Error(), "empty path to config file")
}

func TestLoader_Load_InvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yml")
	content := `app:
  name: test
  version: [invalid yaml structure
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	loader := NewLoader()

	_, err = loader.Load(configPath)

	require.Error(t, err)
}

func TestLoader_Load_ValidationFailed(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_struct.yml")
	content := `app:
  name: test
  version: bad-version
  env: local
server:
  addr: localhost:8080
  tls:
    enable: false
    server-cert-path: ./cert.crt
    server-key-path: ./key.key
  read-timeout: 5s
  write-timeout: 5s
  idle-timeout: 1m
persistence:
  postgres:
    host: localhost
    port: 5432
    sslmode: disable
    dbname: testdb
    auth:
      user: user
      password: pass
    conn:
      max-idles: 5
      max-opens: 10
      max-idle-time: 5m
      max-lifetime: 10m
cache:
  refresh-ttl: 10s
  user-ttl: 30m
  redis:
    host: localhost
    port: 6379
    db: 0
    auth:
      username: user
      password: pass
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	loader := NewLoader()

	_, err = loader.Load(configPath)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestLoader_Load_MissingCacheFields(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "missing_cache_fields.yml")
	content := `app:
  name: test
  version: 1.0.0
  env: local
server:
  addr: localhost:8080
  tls:
    enable: false
    server-cert-path: ./cert.crt
    server-key-path: ./key.key
  read-timeout: 5s
  write-timeout: 5s
  idle-timeout: 1m
persistence:
  postgres:
    host: localhost
    port: 5432
    sslmode: disable
    dbname: testdb
    auth:
      user: user
      password: pass
    conn:
      max-idles: 5
      max-opens: 10
      max-idle-time: 5m
      max-lifetime: 10m
cache:
  redis:
    host: localhost
    port: 6379
    db: 0
    auth:
      username: user
      password: pass
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	loader := NewLoader()

	_, err = loader.Load(configPath)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestLoader_Load_EmptyPath(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load("")

	require.Error(t, err)
	assert.Equal(t, "empty path to config file", err.Error())
}

func TestConstructor_Init(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "constructor_test.yml")

	t.Setenv("APP_ENV", "prod")
	t.Setenv("POSTGRES_HOST", "db.prod")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DBNAME", "proddb")
	t.Setenv("POSTGRES_USER", "admin")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("REDIS_HOST", "redis.prod")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_USERNAME", "admin")
	t.Setenv("REDIS_PASSWORD", "secret")

	content := `app:
  name: prod-app
  version: 2.0.0
  env: ${APP_ENV}
server:
  addr: 0.0.0.0:80
  tls:
    enable: false
    server-cert-path: ./cert.crt
    server-key-path: ./key.key
  read-timeout: 10s
  write-timeout: 10s
  idle-timeout: 2m
persistence:
  postgres:
    host: ${POSTGRES_HOST}
    port: ${POSTGRES_PORT}
    sslmode: disable
    dbname: ${POSTGRES_DBNAME}
    auth:
      user: ${POSTGRES_USER}
      password: ${POSTGRES_PASSWORD}
    conn:
      max-idles: 20
      max-opens: 50
      max-idle-time: 10m
      max-lifetime: 30m
cache:
  refresh-ttl: 15s
  user-ttl: 1h
  redis:
    host: ${REDIS_HOST}
    port: ${REDIS_PORT}
    db: 2
    auth:
      username: ${REDIS_USERNAME}
      password: ${REDIS_PASSWORD}
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	constructor := NewConstructor()
	cfg, err := constructor.Init(configPath)
	require.NoError(t, err)
	assert.Equal(t, "prod-app", cfg.App.Name)
	assert.Equal(t, 15*time.Second, cfg.Cache.RefreshTTL)
	assert.Equal(t, 1*time.Hour, cfg.Cache.UserTTL)

	t.Setenv("APP_CONFIG_PATH", configPath)
	constructor2 := NewConstructor()
	cfg2, err := constructor2.Init("")
	require.NoError(t, err)
	assert.Equal(t, "prod", cfg2.App.Env)

	t.Setenv("APP_CONFIG_PATH", "")

	constructor3 := NewConstructor()
	_, err = constructor3.Init("")
	require.Error(t, err)
}

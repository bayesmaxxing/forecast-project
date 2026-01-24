package config

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type Config struct {
	JWTSecret      []byte
	AllowedOrigin  string
	DBConnString   string
}

// Load loads configuration from environment variables and Google Secret Manager.
// For local development, set USE_LOCAL_SECRETS=true and provide JWT_SECRET as env var.
func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{
		AllowedOrigin: getEnvOrDefault("ALLOWED_ORIGIN", "https://www.samuelsforecasts.com"),
		DBConnString:  os.Getenv("DB_CONNECTION_STRING"),
	}

	// For local development, allow using environment variables directly
	if os.Getenv("USE_LOCAL_SECRETS") == "true" {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			return nil, fmt.Errorf("JWT_SECRET environment variable is required when USE_LOCAL_SECRETS=true")
		}
		if len(jwtSecret) < 32 {
			return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
		}
		cfg.JWTSecret = []byte(jwtSecret)
		return cfg, nil
	}

	// In production, fetch from Google Secret Manager
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is required")
	}

	jwtSecret, err := getSecret(ctx, projectID, "jwt-secret")
	if err != nil {
		return nil, fmt.Errorf("failed to get JWT secret: %w", err)
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT secret must be at least 32 characters")
	}
	cfg.JWTSecret = jwtSecret

	return cfg, nil
}

func getSecret(ctx context.Context, projectID, secretID string) ([]byte, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretID)
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to access secret %s: %w", secretID, err)
	}

	return result.Payload.Data, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

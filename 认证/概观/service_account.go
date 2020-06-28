import (
        "context"
        "fmt"

        "cloud.google.com/go/pubsub"
)

// serviceAccount shows how to use a service account to authenticate.
func serviceAccount() error {
        // Download service account key per https://cloud.google.com/docs/authentication/production.
        // Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
        // This environment variable will be automatically picked up by the client.
        client, err := pubsub.NewClient(context.Background(), "your-project-id")
        if err != nil {
                return fmt.Errorf("pubsub.NewClient: %v", err)
        }
        // Use the authenticated client.
        _ = client

        return nil
}

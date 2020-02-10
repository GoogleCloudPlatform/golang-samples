package notifications

// [START scc_receive_notifications]
import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/jsonpb"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1p1beta1"
)

func parseNotificationMessage(ctx context.Context, msg *pubsub.Message) {
	var notificationMessage = new(securitycenterpb.NotificationMessage)
	u := &jsonpb.Unmarshaler{}
	u.Unmarshal(bytes.NewReader(msg.Data), notificationMessage)

	fmt.Printf("Got finding: %v", notificationMessage.GetFinding())
	msg.Ack()
}

func receiveMessages(projectId string, subscriptionName string) error {
	ctx := context.Background()

	// TODO projectId := "your-project-id"
	// TODO subsriptionName := "your-subscription-name"

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subscriptionName)
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = sub.Receive(cctx, parseNotificationMessage)
	if err != nil {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}

// [END scc_receive_notifications]

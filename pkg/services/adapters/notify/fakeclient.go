package notify

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	port "github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockState struct {
	Notifications []port.Notification
}

type mockClient struct {
	port.NotificationsClient // TODO: implement

	mu    sync.Mutex
	state MockState
}

func NewMock(state MockState) port.NotificationsClient {
	return &mockClient{state: state}
}

func (m *mockClient) SendNotification(ctx context.Context, in *port.NotificationRequest, opts ...grpc.CallOption) (*port.NotificationsResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state.Notifications = append(m.state.Notifications, port.Notification{
		Id:         uuid.NewString(),
		Message:    in.GetMessage(),
		Recipients: in.GetRecipients(),
		CreatedAt:  timestamppb.New(time.Now()),
	})
	return &port.NotificationsResponse{}, nil
}

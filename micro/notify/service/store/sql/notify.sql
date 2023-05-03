CREATE TABLE notifications (
  id UUID PRIMARY KEY,
  message TEXT NOT NULL,
  recipients TEXT[] NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX notifications_created_at_idx ON notifications (created_at DESC);

-- name: CreateNotification :exec
INSERT INTO notifications (message, recipients) VALUES ($1, $2) RETURNING id, created_at;

-- GetNotificationsByRecipientsIds :many
-- sorted by most recent first
SELECT id, message, created_at FROM notifications WHERE $1 = ANY(recipients) ORDER BY created_at DESC LIMIT $2 OFFSET $3;
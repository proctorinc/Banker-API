package resolvers

import (
	"context"
	"time"

	"github.com/proctorinc/banker/internal/db"
)

func (r *accountSyncItemResolver) ID(ctx context.Context, syncItem *db.AccountSyncItem) (string, error) {
	return syncItem.ID.String(), nil
}

func (r *accountSyncItemResolver) Date(ctx context.Context, syncItem *db.AccountSyncItem) (string, error) {
	return syncItem.Date.Format(time.RFC3339), nil
}

func (r *accountSyncItemResolver) UploadSource(ctx context.Context, syncItem *db.AccountSyncItem) (string, error) {
	return string(syncItem.Uploadsource), nil
}

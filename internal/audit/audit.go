package audit

import (
	"context"
	"encoding/json"

	"trustpin/internal/store"
)

type Logger struct {
	Store *store.Postgres
}

func (l Logger) Event(ctx context.Context, e store.AuditEvent, metadata any) error {
	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		e.Metadata = b
	}
	return l.Store.InsertAuditEvent(ctx, e)
}

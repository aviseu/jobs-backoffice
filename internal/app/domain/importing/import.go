package importing

import (
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
)

type importEntry struct {
	startedAt time.Time
	endedAt   null.Time
	error     null.String
	status    aggregator.ImportStatus
	id        uuid.UUID
	channelID uuid.UUID
}

func newImportEntry(id, channelID uuid.UUID, status aggregator.ImportStatus, startedAt time.Time, endedAt null.Time, err null.String) *importEntry {
	i := &importEntry{
		id:        id,
		channelID: channelID,
		status:    status,
		startedAt: startedAt,
		endedAt:   endedAt,
		error:     err,
	}

	return i
}

func (i *importEntry) markAsFailed(err error) {
	i.status = aggregator.ImportStatusFailed
	i.endedAt = null.TimeFrom(time.Now())
	i.error = null.StringFrom(err.Error())
}

func (i *importEntry) markAsFetching() {
	i.status = aggregator.ImportStatusFetching
}

func (i *importEntry) markAsProcessing() {
	i.status = aggregator.ImportStatusProcessing
}

func (i *importEntry) markAsPublishing() {
	i.status = aggregator.ImportStatusPublishing
}

func (i *importEntry) markAsCompleted() {
	i.status = aggregator.ImportStatusCompleted
	i.endedAt = null.TimeFrom(time.Now())
}

func (i *importEntry) toAggregate() *aggregator.Import {
	return &aggregator.Import{
		ID:        i.id,
		ChannelID: i.channelID,
		StartedAt: i.startedAt,
		EndedAt:   i.endedAt,
		Error:     i.error,
		Status:    i.status,
	}
}

func newImportFromAggregator(i *aggregator.Import) *importEntry {
	return newImportEntry(
		i.ID,
		i.ChannelID,
		i.Status,
		i.StartedAt,
		i.EndedAt,
		i.Error,
	)
}

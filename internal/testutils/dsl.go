package testutils

import (
	"bytes"
	"github.com/aviseu/jobs-backoffice/internal/app/application/http"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/scheduling"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v3"
	"log/slog"
	oghttp "net/http"
	"net/http/httptest"
	"time"
)

type DSL struct {
	// Dependencies
	HTTPClient       arbeitnow.HTTPClient
	RequestLogger    *RequestLogger
	Logger           *slog.Logger
	LogBuffer        *bytes.Buffer
	AirbeitnowServer *httptest.Server
	Config           *importing.Config

	// Infrastructure
	AirbeitnowService *arbeitnow.Service
	JobRepository     *JobRepository
	ChannelRepository *ChannelRepository
	ImportRepository  *ImportRepository
	PubSubService     *PubSubService

	// Domains
	ConfiguringService *configuring.Service
	Factory            *importing.Factory
	ImportService      *importing.Service
	SchedulingService  *scheduling.Service

	// Application
	APIServer  oghttp.Handler
	HTTPConfig *http.Config
}

type DSLOptions func(*DSL)

func WithChannelRepositoryError(err error) DSLOptions {
	return func(dsl *DSL) {
		if dsl.ChannelRepository == nil {
			dsl.ChannelRepository = NewChannelRepository()
		}
		dsl.ChannelRepository.FailWith(err)
	}
}

func WithImportRepositoryError(err error) DSLOptions {
	return func(dsl *DSL) {
		if dsl.ImportRepository == nil {
			dsl.ImportRepository = NewImportRepository()
		}
		dsl.ImportRepository.FailWith(err)
	}
}

func WithPubSubServiceError(err error) DSLOptions {
	return func(dsl *DSL) {
		if dsl.PubSubService == nil {
			dsl.PubSubService = NewPubSubService()
		}
		dsl.PubSubService.FailWith(err)
	}
}

func WithArbeitnowEnabled() DSLOptions {
	return func(dsl *DSL) {
		dsl.AirbeitnowServer = NewArbeitnowServer()
		dsl.RequestLogger = NewRequestLogger(oghttp.DefaultClient)

		var ch *aggregator.Channel
		if dsl.ChannelRepository != nil && len(dsl.Channels()) > 0 {
			ch = dsl.Channels()[0]
		}
		dsl.AirbeitnowService = arbeitnow.NewService(
			dsl.RequestLogger,
			arbeitnow.Config{URL: dsl.AirbeitnowServer.URL},
			ch,
		)
	}
}

func WithHTTPConfig(cfg http.Config) DSLOptions {
	return func(dsl *DSL) {
		dsl.HTTPConfig = &cfg
	}
}

type WithChannelOptions func(ch *aggregator.Channel)

func WithChannelID(id uuid.UUID) WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.ID = id
	}
}

func WithChannelIntegration(i aggregator.Integration) WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.Integration = i
	}
}

func WithChannelName(name string) WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.Name = name
	}
}

func WithChannelDeactivated() WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.Status = aggregator.ChannelStatusInactive
	}
}

func WithChannelActivated() WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.Status = aggregator.ChannelStatusActive
	}
}

func WithChannelTimestamps(cat, uat time.Time) WithChannelOptions {
	return func(ch *aggregator.Channel) {
		ch.CreatedAt = cat
		ch.UpdatedAt = uat
	}
}

func WithChannel(opts ...WithChannelOptions) DSLOptions {
	return func(dsl *DSL) {
		if dsl.ChannelRepository == nil {
			dsl.ChannelRepository = NewChannelRepository()
		}
		ch := &aggregator.Channel{
			ID:          uuid.New(),
			Name:        "channel 1",
			Integration: aggregator.IntegrationArbeitnow,
			Status:      aggregator.ChannelStatusActive,
			CreatedAt:   time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		}
		for _, opt := range opts {
			opt(ch)
		}
		dsl.ChannelRepository.Add(
			ch,
		)
	}
}

type WithImportOptions func(i *aggregator.Import)

func WithImportChannelID(chID uuid.UUID) WithImportOptions {
	return func(i *aggregator.Import) {
		i.ChannelID = chID
	}
}

func WithImportID(id uuid.UUID) WithImportOptions {
	return func(i *aggregator.Import) {
		i.ID = id
	}
}

func WithImportError(err string) WithImportOptions {
	return func(i *aggregator.Import) {
		i.Error = null.StringFrom(err)
	}
}

func WithoutImportError() WithImportOptions {
	return func(i *aggregator.Import) {
		i.Error = null.NewString("", false)
	}
}

func WithImportStatus(status aggregator.ImportStatus) WithImportOptions {
	return func(i *aggregator.Import) {
		i.Status = status
	}
}

func WithImportStartedAt(startedAt time.Time) WithImportOptions {
	return func(i *aggregator.Import) {
		i.StartedAt = startedAt
	}
}

func WithImportEndedAt(endedAt time.Time) WithImportOptions {
	return func(i *aggregator.Import) {
		i.EndedAt = null.NewTime(endedAt, true)
	}
}

func WithoutImportEndedAt() WithImportOptions {
	return func(i *aggregator.Import) {
		i.EndedAt = null.NewTime(time.Now(), false)
	}
}

func WithImportJobs(result aggregator.ImportJobResult, count int) WithImportOptions {
	return func(i *aggregator.Import) {
		for j := 0; j < count; j++ {
			i.Jobs = append(i.Jobs, &aggregator.ImportJob{
				ID:     uuid.New(),
				Result: result,
			})
		}
	}
}

func WithImport(opts ...WithImportOptions) DSLOptions {
	return func(dsl *DSL) {
		if dsl.ImportRepository == nil {
			dsl.ImportRepository = NewImportRepository()
		}
		i := &aggregator.Import{
			StartedAt: time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
			EndedAt:   null.NewTime(time.Now(), false),
			Error:     null.NewString("", false),
			Jobs:      make([]*aggregator.ImportJob, 0),
			Status:    aggregator.ImportStatusPending,
			ID:        uuid.New(),
			ChannelID: uuid.New(),
		}
		for _, opt := range opts {
			opt(i)
		}
		dsl.ImportRepository.AddImport(
			i,
		)
	}
}

func NewDSL(opts ...DSLOptions) *DSL {
	dsl := &DSL{}

	for _, opt := range opts {
		opt(dsl)
	}

	if dsl.ChannelRepository == nil {
		dsl.ChannelRepository = NewChannelRepository()
	}
	if dsl.ConfiguringService == nil {
		dsl.ConfiguringService = configuring.NewService(dsl.ChannelRepository)
	}
	if dsl.ImportRepository == nil {
		dsl.ImportRepository = NewImportRepository()
	}
	if dsl.JobRepository == nil {
		dsl.JobRepository = NewJobRepository()
	}
	if dsl.HTTPClient == nil {
		dsl.HTTPClient = NewHTTPClientMock()
	}
	if dsl.Config == nil {
		dsl.Config = &importing.Config{
			Import: struct {
				ResultBufferSize int `split_words:"true" default:"10"`
				ResultWorkers    int `split_words:"true" default:"10"`
			}{
				ResultBufferSize: 10,
				ResultWorkers:    10,
			},
		}
	}
	if dsl.Factory == nil {
		dsl.Factory = importing.NewFactory(dsl.HTTPClient, *dsl.Config)
	}
	if dsl.Logger == nil {
		dsl.LogBuffer, dsl.Logger = NewLogger()
	}
	if dsl.ImportService == nil {
		dsl.ImportService = importing.NewService(dsl.ChannelRepository, dsl.ImportRepository, dsl.JobRepository, dsl.Factory, *dsl.Config, dsl.Config.Import.ResultBufferSize, dsl.Config.Import.ResultWorkers, dsl.Logger)
	}
	if dsl.PubSubService == nil {
		dsl.PubSubService = NewPubSubService()
	}
	if dsl.SchedulingService == nil {
		dsl.SchedulingService = scheduling.NewService(dsl.ImportRepository, dsl.ChannelRepository, dsl.PubSubService, dsl.Logger)
	}

	if dsl.HTTPConfig == nil {
		dsl.HTTPConfig = &http.Config{}
	}

	if dsl.APIServer == nil {
		dsl.APIServer = http.APIRootHandler(dsl.ConfiguringService, dsl.ChannelRepository, dsl.ImportRepository, dsl.SchedulingService, *dsl.HTTPConfig, dsl.Logger)
	}

	return dsl
}

func (dsl *DSL) Channels() []*aggregator.Channel {
	var channels []*aggregator.Channel
	for _, ch := range dsl.ChannelRepository.Channels {
		channels = append(channels, ch)
	}

	return channels
}

func (dsl *DSL) Channel(id uuid.UUID) *aggregator.Channel {
	return dsl.ChannelRepository.Channels[id]
}

func (dsl *DSL) FirstChannel() *aggregator.Channel {
	for _, ch := range dsl.ChannelRepository.Channels {
		return ch
	}

	return nil
}

func (dsl *DSL) Imports() []*aggregator.Import {
	var imports []*aggregator.Import
	for _, i := range dsl.ImportRepository.Imports {
		imports = append(imports, i)
	}

	return imports
}

func (dsl *DSL) FirstImport() *aggregator.Import {
	for _, i := range dsl.ImportRepository.Imports {
		return i
	}

	return nil
}

func (dsl *DSL) PublishedImports() []uuid.UUID {
	return dsl.PubSubService.ImportIDs
}

func (dsl *DSL) LogLines() []string {
	return LogLines(dsl.LogBuffer)
}

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
	JobRepository       *JobRepository
	ChannelRepository   *ChannelRepository
	ImportRepository    *ImportRepository
	PubSubImportService *PubSubImportService
	PubSubJobService    *PubSubJobService

	// Domains
	ConfiguringService *configuring.Service
	ImportService      *importing.Service
	SchedulingService  *scheduling.Service

	// Application
	APIServer    oghttp.Handler
	ImportServer oghttp.Handler
	HTTPConfig   *http.Config
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
		if dsl.PubSubImportService == nil {
			dsl.PubSubImportService = NewPubSubImportService()
		}
		dsl.PubSubImportService.FailWith(err)
	}
}

func WithArbeitnowEnabled() DSLOptions {
	return func(dsl *DSL) {
		dsl.AirbeitnowServer = NewArbeitnowServer()
		dsl.RequestLogger = NewRequestLogger(oghttp.DefaultClient)
		dsl.HTTPClient = dsl.RequestLogger

		if dsl.Config == nil {
			dsl.Config = dsl.defaultConfig()
			dsl.Config.Arbeitnow.URL = dsl.AirbeitnowServer.URL
		}
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

func WithImportMetrics(metricType aggregator.ImportMetricType, count int) WithImportOptions {
	return func(i *aggregator.Import) {
		for j := 0; j < count; j++ {
			i.Metrics = append(i.Metrics, &aggregator.ImportMetric{
				ID:         uuid.New(),
				JobID:      uuid.New(),
				MetricType: metricType,
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
			Metrics:   make([]*aggregator.ImportMetric, 0),
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

type WithJobOptions func(i *aggregator.Job)

func WithJobID(id uuid.UUID) WithJobOptions {
	return func(j *aggregator.Job) {
		j.ID = id
	}
}

func WithJobChannelID(chID uuid.UUID) WithJobOptions {
	return func(j *aggregator.Job) {
		j.ChannelID = chID
	}
}

func WithJobStatus(status aggregator.JobStatus) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Status = status
	}
}

func WithJobPublishStatus(status aggregator.JobPublishStatus) WithJobOptions {
	return func(j *aggregator.Job) {
		j.PublishStatus = status
	}
}

func WithJobURL(url string) WithJobOptions {
	return func(j *aggregator.Job) {
		j.URL = url
	}
}

func WithJobTitle(title string) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Title = title
	}
}

func WithJobDescription(description string) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Description = description
	}
}

func WithJobSource(source string) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Source = source
	}
}

func WithJobLocation(location string) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Location = location
	}
}

func WithJobRemote(isRemote bool) WithJobOptions {
	return func(j *aggregator.Job) {
		j.Remote = isRemote
	}
}

func WithJobPostedAt(postedAt time.Time) WithJobOptions {
	return func(j *aggregator.Job) {
		j.PostedAt = postedAt
	}
}

func WithJobTimestamps(cat, uat time.Time) WithJobOptions {
	return func(j *aggregator.Job) {
		j.CreatedAt = cat
		j.UpdatedAt = uat
	}
}

func WithJob(opts ...WithJobOptions) DSLOptions {
	return func(dsl *DSL) {
		if dsl.JobRepository == nil {
			dsl.JobRepository = NewJobRepository()
		}
		j := &aggregator.Job{
			ID:            uuid.New(),
			ChannelID:     uuid.New(),
			Status:        aggregator.JobStatusActive,
			PublishStatus: aggregator.JobPublishStatusPublished,
			URL:           "https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
			Title:         "Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)",
			Description:   "<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
			Source:        aggregator.IntegrationArbeitnow.String(),
			Location:      "Munich",
			Remote:        true,
			PostedAt:      time.Unix(1739357344, 0),
			CreatedAt:     time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
			UpdatedAt:     time.Date(2025, 1, 1, 0, 2, 0, 0, time.UTC),
		}
		for _, opt := range opts {
			opt(j)
		}
		dsl.JobRepository.Add(j)
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
		dsl.Config = dsl.defaultConfig()
	}
	if dsl.Logger == nil {
		dsl.LogBuffer, dsl.Logger = NewLogger()
	}
	if dsl.PubSubJobService == nil {
		dsl.PubSubJobService = NewPubSubJobService()
	}
	if dsl.ImportService == nil {
		dsl.ImportService = importing.NewService(dsl.ChannelRepository, dsl.ImportRepository, dsl.JobRepository, dsl.HTTPClient, *dsl.Config, dsl.PubSubJobService, dsl.Logger)
	}
	if dsl.PubSubImportService == nil {
		dsl.PubSubImportService = NewPubSubImportService()
	}
	if dsl.SchedulingService == nil {
		dsl.SchedulingService = scheduling.NewService(dsl.ImportRepository, dsl.ChannelRepository, dsl.PubSubImportService, dsl.Logger)
	}

	if dsl.HTTPConfig == nil {
		dsl.HTTPConfig = &http.Config{}
	}

	if dsl.APIServer == nil {
		dsl.APIServer = http.APIRootHandler(dsl.ConfiguringService, dsl.ChannelRepository, dsl.ImportRepository, dsl.SchedulingService, *dsl.HTTPConfig, dsl.Logger)
	}

	if dsl.ImportServer == nil {
		dsl.ImportServer = http.ImportRootHandler(dsl.ImportService, dsl.Logger)
	}

	return dsl
}

func (dsl *DSL) defaultConfig() *importing.Config {
	return &importing.Config{
		Import: struct {
			Metric  importing.ConfigWorker
			Job     importing.ConfigWorker
			Publish importing.ConfigWorker
		}{
			Metric: importing.ConfigWorker{
				BufferSize: 10,
				Workers:    10,
			},
			Job: importing.ConfigWorker{
				BufferSize: 10,
				Workers:    10,
			},
			Publish: importing.ConfigWorker{
				BufferSize: 10,
				Workers:    10,
			},
		},
	}
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
	return dsl.PubSubImportService.ImportIDs
}

func (dsl *DSL) ImportMetrics() []*aggregator.ImportMetric {
	var metrics []*aggregator.ImportMetric
	for _, m := range dsl.ImportRepository.ImportMetrics() {
		metrics = append(metrics, m)
	}

	return metrics
}

func (dsl *DSL) ImportMetric(id uuid.UUID) *aggregator.ImportMetric {
	return dsl.ImportRepository.ImportMetrics()[id]
}

func (dsl *DSL) ImportMetricsByJobID(jobID uuid.UUID) map[aggregator.ImportMetricType]int {
	metrics := make(map[aggregator.ImportMetricType]int)
	for _, m := range dsl.ImportRepository.ImportMetrics() {
		if m.JobID == jobID {
			metrics[m.MetricType]++
		}
	}

	return metrics
}

func (dsl *DSL) Jobs() []*aggregator.Job {
	var jobs []*aggregator.Job
	for _, i := range dsl.JobRepository.Jobs {
		jobs = append(jobs, i)
	}

	return jobs
}

func (dsl *DSL) Job(id uuid.UUID) *aggregator.Job {
	return dsl.JobRepository.Jobs[id]
}

func (dsl *DSL) PublishedJobInformations() []*aggregator.Job {
	return dsl.PubSubJobService.JobInformations
}

func (dsl *DSL) PublishedJobInformation(id uuid.UUID) *aggregator.Job {
	for _, j := range dsl.PubSubJobService.JobInformations {
		if j.ID == id {
			return j
		}
	}

	return nil
}

func (dsl *DSL) PublishedJobMissings() []*aggregator.Job {
	return dsl.PubSubJobService.JobMissings
}

func (dsl *DSL) PublishedJobMissing(id uuid.UUID) *aggregator.Job {
	for _, j := range dsl.PubSubJobService.JobMissings {
		if j.ID == id {
			return j
		}
	}

	return nil
}

func (dsl *DSL) LogLines() []string {
	return LogLines(dsl.LogBuffer)
}

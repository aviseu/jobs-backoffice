package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/configuring"
	"github.com/aviseu/jobs-backoffice/internal/app/domain/importing"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) Test_Success() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	jr := testutils.NewJobRepository()
	js := importing.NewJobService(jr, 10, 10)
	ir := testutils.NewImportRepository()
	c := testutils.NewRequestLogger(http.DefaultClient)
	lbuf, log := testutils.NewLogger()
	cfg := importing.Config{
		Arbeitnow: arbeitnow.Config{
			URL: server.URL,
		},
		Import: struct {
			ResultBufferSize int `split_words:"true" default:"10"`
			ResultWorkers    int `split_words:"true" default:"10"`
		}{
			ResultBufferSize: 10,
			ResultWorkers:    10,
		},
	}
	f := importing.NewFactory(
		c,
		cfg,
	)
	chr := testutils.NewChannelRepository()
	ch := configuring.NewChannel(uuid.New(), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	chr.Add(ch.ToAggregator())

	j1 := importing.NewJob(
		uuid.NewSHA1(ch.ID(), []byte("bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288")),
		ch.ID(),
		aggregator.JobStatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
		"Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)",
		"<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
		ch.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished),
	)
	jr.Add(j1.ToDTO())
	j2 := importing.NewJob(
		uuid.New(),
		ch.ID(),
		aggregator.JobStatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/another",
		"bankkaufmann-fur-front",
		"Das Wichtigste für unseren Kunden: Mitarbeiter",
		ch.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		importing.JobWithPublishStatus(aggregator.JobPublishStatusPublished),
	)
	jr.Add(j2.ToDTO())

	i := importing.NewImport(uuid.New(), ch.ID())
	ir.AddImport(i.ToAggregate())

	s := importing.NewService(chr, ir, js, f, cfg, log)

	// Execute
	err := s.Import(context.Background(), i.ID())

	// Assert
	suite.NoError(err)

	// Assert Jobs
	suite.Len(ir.Imports, 1)
	for _, imp := range ir.Imports {
		i = importing.NewImportFromDTO(imp)
	}
	suite.Equal(ch.ID(), i.ChannelID())
	suite.Equal(2, i.NewJobs())
	suite.Equal(0, i.UpdatedJobs())
	suite.Equal(1, i.NoChangeJobs())
	suite.Equal(1, i.MissingJobs())
	suite.Equal(0, i.FailedJobs())

	// Assert import results
	jobs := ir.ImportJobs()
	suite.Len(jobs, 4)
	suite.Equal(aggregator.ImportJobResultNoChange, jobs[j1.ID()].Result)
	suite.Equal(aggregator.ImportJobResultMissing, jobs[j2.ID()].Result)
	suite.Equal(aggregator.ImportJobResultNew, jobs[uuid.NewSHA1(ch.ID(), []byte("bankkaufmann-fur-front-office-middle-office-back-office-munich-304839"))].Result)
	suite.Equal(aggregator.ImportJobResultNew, jobs[uuid.NewSHA1(ch.ID(), []byte("fund-accountant-wertpapierfonds-munich-310570"))].Result)

	// Assert Logs
	suite.Empty(lbuf)
}

func (suite *ServiceSuite) Test_Execute_ImportRepositoryFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ir := testutils.NewImportRepository()
	ir.FailWith(errors.New("boom!"))
	jr := testutils.NewJobRepository()
	js := importing.NewJobService(jr, 10, 10)
	c := testutils.NewHTTPClientMock()
	cfg := importing.Config{
		Import: struct {
			ResultBufferSize int `split_words:"true" default:"10"`
			ResultWorkers    int `split_words:"true" default:"10"`
		}{
			ResultBufferSize: 10,
			ResultWorkers:    10,
		},
	}
	f := importing.NewFactory(
		c,
		cfg,
	)
	s := importing.NewService(chr, ir, js, f, cfg, log)
	id := uuid.New()

	// Execute
	err := s.Import(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to find import "+id.String())
	suite.ErrorContains(err, "boom!")

	// Assert Logs
	suite.Empty(lbuf)
}

func (suite *ServiceSuite) Test_Execute_ChannelServiceFail() {
	// Prepare
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	chr.FailWith(errors.New("boom!"))
	ir := testutils.NewImportRepository()
	jr := testutils.NewJobRepository()
	js := importing.NewJobService(jr, 10, 10)
	c := testutils.NewHTTPClientMock()
	cfg := importing.Config{
		Import: struct {
			ResultBufferSize int `split_words:"true" default:"10"`
			ResultWorkers    int `split_words:"true" default:"10"`
		}{
			ResultBufferSize: 10,
			ResultWorkers:    10,
		},
	}
	f := importing.NewFactory(
		c,
		cfg,
	)
	s := importing.NewService(chr, ir, js, f, cfg, log)
	i := importing.NewImport(uuid.New(), uuid.New())
	ir.AddImport(i.ToAggregate())

	// Execute
	err := s.Import(context.Background(), i.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to find channel "+i.ChannelID().String())
	suite.ErrorContains(err, "boom!")

	// Assert Logs
	suite.Empty(lbuf)
}

func (suite *ServiceSuite) Test_Execute_GatewayFail() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	lbuf, log := testutils.NewLogger()
	chr := testutils.NewChannelRepository()
	ch := configuring.NewChannel(uuid.MustParse(testutils.ArbeitnowMethodNotFound), "channel 1", aggregator.IntegrationArbeitnow, aggregator.ChannelStatusActive)
	chr.Add(ch.ToAggregator())
	ir := testutils.NewImportRepository()
	jr := testutils.NewJobRepository()
	js := importing.NewJobService(jr, 10, 10)
	c := http.DefaultClient
	cfg := importing.Config{
		Arbeitnow: arbeitnow.Config{
			URL: server.URL,
		},
		Import: struct {
			ResultBufferSize int `split_words:"true" default:"10"`
			ResultWorkers    int `split_words:"true" default:"10"`
		}{
			ResultBufferSize: 10,
			ResultWorkers:    10,
		},
	}
	f := importing.NewFactory(
		c,
		cfg,
	)
	s := importing.NewService(chr, ir, js, f, cfg, log)
	i := importing.NewImport(uuid.New(), ch.ID())
	ir.AddImport(i.ToAggregate())

	// Execute
	err := s.Import(context.Background(), i.ID())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to get jobs page 1 on channel "+ch.ID().String())
	suite.ErrorContains(err, "<title>An Error Occurred: Method Not Allowed</title>")

	// Assert Logs
	suite.Empty(lbuf)
}

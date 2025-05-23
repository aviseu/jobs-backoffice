package arbeitnow_test

import (
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/api/arbeitnow"
	"github.com/stretchr/testify/suite"
)

func TestService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
	service *arbeitnow.Service
}

func (suite *ServiceSuite) Test_GetJobs_Success() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	defer server.Close()
	ch := &aggregator.Channel{
		ID:          uuid.New(),
		Name:        "arbeitnow integration",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	c := testutils.NewRequestLogger(http.DefaultClient)
	s := arbeitnow.NewService(c, arbeitnow.Config{URL: server.URL}, ch)

	// Execute
	jobs, err := s.GetJobs()

	// Assert result
	suite.NoError(err)
	suite.Len(jobs, 3)
	suite.Equal(uuid.NewSHA1(ch.ID, []byte("bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288")), jobs[0].ID)
	suite.Equal("Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)", jobs[0].Title)
	suite.Equal("<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>", jobs[0].Description)
	suite.Equal("Munich", jobs[0].Location)
	suite.True(jobs[0].PostedAt.Equal(time.Unix(1739357344, 0)))
	suite.Equal("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288", jobs[0].URL)
	suite.True(jobs[0].Remote)
	suite.Equal(aggregator.JobStatusActive, jobs[0].Status)
	suite.Equal(aggregator.JobPublishStatusUnpublished, jobs[0].PublishStatus)
	suite.Equal(ch.ID, jobs[0].ChannelID)

	// Assert requests made
	suite.Len(c.Logs, 2)

	suite.Equal("GET", c.Logs[0].Method)
	suite.Equal(server.URL+"/api/job-board-api", c.Logs[0].URL)
	suite.NotEmpty(c.Logs[0].Response)

	suite.Equal("GET", c.Logs[1].Method)
	suite.Equal(server.URL+"/api/job-board-api?page=2", c.Logs[1].URL)
	suite.NotEmpty(c.Logs[1].Response)
}

func (suite *ServiceSuite) Test_GetJobs_BadRequestFailed() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	defer server.Close()
	ch := &aggregator.Channel{
		ID:          uuid.MustParse(testutils.ArbeitnowMethodNotFound),
		Name:        "arbeitnow integration",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	c := testutils.NewRequestLogger(http.DefaultClient)
	s := arbeitnow.NewService(c, arbeitnow.Config{URL: server.URL}, ch)

	// Execute
	jobs, err := s.GetJobs()

	// Assert result
	suite.Nil(jobs)
	suite.Error(err)
	suite.ErrorContains(err, `<h2>The server returned a "405 Method Not Allowed".</h2>`)
	suite.ErrorContains(err, "failed to request with http code 405 and body:")
	suite.ErrorContains(err, "failed to get jobs page 1 on channel 3fae894d-3484-4274-b337-fcd35a9f135c")
}

func (suite *ServiceSuite) Test_GetJobs_ClientError() {
	// Prepare
	chID := uuid.New()
	ch := &aggregator.Channel{
		ID:          chID,
		Name:        "arbeitnow integration",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	m := testutils.NewHTTPClientMock()
	c := testutils.NewRequestLogger(m)
	s := arbeitnow.NewService(c, arbeitnow.Config{}, ch)

	m.On("Do", mock.Anything).Return(nil, errors.New("something bad happened")).Once()

	// Execute
	jobs, err := s.GetJobs()

	// Assert result
	suite.Nil(jobs)
	suite.Error(err)
	suite.ErrorContains(err, "failed to get jobs page 1 on channel "+chID.String())
	suite.ErrorContains(err, "failed to get jobBoard: something bad happened")
}

func (suite *ServiceSuite) Test_GetJobs_InvalidResponse() {
	// Prepare
	chID := uuid.New()
	ch := &aggregator.Channel{
		ID:          chID,
		Name:        "arbeitnow integration",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	m := testutils.NewHTTPClientMock()
	c := testutils.NewRequestLogger(m)
	s := arbeitnow.NewService(c, arbeitnow.Config{}, ch)

	m.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Body:       io.NopCloser(strings.NewReader("<html><title>An Error Occurred</title></html>")),
	}, nil).Once()

	// Execute
	jobs, err := s.GetJobs()

	// Assert result
	suite.Nil(jobs)
	suite.Error(err)
	suite.ErrorContains(err, "failed to get jobs page 1 on channel "+chID.String())
	suite.ErrorContains(err, "failed to decode response body")
	suite.ErrorContains(err, "<html><title>An Error Occurred</title></html>")
}

func (suite *ServiceSuite) Test_GetJobs_FailedResponseWithoutBody() {
	// Prepare
	chID := uuid.New()
	ch := &aggregator.Channel{
		ID:          chID,
		Name:        "arbeitnow integration",
		Integration: aggregator.IntegrationArbeitnow,
		Status:      aggregator.ChannelStatusActive,
	}
	m := testutils.NewHTTPClientMock()
	c := testutils.NewRequestLogger(m)
	s := arbeitnow.NewService(c, arbeitnow.Config{}, ch)

	m.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: http.StatusInternalServerError,
		Status:     http.StatusText(http.StatusInternalServerError),
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil).Once()

	// Execute
	jobs, err := s.GetJobs()

	// Assert result
	suite.Nil(jobs)
	suite.Error(err)
	suite.ErrorContains(err, "failed to get jobs page 1 on channel "+chID.String())
	suite.ErrorContains(err, "failed to request with http code 500 and no body")
}

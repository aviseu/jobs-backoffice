package gateway_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/aviseu/jobs/internal/app/gateway"
	"github.com/aviseu/jobs/internal/app/gateway/arbeitnow"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestGateway(t *testing.T) {
	suite.Run(t, new(GatewaySuite))
}

type GatewaySuite struct {
	suite.Suite
}

func (suite *GatewaySuite) Test_ImportChannel_Success() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	r := testutils.NewJobRepository()
	s := job.NewService(r)
	c := testutils.NewRequestLogger(http.DefaultClient)
	f := gateway.NewFactory(
		s,
		c,
		gateway.Config{
			Arbeitnow: arbeitnow.Config{
				URL: server.URL,
			},
		},
	)
	ch := channel.New(uuid.New(), "channel", channel.IntegrationArbeitnow, channel.StatusActive)
	gw := f.Create(ch)

	j1 := job.New(
		uuid.NewSHA1(ch.ID(), []byte("bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288")),
		ch.ID(),
		job.StatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
		"Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)",
		"<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
		ch.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	r.Add(j1)
	j2 := job.New(
		uuid.New(),
		ch.ID(),
		job.StatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/another",
		"bankkaufmann-fur-front",
		"Das Wichtigste für unseren Kunden: Mitarbeiter",
		ch.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	r.Add(j2)

	// Execute
	err := gw.ImportChannel(context.Background())

	// Assert
	suite.NoError(err)
	suite.Len(r.Jobs, 4)
	suite.Equal(job.PublishStatusPublished, r.Jobs[j1.ID()].PublishStatus())
	suite.Equal(job.PublishStatusUnpublished, r.Jobs[j2.ID()].PublishStatus())
}

func (suite *GatewaySuite) Test_ImportChannel_RepositoryFail() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	r := testutils.NewJobRepository()
	r.FailWith(errors.New("boom!"))
	s := job.NewService(r)
	c := testutils.NewRequestLogger(http.DefaultClient)
	f := gateway.NewFactory(
		s,
		c,
		gateway.Config{
			Arbeitnow: arbeitnow.Config{
				URL: server.URL,
			},
		},
	)
	ch := channel.New(uuid.New(), "channel", channel.IntegrationArbeitnow, channel.StatusActive)
	gw := f.Create(ch)

	// Execute
	err := gw.ImportChannel(context.Background())

	// Assert
	suite.Error(err)
	suite.Equal("failed to sync jobs for channel "+ch.ID().String()+": failed to get existing jobs: boom!", err.Error())
}

func (suite *GatewaySuite) Test_ImportChannel_ServerFail() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	r := testutils.NewJobRepository()
	s := job.NewService(r)
	c := testutils.NewRequestLogger(http.DefaultClient)
	f := gateway.NewFactory(
		s,
		c,
		gateway.Config{
			Arbeitnow: arbeitnow.Config{
				URL: server.URL,
			},
		},
	)
	ch := channel.New(uuid.MustParse(testutils.ArbeitnowMethodNotFound), "channel", channel.IntegrationArbeitnow, channel.StatusActive)
	gw := f.Create(ch)

	// Execute
	err := gw.ImportChannel(context.Background())

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to import channel "+ch.ID().String())
	suite.ErrorContains(err, `<h2>The server returned a "405 Method Not Allowed".</h2>`)
}

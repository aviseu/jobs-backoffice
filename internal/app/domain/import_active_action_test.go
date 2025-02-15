package domain_test

import (
	"context"
	"github.com/aviseu/jobs/internal/app/domain"
	"github.com/aviseu/jobs/internal/app/domain/channel"
	"github.com/aviseu/jobs/internal/app/domain/imports"
	"github.com/aviseu/jobs/internal/app/domain/job"
	"github.com/aviseu/jobs/internal/app/gateway"
	"github.com/aviseu/jobs/internal/app/gateway/arbeitnow"
	"github.com/aviseu/jobs/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

func TestImportActive(t *testing.T) {
	suite.Run(t, new(ImportActiveSuite))
}

type ImportActiveSuite struct {
	suite.Suite
}

func (suite *ImportActiveSuite) Test_ImportActive_Success() {
	// Prepare
	server := testutils.NewArbeitnowServer()
	jr := testutils.NewJobRepository()
	js := job.NewService(jr, 10, 10)
	ir := testutils.NewImportRepository()
	is := imports.NewService(ir)
	c := testutils.NewRequestLogger(http.DefaultClient)
	lbuf, log := testutils.NewLogger()
	f := gateway.NewFactory(
		js,
		is,
		c,
		gateway.Config{
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
		},
		log,
	)
	chr := testutils.NewChannelRepository()
	ch1 := channel.New(uuid.New(), "channel 1", channel.IntegrationArbeitnow, channel.StatusActive)
	chr.Channels[ch1.ID()] = ch1
	ch2 := channel.New(uuid.New(), "channel 2", channel.IntegrationArbeitnow, channel.StatusInactive)
	chr.Channels[ch2.ID()] = ch2
	chs := channel.NewService(chr)

	j1 := job.New(
		uuid.NewSHA1(ch1.ID(), []byte("bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288")),
		ch1.ID(),
		job.StatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
		"Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)",
		"<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
		ch1.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	jr.Add(j1)
	j2 := job.New(
		uuid.New(),
		ch1.ID(),
		job.StatusActive,
		"https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/another",
		"bankkaufmann-fur-front",
		"Das Wichtigste für unseren Kunden: Mitarbeiter",
		ch1.Integration().String(),
		"Munich",
		true,
		time.Unix(1739357344, 0),
		job.WithPublishStatus(job.PublishStatusPublished),
	)
	jr.Add(j2)

	action := domain.NewImportActiveAction(chs, f, log)

	// Execute
	err := action.Execute(context.Background())

	// Assert
	suite.NoError(err)

	// Assert Jobs
	suite.Len(ir.Imports, 1)
	var i *imports.Import
	for _, imp := range ir.Imports {
		i = imp
	}
	suite.Equal(ch1.ID(), i.ChannelID())
	suite.Equal(2, i.NewJobs())
	suite.Equal(0, i.UpdatedJobs())
	suite.Equal(1, i.NoChangeJobs())
	suite.Equal(1, i.MissingJobs())
	suite.Equal(0, i.FailedJobs())

	// Assert import results
	suite.Len(ir.JobResults, 4)
	suite.Equal(imports.JobStatusNoChange, ir.JobResults[j1.ID()].Result())
	suite.Equal(imports.JobStatusMissing, ir.JobResults[j2.ID()].Result())
	suite.Equal(imports.JobStatusNew, ir.JobResults[uuid.NewSHA1(ch1.ID(), []byte("bankkaufmann-fur-front-office-middle-office-back-office-munich-304839"))].Result())
	suite.Equal(imports.JobStatusNew, ir.JobResults[uuid.NewSHA1(ch1.ID(), []byte("fund-accountant-wertpapierfonds-munich-310570"))].Result())

	// Assert Logs
	lines := testutils.LogLines(lbuf)
	suite.Len(lines, 1)
	suite.Contains(lines[0], `"level":"INFO"`)
	suite.Contains(lines[0], "importing channel "+ch1.ID().String()+" [arbeitnow] [name: channel 1]")
}

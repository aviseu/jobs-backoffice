package importing_test

import (
	"context"
	"errors"
	"github.com/aviseu/jobs-backoffice/internal/app/infrastructure/aggregator"
	"github.com/aviseu/jobs-backoffice/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceSuite))
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) Test_Success() {
	// Prepare
	chID := uuid.New()
	j1ID := uuid.NewSHA1(chID, []byte("bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288"))
	j2ID := uuid.New()
	iID := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithArbeitnowEnabled(),
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithJob(
			testutils.WithJobID(j1ID),
			testutils.WithJobChannelID(chID),
			testutils.WithJobStatus(aggregator.JobStatusActive),
			testutils.WithJobPublishStatus(aggregator.JobPublishStatusPublished),
			testutils.WithJobURL("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288"),
			testutils.WithJobTitle("Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)"),
			testutils.WithJobDescription("<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>"),
			testutils.WithJobSource(aggregator.IntegrationArbeitnow.String()),
			testutils.WithJobLocation("Munich"),
			testutils.WithJobRemote(true),
			testutils.WithJobPostedAt(time.Unix(1739357344, 0)),
		),
		testutils.WithJob(
			testutils.WithJobID(j2ID),
			testutils.WithJobChannelID(chID),
			testutils.WithJobStatus(aggregator.JobStatusActive),
			testutils.WithJobPublishStatus(aggregator.JobPublishStatusPublished),
			testutils.WithJobURL("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/another"),
			testutils.WithJobTitle("bankkaufmann-fur-front"),
			testutils.WithJobDescription("Das Wichtigste für unseren Kunden: Mitarbeiter"),
			testutils.WithJobSource(aggregator.IntegrationArbeitnow.String()),
			testutils.WithJobLocation("Berlin"),
			testutils.WithJobRemote(true),
			testutils.WithJobPostedAt(time.Unix(1739357344, 0)),
		),
		testutils.WithImport(
			testutils.WithImportID(iID),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
		),
	)

	// Execute
	err := dsl.ImportService.Import(context.Background(), iID)

	// Assert
	suite.NoError(err)

	// Assert Import
	suite.Len(dsl.Imports(), 1)
	dbImport := dsl.FirstImport()
	suite.Equal(chID, dbImport.ChannelID)
	suite.Equal(2, dbImport.NewJobs())
	suite.Equal(0, dbImport.UpdatedJobs())
	suite.Equal(1, dbImport.NoChangeJobs())
	suite.Equal(1, dbImport.MissingJobs())
	suite.Equal(0, dbImport.FailedJobs())

	// Assert import results
	importJobs := dsl.ImportMetrics()
	suite.Len(importJobs, 4)
	suite.Equal(aggregator.ImportMetricTypeNoChange, dsl.ImportMetric(j1ID).MetricType)
	suite.Equal(aggregator.ImportMetricTypeMissing, dsl.ImportMetric(j2ID).MetricType)
	jNew1ID := uuid.NewSHA1(chID, []byte("bankkaufmann-fur-front-office-middle-office-back-office-munich-304839"))
	suite.Equal(aggregator.ImportMetricTypeNew, dsl.ImportMetric(jNew1ID).MetricType)
	jNew2ID := uuid.NewSHA1(chID, []byte("fund-accountant-wertpapierfonds-munich-310570"))
	suite.Equal(aggregator.ImportMetricTypeNew, dsl.ImportMetric(jNew2ID).MetricType)

	// Assert Job
	jobs := dsl.Jobs()
	suite.Len(jobs, 4)

	// no change
	suite.Equal(j1ID, dsl.Job(j1ID).ID)
	suite.Equal(chID, dsl.Job(j1ID).ChannelID)
	suite.Equal(aggregator.JobStatusActive, dsl.Job(j1ID).Status)
	suite.Equal(aggregator.JobPublishStatusPublished, dsl.Job(j1ID).PublishStatus)
	suite.Equal("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288", dsl.Job(j1ID).URL)
	suite.Equal("Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)", dsl.Job(j1ID).Title)
	suite.Equal("<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>", dsl.Job(j1ID).Description)
	suite.Equal(aggregator.IntegrationArbeitnow.String(), dsl.Job(j1ID).Source)
	suite.Equal("Munich", dsl.Job(j1ID).Location)
	suite.True(dsl.Job(j1ID).Remote)
	suite.Equal(time.Unix(1739357344, 0), dsl.Job(j1ID).PostedAt)

	// missing
	suite.Equal(j2ID, dsl.Job(j2ID).ID)
	suite.Equal(chID, dsl.Job(j2ID).ChannelID)
	suite.Equal(aggregator.JobStatusInactive, dsl.Job(j2ID).Status)
	suite.Equal(aggregator.JobPublishStatusUnpublished, dsl.Job(j2ID).PublishStatus)
	suite.Equal("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/another", dsl.Job(j2ID).URL)
	suite.Equal("bankkaufmann-fur-front", dsl.Job(j2ID).Title)
	suite.Equal("Das Wichtigste für unseren Kunden: Mitarbeiter", dsl.Job(j2ID).Description)
	suite.Equal(aggregator.IntegrationArbeitnow.String(), dsl.Job(j2ID).Source)
	suite.Equal("Berlin", dsl.Job(j2ID).Location)
	suite.True(dsl.Job(j2ID).Remote)
	suite.Equal(time.Unix(1739357344, 0), dsl.Job(j2ID).PostedAt)

	// new 1
	suite.Equal(jNew1ID, dsl.Job(jNew1ID).ID)
	suite.Equal(chID, dsl.Job(jNew1ID).ChannelID)
	suite.Equal(aggregator.JobStatusActive, dsl.Job(jNew1ID).Status)
	suite.Equal(aggregator.JobPublishStatusPublished, dsl.Job(jNew1ID).PublishStatus)
	suite.Equal("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkaufmann-fur-front-office-middle-office-back-office-munich-304839", dsl.Job(jNew1ID).URL)
	suite.Equal("Bankkaufmann (m/w/d) für Front Office | Middle Office | Back Office", dsl.Job(jNew1ID).Title)
	suite.Equal("<p>Das Wichtigste für unseren Kunden: Mitarbeiter, auf die er sich verlassen kann. Und dieses Vertrauen zahlt sich aus. Auch für Sie. Neben kurzen Entscheidungswegen profitieren Sie von ausgezeichneten Entwicklungsmöglichkeiten und einem kollegialen Umfeld. Bewerben Sie sich bei uns als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<p>Sie haben Ihre Ausbildung als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong> bereits beendet und wollen erste Berufserfahrungen sammeln? Oder sind Sie bereits Spezialist im Bereich Bank- und Finanzwesen und suchen einen neuen Wirkungskreis in München? Wir werden Ihnen dabei helfen!</p>\n<p>Bei unserem Kunden handelt sich um ein Unternehmen im Bereich Finanzdienstleistungen mit Hauptsitz in München.</p>\n<p><strong>Starten Sie Ihre Karriere und bewerben Sie sich bei uns!</strong></p>\n<p>Wir bieten Ihnen bei einer renommierten Bank Positionen in den unterschiedlichen Bereichen - Front Office | Middle Office | Back Office:</p>\n<p>...z.B. im <strong>Back Office</strong>:</p>\n<ul>\n<li>Sachbearbeitung im Kontenservice oder Vertragsservice oder Depotservice</li>\n<li>Dokumentensachbearbeiter (m/w/d)</li>\n<li>Unterstützung in der qualifizierten Kreditsachbearbeitung<br>\nWertpapierabwicklung im nationalen oder internationalen Umfeld</li>\n<li>Wertpapiersachbearbeiter (m/w/d)</li>\n<li>Mitarbeiter Meldewesen</li>\n<li>Mitarbeiter im Zahlungsverkehr (m/w/d) </li>\n<li>Mitarbeiter Compliance</li>\n<li>…oder im <strong>Front Office</strong>:</li>\n<li>Kundenbetreuer mit Entwicklungspotential zum Privatkundenberater</li>\n<li>Serviceberater</li>\n<li>Assistenz Private Banking (m/w/d)</li>\n<li>Spezieller Fokus auf Finanzierungsberatung</li>\n<li>Einstieg in die Geschäftskundenberatung</li>\n<li>Direkt in die Privatkundenberatung und Firmenkundenberatung</li>\n<li>Vermögenskundenbetreuer | Private Banker (m/w/d)</li>\n</ul>\n<h2>Qualifikation</h2>\n<p>⭐ Abgeschlossene Bankausbildung oder relevantes wirtschaftliches Studium</p>\n<p>⭐ Erste Berufserfahrungen im Finanzbereich sind von Vorteil</p>\n<p>⭐ IT-Affinität und gute MS-Office Kenntnisse</p>\n<p>⭐ Sehr gute Deutschkenntnisse</p>\n<h2>Benefits</h2>\n<p>Als Ex­perten für Personal­aus­wahl stehen wir Ihnen bei der Suche nach einer neuen Heraus­forderung zur Seite. <strong>OPUS ONE</strong> bietetIhnen eine Viel­zahl an Mög­lich­keiten, sich beruflich zu ver­ändern oder weiter­zu­entwickeln. Pro­fitieren von unseren Kennt­nissen in Ihrem Berufs­feld und Ihrer Branche. Eine Viel­zahl an Jobs ist verfügbar!</p>\n<p><strong>So profitieren Sie mit uns:</strong></p>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen für das Vorstellungsgespräch beim Unternehme</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen)</li>\n<li>Persönliches oder telefonisches Interview mit anschließendem Karrierecoaching</li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen</li>\n<li>Beratung zum Arbeitsvertrag des neuen ArbeitgeberSelbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n</ul>\n<p><strong>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos.</strong></p>\n<p>Werden Sie aktiv! Wir freuen uns darauf, Sie kennen zu lernen! Senden Sie Ihre aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Ihrem Gehaltswunsch sowie Ihrem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong> die Richtige für Sie ist und ob wir Ihnen außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>IHR ANSPRECHPARTNER:</strong></p>\n<p>Herr Florian Fendt</p>\n<p>Tel.: 089 890 648 127</p>\n<p>Find more <a href=\"https://www.arbeitnow.com/english-speaking-jobs\">English Speaking Jobs in Germany</a> on Arbeitnow</a>", dsl.Job(jNew1ID).Description)
	suite.Equal(aggregator.IntegrationArbeitnow.String(), dsl.Job(jNew1ID).Source)
	suite.Equal("Munich", dsl.Job(jNew1ID).Location)
	suite.False(dsl.Job(jNew1ID).Remote)
	suite.Equal(time.Unix(1739357344, 0), dsl.Job(jNew1ID).PostedAt)

	// new 2
	suite.Equal(jNew2ID, dsl.Job(jNew2ID).ID)
	suite.Equal(chID, dsl.Job(jNew2ID).ChannelID)
	suite.Equal(aggregator.JobStatusActive, dsl.Job(jNew2ID).Status)
	suite.Equal(aggregator.JobPublishStatusPublished, dsl.Job(jNew2ID).PublishStatus)
	suite.Equal("https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/fund-accountant-wertpapierfonds-munich-310570", dsl.Job(jNew2ID).URL)
	suite.Equal("Fund Accountant Wertpapierfonds (m/w/d)", dsl.Job(jNew2ID).Title)
	suite.Equal("<p>Unser Kunde gehört zu einem der größten Marktteilnehmer im Bereich der Wertpapierabwicklung und -verwahrung. Hier wird der Fokus daraufgelegt, seinen Kunden einen allumfassenden, individuellen Service anbieten zu können. Flexibilität und Sorgfalt werden hier großgeschrieben. Langjährige Erfahrung und die Kooperation mit vielen Unternehmen im Finanzumfeld zeichnen diesen Bereich des Unternehmens aus. Das Unternehmen bearbeitet die Themengebiete mit seinen qualifizierten Mitarbeitern, einer leistungsfähigen IT-Landschaft und dem gelebten Servicegedanken für alle Geschäftspartner. </p>\n<p>Profitieren Sie von flexiblen Arbeitszeiten mit Homeoffice-Option sowie Aufstiegs- und Weiterbildungsmöglichkeiten. Zudem bietet das Unternehmen familienfreundlichen Arbeitsbedingungen, Sportangebote und abwechslungsreichen Aufgaben. Dies macht unseren Kunden zum Top-Arbeitgeber für Sie. </p>\n<p>Also nutzen Sie die Chance und bewerben Sie sich jetzt!</p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Kontrolle von Differenzen und Absprache mit Fachabteilungen und externen Serviceprovidern</li>\n<li>Verantwortlich für Fondsmigrationen, -verschmelzungen oder -schließungen sowie für die korrekte Verbuchung sämtlicher Geschäftsvorfälle für den Fonds</li>\n<li>Berechnung der Anteilpreise für Publikums und Spezialfonds</li>\n<li>Bearbeitung relevanter Kapitalmaßnahmen </li>\n<li>Tatkräftige Unterstützung bei Projekten</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung im Bank- oder im Investmentfondsbereich oder eine vergleichbare Qualifikation</li>\n<li>Erste Berufserfahrung im Umgang mit Wertpapierfonds und im Finanzproduktbereich</li>\n<li>Sichere Handhabung mit Wertpapierfondsbuchhaltungssystemen</li>\n<li>Analytische und selbständige Arbeitsweise</li>\n<li>Sehr gute Deutsch- und Englischkenntnisse; Französischkenntnisse von Vorteil.</li>\n</ul>\n<h2>Benefits</h2>\n<p>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen)</p>\n<ul>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching</li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen</li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen</li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers</li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Werden Sie aktiv! Wir freuen uns darauf, Sie kennen zu lernen! Senden Sie Ihre aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Ihrem Gehaltswunsch sowie Ihrem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Sie ist und ob wir Ihnen außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>IHR ANSPRECHPARTNER:</strong></p>\n<p>Herr Florian Fendt</p>\n<p>Tel.: 089/ 890 648 127</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>", dsl.Job(jNew2ID).Description)
	suite.Equal(aggregator.IntegrationArbeitnow.String(), dsl.Job(jNew2ID).Source)
	suite.Equal("Munich", dsl.Job(jNew2ID).Location)
	suite.False(dsl.Job(jNew2ID).Remote)
	suite.Equal(time.Unix(1739357344, 0), dsl.Job(jNew2ID).PostedAt)

	// Assert publish
	suite.Len(dsl.PublishedJobs(), 2)
	suite.NotNil(dsl.PublishedJob(jNew1ID)) // new
	suite.NotNil(dsl.PublishedJob(jNew2ID)) // new

	// Assert Logs
	logs := dsl.LogLines()
	suite.Len(logs, 1)
	suite.Contains(logs[0], `"level":"INFO"`)
	suite.Contains(logs[0], "published 2 jobs for import "+iID.String())
}

func (suite *ServiceSuite) Test_Execute_ImportRepositoryFail() {
	// Prepare
	dsl := testutils.NewDSL(
		testutils.WithImportRepositoryError(errors.New("boom")),
	)
	id := uuid.New()

	// Execute
	err := dsl.ImportService.Import(context.Background(), id)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to find import "+id.String())
	suite.ErrorContains(err, "boom")

	// Assert Logs
	suite.Empty(dsl.LogLines())
}

func (suite *ServiceSuite) Test_Execute_ChannelRepositoryFail() {
	// Prepare
	chID := uuid.New()
	iID := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithArbeitnowEnabled(),
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(iID),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
		),
		testutils.WithChannelRepositoryError(errors.New("boom")),
	)

	// Execute
	err := dsl.ImportService.Import(context.Background(), iID)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to find channel "+chID.String())
	suite.ErrorContains(err, "boom")

	// Assert Logs
	suite.Empty(dsl.LogLines())
}

func (suite *ServiceSuite) Test_Execute_GatewayFail() {
	// Prepare
	chID := uuid.MustParse(testutils.ArbeitnowMethodNotFound)
	iID := uuid.New()
	dsl := testutils.NewDSL(
		testutils.WithArbeitnowEnabled(),
		testutils.WithChannel(
			testutils.WithChannelID(chID),
			testutils.WithChannelIntegration(aggregator.IntegrationArbeitnow),
		),
		testutils.WithImport(
			testutils.WithImportID(iID),
			testutils.WithImportChannelID(chID),
			testutils.WithImportStatus(aggregator.ImportStatusPending),
		),
	)

	// Execute
	err := dsl.ImportService.Import(context.Background(), iID)

	// Assert
	suite.Error(err)
	suite.ErrorContains(err, "failed to get jobs page 1 on channel "+chID.String())
	suite.ErrorContains(err, "<title>An Error Occurred: Method Not Allowed</title>")

	// Assert Logs
	suite.Empty(dsl.LogLines())
}

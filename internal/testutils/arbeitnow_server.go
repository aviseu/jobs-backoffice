package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gopkg.in/guregu/null.v3"
)

const (
	pageSize = 2

	ArbeitnowMethodNotFound = "3fae894d-3484-4274-b337-fcd35a9f135c"
)

type jobEntry struct {
	Slug        string   `json:"slug"`
	CompanyName string   `json:"company_name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Remote      bool     `json:"remote"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
	JobTypes    []string `json:"job_types"`
	Location    string   `json:"location"`
	CreatedAt   int64    `json:"created_at"`
}

type JobBoardResponse struct {
	Data  []*jobEntry `json:"data"`
	Links struct {
		Next null.String `json:"next"`
	} `json:"links"`
}

type ArbeitnowServer struct{}

func NewArbeitnowServer() *httptest.Server {
	r := chi.NewRouter()
	var s *httptest.Server

	r.Get("/api/job-board-api", func(w http.ResponseWriter, r *http.Request) {
		page := 1
		if r.URL.Query().Has("page") {
			p, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			page = p
		}

		if r.Header.Get("X-Channel-Id") == ArbeitnowMethodNotFound {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(arbeitnowMethodNotAllowedResponse)
			return
		}

		data := arbeitnowData()
		// paginate data based on page and pageSize and length
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > len(data) {
			end = len(data)
		}
		jobs := data[start:end]
		next := ""

		isLastPage := len(data) <= page*pageSize
		if !isLastPage {
			next = s.URL + "/api/job-board-api?page=" + strconv.Itoa(page+1)
		}

		w.Header().Set("Content-Type", "application/json")
		resp := JobBoardResponse{
			Data: jobs,
			Links: struct {
				Next null.String `json:"next"`
			}{
				Next: null.NewString(next, !isLastPage),
			},
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	s = httptest.NewServer(r)
	return s
}

func arbeitnowData() []*jobEntry {
	return []*jobEntry{
		{
			Slug:        "bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
			CompanyName: "OPUS ONE Recruitment GmbH",
			Title:       "Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)",
			Description: "<p>Unser Kunde ist im Bereich Vermögensverwaltung und Fondmanagement ein führender Finanzdienstleister mit Sitz in München. Als zuverlässiger Partner unabhängiger Vermögensberater und ausgewählter institutioneller Kunden verfügt das Unternehmen über ein Verwaltungsvolumen mehrerer Mrd. EUR. Mit derzeit über 40 Mitarbeitern befasst sich das Unternehmen um alle Vermögensbelange seines Kunden. Nachhaltige Qualität und Kundenzufriedenheit stehen im Mittelpunkt des Unternehmens.</p>\n<p>Wir freuen uns auf Ihre Bewerbung als</p>\n<p><strong>Bankkauffrau im Bereich Zahlungsverkehr und Kontolöschung (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Überprüfung und Dokumentation von Daueraufträgen sowie (Dauer)-Lastschriften.</li>\n<li>Abwicklung des Zahlungsverkehrs im In- und Ausland.</li>\n<li>Bearbeitung von Nachlasskonten im Zusammenhang mit der Kontolöschung.</li>\n<li>Erfassung interner Kostenrechnungen und Kundenbuchungen.</li>\n<li>Überprüfung und Erfassung von Kontolöschungen. </li>\n<li>Durchführung von Tests für bestehende und neu einzuführende Prozesse.</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung als Bankkaufmann (m/w/d) oder vergleichbare kaufmännische Qualifikation.</li>\n<li>Expertise im nationalen und internationalen Zahlungsverkehr.</li>\n<li>Kenntnisse in der Kundenstammdatenpflege.</li>\n<li>Fähigkeit zur selbstständigen Arbeit sowie analytische Herangehensweise</li>\n<li>Anwendungssicher in MS Office, insbesondere Excel von Vorteil.</li>\n<li>Hohes Maß an sorgfältiger und präziser Arbeitsweise</li>\n</ul>\n<h2>Benefits</h2>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen) </li>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching </li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen </li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen </li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers </li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Wir freuen uns darauf, Dich kennen zu lernen! Sende Deine aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Deinem Gehaltswunsch sowie Deinem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Dich ist und ob wir Dir außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>DEIN ANSPRECHPARTNER:</strong></p>\n<p>Frau Elwira Dabrowska | Tel.: 089/890 648 1039</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
			Remote:      true,
			URL:         "https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkauffrau-im-bereich-zahlungsverkehr-und-kontoloschung-munich-290288",
			Tags: []string{
				"Finance",
			},
			JobTypes:  []string{},
			Location:  "Munich",
			CreatedAt: 1739357344,
		},
		{
			Slug:        "bankkaufmann-fur-front-office-middle-office-back-office-munich-304839",
			CompanyName: "OPUS ONE Recruitment GmbH",
			Title:       "Bankkaufmann (m/w/d) für Front Office | Middle Office | Back Office",
			Description: "<p>Das Wichtigste für unseren Kunden: Mitarbeiter, auf die er sich verlassen kann. Und dieses Vertrauen zahlt sich aus. Auch für Sie. Neben kurzen Entscheidungswegen profitieren Sie von ausgezeichneten Entwicklungsmöglichkeiten und einem kollegialen Umfeld. Bewerben Sie sich bei uns als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong></p>\n<h2>Aufgaben</h2>\n<p>Sie haben Ihre Ausbildung als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong> bereits beendet und wollen erste Berufserfahrungen sammeln? Oder sind Sie bereits Spezialist im Bereich Bank- und Finanzwesen und suchen einen neuen Wirkungskreis in München? Wir werden Ihnen dabei helfen!</p>\n<p>Bei unserem Kunden handelt sich um ein Unternehmen im Bereich Finanzdienstleistungen mit Hauptsitz in München.</p>\n<p><strong>Starten Sie Ihre Karriere und bewerben Sie sich bei uns!</strong></p>\n<p>Wir bieten Ihnen bei einer renommierten Bank Positionen in den unterschiedlichen Bereichen - Front Office | Middle Office | Back Office:</p>\n<p>...z.B. im <strong>Back Office</strong>:</p>\n<ul>\n<li>Sachbearbeitung im Kontenservice oder Vertragsservice oder Depotservice</li>\n<li>Dokumentensachbearbeiter (m/w/d)</li>\n<li>Unterstützung in der qualifizierten Kreditsachbearbeitung<br>\nWertpapierabwicklung im nationalen oder internationalen Umfeld</li>\n<li>Wertpapiersachbearbeiter (m/w/d)</li>\n<li>Mitarbeiter Meldewesen</li>\n<li>Mitarbeiter im Zahlungsverkehr (m/w/d) </li>\n<li>Mitarbeiter Compliance</li>\n<li>…oder im <strong>Front Office</strong>:</li>\n<li>Kundenbetreuer mit Entwicklungspotential zum Privatkundenberater</li>\n<li>Serviceberater</li>\n<li>Assistenz Private Banking (m/w/d)</li>\n<li>Spezieller Fokus auf Finanzierungsberatung</li>\n<li>Einstieg in die Geschäftskundenberatung</li>\n<li>Direkt in die Privatkundenberatung und Firmenkundenberatung</li>\n<li>Vermögenskundenbetreuer | Private Banker (m/w/d)</li>\n</ul>\n<h2>Qualifikation</h2>\n<p>⭐ Abgeschlossene Bankausbildung oder relevantes wirtschaftliches Studium</p>\n<p>⭐ Erste Berufserfahrungen im Finanzbereich sind von Vorteil</p>\n<p>⭐ IT-Affinität und gute MS-Office Kenntnisse</p>\n<p>⭐ Sehr gute Deutschkenntnisse</p>\n<h2>Benefits</h2>\n<p>Als Ex­perten für Personal­aus­wahl stehen wir Ihnen bei der Suche nach einer neuen Heraus­forderung zur Seite. <strong>OPUS ONE</strong> bietetIhnen eine Viel­zahl an Mög­lich­keiten, sich beruflich zu ver­ändern oder weiter­zu­entwickeln. Pro­fitieren von unseren Kennt­nissen in Ihrem Berufs­feld und Ihrer Branche. Eine Viel­zahl an Jobs ist verfügbar!</p>\n<p><strong>So profitieren Sie mit uns:</strong></p>\n<ul>\n<li>Sie bewerben sich einmal bei uns und wir übernehmen die Suche nach einem passenden Job für Sie</li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen für das Vorstellungsgespräch beim Unternehme</li>\n<li>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen)</li>\n<li>Persönliches oder telefonisches Interview mit anschließendem Karrierecoaching</li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen</li>\n<li>Beratung zum Arbeitsvertrag des neuen ArbeitgeberSelbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n</ul>\n<p><strong>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos.</strong></p>\n<p>Werden Sie aktiv! Wir freuen uns darauf, Sie kennen zu lernen! Senden Sie Ihre aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Ihrem Gehaltswunsch sowie Ihrem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position als <strong>Bankkaufmann | Bankkauffrau (m/w/d)</strong> die Richtige für Sie ist und ob wir Ihnen außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>IHR ANSPRECHPARTNER:</strong></p>\n<p>Herr Florian Fendt</p>\n<p>Tel.: 089 890 648 127</p>\n<p>Find more <a href=\"https://www.arbeitnow.com/english-speaking-jobs\">English Speaking Jobs in Germany</a> on Arbeitnow</a>",
			Remote:      false,
			URL:         "https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/bankkaufmann-fur-front-office-middle-office-back-office-munich-304839",
			Tags: []string{
				"Finance",
			},
			JobTypes:  []string{},
			Location:  "Munich",
			CreatedAt: 1739357344,
		},

		{
			Slug:        "fund-accountant-wertpapierfonds-munich-310570",
			CompanyName: "OPUS ONE Recruitment GmbH",
			Title:       "Fund Accountant Wertpapierfonds (m/w/d)",
			Description: "<p>Unser Kunde gehört zu einem der größten Marktteilnehmer im Bereich der Wertpapierabwicklung und -verwahrung. Hier wird der Fokus daraufgelegt, seinen Kunden einen allumfassenden, individuellen Service anbieten zu können. Flexibilität und Sorgfalt werden hier großgeschrieben. Langjährige Erfahrung und die Kooperation mit vielen Unternehmen im Finanzumfeld zeichnen diesen Bereich des Unternehmens aus. Das Unternehmen bearbeitet die Themengebiete mit seinen qualifizierten Mitarbeitern, einer leistungsfähigen IT-Landschaft und dem gelebten Servicegedanken für alle Geschäftspartner. </p>\n<p>Profitieren Sie von flexiblen Arbeitszeiten mit Homeoffice-Option sowie Aufstiegs- und Weiterbildungsmöglichkeiten. Zudem bietet das Unternehmen familienfreundlichen Arbeitsbedingungen, Sportangebote und abwechslungsreichen Aufgaben. Dies macht unseren Kunden zum Top-Arbeitgeber für Sie. </p>\n<p>Also nutzen Sie die Chance und bewerben Sie sich jetzt!</p>\n<h2>Aufgaben</h2>\n<ul>\n<li>Kontrolle von Differenzen und Absprache mit Fachabteilungen und externen Serviceprovidern</li>\n<li>Verantwortlich für Fondsmigrationen, -verschmelzungen oder -schließungen sowie für die korrekte Verbuchung sämtlicher Geschäftsvorfälle für den Fonds</li>\n<li>Berechnung der Anteilpreise für Publikums und Spezialfonds</li>\n<li>Bearbeitung relevanter Kapitalmaßnahmen </li>\n<li>Tatkräftige Unterstützung bei Projekten</li>\n</ul>\n<h2>Qualifikation</h2>\n<ul>\n<li>Abgeschlossene Ausbildung im Bank- oder im Investmentfondsbereich oder eine vergleichbare Qualifikation</li>\n<li>Erste Berufserfahrung im Umgang mit Wertpapierfonds und im Finanzproduktbereich</li>\n<li>Sichere Handhabung mit Wertpapierfondsbuchhaltungssystemen</li>\n<li>Analytische und selbständige Arbeitsweise</li>\n<li>Sehr gute Deutsch- und Englischkenntnisse; Französischkenntnisse von Vorteil.</li>\n</ul>\n<h2>Benefits</h2>\n<p>Zugang zum sog. verdeckten Stellenmarkt (nicht ausgeschriebene Stellen)</p>\n<ul>\n<li>Persönliches Interview mit anschließendem individuellem Karrierecoaching</li>\n<li>Eine Vielzahl an Stellen, die kurzfristig besetzt werden müssen</li>\n<li>Persönliche Kontakte zu Entscheidern und hilfreiche Informationen</li>\n<li>Beratung zum Arbeitsvertrag des neuen Arbeitgebers</li>\n<li>Selbstverständlich behandeln wir Ihre persönlichen Daten und Bewerbungsunterlagen absolut vertraulich und diskret</li>\n<li>Unsere Leistung ist für Sie als Bewerber (m/w/d) absolut kostenlos</li>\n</ul>\n<p>Werden Sie aktiv! Wir freuen uns darauf, Sie kennen zu lernen! Senden Sie Ihre aussagekräftigen Bewerbungsunterlagen (mit Angaben zu Ihrem Gehaltswunsch sowie Ihrem nächstmöglichen Eintrittstermin).</p>\n<p>Gemeinsam finden wir heraus, ob diese Position die Richtige für Sie ist und ob wir Ihnen außerdem noch andere Perspektiven anbieten können.</p>\n<p><strong>IHR ANSPRECHPARTNER:</strong></p>\n<p>Herr Florian Fendt</p>\n<p>Tel.: 089/ 890 648 127</p>\n<p>Find <a href=\"https://www.arbeitnow.com/\">Jobs in Germany</a> on Arbeitnow</a>",
			Remote:      false,
			URL:         "https://www.arbeitnow.com/jobs/companies/opus-one-recruitment-gmbh/fund-accountant-wertpapierfonds-munich-310570",
			Tags: []string{
				"Finance",
			},
			JobTypes:  []string{},
			Location:  "Munich",
			CreatedAt: 1739357344,
		},
	}
}

var arbeitnowMethodNotAllowedResponse = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="robots" content="noindex,nofollow,noarchive" />
    <title>An Error Occurred: Method Not Allowed</title>
    <link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 128 128%22><text y=%221.2em%22 font-size=%2296%22>❌</text></svg>" />
    <style>body { background-color: #fff; color: #222; font: 16px/1.5 -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; margin: 0; }
.container { margin: 30px; max-width: 600px; }
h1 { color: #dc3545; font-size: 24px; }
h2 { font-size: 18px; }</style>
</head>
<body>
<div class="container">
    <h1>Oops! An Error Occurred</h1>
    <h2>The server returned a "405 Method Not Allowed".</h2>

    <p>
        Something is broken. Please let us know what you were doing when this error occurred.
        We will fix it as soon as possible. Sorry for any inconvenience caused.
    </p>
</div>
<script>(function(){function c(){var b=a.contentDocument||a.contentWindow.document;if(b){var d=b.createElement('script');d.innerHTML="window.__CF$cv$params={r:'910e8359c9ae3868',t:'MTczOTM4MzU5Mi4wMDAwMDA='};var a=document.createElement('script');a.nonce='';a.src='/cdn-cgi/challenge-platform/scripts/jsd/main.js';document.getElementsByTagName('head')[0].appendChild(a);";b.getElementsByTagName('head')[0].appendChild(d)}}if(document.body){var a=document.createElement('iframe');a.height=1;a.width=1;a.style.position='absolute';a.style.top=0;a.style.left=0;a.style.border='none';a.style.visibility='hidden';document.body.appendChild(a);if('loading'!==document.readyState)c();else if(window.addEventListener)document.addEventListener('DOMContentLoaded',c);else{var e=document.onreadystatechange||function(){};document.onreadystatechange=function(b){e(b);'loading'!==document.readyState&&(document.onreadystatechange=e,c())}}}})();</script></body>
</html>`)

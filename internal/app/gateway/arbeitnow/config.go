package arbeitnow

type Config struct {
	URL string `envconfig:"ARBEITNOW_URL" default:"https://arbeitnow.com"`
}

package arbeitnow

type Config struct {
	URL string `env:"URL" envDefault:"https://arbeitnow.com"`
}

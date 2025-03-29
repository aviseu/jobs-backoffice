package base

type Integration int

const (
	IntegrationArbeitnow Integration = iota
)

var Integrations = map[Integration]string{
	IntegrationArbeitnow: "arbeitnow",
}

func (i Integration) String() string {
	return Integrations[i]
}

func ParseIntegration(s string) (Integration, bool) {
	for _, i := range Integrations {
		if i == s {
			return IntegrationArbeitnow, true
		}
	}

	return -1, false
}

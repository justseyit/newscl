package model

type Provider string

const (
	BBC     Provider = "bbc"
	REUTERS Provider = "reuters"
	INVALID Provider = "invalid"
)

func ToProvider(provider string) Provider {
	switch provider {
	case "bbc":
		return BBC
	case "reuters":
		return REUTERS
	default:
		return INVALID
	}
}

func (p *Provider) String() string {
	switch *p {
	case BBC:
		return "BBC"
	case REUTERS:
		return "Reuters"
	default:
		return "Invalid"
	}
}

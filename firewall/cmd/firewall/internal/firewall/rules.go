package firewall

import (
	"gopkg.in/yaml.v2"
)

type Rules struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Endpoint               string   `yaml:"endpoint"`
	ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
	ForbiddenHeaders       []string `yaml:"forbidden_headers"`
	RequiredHeaders        []string `yaml:"required_headers"`
	MaxRequestLengthBytes  int      `yaml:"max_request_length_bytes"`
	MaxResponseLengthBytes int      `yaml:"max_response_length_bytes"`
	ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
	ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
	ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
}

func ParseRules(data []byte) (*Rules, error) {
	var res Rules
	err := yaml.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

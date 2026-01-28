package scanner

import "time"

type E2EConfig struct {
	Enable                bool
	SlipstreamHealth      string
	SlipstreamFingerprint string
	DNSTTHealth           string
	DNSTTPubKey           string
}

func RunE2E(ip string, timeout time.Duration, cfg E2EConfig) *E2EResult {
	res := &E2EResult{}

	if cfg.SlipstreamHealth != "" {
		ok, err := runSlipstreamE2E(ip, timeout, cfg)
		res.SlipstreamOK = ok
		if err != nil {
			res.Error = err.Error()
			return res
		}
	}

	if !res.SlipstreamOK {
		res.Error = "slipstream e2e failed"
	}

	return res
}

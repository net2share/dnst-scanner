package scanner

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"
)

func runSlipstreamE2E(resolver string, timeout time.Duration, cfg E2EConfig) (bool, error) {
	bin := os.Getenv("DNST_SCANNER_SLIPSTREAM_PATH")
	if bin == "" {
		return false, errors.New("DNST_SCANNER_SLIPSTREAM_PATH not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{
		"--resolver", resolver,
		"--health", cfg.SlipstreamHealth,
	}

	if cfg.SlipstreamFingerprint != "" {
		args = append(args, "--fingerprint", cfg.SlipstreamFingerprint)
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return false, errors.New("slipstream timeout")
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

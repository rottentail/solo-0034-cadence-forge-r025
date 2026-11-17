package inspect

import (
	"cadenceforge/internal/domain"
	"cadenceforge/internal/media"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
)

func Digest(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	h := sha256.New()
	n, err := io.Copy(h, f)
	if err != nil {
		return "", n, err
	}
	return hex.EncodeToString(h.Sum(nil)), n, nil
}

// Artifact accepts probed facts only if they satisfy the delivery contract
// frozen in the render plan. The contract's tolerance and required streams are
// supplied by the plan, not inferred from what the probe happened to observe.
func Artifact(ctx context.Context, prober media.Prober, path string, r domain.Render) (domain.Inspection, error) {
	facts, err := prober.Probe(ctx, path)
	if err != nil {
		return domain.Inspection{}, err
	}
	contract := r.Contract()
	delta := facts.DurationMS - contract.DurationMS
	if tolerance := contract.DurationToleranceMS(); math.Abs(float64(delta)) > float64(tolerance) {
		return domain.Inspection{}, fmt.Errorf("output duration delta %dms exceeds %dms", delta, tolerance)
	}
	if contract.RequireAudio && (!facts.Audio || facts.SampleRate != contract.SampleRate) {
		return domain.Inspection{}, fmt.Errorf("output audio contract mismatch")
	}
	if contract.RequireVideo && (!facts.Video || facts.Width != contract.Width || facts.Height != contract.Height) {
		return domain.Inspection{}, fmt.Errorf("output video contract mismatch")
	}
	digest, n, err := Digest(path)
	if err != nil {
		return domain.Inspection{}, err
	}
	return domain.Inspection{Media: domain.ArtifactMedia{MediaFacts: facts}, Digest: digest, Bytes: n, DurationDeltaMS: delta}, nil
}

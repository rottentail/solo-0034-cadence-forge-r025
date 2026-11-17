package domain

// DeliveryContract is what a successful render must prove about its artifact.
// Every field is a decision made from the frozen plan, recipe and output
// format; none of it can be read off a probe. A probe reports what an object
// happens to contain, while the contract states what the artifact is required
// to contain and how much its measured duration may drift.
//
// The contract is intentionally derived rather than stored: inspection uses the
// same frozen render that encoding used, so the two cannot disagree about the
// target while the persisted state format stays unchanged.
type DeliveryContract struct {
	Format       string
	Width        int
	Height       int
	FPS          int
	SampleRate   int
	DurationMS   int64
	RequireVideo bool
	RequireAudio bool
}

// Contract derives the acceptance requirements from the render's frozen inputs.
func (r Render) Contract() DeliveryContract {
	contract := DeliveryContract{
		Format:     r.Recipe.Format,
		Width:      r.Recipe.Width,
		Height:     r.Recipe.Height,
		FPS:        r.Recipe.FPS,
		SampleRate: r.Recipe.SampleRate,
		DurationMS: r.Plan.DurationMS,
		// MP4 always carries a video stream; WAV is audio-only. Both formats
		// always carry audio because MP4 spans without source audio are
		// rendered as equal-length silence during encoding.
		RequireVideo: r.Recipe.Format == "mp4",
		RequireAudio: true,
	}
	return contract
}

// DurationToleranceMS is the accepted absolute drift between planned and
// measured duration. It mirrors the pre-existing inspection tolerance: 100ms,
// widened to two frames for video so frame quantization is not a failure.
func (c DeliveryContract) DurationToleranceMS() int64 {
	tolerance := int64(100)
	if c.Format == "mp4" {
		tolerance = max(tolerance, int64(2000/c.FPS))
	}
	return tolerance
}

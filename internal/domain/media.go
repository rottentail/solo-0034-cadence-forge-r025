package domain

import "time"

// MediaFacts are the stream properties a probe can observe for an object.
// Source ingestion, encoding decisions and artifact verification all speak
// this vocabulary so they agree on the same measured facts. MediaFacts does
// not assert that an object satisfies a delivery contract; duration tolerance,
// required streams and target dimensions/rates are contract parameters that a
// probe cannot establish on its own.
type MediaFacts struct {
	DurationMS int64  `json:"duration_ms"`
	Video      bool   `json:"video"`
	Audio      bool   `json:"audio"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
	VideoCodec string `json:"video_codec"`
	AudioCodec string `json:"audio_codec"`
}

// SourceMedia is the probe record retained for accepted source bytes. It is
// only populated after the object has been probed successfully.
type SourceMedia struct {
	MediaFacts
}

type Source struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Digest    string      `json:"digest"`
	Bytes     int64       `json:"bytes"`
	Media     SourceMedia `json:"media"`
	CreatedAt time.Time   `json:"created_at"`
}
type Span struct {
	SourceID string  `json:"source_id"`
	InMS     int64   `json:"in_ms"`
	OutMS    int64   `json:"out_ms"`
	GainDB   float64 `json:"gain_db"`
}
type Cue struct {
	SourceID string `json:"source_id"`
	StartMS  int64  `json:"start_ms"`
	EndMS    int64  `json:"end_ms"`
	Text     string `json:"text"`
}
type Sheet struct {
	ID        string    `json:"id"`
	Revision  int       `json:"revision"`
	Name      string    `json:"name"`
	Spans     []Span    `json:"spans"`
	Cues      []Cue     `json:"cues"`
	CreatedAt time.Time `json:"created_at"`
}
type Recipe struct {
	ID         string  `json:"id"`
	Revision   int     `json:"revision"`
	Name       string  `json:"name"`
	Format     string  `json:"format"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	FPS        int     `json:"fps"`
	SampleRate int     `json:"sample_rate"`
	Normalize  bool    `json:"normalize"`
	TargetLUFS float64 `json:"target_lufs"`
}
type Placement struct {
	SourceID    string  `json:"source_id"`
	SourceInMS  int64   `json:"source_in_ms"`
	SourceOutMS int64   `json:"source_out_ms"`
	OutputInMS  int64   `json:"output_in_ms"`
	OutputOutMS int64   `json:"output_out_ms"`
	GainDB      float64 `json:"gain_db"`
}
type Plan struct {
	DurationMS int64       `json:"duration_ms"`
	Placements []Placement `json:"placements"`
	Recipe     Recipe      `json:"recipe"`
	VTT        string      `json:"vtt"`
}

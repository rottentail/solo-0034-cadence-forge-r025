package media

import (
	"cadenceforge/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

type Prober struct{ Binary string }

// Probe measures observable stream facts. Contract judgments (required
// streams, dimension/rate targets, duration tolerance) belong to the caller
// and are deliberately not part of this result.
func (p Prober) Probe(ctx context.Context, path string) (domain.MediaFacts, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	data, err := Run(ctx, p.Binary, []string{"-v", "error", "-protocol_whitelist", "file,pipe", "-format_whitelist", Formats, "-show_format", "-show_streams", "-of", "json", path}, 1<<20)
	if err != nil {
		return domain.MediaFacts{}, err
	}
	var raw struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
		Streams []struct {
			Type        string `json:"codec_type"`
			Codec       string `json:"codec_name"`
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			Rate        string `json:"sample_rate"`
			Channels    int    `json:"channels"`
			Duration    string `json:"duration"`
			Disposition struct {
				Attached int `json:"attached_pic"`
			} `json:"disposition"`
		} `json:"streams"`
	}
	if err = json.Unmarshal(data, &raw); err != nil {
		return domain.MediaFacts{}, err
	}
	duration, _ := strconv.ParseFloat(raw.Format.Duration, 64)
	m := domain.MediaFacts{}
	for _, s := range raw.Streams {
		if s.Type == "video" && s.Disposition.Attached == 0 && !m.Video {
			m.Video = true
			m.Width = s.Width
			m.Height = s.Height
			m.VideoCodec = s.Codec
		}
		if s.Type == "audio" && !m.Audio {
			m.Audio = true
			m.SampleRate, _ = strconv.Atoi(s.Rate)
			m.Channels = s.Channels
			m.AudioCodec = s.Codec
		}
		if duration == 0 {
			d, _ := strconv.ParseFloat(s.Duration, 64)
			if d > duration {
				duration = d
			}
		}
	}
	if math.IsNaN(duration) || math.IsInf(duration, 0) || duration <= 0 || duration > 121 {
		return m, fmt.Errorf("media duration must be >0 and <=121 seconds")
	}
	if !m.Video && !m.Audio {
		return m, fmt.Errorf("no usable audio or video stream")
	}
	if m.Video && (m.Width < 1 || m.Height < 1 || m.Width > 3840 || m.Height > 2160) {
		return m, fmt.Errorf("video exceeds processing dimensions")
	}
	if m.Audio && (m.Channels < 1 || m.Channels > 8 || m.SampleRate < 8000 || m.SampleRate > 192000) {
		return m, fmt.Errorf("unsupported audio layout")
	}
	m.DurationMS = int64(math.Round(duration * 1000))
	return m, nil
}

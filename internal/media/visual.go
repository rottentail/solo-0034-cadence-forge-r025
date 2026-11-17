package media

import (
	"cadenceforge/internal/domain"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Visual creates a bounded inspection image, rather than changing the source.
func Visual(ctx context.Context, binary, path, mode string, m domain.SourceMedia) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	dir, err := os.MkdirTemp("", "cadence-visual-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	output := filepath.Join(dir, "inspection.png")
	args := []string{"-nostdin", "-hide_banner", "-loglevel", "error", "-protocol_whitelist", "file,pipe", "-format_whitelist", Formats, "-i", path}
	switch mode {
	case "storyboard":
		if !m.Video {
			return nil, domain.Invalid("storyboard requires video")
		}
		interval := float64(m.DurationMS) / 4000
		args = append(args, "-an", "-vf", fmt.Sprintf("fps=1/%.6f,scale=160:90:force_original_aspect_ratio=decrease,pad=160:90:(ow-iw)/2:(oh-ih)/2,tile=2x2", interval))
	case "waveform":
		if !m.Audio {
			return nil, domain.Invalid("waveform requires audio")
		}
		args = append(args, "-filter_complex", "[0:a:0]aformat=channel_layouts=mono,showwavespic=s=640x120:colors=0x397d8d[v]", "-map", "[v]")
	default:
		return nil, domain.Invalid("unknown inspection image mode")
	}
	args = append(args, "-frames:v", "1", "-threads", "1", output)
	if _, err = Run(ctx, binary, args, 65536); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(output)
	if err != nil {
		return nil, err
	}
	if len(data) > 2<<20 {
		return nil, fmt.Errorf("inspection image exceeded bound")
	}
	return data, nil
}

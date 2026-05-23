package tui

import (
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

type Progress struct {
	bar *progressbar.ProgressBar
}

func NewProgress() *Progress { return &Progress{} }

func (p *Progress) Start(label string, total int64) {
	p.bar = progressbar.NewOptions64(total,
		progressbar.OptionSetDescription(label),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionThrottle(100_000_000), // 100ms refresh
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { _, _ = os.Stderr.WriteString("\n") }),
		progressbar.OptionFullWidth(),
	)
}

func (p *Progress) Wrap(r io.Reader) io.Reader {
	if p.bar == nil {
		return r
	}
	return &tolerantReader{r: r, bar: p.bar}
}

func (p *Progress) WrapWriter(w io.Writer) io.Writer {
	if p.bar == nil {
		return w
	}
	return &tolerantWriter{w: w, bar: p.bar}
}

func (p *Progress) Finish() {
	if p.bar != nil {
		_ = p.bar.Finish()
	}
}

// tolerantWriter writes to the underlying stream and best-effort updates the
// progress bar. The bar's max is the manifest total_bytes, but the actual byte
// stream includes tar headers/padding, so we ignore overflow errors from the bar.
type tolerantWriter struct {
	w   io.Writer
	bar *progressbar.ProgressBar
}

func (t *tolerantWriter) Write(p []byte) (int, error) {
	n, err := t.w.Write(p)
	if n > 0 {
		_ = t.bar.Add(n)
	}
	return n, err
}

type tolerantReader struct {
	r   io.Reader
	bar *progressbar.ProgressBar
}

func (t *tolerantReader) Read(p []byte) (int, error) {
	n, err := t.r.Read(p)
	if n > 0 {
		_ = t.bar.Add(n)
	}
	return n, err
}

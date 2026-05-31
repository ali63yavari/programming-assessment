package structquery

import "time"

type WaitMode string

const (
	WaitLoad          WaitMode = "load"
	WaitDOMReady      WaitMode = "dom_ready"
	WaitStableDOM     WaitMode = "stable_dom"
	WaitFixedDuration WaitMode = "fixed_duration"
	WaitSelector      WaitMode = "selector"
)

type RenderOptions struct {
	Timeout       time.Duration
	WaitMode      WaitMode
	WaitSelector  string
	WaitDuration  time.Duration
	StableFor     time.Duration
	MaxStableWait time.Duration
	Headless      *bool
}

func DefaultRenderOptions() RenderOptions {
	b := true
	return RenderOptions{
		Timeout:       45 * time.Second,
		WaitMode:      WaitStableDOM,
		StableFor:     800 * time.Millisecond,
		MaxStableWait: 8 * time.Second,
		Headless:      &b,
	}
}

type RenderOptionsProvider interface {
	RenderOptions() RenderOptions
}

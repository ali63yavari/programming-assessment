package structquery

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type RodRenderer struct{}

func (r RodRenderer) Render(
	ctx context.Context,
	url string,
	opts RenderOptions,
) (string, error) {
	opts = normalizeRenderOptions(opts)

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	l := launcher.New()
	if opts.Headless != nil {
		l = l.Headless(*opts.Headless)
	}
	browserURL, err := l.Launch()

	if err != nil {
		return "", fmt.Errorf("launch browser: %w", err)
	}

	browser := rod.New().
		ControlURL(browserURL).
		Context(ctx)

	if err := browser.Connect(); err != nil {
		return "", fmt.Errorf("connect browser: %w", err)
	}

	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return "", fmt.Errorf("create page: %w", err)
	}

	if err := page.Navigate(url); err != nil {
		return "", fmt.Errorf("navigate to %s: %w", url, err)
	}

	if err := waitForPage(page, opts); err != nil {
		return "", err
	}

	html, err := page.HTML()
	if err != nil {
		return "", fmt.Errorf("extract rendered html: %w", err)
	}

	return html, nil
}

func normalizeRenderOptions(opts RenderOptions) RenderOptions {
	defaults := DefaultRenderOptions()

	if opts.Timeout == 0 {
		opts.Timeout = defaults.Timeout
	}

	if opts.WaitMode == "" {
		opts.WaitMode = defaults.WaitMode
	}

	if opts.StableFor == 0 {
		opts.StableFor = defaults.StableFor
	}

	if opts.MaxStableWait == 0 {
		opts.MaxStableWait = defaults.MaxStableWait
	}

	return opts
}

func waitForPage(page *rod.Page, opts RenderOptions) error {
	switch opts.WaitMode {
	case WaitLoad:
		return page.WaitLoad()

	case WaitDOMReady:
		return waitForDOMReady(page)

	case WaitFixedDuration:
		if err := page.WaitLoad(); err != nil {
			return fmt.Errorf("wait load: %w", err)
		}

		d := opts.WaitDuration
		if d == 0 {
			d = 2 * time.Second
		}

		time.Sleep(d)
		return nil

	case WaitSelector:
		if opts.WaitSelector == "" {
			return fmt.Errorf("wait selector is required when wait mode is selector")
		}

		if err := page.WaitLoad(); err != nil {
			return fmt.Errorf("wait load: %w", err)
		}

		_, err := page.Element(opts.WaitSelector)
		if err != nil {
			return fmt.Errorf("wait selector %q: %w", opts.WaitSelector, err)
		}

		return nil

	case WaitStableDOM:
		if err := page.WaitLoad(); err != nil {
			return fmt.Errorf("wait load: %w", err)
		}

		return waitForStableDOM(page, opts.StableFor, opts.MaxStableWait)

	default:
		return fmt.Errorf("unsupported wait mode: %s", opts.WaitMode)
	}
}

func waitForDOMReady(page *rod.Page) error {
	eOpt := rod.EvalOptions{
		JS: `() => document.readyState === "interactive" || document.
readyState === "complete"`,
	}
	err := page.Wait(&eOpt)
	if err != nil {
		return fmt.Errorf("wait dom ready: %w", err)
	}

	return nil
}

func waitForStableDOM(
	page *rod.Page,
	stableFor time.Duration,
	maxWait time.Duration,
) error {
	start := time.Now()

	var lastLength int
	var stableSince time.Time

	for {
		if time.Since(start) > maxWait {
			return nil
		}

		var currentLength int

		_, err := page.Eval(`() => document.documentElement.outerHTML.length`)
		if err != nil {
			return fmt.Errorf("read dom length: %w", err)
		}

		if currentLength == lastLength {
			if stableSince.IsZero() {
				stableSince = time.Now()
			}

			if time.Since(stableSince) >= stableFor {
				return nil
			}
		} else {
			lastLength = currentLength
			stableSince = time.Time{}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

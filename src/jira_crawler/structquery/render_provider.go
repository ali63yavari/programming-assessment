package structquery

import "reflect"

func ResolveRenderOptions(out any) RenderOptions {
	opts := DefaultRenderOptions()

	if provider, ok := out.(RenderOptionsProvider); ok {
		return mergeRenderOptions(opts, provider.RenderOptions())
	}

	value := reflect.ValueOf(out)
	if value.Kind() == reflect.Pointer && !value.IsNil() {
		elem := value.Elem()

		if elem.CanInterface() {
			if provider, ok := elem.Interface().(RenderOptionsProvider); ok {
				return mergeRenderOptions(opts, provider.RenderOptions())
			}
		}
	}

	return opts
}

func mergeRenderOptions(defaults RenderOptions, custom RenderOptions) RenderOptions {
	if custom.Timeout == 0 {
		custom.Timeout = defaults.Timeout
	}

	if custom.WaitMode == "" {
		custom.WaitMode = defaults.WaitMode
	}

	if custom.StableFor == 0 {
		custom.StableFor = defaults.StableFor
	}

	if custom.MaxStableWait == 0 {
		custom.MaxStableWait = defaults.MaxStableWait
	}

	if custom.WaitDuration == 0 {
		custom.WaitDuration = defaults.WaitDuration
	}

	if custom.WaitSelector == "" {
		custom.WaitSelector = defaults.WaitSelector
	}

	if custom.Headless == nil {
		custom.Headless = defaults.Headless
	}

	return custom
}

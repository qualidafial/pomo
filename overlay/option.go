package overlay

type Option func(*overlayOptions)

type overlayOptions struct {
	maxWidth  int
	maxHeight int
}

func WithMaxSize(width, height int) Option {
	return func(o *overlayOptions) {
		o.maxWidth = width
		o.maxHeight = height
	}
}

func defaultOptions() overlayOptions {
	return overlayOptions{}
}

func applyOptions(opts []Option) overlayOptions {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

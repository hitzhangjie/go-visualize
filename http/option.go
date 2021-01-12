package http

type options struct {
	serialization Serialization
	headers       map[string]string
}

// Option http请求选项
type Option func(*options)

// WithJSON 指定JSON序列化
func WithJSON() Option {
	return func(opts *options) {
		opts.serialization = JSON
	}
}

// WithFORM 指定FORM序列化
func WithFORM() Option {
	return func(opts *options) {
		opts.serialization = FORM
	}
}

// WithHeader 指定请求头
func WithHeader(header, value string) Option {
	return func(opts *options) {
		if len(header) == 0 || len(value) == 0 {
			return
		}
		if opts.headers == nil {
			opts.headers = map[string]string{}
		}
		opts.headers[header] = value
	}
}

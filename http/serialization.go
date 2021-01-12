package http

// Serialization 序列化方式
type Serialization string

const (
	// JSON JSON序列化
	JSON = Serialization("JSON")

	// FORM FORM序列化
	FORM = Serialization("FORM")
)

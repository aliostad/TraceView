package tracing

type Config struct {
	TimestampFieldNames     []string
	MessageFieldNames       []string
	LevelFieldNames         []string
	CorrelationIdFieldNames []string
	IndexableFieldNames     []string
	KeepOriginalPayload     bool
}

func EmptyConfig() *Config {
	return &Config{}
}

package section

type (
	Processor struct {
		WebServer ProcessorWebServer `split_words:"true"`
	}

	ProcessorWebServer struct {
		ListenPort uint32 `split_words:"true" default:"8080"`
	}
)

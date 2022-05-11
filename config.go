package kid

type Config struct {
	// ErrorHandler is executed when an error is returned from kid.HandlerFunc.
	//
	// Default: DefaultErrorHandler
	ErrorHandler ErrorHandlerFunc
}

func setDefaultConfig(k *Kid) {
	if k.config.ErrorHandler == nil {
		k.config.ErrorHandler = DefaultErrorHandler
	}
}

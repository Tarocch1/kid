package kid

type Config struct {
	// ErrorHandler is executed when an error is returned from kid.HandlerFunc.
	//
	// Default: DefaultErrorHandler
	ErrorHandler ErrorHandlerFunc
}

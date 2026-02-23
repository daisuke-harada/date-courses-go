package logger

var debug bool

func setDebug(enable bool) {
	debug = enable
}

type Config struct {
	Level            string
	Development      bool
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
}

func development() Config {
	return Config{
		Level:            "debug",
		Development:      true,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func production() Config {
	return Config{
		Level:            "info",
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newConfig() Config {
	if debug {
		return development()
	}
	return production()
}

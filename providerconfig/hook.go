package internal

type SignalHookName string

const (
	TraceHook  SignalHookName = "trace"
	MetricHook SignalHookName = "metric"
	logHook    SignalHookName = "log"
)

type ShutdownHooks map[SignalHookName]func()

func NewShutdownHooks(fns ...func() (SignalHookName, func())) ShutdownHooks {
	sdh := make(ShutdownHooks, len(fns))
	for _, fn := range fns {
		name, hook := fn()
		sdh[name] = hook
	}
	return sdh
}

func ShutDownPair(name SignalHookName, fn func()) func() (SignalHookName, func()) {
	return func() (SignalHookName, func()) {
		return name, fn
	}
}

func (h ShutdownHooks) ShutdownAll() {
	for _, hook := range h {
		hook()
	}
}

func (h ShutdownHooks) ShutdownByType(hookType SignalHookName) bool {
	if hook, exists := h[hookType]; exists {
		hook()
		return true
	}
	return false
}

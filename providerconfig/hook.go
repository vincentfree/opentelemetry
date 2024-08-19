package providerconfig

type SignalHookName string

const (
	TraceHook  SignalHookName = "trace"
	MetricHook SignalHookName = "metric"
	logHook    SignalHookName = "log"
)

type ShutdownHook func()
type ShutdownHooks map[SignalHookName]ShutdownHook

func NewShutdownHooks(fns ...func() (SignalHookName, ShutdownHook)) ShutdownHooks {
	sdh := make(ShutdownHooks, len(fns))
	for _, fn := range fns {
		name, hook := fn()
		sdh[name] = hook
	}
	return sdh
}

func ShutDownPair(name SignalHookName, fn ShutdownHook) func() (SignalHookName, ShutdownHook) {
	return func() (SignalHookName, ShutdownHook) {
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

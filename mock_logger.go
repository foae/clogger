package clogger

var (
	noopLogger = &mockLogger{
		streamEventFn:        func(ev *Event) {},
		setEventOptionsFn:    func(opts ...EventOption) {},
		eventOptionsFn:       func() []EventOption { return make([]EventOption, 0) },
		streamLogEntryFn:     func(e *LogEntry) {},
		setLogEntryOptionsFn: func(opts ...LogEntryOption) {},
		logEntryOptionsFn:    func() []LogEntryOption { return make([]LogEntryOption, 0) },
		setEncoderFn:         func(enc Encoder) {},
	}
)

type mockLogger struct {
	streamLogEntryFn     func(e *LogEntry)
	setLogEntryOptionsFn func(opts ...LogEntryOption)
	logEntryOptionsFn    func() []LogEntryOption

	streamEventFn     func(ev *Event)
	setEventOptionsFn func(opts ...EventOption)
	eventOptionsFn    func() []EventOption

	setEncoderFn func(enc Encoder)
}

func (ml *mockLogger) StreamLogEntry(e *LogEntry) {
	ml.streamLogEntryFn(e)
}
func (ml *mockLogger) SetLogEntryOptions(opts ...LogEntryOption) {
	ml.setLogEntryOptionsFn(opts...)
}
func (ml *mockLogger) LogEntryOptions() []LogEntryOption {
	return ml.logEntryOptionsFn()
}
func (ml *mockLogger) StreamEvent(ev *Event) {
	ml.streamEventFn(ev)
}
func (ml *mockLogger) SetEventOptions(opts ...EventOption) {
	ml.setEventOptionsFn(opts...)
}
func (ml *mockLogger) EventOptions() []EventOption {
	return ml.eventOptionsFn()
}
func (ml *mockLogger) SetEncoder(enc Encoder) {
	ml.setEncoderFn(enc)
}

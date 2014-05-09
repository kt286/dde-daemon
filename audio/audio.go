package main

import "dlib"
import "dlib/dbus"
import "dlib/logger"
import "dlib/pulse"
import "os"

var Logger = logger.NewLogger("com.deepin.daemon.Audio")

type Audio struct {
	init bool
	core *pulse.Context

	Sinks         []*Sink
	Sources       []*Source
	SinkInputs    []*SinkInput
	DefaultSink   string
	DefaultSource string
}

func NewSink(core *pulse.Sink) *Sink {
	s := &Sink{core: core}
	s.update()
	return s
}
func NewSource(core *pulse.Source) *Source {
	s := &Source{core: core}
	s.update()
	return s
}
func NewSinkInput(core *pulse.SinkInput) *SinkInput {
	s := &SinkInput{core: core}
	s.update()
	return s
}
func NewAudio(core *pulse.Context) *Audio {
	a := &Audio{core: core}
	a.update()
	a.initEventHandlers()
	return a
}

func (a *Audio) SetDefaultSink(name string) {
	a.core.SetDefaultSink(name)
}
func (a *Audio) SetDefaultSource(name string) {
	a.core.SetDefaultSource(name)
}

type SourceOutputTest struct {
	Volume float64
}

func (*SourceOutputTest) Tick() {
}

type Sink struct {
	core        *pulse.Sink
	Name        string
	Description string

	BaseVolume float64
	Volume     float64
	Balance    float64
	Mute       bool

	Ports      []string
	ActivePort string
}

func (s *Sink) SetVolume(v float64) {
	s.core.SetAvgVolume(v)
}
func (s *Sink) SetBalance(v float64) {
	s.core.SetBalance(v)
}
func (s *Sink) SetMute(v bool) {
	s.core.SetMute(v)
}
func (s *Sink) SetPort(name string) {
	s.core.SetPort(name)
}

type SinkInput struct {
	core   *pulse.SinkInput
	Name   string
	Icon   string
	Volume float64
	Mute   bool
}

func (s *SinkInput) SetVolume(v float64) {
	s.core.SetAvgVolume(v)
}
func (s *SinkInput) SetMute(v bool) {
	s.core.SetMute(v)
}

type Source struct {
	core        *pulse.Source
	Name        string
	Description string

	BaseVolume float64
	Volume     float64
	Balance    float64
	Mute       bool

	ActivePort string
	Ports      []string
}

func (s *Source) SetPort(name string) {
	s.core.SetPort(name)
}

func main() {
	defer Logger.EndTracing()
	if !dlib.UniqueOnSession("com.deepin.daemon.Audio") {
		Logger.Warning("There already has an Audio daemon running.")
		return
	}

	ctx := pulse.GetContext()
	audio := NewAudio(ctx)
	tester := &SourceOutputTest{}

	dbus.InstallOnSession(tester)

	if err := dbus.InstallOnSession(audio); err != nil {
		Logger.Error("Failed InstallOnSession:", err)
		return
	}
	if err := dbus.InstallOnSession(tester); err != nil {
		Logger.Error("Failed InstallOnSession:", err)
	}

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		Logger.Error("dbus.Wait recieve an error:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
	dbus.Wait()
}

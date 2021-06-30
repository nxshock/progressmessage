package progressmessage

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var (
	defaultWriter           = os.Stderr
	defaultUpdatePeriod     = time.Second
	defaultMessageDelimiter = "\r"
)

type ProgressMessage struct {
	// message format
	format string

	// list of message params
	params []interface{}

	// ticker for updating message
	ticker *time.Ticker

	// chan for stopping updates
	stopChan chan struct{}

	// previous message to prevent output of same message more than one time
	prevMessage string

	writer    io.Writer
	isStarted bool

	// message update period
	updatePeriod time.Duration

	mu sync.RWMutex
}

func New(format string) *ProgressMessage {
	pm := &ProgressMessage{
		format: format,

		stopChan:     make(chan struct{}),
		writer:       defaultWriter,
		updatePeriod: defaultUpdatePeriod}

	return pm
}

func (pm *ProgressMessage) ChangeWriter(newWriter io.Writer) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.writer = newWriter
}

func (pm *ProgressMessage) ChangeFormat(newFormat string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.format = newFormat
}

func (pm *ProgressMessage) Update(params ...interface{}) {
	// TODO: use mutex?

	pm.params = params
}

func (pm *ProgressMessage) ChangeUpdatePeriod(newUpdatePeriod time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.isStarted {
		defer func() {
			pm.Stop()
			pm.Start()
		}()
	}

	pm.updatePeriod = newUpdatePeriod
}

func (pm *ProgressMessage) Start() {
	pm.isStarted = true

	go func() {
		pm.ticker = time.NewTicker(pm.updatePeriod)

		for {
			select {
			case <-pm.ticker.C:
				if pm.params == nil {
					continue
				}

				m := fmt.Sprintf(pm.format, pm.params...)
				if m == pm.prevMessage {
					continue
				}
				pm.writer.Write([]byte(defaultMessageDelimiter + m))
				pm.prevMessage = m
			case <-pm.stopChan:
				pm.ticker.Stop()
				break
			}
		}
	}()
}

func (pm *ProgressMessage) Stop() {
	pm.stopChan <- struct{}{}

	pm.isStarted = false
}

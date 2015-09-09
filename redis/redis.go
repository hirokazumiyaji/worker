package redis

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/garyburd/redigo/redis"
)

type HandlerFunc func([]byte) error

type Worker struct {
	wg         *sync.WaitGroup
	connection redis.Conn
	queueName  string
	handler    HandlerFunc
	signal     chan os.Signal
}

func NewWorker(wg *sync.WaitGroup, connection redis.Conn, queueName string, handler HandlerFunc) *Worker {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	return &Worker{
		wg:         wg,
		connection: connection,
		queueName:  queueName,
		handler:    handler,
		signal:     sig,
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go w.backgroundLoop()
}

func (w *Worker) Shutdown() {
	defer w.wg.Done()
	w.connection.Close()
}

func (w *Worker) backgroundLoop() {
	errorCount := 0
loop:
	for {
		select {
		case _ = <-w.signal:
			signal.Stop(w.signal)
			break loop
		default:
			var data interface{}
			data, err := w.connection.Do("BRPOP", w.queueName, 1)
			if err != nil {
				errorCount += 1
				continue
			}

			errorCount = 0
			if data == nil {
				continue
			}

			values, err := redis.Values(data, err)
			if err != nil {
				errorCount += 1
				continue
			}

			value := values[1].([]byte)
			if err := w.handler(values[1].([]byte)); err != nil {
				w.connection.Do("LPUSH", fmt.Sprintf("%s:FAIL", w.queueName), value)
			}
		}
	}

	w.Shutdown()
}

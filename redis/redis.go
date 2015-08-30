package redis

import (
	"github.com/garyburd/redigo/redis"
)

type RedisConfig struct {
	Host     string
	Port     int
	DB       int
	Password string
	Socker   string
}

type Worker struct {
	conn      *redis.Conn
	queueName string
}

func NewWorker(config *RedisConfig, queueName string) *Worker {
	connection := nil
	return &Worker{
		conn:      connection,
		queueName: queueName,
	}
}

func (w *Worker) Start() error {
}

func (w *Worker) Shutdown() error {
	if err := w.conn.Close(); err != nil {
		return err
	}
}

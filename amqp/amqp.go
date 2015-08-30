package amqp

import (
	"github.com/hirokazumiyaji/worker/worker"
	"github.com/streadway/amqp"
)

type Worker struct {
	worker.Worker
	queueName  string
	brokerURL  string
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewWorker(brokerURL, queueName string) *Worker {
	return &Worker{
		queueName: queueName,
		brokerURL: brokerURL,
	}
}

func (w *Worker) Start() error {
	if w.connection, err = amqp.Dial(w.brokerURL); err != nil {
		return err
	}

	if w.channel, err = w.connection.Channel(); err != nil {
		return err
	}

	if w.queue, err = w.channel.QueueDeclare(
		w.queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	return nil
}

func (w *Worker) Shutdown() error {
	if err := w.channel.Cancel(); err != nil {
		return err
	}
	if err := w.connection.Close(); err != nil {
		return err
	}
	return nil
}

package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hirokazumiyaji/woker/woker"
)

type Worker struct {
	worker.Worker
	conn *sqs.SQS
}

func NewWorker(config *aws.Config) *Worker {
	return &Worker{
		conn: sqs.New(config),
	}
}

func (w *Worker) Start() error {
	return nil
}

func (w *Worker) Shutdown() error {
}

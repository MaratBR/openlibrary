package httplim

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type httpRequestEnvelope struct {
	responseQueue chan httpResponseEnvelope
	req           *http.Request
}

type httpResponseEnvelope struct {
	err  error
	resp *http.Response
}

const (
	stateIdle = iota
	stateWorking
	stateClosed
)

type HttpWorker struct {
	workers             int
	client              *http.Client
	Limiter             *rate.Limiter
	requestsQueue       chan httpRequestEnvelope
	responseChannelPool sync.Pool
	wg                  sync.WaitGroup
	state               int
}

func NewHttpWorker(workers int, client *http.Client, limiter *rate.Limiter) *HttpWorker {
	return &HttpWorker{
		workers:             workers,
		client:              client,
		Limiter:             limiter,
		requestsQueue:       make(chan httpRequestEnvelope, workers),
		responseChannelPool: sync.Pool{New: func() any { return make(chan httpResponseEnvelope, 1) }},
	}
}

func (w *HttpWorker) Do(req *http.Request) (*http.Response, error) {
	responseQueue := w.responseChannelPool.Get().(chan httpResponseEnvelope)
	defer w.responseChannelPool.Put(responseQueue)
	envelope := httpRequestEnvelope{
		responseQueue: responseQueue,
		req:           req,
	}
	w.requestsQueue <- envelope
	response := <-responseQueue
	return response.resp, response.err
}

func (w *HttpWorker) Close() {
	if w.state == stateClosed || w.state == stateIdle {
		return
	}
	close(w.requestsQueue)
	w.wg.Wait()
	w.state = stateClosed
}

func (w *HttpWorker) worker(id int) {
	defer w.wg.Done()
	ctx := context.Background()

	slog.Info("starting worker", "worker_id", id)
	for envelope := range w.requestsQueue {
		err := w.Limiter.Wait(ctx)
		if err != nil {
			slog.Error("error while waiting for rate limiter", "err", err, "worker_id", id)
			break
		}

		resp, err := w.client.Do(envelope.req)
		if resp != nil && resp.StatusCode >= 400 {
			slog.Error("got non-200 response", "worker_id", id, "status", resp.StatusCode)
		}

		envelope.responseQueue <- httpResponseEnvelope{
			err:  err,
			resp: resp,
		}
	}

	slog.Info("worker exited", "worker_id", id)
}

func (w *HttpWorker) Run() {
	if w.state == stateWorking {
		return
	}
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go w.worker(i)
	}
	w.state = stateWorking
}

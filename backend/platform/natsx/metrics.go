package natsx

import (
	"errors"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	handled   *prometheus.HistogramVec
	unhandled *prometheus.CounterVec
	errors    *prometheus.CounterVec
	published *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		handled: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "natsx_consumer_handle_duration_seconds",
			Help:    "Duration of handled NATS messages",
			Buckets: prometheus.ExponentialBuckets(1, 2, 10),
		}, []string{"stream", "consumer", "outcome"}),
		unhandled: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "natsx_consumer_messages_unhandled_total",
			Help: "Number of unhandled NATS messages",
		}, []string{"stream", "consumer", "reason"}),
		errors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "natsx_consumer_errors_total",
			Help: "Number of NATS errors",
		}, []string{"stream", "consumer", "op"}),
		published: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "natsx_publisher_publish_duration_seconds",
			Help: "Duration of NATS message publish",
		}, []string{"stream", "subject", "outcome"}),
	}
	reg.MustRegister(
		m.handled,
		m.unhandled,
		m.errors,
		m.published,
	)
	return m
}

type consumerMetrics struct {
	handled   prometheus.ObserverVec
	unhandled *prometheus.CounterVec
	errors    *prometheus.CounterVec
}

func (m *Metrics) forConsumer(cfg *ConsumerConfig) consumerMetrics {
	labels := prometheus.Labels{"stream": cfg.StreamName, "consumer": cfg.DurableName}
	c := consumerMetrics{
		handled:   m.handled.MustCurryWith(labels),
		unhandled: m.unhandled.MustCurryWith(labels),
		errors:    m.errors.MustCurryWith(labels),
	}
	for _, o := range []handlerResult{HandlerSuccess, HandlerRetriable, HandlerFatal} {
		c.handled.WithLabelValues(string(o))
	}
	for _, o := range []unhandledReason{UnhandledMalformed, UnhandledDraining} {
		c.unhandled.WithLabelValues(string(o))
	}
	for _, o := range []errorOp{ErrorOpAck, ErrorOpHeartbeat, ErrorOpConsume} {
		c.errors.WithLabelValues(string(o))
	}

	return c
}

type publisherMetrics struct {
	published prometheus.ObserverVec
}

func (m *Metrics) forPublisher(cfg *PublisherConfig) publisherMetrics {
	labels := prometheus.Labels{"stream": cfg.StreamName, "subject": cfg.Subject}
	c := publisherMetrics{
		published: m.published.MustCurryWith(labels),
	}
	for _, o := range []handlerResult{HandlerSuccess, HandlerRetriable, HandlerFatal} {
		c.published.WithLabelValues(string(o))
	}

	return c
}

type handlerResult string

const (
	HandlerSuccess   handlerResult = "success"
	HandlerRetriable handlerResult = "retriable"
	HandlerFatal     handlerResult = "fatal"
)

func handlerOutcome(err error) handlerResult {
	if err == nil {
		return HandlerSuccess
	}
	if errors.Is(err, queue.ErrFatal) {
		return HandlerFatal
	}
	return HandlerRetriable
}

type unhandledReason string

const (
	UnhandledMalformed unhandledReason = "malformed"
	UnhandledDraining  unhandledReason = "draining"
)

type errorOp string

const (
	ErrorOpAck       errorOp = "ack"
	ErrorOpHeartbeat errorOp = "heartbeat"
	ErrorOpConsume   errorOp = "consume"
)

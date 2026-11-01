package natsx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/RobertsMJ/simc-cloud/backend/transport/queue"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/suite"
)

type ConsumerTestSuite struct {
	suite.Suite

	handler  *recordingHandler
	shutdown *fakeShutdowner
	consumer *Consumer[testPayload]
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (s *ConsumerTestSuite) SetupTest() {
	s.handler = &recordingHandler{}
	s.shutdown = &fakeShutdowner{called: make(chan struct{}, 1)}
	s.consumer = NewConsumer[testPayload](
		slog.Default(),
		&ConsumerConfig{StreamName: "stream", DurableName: "consumer"},
		nil,
		s.handler,
		s.shutdown.Shutdown,
	)
	s.consumer.heartbeatInterval = time.Minute
}

type testPayload struct {
	Name string `json:"name"`
}

type recordingHandler struct {
	mu       sync.Mutex
	received []testPayload
	err      error
}

func (h *recordingHandler) Handle(_ context.Context, msg testPayload) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.received = append(h.received, msg)
	return h.err
}

func (h *recordingHandler) payloads() []testPayload {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]testPayload(nil), h.received...)
}

type fakeMsg struct {
	jetstream.Msg
	data []byte

	mu         sync.Mutex
	calls      []string
	inProgress int
}

func (m *fakeMsg) Data() []byte { return m.data }

func (m *fakeMsg) record(call string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, call)
}

func (m *fakeMsg) Ack() error { m.record("ack"); return nil }
func (m *fakeMsg) Nak() error { m.record("nak"); return nil }
func (m *fakeMsg) TermWithReason(string) error {
	m.record("term")
	return nil
}

func (m *fakeMsg) InProgress() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inProgress++
	return nil
}

func (m *fakeMsg) progressCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.inProgress
}

func (m *fakeMsg) recorded() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.calls...)
}

func (s *ConsumerTestSuite) TestOnMsg() {
	tests := []struct {
		name        string
		data        string
		handlerErr  error
		draining    bool
		wantCalls   []string
		wantPayload []testPayload
	}{
		{
			name:        "valid message is acked",
			data:        `{"name":"alpha"}`,
			wantCalls:   []string{"ack"},
			wantPayload: []testPayload{{Name: "alpha"}},
		},
		{
			name:      "invalid json is terminated without handling",
			data:      `{not json`,
			wantCalls: []string{"term"},
		},
		{
			name:        "fatal handler error is terminated",
			data:        `{"name":"alpha"}`,
			handlerErr:  fmt.Errorf("boom: %w", queue.ErrFatal),
			wantCalls:   []string{"term"},
			wantPayload: []testPayload{{Name: "alpha"}},
		},
		{
			name:        "other handler error is nacked",
			data:        `{"name":"alpha"}`,
			handlerErr:  errors.New("transient"),
			wantCalls:   []string{"nak"},
			wantPayload: []testPayload{{Name: "alpha"}},
		},
		{
			name:      "draining consumer nacks without handling",
			data:      `{"name":"alpha"}`,
			draining:  true,
			wantCalls: []string{"nak"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()
			s.handler.err = tt.handlerErr
			s.consumer.draining.Store(tt.draining)
			msg := &fakeMsg{data: []byte(tt.data)}

			s.consumer.onMsg(context.Background(), msg)

			s.Equal(tt.wantCalls, msg.recorded())
			s.Equal(tt.wantPayload, s.handler.payloads())
		})
	}
}

type fakeShutdowner struct {
	called chan struct{}
}

func (f *fakeShutdowner) Shutdown() error {
	f.called <- struct{}{}
	return nil
}

type handlerFunc func(ctx context.Context, msg testPayload) error

func (f handlerFunc) Handle(ctx context.Context, msg testPayload) error { return f(ctx, msg) }

type fakeJetStream struct {
	jetstream.JetStream
	consumer jetstream.Consumer
	err      error
}

func (j *fakeJetStream) Consumer(context.Context, string, string) (jetstream.Consumer, error) {
	return j.consumer, j.err
}

type fakePullConsumer struct {
	jetstream.Consumer
	ackWait time.Duration
	sub     *fakeSubscription
	handler jetstream.MessageHandler
}

func (c *fakePullConsumer) CachedInfo() *jetstream.ConsumerInfo {
	return &jetstream.ConsumerInfo{Config: jetstream.ConsumerConfig{AckWait: c.ackWait}}
}

func (c *fakePullConsumer) Consume(h jetstream.MessageHandler, _ ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	c.handler = h
	return c.sub, nil
}

type fakeSubscription struct {
	mu           sync.Mutex
	drained      bool
	closeOnDrain bool
	closed       chan struct{}
}

func newFakeSubscription(closeOnDrain bool) *fakeSubscription {
	return &fakeSubscription{closeOnDrain: closeOnDrain, closed: make(chan struct{})}
}

func (f *fakeSubscription) Stop() {}

func (f *fakeSubscription) Drain() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.drained = true
	if f.closeOnDrain {
		close(f.closed)
	}
}

func (f *fakeSubscription) Closed() <-chan struct{} { return f.closed }

func (f *fakeSubscription) wasDrained() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.drained
}

func (s *ConsumerTestSuite) bind(pc *fakePullConsumer) {
	s.consumer.js = &fakeJetStream{consumer: pc}
}

func (s *ConsumerTestSuite) waitFor(ch <-chan struct{}) {
	s.T().Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		s.FailNow("timed out waiting for signal")
	}
}

func (s *ConsumerTestSuite) TestStart_BindFailure() {
	bindErr := errors.New("no such consumer")
	s.consumer.js = &fakeJetStream{err: bindErr}

	err := s.consumer.Start(context.Background())

	s.ErrorIs(err, bindErr)
	s.ErrorContains(err, "failed to bind consumer")
}

func (s *ConsumerTestSuite) TestStart_DeliversMessagesToHandler() {
	pc := &fakePullConsumer{ackWait: time.Minute, sub: newFakeSubscription(true)}
	s.bind(pc)
	s.Require().NoError(s.consumer.Start(context.Background()))
	msg := &fakeMsg{data: []byte(`{"name":"alpha"}`)}

	pc.handler(msg)

	s.Equal([]testPayload{{Name: "alpha"}}, s.handler.payloads())
	s.Equal([]string{"ack"}, msg.recorded())
}

func (s *ConsumerTestSuite) TestStop_DrainsSubscription() {
	sub := newFakeSubscription(true)
	s.bind(&fakePullConsumer{ackWait: time.Minute, sub: sub})
	s.Require().NoError(s.consumer.Start(context.Background()))

	s.NoError(s.consumer.Stop(context.Background()))

	s.True(sub.wasDrained())
	s.True(s.consumer.draining.Load())
}

func (s *ConsumerTestSuite) TestStop_CancelsInFlightHandlerAtDeadline() {
	pc := &fakePullConsumer{ackWait: time.Minute, sub: newFakeSubscription(false)}
	s.bind(pc)
	entered := make(chan struct{})
	s.consumer.handler = handlerFunc(func(ctx context.Context, _ testPayload) error {
		close(entered)
		<-ctx.Done()
		return ctx.Err()
	})
	s.Require().NoError(s.consumer.Start(context.Background()))
	msg := &fakeMsg{data: []byte(`{"name":"alpha"}`)}
	done := make(chan struct{})
	go func() {
		defer close(done)
		pc.handler(msg)
	}()
	s.waitFor(entered)
	stopCtx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	s.NoError(s.consumer.Stop(stopCtx))

	s.waitFor(done)
	s.Equal([]string{"nak"}, msg.recorded())
}

func (s *ConsumerTestSuite) TestHeartbeat_ExtendsAckWhileHandling() {
	s.consumer.heartbeatInterval = 5 * time.Millisecond
	release := make(chan struct{})
	s.consumer.handler = handlerFunc(func(context.Context, testPayload) error {
		<-release
		return nil
	})
	msg := &fakeMsg{data: []byte(`{"name":"alpha"}`)}
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.consumer.onMsg(context.Background(), msg)
	}()

	s.Eventually(func() bool { return msg.progressCount() > 0 }, time.Second, 5*time.Millisecond)
	close(release)

	s.waitFor(done)
	s.Equal([]string{"ack"}, msg.recorded())
}

func (s *ConsumerTestSuite) TestOnConsumeErr() {
	tests := []struct {
		name         string
		err          error
		wantShutdown bool
	}{
		{name: "consumer deleted shuts down", err: jetstream.ErrConsumerDeleted, wantShutdown: true},
		{name: "wrapped bad request shuts down", err: fmt.Errorf("pull: %w", jetstream.ErrBadRequest), wantShutdown: true},
		{name: "missing heartbeat is survivable", err: jetstream.ErrNoHeartbeat},
		{name: "leadership change is survivable", err: jetstream.ErrConsumerLeadershipChanged},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.SetupTest()

			s.consumer.onConsumeErr(nil, tt.err)

			s.Equal(tt.wantShutdown, len(s.shutdown.called) == 1)
		})
	}
}

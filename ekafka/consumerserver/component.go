package consumerserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/gotomicro/ego-component/ekafka"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server"
	"github.com/segmentio/kafka-go"
)

// Interface check
var _ server.Server = (*Component)(nil)

// PackageName is the name of this component.
const PackageName = "component.ekafka.consumerserver"

type consumptionMode int

const (
	consumptionModeManual consumptionMode = iota + 1
	consumptionModeSingle
)

// Component starts an Ego server for message consuming.
type Component struct {
	ServerCtx             context.Context
	stopServer            context.CancelFunc
	config                *config
	name                  string
	ekafkaComponent       *ekafka.Component
	logger                *elog.Component
	mode                  consumptionMode
	onEachMessageHandler  OnEachMessageHandler
	onStartMessageHandler OnStartHandler
	consumptionErrors     chan<- error
}

// PackageName returns the package name.
func (cmp *Component) PackageName() string {
	return PackageName
}

// Info returns server info, used by governor and consumer balancer.
func (cmp *Component) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithKind(constant.ServiceProvider),
	)
	return &info
}

// GracefulStop stops the server.
func (cmp *Component) GracefulStop(ctx context.Context) error {
	cmp.stopServer()
	return nil
}

// Stop stops the server.
func (cmp *Component) Stop() error {
	cmp.stopServer()
	return nil
}

// Init ...
func (cmp *Component) Init() error {
	return nil
}

// Name returns the name of this instance.
func (cmp *Component) Name() string {
	return cmp.name
}

// Start will start consuming.
func (cmp *Component) Start() error {
	switch cmp.mode {
	case consumptionModeManual:
		return cmp.startManualMode()
	case consumptionModeSingle:
		return cmp.startSingleMode()
	default:
		return fmt.Errorf("undefined consumption mode: %v", cmp.mode)
	}
}

// GetConsumer returns the default consumer.
func (cmp *Component) GetConsumer() *ekafka.Consumer {
	return cmp.ekafkaComponent.Consumer(cmp.config.ConsumerName)
}

// OnEachMessage registers a single message handler.
func (cmp *Component) OnEachMessage(consumptionErrors chan<- error, handler OnEachMessageHandler) error {
	cmp.consumptionErrors = consumptionErrors
	cmp.mode = consumptionModeSingle
	cmp.onEachMessageHandler = handler
	return nil
}

// OnStart registers a manual message handler.
func (cmp *Component) OnStart(handler OnStartHandler) error {
	cmp.mode = consumptionModeManual
	cmp.onStartMessageHandler = handler
	return nil
}

func isErrorUnrecoverable(err error) bool {
	if kafkaError, ok := err.(kafka.Error); ok {
		if kafkaError.Temporary() {
			return false
		}
	}
	return true
}

func (cmp *Component) startManualMode() error {
	consumer := cmp.GetConsumer()

	if cmp.onStartMessageHandler == nil {
		return errors.New("you must define a MessageHandler first")
	}

	handlerExit := make(chan error)
	go func() {
		handlerExit <- cmp.onStartMessageHandler(cmp.ServerCtx, consumer)
		close(handlerExit)
	}()

	var originErr error
	select {
	case originErr := <-handlerExit:
		if originErr != nil {
			cmp.logger.Error("terminating consumer because an error", elog.FieldErr(originErr))
		} else {
			cmp.logger.Info("message handler exited without any error, terminating consumer server")
		}
		cmp.stopServer()
	case <-cmp.ServerCtx.Done():
		originErr := cmp.ServerCtx.Err()
		cmp.logger.Error("terminating consumer because a context error", elog.FieldErr(originErr))

		err := <-handlerExit
		if err != nil {
			cmp.logger.Error("terminating consumer because an error", elog.FieldErr(err))
		} else {
			cmp.logger.Info("message handler exited without any error")
		}
	}

	err := cmp.closeConsumer(consumer)
	if err != nil {
		return fmt.Errorf("encountered an error while closing consumer: %w", err)
	}

	if errors.Is(originErr, context.Canceled) {
		return nil
	}

	return originErr
}

func (cmp *Component) startSingleMode() error {
	consumer := cmp.GetConsumer()

	if cmp.onEachMessageHandler == nil {
		return errors.New("you must define a MessageHandler first")
	}

	unrecoverableError := make(chan error)
	go func() {
		for {
			if cmp.ServerCtx.Err() != nil {
				return
			}

			message, err := consumer.FetchMessage(cmp.ServerCtx)
			if err != nil {
				cmp.consumptionErrors <- err
				cmp.logger.Error("encountered an error while fetching message", elog.FieldErr(err))

				// If this error is unrecoverable, stop consuming.
				if isErrorUnrecoverable(err) {
					unrecoverableError <- err
					return
				}
				// Otherwise, try to fetch message again.
				continue
			}

			err = cmp.onEachMessageHandler(cmp.ServerCtx, message)
			if err != nil {
				cmp.logger.Error("encountered an error while handling message", elog.FieldErr(err))
				cmp.consumptionErrors <- err
				// Any error that returned from the handler should be considered as an unrecoverable
				// error, developers should write their own retry logic in the handler.
				unrecoverableError <- err
				return
			}

		COMMIT:

			err = consumer.CommitMessages(cmp.ServerCtx, message)
			if err != nil {
				cmp.consumptionErrors <- err
				cmp.logger.Error("encountered an error while committing message", elog.FieldErr(err))

				// If this error is unrecoverable, stop retry and consuming.
				if isErrorUnrecoverable(err) {
					unrecoverableError <- err
					return
				}

				if cmp.ServerCtx.Err() != nil {
					return
				}

				// Try to commit this message again.
				cmp.logger.Debug("try to commit message again")
				goto COMMIT
			}
		}
	}()

	select {
	case <-cmp.ServerCtx.Done():
		rootErr := cmp.ServerCtx.Err()
		cmp.logger.Error("terminating consumer because a context error", elog.FieldErr(rootErr))

		err := cmp.closeConsumer(consumer)
		if err != nil {
			return fmt.Errorf("encountered an error while closing consumer: %w", err)
		}

		if errors.Is(rootErr, context.Canceled) {
			return nil
		}

		return rootErr
	case originErr := <-unrecoverableError:
		if originErr == nil {
			panic("unrecoverableError should receive an error instead of nil")
		}

		cmp.logger.Fatal("stopping server because of an unrecoverable error", elog.FieldErr(originErr))
		cmp.Stop()

		err := cmp.closeConsumer(consumer)
		if err != nil {
			return fmt.Errorf("exiting due to an unrecoverable error, but encountered an error while closing consumer: %w", err)
		}
		return originErr
	}
}

func (cmp *Component) closeConsumer(consumer *ekafka.Consumer) error {
	if err := consumer.Close(); err != nil {
		cmp.logger.Fatal("failed to close kafka writer", elog.FieldErr(err))
		return err
	}
	cmp.logger.Info("consumer server terminated")
	return nil
}

// NewConsumerServerComponent creates a new server instance.
func NewConsumerServerComponent(name string, config *config, ekafkaComponent *ekafka.Component, logger *elog.Component) *Component {
	serverCtx, stopServer := context.WithCancel(context.Background())
	return &Component{
		ServerCtx:       serverCtx,
		stopServer:      stopServer,
		name:            name,
		config:          config,
		ekafkaComponent: ekafkaComponent,
		logger:          logger,
		mode:            consumptionModeSingle,
	}
}

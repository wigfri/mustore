package mailqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wigfri/mustore/domain/services"
	"github.com/wigfri/mustore/services/config"
	"github.com/wigfri/mustore/services/mail"
)

const mailQueueName = "mustore.mail"

type MailJob struct {
	Kind string `json:"kind"`
	To   string `json:"to"`
	Code string `json:"code"`
}

const (
	kindVerification = "verification"
	kindLogin        = "login"
)

type Dispatcher struct {
	cfg    *config.Config
	logger services.Logger
	sender *mail.Sender
}

func NewDispatcher(cfg *config.Config, log services.Logger) *Dispatcher {
	return &Dispatcher{
		cfg:    cfg,
		logger: log,
		sender: mail.NewSender(cfg),
	}
}

var _ services.MailQueue = (*Dispatcher)(nil)

func (d *Dispatcher) PublishVerificationEmail(to, code string) error {
	return d.publish(MailJob{Kind: kindVerification, To: to, Code: code})
}

func (d *Dispatcher) PublishLoginCodeEmail(to, code string) error {
	return d.publish(MailJob{Kind: kindLogin, To: to, Code: code})
}

func (d *Dispatcher) StartConsumer() {
	if !d.cfg.RabbitMQEnabled() || d.cfg.MailSkipSend() {
		return
	}
	ctx := context.Background()
	go d.consumeLoop(ctx)
}

func (d *Dispatcher) publish(job MailJob) error {
	if d.cfg.MailSkipSend() {
		d.logger.Info("mail skipped (dev)", "kind", job.Kind, "to", job.To, "code", job.Code)
		return nil
	}
	if d.cfg.RabbitMQEnabled() {
		return d.publishAMQP(job)
	}
	return d.deliver(job)
}

func (d *Dispatcher) publishAMQP(job MailJob) error {
	conn, err := amqp.Dial(d.cfg.RabbitMQURL())
	if err != nil {
		return fmt.Errorf("rabbitmq dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(mailQueueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("rabbitmq declare queue: %w", err)
	}

	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	if err := ch.PublishWithContext(context.Background(), "", mailQueueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	}); err != nil {
		return fmt.Errorf("rabbitmq publish: %w", err)
	}
	return nil
}

func (d *Dispatcher) deliver(job MailJob) error {
	switch job.Kind {
	case kindVerification:
		return d.sender.SendVerification(job.To, job.Code)
	case kindLogin:
		return d.sender.SendLoginCode(job.To, job.Code)
	default:
		return fmt.Errorf("unknown mail job kind: %s", job.Kind)
	}
}

func (d *Dispatcher) consumeLoop(ctx context.Context) {
	backoff := time.Second * 2
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := d.consumeOnce(ctx); err != nil {
			d.logger.Error("mail consumer stopped", "error", err.Error())
			time.Sleep(backoff)
		}
	}
}

func (d *Dispatcher) consumeOnce(ctx context.Context) error {
	conn, err := amqp.Dial(d.cfg.RabbitMQURL())
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	defer ch.Close()

	if _, err := ch.QueueDeclare(mailQueueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare: %w", err)
	}
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("qos: %w", err)
	}

	msgs, err := ch.Consume(mailQueueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	d.logger.Info("mail rabbitmq consumer started", "queue", mailQueueName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}
			var job MailJob
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				d.logger.Error("mail job json", "error", err.Error())
				_ = msg.Nack(false, false)
				continue
			}
			if err := d.deliver(job); err != nil {
				d.logger.Error("mail send failed", "error", err.Error())
				_ = msg.Nack(false, true)
				continue
			}
			_ = msg.Ack(false)
		}
	}
}

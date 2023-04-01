package domain

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/internal/model/outbox"
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/protobuf/proto"
)

var (
	ErrRecordNotExists = errors.New("record is not found in outbox")
)

// CreateStatusMessage creates a new outbox.Message from loms.OrderStatus and adds it to the outbox.
func (d *Domain) CreateStatusMessage(ctx context.Context, orderID int64, status model.Status) error {
	statuspb := &loms.OrderStatus{
		OrderId: orderID,
		Status:  loms.Status(status),
	}

	raw, err := proto.Marshal(statuspb)
	if err != nil {
		return errors.WithMessage(err, "marshalling order status")
	}

	err = d.LOMSRepo.AddMessageWithKey(ctx, strconv.FormatInt(orderID, 10), raw)
	if err != nil {
		return errors.WithMessage(err, "saving message to outbox")
	}

	return nil
}

// SendMessage sends outbox.Message to the broker.
func (d *Domain) SendMessage(ctx context.Context, msg outbox.Message) error {
	err := d.LOMSRepo.UpdateMessageStatus(ctx, msg.ID, outbox.Processing)
	if err != nil {
		return err
	}

	d.StatusSender.SendWithKey(msg.ID, msg.Key, msg.Payload)

	return nil
}

// OnSendSuccess is used when the broker confirmed that the message was successfully sent.
// It deletes such message from outbox.
func (d *Domain) OnSendSuccess(id int64) error {
	// context.Background is used here as we might need to update outbox statuses during service shutdown.
	return d.LOMSRepo.DeleteMessage(context.Background(), id)
}

// OnSendFail is used when the broker did not send the message.
// It updates such message's status to outbox.Failed. The message remains in the outbox.
func (d *Domain) OnSendFail(id int64) error {
	// context.Background is used here as we might need to update outbox statuses during service shutdown.
	return d.LOMSRepo.UpdateMessageStatus(context.Background(), id, outbox.Failed)
}

// MonitorUnsent retrieves messages that should be sent to the broker and tries to submit them.
func (d *Domain) MonitorUnsent(ctx context.Context, errChan chan error) {
	ticker := time.NewTicker(d.config.SendInterval)
	for {
		select {
		case <-ticker.C:
			messages, err := d.LOMSRepo.ListUnsent(ctx)
			if err != nil {
				errChan <- err
				return
			}

			for _, message := range messages {
				err = d.SendMessage(ctx, message)
				if err != nil {
					errChan <- err
				}
			}

		case <-ctx.Done():
			close(errChan)
			return
		}
	}
}

// MonitorSenderResult checks Success and Errors channels provided by sender.
func (d *Domain) MonitorSenderResult(ctx context.Context, successChan, failChan chan int64, errChan chan error) {
	for {
		select {
		case id := <-successChan:
			err := d.OnSendSuccess(id)
			if err != nil {
				errChan <- err
			}

		case id := <-failChan:
			err := d.OnSendFail(id)
			if err != nil {
				errChan <- err
			}

		case <-ctx.Done():
			close(errChan)
			return
		}
	}
}

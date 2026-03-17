package jrnl

import (
	"fmt"
	"iter"

	"github.com/emersion/go-imap/v2/imapclient"
)

type (
	Receiver struct {
		server   string
		port     int
		user     string
		password string
		mailbox  string

		allowedSender             []string
		deleteFromUnallowedSender bool

		client *imapclient.Client
	}

	Mail struct {
		Subject string
		Body    string
	}
)

// NewReceiver creates a new receiver.
func NewReceiver(server string, port int, user, password, mailbox string) Receiver {
	return Receiver{
		server:   server,
		port:     port,
		user:     user,
		password: password,
		mailbox:  mailbox,
	}
}

// Start starts the receiver.
// It connects to the IMAP server and authenticates the user.
func (r *Receiver) Start() error {
	address := fmt.Sprintf("%s:%d", r.server, r.port)
	c, err := imapclient.DialTLS(address, nil)
	if err != nil {
		return fmt.Errorf("connect to IMAP server: %w", err)
	}

	if err := c.Login(r.user, r.password).Wait(); err != nil {
		_ = c.Logout().Wait()
		return fmt.Errorf("authenticate: %w", err)
	}

	if _, err := c.Select(r.mailbox, nil).Wait(); err != nil {
		_ = c.Logout().Wait()
		return fmt.Errorf("select mailbox %q: %w", r.mailbox, err)
	}

	r.client = c
	return nil
}

// Stop stops the receiver.
// It disconnects from the IMAP server.
func (r *Receiver) Stop() error {
	if r.client == nil {
		return nil
	}
	err := r.client.Logout().Wait()
	r.client = nil
	return err
}

// GetMails returns a sequence of mails in the inbox.
func (r *Receiver) GetMails() iter.Seq[Mail] {
	return nil
}

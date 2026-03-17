package jrnl

import "iter"

type (
	Receiver struct {
		server   string
		port     int
		user     string
		password string
		mailbox  string

		allowedSender             []string
		deleteFromUnallowedSender bool
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
func (r *Receiver) Start() error {}

// Stop stops the receiver.
// It disconnects from the IMAP server.
func (r *Receiver) Stop() error {}

// GetMails returns a sequence of mails in the inbox.
func (r *Receiver) GetMails() iter.Seq[Mail] {
	return nil
}

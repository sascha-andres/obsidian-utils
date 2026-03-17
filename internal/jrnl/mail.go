package jrnl

import (
	"fmt"
	"iter"
	"strconv"

	"github.com/emersion/go-imap/v2"
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
		MailID  string
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

	if r.mailbox == "" {
		r.mailbox = "INBOX"
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

// Move moves the given mail to the destination folder.
func (r *Receiver) Move(m Mail, destination string) error {
	if r.client == nil {
		return fmt.Errorf("not connected")
	}

	uid, err := strconv.ParseUint(m.MailID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid mail ID %q: %w", m.MailID, err)
	}

	uidSet := imap.UIDSetNum(imap.UID(uid))
	if _, err := r.client.Move(uidSet, destination).Wait(); err != nil {
		return fmt.Errorf("move to %q: %w", destination, err)
	}
	return nil
}

// GetMails returns a sequence of mails in the inbox.
func (r *Receiver) GetMails() iter.Seq[Mail] {
	return func(yield func(Mail) bool) {
		if r.client == nil {
			return
		}

		var seqSet imap.SeqSet
		seqSet.AddRange(1, 0) // 1:* — all messages

		fetchOptions := &imap.FetchOptions{
			UID:      true,
			Envelope: true,
			BodySection: []*imap.FetchItemBodySection{
				{Specifier: imap.PartSpecifierText},
			},
		}

		cmd := r.client.Fetch(seqSet, fetchOptions)
		defer cmd.Close()

		for {
			msg := cmd.Next()
			if msg == nil {
				break
			}

			buf, err := msg.Collect()
			if err != nil {
				break
			}

			var subject string
			if buf.Envelope != nil {
				subject = buf.Envelope.Subject
			}

			var body string
			if raw := buf.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText}); raw != nil {
				body = string(raw)
			} else if len(buf.BodySection) > 0 {
				body = string(buf.BodySection[0].Bytes)
			}

			m := Mail{
				MailID:  fmt.Sprintf("%d", buf.UID),
				Subject: subject,
				Body:    body,
			}
			if !yield(m) {
				return
			}
		}
	}
}

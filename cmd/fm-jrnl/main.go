package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/sascha-andres/obsidian-utils/internal"
	"github.com/sascha-andres/obsidian-utils/internal/jrnl"
	"github.com/sascha-andres/reuse/flag"
)

var (
	server, user, password, mailbox string
	port                            int
	logLevel                        string
	listMailsOnly                   bool
)

func init() {
	flag.SetEnvPrefix("FM_JRNL")
	flag.StringVar(&server, "server", "imap.gmail.com", "IMAP server")
	flag.StringVar(&user, "user", "", "IMAP user")
	flag.StringVar(&password, "password", "", "IMAP password")
	flag.IntVar(&port, "port", 993, "IMAP port")
	flag.BoolVar(&listMailsOnly, "list-mails-only", false, "list mails only")
	flag.StringVar(&logLevel, "log-level", "info", "pass log level (debug/info/warn/error)")
}

func main() {
	flag.Parse()

	logger := internal.CreateLogger(logLevel, "OBS_UTIL_ICAL")

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		logger.Error("failed to read build info")
	} else {
		logger = logger.With(slog.String("build", bi.Main.Version))
	}

	if err := run(logger); err != nil {
		logger.Error("failed to run fm-jrnl", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	logger.Info("Starting fm-jrnl")
	defer logger.Info("Done with fm-jrnl")

	r, err := jrnl.NewReceiver(server, port, user, password, mailbox)
	if err != nil {
		return err
	}
	err = r.Start()
	if err != nil {
		logger.Error("failed to start receiver", "error", err)
		return err
	}
	defer func() {
		if err := r.Stop(); err != nil {
			logger.Error("failed to stop receiver", "error", err)
		}
	}()

	for mail := range r.GetMails() {
		logger.Info("Received mail", "mail", mail)

		if listMailsOnly {
			continue
		}

		// TODO
	}

	return nil
}

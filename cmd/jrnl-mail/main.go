package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/obsidian-utils/internal"
	"github.com/sascha-andres/obsidian-utils/internal/jrnl"
)

var (
	folder, forDate, dailyFolder string
	headline                     = "## Other stuff"
	logLevel                     string
	dryRun                       bool

	server, user, password, mailbox string
	port                            int
	recipients                      string
)

func init() {
	internal.AddCommonFlagPrefixes()
	flag.SetEnvPrefix("OBS_UTIL_JRNL_MAIL")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.StringVar(&folder, "folder", "", "base path to obsidian vault")
	flag.StringVar(&dailyFolder, "daily-folder", "", "where to store the daily note inside the vault")
	flag.StringVar(&forDate, "for-date", time.Now().Format(time.DateOnly), "date for which to create the daily note (2006-01-02)")
	flag.StringVar(&headline, "headline", headline, fmt.Sprintf("headline under which to place the journal note (default: %s)", headline))
	flag.BoolVar(&dryRun, "dry-run", false, "pass to not edit file but to print added line with some context")

	flag.StringVar(&server, "server", "", "server")
	flag.IntVar(&port, "port", 993, "port")
	flag.StringVar(&user, "user", "", "user")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&mailbox, "mailbox", "", "mailbox")
	flag.StringVar(&recipients, "recipients", "jrnl@mailbox.org", "recipients")
}

func main() {
	flag.Parse()

	logger := internal.CreateLogger("OBS_UTIL_DAILY", logLevel)

	logger.Debug("start looking at mails")
	defer logger.Debug("done looking at mails")

	if err := run(logger); err != nil {
		logger.Error("error running daily", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	m := jrnl.NewReceiver(server, port, user, password, mailbox)

	allowedRecepients := strings.Split(recipients, ",")

	err := m.Start()
	if err != nil {
		return err
	}
	defer func() {
		err := m.Stop()
		if err != nil {
			logger.Error("error stopping receiver", "err", err)
		}
	}()

	for m := range m.GetMails() {
		if len(allowedRecepients) > 0 {
			found := false
			for _, recipient := range allowedRecepients {
				if recipient == m.Receiver {
					found = true
					break
				}
			}
			if !found {
				continue
			}
			logger.Info("processing mail", "m", m)
		}
	}

	return nil
}

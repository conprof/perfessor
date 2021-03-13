package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/conprof/perfessor/pkg/perfessor"
	"github.com/conprof/perfessor/pkg/shipper"
	"github.com/go-kit/kit/log/level"
	"github.com/thanos-io/thanos/pkg/logging"
	labelpb "github.com/thanos-io/thanos/pkg/store/labelpb"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Perfessor - continuous profiling utility")

	// Application settings
	loglevel := app.Flag("log.level", "Log filtering level").Default("info").Enum("debug", "warn", "info", "error")
	debugName := app.Flag("log.debug.name", "Name to add as prefix to log lines.").Hidden().String()
	logFormat := app.Flag("log.format", "Log format to use. Possible options: logfmt or json.").Default(logging.LogFormatLogfmt).Enum(logging.LogFormatLogfmt, logging.LogFormatJSON)
	freq := app.Flag("frequency", "Time between shipping profile and starting the next profile").Default("2m").Duration()

	// Profiling settings
	duration := app.Flag("perf.duration", "Duration of each profile").Default("10s").Duration()
	pnames := app.Flag("perf.pname", "pname of process to record").Strings()

	// Shipping settings
	storeAddress := app.Flag("ship.store", "Address of writable profile store.").Required().String()
	bearerToken := app.Flag("ship.bearer-token", "Bearer token to authenticate with store.").String()
	insecure := app.Flag("ship.insecure", "Send gRPC requests via plaintext instead of TLS.").Default("false").Bool()
	labels := app.Flag("ship.label", "label to apply to profiles. (kv pair in the format of k=v").Strings()

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Usage(os.Args[1:])
		log.Fatal("bad arguments")
	}

	logger := logging.NewLogger(*loglevel, *logFormat, *debugName)

	defaultLabels := make([]labelpb.Label, 0, len(*labels))
	for _, pair := range *labels {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}
		defaultLabels = append(defaultLabels, labelpb.Label{
			Name:  kv[0],
			Value: kv[1],
		})
	}

	// TODO move shipper creation into perfessor.Run
	shipper, err := shipper.NewShipper(
		*storeAddress,
		&shipper.Options{
			BearerToken:   *bearerToken,
			Insecure:      *insecure,
			DefaultLabels: defaultLabels,
		},
	)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create shipper", "error", err)
		os.Exit(1)
	}

	config := &perfessor.Config{
		Freq:      *freq,
		Processes: *pnames,
		Duration:  *duration,
		Shipper:   shipper,
		Logger:    logger,
	}

	if err := perfessor.Run(context.Background(), config); err != nil {
		level.Error(logger).Log("msg", "error after Run", "error", err)
		os.Exit(1)
	}
}

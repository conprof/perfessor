package perfessor

import (
	"context"
	"strings"
	"time"

	"github.com/conprof/perfessor/pkg/profile"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/mitchellh/go-ps"
	"github.com/thanos-io/thanos/pkg/store/labelpb"
)

// Shipper interface to ship profiles to external storage
type Shipper interface {
	Ship(ctx context.Context, profile []byte, labels ...labelpb.Label) error
}

// Config is the perfessor configuration
type Config struct {
	Processes []string
	Freq      time.Duration
	Duration  time.Duration
	Shipper   Shipper
	Logger    log.Logger
}

// Run takes periodic perf profiles
func Run(ctx context.Context, cfg *Config) error {
	level.Info(cfg.Logger).Log("msg", "starting profiling", "freq", cfg.Freq, "duration", cfg.Duration)
	duration := int(cfg.Duration.Seconds())
	firstRun := true
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !firstRun {
				time.Sleep(cfg.Freq)
			}
			firstRun = false

			// List all processes
			p, err := Filter(cfg.Processes)
			if err != nil {
				level.Error(cfg.Logger).Log("msg", "list pids failed", "error", err)
				continue
			}

			if len(p) == 0 {
				level.Error(cfg.Logger).Log("msg", "no processes found that match", "procs", strings.Join(cfg.Processes, ","))
				continue
			}

			// TODO can we perf processes in parallel
			for _, proc := range p {
				name := proc.Executable()
				pid := proc.Pid()

				err := profile.Record(pid, duration, cfg.Logger)
				if err != nil {
					level.Error(cfg.Logger).Log("msg", "record failed", "pname", name, "error", err)
					continue
				}

				pprof, err := profile.Convert(pid, cfg.Logger)
				if err != nil {
					level.Error(cfg.Logger).Log("msg", "convert failed", "pname", name, "error", err)
					continue
				}

				err = cfg.Shipper.Ship(ctx, pprof, labelpb.Label{
					Name:  "pname",
					Value: name,
				})
				if err != nil {
					level.Error(cfg.Logger).Log("msg", "failed to ship profile", "pname", name, "error", err)
					continue
				}
			}
		}
	}
}

// Filter out process names
func Filter(proc []string) ([]ps.Process, error) {
	procs, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	if len(proc) == 0 { // empty list return all
		return procs, nil
	}

	list := []ps.Process{}
	for _, p := range procs {
		for _, n := range proc {
			if p.Executable() == n {
				list = append(list, p)
			}
		}
	}

	return list, nil
}

package profile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Record records a CPU profile using perf
func Record(pid, duration int, logger log.Logger) error {
	p := strconv.Itoa(pid)
	d := strconv.Itoa(duration)
	perf := exec.Command("perf", "record", "-g", "-p", p, "--", "sleep", d)
	level.Debug(logger).Log("msg", perf.String())
	var b bytes.Buffer
	perf.Stderr = &b
	err := perf.Run()
	if err != nil {
		return fmt.Errorf("perf failed: %v, stderr: %v", err, string(b.Bytes()))
	}

	level.Info(logger).Log("msg", "perf taken")
	return nil
}

// Convert a perf.data file into a pprof proto
func Convert(pid int, logger log.Logger) ([]byte, error) {
	err := os.Setenv("PPROF_BINARY_PATH", filepath.Join("/proc", strconv.Itoa(pid), "root"))
	if err != nil {
		return nil, fmt.Errorf("failed to set env var PPROF_BINARY_PATH=%v: %v", filepath.Join("/proc", strconv.Itoa(pid), "root"), err)
	}

	pprof := exec.Command("pprof", "-proto", "perf.data")
	level.Debug(logger).Log("msg", pprof.String())
	var b bytes.Buffer
	var e bytes.Buffer
	pprof.Stdout = &b
	pprof.Stderr = &e
	err = pprof.Run()
	if err != nil {
		return nil, fmt.Errorf("pprof failed: %v, stderr: %v", err, string(e.Bytes()))
	}

	level.Debug(logger).Log("msg", "pprof -proto finished")
	return b.Bytes(), nil
}

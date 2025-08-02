package runner

import (
	"flag"
	"time"

	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/internal/updater"
	"github.com/projectdiscovery/gologger"
)

// Options defines the values needed to execute the Runner.
func Options() *common.Options {
	opt := &common.Options{}

	flag.StringVar(&opt.File, "f", "", "")
	flag.StringVar(&opt.File, "file", "", "")

	flag.StringVar(&opt.Address, "a", "", "")
	flag.StringVar(&opt.Address, "address", "", "")

	flag.StringVar(&opt.Auth, "A", "", "")
	flag.StringVar(&opt.Auth, "auth", "", "")

	flag.BoolVar(&opt.Check, "c", false, "")
	flag.BoolVar(&opt.Check, "check", false, "")

	flag.StringVar(&opt.CC, "only-cc", "", "")

	flag.DurationVar(&opt.Timeout, "t", 30*time.Second, "")
	flag.DurationVar(&opt.Timeout, "timeout", 30*time.Second, "")

	flag.IntVar(&opt.Rotate, "r", 1, "")
	flag.IntVar(&opt.Rotate, "rotate", 1, "")

	flag.BoolVar(&opt.RotateOnErr, "rotate-on-error", false, "")
	flag.BoolVar(&opt.RemoveOnErr, "remove-on-error", false, "")

	flag.StringVar(&opt.Method, "m", "sequent", "")
	flag.StringVar(&opt.Method, "method", "sequent", "")

	flag.BoolVar(&opt.Sync, "s", false, "")
	flag.BoolVar(&opt.Sync, "sync", false, "")

	flag.BoolVar(&opt.Verbose, "v", false, "")
	flag.BoolVar(&opt.Verbose, "verbose", false, "")

	flag.BoolVar(&opt.Daemon, "d", false, "")
	flag.BoolVar(&opt.Daemon, "daemon", false, "")

	flag.StringVar(&opt.Output, "o", "", "")
	flag.StringVar(&opt.Output, "output", "", "")

	flag.StringVar(&opt.OutputFormat, "output-format", "", "")

	flag.BoolVar(&doUpdate, "u", false, "")
	flag.BoolVar(&doUpdate, "update", false, "")

	flag.BoolVar(&version, "V", false, "")
	flag.BoolVar(&version, "version", false, "")

	flag.BoolVar(&opt.Watch, "w", false, "")
	flag.BoolVar(&opt.Watch, "watch", false, "")

	flag.IntVar(&opt.Goroutine, "g", 50, "")
	flag.IntVar(&opt.Goroutine, "goroutine", 50, "")

	flag.IntVar(&opt.MaxErrors, "max-errors", 3, "")
	flag.IntVar(&opt.MaxRedirects, "max-redirs", 10, "")
	flag.IntVar(&opt.MaxRetries, "max-retries", 0, "")

	flag.Usage = func() {
		showBanner()
		showUsage()
	}
	flag.Parse()

	if version {
		showVersion()
	}
	showBanner()

	if doUpdate {
		if err := updater.New(); err != nil {
			gologger.Fatal().Msgf("Error! %s.", err)
		}
	}

	if err := validate(opt); err != nil {
		gologger.Fatal().Msgf("Error! %s.", err)
	}

	return opt
}

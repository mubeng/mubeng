package runner

import (
	"errors"
	"fmt"
	"github.com/mubeng/mubeng/pkg/mubeng"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mubeng/mubeng/common"
	"github.com/mubeng/mubeng/internal/proxymanager"
)

// validate user-supplied option values before Runner.
func validate(opt *common.Options) error {
	var err error

	if hasStdin() {
		tmp, err := os.CreateTemp("", "mubeng-stdin-*")
		if err != nil {
			return err
		}
		defer tmp.Close()

		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		if _, err := tmp.Write(data); err != nil {
			return err
		}

		opt.File = tmp.Name()

		defer os.Remove(opt.File)
	}

	if opt.File == "" {
		return errors.New("no proxy file provided")
	}

	opt.File, err = filepath.Abs(opt.File)
	if err != nil {
		return err
	}
	// get server type
	mubeng.ServerType = opt.Type
	opt.ProxyManager, err = proxymanager.New(opt.File)
	if err != nil {
		return err
	}

	validMethod := map[string]bool{
		"sequent": true,
		"random":  true,
	}

	if opt.Address != "" && !opt.Check {
		if !validMethod[opt.Method] {
			return fmt.Errorf("unknown method for %q", opt.Method)
		}

		if opt.Auth != "" {
			auth := strings.SplitN(opt.Auth, ":", 2)
			if len(auth) != 2 {
				return errors.New("invalid proxy authorization format")
			}
		}
	}

	if opt.CC != "" {
		opt.Countries = strings.Split(opt.CC, ",")
	}

	if opt.Output != "" {
		opt.Output, err = filepath.Abs(opt.Output)
		if err != nil {
			return err
		}

		opt.Result, err = os.OpenFile(opt.Output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

package observer

import (
	"github.com/hellflame/argparse"
	. "github.com/usetheplatform/ci-system/pkg/common"
)

type Options struct {
	Path              string
	DispatcherAddress string
}

func GetOptions() (*Options, error) {
	parser := argparse.NewParser("observer", "observes the repo", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	path := parser.String("p", "path", &argparse.Option{
		Required: true,
	})

	dispatcherAddress := parser.String("da", "dispatcher-address", &argparse.Option{
		Required: false,
		HintInfo: "dispatcher host:port",
		Default:  "localhost:8080",
	})

	err := parser.Parse(nil)

	CheckIfError(err)

	return &Options{
		Path:              *path,
		DispatcherAddress: *dispatcherAddress,
	}, nil
}

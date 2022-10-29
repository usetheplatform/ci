package dispatcher

import (
	"github.com/hellflame/argparse"
	. "github.com/usetheplatform/ci-system/pkg/common"
)

type Options struct {
	Host string
	Port string
}

func GetOptions() (*Options, error) {
	parser := argparse.NewParser("dispatcher", "delegates tasks to job runners", &argparse.ParserConfig{
		DisableDefaultShowHelp: true,
	})

	host := parser.String("h", "host", &argparse.Option{
		Default:  "localhost",
		Required: false,
	})

	port := parser.String("p", "port", &argparse.Option{
		Default:  "8080",
		Required: false,
	})

	err := parser.Parse(nil)
	CheckIfError(err)

	return &Options{
		Host: *host,
		Port: *port,
	}, nil
}

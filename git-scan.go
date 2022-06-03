package main

import (
	"context"
	"flag"
	"predorock/gitscan/git"
	"predorock/gitscan/logstreamer"
	scanner "predorock/gitscan/scan"

	"github.com/cristalhq/acmd"
)

var log = logstreamer.GetInstance().Log

type commonFlags struct {
	parallel bool
}

// NOTE: should be added before flag.FlagSet method Parse().
func withCommonFlags(fs *flag.FlagSet) *commonFlags {
	c := &commonFlags{}
	fs.BoolVar(&c.parallel, "parallel", false, "run scan in parallel")
	//fs.StringVar(&c.Dir, "dir", ".", "directory to process")
	return c
}

func main() {

	cmds := []acmd.Command{
		{
			Name:        "update",
			Description: "fetch and pull on current repo branch",
			Do: func(ctx context.Context, args []string) error {
				cmdFlags := flag.NewFlagSet("update", flag.ContinueOnError)
				common := withCommonFlags(cmdFlags)

				if err := cmdFlags.Parse(args); err != nil {
					return err
				}

				for _, f := range cmdFlags.Args() {
					log.Printf("Search for repos in folder %s\n", f)

					if common.parallel {
						scanner.ScanParallel(f, git.UpdateCommand)
					} else {
						scanner.Scan(f, git.UpdateCommand)
					}

				}
				return nil
			},
		},
	}

	// all the acmd.Config fields are optional
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "gitscan",
		AppDescription: "Do git stuff with multiple repos at the same time\nby Marco Predari <predorock@gmail.com>",
		Version:        "v0.0.1",
		// Context - if nil `signal.Notify` will be used
		// Args - if nil `os.Args[1:]` will be used
		// Usage - if nil default print will be used
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}

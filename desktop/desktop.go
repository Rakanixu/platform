package desktop

import (
	"log"

	"github.com/micro/cli"
	"os"
	"os/exec"
	"sync"
)

func all(ctx *cli.Context) {
	var wg sync.WaitGroup
	cmds := ctx.App.Commands

	log.Print(ctx.Args())
	for _, cmd := range cmds {
		if cmd.Name != "help" && len(cmd.Subcommands) > 0 {
			for _, subcmd := range cmd.Subcommands {
				//time.Sleep(time.Second)
				wg.Add(1)
				log.Print(cmd.Name, subcmd.Name)
				c := exec.Command("./kazoup", "--registry=mdns", cmd.Name, subcmd.Name)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Start(); err != nil {
					log.Print(err.Error())
					wg.Done()
				}
			}
		}
		if cmd.Name != "help" && len(cmd.Subcommands) == 0 && cmd.Name != "desktop" {

			wg.Add(1)
			log.Print(cmd.Name)
			c := exec.Command("./kazoup", "--registry=mdns", cmd.Name)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Start(); err != nil {
				log.Print(err.Error())
				wg.Done()
			}
		}
	}
	wg.Wait()

}

func desktopCommands() []cli.Command {
	return []cli.Command{{
		Name:   "desktop",
		Usage:  "Run desktop service",
		Action: all,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "desktop",
			Usage:       "Desktop commands",
			Subcommands: desktopCommands(),
		},
	}
}

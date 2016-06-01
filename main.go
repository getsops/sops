package main /* import "gozilla.io/sops" */

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "sops"
	app.Usage = "Secrets management stinks, use some sops!"
	app.UsageText = "sops <file>"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "decrypt, d",
			Usage: "decrypt <file> and print it to stdout",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			return cli.NewExitError("error: <file> not specified", 1)
		}
		fileName := c.Args()[0]
		if c.Bool("decrypt") {
			if err := DecryptFile(fileName); err != nil {
				return cli.NewExitError(fmt.Sprintf("Error decrypting %s: %v", fileName, err), 1)
			}
		}
		return nil
	}

	app.Run(os.Args)
}

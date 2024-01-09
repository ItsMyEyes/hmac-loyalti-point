package cli

import (
	"os"

	"github.com/urfave/cli"
)

func InitCLi() {
	app := cli.NewApp()
	app.Name = "hmac-loyalti"
	app.UsageText = "hmac-loyalti [global options] command [command options] [arguments...]"
	app.Description = "hmac-loyalti is a web server for microservice architecture"
	app.Authors = []cli.Author{
		{
			Name:  "Indo Artha Graha",
			Email: "dev@iat.id",
		},
		{
			Name:  "IAT Dev Team",
			Email: "dev@iat.id",
		},
	}
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:        "http",
			Description: "http server for hmac-loyalti",
			Usage:       "run http server for hmac-loyalti",
			UsageText:   "hmac-loyalti http -c config.yaml",
			Aliases:     []string{"h"},
			Action:      GenerateHmacHttp,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "config, c",
					Value:    "config.yaml",
					Usage:    "config file path",
					Required: true,
				},
			},
		},
		{
			Name:        "generate",
			Description: "generate hmac-loyalti",
			Aliases:     []string{"g"},
			Usage:       "generate hmac-loyalti",
			UsageText:   "hmac-loyalti generate -c config.yaml",
			HelpName:    "generate",
			Action:      GenerateHmacCli,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "config, c",
					Value:    "config.yaml",
					Usage:    "config file path",
					Required: true,
				},
				cli.StringFlag{
					Name:     "path, p",
					Value:    "/core/v1/transaction",
					Usage:    "Path for your request [example] /core/v1/transaction",
					Required: true,
				},
				cli.StringFlag{
					Name:     "method, m",
					Value:    "POST",
					Usage:    "Method for your request",
					Required: true,
				},
				cli.StringFlag{
					Name:  "body, b",
					Value: "{nothingHereBrooo}",
					Usage: "Body for your request",
				},
				cli.StringFlag{
					Name:  "timestamp, t",
					Value: "2019-01-01 00:00:00",
					Usage: "Body for your request",
				},
			},
		},
	}

	app.Run(os.Args)
}

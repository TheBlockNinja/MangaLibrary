package main

import "github.com/urfave/cli/v2"

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "port-number",
		Usage: "server-port-number",
		EnvVars: []string{
			"PORT_NUMBER",
		},
		Required:  false,
		Hidden:    false,
		TakesFile: false,
		Value:     "8080",
	},
	&cli.StringFlag{
		Name:  "mysql-host",
		Usage: "mysql DB Hostname",
		EnvVars: []string{
			"MYSQL_HOST",
		},
		Required:  false,
		Hidden:    false,
		TakesFile: false,
		Value:     "mysql",
	},
	&cli.StringFlag{
		Name:  "mysql-port",
		Usage: "mysql DB port",
		EnvVars: []string{
			"MYSQL_PORT",
		},
		Required:  false,
		Hidden:    false,
		TakesFile: false,
		Value:     "3306",
	},
	&cli.StringFlag{
		Name:  "mysql-user",
		Usage: "mysql DB user",
		EnvVars: []string{
			"MYSQL_USER",
		},
		Required:  false,
		Hidden:    false,
		TakesFile: false,
		Value:     "root",
	},
	&cli.StringFlag{
		Name:  "mysql-password",
		Usage: "mysql DB password",
		EnvVars: []string{
			"MYSQL_PASSWORD",
		},
		Required:  false,
		Hidden:    false,
		TakesFile: false,
		Value:     "password",
	},
	&cli.IntFlag{
		Name:  "number-of-workers",
		Usage: "Number of workers to process jobs",
		EnvVars: []string{
			"WORKERS",
		},
		Value: 2,
	},
	&cli.StringFlag{
		Name:  "log-level",
		Usage: "current_log_level",
		EnvVars: []string{
			"LOG_LEVEL",
		},
		Value: "info",
	},
	&cli.BoolFlag{
		Name:  "log-production",
		Usage: "current_log_level",
		EnvVars: []string{
			"PRODUCTION",
		},
		Value: false,
	},
}

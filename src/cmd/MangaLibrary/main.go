package main

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/jobs"
	"MangaLibrary/src/internal/server"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:  "WebSearch(Web)",
		Usage: "",
		Action: func(c *cli.Context) error {
			err := app(c)
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
		},
		Flags: flags,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
func app(c *cli.Context) error {
	loggerConfig := zap.NewDevelopmentConfig()
	if c.Bool("log-production") {
		loggerConfig = zap.NewProductionConfig()
	}
	if err := loggerConfig.Level.UnmarshalText([]byte(c.String("log-level"))); err != nil {
		return err
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		return err
	}

	db, err := dao.GetDB(c)
	if err != nil {
		logger.Error("Failed loading DB", zap.Error(err))
		return err
	}
	webComponentDao := &dao.WebComponentDAO{DB: db}
	jobDao := &dao.JobDAO{DB: db}
	pj := jobs.ProcessJobs{
		Logger:  logger,
		MaxJobs: c.Int("number-of-workers"),
		JobDAO:  jobDao,
		CompDao: webComponentDao,
	}
	err = pj.ResetFailedJobs()
	if err != nil {
		logger.Error("Failed loading DB", zap.Error(err))
		return err
	}
	go func() {
		pj.Start()
	}()

	return server.StartServer(db, logger)
}

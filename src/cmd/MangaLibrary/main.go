package main

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/jobs"
	"MangaLibrary/src/internal/server"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewProduction()
	//webDriver.NewWebDriver()
	app := &cli.App{
		Name:  "WebSearch(Web)",
		Usage: "",
		Action: func(c *cli.Context) error {
			err := app(c, logger)
			if err != nil {
				logger.Error("could not start app", zap.Error(err))
				return err
			}
			return nil
		},
	}
	if err != app.Run(os.Args) {
		logger.Error("Failed to run app", zap.Error(err))
	}

}
func app(c *cli.Context, logger *zap.Logger) error {
	db, err := dao.GetDB()
	if err != nil {
		logger.Error("Failed loading DB", zap.Error(err))
		return err
	}
	webComponentDao := &dao.WebComponentDAO{DB: db}
	jobDao := &dao.JobDAO{DB: db}
	pj := jobs.ProcessJobs{
		Logger:  logger,
		MaxJobs: 2,
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

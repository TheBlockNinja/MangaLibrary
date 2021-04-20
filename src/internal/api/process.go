package api

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"MangaLibrary/src/internal/jobs"
	"MangaLibrary/src/internal/timezone"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/TheBlockNinja/WebParser"
	"go.uber.org/zap"
)

func StartJob(c *jobs.Component, job *jobs.Job, logger *zap.Logger, jobDAO *dao.JobDAO) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	jStr, _ := json.Marshal(job)
	cStr, _ := json.Marshal(c)
	fmt.Printf("JOB REQUEST: %s\nComponents:%s\n", string(jStr), string(cStr))
	go func() {
		logger.Info("working on job...")
		if job.CurrentJobData != "" && job.CurrentJobData != "{}" {
			err = json.Unmarshal([]byte(job.CurrentJobData), &c)
			if err != nil {
				logger.Error("failed processing job", zap.Error(err))
				job.TotalProgress = -1
				wg.Done()
				logger.Error("failed")
				job.Message = "failed"
				return
			}

		}
		currentTime, _ := timezone.GetTime("PST")
		job.EstFinish = currentTime
		job.StartTime = currentTime
		job.LastUpdate = currentTime
		switch job.Ctx.Type {
		case "download":
			logger.Info("downloading job")
			Parser := &WebParser.Parser{Logger: logger}
			testBook := &dto.Books{}
			err = c.Download(Parser, job.Ctx.BasePath, job, jobDAO, testBook)
			break
		case "pdf":
			break
		case "process":
			logger.Info("processing job")
			err = ProcessJob(c, job, logger, jobDAO)
			break
		default:
			err = fmt.Errorf("invalid ctx type")
		}
		if err != nil {
			logger.Error("failed processing job", zap.Error(err))
			job.TotalProgress = -1
		} else {
			bytes, errs := json.Marshal(c)
			if errs != nil {
				err = errs
			} else {
				job.CurrentJobData = string(bytes)
			}
		}
		logger.Info("done")
		job.Message = "success"
		wg.Done()

	}()
	go func() {
		err := UpdateJob(job, logger)
		if err != nil {
			//dto update failed job
			logger.Error("failed processing job", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Wait()
	return err
}
func ProcessJob(c *jobs.Component, job *jobs.Job, logger *zap.Logger, jobDAO *dao.JobDAO) error {
	newParser := WebParser.NewParser(logger)
	err := c.LoadSiteData(newParser, job, "", jobDAO)
	if err != nil {
		return err
	}
	return nil
}
func UpdateJob(job *jobs.Job, logger *zap.Logger) error {
	for job.CurrentProgress < job.TotalProgress || job.TotalProgress == 0 {
		time.Sleep(5 * time.Second)
		logger.Info(fmt.Sprintf("%s (%d/%d) : Finish %s", job.Name, job.CurrentProgress, job.TotalProgress, job.EstFinish.String()))
		if job.Message == "failed" {
			return nil
		}
	}

	logger.Info(fmt.Sprintf("%s (%d/%d) : Finish %s", job.Name, job.CurrentProgress, job.TotalProgress, job.EstFinish.String()))
	return nil
}

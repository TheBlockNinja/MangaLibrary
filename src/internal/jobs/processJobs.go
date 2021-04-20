package jobs

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ProcessJobs struct {
	Logger  *zap.Logger
	MaxJobs int
	JobDAO  *dao.JobDAO
	CompDao *dao.WebComponentDAO
}

func (p *ProcessJobs) ResetFailedJobs() error {
	for _, t := range []string{"in-progress", "downloading", "processing", "waiting"} {
		err := p.resetJobType(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessJobs) resetJobType(jobName string) error {
	jobs, err := p.JobDAO.GetJobsWith(jobName)
	if err != nil {
		return err
	}
	for _, j := range jobs {
		j.Message = "waiting"
		j.TotalProgress = 0
		j.CurrentProgress = 0
		err = p.JobDAO.UpdateJob(j)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProcessJobs) GrabJobs(jobs chan *dto.Job) error {
	for {
		time.Sleep(time.Second * 15)
		js, err := p.JobDAO.GetJobsWith("waiting")
		if err != nil {
			close(jobs)
			return err
		}
		for _, j := range js {
			p.Logger.Info(fmt.Sprintf("grabbed job %s", j.Name))
			j.Message = "queued"
			err = p.JobDAO.UpdateJob(j)
			if err != nil {
				close(jobs)
				return err
			}
			jobs <- j
		}
	}
}

func (p *ProcessJobs) Worker(jobs chan *dto.Job) error {
	jobDto := <-jobs

	p.Logger.Info(fmt.Sprintf("processing job %s", jobDto.Name))
	job, err := ConvertDtoToJob(jobDto)
	if err != nil {
		return err
	}
	types := []string{}
	err = json.Unmarshal([]byte(job.Ctx.Type), &types)
	if err != nil {
		jobDto.Message = "failed"
		err = p.JobDAO.UpdateJob(jobDto)
		return err
	}
	comps, err := p.CompDao.GetComponentsForSite(job.SiteID)
	if err != nil {
		jobDto.Message = "failed"
		err = p.JobDAO.UpdateJob(jobDto)
		return err
	}
	comp := WebComponentsToComponent(comps)
	jobDto.Message = "in-progress"
	err = p.JobDAO.UpdateJob(jobDto)

	err = job.StartJob(comp, types, p.Logger, p.JobDAO)
	if err != nil {
		p.Logger.Error("job failed to run", zap.Error(err))
		jobDto, err = job.ToDTO()
		if err != nil {
			return err
		}
		jobDto.Message = "failed"
		_ = p.JobDAO.UpdateJob(jobDto)
		return err
	}
	jobDto, err = job.ToDTO()
	if err != nil {
		return err
	}
	jobDto.Message = "Success"
	_ = p.JobDAO.UpdateJob(jobDto)
	return nil
}

func (p *ProcessJobs) Workers(jobs chan *dto.Job) {
	for {
		err := p.Worker(jobs)
		if err != nil {
			p.Logger.Error("worker failed", zap.Error(err))
			time.Sleep(60 * time.Second)
		}
		time.Sleep(5 * time.Second)
	}
}

func (p *ProcessJobs) Start() {
	jobs := make(chan *dto.Job, p.MaxJobs)
	go func() {
		err := p.GrabJobs(jobs)
		if err != nil {
			p.Logger.Error("failed grabbing jobs", zap.Error(err))
		}
	}()
	for i := 0; i < p.MaxJobs; i++ {
		go p.Workers(jobs)
		time.Sleep(1 * time.Second)
	}
}

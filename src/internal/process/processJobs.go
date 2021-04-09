package process

import (
	"MangaLibrary/src/internal/api"
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/driver"
	"MangaLibrary/src/internal/dto"
	jobs2 "MangaLibrary/src/internal/jobs"
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

func (p *ProcessJobs) ResetFailedJobs() {

}
func (p *ProcessJobs) Start() {
	jobs := make(chan *dto.Job, p.MaxJobs)
	go func() {
		for {
			js, err := p.JobDAO.GetJobsWith("waiting")
			if err != nil {
				return
			}
			for _, j := range js {
				p.Logger.Info(fmt.Sprintf("grabbed job %s", j.Name))
				j.Message = "queued"
				err = p.JobDAO.UpdateJob(j)
				if err != nil {
					return
				}
				jobs <- j
				time.Sleep(time.Second)
			}
		}
	}()
	go func() {
		for t := range jobs {
			p.Logger.Info(fmt.Sprintf("processing job %s", t.Name))
			job, err := jobs2.ConvertDtoToJob(t)
			if err != nil {
				continue
			}
			types := []string{}
			err = json.Unmarshal([]byte(job.Ctx.Type), &types)
			if err != nil {
				t.Message = "failed"
				err = p.JobDAO.UpdateJob(t)
				continue
			}

			t.Message = "in-progress"
			err = p.JobDAO.UpdateJob(t)

			comps, err := p.CompDao.GetComponentsForSite(job.SiteID)
			if err != nil {
				t.Message = "failed"
				err = p.JobDAO.UpdateJob(t)
				continue
			}
			comp := driver.WebComponentsToComponent(comps)
			for _, currentType := range types {
				job.Ctx.Type = currentType
				p.Logger.Info(fmt.Sprintf("starting job for context type %s", currentType))
				err := api.StartJob(comp, job, p.Logger, p.JobDAO)
				if err != nil {
					break
				}
			}

		}
	}()
}

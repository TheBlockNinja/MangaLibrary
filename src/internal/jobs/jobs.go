package jobs

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"MangaLibrary/src/internal/timezone"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheBlockNinja/WebParser"
	"go.uber.org/zap"
)

type Job struct {
	Id              int         `json:"id"`
	User            int         `json:"user"`
	SiteID          int         `json:"site_id"`
	Name            string      `db:"name"`
	Ctx             *JobContext `db:"job_context"`
	StartTime       time.Time   `db:"start_time"`
	CurrentProgress int         `db:"current"`
	TotalProgress   int         `db:"total_progress"`
	EstFinish       time.Time   `db:"est_finish"`
	LastUpdate      time.Time   `json:"last_update"`
	Message         string      `db:"message"`
	ForceStop       bool        `json:"force_stop"`
	CurrentJobData  string      `json:"current_job"`
}

type JobContext struct {
	URL      string
	BasePath string
	Type     string
}

func (j *Job) UpdateTime() {
	otherTime := j.EstFinish.Second()
	currentTime, _ := timezone.GetTime("PST")
	updateTime := Abs(currentTime.Second() - j.LastUpdate.Second())
	remaining := Abs(j.TotalProgress - j.CurrentProgress)
	willFinish := currentTime.Add(time.Second * time.Duration(((updateTime*remaining)+otherTime)/2))
	j.EstFinish = willFinish
	j.LastUpdate = currentTime
}

func (j *Job) UpdateDB(dao *dao.JobDAO) {
	dtoJob, err := j.ToDTO()
	if err != nil {
		fmt.Printf("error formatting to DTO %s\n", err.Error())
		return
	}
	if dao != nil {
		err = dao.UpdateJob(dtoJob)
		if err != nil {
			return
		}
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ConvertDtoToJob(jobDto *dto.Job) (*Job, error) {

	startT, err := timezone.TimeParse(jobDto.StartTime)
	if err != nil {
		return nil, err
	}
	estTime, err := timezone.TimeParse(jobDto.EstFinish)
	if err != nil {
		return nil, err
	}

	ctx := &JobContext{}
	err = json.Unmarshal([]byte(jobDto.JobContext), ctx)
	if err != nil {
		return nil, err
	}
	currentTime, _ := timezone.GetTime("PST")
	job := &Job{
		Id:              jobDto.ID,
		User:            jobDto.UserID,
		SiteID:          jobDto.SiteID,
		Name:            jobDto.Name,
		Ctx:             ctx,
		StartTime:       startT,
		CurrentProgress: jobDto.CurrentProgress,
		TotalProgress:   jobDto.TotalProgress,
		EstFinish:       estTime,
		LastUpdate:      currentTime,
		Message:         jobDto.Message,
		ForceStop:       false,
		CurrentJobData:  jobDto.CurrentJobData,
	}
	return job, nil
}

func (j *Job) ToDTO() (*dto.Job, error) {
	cxt, err := json.Marshal(j.Ctx)
	if err != nil {
		return nil, err
	}

	jDTO := &dto.Job{
		ID:              j.Id,
		UserID:          j.User,
		SiteID:          j.SiteID,
		Name:            j.Name,
		JobContext:      string(cxt),
		StartTime:       j.StartTime.Format(timezone.BasicFormat),
		CurrentProgress: j.CurrentProgress,
		TotalProgress:   j.TotalProgress,
		EstFinish:       j.StartTime.Format(timezone.BasicFormat),
		Message:         j.Message,
		CurrentJobData:  j.CurrentJobData,
	}
	return jDTO, err
}

func (j *Job) StartJob(c *Component, types []string, logger *zap.Logger, jobDAO *dao.JobDAO) error {
	for _, currentType := range types {
		j.Message = currentType + "ing"
		jobDto, err := j.ToDTO()
		if err != nil {
			return err
		}
		if jobDAO != nil {
			err = jobDAO.UpdateJob(jobDto)
			if err != nil {
				return err
			}
		}
		j.Ctx.Type = currentType
		err = j.loadJobData(c)
		if err != nil {
			return err
		}
		err = j.processJobType(c, currentType, logger, jobDAO)
		if err != nil {
			return err
		}
	}
	jobDto, err := j.ToDTO()
	if err != nil {
		return err
	}
	jobDto.Message = "success"
	if jobDAO != nil {
		err = jobDAO.UpdateJob(jobDto)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *Job) loadJobData(c *Component) error {
	if j.CurrentJobData != "" && j.CurrentJobData != "{}" {
		err := json.Unmarshal([]byte(j.CurrentJobData), &c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (j *Job) processJobType(c *Component, jobType string, logger *zap.Logger, jobDAO *dao.JobDAO) error {
	switch jobType {
	case "download":
		Parser := &WebParser.Parser{Logger: logger}
		testBook := &dto.Books{}
		err := c.Download(Parser, j.Ctx.BasePath, j, jobDAO, testBook)
		if err != nil {
			return err
		}
	case "pdf":
		break
	case "process":
		err := ProcessJob(c, j, logger, jobDAO)
		if err != nil {
			return err
		}
		break
	default:
		return fmt.Errorf("invalid ctx type")
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	} else {
		j.CurrentJobData = string(bytes)
	}
	return nil
}

func ProcessJob(c *Component, job *Job, logger *zap.Logger, jobDAO *dao.JobDAO) error {
	newParser := WebParser.NewParser(logger)
	err := c.LoadSiteData(newParser, job, "", jobDAO)
	if err != nil {
		return err
	}
	return nil
}

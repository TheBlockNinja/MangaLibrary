package jobs

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"encoding/json"
	"time"
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
	updateTime := Abs(time.Now().Second() - j.LastUpdate.Second())

	remaining := Abs(j.TotalProgress - j.CurrentProgress)
	willFinish := time.Now().Add(time.Second * time.Duration(((updateTime*remaining)+otherTime)/2))
	//	Parser.Logger.Info(fmt.Sprintf("Seconds since last update %d", updateTime))
	j.EstFinish = willFinish
	//progress.RemainingTime = willFinish.Format("2006-01-02 15:04:05")
	j.LastUpdate = time.Now()
}

func (j *Job) UpdateDB(dao *dao.JobDAO) {
	dtoJob, err := j.ToDTO()
	if err != nil {
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
	startT, err := time.Parse("2006-01-02 3:04PM", jobDto.StartTime)
	if err != nil {
		return nil, err
	}
	estTime, err := time.Parse("2006-01-02 3:04PM", jobDto.StartTime)
	if err != nil {
		return nil, err
	}

	ctx := &JobContext{}
	err = json.Unmarshal([]byte(jobDto.JobContext), ctx)
	if err != nil {
		return nil, err
	}

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
		LastUpdate:      time.Now(),
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
		StartTime:       j.StartTime.Format("2006-01-02 3:04PM"),
		CurrentProgress: j.CurrentProgress,
		TotalProgress:   j.TotalProgress,
		EstFinish:       j.StartTime.Format("2006-01-02 3:04PM"),
		Message:         j.Message,
		CurrentJobData:  j.CurrentJobData,
	}
	return jDTO, err
}

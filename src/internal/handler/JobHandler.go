package handler

import (
	"MangaLibrary/src/internal/dao"
	"MangaLibrary/src/internal/dto"
	"MangaLibrary/src/internal/jobs"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type JobRequest struct {
	JobName  string            `json:"name"`
	SiteName string            `json:"site_name"`
	Search   map[string]string `json:"search"`
	Types    []string          `json:"types"`
}
type JobHandler struct {
	Site   *dao.SitesDAO
	Job    *dao.JobDAO
	User   *dao.UsersDAO
	Logger *zap.Logger
}

func (h JobHandler) NewJob(w http.ResponseWriter, r *http.Request) {
	jobReq := &JobRequest{}
	err := json.NewDecoder(r.Body).Decode(jobReq)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	site, err := h.Site.GetSite(jobReq.SiteName)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	if site.MinAge > user.Age {
		SendError("to young to view this data", w)
		return
	}
	types, err := json.Marshal(jobReq.Types)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	ctx := jobs.JobContext{
		URL:      site.GetURL(jobReq.Search),
		BasePath: site.BasePath + "",
		Type:     string(types),
	}
	strCtx, err := json.Marshal(ctx)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	newJob := &dto.Job{
		UserID:          user.ID,
		SiteID:          site.ID,
		Name:            jobReq.JobName,
		JobContext:      string(strCtx),
		StartTime:       time.Now().Format("2006-01-02 3:04PM"),
		CurrentProgress: 0,
		TotalProgress:   0,
		EstFinish:       time.Now().Format("2006-01-02 3:04PM"),
		Message:         "waiting",
		CurrentJobData:  "{}",
	}
	err = h.Job.NewJob(newJob)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	SendData(newJob, "Added New Job", w)
}
func (h JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	user, err := h.User.GetUser(r.Header.Get("Authentication-Key"))
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	userJobs, err := h.Job.GetJobsForUser(user.ID)
	if err != nil {
		SendError(err.Error(), w)
		return
	}
	for i := range userJobs {
		userJobs[i].CurrentJobData = "censored"
	}
	SendData(userJobs, "Found Jobs", w)
}

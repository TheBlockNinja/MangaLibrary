package dao

import (
	"MangaLibrary/src/internal/dto"

	"github.com/jmoiron/sqlx"
)

type JobDAO struct {
	DB *sqlx.DB
}

func (s *JobDAO) NewJob(job *dto.Job) error {
	insertStmt := "insert into manga_library.jobs(user_id,site_id,name,job_context,current,total,start_time,est_finish,message,job_data) VALUES(?,?,?,?,?,?,?,?,?,?)"
	results, err := s.DB.Exec(insertStmt, job.UserID, job.SiteID, job.Name, job.JobContext, job.CurrentProgress, job.TotalProgress, job.StartTime, job.EstFinish, job.Message, job.CurrentJobData)
	if err != nil {
		return err
	}
	id, _ := results.LastInsertId()
	job.ID = int(id)
	return nil
}

func (s *JobDAO) GetJob(id int) (*dto.Job, error) {
	job := &dto.Job{}
	err := s.DB.Get(job, "select * from manga_library.jobs where id = ?", id)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobDAO) GetJobsForUser(userID int) ([]*dto.Job, error) {
	var jobList []*dto.Job
	rows, err := s.DB.Queryx("select * from manga_library.jobs where user_id = ? order by start_time", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := &dto.Job{}
		err = rows.StructScan(webComponent)
		if err != nil {
			return nil, err
		}
		jobList = append(jobList, webComponent)
	}
	return jobList, nil
}

func (s *JobDAO) GetJobsInProgress(userID int) ([]*dto.Job, error) {
	var jobList []*dto.Job
	rows, err := s.DB.Queryx("select * from manga_library.jobs where user_id = ? and current < total order by id", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := &dto.Job{}
		err = rows.StructScan(webComponent)
		if err != nil {
			return nil, err
		}
		jobList = append(jobList, webComponent)
	}
	return jobList, nil
}

func (s *JobDAO) GetJobsWith(message string) ([]*dto.Job, error) {
	var jobList []*dto.Job
	rows, err := s.DB.Queryx("select * from manga_library.jobs where message like ? order by start_time", "%"+message+"%")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		webComponent := &dto.Job{}
		err = rows.StructScan(webComponent)
		if err != nil {
			return nil, err
		}
		jobList = append(jobList, webComponent)
	}
	return jobList, nil
}

func (s *JobDAO) UpdateJob(job *dto.Job) error {
	_, err := s.DB.Exec("update manga_library.jobs set current = ?, total = ?, est_finish = ?,message = ?,job_data = ? where id = ? and user_id = ?",
		job.CurrentProgress, job.TotalProgress, job.EstFinish, job.Message, job.CurrentJobData, job.ID, job.UserID)
	if err != nil {
		return err
	}
	return nil
}

package handler

import (
	"encoding/json"
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/robfig/cron"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

type CronWrapper struct {
	Cron *cron.Cron
	Id   string
}

func (s *Scheduler) createTask(ctx context.Context, req *proto.CreateScheduledTaskRequest) (*proto.CreateScheduledTaskResponse, error) {
	taskRecognise := false

	if req.Task.Action == globals.StartScanTask {
		i := "@every " + strconv.Itoa(int(req.Schedule.IntervalSeconds)) + "s"

		c := cron.New()
		c.AddFunc(i, func() {
			var ds *datasource_proto.Endpoint
			dbC := db_proto.NewDBClient(globals.DB_SERVICE_NAME, s.Client)

			rsp, err := dbC.Read(context.Background(), &db_proto.ReadRequest{
				Index: "datasources",
				Type:  "datasource",
				Id:    req.Task.Id,
			})

			if err != nil {
				log.Println("ERROR", err)
				return
			}

			if err := json.Unmarshal([]byte(rsp.Result), &ds); err != nil {
				log.Println("ERROR", err)
				return
			}

			if len(ds.Id) > 0 {
				// Datasource is in db, so try run a scan if is not already running
				if !ds.CrawlerRunning {
					dsC := datasource_proto.NewDataSourceClient(globals.DATASOURCE_SERVICE_NAME, s.Client)

					// ctx it is not passed because circuit breaker, timeout will exceed when executed in the future
					_, err = dsC.Scan(context.Background(), &datasource_proto.ScanRequest{
						Id: req.Task.Id,
					})
					if err != nil {
						log.Println("ERROR", err)
						return
					}
				}
			} else {
				// Datasource is not in DB anymore
				// We supose is because user remove it, so lets remove the task associated with it
				// Datasources and tasks will be eventually consistent.
				for k, v := range s.Crons {
					if v.Id == req.Task.Id {
						// Stop cron task and delete CronWrapper instance from Scheduler
						v.Cron.Stop()
						s.Crons = append(s.Crons[:k], s.Crons[k+1:]...)
					}
				}
			}
		})
		c.Start()
		s.Crons = append(s.Crons, &CronWrapper{
			Id:   req.Task.Id,
			Cron: c,
		})

		taskRecognise = true
	}

	if !taskRecognise {
		return nil, errors.New("Task not recognise")
	}

	return &proto.CreateScheduledTaskResponse{}, nil
}

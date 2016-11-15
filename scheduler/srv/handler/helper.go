package handler

import (
	"encoding/json"
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"golang.org/x/net/context"
	"gopkg.in/robfig/cron.v2" // This is like vendoring pre GO 1.5 - github.com/robfig/cron
	"log"
	"strconv"
)

type CronWrapper struct {
	Cron   *cron.Cron
	CronId cron.EntryID
	Id     string
}

// TODO: FIXME: cron task are keep in memory, so two instance of this srv running will fail
// We need to sync the data between instances.
func (s *Scheduler) createTask(ctx context.Context, req *proto.CreateScheduledTaskRequest) (*proto.CreateScheduledTaskResponse, error) {
	taskRecognise := false

	if req.Task.Action == globals.StartScanTask {
		i := "@every " + strconv.Itoa(int(req.Schedule.IntervalSeconds)) + "s"

		c := cron.New()

		cId, err := c.AddFunc(i, func() {
			var ds *datasource_proto.Endpoint
			dbC := db_proto.NewDBClient(globals.DB_SERVICE_NAME, s.Client)

			rsp, err := dbC.Read(globals.NewSystemContext(), &db_proto.ReadRequest{
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
					_, err = dsC.Scan(globals.NewSystemContext(), &datasource_proto.ScanRequest{
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
						// Stop cron task if running
						v.Cron.Stop()
						// Remove cron task
						c.Remove(v.CronId)
						// Delete CronWrapper instance from Scheduler
						s.Crons = append(s.Crons[:k], s.Crons[k+1:]...)
					}
				}
			}
		})
		if err != nil {
			log.Println("ERROR", err)
			return nil, err
		}

		c.Start()
		s.Crons = append(s.Crons, &CronWrapper{
			Id:     req.Task.Id,
			CronId: cId,
			Cron:   c,
		})

		taskRecognise = true
	}

	if !taskRecognise {
		return nil, errors.New("Task not recognise")
	}

	return &proto.CreateScheduledTaskResponse{}, nil
}

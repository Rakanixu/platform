package handler

import (
	"encoding/json"
	"errors"
	"github.com/jasonlvhit/gocron"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"log"
)

//"github.com/jasonlvhit/gocron" It is not good enought. We can not remove a single instance task
//  As is implemented, if we remove the registerScanTask, we would remove all of them
func (s *Scheduler) createTask(ctx context.Context, req *proto.CreateScheduledTaskRequest) (*proto.CreateScheduledTaskResponse, error) {
	taskRecognise := false

	if req.Task.Action == globals.StartScanTask {
		// ctx it is not passed because circuit breaker, timeout will exceed when executed in the future
		gocron.Every(uint64(req.Schedule.IntervalSeconds)).Seconds().Do(registerScanTask, s, req)
		// Non blocking routine
		go func() {
			<-gocron.Start()
		}()
		taskRecognise = true
	}

	if !taskRecognise {
		return nil, errors.New("Task not recognise")
	}

	return &proto.CreateScheduledTaskResponse{}, nil
}

// registerScanTask trigger an scan for the datasource id passed in the request
// It checks the datasource exists on DB, therefore execute it
// If datasources is not found task must be removed from scheduler
// TODO Remove task from sheduler when
func registerScanTask(s *Scheduler, req *proto.CreateScheduledTaskRequest) {
	var ds *datasource_proto.Endpoint
	dbC := db_proto.NewDBClient("", s.Client)

	rsp, err := dbC.Read(context.Background(), &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    req.Task.Id,
	})
	if err != nil {
		// TODO:
		// DB.Read fails if record not found. Here we 'should' have to remove the task from scheduler (NO!)
		// We have to return empty struct when record not found, and work with data afterwards

		//gocron.Remove() --> remove task by function name, this is shit

		log.Println("ERROR", err)
		return
	}

	if err := json.Unmarshal([]byte(rsp.Result), &ds); err != nil {
		log.Println("ERROR", err)
		return
	}

	// Run the scheduled task only if is not already running
	if !ds.CrawlerRunning {
		dsC := datasource_proto.NewDataSourceClient("", s.Client)
		_, err = dsC.Scan(context.Background(), &datasource_proto.ScanRequest{
			Id: req.Task.Id,
		})
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	}
}

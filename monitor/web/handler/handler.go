package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"golang.org/x/net/context"

	monitor "github.com/kazoup/platform/monitor/srv/proto/monitor"
	proto "github.com/micro/go-os/monitor/proto"
)

var (
	opts          *ace.Options
	monitorClient monitor.MonitorClient
)

func Init(dir string, m monitor.MonitorClient) {
	monitorClient = m

	opts = ace.InitializeOptions(nil)
	opts.BaseDir = dir
	opts.DynamicReload = true
	opts.FuncMap = template.FuncMap{
		"TimeAgo": func(t int64) string {
			return timeAgo(t)
		},
		"Colour": func(s string) string {
			return colour(s)
		},
	}
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	basePath := hostPath(r)

	opts.FuncMap["URL"] = func(path string) string {
		return filepath.Join(basePath, path)
	}

	tpl, err := ace.Load("layout", tmpl, opts)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
		return
	}

	if err := tpl.Execute(w, data); err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", 302)
	}
}

// The index page
func Index(w http.ResponseWriter, r *http.Request) {
	rsp, err := monitorClient.Services(context.TODO(), &monitor.ServicesRequest{})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	// sort
	sort.Sort(&sortedServices{rsp.Services})

	type serv struct {
		Name     string
		Nodes    int
		Versions int
	}

	data := make(map[string]serv)

	for _, service := range rsp.Services {
		st, ok := data[service.Name]
		if !ok {
			data[service.Name] = serv{
				Name:     service.Name,
				Nodes:    len(service.Nodes),
				Versions: 1,
			}
			continue
		}

		st.Nodes += len(service.Nodes)
		st.Versions++
		data[service.Name] = st
	}

	// render
	render(w, r, "index", map[string]interface{}{
		"Services": data,
	})
}

func Service(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	if len(service) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	// TODO: limit/offset
	rsp, err := monitorClient.Services(context.TODO(), &monitor.ServicesRequest{
		Service: service,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	hrsp, err := monitorClient.HealthChecks(context.TODO(), &monitor.HealthChecksRequest{
		Service: service,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	srsp, err := monitorClient.Status(context.TODO(), &monitor.StatusRequest{
		Service: service,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	trsp, err := monitorClient.Stats(context.TODO(), &monitor.StatsRequest{
		Service: service,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	type stat struct {
		Name      string
		Status    string
		Info      string
		Id        string
		Cpu       string
		Mem       string
		Ctx       uint64
		Version   string
		Timestamp int64
		Threads   uint64
		Interval  int64
		InBlock   uint64
		OuBlock   uint64
	}

	status := make(map[string]*proto.Status)
	stats := make(map[string][]*stat)

	for _, st := range srsp.Statuses {
		status[st.Service.Nodes[0].Id] = st
	}

	for _, st := range trsp.Stats {
		cpu := float64(st.Cpu.UserTime+st.Cpu.SystemTime) / float64(st.Interval*1e9)

		var nstatus, ninfo string

		if st, ok := status[st.Service.Nodes[0].Id]; ok {
			nstatus = st.Status.String()
			ninfo = st.Info
		}

		stats[st.Service.Version] = append(stats[st.Service.Version], &stat{
			Name:      st.Service.Name,
			Status:    nstatus,
			Info:      ninfo,
			Id:        st.Service.Nodes[0].Id,
			Cpu:       fmt.Sprintf("%.2f %%", cpu*100.0),
			Mem:       fmt.Sprintf("%d mb", st.Memory.MaxRss/uint64(1024*1024)),
			Ctx:       st.Cpu.InvCtxSwitch + st.Cpu.VolCtxSwitch,
			InBlock:   st.Disk.InBlock,
			OuBlock:   st.Disk.OuBlock,
			Threads:   st.Runtime.NumThreads,
			Timestamp: st.Timestamp,
			Interval:  st.Interval,
		})
	}

	// render
	render(w, r, "service", map[string]interface{}{
		"Name":         service,
		"Services":     rsp.Services,
		"Status":       srsp.Statuses,
		"Stats":        stats,
		"Healthchecks": hrsp.Healthchecks,
	})
}

func ServiceStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	if len(service) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	rsp, err := monitorClient.Status(context.TODO(), &monitor.StatusRequest{
		Service: service,
		Verbose: true,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	// sort
	sort.Sort(&sortedStatuses{rsp.Statuses})

	versions := make(map[string][]*proto.Status)

	for _, status := range rsp.Statuses {
		st, ok := versions[status.Service.Version]
		if !ok {
			versions[status.Service.Version] = []*proto.Status{status}
			continue
		}
		st = append(st, status)
		versions[status.Service.Version] = st
	}

	// render
	render(w, r, "serviceStatus", map[string]interface{}{
		"Name":   service,
		"Status": versions,
	})
}

func ServiceStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]

	if len(service) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	// TODO: limit/offset
	rsp, err := monitorClient.Stats(context.TODO(), &monitor.StatsRequest{
		Service: service,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	type stat struct {
		Name      string
		Id        string
		Cpu       string
		Mem       string
		Ctx       uint64
		Version   string
		Timestamp int64
		Threads   uint64
		Interval  int64
		InBlock   uint64
		OuBlock   uint64
	}

	stats := make(map[string]*stat)

	for _, st := range rsp.Stats {
		cpu := float64(st.Cpu.UserTime+st.Cpu.SystemTime) / float64(st.Interval*1e9)

		stats[st.Service.Nodes[0].Id] = &stat{
			Name:      st.Service.Name,
			Id:        st.Service.Nodes[0].Id,
			Cpu:       fmt.Sprintf("%.2f %%", cpu*100.0),
			Mem:       fmt.Sprintf("%d mb", st.Memory.MaxRss/uint64(1024*1024)),
			Version:   st.Service.Version,
			Ctx:       st.Cpu.InvCtxSwitch + st.Cpu.VolCtxSwitch,
			InBlock:   st.Disk.InBlock,
			OuBlock:   st.Disk.OuBlock,
			Threads:   st.Runtime.NumThreads,
			Timestamp: st.Timestamp,
			Interval:  st.Interval,
		}
	}

	render(w, r, "serviceStats", map[string]interface{}{
		"Name":  service,
		"Stats": stats,
	})
}

func Stats(w http.ResponseWriter, r *http.Request) {
	// TODO: limit/offset
	rsp, err := monitorClient.Stats(context.TODO(), &monitor.StatsRequest{})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	type stat struct {
		Name  string
		Cpu   string
		Mem   string
		Nodes int

		cpu      float64
		mem      uint64
		interval int64
	}

	stats := make(map[string]*stat)

	for _, st := range rsp.Stats {
		estat, ok := stats[st.Service.Name]
		if !ok {
			cpu := float64(st.Cpu.UserTime+st.Cpu.SystemTime) / float64(st.Interval*1e9)

			stats[st.Service.Name] = &stat{
				Name:     st.Service.Name,
				Nodes:    1,
				Cpu:      fmt.Sprintf("%.2f %%", cpu*100.0),
				Mem:      fmt.Sprintf("%d mb", st.Memory.MaxRss/uint64(1024*1024)),
				cpu:      float64(st.Cpu.UserTime + st.Cpu.SystemTime),
				mem:      st.Memory.MaxRss,
				interval: st.Interval,
			}
			continue
		}

		// update stats
		estat.Nodes++
		estat.mem += st.Memory.MaxRss
		estat.cpu += float64(st.Cpu.UserTime + st.Cpu.SystemTime)

		// update string vals
		estat.Mem = fmt.Sprintf("%d mb", estat.mem/uint64(1024*1024))
		estat.Cpu = fmt.Sprintf("%.2f %%", (estat.cpu / float64(estat.interval*1e9) * 100.0))
	}

	render(w, r, "stats", map[string]interface{}{
		"Stats": stats,
	})
}

func Status(w http.ResponseWriter, r *http.Request) {
	rsp, err := monitorClient.Status(context.TODO(), &monitor.StatusRequest{})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	// sort
	sort.Sort(&sortedStatuses{rsp.Statuses})

	type stat struct {
		Name     string
		Versions []*proto.Status
	}

	services := make(map[string]*stat)

	for _, status := range rsp.Statuses {
		st, ok := services[status.Service.Name]
		if !ok {
			services[status.Service.Name] = &stat{
				Name:     status.Service.Name,
				Versions: []*proto.Status{status},
			}
			continue
		}
		st.Versions = append(st.Versions, status)
		services[status.Service.Name] = st
	}

	// render
	render(w, r, "status", map[string]interface{}{
		"Status": services,
	})
}

func Healthchecks(w http.ResponseWriter, r *http.Request) {
	// TODO: limit/offset
	rsp, err := monitorClient.HealthChecks(context.TODO(), &monitor.HealthChecksRequest{})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	healthchecks := make(map[string]*proto.HealthCheck)

	for _, check := range rsp.Healthchecks {
		echeck, ok := healthchecks[check.Id]
		if !ok {
			healthchecks[check.Id] = check
			continue
		}
		if len(check.Error) > 0 && (len(echeck.Error) == 0 || check.Timestamp > echeck.Timestamp) {
			healthchecks[check.Id] = check

		} else if check.Timestamp > echeck.Timestamp {
			healthchecks[check.Id] = check
		}
	}

	render(w, r, "healthchecks", map[string]interface{}{
		"Healthchecks": healthchecks,
	})
}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		http.Redirect(w, r, "/", 302)
		return
	}

	// TODO: limit/offset
	rsp, err := monitorClient.HealthChecks(context.TODO(), &monitor.HealthChecksRequest{
		Id: id,
	})
	if err != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	render(w, r, "healthcheck", map[string]interface{}{
		"Id":           id,
		"Healthchecks": rsp.Healthchecks,
	})
}

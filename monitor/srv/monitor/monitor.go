package monitor

import (
	"errors"
	"fmt"
	"sync"
	"time"

	proto "github.com/micro/go-os/monitor/proto"
	"golang.org/x/net/context"
)

type monitor struct {
	sync.RWMutex
	healthChecks map[string][]*proto.HealthCheck
	services     map[string]*proto.Service
	status       map[string]*proto.Status
	stats        map[string]*proto.Stats
}

var (
	DefaultMonitor   = newMonitor()
	ErrNotFound      = errors.New("not found")
	HealthCheckTopic = "micro.monitor.healthcheck"
	StatusTopic      = "micro.monitor.status"
	StatsTopic       = "micro.monitor.stats"
	TickInterval     = time.Duration(time.Minute)
)

func newMonitor() *monitor {
	return &monitor{
		healthChecks: make(map[string][]*proto.HealthCheck),
		services:     make(map[string]*proto.Service),
		status:       make(map[string]*proto.Status),
		stats:        make(map[string]*proto.Stats),
	}
}

func filter(hc []*proto.HealthCheck, status proto.HealthCheck_Status, limit, offset int) []*proto.HealthCheck {
	var hcs []*proto.HealthCheck

	if len(hc) < offset {
		return hcs
	}

	if (limit + offset) > len(hc) {
		limit = len(hc) - offset
	}

	for i := 0; i < limit; i++ {
		if status == proto.HealthCheck_UNKNOWN {
			hcs = append(hcs, hc[offset])
		} else if hc[offset].Status == status {
			hcs = append(hcs, hc[offset])
		}
		offset++
	}

	return hcs
}

func filterStats(st []*proto.Stats, limit, offset int) []*proto.Stats {
	var stats []*proto.Stats

	if len(st) < offset {
		return stats
	}

	if (limit + offset) > len(st) {
		limit = len(st) - offset
	}

	for i := 0; i < limit; i++ {
		stats = append(stats, st[offset])
		offset++
	}

	return stats
}

func filterStatus(st []*proto.Status, limit, offset int) []*proto.Status {
	var status []*proto.Status

	if len(st) < offset {
		return status
	}

	if (limit + offset) > len(st) {
		limit = len(st) - offset
	}

	for i := 0; i < limit; i++ {
		status = append(status, st[offset])
		offset++
	}

	return status
}

func (m *monitor) reap() {
	m.Lock()
	defer m.Unlock()

	t := time.Now().Unix()

	services := make(map[string]*proto.Service)

	// reap healthchecks
	for id, hc := range m.healthChecks {
		var checks []*proto.HealthCheck
		for _, check := range hc {
			if t > (check.Timestamp+check.Interval) && t > (check.Timestamp+check.Ttl) {
				continue
			}
			checks = append(checks, check)

			// create new service list

			// TODO: maybe hold onto it so we have history
			if check.Service != nil && len(check.Service.Name) > 0 {
				if len(check.Service.Nodes) > 0 {
					services[check.Service.Nodes[0].Id] = check.Service
				}
			}
		}
		m.healthChecks[id] = checks
	}

	// reap status
	for id, status := range m.status {
		// expired
		if t > (status.Timestamp+status.Interval) && t > (status.Timestamp+status.Ttl) {
			delete(m.status, id)
			continue
		}

		// past interval
		if d := t - (status.Timestamp + status.Interval); d > 0 {
			status.Status = proto.Status_UNKNOWN
			status.Info = fmt.Sprintf("Last update %v ago", time.Duration(d)*time.Second)
		}

		// incase its not seen or something
		services[status.Service.Nodes[0].Id] = status.Service
	}

	// reap stats
	for id, stats := range m.stats {
		// expired
		if t > (stats.Timestamp+stats.Interval) && t > (stats.Timestamp+stats.Ttl) {
			delete(m.stats, id)
			continue
		}

		// incase its not seen or something
		services[stats.Service.Nodes[0].Id] = stats.Service
	}

	m.services = services
}

func (m *monitor) run() {
	t := time.NewTicker(TickInterval)

	for _ = range t.C {
		m.reap()
	}
}

func (m *monitor) HealthChecks(service, id string, status proto.HealthCheck_Status, limit, offset int) ([]*proto.HealthCheck, error) {
	m.RLock()
	defer m.RUnlock()

	if len(service) == 0 && len(id) == 0 {
		var hcs []*proto.HealthCheck
		for _, hc := range m.healthChecks {
			hcs = append(hcs, hc...)
		}
		return filter(hcs, status, limit, offset), nil
	}

	if len(service) > 0 {
		l := len(id)

		var hcs []*proto.HealthCheck
		for _, hc := range m.healthChecks {
			for _, ihc := range hc {
				if ihc.Service.Name != service {
					continue
				}
				if l > 0 && ihc.Id != id {
					continue
				}
				hcs = append(hcs, ihc)
			}
		}
		return filter(hcs, status, limit, offset), nil
	}

	hcs, ok := m.healthChecks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return filter(hcs, status, limit, offset), nil
}

func (m *monitor) Services(s string) ([]*proto.Service, error) {
	m.RLock()
	defer m.RUnlock()

	toService := make(map[string][]*proto.Service)

	for _, service := range m.services {
		if len(s) > 0 && service.Name != s {
			continue
		}

		cp := &proto.Service{}

		*cp = *service

		sr, ok := toService[cp.Name]
		if !ok {
			toService[cp.Name] = []*proto.Service{cp}
			continue
		}

		// insert nodes into service version
		var seen bool
		for _, srv := range sr {
			if srv.Version == cp.Version {
				srv.Nodes = append(srv.Nodes, cp.Nodes...)
				seen = true
				break
			}
		}
		if !seen {
			toService[cp.Name] = append(toService[cp.Name], cp)
		}
	}

	var services []*proto.Service
	for _, service := range toService {
		services = append(services, service...)
	}
	return services, nil
}

func (m *monitor) Stats(service, id string, limit, offset int) ([]*proto.Stats, error) {
	m.RLock()
	defer m.RUnlock()

	var stats []*proto.Stats

	// service node
	if len(service) > 0 && len(id) > 0 {
		stat, ok := m.stats[id]
		if !ok || stat.Service.Name != service {
			return nil, ErrNotFound
		}
		return []*proto.Stats{stat}, nil
	}

	// return service stats
	if len(service) > 0 && len(id) == 0 {
		for _, stat := range m.stats {
			if stat.Service.Name != service {
				continue
			}
			stats = append(stats, stat)
		}
		return filterStats(stats, limit, offset), nil
	}

	// single node
	if len(service) == 0 && len(id) > 0 {
		if stat, ok := m.stats[id]; ok {
			return []*proto.Stats{stat}, nil
		}
		return nil, ErrNotFound
	}

	// all services
	for _, stat := range m.stats {
		stats = append(stats, stat)
	}

	return filterStats(stats, limit, offset), nil
}
func (m *monitor) Status(service, id string, limit, offset int, verbose bool) ([]*proto.Status, error) {
	m.RLock()
	defer m.RUnlock()

	// return single sevice node
	if len(service) > 0 && len(id) > 0 {
		stat, ok := m.status[id]
		if !ok || stat.Service.Name != service {
			return nil, ErrNotFound
		}
		return []*proto.Status{stat}, nil
	}

	// return single node regardless of service
	if len(service) == 0 && len(id) > 0 {
		if stat, ok := m.status[id]; ok {
			return []*proto.Status{stat}, nil
		}
		return nil, ErrNotFound
	}

	mstat := make(map[string]*proto.Status)

	toSlice := func(statusMap map[string]*proto.Status) ([]*proto.Status, error) {
		var status []*proto.Status
		for _, stat := range statusMap {
			status = append(status, stat)
		}
		return filterStatus(status, limit, offset), nil
	}

	updateMap := func(status *proto.Status) {
		st, ok := mstat[status.Service.Name+status.Service.Version]
		if !ok {
			stat := &proto.Status{
				Service: &proto.Service{},
				Status:  status.Status,
				Info:    status.Info,
			}
			*stat.Service = *status.Service
			mstat[status.Service.Name+status.Service.Version] = stat
			return
		}

		// append nodes
		st.Service.Nodes = append(st.Service.Nodes, status.Service.Nodes...)
	}

	// return all service nodes
	if len(service) > 0 {
		// break out nodes
		if verbose {
			var statuses []*proto.Status
			for _, status := range m.status {
				if status.Service.Name != service {
					continue
				}
				statuses = append(statuses, status)
			}
			return filterStatus(statuses, limit, offset), nil
		}

		// less verbose
		for _, status := range m.status {
			if status.Service.Name != service {
				continue
			}
			updateMap(status)
		}

		return toSlice(mstat)
	}

	// service and id is blank... return all services

	// break out nodes
	if verbose {
		return toSlice(m.status)
	}

	// return less verbose
	for _, status := range m.status {
		updateMap(status)
	}

	return toSlice(mstat)
}

func (m *monitor) ProcessHealthCheck(ctx context.Context, hc *proto.HealthCheck) error {
	m.Lock()
	defer m.Unlock()

	// Add to status or services if we don't have them
	if hc.Service != nil && len(hc.Service.Name) > 0 {
		if len(hc.Service.Nodes) > 0 {
			// add to service list
			m.services[hc.Service.Nodes[0].Id] = hc.Service

			// add to status list
			if _, ok := m.status[hc.Service.Nodes[0].Id]; !ok {
				m.status[hc.Service.Nodes[0].Id] = &proto.Status{
					Service: hc.Service,
					Status:  proto.Status_UNKNOWN,
				}
			}
		}
	}

	hcs, ok := m.healthChecks[hc.Id]
	if !ok {
		m.healthChecks[hc.Id] = append(m.healthChecks[hc.Id], hc)
		return nil
	}

	for i, h := range hcs {
		if len(hc.Service.Nodes) > 0 && (h.Service.Nodes[0].Id == hc.Service.Nodes[0].Id) {
			hcs[i] = hc
			m.healthChecks[hc.Id] = hcs
			return nil
		}
	}

	hcs = append(hcs, hc)
	m.healthChecks[hc.Id] = hcs
	return nil
}

func (m *monitor) ProcessStatus(ctx context.Context, st *proto.Status) error {
	m.Lock()
	defer m.Unlock()

	// add to service list
	m.services[st.Service.Nodes[0].Id] = st.Service
	// add status
	m.status[st.Service.Nodes[0].Id] = st

	return nil
}

func (m *monitor) ProcessStats(ctx context.Context, st *proto.Stats) error {
	m.Lock()
	defer m.Unlock()

	// add to service list
	m.services[st.Service.Nodes[0].Id] = st.Service
	// add stats
	m.stats[st.Service.Nodes[0].Id] = st

	return nil
}

func (m *monitor) Run() {
	go m.run()
}

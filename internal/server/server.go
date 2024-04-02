package server

import (
	"github.com/sadbanzai/sysmon/internal/collector"
	pb "github.com/sadbanzai/sysmon/internal/pb"
	"log"
	"sync"
	"time"
)

type Server struct {
	pb.UnimplementedSysmonServer
	//Collectors           collector.Collectors
	CollectorLoadAverage collector.Collector
	CollectorCpuStat     collector.Collector
}

func (s *Server) StreamStats(req *pb.StatsRequest, stream pb.Sysmon_StreamStatsServer) error {
	for {
		stats, err := s.GetStats(int(req.GetAverageM()))
		if err != nil {
			return err
		}

		if err := stream.Send(stats); err != nil {
			return err
		}

		time.Sleep(time.Duration(req.GetIntervalN()) * time.Second)
	}
}

func (s *Server) GetStats(m int) (*pb.StatsResponse, error) {
	stats := make(map[string]collector.Items)
	var err error
	stats[s.CollectorLoadAverage.Name], err = s.CollectorLoadAverage.Average(m)
	if err != nil {
		log.Printf("%s", err)
	}
	stats[s.CollectorCpuStat.Name], err = s.CollectorCpuStat.Average(m)
	if err != nil {
		log.Printf("%s", err)
	}
	cpuStat := &pb.CpuStat{User: 0, System: 0, Idle: 0}
	loadAverage := &pb.LoadAverage{OneMinute: 0, FiveMinutes: 0, FifteenMinutes: 0}

	loadAverage.OneMinute = stats[s.CollectorLoadAverage.Name]["avg1"]
	loadAverage.FiveMinutes = stats[s.CollectorLoadAverage.Name]["avg5"]
	loadAverage.FifteenMinutes = stats[s.CollectorLoadAverage.Name]["avg15"]

	cpuStat.User = stats[s.CollectorCpuStat.Name]["user"]
	cpuStat.System = stats[s.CollectorCpuStat.Name]["system"]
	cpuStat.Idle = stats[s.CollectorCpuStat.Name]["idle"]

	response := &pb.StatsResponse{
		CpuStat:     cpuStat,
		LoadAverage: loadAverage,
	}

	return response, nil
}

func (s *Server) Start() {
	doneCh := make(chan struct{})
	wg := sync.WaitGroup{}
	if s.CollectorLoadAverage.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.CollectorLoadAverage.CollectData(doneCh)
		}()
	}
	if s.CollectorCpuStat.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.CollectorCpuStat.CollectData(doneCh)
		}()
	}
}

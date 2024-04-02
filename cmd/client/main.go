package main

import (
	"context"
	"fmt"
	"github.com/gosuri/uilive"
	pb "github.com/sadbanzai/sysmon/internal/pb"
	"google.golang.org/grpc"
	"log"
	"text/tabwriter"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSysmonClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.StreamStats(ctx, &pb.StatsRequest{IntervalN: 10, AverageM: 2})
	if err != nil {
		log.Fatalf("Failed to get load average: %v", err)
	}

	writer := uilive.New()
	writer.Start()

	w := tabwriter.NewWriter(writer, 0, 0, 1, ' ', tabwriter.TabIndent)

	for {
		stats, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive load average: %v", err)
		}

		fmt.Fprint(w, "\033[H\033[2J")

		fmt.Fprintf(w, "%s %.2f, %.2f, %.2f\n",
			"load average:",
			stats.LoadAverage.OneMinute,
			stats.LoadAverage.FiveMinutes,
			stats.LoadAverage.FifteenMinutes)
		fmt.Fprintf(w, "%s %.2f us, %.2f sy, %.2f id\n",
			"cpu stats:",
			stats.CpuStat.User,
			stats.CpuStat.System,
			stats.CpuStat.Idle)
		w.Flush()
	}
}

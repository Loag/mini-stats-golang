package client

import (
	"context"

	"time"

	pb "github.com/Loag/mini-stats-proto/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Instance struct {
	Name        string
	MessageType string
	Value       float64
	Time        int64
}

func (i *Instance) ToMsg() *pb.IngestRequest {

	return &pb.IngestRequest{
		Name:       i.Name,
		Value:      i.Value,
		Time:       uint64(i.Time),
		MetricType: get_metric_type(i.MessageType),
	}
}

func get_metric_type(input string) pb.MetricType {
	switch input {
	case "GAUGE":
		return pb.MetricType_GAUGE
	case "COUNTER":
		return pb.MetricType_COUNTER
	default:
		return pb.MetricType_UNSPECIFIED
	}
}

type Metric interface {
	GetValue() Instance
}

type MiniStatsClientOptions struct {
	Debug    bool
	ApiKey   string
	Endpoint string
	Interval int
}

type MiniStatsClient struct {
	client   pb.IngestServiceClient
	metrics  []Metric
	debug    bool
	key      string
	endpoint string
	interval int
}

func New(opts MiniStatsClientOptions) *MiniStatsClient {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if opts.Debug {
		log.Info().Interface("dict", opts).Msg("opts set")
	}

	return &MiniStatsClient{
		debug:    opts.Debug,
		key:      opts.ApiKey,
		endpoint: opts.Endpoint,
		interval: opts.Interval,
	}
}

func (m *MiniStatsClient) AddMetric(metric Metric) *MiniStatsClient {
	m.metrics = append(m.metrics, metric)
	return m
}

func (m *MiniStatsClient) Start() {

	if m.debug {
		log.Debug().Msgf("starting client at: %s", m.endpoint)
	}

	conn, err := grpc.NewClient(m.endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Err(err).Msg("failed to setup connection to ministats")
	}
	defer conn.Close()
	c := pb.NewIngestServiceClient(conn)

	ticker := time.NewTicker(time.Duration(m.interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if m.debug {
			log.Info().Msg("running cycle")
		}

		for _, metric := range m.metrics {

			if m.debug {
				log.Info().Msg("iterating set metrics")
			}

			instance := metric.GetValue()

			_, err := c.Ingest(context.Background(), instance.ToMsg())

			if err != nil {
				log.Err(err).Msgf("unable to metric with name: %s", instance.Name)
			}
		}
	}
}

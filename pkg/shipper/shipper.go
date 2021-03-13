package shipper

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/conprof/conprof/pkg/store/storepb"
	"github.com/prometheus/prometheus/pkg/timestamp"
	labelpb "github.com/thanos-io/thanos/pkg/store/labelpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type perRequestBearerToken struct {
	token    string
	insecure bool
}

func (t *perRequestBearerToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t *perRequestBearerToken) RequireTransportSecurity() bool {
	return !t.insecure
}

// Shipper is a wrapper around a WritableProfileStoreClient
type Shipper struct {
	client        storepb.WritableProfileStoreClient
	defaultLabels []labelpb.Label
}

// Options for the Shipper
type Options struct {
	BearerToken   string
	Insecure      bool
	DefaultLabels []labelpb.Label
}

// Ship a profile
func (s *Shipper) Ship(ctx context.Context, profile []byte, labels ...labelpb.Label) error {
	_, err := s.client.Write(ctx, &storepb.WriteRequest{
		ProfileSeries: []storepb.ProfileSeries{
			{
				Labels: append(s.defaultLabels, labels...),
				Samples: []storepb.Sample{
					{
						Value:     profile,
						Timestamp: timestamp.FromTime(time.Now()),
					},
				},
			},
		},
	})
	return err
}

// NewShipper returns a new shipper to write profiles to remote storage
func NewShipper(upstreamAddress string, o *Options) (*Shipper, error) {
	opts := []grpc.DialOption{}
	if o.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	if o.BearerToken != "" {
		opts = append(opts, grpc.WithPerRPCCredentials(&perRequestBearerToken{
			token:    o.BearerToken,
			insecure: o.Insecure,
		}))
	}

	conn, err := grpc.Dial(upstreamAddress, opts...)
	if err != nil {
		return nil, err
	}

	labels := []labelpb.Label{
		{
			Name:  "__name__",
			Value: "perf",
		},
	}
	if o.DefaultLabels != nil {
		labels = append(labels, o.DefaultLabels...)
	}

	return &Shipper{
		client:        storepb.NewWritableProfileStoreClient(conn),
		defaultLabels: labels,
	}, nil
}

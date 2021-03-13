module github.com/conprof/perfessor

go 1.15

require (
	github.com/conprof/conprof v0.0.0-20210311120814-17eb689b9725
	github.com/go-kit/kit v0.10.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/thanos-io/thanos v0.18.0
	google.golang.org/grpc v1.34.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.8.2-0.20200922180708-b0145884d381

# perfessor - Continuous Profiling Sidecar


### About

Perfessor is a continuous profiling agent that can profile running programs using [perf](http://www.brendangregg.com/perf.html)

It then converts those profiles from perf format into [pprof](https://github.com/google/pprof) format, and ships those profiles 
to a supported backend such as Conprof.

### Install

`go get github.com/conprof/perfessor`

### Build

`make docker`

### Run

`docker run --rm -it --cap-add SYS_ADMIN perfessor:latest`

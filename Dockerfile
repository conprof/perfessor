FROM python:3-slim-stretch
FROM golang:1.16.2-stretch

# install perf
RUN apt-get update && \
    apt-get install -y --no-install-recommends linux-perf && \
    ln -s /usr/bin/perf_* /usr/bin/perf && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# build pprof
RUN go get -u github.com/google/pprof

# build perf data converter
RUN git clone https://github.com/google/perf_data_converter.git
RUN apt update
RUN apt install -y g++ libelf-dev libcap-dev
RUN apt install -y curl gnupg
RUN curl -fsSL https://bazel.build/bazel-release.pub.gpg | gpg --dearmor > bazel.gpg
RUN mv bazel.gpg /etc/apt/trusted.gpg.d/
RUN echo "deb [arch=amd64] https://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list
RUN apt update && apt install -y bazel
WORKDIR perf_data_converter
RUN bazel build src:perf_to_profile
RUN mv bazel-bin/src/perf_to_profile /bin

# copy in perfessor binary
COPY perfessor /bin

CMD ["/bin/perfessor"]

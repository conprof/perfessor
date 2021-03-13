FROM ubuntu:groovy
RUN apt update
RUN apt install -y linux-tools-common linux-tools-generic wget flex bison make
RUN wget http://archive.ubuntu.com/ubuntu/pool/main/l/linux/linux_5.4.0.orig.tar.gz
RUN tar -xzf linux_5.4.0.orig.tar.gz 
WORKDIR linux-5.4/tools/perf
RUN make -C .
RUN make install

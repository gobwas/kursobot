# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /root/kursobot

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/constabulary/gb/...
RUN cd /root/outyet && gb vendor restore
RUN cd /root/outyet && gb build all

# Run the outyet command by default when the container starts.
ENTRYPOINT /root/outyet/bin/kursobot -config=/usr/local/etc/kursobot.conf

# Document that the service listens on port 8080.
EXPOSE 8443
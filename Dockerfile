FROM golang
ENV IMAGR_PASSWORD="password"
RUN go get -u github.com/groob/imagr-server
ENTRYPOINT ["imagr-server", "/imagr_repo"]

VOLUME ["/imagr_repo"]
EXPOSE 3000


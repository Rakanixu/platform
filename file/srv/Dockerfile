FROM alpine:3.4
# Let's roll
RUN	apk update && \
	apk upgrade && \
	apk add --update tzdata && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

FROM scratch
ADD srv /
ENTRYPOINT [ "/srv" ]

FROM alpine:3.4
# Let's roll
RUN	apk update && \
	apk upgrade && \
	apk add --update tzdata && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

ADD doc_templates /doc_templates
ADD srv /
ENTRYPOINT [ "/srv" ]

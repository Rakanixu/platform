FROM alpine:3.4
# Let's roll
RUN	apk update && \
	apk upgrade && \
	apk add --update tzdata
ENTRYPOINT [ "/web" ]

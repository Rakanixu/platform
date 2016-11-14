FROM alpine:3.4
# Let's roll
RUN	apk update && \
	apk upgrade && \
	apk add --update tzdata
ADD src src
ADD ui /
ENTRYPOINT [ "/ui" ]



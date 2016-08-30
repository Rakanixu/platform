FROM alpine:3.4
# Let's roll
RUN	apk update && \
	apk upgrade && \
	apk add --update tzdata
ADD app app
ADD web /
ENTRYPOINT [ "/web" ]



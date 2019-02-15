FROM 1.11.5-alpine3.9 as builder

ADD . /app

# TODO: add build stage

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["logour"]

FROM redis:7.2.5-alpine

ARG REDIS_PORT

WORKDIR /redis-workdir

RUN apk update
RUN apk --no-cache add wget envsubst

RUN wget https://raw.githubusercontent.com/redis/redis/7.2.5/redis.conf
RUN mv redis.conf redis.default.conf

COPY redis.conf .

ENV REDIS_PORT $REDIS_PORT
RUN envsubst < redis.conf > updated_redis.conf
RUN mv updated_redis.conf redis.conf

EXPOSE $REDIS_PORT

CMD redis-server ./redis.conf
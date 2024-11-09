FROM postgres:14

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=12345
ENV POSTGRES_DB=campaign-consumer-api

EXPOSE 5432:5432

ARG LOCAL_SQL=local/sql
ARG MIGRATIONS_SQL=migrations/sql

COPY $LOCAL_SQL /docker-entrypoint-initdb.d/
COPY $MIGRATIONS_SQL /docker-entrypoint-initdb.d/
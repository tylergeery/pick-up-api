FROM postgres:9.6.0

ENV POSTGRES_USER 'raccoon'
ENV POSTGRES_DB 'pickup'

COPY psql_dump.sql /docker-entrypoint-initdb.d/

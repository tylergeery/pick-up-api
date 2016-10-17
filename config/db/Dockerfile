FROM postgres:9.6.0

ENV POSTGRES_USER 'raccoon'

COPY psql_dump.sql /docker-entrypoint-initdb.d/

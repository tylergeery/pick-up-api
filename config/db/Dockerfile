FROM postgres:9.6.0

ENV POSTGRES_USER 'raccoon'
ENV POSTGRES_DB 'pickup'

# Copies all seed files to image for initialization
COPY seed /docker-entrypoint-initdb.d/

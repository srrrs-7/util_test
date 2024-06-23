FROM postgres:16.1

COPY ./postgres /postgres

RUN apt-get update && \
    apt-get install -y git make gcc postgresql-server-dev-16

RUN cd /tmp && \
    git clone https://github.com/pgvector/pgvector.git && \
    cd pgvector && \
    make && \
    make install && \
    cd ../ && rm -rf pgvector
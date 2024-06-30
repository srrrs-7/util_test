FROM ubuntu

RUN apt-get -y update && \
    apt-get -y install curl

RUN curl --proto '=https' --tlsv1.2 -sSf https://install.surrealdb.com | sh -s -- --nightly

COPY ./surrealdb /surrealdb

ENTRYPOINT surreal import --conn $SURREALDB_ENDPOINT --user $SURREALDB_USER --pass $SURREALDB_PASS --ns test --db test /surrealdb/init.surql && \
    surreal sql -e $SURREALDB_ENDPOINT --user $SURREALDB_USER --pass $SURREALDB_PASS --ns test --db test
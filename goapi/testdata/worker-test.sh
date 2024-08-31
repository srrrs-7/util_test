#!/bin/sh

for i in {1..500};
do
    curl --location 'http://localhost:8080/domain/v1/user/1/create' \
    --header 'Content-Type: application/json' \
    --data '{
        "userId":"1"
    }'
done
#!/bin/sh

# health check
curl -XGET 'localhost:9200/'

# create index
curl -X POST "http://localhost:9200/books" \
     -H 'Content-Type: application/json' \
     -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "english"
      },
      "author": {
        "type": "keyword" 
      },
      "published_year": {
        "type": "date"
      },
      "price": {
        "type": "double"
      }
    }
  }
}'
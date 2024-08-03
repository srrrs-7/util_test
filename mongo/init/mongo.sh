#!/bin/sh



mongosh --username root --password root << EOF
db.createCollection("test_collection",　{　capped: true,　size: 1000000,　max: 1000　});

use test_collection;

db.test_collection.insertMany([{ id: 1, name: "Jane Doe", age: 25 },　{ id: 2, name: "Bob Smith", age: 40 }]);

db.test_collection.findOne({ age: {\$gt: 30} })
EOF

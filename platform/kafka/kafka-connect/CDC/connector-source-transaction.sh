#!/bin/bash

docker compose exec debezium curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '
 {
  "name": "mysql-connector",
  "config": {
    "connector.class": "io.debezium.connector.mysql.MySqlConnector",
    "tasks.max": "1",
    "database.hostname": "host.docker.internal",
    "database.port": "4001",
    "database.user": "root",
    "database.password": "root_password",
    "database.server.name": "transaction_db",
    "topic.prefix": "trx",
    "database.include.list": "transaction",
    "schema.history.internal.kafka.bootstrap.servers": "redpanda-0:9092,redpanda-1:9092,redpanda-2:9092",
    "schema.history.internal.kafka.topic": "schemahistory.dbz.transaction",
    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter.schemas.enable": "false",
    "decimal.handling.mode": "string",
    "database.server.id": "123457"
  }
}'
#!/bin/bash

docker compose exec debezium curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '
{
  "name": "mysql-mail-connector",
  "config": {
    "connector.class": "io.debezium.connector.jdbc.JdbcSinkConnector",
    "tasks.max": "1",  // Number of tasks for parallelism
    "topics": "dbz.user_management.users",  // Kafka topic name that Debezium is writing to
    "connection.url": "jdbc:mysql://host.docker.internal:4002/mail",
    "connection.username": "root",
    "connection.password": "root_password",
    "insert.mode": "upsert",  // Use insert if you are only inserting new data
    "schema.evolution": "basic",
    "primary.key.fields": "id",  // Primary key field in the target table
    "primary.key.mode": "record_key",  // Use the record key for primary key in MySQL
    "database.time_zone": "UTC",
    "delete.enabled": "true"
  }
}'
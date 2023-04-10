# Setup manually
Run `docker-compose up` in **setup** folder to create 3 shards with 3 replicas ring technology

# Clickhouse Cli
Best tool to manage clickhouse cluster

## Cli
To Create/Delete tables in Clickhouse database

Host: ``
### API

The api are:
```
    /api/tables/create      Create database and tables with full requirement engine, cluster, sharding,...
    /api/tables/drop        Drop table in database
    /api/column/add         Add column to designated table both distributed and replicas
    /api/column/drop        Drop column on designated table
```

For body request of each api, check in `example` folder

## Airflow
To manipulate data on Clickhouse 

### InsertData

Usage: `./out/airflow -action=InsertData`

The arguments are:
```
    -table          Description fields/datatypes in table 
    -query          Data to insert 
```

For file examples, can check in `example` folder

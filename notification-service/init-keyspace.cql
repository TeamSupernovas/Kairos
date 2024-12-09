-- init-keyspace.cql
CREATE KEYSPACE IF NOT EXISTS notifications 
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

CREATE TABLE IF NOT EXISTS feed_logs(
    id SERIAL PRIMARY KEY,
    client_id VARCHAR (60) NOT NULL,
    portions SMALLINT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    CONSTRAINT fk_feeder
      FOREIGN KEY(client_id) 
	  REFERENCES feeders(client_id)
);
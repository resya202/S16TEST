CREATE TABLE delegation_hourly (
  id SERIAL PRIMARY KEY,
  validator_addr TEXT NOT NULL,
  delegator_addr TEXT NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
  amount_uatom BIGINT NOT NULL,
  change_uatom BIGINT NOT NULL
);

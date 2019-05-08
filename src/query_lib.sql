/*
  need to create the trigger function to extract peices of information from the
  json format
*/
WITH dt_json AS (
  SELECT json_array_elements(response -> 'messages') AS message
  FROM response_raw
  limit 10
)
SELECT *
FROM json_to_recordset(SELECT message FROM dt_json) AS x(id BIGINT, body TEXT)

-- get message body with sentiment
SELECT msg.id, msg.body, msg.entities -> 'sentiment' -> 'basic'
FROM response_raw,
json_to_recordset(response_raw.response -> 'messages') AS msg(id BIGINT, body TEXT, entities JSON);

-- add symbols to above query
-- SELECT msg.id, msg.body, msg.entities -> 'sentiment' -> 'basic'
-- FROM response_raw,
-- json_to_recordset(response_raw.response -> 'messages') AS msg(id BIGINT, body TEXT, entities JSON);

/*
  get the data where new key came through

docker exec -it stocks_db psql -U postgres -d stocks -c "\copy (
  select json_array_elements(response -> 'messages')
  from symbolstream.response_raw
  where id = 'c4aeb47e-e8f9-4e32-a749-4196c410f195'
) TO STDOUT" > response_data.json

*/

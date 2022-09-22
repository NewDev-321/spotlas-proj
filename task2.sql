
/******************************************************************************

                      Task 2:  Write a Query
Create an endpoint which returns spots in a circle or square area. This task must be
completed in Golang.
1.Endpoint should receive 4 parameters
	‣ Latitude
	‣ Longitude
	‣ Radius (in meters)
	‣ Type (circle or square)
2.Find all spots in the table (spots.sql) using the received parameters.
3.Order results by distance.
	‣ If distance between two spots is smaller than 50m, then order by rating.
4.Endpoint should return an array of objects containing all fields in the data set.

*******************************************************************************/


CREATE OR REPLACE FUNCTION get_spots(lat float8, long float8, radius float8, shape varchar) RETURNS TABLE (
	id uuid,
		name varchar,
		website varchar,
		coordinates geography,
		description varchar,
		rating float8,
	distance float8
) AS $$
BEGIN
	RETURN QUERY SELECT
			*, ST_Distance("MY_TABLE".coordinates, CONCAT('POINT(',long::varchar(255),' ',lat::varchar(255),')')::geography) as distance
	FROM "MY_TABLE"
WHERE
	CASE
	WHEN shape = 'square' THEN
		(ST_X("MY_TABLE".coordinates::geometry) BETWEEN (long - (radius/111259.54243971)) AND (long + (radius/111259.54243971))) AND
		(ST_Y("MY_TABLE".coordinates::geometry) BETWEEN (lat - (radius/111259.54243971)) AND (lat + (radius/111259.54243971)))
	WHEN shape = 'circle' THEN
		ST_Distance("MY_TABLE".coordinates, CONCAT('POINT(',long::varchar(255),' ',lat::varchar(255),')')::geography) < radius
	END
ORDER BY
	"MY_TABLE".coordinates <-> CONCAT('POINT(',long::varchar(255),' ',lat::varchar(255),')')::geography,
	"MY_TABLE".rating ASC;
END ; $$ LANGUAGE 'plpgsql';


-- Call to the function
SELECT * FROM get_spots(100.12, 100.3, 5000000.0, 'circle');

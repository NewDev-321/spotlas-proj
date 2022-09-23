
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

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pgAdmin1234"
	dbname   = "dev"
)

type MY_TABLE struct {
	id          string
	name        string
	website     sql.NullString
	coordinates string
	description sql.NullString
	rating      float64
}

// function to retrieve spots
func retrieveSpots(lat float64, long float64, radius float64, shape string) []MY_TABLE {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	// query execution
	queryStr := `SELECT * FROM "MY_TABLE"
			WHERE CASE
				WHEN $4 = 'square' THEN
					(ST_X("MY_TABLE".coordinates::geometry) BETWEEN ($1::float8 - ($3::float8 / 111259.54243971)) AND ($1::float8 + ($3::float8 / 111259.54243971))) AND
					(ST_Y("MY_TABLE".coordinates::geometry) BETWEEN ($2::float8 - ($3::float8 / 111259.54243971)) AND ($2::float8 + ($3::float8 / 111259.54243971)))
				WHEN $4 = 'circle' THEN
					ST_Distance("MY_TABLE".coordinates, CONCAT('POINT(',$1::varchar(255),' ',$2::varchar(255),')')::geography) > $3::float8
				END
				ORDER BY
					"MY_TABLE".coordinates <-> CONCAT('POINT(',$1::varchar(255),' ',$2::varchar(255),')')::geography,
					"MY_TABLE".rating ASC`

	rowsRs, err := db.Query(queryStr, long, lat, radius, shape)
	if err != nil {
		panic(err)
		return []MY_TABLE{}
	}
	defer rowsRs.Close()
	// creates placeholder of the table
	tbl := make([]MY_TABLE, 0)

	// we loop through the values of rows
	for rowsRs.Next() {
		row := MY_TABLE{}
		err := rowsRs.Scan(&row.id, &row.name, &row.website, &row.coordinates, &row.description, &row.rating)
		if err != nil {
			panic(err)
			return []MY_TABLE{}
		}
		tbl = append(tbl, row)
	}
	return tbl
}

func main() {

	// call to the function
	result := retrieveSpots(300.21, 200.05, 200100.12, "circle")

	// loop and display the result in the browser
	for _, row := range result {
		fmt.Println(row.id, row.name)
	}

}

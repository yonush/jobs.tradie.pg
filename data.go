package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gookit/validate"
	_ "github.com/jackc/pgx/v5/stdlib" //use pgx in database/sql mode
)

type Jobs struct {
	Jobid     int    `json:"jobid"  validate:"required|max:2017"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Name      struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Address string   `json:"address"`
	Phone   string   `json:"phone"`
	Email   string   `json:"email"`
	Notes   []string `json:"notes,omitempty"`
}

// import the JSON data into a collection
func loadFromJson(db *sql.DB, filename string) error {
	// Create table as required, along with attribute constraints
	sql := `CREATE TABLE IF NOT EXISTS jobs
	(
		jobid INT PRIMARY KEY NOT NULL,
		status VARCHAR (20) NOT NULL,
		datum timestamp,		
		firstname VARCHAR (20) NOT NULL,
		lastname VARCHAR (20) NOT NULL,
		address VARCHAR (50) NOT NULL,
		phone VARCHAR (12) NOT NULL,
		email VARCHAR (50) NOT NULL,
		notes text
	);`

	//create the jobs table seeing as it possible does not exist
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("jobs table has been successfully created.")

	// open the json file for importing in PG database
	data, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer data.Close()

	// prepare the SQL for multiple inserts
	stmt, err := db.Prepare("INSERT INTO jobs VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)")
	if err != nil {
		log.Fatal(err)
	}

	/* read the first [ character otherwise you will get the following error
	 * json: cannot unmarshal array into Go value of type
	 */
	var msg Jobs
	var notes string
	dec := json.NewDecoder(bufio.NewReader(data))

	dec.Token()
	for dec.More() {
		msg = Jobs{} //make sure its empty
		err := dec.Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}

		// valdate the data from the json file
		v := validate.Struct(msg)
		if v.Validate() {
			notes = ""
			notes = strings.Join(msg.Notes, "#") //convert notes into a string
			//s := strings.Split(myString, "?")[0] string
			_, err = stmt.Exec(msg.Jobid, msg.Status, msg.Timestamp,
				msg.Name.First, msg.Name.Last,
				msg.Address, msg.Phone, msg.Email, notes)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// create temp file to notify data imported
	//can use database directly but this is an example
	// https://golangbyexample.com/touch-file-golang/
	file, err := os.Create("./imported.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return err
}

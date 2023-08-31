package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// User and Segment struct to convert from or to JSON
type User struct {
	Id        int    `json:"user_id"`
	User_name string `json:"user_name"`
	Slug      []Segment
}

type Segment struct {
	Slug string `json:"slug"`
}

var DB *sql.DB

func ConnectDB() {
	var err error
	host := loadEnvVar("DBhost")
	port := loadEnvVar("DBport")
	user := loadEnvVar("DBuser")
	password := loadEnvVar("DBpassword")
	dbname := loadEnvVar("DBname")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	DB, err = sql.Open("postgres", psqlconn)
	checkErr(err)

	err = DB.Ping()
	checkErr(err)
}

// Get a user by id
func GetUser(id int) string {
	// Select a user and join the segment table to get all the segments a user is in
	sql := "SELECT su.user_id, u.user_name, s.slug FROM users u " +
		"LEFT JOIN segment_user su ON su.user_id = u.id " +
		"LEFT JOIN segment s ON su.segment_id = s.id " +
		"WHERE u.id = $1 AND (su.ttl > now() OR su.ttl IS NULL) "
	rows, err := DB.Query(sql, id)
	defer rows.Close()

	checkErr(err)
	var user User
	// read rows from the result of the query
	for rows.Next() {
		var seg Segment
		// scan rows and put each field into its respective field of the struct
		rows.Scan(&user.Id, &user.User_name, &seg.Slug)
		// add the segments to the Slug field
		user.Slug = append(user.Slug, seg)
	}
	// convert user struct to JSON
	responseData, err := json.Marshal(user)
	// return the JSON string
	return string(responseData)
}

// Create a new segment
func CreateSegment(slug string) int {
	// insert the segment to the database table
	sql := "INSERT INTO segment (slug) VALUES ($1)"
	_, err := DB.Exec(sql, slug)

	//return 0 on error 1 on success
	if err != nil {
		return 0
	} else {
		return 1
	}
}

// Delete a segment from the database table
func DeleteSegment(slug string) int {
	sql := "DELETE FROM segment WHERE slug = $1"
	_, err := DB.Exec(sql, slug)

	//return 0 on error 1 on success
	if err != nil {
		return 0
	} else {
		return 1
	}
}

// function to write a entry to the history log
func writeLog(user int, segment int, action string) {
	sql := "INSERT INTO eventlog (action,created_at,user_id,segment_id) VALUES ($1,now(),$2,$3)"
	DB.Exec(sql, action, user, segment)
}

// Add a user to one segment
func AddUserToSegment(user int, segment int) int {
	sql := "INSERT INTO segment_user (segment_id,user_id,created_at) VALUES ($1,$2,now())"
	_, err := DB.Exec(sql, segment, user)

	if err != nil {
		return 0
	} else {
		writeLog(user, segment, "add")
		return 1
	}
}

// ttl variant to add a user, ttl is days in the future
func AddUserToSegmentTTL(user int, segment int, ttl string) int {
	sql := "INSERT INTO segment_user (segment_id,user_id,created_at,ttl) VALUES ($1,$2,now(),now()+ interval '" + ttl + " days')"
	_, err := DB.Exec(sql, segment, user)

	if err != nil {
		log.Println(err)
		return 0
	} else {
		writeLog(user, segment, "add")
		return 1
	}
}

// Delete a user from a segment
func RemoveUserFromSegment(user int, segment int) int {
	sql := "DELETE FROM segment_user WHERE segment_id = $1 AND user_id = $2"
	_, err := DB.Exec(sql, segment, user)

	if err != nil {
		return 0
	} else {
		writeLog(user, segment, "remove")
		return 1
	}
}

// helper function to get the id of a segment by the slug
func GetSegmentBySlug(slug string) int {
	var id int

	sql := "SELECT id FROM segment WHERE slug = $1"
	rows, err := DB.Query(sql, slug)
	for rows.Next() {
		rows.Scan(&id)
	}
	defer rows.Close()
	checkErr(err)
	return id
}

// Add random users to a segment based on a percentage
func AddRandomUserPercent(percent int, segment_id int) {
	var amount int
	var user_id int

	// Select the amount of users in the table
	sql := "SELECT count(id) FROM users"
	rows, _ := DB.Query(sql)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&amount)
	}

	// calculate the user limit and select random users
	userLimit := int(amount * percent / 100)
	sql = "SELECT id FROM users ORDER BY random() LIMIT $1"

	rows2, _ := DB.Query(sql, userLimit)
	defer rows2.Close()
	for rows2.Next() {
		rows2.Scan(&user_id)
		// add those random users to the segment
		AddUserToSegment(user_id, segment_id)
		writeLog(user_id, segment_id, "add")
	}

}

func loadEnvVar(envVar string) string {
	err := godotenv.Load(".env")
	checkErr(err)
	return os.Getenv(envVar)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

package SqlHelper 

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"time"
	_ "embed"
)


var db *sql.DB
var DatabasePath = "database.db"

//go:embed create_tables.sql
var createTablesQuery string

/* Message Struct - like what we will have in the db */
type Message struct {
	Date int64
	Device string
	Checksum string
	Mode int
	Message string
}

/* configuration key struct */
type Config struct {
	Key string
	Value string
}

func init() {
	InitializeDB()
}

/**
 * Panic if the checked error is not 'nil'
 */
func check (err error) {
	if (err != nil) {
		panic(err)
	}
}

/**
 * Create the database tables if not exist
 */ 
func createTables() {
	// create tables if not exist
	_, err := db.Exec(createTablesQuery)
	check(err)
}

/**
 * initialize the database	 
 */
func InitializeDB() {
	var err error = nil

	// create the database (or connect if already exist)
	db, err = sql.Open("sqlite", DatabasePath)
	check(err)

	// create database tables 
	createTables()	
}

/* Push a new message inside the database */
func InsertNewMessage(message Message) {
	stmt, err := db.Prepare("INSERT INTO messages (date, device, checksum, mode, message) VALUES (?, ?, ?, ?, ?) ")
	check(err)

	_, err = stmt.Exec(message.Date, message.Device, message.Checksum, message.Mode, message.Message)
	check(err)
}

/* Get the last message from the database (by timestamp) */
func QueryLastMessage() Message {
	var message Message

	stmt, err := db.Prepare("SELECT * FROM messages ORDER BY date DESC LIMIT 1")
	check(err)

	rows, err := stmt.Query()
	check(err)
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&message.Date, &message.Device, &message.Checksum, &message.Mode, &message.Message)
		check(err)

	}

	return message
}

/* delete a message from the db based on timestamp */
func DeleteMessageByTimestamp(timestamp int64) {
	stmt, err := db.Prepare("DELETE FROM messages WHERE date = ?")
	check(err)

	_, err = stmt.Exec(timestamp)
	check(err)
}

/* Push a new configuration key to the database */
func InsertConfig(config Config) {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO config (key, value) VALUES (?, ?) ")
	check(err)

	_, err = stmt.Exec(config.Key, config.Value)
	check(err)
}

/* Get configuration by key */
func QueryConfigByKey(key string) Config {
	var config Config

	stmt, err := db.Prepare("SELECT * FROM config WHERE key = ?")
	check(err)

	rows, err := stmt.Query(key)
	defer rows.Close()
	check(err)
	

	//Iterate through result set
	for rows.Next() {
		err := rows.Scan(&config.Key, &config.Value)
		check(err)
	}

	return config
}

/* update some configuration key */
func UpdateConfig(config Config) {
	stmt, err := db.Prepare("UPDATE config SET value = ? where key = ?")
	check(err)

	_, err = stmt.Exec(config.Key, config.Value)
	check(err)
}

/* delete config value based on key */
func DeleteConfigByKey(key string) {
	stmt, err := db.Prepare("DELETE FROM config WHERE key = ?")
	check(err)

	_, err = stmt.Exec(key)
	check(err)
}

func GetTimestamp() int64 {
	now := time.Now()      // current local time
	sec := now.Unix()      // number of seconds since January 1, 1970 UTC

	return sec 	
}

func main() {
	InitializeDB()
}
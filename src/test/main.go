package main

func main() {
	// connect to Database
	db := ConnectToDatabase()

	// create tables
	CreateTables(db)
}

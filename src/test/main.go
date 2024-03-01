package main

func main() {
	// connect to Database
	db := ConnectToDatabase()

	// create the user_basic table
	CreateUserTable(db)

	// test User table
}

package main

func main() {
	// connect to Database
	db := ConnectToDatabase()

	// create the user_basic table
	// CreateUserTable(db)

	// create the Relation table
	CreateRelationTable(db)
}

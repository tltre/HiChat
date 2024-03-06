package main

func main() {
	// 1. MySQL Initial
	/*
		// connect to Database
		db := ConnectToDatabase()

		// create tables
		CreateTables(db)
	*/
	// 2. Redis Initial
	redisClient := ConnectToRedis()
	CheckIfConnectRedis(redisClient)
}

package main

func main() {
	db := connectDB()
	defer db.Close()
}

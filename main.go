package main

//	@title			Swagger Petstore
//	@version		1.0.7
//	@description	This is a sample server Petstore server.

//	@host		localhost:8080
//	@BasePath	/v2

//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization

//	@tag.name			pet
//	@tag.description	Everything about your Pets
//	@tag.name			store	
//	@tag.description	Access to Petstore orders
//	@tag.name			user	
//	@tag.description	Operations about users

// main runs the server on the given address.
func main() {
	NewServer(":8080").Serve()
}

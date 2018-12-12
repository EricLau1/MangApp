package main 

import(
	"net/http"
	"fmt"
	"./routes"
	"./utils"
)

func main() {
	
	fmt.Println("Listening port 8080")

	utils.LoadTemplates("templates/*.html")

	router := routes.NewRouter()

	http.Handle("/", router)

	http.ListenAndServe(":8080", nil)

}
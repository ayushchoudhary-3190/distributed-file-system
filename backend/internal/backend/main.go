package backend

import (

	"github.com/gofiber/fiber/v2"
)

func main(){

	app := fiber.New(fiber.Config{
		AppName:"distributed file system",
		Views: engine,
	})
	
	api:=app.Group("/api/metaservice")
	api.Post("/create",func(c *fiber.Ctx)error{
		
	})

}
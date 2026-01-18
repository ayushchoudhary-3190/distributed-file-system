package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type File_info struct{
	FileName   string
	OwnerID    string
	ChunkCount int64
	FileSize   int64
}

func main(){

	app := fiber.New(fiber.Config{
		AppName:"distributed file system",
	})
	
	// Metaservice endpoints
	api:=app.Group("/api")
	api.Post("/createfile",func(c *fiber.Ctx)error{
		fileName:= c.Get("fileName")
		ownerID:=c.Get("userID")
		chunkCountstr:=c.Get("chunkCount")
		fileSizestr:= c.Get("fileSize")

		chunkCount,err:= StringToInt(chunkCountstr)
		if err!=nil{
			return err
		}
		fileSize,err:= StringToInt(fileSizestr)
		if err!=nil{
			return err
		}

		NewFile:= &File_info{
			FileName: fileName,
			OwnerID: ownerID,
			ChunkCount: chunkCount,
			FileSize: fileSize,
		}

		fileID,err:= client.AddFile(NewFile)

		return c.JSON(fiber.Map{
			"fileID":fileID,
		})
	})





	api.Post("/deletefile",func(c *fiber.Ctx)error{ 
			fileID := c.Get("fileid")
			if fileID == "" {
				return fiber.NewError(
					fiber.StatusBadRequest,
					"missing fileid header",
				)
			}

			if err := client.DeleteFile(fileID); err != nil {
				return fiber.NewError(
					fiber.StatusInternalServerError,
					err.Error(),
				)
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "file successfully deleted",
				"fileId":  fileID,
			})
	})



	api.Get("/getuserfiles",func(c *fiber.Ctx)error{
		ownerID:= c.Get("ownerid")
		files, err := client.GetUserFiles(ownerID)
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				err.Error(),
			)
		}

		return c.Status(fiber.StatusOK).JSON(files)
	})




	api.Get("/getfile",func(c *fiber.Ctx)error{
		ownerID:= c.Get("ownerid")
		fileID:= c.Get("fileid")
		file, err:= client.GetFile(ownerID,fileID)
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				err.Error(),
			)
		}

		return c.Status(fiber.StatusOK).JSON(file)
	})
}
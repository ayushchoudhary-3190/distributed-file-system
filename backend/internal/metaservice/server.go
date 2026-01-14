package main

import (
	"context"
	"log"

	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/metaservice"
	"github.com/ayushchoudhary-3190/Distributed_file_system/pb/github.com/ayushchoudhary-3190/grpc_project/backend/pb"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main(){

	dsn := ""
	var err error
	DB,err = gorm.Open(sqlite.Open(dsn),gorm.config{})
	if err!=nil{
		log.Fatal("failed to connect database",err)
	}

	err:=DB.AutoMigrate(&metaservice.File_table{},&metaservice.Chunk_table{},&metaservice.Node_table{})
	if err!=nil{
		log.Fatal("failed to create tables")
	}

	server:=grpc.NewServer()
	metaserver:= &metaservice.Metaserver{
		DB :DB,
	}

	pb.RegisterMetaServiceServer(server,metaserver)
}
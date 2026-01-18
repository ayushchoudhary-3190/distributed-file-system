package metaservice

import (
	"context"
	"log"

	"github.com/ayushchoudhary-3190/Distributed_file_system/pb/github.com/ayushchoudhary-3190/grpc_project/backend/pb"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Metaserver struct {
	DB *gorm.DB
}

type Metaservice struct{
	pb.UnimplementedMetaServiceServer
}

func main(){

	dsn := ""
	var err error
	DB,err = gorm.Open(sqlite.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("failed to connect database",err)
	}

	err = DB.AutoMigrate(&File_table{},&Chunk_table{},&Node_table{})
	if err!=nil{
		log.Fatal("failed to create tables")
	}

	server:=grpc.NewServer()
	metaserver:= &Metaserver{
		DB :DB,
	}

	pb.RegisterMetaServiceServer(server,&Metaservice{})
}
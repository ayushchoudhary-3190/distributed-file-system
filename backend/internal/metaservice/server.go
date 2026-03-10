package metaservice

import (
	"log"

	"github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	datanodeservice "github.com/ayushchoudhary-3190/Distributed_file_system/internal/datanodes"
)

var DB *gorm.DB

type Metaserver struct {
	DB *gorm.DB
}

type Metaservice struct {
	pb.UnimplementedMetaServiceServer
}

func main() {

	dsn := ""
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	err = DB.AutoMigrate(&File_table{}, &datanodeservice.Chunk_table{}, &datanodeservice.Node_table{})
	if err != nil {
		log.Fatal("failed to create tables")
	}

	server := grpc.NewServer()
	_ = Metaserver{
		DB: DB,
	}

	pb.RegisterMetaServiceServer(server, &Metaservice{})
}

package client

import (
	"log"

	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/metaservice"
	"github.com/google/uuid"
	"gorm.io/gorm"
)



func AddFile(newFile *metaservice.NewFile_Params) (uuid.UUID,string) {
	fileID,err:= metaservice.CreateFile(newFile)
	if err!=nil{
		log.Fatal("failed to create metaservice entry")
	}

}

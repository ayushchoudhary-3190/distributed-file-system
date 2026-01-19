package metaservice

import (
	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/client"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type NewFile_Params struct{
	newFile *File_info
	db 		*gorm.DB
}

type File_info struct{
	FileName   string
	OwnerID    string
	ChunkCount int64
	ChunkArray []string
	FileSize   int64
}


//  function to add a new file to the metaservice table
func  CreateFile(newFileInfo *NewFile_Params) (uuid.UUID, error) {
	fileID:= uuid.New()
	err := newFileInfo.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			INSERT INTO file_tables(fileid,file_name,ownerid,chunkcount,chunkarray,filesize) VALUES(?,?,?,?,?,?)
		`,fileID, newFileInfo.newFile.FileName, newFileInfo.newFile.OwnerID, newFileInfo.newFile.ChunkCount, newFileInfo.newFile.ChunkArray, newFileInfo.newFile.FileSize)
		if result.Error != nil {
			return result.Error
		}

		return result.Error
	})
	if err != nil {
		return  uuid.Nil,err
	}
	return fileID, nil
}





package metaservice

import (
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

type NewFileEntry struct {
	FileName   string
	OwnerID    string
	ChunkCount int64
	ChunkArray []string
	FileSize   int64
}
func NewFileEntity(fileName string, ownerID string, chunkCount int64, chunkArray []string, fileSize int64) *NewFileEntry {
	return &NewFileEntry{
		FileName:   fileName,
		OwnerID:    ownerID,
		ChunkCount: chunkCount,
		ChunkArray: chunkArray,
		FileSize:   fileSize,
	}
}
//  function to add a new file to the metaservice table
func (DB *Database) CreateFile(newFile *NewFileEntry) (int64, error) {
	var insertedID int64
	err := DB.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			INSERT INTO file_tables(file_name,ownerid,chunkcount,chunkarray,filesize) VALUES(?,?,?,?,?)
		`, newFile.FileName, newFile.OwnerID, newFile.ChunkCount, newFile.ChunkArray, newFile.FileSize)
		if result.Error != nil {
			return result.Error
		}

		insertedID = result.RowsAffected
		return nil
	})
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}





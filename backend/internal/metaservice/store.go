package metaservice

type FileStore interface{
	CreateFile(fileName string,ownerName string) (int64, error)
	DeleteFile(fileID int64 ) (string,error)
	GetAllFiles(ownerID string) ([]string,error)
	GetFile(fileID int64) ([]string,error)
}

type ChunkStore interface{
	Allocator(fileID int64, fileSize int64 ) ([]string,error)
	GetChunks(fileID int64) ([]string,error)
}

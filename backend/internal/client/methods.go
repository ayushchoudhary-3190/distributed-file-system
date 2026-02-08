package client

import (
	"log"

	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/metaservice"
	"github.com/ayushchoudhary-3190/Distributed_file_system/pb/github.com/ayushchoudhary-3190/grpc_project/backend/pb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


// Client function to add a new file 
func AddFile(newFile *pb.UploadFileRequest) (string,string) {
	fileID:= uuid.New()

	// Allocate metadata to chunks
	chunkMetadata:= metaservice.AllocateChunks(fileID,newFile.fileSize,count)
	// write chunks to disk and datanode table
	err:= datanodeservice.WriteChunks(chunkMetadata,file,offset,size)
	if err!=nil{
		log.Fatal("failed to write chunk to disk")
		return "Server Error :" , "Failure in writing chunks on the disk"
	}

	// Write metadata to metaservice
	path,res:= metaservice.UploadRequest(newFile.Filename,newFile.size,newFile.chunks)
	return path,res
}

// Client function to delete a file
func DeleteFile(filePath string) string{

	// Call Datanodeservice delete function 
	res:= datanodeservice.DeleteChunk(fileID)

	// Call metaservice delete function
	res,err:= metaservice.DeleteRequest(filePath)
	if err!=nil{
		return "failed to delete file"
	}

	return "file Deleted successfully"
}

func ListFiles(owner string ,ownerID string) *ListFilesResponse{
	res := metaservice.ListFiles(owner,ownerID)
	if res.err != nil{
		log.Fatal("failed to retrieve files")
		return &ListFilesResponse{}
	}

	return res
}

func GetFile(ownerID string, path string) []byte{
	res:= metaservice.GetFile(ownerID,path)
	if res.err != nil{
		log.Fatal("failed to retrieve file")
		return &pb.GetFileResponse{}
	}
	return res;
}

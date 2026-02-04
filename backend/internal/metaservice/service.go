package metaservice

import (
	"context"

	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/client"
	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/metaservice"
	pb "github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetaServer struct {
	pb.UnimplementedMetaServiceServer
	DB *gorm.DB
}

// function to add a new file to the metaservice table
func (s *MetaServer) UploadRequest(ctx *context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, string) {
	//insert file metadata in metadata table
	tx := s.DB.Begin()

	// Extract chunk IDs from request
	chunkArray := make([]string, len(req.Chunks))
	for i, chunk := range req.Chunks {
		chunkArray[i] = chunk.Chunkid
	}

	file := metaservice.File_table{
		FileID:     req.Fileid,
		FileName:   req.Filename,
		OwnerID:    req.Ownerid,
		ChunkCount: req.Chunkcount,
		ChunkArray: chunkArray,
		FileSize:   req.Filesize,
	}

	// Create file record in database
	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		response := &pb.UploadFileResponse{
			Path:     req.Filename,
			Response: "Failed to upload file metadata",
		}
		return response, err.Error()
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.UploadFileResponse{
			Path:     req.Filename,
			Response: "Failed to commit transaction",
		}
		return response, err.Error()
	}

	// Return success response
	response := &pb.UploadFileResponse{
		Path:     req.Filename,
		Response: "File uploaded successfully",
	}
	return response, " "
}

// function to delete a file from the metaservice table
func (s *MetaServer) DeleteRequest(ctx *context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, string) {
	// Start transaction
	tx := s.DB.Begin()

	// Find and delete the file with the given path
	result := tx.Where("file_name = ?", req.Path).Delete(&metaservice.File_table{})

	if result.Error != nil {
		tx.Rollback()
		response := &pb.DeleteFileResponse{
			Message: "Failed to delete file",
			Error:   result.Error.Error(),
		}
		return response, result.Error.Error()
	}

	// Check if any record was actually deleted
	if result.RowsAffected == 0 {
		tx.Rollback()
		response := &pb.DeleteFileResponse{
			Message: "File not found",
			Error:   "No file exists with the given path",
		}
		return response, "No file exists with the given path"
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.DeleteFileResponse{
			Message: "Failed to commit transaction",
			Error:   err.Error(),
		}
		return response, err.Error()
	}

	// Return success response
	response := &pb.DeleteFileResponse{
		Message: "File deleted successfully",
		Error:   "",
	}
	return response, " "
}

// function to list files belonging to a specific owner
func (s *MetaServer) ListFiles(ctx *context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, string) {
	var files []metaservice.File_table

	// Query files by owner_id
	result := s.DB.Where("owner_id = ?", req.OwnerId).Find(&files)

	if result.Error != nil {
		response := &pb.ListFilesResponse{
			Owner:   req.Owner,
			OwnerId: req.OwnerId,
			Count:   0,
			Err:     result.Error.Error(),
		}
		return response, result.Error.Error()
	}

	// Create FileListItem array with only file names
	fileList := make([]*pb.FileListItem, len(files))
	for i, file := range files {
		fileList[i] = &pb.FileListItem{
			Filename: file.FileName,
		}
	}

	// Return success response
	response := &pb.ListFilesResponse{
		Filename: fileList,
		Owner:    req.Owner,
		OwnerId:  req.OwnerId, // keeping as 0 since owner_id is string in proto but int64 in response
		Count:    int64(len(fileList)),
		Err:      "",
	}
	return response, " "
}

package metaservice

import (
	"context"
	"sync"

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

// ChunkDataWithIndex holds chunk data along with its index for ordered reconstruction
type ChunkDataWithIndex struct {
	Index int32
	Data  []byte
}

// function to add a new file to the metaservice table
func (s *MetaServer) UploadRequest(ctx *context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
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
		return response, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.UploadFileResponse{
			Path:     req.Filename,
			Response: "Failed to commit transaction",
		}
		return response, err
	}

	// Return success response
	response := &pb.UploadFileResponse{
		Path:     req.Filename,
		Response: "File uploaded successfully",
	}
	return response, nil
}

// function to delete a file from the metaservice table
func (s *MetaServer) DeleteRequest(ctx *context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
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
		return response, result.Error
	}

	// Check if any record was actually deleted
	if result.RowsAffected == 0 {
		tx.Rollback()
		response := &pb.DeleteFileResponse{
			Message: "File not found",
			Error:   "file not deleted",
		}
		return response, nil
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response := &pb.DeleteFileResponse{
			Message: "Failed to commit transaction",
			Error:   err.Error(),
		}
		return response, err
	}

	// Return success response
	response := &pb.DeleteFileResponse{
		Message: "File deleted successfully",
		Error:   "",
	}
	return response, nil
}

// function to list files belonging to a specific owner
func (s *MetaServer) ListFiles(ctx *context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
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
		return response, result.Error
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
		OwnerId:  req.OwnerId,
		Count:    int64(len(fileList)),
		Err:      "",
	}
	return response, nil
}

// function to get file by owner_id and path and reconstruct from chunks using new workflow
func (s *MetaServer) GetFile(ctx *context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	var file metaservice.File_table

	// Query file by owner_id and file_name (path)
	result := s.DB.Where("owner_id = ? AND file_name = ?", req.OwnerId, req.Path).First(&file)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			response := &pb.GetFileResponse{
				OwnerId:  req.OwnerId,
				Path:     req.Path,
				Size:     0,
				Err:      "File not found",
				FileData: nil,
			}
			return response, result.Error
		}
		response := &pb.GetFileResponse{
			OwnerId:  req.OwnerId,
			Path:     req.Path,
			Size:     0,
			Err:      result.Error.Error(),
			FileData: nil,
		}
		return response, result.Error
	}

	// Get chunkIDs array from table using ownerid and path from parameter
	chunkIDs := file.ChunkArray

	// Use that chunkIDs array inside reconstructFileFromId function
	fileData := s.reconstructFileFromId(chunkIDs)

	// Return success response with reconstructed file data
	response := &pb.GetFileResponse{
		OwnerId:  file.OwnerID,
		Path:     req.Path,
		Size:     file.FileSize,
		Err:      "",
		FileData: fileData,
	}
	return response, nil
}

// reconstructFileFromId reconstructs file from chunk IDs using location-based approach
func (s *MetaServer) reconstructFileFromId(chunkIDs []string) []byte {
	// Call getChunksLocation function to get locations for all chunk IDs
	locations := s.getChunksLocation(chunkIDs)

	// Use these locations inside readChunks function for each address
	return s.readChunks(chunkIDs, locations)
}

// getChunksLocation function to get locations for chunk IDs (placeholder implementation)
func (s *MetaServer) getChunksLocation(chunkIDs []string) map[string]string {
	// TODO: Implement actual location retrieval
	// This should return a map of chunkID -> location/address
	// For now, return empty map
	locations := make(map[string]string)
	for _, chunkID := range chunkIDs {
		locations[chunkID] = "" // placeholder
	}
	return locations
}

// readChunks function to read chunks from their respective addresses and append them (placeholder implementation)
func (s *MetaServer) readChunks(chunkIDs []string, locations map[string]string) []byte {
	// TODO: Implement actual chunk reading from different addresses
	// This should read chunks from their locations and append them in order
	// For now, return empty bytes
	fileData := []byte{}
	for _, chunkID := range chunkIDs {
		// Read chunk from its location
		chunkData := s.readChunkFromLocation(chunkID, locations[chunkID])
		// Append chunk to file data
		fileData = append(fileData, chunkData...)
	}
	return fileData
}

// readChunkFromLocation function to read a chunk from a specific location (placeholder implementation)
func (s *MetaServer) readChunkFromLocation(chunkID string, location string) []byte {
	// TODO: Implement actual chunk reading from specific location
	// For now, return empty bytes
	return []byte{}
}

package metaservice

import (
	"context"
	"log"
	"sync"

	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/client"
	datanodeservice "github.com/ayushchoudhary-3190/Distributed_file_system/internal/datanodes"
	"github.com/ayushchoudhary-3190/Distributed_file_system/internal/metaservice"
	pb "github.com/ayushchoudhary-3190/Distributed_file_system/pb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetaServer struct {
	pb.UnimplementedMetaServiceServer
	DB             *gorm.DB
	dataNodeClient pb.DataNodeServiceClient
}

// ChunkDataWithIndex holds chunk data along with its index for ordered reconstruction
type ChunkDataWithIndex struct {
	Index int32
	Data  []byte
}

// function to add a new file to the metaservice table
func (s *MetaServer) UploadRequest(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
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
func (s *MetaServer) DeleteRequest(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
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
func (s *MetaServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
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
func (s *MetaServer) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	var file metaservice.File_table

	// Query file by owner_id and file_name (path)
	result := s.DB.Where("owner_id = ? AND file_name = ?", req.OwnerId, req.Path).First(&file)

	if result.Error != nil {
		response := &pb.GetFileResponse{
			OwnerId: req.OwnerId,
			Path:    req.Path,
			Size:    0,
			File:    nil,
		}
		return response, result.Error
	}

	// Get chunkIDs array from table using ownerid and path from parameter
	chunkIDs := file.ChunkArray

	// Use that chunkIDs array inside reconstructFileFromId function
	fileData ,err:= s.reconstructFileFromId(chunkIDs)

	if err!=nil{
		log.Fatal("failed to reconstruct file from chunk ids")
		response:= &pb.GetFileResponse{
			OwnerId: req.OwnerId,
			Path: req.Path,
			Size: 0,
			File: nil,
		}
		return response , err
	}

	// Return success response with reconstructed file data
	response := &pb.GetFileResponse{
		OwnerId: file.OwnerID,
		Path:    req.Path,
		Size:    file.FileSize,
		File:    fileData,
	}
	return response, nil
}

// reconstructFileFromId reconstructs file from chunk IDs using location-based approach
func (s *MetaServer) reconstructFileFromId(chunkIDs []string) ([]byte, error) { //// metaservice function
	// Get locations for chunk IDs
	locations := s.getChunksLocation(chunkIDs)

	var fileData []byte

	// For each chunkID, call ReadChunks via DataNode client
	for _, chunkID := range chunkIDs {
		// Prepare gRPC request
		req := &pb.ChunkReadRequest{
			ChunkId: chunkID,
			Offset:  0,
			Length:  65536, // appropriate chunk size
		}

		// Call DataNode via client - this is the correct way
		stream, err := s.dataNodeClient.ReadChunks(context.Background(), req)
		if err != nil {
			return fileData, err // Return error from response
		}

		// Receive ALL responses from stream in a loop
		for {
			resp, err := stream.Recv()
			if err != nil {
				break // Stream ended
			}

			// Append only the DATA part from response
			fileData = append(fileData, resp.Data...)

			// Check for EOF - if true, no more data for this chunk
			if resp.Eof {
				break
			}
		}
	}

	return fileData, nil // Return data bytes and nil error on success
}

// getChunksLocation function to get locations for chunk IDs (placeholder implementation)
func (s *MetaServer) getChunksLocation(chunkIDs []string) map[string]string { ////metaservice fucntion
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
//func (s *MetaServer) readChunks(chunkIDs []string, locations map[string]string) []byte {  /////datanodeservice function
// TODO: Implement actual chunk reading from different addresses
// This should read chunks from their locations and append them in order
// For now, return empty bytes
//	fileData := []byte{}
//	for _, chunkID := range chunkIDs {
//		// Read chunk from its location
//		chunkData := s.readChunkFromLocation(chunkID, locations[chunkID])
//		// Append chunk to file data
//		fileData = append(fileData, chunkData...)
//	}
//	return fileData
//}

// readChunkFromLocation function to read a chunk from a specific location (placeholder implementation)
func (s *MetaServer) readChunkFromLocation(chunkID string, location string) []byte { ////datanodeservice function
	// TODO: Implement actual chunk reading from specific location
	// For now, return empty bytes
	return []byte{}
}

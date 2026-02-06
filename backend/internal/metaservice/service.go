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
		OwnerId:  req.OwnerId,
		Count:    int64(len(fileList)),
		Err:      "",
	}
	return response, " "
}

// function to get file by owner_id and path and reconstruct from chunks using parallel processing
func (s *MetaServer) GetFile(ctx *context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, string) {
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
			return response, "File not found"
		}
		response := &pb.GetFileResponse{
			OwnerId:  req.OwnerId,
			Path:     req.Path,
			Size:     0,
			Err:      result.Error.Error(),
			FileData: nil,
		}
		return response, result.Error.Error()
	}

	// Get chunk array from table
	chunkIDs := file.ChunkArray

	// Reconstruct file from chunks using parallel processing
	fileData := s.reconstructFileFromChunks(chunkIDs)

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

// reconstructFileFromChunks reconstructs file from chunks using parallel processing
func (s *MetaServer) reconstructFileFromChunks(chunkIDs []string) []byte {
	// Create channel with 128MB buffer (128 * 1024 * 1024 bytes)
	chunkChannel := make(chan *ChunkDataWithIndex, 128*1024*1024)
	var wg sync.WaitGroup

	// Start workers for each chunk
	for index, chunkID := range chunkIDs {
		wg.Add(1)
		go func(idx int, cid string) {
			defer wg.Done()
			// Read chunk from disk
			chunkData := s.readChunk(cid)
			// Add chunk to channel with index
			chunkChannel <- &ChunkDataWithIndex{
				Index: int32(idx),
				Data:  chunkData,
			}
		}(index, chunkID)
	}

	// Start a goroutine to close channel when all workers are done
	go func() {
		wg.Wait()
		close(chunkChannel)
	}()

	// Collect chunks from channel and reconstruct file in order
	chunkMap := make(map[int32][]byte)
	for chunkData := range chunkChannel {
		chunkMap[chunkData.Index] = chunkData.Data
	}

	// Reconstruct file data in correct order
	fileData := []byte{}
	for i := 0; i < len(chunkIDs); i++ {
		if chunkData, exists := chunkMap[int32(i)]; exists {
			fileData = append(fileData, chunkData...)
		}
	}

	return fileData
}

// readChunk function to read chunk from disk (placeholder implementation)
func (s *MetaServer) readChunk(chunkID string) []byte {
	// TODO: Implement actual chunk reading from disk
	// For now, return empty bytes
	// This will be implemented later
	return []byte{}
}

// appendChunk function to append chunk data to file data (placeholder implementation)
func (s *MetaServer) appendChunk(existingData []byte, newChunkData []byte) []byte {
	// TODO: Implement proper chunk appending logic
	// For now, simply append the new chunk data
	// This will be implemented later
	return append(existingData, newChunkData...)
}

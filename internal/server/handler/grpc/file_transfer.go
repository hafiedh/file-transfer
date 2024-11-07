package grpc

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "hafiedh.com/downloader/gen/go/file_transfer/v2"
	"hafiedh.com/downloader/internal/usecase/downloader"
)

func NewFileTransferHandler(fileTransferService downloader.FileTransfer) pb.FileTransferServiceServer {
	return &fileTransferHandler{
		fileTransferService: fileTransferService,
	}
}

func (h *fileTransferHandler) UploadFile(stream pb.FileTransferService_UploadFileServer) error {
	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed to receive file info: %v", err)
	}

	fileInfo := req.GetFileInfo()
	if fileInfo == nil {
		return status.Error(codes.InvalidArgument, "first message must contain file info")
	}

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive chunk: %v", err)
		}

		chunk := req.GetChunk()
		if chunk == nil {
			continue
		}

		if _, err := writer.Write(chunk); err != nil {
			return status.Errorf(codes.Internal, "failed to write chunk: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return status.Errorf(codes.Internal, "failed to flush buffer: %v", err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: fileInfo.Filename,
		Size:     fileInfo.Size,
	}

	files := map[string][]*multipart.FileHeader{
		fileInfo.Key: {fileHeader},
	}

	resp, err := h.fileTransferService.UploadFile(ctx, files, nil)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upload file: %v", err)
	}

	grpcResp := &pb.UploadFileResponse{
		Message: resp.Message,
		Status:  int32(resp.Status),
		Errors:  resp.Errors,
	}

	if resp.Status != http.StatusOK {
		return stream.SendAndClose(grpcResp)
	}

	return stream.SendAndClose(grpcResp)
}

func (h *fileTransferHandler) GetPresignedURL(ctx context.Context, req *pb.GetPresignedURLRequest) (*pb.GetPresignedURLResponse, error) {
	resp, err := h.fileTransferService.PresignedFile(ctx, req.FileId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get presigned URL: %v", err)
	}
	urlParsed, err := url.Parse(resp.Data.Url)
	if err != nil {
		slog.Error("[repo][google][URLParser] error while parse url err: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to parse presigned URL: %v", err)
	}
	decodePath := urlParsed.Path
	if decodePath, err = url.QueryUnescape(decodePath); err != nil {
		slog.Error("[repo][google][URLParser] error while decode url err: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to decode presigned URL: %v", err)
	}
	cleanUrl := urlParsed.Scheme + "://" + urlParsed.Host + decodePath + "?" + urlParsed.RawQuery

	grpcResp := &pb.GetPresignedURLResponse{
		Message: resp.Message,
		Status:  int32(resp.Status),
		Errors:  resp.Errors,
		Data: &pb.Data{
			Url: cleanUrl,
		},
	}
	return grpcResp, nil
}

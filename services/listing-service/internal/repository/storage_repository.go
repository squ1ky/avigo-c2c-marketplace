package repository

import (
	"context"
	"io"
	"time"
)

type StorageRepository interface {
	UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error)
	GetPresignedDownloadURL(ctx context.Context, objectID string, expiryDuration time.Duration) (string, error)
	GetPresignedUploadURL(ctx context.Context, objectID string, expiryDuration time.Duration) (string, error)
	DeleteFile(ctx context.Context, objectID string) error
	DeleteFiles(ctx context.Context, objectIDs []string) error
	GetFileInfo(ctx context.Context, objectID string) (*FileInfo, error)
	CopyFile(ctx context.Context, srcObjectID, dstObjectID string) error
}

type UploadFileRequest struct {
	ObjectID    string
	Data        io.Reader
	Size        int64
	ContentType string
	Metadata    map[string]string
}

type UploadFileResponse struct {
	ObjectID    string
	Size        int64
	ContentType string
	ETag        string
	UploadedAt  time.Time
}

type FileInfo struct {
	ObjectID     string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
	Metadata     map[string]string
}

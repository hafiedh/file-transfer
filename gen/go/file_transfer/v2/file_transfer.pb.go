// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: file_transfer/v2/file_transfer.proto

package v2

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*UploadFileRequest_FileInfo
	//	*UploadFileRequest_Chunk
	Data isUploadFileRequest_Data `protobuf_oneof:"data"`
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{0}
}

func (m *UploadFileRequest) GetData() isUploadFileRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *UploadFileRequest) GetFileInfo() *FileInfo {
	if x, ok := x.GetData().(*UploadFileRequest_FileInfo); ok {
		return x.FileInfo
	}
	return nil
}

func (x *UploadFileRequest) GetChunk() []byte {
	if x, ok := x.GetData().(*UploadFileRequest_Chunk); ok {
		return x.Chunk
	}
	return nil
}

type isUploadFileRequest_Data interface {
	isUploadFileRequest_Data()
}

type UploadFileRequest_FileInfo struct {
	FileInfo *FileInfo `protobuf:"bytes,1,opt,name=file_info,json=fileInfo,proto3,oneof"`
}

type UploadFileRequest_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}

func (*UploadFileRequest_FileInfo) isUploadFileRequest_Data() {}

func (*UploadFileRequest_Chunk) isUploadFileRequest_Data() {}

type FileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Key      string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Size     int64  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *FileInfo) Reset() {
	*x = FileInfo{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileInfo) ProtoMessage() {}

func (x *FileInfo) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileInfo.ProtoReflect.Descriptor instead.
func (*FileInfo) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{1}
}

func (x *FileInfo) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *FileInfo) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FileInfo) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

type UploadProgress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename         string  `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	BytesTransferred int64   `protobuf:"varint,2,opt,name=bytes_transferred,json=bytesTransferred,proto3" json:"bytes_transferred,omitempty"`
	TotalBytes       int64   `protobuf:"varint,3,opt,name=total_bytes,json=totalBytes,proto3" json:"total_bytes,omitempty"`
	Percentage       float64 `protobuf:"fixed64,4,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

func (x *UploadProgress) Reset() {
	*x = UploadProgress{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadProgress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadProgress) ProtoMessage() {}

func (x *UploadProgress) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadProgress.ProtoReflect.Descriptor instead.
func (*UploadProgress) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{2}
}

func (x *UploadProgress) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *UploadProgress) GetBytesTransferred() int64 {
	if x != nil {
		return x.BytesTransferred
	}
	return 0
}

func (x *UploadProgress) GetTotalBytes() int64 {
	if x != nil {
		return x.TotalBytes
	}
	return 0
}

func (x *UploadProgress) GetPercentage() float64 {
	if x != nil {
		return x.Percentage
	}
	return 0
}

type UploadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string          `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Status  int32           `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	Errors  []string        `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty"`
	Results []*UploadResult `protobuf:"bytes,4,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *UploadFileResponse) Reset() {
	*x = UploadFileResponse{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileResponse) ProtoMessage() {}

func (x *UploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileResponse.ProtoReflect.Descriptor instead.
func (*UploadFileResponse) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{3}
}

func (x *UploadFileResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *UploadFileResponse) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *UploadFileResponse) GetErrors() []string {
	if x != nil {
		return x.Errors
	}
	return nil
}

func (x *UploadFileResponse) GetResults() []*UploadResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type UploadResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Url      string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Error    string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UploadResult) Reset() {
	*x = UploadResult{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResult) ProtoMessage() {}

func (x *UploadResult) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResult.ProtoReflect.Descriptor instead.
func (*UploadResult) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{4}
}

func (x *UploadResult) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *UploadResult) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *UploadResult) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type GetPresignedURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId int64 `protobuf:"varint,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *GetPresignedURLRequest) Reset() {
	*x = GetPresignedURLRequest{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPresignedURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPresignedURLRequest) ProtoMessage() {}

func (x *GetPresignedURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPresignedURLRequest.ProtoReflect.Descriptor instead.
func (*GetPresignedURLRequest) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{5}
}

func (x *GetPresignedURLRequest) GetFileId() int64 {
	if x != nil {
		return x.FileId
	}
	return 0
}

type GetPresignedURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data    *Data    `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Message string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Errors  []string `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty"`
	Status  int32    `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetPresignedURLResponse) Reset() {
	*x = GetPresignedURLResponse{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPresignedURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPresignedURLResponse) ProtoMessage() {}

func (x *GetPresignedURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPresignedURLResponse.ProtoReflect.Descriptor instead.
func (*GetPresignedURLResponse) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{6}
}

func (x *GetPresignedURLResponse) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GetPresignedURLResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GetPresignedURLResponse) GetErrors() []string {
	if x != nil {
		return x.Errors
	}
	return nil
}

func (x *GetPresignedURLResponse) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *Data) Reset() {
	*x = Data{}
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_file_transfer_v2_file_transfer_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_file_transfer_v2_file_transfer_proto_rawDescGZIP(), []int{7}
}

func (x *Data) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_file_transfer_v2_file_transfer_proto protoreflect.FileDescriptor

var file_file_transfer_v2_file_transfer_proto_rawDesc = []byte{
	0x0a, 0x24, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2f,
	0x76, 0x32, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x22, 0x6e, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a,
	0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x2e, 0x76, 0x32, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x4c, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x9a, 0x01, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x10, 0x62, 0x79, 0x74, 0x65, 0x73, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x72,
	0x65, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x62, 0x79, 0x74, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x61, 0x67, 0x65, 0x22, 0x98, 0x01, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x12, 0x38, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x52,
	0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x31, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e,
	0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x8f, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2a, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e,
	0x76, 0x32, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x18, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x32, 0xdc, 0x01, 0x0a, 0x13, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x0a, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x23, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x68, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65,
	0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x12, 0x28, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x65, 0x74,
	0x50, 0x72, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x65, 0x73, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x12, 0x5a, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x2f, 0x76, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_transfer_v2_file_transfer_proto_rawDescOnce sync.Once
	file_file_transfer_v2_file_transfer_proto_rawDescData = file_file_transfer_v2_file_transfer_proto_rawDesc
)

func file_file_transfer_v2_file_transfer_proto_rawDescGZIP() []byte {
	file_file_transfer_v2_file_transfer_proto_rawDescOnce.Do(func() {
		file_file_transfer_v2_file_transfer_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_transfer_v2_file_transfer_proto_rawDescData)
	})
	return file_file_transfer_v2_file_transfer_proto_rawDescData
}

var file_file_transfer_v2_file_transfer_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_file_transfer_v2_file_transfer_proto_goTypes = []any{
	(*UploadFileRequest)(nil),       // 0: file_transfer.v2.UploadFileRequest
	(*FileInfo)(nil),                // 1: file_transfer.v2.FileInfo
	(*UploadProgress)(nil),          // 2: file_transfer.v2.UploadProgress
	(*UploadFileResponse)(nil),      // 3: file_transfer.v2.UploadFileResponse
	(*UploadResult)(nil),            // 4: file_transfer.v2.UploadResult
	(*GetPresignedURLRequest)(nil),  // 5: file_transfer.v2.GetPresignedURLRequest
	(*GetPresignedURLResponse)(nil), // 6: file_transfer.v2.GetPresignedURLResponse
	(*Data)(nil),                    // 7: file_transfer.v2.Data
}
var file_file_transfer_v2_file_transfer_proto_depIdxs = []int32{
	1, // 0: file_transfer.v2.UploadFileRequest.file_info:type_name -> file_transfer.v2.FileInfo
	4, // 1: file_transfer.v2.UploadFileResponse.results:type_name -> file_transfer.v2.UploadResult
	7, // 2: file_transfer.v2.GetPresignedURLResponse.data:type_name -> file_transfer.v2.Data
	0, // 3: file_transfer.v2.FileTransferService.UploadFile:input_type -> file_transfer.v2.UploadFileRequest
	5, // 4: file_transfer.v2.FileTransferService.GetPresignedURL:input_type -> file_transfer.v2.GetPresignedURLRequest
	3, // 5: file_transfer.v2.FileTransferService.UploadFile:output_type -> file_transfer.v2.UploadFileResponse
	6, // 6: file_transfer.v2.FileTransferService.GetPresignedURL:output_type -> file_transfer.v2.GetPresignedURLResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_file_transfer_v2_file_transfer_proto_init() }
func file_file_transfer_v2_file_transfer_proto_init() {
	if File_file_transfer_v2_file_transfer_proto != nil {
		return
	}
	file_file_transfer_v2_file_transfer_proto_msgTypes[0].OneofWrappers = []any{
		(*UploadFileRequest_FileInfo)(nil),
		(*UploadFileRequest_Chunk)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_file_transfer_v2_file_transfer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_file_transfer_v2_file_transfer_proto_goTypes,
		DependencyIndexes: file_file_transfer_v2_file_transfer_proto_depIdxs,
		MessageInfos:      file_file_transfer_v2_file_transfer_proto_msgTypes,
	}.Build()
	File_file_transfer_v2_file_transfer_proto = out.File
	file_file_transfer_v2_file_transfer_proto_rawDesc = nil
	file_file_transfer_v2_file_transfer_proto_goTypes = nil
	file_file_transfer_v2_file_transfer_proto_depIdxs = nil
}
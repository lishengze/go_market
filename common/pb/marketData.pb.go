// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: marketData.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Symbols defined in public import of google/protobuf/timestamp.proto.

type Timestamp = timestamppb.Timestamp

type PriceVolume struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Price  string `protobuf:"bytes,1,opt,name=price,proto3" json:"price,omitempty"`
	Volume string `protobuf:"bytes,2,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (x *PriceVolume) Reset() {
	*x = PriceVolume{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PriceVolume) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PriceVolume) ProtoMessage() {}

func (x *PriceVolume) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PriceVolume.ProtoReflect.Descriptor instead.
func (*PriceVolume) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{0}
}

func (x *PriceVolume) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *PriceVolume) GetVolume() string {
	if x != nil {
		return x.Volume
	}
	return ""
}

type Depth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp    *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 交易所更新时间
	Exchange     string                 `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol       string                 `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Asks         []*PriceVolume         `protobuf:"bytes,4,rep,name=asks,proto3" json:"asks,omitempty"`
	Bids         []*PriceVolume         `protobuf:"bytes,5,rep,name=bids,proto3" json:"bids,omitempty"`
	MpuTimestamp *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=mpu_timestamp,json=mpuTimestamp,proto3" json:"mpu_timestamp,omitempty"` // mpu 服务器时间
}

func (x *Depth) Reset() {
	*x = Depth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Depth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Depth) ProtoMessage() {}

func (x *Depth) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Depth.ProtoReflect.Descriptor instead.
func (*Depth) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{1}
}

func (x *Depth) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Depth) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *Depth) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Depth) GetAsks() []*PriceVolume {
	if x != nil {
		return x.Asks
	}
	return nil
}

func (x *Depth) GetBids() []*PriceVolume {
	if x != nil {
		return x.Bids
	}
	return nil
}

func (x *Depth) GetMpuTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.MpuTimestamp
	}
	return nil
}

type Trade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 交易所更新时间
	Exchange  string                 `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol    string                 `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Price     string                 `protobuf:"bytes,4,opt,name=price,proto3" json:"price,omitempty"`
	Volume    string                 `protobuf:"bytes,5,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (x *Trade) Reset() {
	*x = Trade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trade) ProtoMessage() {}

func (x *Trade) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trade.ProtoReflect.Descriptor instead.
func (*Trade) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{2}
}

func (x *Trade) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Trade) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *Trade) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Trade) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *Trade) GetVolume() string {
	if x != nil {
		return x.Volume
	}
	return ""
}

type Kline struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp  *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 时间
	Exchange   string                 `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Symbol     string                 `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Open       string                 `protobuf:"bytes,4,opt,name=open,proto3" json:"open,omitempty"`
	High       string                 `protobuf:"bytes,5,opt,name=high,proto3" json:"high,omitempty"`
	Low        string                 `protobuf:"bytes,6,opt,name=low,proto3" json:"low,omitempty"`
	Close      string                 `protobuf:"bytes,7,opt,name=close,proto3" json:"close,omitempty"`
	Volume     string                 `protobuf:"bytes,8,opt,name=volume,proto3" json:"volume,omitempty"` // 量
	Value      string                 `protobuf:"bytes,9,opt,name=value,proto3" json:"value,omitempty"`   // 额
	Resolution uint32                 `protobuf:"varint,10,opt,name=resolution,proto3" json:"resolution,omitempty"`
}

func (x *Kline) Reset() {
	*x = Kline{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kline) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kline) ProtoMessage() {}

func (x *Kline) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kline.ProtoReflect.Descriptor instead.
func (*Kline) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{3}
}

func (x *Kline) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Kline) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *Kline) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Kline) GetOpen() string {
	if x != nil {
		return x.Open
	}
	return ""
}

func (x *Kline) GetHigh() string {
	if x != nil {
		return x.High
	}
	return ""
}

func (x *Kline) GetLow() string {
	if x != nil {
		return x.Low
	}
	return ""
}

func (x *Kline) GetClose() string {
	if x != nil {
		return x.Close
	}
	return ""
}

func (x *Kline) GetVolume() string {
	if x != nil {
		return x.Volume
	}
	return ""
}

func (x *Kline) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *Kline) GetResolution() uint32 {
	if x != nil {
		return x.Resolution
	}
	return 0
}

type HistKlineData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol    string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Exchange  string   `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	StartTime uint64   `protobuf:"varint,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime   uint64   `protobuf:"varint,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Count     uint32   `protobuf:"varint,5,opt,name=count,proto3" json:"count,omitempty"`
	Frequency uint32   `protobuf:"varint,6,opt,name=frequency,proto3" json:"frequency,omitempty"`
	KlineData []*Kline `protobuf:"bytes,7,rep,name=kline_data,json=klineData,proto3" json:"kline_data,omitempty"`
}

func (x *HistKlineData) Reset() {
	*x = HistKlineData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistKlineData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistKlineData) ProtoMessage() {}

func (x *HistKlineData) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistKlineData.ProtoReflect.Descriptor instead.
func (*HistKlineData) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{4}
}

func (x *HistKlineData) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *HistKlineData) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *HistKlineData) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *HistKlineData) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *HistKlineData) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *HistKlineData) GetFrequency() uint32 {
	if x != nil {
		return x.Frequency
	}
	return 0
}

func (x *HistKlineData) GetKlineData() []*Kline {
	if x != nil {
		return x.KlineData
	}
	return nil
}

type EmptyReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyReq) Reset() {
	*x = EmptyReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyReq) ProtoMessage() {}

func (x *EmptyReq) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyReq.ProtoReflect.Descriptor instead.
func (*EmptyReq) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{5}
}

type EmptyRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyRsp) Reset() {
	*x = EmptyRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRsp) ProtoMessage() {}

func (x *EmptyRsp) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRsp.ProtoReflect.Descriptor instead.
func (*EmptyRsp) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{6}
}

// 增加的两个是pms,还有前置服务器的 需要用到的行情接口；
type ReqHishKlineInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol    string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Exchange  string `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	StartTime uint64 `protobuf:"varint,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime   uint64 `protobuf:"varint,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Count     uint32 `protobuf:"varint,5,opt,name=count,proto3" json:"count,omitempty"`
	Frequency uint32 `protobuf:"varint,6,opt,name=frequency,proto3" json:"frequency,omitempty"`
}

func (x *ReqHishKlineInfo) Reset() {
	*x = ReqHishKlineInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqHishKlineInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqHishKlineInfo) ProtoMessage() {}

func (x *ReqHishKlineInfo) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqHishKlineInfo.ProtoReflect.Descriptor instead.
func (*ReqHishKlineInfo) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{7}
}

func (x *ReqHishKlineInfo) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *ReqHishKlineInfo) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *ReqHishKlineInfo) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *ReqHishKlineInfo) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

func (x *ReqHishKlineInfo) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *ReqHishKlineInfo) GetFrequency() uint32 {
	if x != nil {
		return x.Frequency
	}
	return 0
}

type ReqTradeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol   string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Exchange string `protobuf:"bytes,2,opt,name=exchange,proto3" json:"exchange,omitempty"`
	Time     uint64 `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *ReqTradeInfo) Reset() {
	*x = ReqTradeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_marketData_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqTradeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqTradeInfo) ProtoMessage() {}

func (x *ReqTradeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_marketData_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqTradeInfo.ProtoReflect.Descriptor instead.
func (*ReqTradeInfo) Descriptor() ([]byte, []int) {
	return file_marketData_proto_rawDescGZIP(), []int{8}
}

func (x *ReqTradeInfo) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *ReqTradeInfo) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *ReqTradeInfo) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_marketData_proto protoreflect.FileDescriptor

var file_marketData_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x0b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x56,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x76, 0x6f, 0x6c,
	0x75, 0x6d, 0x65, 0x22, 0x9e, 0x02, 0x0a, 0x05, 0x44, 0x65, 0x70, 0x74, 0x68, 0x12, 0x38, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x32, 0x0a, 0x04, 0x61,
	0x73, 0x6b, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x04, 0x61, 0x73, 0x6b, 0x73, 0x12,
	0x32, 0x0a, 0x04, 0x62, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x2e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x56, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x52, 0x04, 0x62,
	0x69, 0x64, 0x73, 0x12, 0x3f, 0x0a, 0x0d, 0x6d, 0x70, 0x75, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x6d, 0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x22, 0xa3, 0x01, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x64, 0x65, 0x12, 0x38,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x72, 0x69,
	0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x22, 0x93, 0x02, 0x0a, 0x05, 0x4b,
	0x6c, 0x69, 0x6e, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a,
	0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f,
	0x77, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6c, 0x6f,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0xea, 0x01, 0x0a, 0x0d, 0x48, 0x69, 0x73, 0x74, 0x4b, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x66, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x79, 0x12, 0x37, 0x0a, 0x0a, 0x6b, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x4b, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x09, 0x6b, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x0a, 0x0a,
	0x08, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x22, 0x0a, 0x0a, 0x08, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x52, 0x73, 0x70, 0x22, 0xb4, 0x01, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x48, 0x69, 0x73,
	0x68, 0x4b, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x66, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x66, 0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x22, 0x56, 0x0a, 0x0c,
	0x52, 0x65, 0x71, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x32, 0xc1, 0x01, 0x0a, 0x0d, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x14, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x48, 0x69, 0x73, 0x74, 0x4b, 0x6c, 0x69, 0x6e, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x23,
	0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x2e, 0x52, 0x65, 0x71, 0x48, 0x69, 0x73, 0x68, 0x4b, 0x6c, 0x69, 0x6e, 0x65, 0x49,
	0x6e, 0x66, 0x6f, 0x1a, 0x20, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72,
	0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x4b, 0x6c, 0x69, 0x6e,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x54, 0x72, 0x61, 0x64, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x2e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x2e,
	0x52, 0x65, 0x71, 0x54, 0x72, 0x61, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x18, 0x2e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x33, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x54, 0x72, 0x61, 0x64, 0x65, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62,
	0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_marketData_proto_rawDescOnce sync.Once
	file_marketData_proto_rawDescData = file_marketData_proto_rawDesc
)

func file_marketData_proto_rawDescGZIP() []byte {
	file_marketData_proto_rawDescOnce.Do(func() {
		file_marketData_proto_rawDescData = protoimpl.X.CompressGZIP(file_marketData_proto_rawDescData)
	})
	return file_marketData_proto_rawDescData
}

var file_marketData_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_marketData_proto_goTypes = []interface{}{
	(*PriceVolume)(nil),           // 0: Proto3.MarketData.PriceVolume
	(*Depth)(nil),                 // 1: Proto3.MarketData.Depth
	(*Trade)(nil),                 // 2: Proto3.MarketData.Trade
	(*Kline)(nil),                 // 3: Proto3.MarketData.Kline
	(*HistKlineData)(nil),         // 4: Proto3.MarketData.HistKlineData
	(*EmptyReq)(nil),              // 5: Proto3.MarketData.EmptyReq
	(*EmptyRsp)(nil),              // 6: Proto3.MarketData.EmptyRsp
	(*ReqHishKlineInfo)(nil),      // 7: Proto3.MarketData.ReqHishKlineInfo
	(*ReqTradeInfo)(nil),          // 8: Proto3.MarketData.ReqTradeInfo
	(*timestamppb.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_marketData_proto_depIdxs = []int32{
	9, // 0: Proto3.MarketData.Depth.timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: Proto3.MarketData.Depth.asks:type_name -> Proto3.MarketData.PriceVolume
	0, // 2: Proto3.MarketData.Depth.bids:type_name -> Proto3.MarketData.PriceVolume
	9, // 3: Proto3.MarketData.Depth.mpu_timestamp:type_name -> google.protobuf.Timestamp
	9, // 4: Proto3.MarketData.Trade.timestamp:type_name -> google.protobuf.Timestamp
	9, // 5: Proto3.MarketData.Kline.timestamp:type_name -> google.protobuf.Timestamp
	3, // 6: Proto3.MarketData.HistKlineData.kline_data:type_name -> Proto3.MarketData.Kline
	7, // 7: Proto3.MarketData.MarketService.RequestHistKlineData:input_type -> Proto3.MarketData.ReqHishKlineInfo
	8, // 8: Proto3.MarketData.MarketService.RequestTradeData:input_type -> Proto3.MarketData.ReqTradeInfo
	4, // 9: Proto3.MarketData.MarketService.RequestHistKlineData:output_type -> Proto3.MarketData.HistKlineData
	2, // 10: Proto3.MarketData.MarketService.RequestTradeData:output_type -> Proto3.MarketData.Trade
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_marketData_proto_init() }
func file_marketData_proto_init() {
	if File_marketData_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_marketData_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PriceVolume); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Depth); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trade); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kline); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistKlineData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyRsp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqHishKlineInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_marketData_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqTradeInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_marketData_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_marketData_proto_goTypes,
		DependencyIndexes: file_marketData_proto_depIdxs,
		MessageInfos:      file_marketData_proto_msgTypes,
	}.Build()
	File_marketData_proto = out.File
	file_marketData_proto_rawDesc = nil
	file_marketData_proto_goTypes = nil
	file_marketData_proto_depIdxs = nil
}
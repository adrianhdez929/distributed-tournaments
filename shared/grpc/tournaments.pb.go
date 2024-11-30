// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: grpc/tournaments.proto

package tournaments

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

// TournamentStatus represents the current state of a tournament
type TournamentStatus int32

const (
	TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED TournamentStatus = 0
	TournamentStatus_TOURNAMENT_STATUS_IN_PROGRESS TournamentStatus = 1
	TournamentStatus_TOURNAMENT_STATUS_COMPLETED   TournamentStatus = 2
)

// Enum value maps for TournamentStatus.
var (
	TournamentStatus_name = map[int32]string{
		0: "TOURNAMENT_STATUS_NOT_STARTED",
		1: "TOURNAMENT_STATUS_IN_PROGRESS",
		2: "TOURNAMENT_STATUS_COMPLETED",
	}
	TournamentStatus_value = map[string]int32{
		"TOURNAMENT_STATUS_NOT_STARTED": 0,
		"TOURNAMENT_STATUS_IN_PROGRESS": 1,
		"TOURNAMENT_STATUS_COMPLETED":   2,
	}
)

func (x TournamentStatus) Enum() *TournamentStatus {
	p := new(TournamentStatus)
	*p = x
	return p
}

func (x TournamentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TournamentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_tournaments_proto_enumTypes[0].Descriptor()
}

func (TournamentStatus) Type() protoreflect.EnumType {
	return &file_grpc_tournaments_proto_enumTypes[0]
}

func (x TournamentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TournamentStatus.Descriptor instead.
func (TournamentStatus) EnumDescriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{0}
}

// Tournament represents a tournament entity
type Tournament struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description     string           `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	StartTimestamp  string           `protobuf:"bytes,4,opt,name=start_timestamp,json=startTimestamp,proto3" json:"start_timestamp,omitempty"` // ISO 8601 format
	Status          TournamentStatus `protobuf:"varint,5,opt,name=status,proto3,enum=tournament.TournamentStatus" json:"status,omitempty"`
	MaxParticipants int32            `protobuf:"varint,6,opt,name=max_participants,json=maxParticipants,proto3" json:"max_participants,omitempty"`
	Game            string           `protobuf:"bytes,7,opt,name=game,proto3" json:"game,omitempty"`
	Players         []*Player        `protobuf:"bytes,8,rep,name=players,proto3" json:"players,omitempty"`
}

func (x *Tournament) Reset() {
	*x = Tournament{}
	mi := &file_grpc_tournaments_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tournament) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tournament) ProtoMessage() {}

func (x *Tournament) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tournament.ProtoReflect.Descriptor instead.
func (*Tournament) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{0}
}

func (x *Tournament) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tournament) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tournament) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Tournament) GetStartTimestamp() string {
	if x != nil {
		return x.StartTimestamp
	}
	return ""
}

func (x *Tournament) GetStatus() TournamentStatus {
	if x != nil {
		return x.Status
	}
	return TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED
}

func (x *Tournament) GetMaxParticipants() int32 {
	if x != nil {
		return x.MaxParticipants
	}
	return 0
}

func (x *Tournament) GetGame() string {
	if x != nil {
		return x.Game
	}
	return ""
}

func (x *Tournament) GetPlayers() []*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

// Participant represents a tournament participant
type Player struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	AgentName string `protobuf:"bytes,3,opt,name=agent_name,json=agentName,proto3" json:"agent_name,omitempty"`
}

func (x *Player) Reset() {
	*x = Player{}
	mi := &file_grpc_tournaments_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Player) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Player) ProtoMessage() {}

func (x *Player) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Player.ProtoReflect.Descriptor instead.
func (*Player) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{1}
}

func (x *Player) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Player) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Player) GetAgentName() string {
	if x != nil {
		return x.AgentName
	}
	return ""
}

// CreateTournamentRequest is the request for creating a new tournament
type CreateTournamentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description    string    `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	StartTimestamp string    `protobuf:"bytes,3,opt,name=start_timestamp,json=startTimestamp,proto3" json:"start_timestamp,omitempty"`
	PlayerCount    int32     `protobuf:"varint,4,opt,name=player_count,json=playerCount,proto3" json:"player_count,omitempty"`
	Game           string    `protobuf:"bytes,5,opt,name=game,proto3" json:"game,omitempty"`
	Players        []*Player `protobuf:"bytes,6,rep,name=players,proto3" json:"players,omitempty"`
}

func (x *CreateTournamentRequest) Reset() {
	*x = CreateTournamentRequest{}
	mi := &file_grpc_tournaments_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTournamentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTournamentRequest) ProtoMessage() {}

func (x *CreateTournamentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTournamentRequest.ProtoReflect.Descriptor instead.
func (*CreateTournamentRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{2}
}

func (x *CreateTournamentRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateTournamentRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateTournamentRequest) GetStartTimestamp() string {
	if x != nil {
		return x.StartTimestamp
	}
	return ""
}

func (x *CreateTournamentRequest) GetPlayerCount() int32 {
	if x != nil {
		return x.PlayerCount
	}
	return 0
}

func (x *CreateTournamentRequest) GetGame() string {
	if x != nil {
		return x.Game
	}
	return ""
}

func (x *CreateTournamentRequest) GetPlayers() []*Player {
	if x != nil {
		return x.Players
	}
	return nil
}

// CreateTournamentResponse is the response for creating a new tournament
type CreateTournamentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tournament *Tournament `protobuf:"bytes,1,opt,name=tournament,proto3" json:"tournament,omitempty"`
}

func (x *CreateTournamentResponse) Reset() {
	*x = CreateTournamentResponse{}
	mi := &file_grpc_tournaments_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateTournamentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTournamentResponse) ProtoMessage() {}

func (x *CreateTournamentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTournamentResponse.ProtoReflect.Descriptor instead.
func (*CreateTournamentResponse) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{3}
}

func (x *CreateTournamentResponse) GetTournament() *Tournament {
	if x != nil {
		return x.Tournament
	}
	return nil
}

// GetTournamentRequest is the request for getting a tournament by ID
type GetTournamentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetTournamentRequest) Reset() {
	*x = GetTournamentRequest{}
	mi := &file_grpc_tournaments_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTournamentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTournamentRequest) ProtoMessage() {}

func (x *GetTournamentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTournamentRequest.ProtoReflect.Descriptor instead.
func (*GetTournamentRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{4}
}

func (x *GetTournamentRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetTournamentResponse is the response containing the requested tournament
type GetTournamentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tournament *Tournament `protobuf:"bytes,1,opt,name=tournament,proto3" json:"tournament,omitempty"`
}

func (x *GetTournamentResponse) Reset() {
	*x = GetTournamentResponse{}
	mi := &file_grpc_tournaments_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetTournamentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTournamentResponse) ProtoMessage() {}

func (x *GetTournamentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTournamentResponse.ProtoReflect.Descriptor instead.
func (*GetTournamentResponse) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{5}
}

func (x *GetTournamentResponse) GetTournament() *Tournament {
	if x != nil {
		return x.Tournament
	}
	return nil
}

// ListTournamentsRequest is the request for listing tournaments
type ListTournamentsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageSize  int32            `protobuf:"varint,1,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken string           `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	Status    TournamentStatus `protobuf:"varint,3,opt,name=status,proto3,enum=tournament.TournamentStatus" json:"status,omitempty"`
}

func (x *ListTournamentsRequest) Reset() {
	*x = ListTournamentsRequest{}
	mi := &file_grpc_tournaments_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTournamentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTournamentsRequest) ProtoMessage() {}

func (x *ListTournamentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTournamentsRequest.ProtoReflect.Descriptor instead.
func (*ListTournamentsRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{6}
}

func (x *ListTournamentsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListTournamentsRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListTournamentsRequest) GetStatus() TournamentStatus {
	if x != nil {
		return x.Status
	}
	return TournamentStatus_TOURNAMENT_STATUS_NOT_STARTED
}

// ListTournamentsResponse is the response containing a list of tournaments
type ListTournamentsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tournaments   []*Tournament `protobuf:"bytes,1,rep,name=tournaments,proto3" json:"tournaments,omitempty"`
	NextPageToken string        `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListTournamentsResponse) Reset() {
	*x = ListTournamentsResponse{}
	mi := &file_grpc_tournaments_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListTournamentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListTournamentsResponse) ProtoMessage() {}

func (x *ListTournamentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tournaments_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListTournamentsResponse.ProtoReflect.Descriptor instead.
func (*ListTournamentsResponse) Descriptor() ([]byte, []int) {
	return file_grpc_tournaments_proto_rawDescGZIP(), []int{7}
}

func (x *ListTournamentsResponse) GetTournaments() []*Tournament {
	if x != nil {
		return x.Tournaments
	}
	return nil
}

func (x *ListTournamentsResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

var File_grpc_tournaments_proto protoreflect.FileDescriptor

var file_grpc_tournaments_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x6e, 0x74, 0x22, 0x9e, 0x02, 0x0a, 0x0a, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x34, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x6d, 0x61, 0x78, 0x5f,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0f, 0x6d, 0x61, 0x78, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61,
	0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x73, 0x22, 0x4b, 0x0a, 0x06, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x22, 0xdd, 0x01, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x21, 0x0a,
	0x0c, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x67, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65,
	0x72, 0x73, 0x22, 0x52, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x75, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36,
	0x0a, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x74, 0x6f, 0x75, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x26, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4f,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x6f,
	0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x0a, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x22,
	0x8a, 0x01, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70,
	0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x34, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x7b, 0x0a, 0x17,
	0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x74, 0x6f, 0x75, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74,
	0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0b, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74,
	0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2a, 0x79, 0x0a, 0x10, 0x54, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a,
	0x1d, 0x54, 0x4f, 0x55, 0x52, 0x4e, 0x41, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x54, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x21, 0x0a, 0x1d, 0x54, 0x4f, 0x55, 0x52, 0x4e, 0x41, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45, 0x53,
	0x53, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x54, 0x4f, 0x55, 0x52, 0x4e, 0x41, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x4c, 0x45, 0x54,
	0x45, 0x44, 0x10, 0x02, 0x32, 0xaa, 0x02, 0x0a, 0x11, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x10, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23,
	0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x56, 0x0a, 0x0d, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x74,
	0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x5c, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x22, 0x2e, 0x74, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x74, 0x6f, 0x75,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x54, 0x6f, 0x75, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x14, 0x5a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x6f, 0x75, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_tournaments_proto_rawDescOnce sync.Once
	file_grpc_tournaments_proto_rawDescData = file_grpc_tournaments_proto_rawDesc
)

func file_grpc_tournaments_proto_rawDescGZIP() []byte {
	file_grpc_tournaments_proto_rawDescOnce.Do(func() {
		file_grpc_tournaments_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_tournaments_proto_rawDescData)
	})
	return file_grpc_tournaments_proto_rawDescData
}

var file_grpc_tournaments_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_grpc_tournaments_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_grpc_tournaments_proto_goTypes = []any{
	(TournamentStatus)(0),            // 0: tournament.TournamentStatus
	(*Tournament)(nil),               // 1: tournament.Tournament
	(*Player)(nil),                   // 2: tournament.Player
	(*CreateTournamentRequest)(nil),  // 3: tournament.CreateTournamentRequest
	(*CreateTournamentResponse)(nil), // 4: tournament.CreateTournamentResponse
	(*GetTournamentRequest)(nil),     // 5: tournament.GetTournamentRequest
	(*GetTournamentResponse)(nil),    // 6: tournament.GetTournamentResponse
	(*ListTournamentsRequest)(nil),   // 7: tournament.ListTournamentsRequest
	(*ListTournamentsResponse)(nil),  // 8: tournament.ListTournamentsResponse
}
var file_grpc_tournaments_proto_depIdxs = []int32{
	0,  // 0: tournament.Tournament.status:type_name -> tournament.TournamentStatus
	2,  // 1: tournament.Tournament.players:type_name -> tournament.Player
	2,  // 2: tournament.CreateTournamentRequest.players:type_name -> tournament.Player
	1,  // 3: tournament.CreateTournamentResponse.tournament:type_name -> tournament.Tournament
	1,  // 4: tournament.GetTournamentResponse.tournament:type_name -> tournament.Tournament
	0,  // 5: tournament.ListTournamentsRequest.status:type_name -> tournament.TournamentStatus
	1,  // 6: tournament.ListTournamentsResponse.tournaments:type_name -> tournament.Tournament
	3,  // 7: tournament.TournamentService.CreateTournament:input_type -> tournament.CreateTournamentRequest
	5,  // 8: tournament.TournamentService.GetTournament:input_type -> tournament.GetTournamentRequest
	7,  // 9: tournament.TournamentService.ListTournaments:input_type -> tournament.ListTournamentsRequest
	4,  // 10: tournament.TournamentService.CreateTournament:output_type -> tournament.CreateTournamentResponse
	6,  // 11: tournament.TournamentService.GetTournament:output_type -> tournament.GetTournamentResponse
	8,  // 12: tournament.TournamentService.ListTournaments:output_type -> tournament.ListTournamentsResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_grpc_tournaments_proto_init() }
func file_grpc_tournaments_proto_init() {
	if File_grpc_tournaments_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_tournaments_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_tournaments_proto_goTypes,
		DependencyIndexes: file_grpc_tournaments_proto_depIdxs,
		EnumInfos:         file_grpc_tournaments_proto_enumTypes,
		MessageInfos:      file_grpc_tournaments_proto_msgTypes,
	}.Build()
	File_grpc_tournaments_proto = out.File
	file_grpc_tournaments_proto_rawDesc = nil
	file_grpc_tournaments_proto_goTypes = nil
	file_grpc_tournaments_proto_depIdxs = nil
}

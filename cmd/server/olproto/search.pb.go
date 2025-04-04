// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v5.28.0
// source: search.proto

package olproto

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

type ProtoTagsCategory int32

const (
	ProtoTagsCategory_OTHER    ProtoTagsCategory = 0
	ProtoTagsCategory_REL      ProtoTagsCategory = 1
	ProtoTagsCategory_REL_TYPE ProtoTagsCategory = 2
	ProtoTagsCategory_FANDOM   ProtoTagsCategory = 3
	ProtoTagsCategory_WARNING  ProtoTagsCategory = 4
)

// Enum value maps for ProtoTagsCategory.
var (
	ProtoTagsCategory_name = map[int32]string{
		0: "OTHER",
		1: "REL",
		2: "REL_TYPE",
		3: "FANDOM",
		4: "WARNING",
	}
	ProtoTagsCategory_value = map[string]int32{
		"OTHER":    0,
		"REL":      1,
		"REL_TYPE": 2,
		"FANDOM":   3,
		"WARNING":  4,
	}
)

func (x ProtoTagsCategory) Enum() *ProtoTagsCategory {
	p := new(ProtoTagsCategory)
	*p = x
	return p
}

func (x ProtoTagsCategory) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProtoTagsCategory) Descriptor() protoreflect.EnumDescriptor {
	return file_search_proto_enumTypes[0].Descriptor()
}

func (ProtoTagsCategory) Type() protoreflect.EnumType {
	return &file_search_proto_enumTypes[0]
}

func (x ProtoTagsCategory) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProtoTagsCategory.Descriptor instead.
func (ProtoTagsCategory) EnumDescriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{0}
}

type ProtoAgeRating int32

const (
	ProtoAgeRating_UNKNOWN ProtoAgeRating = 0
	ProtoAgeRating_G       ProtoAgeRating = 1
	ProtoAgeRating_PG      ProtoAgeRating = 2
	ProtoAgeRating_PG13    ProtoAgeRating = 3
	ProtoAgeRating_R       ProtoAgeRating = 4
	ProtoAgeRating_NC17    ProtoAgeRating = 5
)

// Enum value maps for ProtoAgeRating.
var (
	ProtoAgeRating_name = map[int32]string{
		0: "UNKNOWN",
		1: "G",
		2: "PG",
		3: "PG13",
		4: "R",
		5: "NC17",
	}
	ProtoAgeRating_value = map[string]int32{
		"UNKNOWN": 0,
		"G":       1,
		"PG":      2,
		"PG13":    3,
		"R":       4,
		"NC17":    5,
	}
)

func (x ProtoAgeRating) Enum() *ProtoAgeRating {
	p := new(ProtoAgeRating)
	*p = x
	return p
}

func (x ProtoAgeRating) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProtoAgeRating) Descriptor() protoreflect.EnumDescriptor {
	return file_search_proto_enumTypes[1].Descriptor()
}

func (ProtoAgeRating) Type() protoreflect.EnumType {
	return &file_search_proto_enumTypes[1]
}

func (x ProtoAgeRating) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProtoAgeRating.Descriptor instead.
func (ProtoAgeRating) EnumDescriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{1}
}

type ProtoDefinedTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	IsAdult     bool              `protobuf:"varint,4,opt,name=is_adult,json=isAdult,proto3" json:"is_adult,omitempty"`
	IsSpoiler   bool              `protobuf:"varint,5,opt,name=is_spoiler,json=isSpoiler,proto3" json:"is_spoiler,omitempty"`
	Category    ProtoTagsCategory `protobuf:"varint,6,opt,name=category,proto3,enum=ProtoTagsCategory" json:"category,omitempty"`
}

func (x *ProtoDefinedTag) Reset() {
	*x = ProtoDefinedTag{}
	mi := &file_search_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProtoDefinedTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoDefinedTag) ProtoMessage() {}

func (x *ProtoDefinedTag) ProtoReflect() protoreflect.Message {
	mi := &file_search_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoDefinedTag.ProtoReflect.Descriptor instead.
func (*ProtoDefinedTag) Descriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{0}
}

func (x *ProtoDefinedTag) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ProtoDefinedTag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProtoDefinedTag) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ProtoDefinedTag) GetIsAdult() bool {
	if x != nil {
		return x.IsAdult
	}
	return false
}

func (x *ProtoDefinedTag) GetIsSpoiler() bool {
	if x != nil {
		return x.IsSpoiler
	}
	return false
}

func (x *ProtoDefinedTag) GetCategory() ProtoTagsCategory {
	if x != nil {
		return x.Category
	}
	return ProtoTagsCategory_OTHER
}

type ProtoBookSearchItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name       string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	TagIds     []int64        `protobuf:"varint,3,rep,packed,name=tag_ids,json=tagIds,proto3" json:"tag_ids,omitempty"`
	Cover      string         `protobuf:"bytes,4,opt,name=cover,proto3" json:"cover,omitempty"`
	Words      uint32         `protobuf:"varint,5,opt,name=words,proto3" json:"words,omitempty"`
	Chapters   uint32         `protobuf:"varint,6,opt,name=chapters,proto3" json:"chapters,omitempty"`
	Favorites  uint32         `protobuf:"varint,7,opt,name=favorites,proto3" json:"favorites,omitempty"`
	AuthorName string         `protobuf:"bytes,8,opt,name=author_name,json=authorName,proto3" json:"author_name,omitempty"`
	AuthorId   string         `protobuf:"bytes,9,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	AgeRating  ProtoAgeRating `protobuf:"varint,10,opt,name=age_rating,json=ageRating,proto3,enum=ProtoAgeRating" json:"age_rating,omitempty"`
	Summary    string         `protobuf:"bytes,11,opt,name=summary,proto3" json:"summary,omitempty"`
	CreatedAt  uint32         `protobuf:"varint,12,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt  uint32         `protobuf:"varint,13,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *ProtoBookSearchItem) Reset() {
	*x = ProtoBookSearchItem{}
	mi := &file_search_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProtoBookSearchItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoBookSearchItem) ProtoMessage() {}

func (x *ProtoBookSearchItem) ProtoReflect() protoreflect.Message {
	mi := &file_search_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoBookSearchItem.ProtoReflect.Descriptor instead.
func (*ProtoBookSearchItem) Descriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{1}
}

func (x *ProtoBookSearchItem) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ProtoBookSearchItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProtoBookSearchItem) GetTagIds() []int64 {
	if x != nil {
		return x.TagIds
	}
	return nil
}

func (x *ProtoBookSearchItem) GetCover() string {
	if x != nil {
		return x.Cover
	}
	return ""
}

func (x *ProtoBookSearchItem) GetWords() uint32 {
	if x != nil {
		return x.Words
	}
	return 0
}

func (x *ProtoBookSearchItem) GetChapters() uint32 {
	if x != nil {
		return x.Chapters
	}
	return 0
}

func (x *ProtoBookSearchItem) GetFavorites() uint32 {
	if x != nil {
		return x.Favorites
	}
	return 0
}

func (x *ProtoBookSearchItem) GetAuthorName() string {
	if x != nil {
		return x.AuthorName
	}
	return ""
}

func (x *ProtoBookSearchItem) GetAuthorId() string {
	if x != nil {
		return x.AuthorId
	}
	return ""
}

func (x *ProtoBookSearchItem) GetAgeRating() ProtoAgeRating {
	if x != nil {
		return x.AgeRating
	}
	return ProtoAgeRating_UNKNOWN
}

func (x *ProtoBookSearchItem) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *ProtoBookSearchItem) GetCreatedAt() uint32 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *ProtoBookSearchItem) GetUpdatedAt() uint32 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type ProtoSearchResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items      []*ProtoBookSearchItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Tags       []*ProtoDefinedTag     `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags,omitempty"`
	TotalPages uint32                 `protobuf:"varint,3,opt,name=total_pages,json=totalPages,proto3" json:"total_pages,omitempty"`
	Page       uint32                 `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
	PageSize   uint32                 `protobuf:"varint,5,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	Took       uint32                 `protobuf:"varint,6,opt,name=took,proto3" json:"took,omitempty"`
	CacheKey   string                 `protobuf:"bytes,7,opt,name=cache_key,json=cacheKey,proto3" json:"cache_key,omitempty"`
	CacheTook  uint32                 `protobuf:"varint,8,opt,name=cache_took,json=cacheTook,proto3" json:"cache_took,omitempty"`
	CacheHit   bool                   `protobuf:"varint,9,opt,name=cache_hit,json=cacheHit,proto3" json:"cache_hit,omitempty"`
	Filter     *ProtoSearchFilter     `protobuf:"bytes,10,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *ProtoSearchResult) Reset() {
	*x = ProtoSearchResult{}
	mi := &file_search_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProtoSearchResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoSearchResult) ProtoMessage() {}

func (x *ProtoSearchResult) ProtoReflect() protoreflect.Message {
	mi := &file_search_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoSearchResult.ProtoReflect.Descriptor instead.
func (*ProtoSearchResult) Descriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{2}
}

func (x *ProtoSearchResult) GetItems() []*ProtoBookSearchItem {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ProtoSearchResult) GetTags() []*ProtoDefinedTag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *ProtoSearchResult) GetTotalPages() uint32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

func (x *ProtoSearchResult) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ProtoSearchResult) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ProtoSearchResult) GetTook() uint32 {
	if x != nil {
		return x.Took
	}
	return 0
}

func (x *ProtoSearchResult) GetCacheKey() string {
	if x != nil {
		return x.CacheKey
	}
	return ""
}

func (x *ProtoSearchResult) GetCacheTook() uint32 {
	if x != nil {
		return x.CacheTook
	}
	return 0
}

func (x *ProtoSearchResult) GetCacheHit() bool {
	if x != nil {
		return x.CacheHit
	}
	return false
}

func (x *ProtoSearchResult) GetFilter() *ProtoSearchFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type ProtoSearchFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WordsMin           *int32   `protobuf:"varint,1,opt,name=words_min,json=wordsMin,proto3,oneof" json:"words_min,omitempty"`
	WordsMax           *int32   `protobuf:"varint,2,opt,name=words_max,json=wordsMax,proto3,oneof" json:"words_max,omitempty"`
	ChaptersMin        *int32   `protobuf:"varint,3,opt,name=chapters_min,json=chaptersMin,proto3,oneof" json:"chapters_min,omitempty"`
	ChaptersMax        *int32   `protobuf:"varint,4,opt,name=chapters_max,json=chaptersMax,proto3,oneof" json:"chapters_max,omitempty"`
	WordsPerChapterMin *int32   `protobuf:"varint,5,opt,name=words_per_chapter_min,json=wordsPerChapterMin,proto3,oneof" json:"words_per_chapter_min,omitempty"`
	WordsPerChapterMax *int32   `protobuf:"varint,6,opt,name=words_per_chapter_max,json=wordsPerChapterMax,proto3,oneof" json:"words_per_chapter_max,omitempty"`
	FavoritesMin       *int32   `protobuf:"varint,7,opt,name=favorites_min,json=favoritesMin,proto3,oneof" json:"favorites_min,omitempty"`
	FavoritesMax       *int32   `protobuf:"varint,8,opt,name=favorites_max,json=favoritesMax,proto3,oneof" json:"favorites_max,omitempty"`
	IncludeTags        []int64  `protobuf:"varint,9,rep,packed,name=include_tags,json=includeTags,proto3" json:"include_tags,omitempty"`
	ExcludeTags        []int64  `protobuf:"varint,10,rep,packed,name=exclude_tags,json=excludeTags,proto3" json:"exclude_tags,omitempty"`
	IncludeUsers       []string `protobuf:"bytes,11,rep,name=include_users,json=includeUsers,proto3" json:"include_users,omitempty"`
	ExcludeUsers       []string `protobuf:"bytes,12,rep,name=exclude_users,json=excludeUsers,proto3" json:"exclude_users,omitempty"`
}

func (x *ProtoSearchFilter) Reset() {
	*x = ProtoSearchFilter{}
	mi := &file_search_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProtoSearchFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtoSearchFilter) ProtoMessage() {}

func (x *ProtoSearchFilter) ProtoReflect() protoreflect.Message {
	mi := &file_search_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtoSearchFilter.ProtoReflect.Descriptor instead.
func (*ProtoSearchFilter) Descriptor() ([]byte, []int) {
	return file_search_proto_rawDescGZIP(), []int{3}
}

func (x *ProtoSearchFilter) GetWordsMin() int32 {
	if x != nil && x.WordsMin != nil {
		return *x.WordsMin
	}
	return 0
}

func (x *ProtoSearchFilter) GetWordsMax() int32 {
	if x != nil && x.WordsMax != nil {
		return *x.WordsMax
	}
	return 0
}

func (x *ProtoSearchFilter) GetChaptersMin() int32 {
	if x != nil && x.ChaptersMin != nil {
		return *x.ChaptersMin
	}
	return 0
}

func (x *ProtoSearchFilter) GetChaptersMax() int32 {
	if x != nil && x.ChaptersMax != nil {
		return *x.ChaptersMax
	}
	return 0
}

func (x *ProtoSearchFilter) GetWordsPerChapterMin() int32 {
	if x != nil && x.WordsPerChapterMin != nil {
		return *x.WordsPerChapterMin
	}
	return 0
}

func (x *ProtoSearchFilter) GetWordsPerChapterMax() int32 {
	if x != nil && x.WordsPerChapterMax != nil {
		return *x.WordsPerChapterMax
	}
	return 0
}

func (x *ProtoSearchFilter) GetFavoritesMin() int32 {
	if x != nil && x.FavoritesMin != nil {
		return *x.FavoritesMin
	}
	return 0
}

func (x *ProtoSearchFilter) GetFavoritesMax() int32 {
	if x != nil && x.FavoritesMax != nil {
		return *x.FavoritesMax
	}
	return 0
}

func (x *ProtoSearchFilter) GetIncludeTags() []int64 {
	if x != nil {
		return x.IncludeTags
	}
	return nil
}

func (x *ProtoSearchFilter) GetExcludeTags() []int64 {
	if x != nil {
		return x.ExcludeTags
	}
	return nil
}

func (x *ProtoSearchFilter) GetIncludeUsers() []string {
	if x != nil {
		return x.IncludeUsers
	}
	return nil
}

func (x *ProtoSearchFilter) GetExcludeUsers() []string {
	if x != nil {
		return x.ExcludeUsers
	}
	return nil
}

var File_search_proto protoreflect.FileDescriptor

var file_search_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1,
	0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x54,
	0x61, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x61,
	0x64, 0x75, 0x6c, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x41, 0x64,
	0x75, 0x6c, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x73, 0x70, 0x6f, 0x69, 0x6c, 0x65,
	0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x70, 0x6f, 0x69, 0x6c,
	0x65, 0x72, 0x12, 0x2e, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x54, 0x61, 0x67, 0x73,
	0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x22, 0xfe, 0x02, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x6f, 0x6f, 0x6b,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x74, 0x61, 0x67, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x06, 0x74, 0x61, 0x67, 0x49, 0x64, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x77, 0x6f,
	0x72, 0x64, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x1c, 0x0a, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x12, 0x1f, 0x0a,
	0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x0a, 0x61,
	0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0f, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x09, 0x61, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75,
	0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x22, 0xd0, 0x02, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x42, 0x6f, 0x6f, 0x6b, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x24, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x44, 0x65, 0x66, 0x69, 0x6e,
	0x65, 0x64, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x6f, 0x6f, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x6f, 0x6f,
	0x6b, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x63, 0x68, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x74, 0x6f, 0x6f, 0x6b, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x54, 0x6f, 0x6f, 0x6b, 0x12, 0x1b, 0x0a,
	0x09, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x68, 0x69, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x63, 0x61, 0x63, 0x68, 0x65, 0x48, 0x69, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x91, 0x05, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x09,
	0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x00, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x20,
	0x0a, 0x09, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x01, 0x52, 0x08, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x26, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x5f, 0x6d, 0x69, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x48, 0x02, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65,
	0x72, 0x73, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x70,
	0x74, 0x65, 0x72, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x48, 0x03,
	0x52, 0x0b, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x36, 0x0a, 0x15, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x63, 0x68,
	0x61, 0x70, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x69, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x04, 0x52, 0x12, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x50, 0x65, 0x72, 0x43, 0x68, 0x61, 0x70, 0x74,
	0x65, 0x72, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x36, 0x0a, 0x15, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x61,
	0x78, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x48, 0x05, 0x52, 0x12, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x50, 0x65, 0x72, 0x43, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x78, 0x88, 0x01, 0x01,
	0x12, 0x28, 0x0a, 0x0d, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x5f, 0x6d, 0x69,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x48, 0x06, 0x52, 0x0c, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x73, 0x4d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x28, 0x0a, 0x0d, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x07, 0x52, 0x0c, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x4d, 0x61,
	0x78, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f,
	0x74, 0x61, 0x67, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x69, 0x6e, 0x63, 0x6c,
	0x75, 0x64, 0x65, 0x54, 0x61, 0x67, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x78, 0x63, 0x6c, 0x75,
	0x64, 0x65, 0x5f, 0x74, 0x61, 0x67, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x65,
	0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x54, 0x61, 0x67, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x6e,
	0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x0c, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x23, 0x0a, 0x0d, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x18, 0x0c, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x6d,
	0x69, 0x6e, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x6d, 0x61, 0x78,
	0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x5f, 0x6d, 0x69,
	0x6e, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x5f, 0x6d,
	0x61, 0x78, 0x42, 0x18, 0x0a, 0x16, 0x5f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x70, 0x65, 0x72,
	0x5f, 0x63, 0x68, 0x61, 0x70, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x69, 0x6e, 0x42, 0x18, 0x0a, 0x16,
	0x5f, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x63, 0x68, 0x61, 0x70, 0x74,
	0x65, 0x72, 0x5f, 0x6d, 0x61, 0x78, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x73, 0x5f, 0x6d, 0x69, 0x6e, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x66, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x5f, 0x6d, 0x61, 0x78, 0x2a, 0x4e, 0x0a, 0x11, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x54, 0x61, 0x67, 0x73, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12,
	0x09, 0x0a, 0x05, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x52, 0x45,
	0x4c, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4c, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10,
	0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x4e, 0x44, 0x4f, 0x4d, 0x10, 0x03, 0x12, 0x0b, 0x0a,
	0x07, 0x57, 0x41, 0x52, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x2a, 0x47, 0x0a, 0x0e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x41, 0x67, 0x65, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x0b, 0x0a, 0x07,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x05, 0x0a, 0x01, 0x47, 0x10, 0x01,
	0x12, 0x06, 0x0a, 0x02, 0x50, 0x47, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x47, 0x31, 0x33,
	0x10, 0x03, 0x12, 0x05, 0x0a, 0x01, 0x52, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x43, 0x31,
	0x37, 0x10, 0x05, 0x42, 0x14, 0x5a, 0x12, 0x63, 0x6d, 0x64, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x6f, 0x6c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_search_proto_rawDescOnce sync.Once
	file_search_proto_rawDescData = file_search_proto_rawDesc
)

func file_search_proto_rawDescGZIP() []byte {
	file_search_proto_rawDescOnce.Do(func() {
		file_search_proto_rawDescData = protoimpl.X.CompressGZIP(file_search_proto_rawDescData)
	})
	return file_search_proto_rawDescData
}

var file_search_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_search_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_search_proto_goTypes = []any{
	(ProtoTagsCategory)(0),      // 0: ProtoTagsCategory
	(ProtoAgeRating)(0),         // 1: ProtoAgeRating
	(*ProtoDefinedTag)(nil),     // 2: ProtoDefinedTag
	(*ProtoBookSearchItem)(nil), // 3: ProtoBookSearchItem
	(*ProtoSearchResult)(nil),   // 4: ProtoSearchResult
	(*ProtoSearchFilter)(nil),   // 5: ProtoSearchFilter
}
var file_search_proto_depIdxs = []int32{
	0, // 0: ProtoDefinedTag.category:type_name -> ProtoTagsCategory
	1, // 1: ProtoBookSearchItem.age_rating:type_name -> ProtoAgeRating
	3, // 2: ProtoSearchResult.items:type_name -> ProtoBookSearchItem
	2, // 3: ProtoSearchResult.tags:type_name -> ProtoDefinedTag
	5, // 4: ProtoSearchResult.filter:type_name -> ProtoSearchFilter
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_search_proto_init() }
func file_search_proto_init() {
	if File_search_proto != nil {
		return
	}
	file_search_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_search_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_search_proto_goTypes,
		DependencyIndexes: file_search_proto_depIdxs,
		EnumInfos:         file_search_proto_enumTypes,
		MessageInfos:      file_search_proto_msgTypes,
	}.Build()
	File_search_proto = out.File
	file_search_proto_rawDesc = nil
	file_search_proto_goTypes = nil
	file_search_proto_depIdxs = nil
}

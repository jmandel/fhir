//    Copyright 2019 Google Inc.
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        https://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: proto/google/fhir/proto/r4/core/resources/related_person.proto

package related_person_go_proto

import (
	any "github.com/golang/protobuf/ptypes/any"
	_ "github.com/google/fhir/go/proto/google/fhir/proto/annotations_go_proto"
	codes_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/codes_go_proto"
	datatypes_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/datatypes_go_proto"
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

// Auto-generated from StructureDefinition for RelatedPerson, last updated
// 2019-11-01T09:29:23.356+11:00. A person that is related to a patient, but who
// is not a direct target of care. See
// http://hl7.org/fhir/StructureDefinition/RelatedPerson
type RelatedPerson struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Logical id of this artifact
	Id *datatypes_go_proto.Id `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Metadata about the resource
	Meta *datatypes_go_proto.Meta `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	// A set of rules under which this content was created
	ImplicitRules *datatypes_go_proto.Uri `protobuf:"bytes,3,opt,name=implicit_rules,json=implicitRules,proto3" json:"implicit_rules,omitempty"`
	// Language of the resource content
	Language *datatypes_go_proto.Code `protobuf:"bytes,4,opt,name=language,proto3" json:"language,omitempty"`
	// Text summary of the resource, for human interpretation
	Text *datatypes_go_proto.Narrative `protobuf:"bytes,5,opt,name=text,proto3" json:"text,omitempty"`
	// Contained, inline Resources
	Contained []*any.Any `protobuf:"bytes,6,rep,name=contained,proto3" json:"contained,omitempty"`
	// Additional content defined by implementations
	Extension []*datatypes_go_proto.Extension `protobuf:"bytes,8,rep,name=extension,proto3" json:"extension,omitempty"`
	// Extensions that cannot be ignored
	ModifierExtension []*datatypes_go_proto.Extension `protobuf:"bytes,9,rep,name=modifier_extension,json=modifierExtension,proto3" json:"modifier_extension,omitempty"`
	// A human identifier for this person
	Identifier []*datatypes_go_proto.Identifier `protobuf:"bytes,10,rep,name=identifier,proto3" json:"identifier,omitempty"`
	// Whether this related person's record is in active use
	Active *datatypes_go_proto.Boolean `protobuf:"bytes,11,opt,name=active,proto3" json:"active,omitempty"`
	// The patient this person is related to
	Patient *datatypes_go_proto.Reference `protobuf:"bytes,12,opt,name=patient,proto3" json:"patient,omitempty"`
	// The nature of the relationship
	Relationship []*datatypes_go_proto.CodeableConcept `protobuf:"bytes,13,rep,name=relationship,proto3" json:"relationship,omitempty"`
	// A name associated with the person
	Name []*datatypes_go_proto.HumanName `protobuf:"bytes,14,rep,name=name,proto3" json:"name,omitempty"`
	// A contact detail for the person
	Telecom []*datatypes_go_proto.ContactPoint `protobuf:"bytes,15,rep,name=telecom,proto3" json:"telecom,omitempty"`
	Gender  *RelatedPerson_GenderCode          `protobuf:"bytes,16,opt,name=gender,proto3" json:"gender,omitempty"`
	// The date on which the related person was born
	BirthDate *datatypes_go_proto.Date `protobuf:"bytes,17,opt,name=birth_date,json=birthDate,proto3" json:"birth_date,omitempty"`
	// Address where the related person can be contacted or visited
	Address []*datatypes_go_proto.Address `protobuf:"bytes,18,rep,name=address,proto3" json:"address,omitempty"`
	// Image of the person
	Photo []*datatypes_go_proto.Attachment `protobuf:"bytes,19,rep,name=photo,proto3" json:"photo,omitempty"`
	// Period of time that this relationship is considered valid
	Period        *datatypes_go_proto.Period     `protobuf:"bytes,20,opt,name=period,proto3" json:"period,omitempty"`
	Communication []*RelatedPerson_Communication `protobuf:"bytes,21,rep,name=communication,proto3" json:"communication,omitempty"`
}

func (x *RelatedPerson) Reset() {
	*x = RelatedPerson{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelatedPerson) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelatedPerson) ProtoMessage() {}

func (x *RelatedPerson) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelatedPerson.ProtoReflect.Descriptor instead.
func (*RelatedPerson) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescGZIP(), []int{0}
}

func (x *RelatedPerson) GetId() *datatypes_go_proto.Id {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *RelatedPerson) GetMeta() *datatypes_go_proto.Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *RelatedPerson) GetImplicitRules() *datatypes_go_proto.Uri {
	if x != nil {
		return x.ImplicitRules
	}
	return nil
}

func (x *RelatedPerson) GetLanguage() *datatypes_go_proto.Code {
	if x != nil {
		return x.Language
	}
	return nil
}

func (x *RelatedPerson) GetText() *datatypes_go_proto.Narrative {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *RelatedPerson) GetContained() []*any.Any {
	if x != nil {
		return x.Contained
	}
	return nil
}

func (x *RelatedPerson) GetExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.Extension
	}
	return nil
}

func (x *RelatedPerson) GetModifierExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.ModifierExtension
	}
	return nil
}

func (x *RelatedPerson) GetIdentifier() []*datatypes_go_proto.Identifier {
	if x != nil {
		return x.Identifier
	}
	return nil
}

func (x *RelatedPerson) GetActive() *datatypes_go_proto.Boolean {
	if x != nil {
		return x.Active
	}
	return nil
}

func (x *RelatedPerson) GetPatient() *datatypes_go_proto.Reference {
	if x != nil {
		return x.Patient
	}
	return nil
}

func (x *RelatedPerson) GetRelationship() []*datatypes_go_proto.CodeableConcept {
	if x != nil {
		return x.Relationship
	}
	return nil
}

func (x *RelatedPerson) GetName() []*datatypes_go_proto.HumanName {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *RelatedPerson) GetTelecom() []*datatypes_go_proto.ContactPoint {
	if x != nil {
		return x.Telecom
	}
	return nil
}

func (x *RelatedPerson) GetGender() *RelatedPerson_GenderCode {
	if x != nil {
		return x.Gender
	}
	return nil
}

func (x *RelatedPerson) GetBirthDate() *datatypes_go_proto.Date {
	if x != nil {
		return x.BirthDate
	}
	return nil
}

func (x *RelatedPerson) GetAddress() []*datatypes_go_proto.Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *RelatedPerson) GetPhoto() []*datatypes_go_proto.Attachment {
	if x != nil {
		return x.Photo
	}
	return nil
}

func (x *RelatedPerson) GetPeriod() *datatypes_go_proto.Period {
	if x != nil {
		return x.Period
	}
	return nil
}

func (x *RelatedPerson) GetCommunication() []*RelatedPerson_Communication {
	if x != nil {
		return x.Communication
	}
	return nil
}

// male | female | other | unknown
type RelatedPerson_GenderCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value     codes_go_proto.AdministrativeGenderCode_Value `protobuf:"varint,1,opt,name=value,proto3,enum=google.fhir.r4.core.AdministrativeGenderCode_Value" json:"value,omitempty"`
	Id        *datatypes_go_proto.String                    `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Extension []*datatypes_go_proto.Extension               `protobuf:"bytes,3,rep,name=extension,proto3" json:"extension,omitempty"`
}

func (x *RelatedPerson_GenderCode) Reset() {
	*x = RelatedPerson_GenderCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelatedPerson_GenderCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelatedPerson_GenderCode) ProtoMessage() {}

func (x *RelatedPerson_GenderCode) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelatedPerson_GenderCode.ProtoReflect.Descriptor instead.
func (*RelatedPerson_GenderCode) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescGZIP(), []int{0, 0}
}

func (x *RelatedPerson_GenderCode) GetValue() codes_go_proto.AdministrativeGenderCode_Value {
	if x != nil {
		return x.Value
	}
	return codes_go_proto.AdministrativeGenderCode_INVALID_UNINITIALIZED
}

func (x *RelatedPerson_GenderCode) GetId() *datatypes_go_proto.String {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *RelatedPerson_GenderCode) GetExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.Extension
	}
	return nil
}

// A language which may be used to communicate with about the patient's health
type RelatedPerson_Communication struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique id for inter-element referencing
	Id *datatypes_go_proto.String `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Additional content defined by implementations
	Extension []*datatypes_go_proto.Extension `protobuf:"bytes,2,rep,name=extension,proto3" json:"extension,omitempty"`
	// Extensions that cannot be ignored even if unrecognized
	ModifierExtension []*datatypes_go_proto.Extension `protobuf:"bytes,3,rep,name=modifier_extension,json=modifierExtension,proto3" json:"modifier_extension,omitempty"`
	// The language which can be used to communicate with the patient about his
	// or her health
	Language *datatypes_go_proto.CodeableConcept `protobuf:"bytes,4,opt,name=language,proto3" json:"language,omitempty"`
	// Language preference indicator
	Preferred *datatypes_go_proto.Boolean `protobuf:"bytes,5,opt,name=preferred,proto3" json:"preferred,omitempty"`
}

func (x *RelatedPerson_Communication) Reset() {
	*x = RelatedPerson_Communication{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelatedPerson_Communication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelatedPerson_Communication) ProtoMessage() {}

func (x *RelatedPerson_Communication) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelatedPerson_Communication.ProtoReflect.Descriptor instead.
func (*RelatedPerson_Communication) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescGZIP(), []int{0, 1}
}

func (x *RelatedPerson_Communication) GetId() *datatypes_go_proto.String {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *RelatedPerson_Communication) GetExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.Extension
	}
	return nil
}

func (x *RelatedPerson_Communication) GetModifierExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.ModifierExtension
	}
	return nil
}

func (x *RelatedPerson_Communication) GetLanguage() *datatypes_go_proto.CodeableConcept {
	if x != nil {
		return x.Language
	}
	return nil
}

func (x *RelatedPerson_Communication) GetPreferred() *datatypes_go_proto.Boolean {
	if x != nil {
		return x.Preferred
	}
	return nil
}

var File_proto_google_fhir_proto_r4_core_resources_related_person_proto protoreflect.FileDescriptor

var file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDesc = []byte{
	0x0a, 0x3e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66,
	0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x34, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x72, 0x65, 0x6c, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x13, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x29, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66,
	0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x34, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x64,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x72, 0x34, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x0f, 0x0a, 0x0d, 0x52, 0x65,
	0x6c, 0x61, 0x74, 0x65, 0x64, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x64,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72,
	0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d,
	0x65, 0x74, 0x61, 0x12, 0x3f, 0x0a, 0x0e, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x5f,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x55, 0x72, 0x69, 0x52, 0x0d, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x52,
	0x75, 0x6c, 0x65, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x4e, 0x61, 0x72, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x32, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x64, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x64, 0x12, 0x3c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x4d, 0x0a, 0x12, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x5f, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x6d,
	0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x3f, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68,
	0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x12, 0x34, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e,
	0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x12, 0x4d, 0x0a, 0x07, 0x70, 0x61, 0x74, 0x69, 0x65,
	0x6e, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x13, 0xf0, 0xd0, 0x87, 0xeb, 0x04, 0x01,
	0xf2, 0xff, 0xfc, 0xc2, 0x06, 0x07, 0x50, 0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x70,
	0x61, 0x74, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x48, 0x0a, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x63, 0x65,
	0x70, 0x74, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70,
	0x12, 0x32, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x07, 0x74, 0x65, 0x6c, 0x65, 0x63, 0x6f, 0x6d, 0x18,
	0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66,
	0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x63, 0x74, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x07, 0x74, 0x65, 0x6c, 0x65, 0x63, 0x6f,
	0x6d, 0x12, 0x45, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e,
	0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x50,
	0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x38, 0x0a, 0x0a, 0x62, 0x69, 0x72, 0x74,
	0x68, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x52, 0x09, 0x62, 0x69, 0x72, 0x74, 0x68, 0x44, 0x61,
	0x74, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x12, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69,
	0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x35, 0x0a, 0x05, 0x70, 0x68,
	0x6f, 0x74, 0x6f, 0x18, 0x13, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x74,
	0x6f, 0x12, 0x33, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e,
	0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x06,
	0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x56, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x50, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0d, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0xb4,
	0x02, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x49, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x33, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x76,
	0x65, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68,
	0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x3a, 0x70, 0xc0, 0x9f, 0xe3, 0xb6, 0x05, 0x01, 0x8a, 0xf9, 0x83, 0xb2, 0x05,
	0x32, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x68, 0x6c, 0x37, 0x2e, 0x6f, 0x72, 0x67, 0x2f,
	0x66, 0x68, 0x69, 0x72, 0x2f, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x74, 0x2f, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x2d, 0x67, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x9a, 0xb5, 0x8e, 0x93, 0x06, 0x2c, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f,
	0x68, 0x6c, 0x37, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x53, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x1a, 0xcf, 0x02, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69,
	0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x4d, 0x0a, 0x12, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x5f, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x11,
	0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x48, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69,
	0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x61, 0x62,
	0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x63, 0x65, 0x70, 0x74, 0x42, 0x06, 0xf0, 0xd0, 0x87, 0xeb, 0x04,
	0x01, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x70,
	0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x52, 0x09, 0x70, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x3a, 0x41, 0xc0, 0x9f, 0xe3, 0xb6, 0x05, 0x03, 0xb2,
	0xfe, 0xe4, 0x97, 0x06, 0x35, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x68, 0x6c, 0x37, 0x2e,
	0x6f, 0x72, 0x67, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75,
	0x72, 0x65, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x65, 0x64, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x4a, 0x04, 0x08, 0x07, 0x10, 0x08,
	0x42, 0x7e, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66,
	0x68, 0x69, 0x72, 0x2e, 0x72, 0x34, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x50, 0x01, 0x5a, 0x5b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x72, 0x34, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x5f, 0x67, 0x6f, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x98, 0xc6, 0xb0, 0xb5, 0x07, 0x04,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescOnce sync.Once
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescData = file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDesc
)

func file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescGZIP() []byte {
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescOnce.Do(func() {
		file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescData)
	})
	return file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDescData
}

var file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_google_fhir_proto_r4_core_resources_related_person_proto_goTypes = []interface{}{
	(*RelatedPerson)(nil),                              // 0: google.fhir.r4.core.RelatedPerson
	(*RelatedPerson_GenderCode)(nil),                   // 1: google.fhir.r4.core.RelatedPerson.GenderCode
	(*RelatedPerson_Communication)(nil),                // 2: google.fhir.r4.core.RelatedPerson.Communication
	(*datatypes_go_proto.Id)(nil),                      // 3: google.fhir.r4.core.Id
	(*datatypes_go_proto.Meta)(nil),                    // 4: google.fhir.r4.core.Meta
	(*datatypes_go_proto.Uri)(nil),                     // 5: google.fhir.r4.core.Uri
	(*datatypes_go_proto.Code)(nil),                    // 6: google.fhir.r4.core.Code
	(*datatypes_go_proto.Narrative)(nil),               // 7: google.fhir.r4.core.Narrative
	(*any.Any)(nil),                                    // 8: google.protobuf.Any
	(*datatypes_go_proto.Extension)(nil),               // 9: google.fhir.r4.core.Extension
	(*datatypes_go_proto.Identifier)(nil),              // 10: google.fhir.r4.core.Identifier
	(*datatypes_go_proto.Boolean)(nil),                 // 11: google.fhir.r4.core.Boolean
	(*datatypes_go_proto.Reference)(nil),               // 12: google.fhir.r4.core.Reference
	(*datatypes_go_proto.CodeableConcept)(nil),         // 13: google.fhir.r4.core.CodeableConcept
	(*datatypes_go_proto.HumanName)(nil),               // 14: google.fhir.r4.core.HumanName
	(*datatypes_go_proto.ContactPoint)(nil),            // 15: google.fhir.r4.core.ContactPoint
	(*datatypes_go_proto.Date)(nil),                    // 16: google.fhir.r4.core.Date
	(*datatypes_go_proto.Address)(nil),                 // 17: google.fhir.r4.core.Address
	(*datatypes_go_proto.Attachment)(nil),              // 18: google.fhir.r4.core.Attachment
	(*datatypes_go_proto.Period)(nil),                  // 19: google.fhir.r4.core.Period
	(codes_go_proto.AdministrativeGenderCode_Value)(0), // 20: google.fhir.r4.core.AdministrativeGenderCode.Value
	(*datatypes_go_proto.String)(nil),                  // 21: google.fhir.r4.core.String
}
var file_proto_google_fhir_proto_r4_core_resources_related_person_proto_depIdxs = []int32{
	3,  // 0: google.fhir.r4.core.RelatedPerson.id:type_name -> google.fhir.r4.core.Id
	4,  // 1: google.fhir.r4.core.RelatedPerson.meta:type_name -> google.fhir.r4.core.Meta
	5,  // 2: google.fhir.r4.core.RelatedPerson.implicit_rules:type_name -> google.fhir.r4.core.Uri
	6,  // 3: google.fhir.r4.core.RelatedPerson.language:type_name -> google.fhir.r4.core.Code
	7,  // 4: google.fhir.r4.core.RelatedPerson.text:type_name -> google.fhir.r4.core.Narrative
	8,  // 5: google.fhir.r4.core.RelatedPerson.contained:type_name -> google.protobuf.Any
	9,  // 6: google.fhir.r4.core.RelatedPerson.extension:type_name -> google.fhir.r4.core.Extension
	9,  // 7: google.fhir.r4.core.RelatedPerson.modifier_extension:type_name -> google.fhir.r4.core.Extension
	10, // 8: google.fhir.r4.core.RelatedPerson.identifier:type_name -> google.fhir.r4.core.Identifier
	11, // 9: google.fhir.r4.core.RelatedPerson.active:type_name -> google.fhir.r4.core.Boolean
	12, // 10: google.fhir.r4.core.RelatedPerson.patient:type_name -> google.fhir.r4.core.Reference
	13, // 11: google.fhir.r4.core.RelatedPerson.relationship:type_name -> google.fhir.r4.core.CodeableConcept
	14, // 12: google.fhir.r4.core.RelatedPerson.name:type_name -> google.fhir.r4.core.HumanName
	15, // 13: google.fhir.r4.core.RelatedPerson.telecom:type_name -> google.fhir.r4.core.ContactPoint
	1,  // 14: google.fhir.r4.core.RelatedPerson.gender:type_name -> google.fhir.r4.core.RelatedPerson.GenderCode
	16, // 15: google.fhir.r4.core.RelatedPerson.birth_date:type_name -> google.fhir.r4.core.Date
	17, // 16: google.fhir.r4.core.RelatedPerson.address:type_name -> google.fhir.r4.core.Address
	18, // 17: google.fhir.r4.core.RelatedPerson.photo:type_name -> google.fhir.r4.core.Attachment
	19, // 18: google.fhir.r4.core.RelatedPerson.period:type_name -> google.fhir.r4.core.Period
	2,  // 19: google.fhir.r4.core.RelatedPerson.communication:type_name -> google.fhir.r4.core.RelatedPerson.Communication
	20, // 20: google.fhir.r4.core.RelatedPerson.GenderCode.value:type_name -> google.fhir.r4.core.AdministrativeGenderCode.Value
	21, // 21: google.fhir.r4.core.RelatedPerson.GenderCode.id:type_name -> google.fhir.r4.core.String
	9,  // 22: google.fhir.r4.core.RelatedPerson.GenderCode.extension:type_name -> google.fhir.r4.core.Extension
	21, // 23: google.fhir.r4.core.RelatedPerson.Communication.id:type_name -> google.fhir.r4.core.String
	9,  // 24: google.fhir.r4.core.RelatedPerson.Communication.extension:type_name -> google.fhir.r4.core.Extension
	9,  // 25: google.fhir.r4.core.RelatedPerson.Communication.modifier_extension:type_name -> google.fhir.r4.core.Extension
	13, // 26: google.fhir.r4.core.RelatedPerson.Communication.language:type_name -> google.fhir.r4.core.CodeableConcept
	11, // 27: google.fhir.r4.core.RelatedPerson.Communication.preferred:type_name -> google.fhir.r4.core.Boolean
	28, // [28:28] is the sub-list for method output_type
	28, // [28:28] is the sub-list for method input_type
	28, // [28:28] is the sub-list for extension type_name
	28, // [28:28] is the sub-list for extension extendee
	0,  // [0:28] is the sub-list for field type_name
}

func init() { file_proto_google_fhir_proto_r4_core_resources_related_person_proto_init() }
func file_proto_google_fhir_proto_r4_core_resources_related_person_proto_init() {
	if File_proto_google_fhir_proto_r4_core_resources_related_person_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelatedPerson); i {
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
		file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelatedPerson_GenderCode); i {
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
		file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelatedPerson_Communication); i {
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
			RawDescriptor: file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_google_fhir_proto_r4_core_resources_related_person_proto_goTypes,
		DependencyIndexes: file_proto_google_fhir_proto_r4_core_resources_related_person_proto_depIdxs,
		MessageInfos:      file_proto_google_fhir_proto_r4_core_resources_related_person_proto_msgTypes,
	}.Build()
	File_proto_google_fhir_proto_r4_core_resources_related_person_proto = out.File
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_rawDesc = nil
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_goTypes = nil
	file_proto_google_fhir_proto_r4_core_resources_related_person_proto_depIdxs = nil
}
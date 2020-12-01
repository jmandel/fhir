//    Copyright 2020 Google Inc.
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
// source: proto/google/fhir/proto/r5/core/resources/manufactured_item_definition.proto

package manufactured_item_definition_go_proto

import (
	any "github.com/golang/protobuf/ptypes/any"
	_ "github.com/google/fhir/go/proto/google/fhir/proto/annotations_go_proto"
	datatypes_go_proto "github.com/google/fhir/go/proto/google/fhir/proto/r5/core/datatypes_go_proto"
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

// Auto-generated from StructureDefinition for ManufacturedItemDefinition, last
// updated 2019-12-31T21:03:40.621+11:00. The definition and characteristics of
// a medicinal manufactured item, such as a tablet or capsule, as contained in a
// packaged medicinal product. See
// http://hl7.org/fhir/StructureDefinition/ManufacturedItemDefinition
type ManufacturedItemDefinition struct {
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
	// Unique identifier
	Identifier []*datatypes_go_proto.Identifier `protobuf:"bytes,10,rep,name=identifier,proto3" json:"identifier,omitempty"`
	// Dose form as manufactured and before any transformation into the
	// pharmaceutical product
	ManufacturedDoseForm *datatypes_go_proto.CodeableConcept `protobuf:"bytes,11,opt,name=manufactured_dose_form,json=manufacturedDoseForm,proto3" json:"manufactured_dose_form,omitempty"`
	// The “real world” units in which the quantity of the manufactured item is
	// described
	UnitOfPresentation *datatypes_go_proto.CodeableConcept `protobuf:"bytes,12,opt,name=unit_of_presentation,json=unitOfPresentation,proto3" json:"unit_of_presentation,omitempty"`
	// Manufacturer of the item (Note that this should be named "manufacturer" but
	// it currently causes technical issues)
	Manufacturer []*datatypes_go_proto.Reference `protobuf:"bytes,13,rep,name=manufacturer,proto3" json:"manufacturer,omitempty"`
	// The ingredients that make up this manufactured item
	Ingredient     []*datatypes_go_proto.Reference              `protobuf:"bytes,14,rep,name=ingredient,proto3" json:"ingredient,omitempty"`
	Characteristic []*ManufacturedItemDefinition_Characteristic `protobuf:"bytes,15,rep,name=characteristic,proto3" json:"characteristic,omitempty"`
}

func (x *ManufacturedItemDefinition) Reset() {
	*x = ManufacturedItemDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManufacturedItemDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManufacturedItemDefinition) ProtoMessage() {}

func (x *ManufacturedItemDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManufacturedItemDefinition.ProtoReflect.Descriptor instead.
func (*ManufacturedItemDefinition) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescGZIP(), []int{0}
}

func (x *ManufacturedItemDefinition) GetId() *datatypes_go_proto.Id {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetMeta() *datatypes_go_proto.Meta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetImplicitRules() *datatypes_go_proto.Uri {
	if x != nil {
		return x.ImplicitRules
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetLanguage() *datatypes_go_proto.Code {
	if x != nil {
		return x.Language
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetText() *datatypes_go_proto.Narrative {
	if x != nil {
		return x.Text
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetContained() []*any.Any {
	if x != nil {
		return x.Contained
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.Extension
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetModifierExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.ModifierExtension
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetIdentifier() []*datatypes_go_proto.Identifier {
	if x != nil {
		return x.Identifier
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetManufacturedDoseForm() *datatypes_go_proto.CodeableConcept {
	if x != nil {
		return x.ManufacturedDoseForm
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetUnitOfPresentation() *datatypes_go_proto.CodeableConcept {
	if x != nil {
		return x.UnitOfPresentation
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetManufacturer() []*datatypes_go_proto.Reference {
	if x != nil {
		return x.Manufacturer
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetIngredient() []*datatypes_go_proto.Reference {
	if x != nil {
		return x.Ingredient
	}
	return nil
}

func (x *ManufacturedItemDefinition) GetCharacteristic() []*ManufacturedItemDefinition_Characteristic {
	if x != nil {
		return x.Characteristic
	}
	return nil
}

// General characteristics of this item
type ManufacturedItemDefinition_Characteristic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique id for inter-element referencing
	Id *datatypes_go_proto.Id `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Additional content defined by implementations
	Extension []*datatypes_go_proto.Extension `protobuf:"bytes,2,rep,name=extension,proto3" json:"extension,omitempty"`
	// Extensions that cannot be ignored even if unrecognized
	ModifierExtension []*datatypes_go_proto.Extension `protobuf:"bytes,3,rep,name=modifier_extension,json=modifierExtension,proto3" json:"modifier_extension,omitempty"`
	// A code expressing the type of characteristic
	Code  *datatypes_go_proto.CodeableConcept               `protobuf:"bytes,4,opt,name=code,proto3" json:"code,omitempty"`
	Value *ManufacturedItemDefinition_Characteristic_ValueX `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ManufacturedItemDefinition_Characteristic) Reset() {
	*x = ManufacturedItemDefinition_Characteristic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManufacturedItemDefinition_Characteristic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManufacturedItemDefinition_Characteristic) ProtoMessage() {}

func (x *ManufacturedItemDefinition_Characteristic) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManufacturedItemDefinition_Characteristic.ProtoReflect.Descriptor instead.
func (*ManufacturedItemDefinition_Characteristic) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ManufacturedItemDefinition_Characteristic) GetId() *datatypes_go_proto.Id {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic) GetExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.Extension
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic) GetModifierExtension() []*datatypes_go_proto.Extension {
	if x != nil {
		return x.ModifierExtension
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic) GetCode() *datatypes_go_proto.CodeableConcept {
	if x != nil {
		return x.Code
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic) GetValue() *ManufacturedItemDefinition_Characteristic_ValueX {
	if x != nil {
		return x.Value
	}
	return nil
}

// A value for the characteristic
type ManufacturedItemDefinition_Characteristic_ValueX struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Choice:
	//	*ManufacturedItemDefinition_Characteristic_ValueX_Coding
	//	*ManufacturedItemDefinition_Characteristic_ValueX_Quantity
	//	*ManufacturedItemDefinition_Characteristic_ValueX_StringValue
	//	*ManufacturedItemDefinition_Characteristic_ValueX_Date
	//	*ManufacturedItemDefinition_Characteristic_ValueX_Boolean
	//	*ManufacturedItemDefinition_Characteristic_ValueX_Attachment
	Choice isManufacturedItemDefinition_Characteristic_ValueX_Choice `protobuf_oneof:"choice"`
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) Reset() {
	*x = ManufacturedItemDefinition_Characteristic_ValueX{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManufacturedItemDefinition_Characteristic_ValueX) ProtoMessage() {}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) ProtoReflect() protoreflect.Message {
	mi := &file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManufacturedItemDefinition_Characteristic_ValueX.ProtoReflect.Descriptor instead.
func (*ManufacturedItemDefinition_Characteristic_ValueX) Descriptor() ([]byte, []int) {
	return file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescGZIP(), []int{0, 0, 0}
}

func (m *ManufacturedItemDefinition_Characteristic_ValueX) GetChoice() isManufacturedItemDefinition_Characteristic_ValueX_Choice {
	if m != nil {
		return m.Choice
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetCoding() *datatypes_go_proto.Coding {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_Coding); ok {
		return x.Coding
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetQuantity() *datatypes_go_proto.Quantity {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_Quantity); ok {
		return x.Quantity
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetStringValue() *datatypes_go_proto.String {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_StringValue); ok {
		return x.StringValue
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetDate() *datatypes_go_proto.Date {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_Date); ok {
		return x.Date
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetBoolean() *datatypes_go_proto.Boolean {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_Boolean); ok {
		return x.Boolean
	}
	return nil
}

func (x *ManufacturedItemDefinition_Characteristic_ValueX) GetAttachment() *datatypes_go_proto.Attachment {
	if x, ok := x.GetChoice().(*ManufacturedItemDefinition_Characteristic_ValueX_Attachment); ok {
		return x.Attachment
	}
	return nil
}

type isManufacturedItemDefinition_Characteristic_ValueX_Choice interface {
	isManufacturedItemDefinition_Characteristic_ValueX_Choice()
}

type ManufacturedItemDefinition_Characteristic_ValueX_Coding struct {
	Coding *datatypes_go_proto.Coding `protobuf:"bytes,1,opt,name=coding,proto3,oneof"`
}

type ManufacturedItemDefinition_Characteristic_ValueX_Quantity struct {
	Quantity *datatypes_go_proto.Quantity `protobuf:"bytes,2,opt,name=quantity,proto3,oneof"`
}

type ManufacturedItemDefinition_Characteristic_ValueX_StringValue struct {
	StringValue *datatypes_go_proto.String `protobuf:"bytes,3,opt,name=string_value,json=string,proto3,oneof"`
}

type ManufacturedItemDefinition_Characteristic_ValueX_Date struct {
	Date *datatypes_go_proto.Date `protobuf:"bytes,4,opt,name=date,proto3,oneof"`
}

type ManufacturedItemDefinition_Characteristic_ValueX_Boolean struct {
	Boolean *datatypes_go_proto.Boolean `protobuf:"bytes,5,opt,name=boolean,proto3,oneof"`
}

type ManufacturedItemDefinition_Characteristic_ValueX_Attachment struct {
	Attachment *datatypes_go_proto.Attachment `protobuf:"bytes,6,opt,name=attachment,proto3,oneof"`
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_Coding) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_Quantity) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_StringValue) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_Date) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_Boolean) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

func (*ManufacturedItemDefinition_Characteristic_ValueX_Attachment) isManufacturedItemDefinition_Characteristic_ValueX_Choice() {
}

var File_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto protoreflect.FileDescriptor

var file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDesc = []byte{
	0x0a, 0x4c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66,
	0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x35, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x75,
	0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x64, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x29,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69,
	0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x72, 0x35, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaa, 0x0e, 0x0a, 0x1a, 0x4d,
	0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x44,
	0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66,
	0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x64, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2d, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72,
	0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74,
	0x61, 0x12, 0x3f, 0x0a, 0x0e, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x5f, 0x72, 0x75,
	0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x55, 0x72, 0x69, 0x52, 0x0d, 0x69, 0x6d, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x52, 0x75, 0x6c,
	0x65, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68,
	0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x52,
	0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4e, 0x61,
	0x72, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x32, 0x0a,
	0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x64, 0x12, 0x3c, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x08,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68,
	0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x4d, 0x0a, 0x12, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x5f, 0x65, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x6d, 0x6f, 0x64,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3f,
	0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72,
	0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12,
	0x62, 0x0a, 0x16, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x5f,
	0x64, 0x6f, 0x73, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x6f,
	0x6e, 0x63, 0x65, 0x70, 0x74, 0x42, 0x06, 0xf0, 0xd0, 0x87, 0xeb, 0x04, 0x01, 0x52, 0x14, 0x6d,
	0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x44, 0x6f, 0x73, 0x65, 0x46,
	0x6f, 0x72, 0x6d, 0x12, 0x56, 0x0a, 0x14, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x6f, 0x66, 0x5f, 0x70,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e,
	0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x61, 0x62, 0x6c, 0x65,
	0x43, 0x6f, 0x6e, 0x63, 0x65, 0x70, 0x74, 0x52, 0x12, 0x75, 0x6e, 0x69, 0x74, 0x4f, 0x66, 0x50,
	0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x56, 0x0a, 0x0c, 0x6d,
	0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x72, 0x18, 0x0d, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e,
	0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x42, 0x12, 0xf2, 0xff, 0xfc, 0xc2, 0x06, 0x0c, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75,
	0x72, 0x65, 0x72, 0x12, 0x50, 0x0a, 0x0a, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x42, 0x10, 0xf2, 0xff, 0xfc, 0xc2, 0x06, 0x0a, 0x49,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x69, 0x6e, 0x67, 0x72, 0x65,
	0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x66, 0x0a, 0x0e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3e, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64,
	0x49, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x52, 0x0e, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x1a, 0xe1, 0x05,
	0x0a, 0x0e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63,
	0x12, 0x27, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x09, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4d, 0x0a, 0x12, 0x6d, 0x6f, 0x64, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69,
	0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x52, 0x11, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x45, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x40, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68,
	0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x65, 0x61,
	0x62, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x63, 0x65, 0x70, 0x74, 0x42, 0x06, 0xf0, 0xd0, 0x87, 0xeb,
	0x04, 0x01, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x5b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x45, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x61,
	0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x44, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x69, 0x73, 0x74, 0x69, 0x63, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x58, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0xf9, 0x02, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x58,
	0x12, 0x35, 0x0a, 0x06, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72,
	0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52,
	0x06, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x3b, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x48, 0x00, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x12, 0x3b, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x72, 0x69, 0x6e,
	0x67, 0x12, 0x2f, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35,
	0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69,
	0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x65, 0x61,
	0x6e, 0x48, 0x00, 0x52, 0x07, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e, 0x12, 0x41, 0x0a, 0x0a,
	0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72,
	0x35, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e,
	0x74, 0x48, 0x00, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x61, 0x63, 0x68, 0x6d, 0x65, 0x6e, 0x74, 0x3a,
	0x06, 0xa0, 0x83, 0x83, 0xe8, 0x06, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x63, 0x68, 0x6f, 0x69, 0x63,
	0x65, 0x3a, 0x4e, 0xc0, 0x9f, 0xe3, 0xb6, 0x05, 0x03, 0xb2, 0xfe, 0xe4, 0x97, 0x06, 0x42, 0x68,
	0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x68, 0x6c, 0x37, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x66, 0x68,
	0x69, 0x72, 0x2f, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x44, 0x65, 0x66, 0x69,
	0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x4d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75,
	0x72, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x4a, 0x04, 0x08, 0x07, 0x10, 0x08, 0x42, 0x8c, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x66, 0x68, 0x69, 0x72, 0x2e, 0x72, 0x35, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x50, 0x01, 0x5a, 0x69, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68, 0x69, 0x72, 0x2f, 0x67, 0x6f,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x66, 0x68,
	0x69, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x35, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x75, 0x66,
	0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x64, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x67, 0x6f, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x98, 0xc6, 0xb0, 0xb5, 0x07, 0x05, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescOnce sync.Once
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescData = file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDesc
)

func file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescGZIP() []byte {
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescOnce.Do(func() {
		file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescData)
	})
	return file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDescData
}

var file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_goTypes = []interface{}{
	(*ManufacturedItemDefinition)(nil),                       // 0: google.fhir.r5.core.ManufacturedItemDefinition
	(*ManufacturedItemDefinition_Characteristic)(nil),        // 1: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic
	(*ManufacturedItemDefinition_Characteristic_ValueX)(nil), // 2: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX
	(*datatypes_go_proto.Id)(nil),                            // 3: google.fhir.r5.core.Id
	(*datatypes_go_proto.Meta)(nil),                          // 4: google.fhir.r5.core.Meta
	(*datatypes_go_proto.Uri)(nil),                           // 5: google.fhir.r5.core.Uri
	(*datatypes_go_proto.Code)(nil),                          // 6: google.fhir.r5.core.Code
	(*datatypes_go_proto.Narrative)(nil),                     // 7: google.fhir.r5.core.Narrative
	(*any.Any)(nil),                                          // 8: google.protobuf.Any
	(*datatypes_go_proto.Extension)(nil),                     // 9: google.fhir.r5.core.Extension
	(*datatypes_go_proto.Identifier)(nil),                    // 10: google.fhir.r5.core.Identifier
	(*datatypes_go_proto.CodeableConcept)(nil),               // 11: google.fhir.r5.core.CodeableConcept
	(*datatypes_go_proto.Reference)(nil),                     // 12: google.fhir.r5.core.Reference
	(*datatypes_go_proto.Coding)(nil),                        // 13: google.fhir.r5.core.Coding
	(*datatypes_go_proto.Quantity)(nil),                      // 14: google.fhir.r5.core.Quantity
	(*datatypes_go_proto.String)(nil),                        // 15: google.fhir.r5.core.String
	(*datatypes_go_proto.Date)(nil),                          // 16: google.fhir.r5.core.Date
	(*datatypes_go_proto.Boolean)(nil),                       // 17: google.fhir.r5.core.Boolean
	(*datatypes_go_proto.Attachment)(nil),                    // 18: google.fhir.r5.core.Attachment
}
var file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_depIdxs = []int32{
	3,  // 0: google.fhir.r5.core.ManufacturedItemDefinition.id:type_name -> google.fhir.r5.core.Id
	4,  // 1: google.fhir.r5.core.ManufacturedItemDefinition.meta:type_name -> google.fhir.r5.core.Meta
	5,  // 2: google.fhir.r5.core.ManufacturedItemDefinition.implicit_rules:type_name -> google.fhir.r5.core.Uri
	6,  // 3: google.fhir.r5.core.ManufacturedItemDefinition.language:type_name -> google.fhir.r5.core.Code
	7,  // 4: google.fhir.r5.core.ManufacturedItemDefinition.text:type_name -> google.fhir.r5.core.Narrative
	8,  // 5: google.fhir.r5.core.ManufacturedItemDefinition.contained:type_name -> google.protobuf.Any
	9,  // 6: google.fhir.r5.core.ManufacturedItemDefinition.extension:type_name -> google.fhir.r5.core.Extension
	9,  // 7: google.fhir.r5.core.ManufacturedItemDefinition.modifier_extension:type_name -> google.fhir.r5.core.Extension
	10, // 8: google.fhir.r5.core.ManufacturedItemDefinition.identifier:type_name -> google.fhir.r5.core.Identifier
	11, // 9: google.fhir.r5.core.ManufacturedItemDefinition.manufactured_dose_form:type_name -> google.fhir.r5.core.CodeableConcept
	11, // 10: google.fhir.r5.core.ManufacturedItemDefinition.unit_of_presentation:type_name -> google.fhir.r5.core.CodeableConcept
	12, // 11: google.fhir.r5.core.ManufacturedItemDefinition.manufacturer:type_name -> google.fhir.r5.core.Reference
	12, // 12: google.fhir.r5.core.ManufacturedItemDefinition.ingredient:type_name -> google.fhir.r5.core.Reference
	1,  // 13: google.fhir.r5.core.ManufacturedItemDefinition.characteristic:type_name -> google.fhir.r5.core.ManufacturedItemDefinition.Characteristic
	3,  // 14: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.id:type_name -> google.fhir.r5.core.Id
	9,  // 15: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.extension:type_name -> google.fhir.r5.core.Extension
	9,  // 16: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.modifier_extension:type_name -> google.fhir.r5.core.Extension
	11, // 17: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.code:type_name -> google.fhir.r5.core.CodeableConcept
	2,  // 18: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.value:type_name -> google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX
	13, // 19: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.coding:type_name -> google.fhir.r5.core.Coding
	14, // 20: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.quantity:type_name -> google.fhir.r5.core.Quantity
	15, // 21: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.string_value:type_name -> google.fhir.r5.core.String
	16, // 22: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.date:type_name -> google.fhir.r5.core.Date
	17, // 23: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.boolean:type_name -> google.fhir.r5.core.Boolean
	18, // 24: google.fhir.r5.core.ManufacturedItemDefinition.Characteristic.ValueX.attachment:type_name -> google.fhir.r5.core.Attachment
	25, // [25:25] is the sub-list for method output_type
	25, // [25:25] is the sub-list for method input_type
	25, // [25:25] is the sub-list for extension type_name
	25, // [25:25] is the sub-list for extension extendee
	0,  // [0:25] is the sub-list for field type_name
}

func init() { file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_init() }
func file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_init() {
	if File_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManufacturedItemDefinition); i {
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
		file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManufacturedItemDefinition_Characteristic); i {
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
		file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManufacturedItemDefinition_Characteristic_ValueX); i {
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
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*ManufacturedItemDefinition_Characteristic_ValueX_Coding)(nil),
		(*ManufacturedItemDefinition_Characteristic_ValueX_Quantity)(nil),
		(*ManufacturedItemDefinition_Characteristic_ValueX_StringValue)(nil),
		(*ManufacturedItemDefinition_Characteristic_ValueX_Date)(nil),
		(*ManufacturedItemDefinition_Characteristic_ValueX_Boolean)(nil),
		(*ManufacturedItemDefinition_Characteristic_ValueX_Attachment)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_goTypes,
		DependencyIndexes: file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_depIdxs,
		MessageInfos:      file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_msgTypes,
	}.Build()
	File_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto = out.File
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_rawDesc = nil
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_goTypes = nil
	file_proto_google_fhir_proto_r5_core_resources_manufactured_item_definition_proto_depIdxs = nil
}

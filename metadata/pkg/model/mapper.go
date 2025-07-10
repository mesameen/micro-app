package model

import "github.com/mesameen/micro-app/src/api/gen"

// MetadataToProto converts a Metada struct into a
// generated proto counterpart
func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

// MetadaFromProto coverts a generated proto counterpart
// into a Metadata struct
func MetadaFromProto(m *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}

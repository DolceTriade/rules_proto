package protoc

import (
	"fmt"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// OtherProtoLibrary implements the ProtoLibrary interface from an existing
// ProtoLibrary rule generated by an other extension (typically the
// bazel-gazelle/proto extension).
type OtherProtoLibrary struct {
	source *rule.File
	rule   *rule.Rule
	files  []*File
}

// NewOtherProtoLibrary constructs a new ProtoLibrary implementation.
func NewOtherProtoLibrary(source *rule.File, rule *rule.Rule, files ...*File) *OtherProtoLibrary {
	return &OtherProtoLibrary{source, rule, files}
}

// Name implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Name() string {
	return s.rule.Name()
}

// BaseName implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) BaseName() string {
	name := s.rule.Name()
	if !strings.HasSuffix(name, "_proto") {
		panic(fmt.Sprintf("Unexpected proto_library name %q (it should always end in '_proto')", name))
	}
	return name[0 : len(name)-len("_proto")]
}

// Rule implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Rule() *rule.Rule {
	return s.rule
}

// Files implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Files() []*File {
	return s.files
}

// Deps implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Deps() []string {
	return s.rule.AttrStrings("deps")
}

// Imports implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Imports() []string {
	// Not supposed to be using this private attr, but...
	importRaw := s.rule.PrivateAttr(config.GazelleImportsKey)
	if v, ok := importRaw.([]string); ok {
		return v
	}
	return nil
}

// Srcs returns the srcs attribute
func (s *OtherProtoLibrary) Srcs() []string {
	return s.rule.AttrStrings("srcs")
	// srcs := make([]string, len(s.files))
	// for _, f := range s.files {
	// 	srcs = append(srcs, path.Join(f.Dir, f.Basename))
	// }
	// return srcs
}

// StripImportPrefix implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) StripImportPrefix() string {
	prefix := s.rule.AttrString("strip_import_prefix")
	if prefix != "" {
		return prefix
	}
	return GetKeptFileRuleAttrString(s.source, s.rule, "strip_import_prefix")
}

package rules_python

import (
	"fmt"
	"path"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func TestImports(t *testing.T) {
	kind := "mykind"
	suffix := "_suffix"
	pkg := "mypkg"
	protoName := "test"
	cases := []struct {
		Name        string
		Outputs     []string
		WantImports []resolve.ImportSpec
	}{{
		Name: "Empty",
		// If for some reason, no python files were output...
		Outputs: []string{},
		// Always include the output from the proto_library
		WantImports: []resolve.ImportSpec{{Lang: kind, Imp: fmt.Sprintf("%s/%s", pkg, protoName)}},
	}, {
		Name:    "One output",
		Outputs: []string{path.Join(pkg, "test_pb2.py")},
		WantImports: []resolve.ImportSpec{
			{Lang: kind, Imp: fmt.Sprintf("%s/%s", pkg, protoName)},
			{Lang: "py", Imp: fmt.Sprintf("%s.%s_pb2", pkg, protoName)},
		},
	}}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			py := &PyLibrary{
				KindName:       kind,
				RuleNameSuffix: suffix,
				Outputs:        c.Outputs,
				Resolver:       protoc.ResolveDepsAttr("deps", true),
			}
			protoLib := protoc.NewOtherProtoLibrary(&rule.File{}, rule.NewRule("proto_library", protoName+"_proto"), protoc.NewFile(pkg, protoName))
			r := rule.NewRule(kind, "test"+suffix)
			r.SetPrivateAttr(protoc.ProtoLibraryKey, protoLib)
			imps := py.Imports(nil, r, &rule.File{Pkg: pkg})
			if diff := cmp.Diff(imps, c.WantImports); diff != "" {
				t.Fatalf("import mismatch: (-got, +want): %s", diff)
			}
			// TODO: How to 
		})
	}
}

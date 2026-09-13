package gogen_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/omniaura/agentflow/pkg/ast"
	"github.com/omniaura/agentflow/pkg/gen/gogen"
)

func TestLoopsCompileAndRender(t *testing.T) {
	source := ".title system\n<*extensions []Extension as extension><!extension.name>:<!extension.detail.text>;<*extension.scores []int as score><!score>,</extension.scores></extensions>\n.title user\n<*names []string as name><!name>;</names><*counts []int as count><!count>;</counts>"
	f, err := ast.NewFile("prompts.af", []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := gogen.GenFile(&b, f, "example"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	fixture := "package example\nimport \"testing\"\ntype Detail struct {Text string}; type Extension struct {Name string; Detail Detail; Scores []int}\nfunc TestRender(t *testing.T) {\ns:=System{Extensions:[]Extension{{Name:\"writing\",Detail:Detail{Text:\"concise\"},Scores:[]int{1,2}}}};if got:=s.String();got!=\"writing:concise;1,2,\" {t.Fatal(got)}\nu:=User{Names:[]string{\"A\",\"B\"},Counts:[]int{0,-2}};if got:=u.String();got!=\"A;B;0;-2;\" {t.Fatal(got)};u=User{};if got:=u.String();got!=\"\" {t.Fatal(got)}\n}\n"
	files := map[string]string{"go.mod": "module example\ngo 1.25.0\n", "prompts.go": b.String(), "prompts_test.go": fixture}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated code failed: %v\n%s\n%s", err, out, b.String())
	}
}

func TestInvalidLoops(t *testing.T) {
	for _, source := range []string{
		"<*items int as item></items>",
		"<*items []string as item></items><!item>",
		"<*items []string as item><!value string;panic()></items>",
		"<*items []string as item></wrong>",
		"<*items []string as item>",
		"<*items []string as item><else></items>",
		"<*items []string as item><*others []int as item></others></items>",
		"<*items []string as item;panic></items>",
	} {
		t.Run(source, func(t *testing.T) {
			f, err := ast.NewFile("invalid.af", []byte(source))
			if err != nil {
				return
			}
			var b bytes.Buffer
			err = gogen.GenFile(&b, f, "example")
			if err == nil {
				t.Fatalf("expected error; got %s", b.String())
			}
			if !strings.Contains(err.Error(), "invalid") {
				t.Fatalf("missing source context: %v", err)
			}
		})
	}
}

func TestArbitraryValuesAndNestedSlices(t *testing.T) {
	source := ".title values\n<*values []any as value><!value>;</values><*matrix [][]int as row><*row []int as number><!number>,</row>;</matrix>"
	f, err := ast.NewFile("values.af", []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err := gogen.GenFile(&b, f, "example"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	fixture := "package example\nimport \"testing\"\nfunc TestValues(t *testing.T){v:=Values{Values:[]any{\"text\",42,struct{Name string}{\"struct\"},nil},Matrix:[][]int{{1,2},nil,{3}}};if got:=v.String();got!=\"text;42;{struct};<nil>;1,2,;;3,;\"{t.Fatal(got)}}"
	for name, content := range map[string]string{"go.mod": "module example\ngo 1.25.0\n", "values.go": b.String(), "values_test.go": fixture} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s\n%s", err, out, b.String())
	}
}

func TestLoopsWithImplicitObjectsAndLexicalConditions(t *testing.T) {
	source := ".title mixed\n<!account.name>:<*items []Item as item><?item.enabled bool><!item.name><*item.tags []string as tag><!item.name>=<!tag>;</item.tags><else>off</item.enabled></items><!account.count int>"
	f, err := ast.NewFile("mixed.af", []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	var b bytes.Buffer
	if err = gogen.GenFile(&b, f, "example"); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	fixture := "package example\nimport \"testing\"\ntype Item struct{Name string; Enabled bool; Tags []string}\nfunc TestMixed(t *testing.T){v:=Mixed{Items:[]Item{{Name:\"A\",Enabled:true,Tags:[]string{\"x\",\"y\"}},{Name:\"B\"}}};v.Account.Name=\"owner\";v.Account.Count=2;if got:=v.String();got!=\"owner:AA=x;A=y;off2\"{t.Fatal(got)}}"
	for name, content := range map[string]string{"go.mod": "module example\ngo 1.25.0\n", "mixed.go": b.String(), "mixed_test.go": fixture} {
		if err = os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s\n%s", err, out, b.String())
	}
}

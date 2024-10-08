/*
Copyright © 2024 Omni Aura peyton@omniaura.co

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package js_test

import (
	"strings"
	"testing"

	"github.com/omniaura/agentflow/cfg"
	"github.com/omniaura/agentflow/pkg/assert/require"
	"github.com/omniaura/agentflow/pkg/ast"
	"github.com/omniaura/agentflow/pkg/gen/js"
	"github.com/omniaura/agentflow/pkg/logger"
	"github.com/omniaura/agentflow/tests/testdata"
)

func TestMain(m *testing.M) {
	cfg.TestMode()
	logger.Setup()
	m.Run()
}

type TestCase struct {
	Name     string
	Filename string
	Content  string
	Want     string
}

func TestGenerate(t *testing.T) {
	cases := []TestCase{
		{
			Name:     "no vars no title",
			Filename: "no_vars_no_title.af",
			Content:  testdata.NoVarsNoTitle,
			Want: "/**\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function noVarsNoTitle() {\n" +
				"	return `say hello to the user!`;\n}\n",
		},
		{
			Name:     "single prompt",
			Filename: "hello1.af",
			Content:  testdata.OneVarNoTitle,
			Want: "/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function hello1(username) {\n" +
				"	return `say hello to ${username}`;\n}\n",
		},
		{
			Name:     "single prompt with title",
			Filename: "hello2.af",
			Content:  testdata.OneVarWithTitle,
			Want: "/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function helloUser(username) {\n" +
				"	return `say hello to ${username}`;\n}\n",
		},
		{
			Name:     "two prompts with titles",
			Filename: "hello3.af",
			Content:  testdata.TwoPromptsWithVars,
			Want: "/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function helloUser(username) {\n" +
				"	return `say hello to ${username}`;\n}\n" +
				"\n" +
				"/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function goodbyeUser(username) {\n" +
				"	return `say goodbye to ${username}`;\n}\n",
		},
		{
			Name:     "two prompts with extra newline",
			Filename: "hello4.af",
			Content:  testdata.TwoPromptsExtraNewline,
			Want: "/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function helloUser(username) {\n" +
				"	return `say hello to ${username}\n`;\n}\n" +
				"\n" +
				"/**\n" +
				" * @param {string} username\n" +
				" * @returns {string}\n" +
				" */\n" +
				"export function goodbyeUser(username) {\n" +
				"	return `say goodbye to ${username}\n`;\n}\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			file, err := ast.NewFile(tc.Filename, []byte(tc.Content))
			require.NoError(t, err)
			var buf strings.Builder
			js.GenFile(&buf, file)
			got := buf.String()
			if got != tc.Want {
				var sb strings.Builder
				sb.WriteRune('\n')
				require.WantGotBoldQuotes(&sb, tc.Want, got)
				t.Error(sb.String())
			}
		})
	}
}

// nolint lll
package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	that := assert.New(t)
	// {{.}}	Renders the root element
	that.Equal(M2(`Hello "bingoo"`, nil), M2(TmplRenderText(`Hello "{{.}}"`, "bingoo")))

	type Item struct {
		Name string
		Desc string
	}

	// {{.Title}}	Renders the “Title”-field in a nested element
	that.Equal(M2(`You have a task named "bingoo" with description:"welcome"`, nil), M2(
		TmplRenderText(`You have a task named "{{.Name}}" with description:"{{.Desc}}"`, Item{"bingoo", "welcome"})))

	that.Equal(M2(`You have a task named "bingoo" with description:"welcome"`, nil), M2(
		TmplRenderText(`You have a task named "{{.Name}}" with description:"{{.Desc}}"`, map[string]interface{}{"Name": "bingoo", "Desc": "welcome"})))

	type Todo struct {
		Title string
		Done  bool
	}

	type TodoPageData struct {
		PageTitle string
		Todos     []Todo
	}

	that.Equal(M2(`<h1>My TODO list</h1><ul><li>Task 1</li><li class="done">Task 2</li><li class="done">Task 3</li></ul>`, nil),
		M2(TmplRenderText(`<h1>{{.PageTitle}}</h1><ul>{{range .Todos}}{{if .Done}}<li class="done">{{.Title}}</li>{{else}}<li>{{.Title}}</li>{{end}}{{end}}</ul>`,
			TodoPageData{PageTitle: "My TODO list",
				Todos: []Todo{
					{Title: "Task 1", Done: false},
					{Title: "Task 2", Done: true},
					{Title: "Task 3", Done: true},
				}})))

	that.Equal(M2(`For k=001,v=S,k=002,v=H,k=003,v=C,k=004,v=D,`, nil),
		M2(TmplRenderText(`For {{range $k,$v := .}}k={{printf "%03d" $k}},v={{$v}},{{end}}`, map[int]string{1: "S", 2: "H", 3: "C", 4: "D"})))

	that.Equal(M2(`Repeat Ape ate Apple`, nil), M2(TmplRenderText(
		`Repeat {{define "T1"}}Apple{{end}}{{define "T2"}}Ape{{end}}{{template "T2"}} ate {{template "T1"}}`,
		map[string]string{"a": "S", "b": "H", "c": "C", "d": "D"})))

	// Variables in Templates
	that.Equal(M2(`It is day number 12 of the March`, nil), M2(TmplRenderText(`{{$number := .}}It is day number {{$number}} of the March`, 12)))
}

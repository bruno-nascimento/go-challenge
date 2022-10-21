package printer

import (
	"fmt"
	"strings"
)

var header = `╔╦╗╔═╗╔╦╗╦═╗╦╔═╗╔═╗   ╔═╗╦ ╦╔╦╗╔╦╗╔═╗╦═╗╦╔═╗╔═╗╦═╗
 ║║║║╣  ║ ╠╦╝║║  ╚═╗   ╚═╗║ ║║║║║║║╠═╣╠╦╝║╔═╝║╣ ╠╦╝
 ╩ ╩╚═╝ ╩ ╩╚═╩╚═╝╚═╝   ╚═╝╚═╝╩ ╩╩ ╩╩ ╩╩╚═╩╚═╝╚═╝╩╚═ V0.1.0`

var tableHeader = `
+---------------------------------------+---------------------------------------+
| Level name				| Total value				|
+---------------------------------------+---------------------------------------+`

var tableRow = `
  %s				| %d`

var tableFooter = `
+---------------------------------------+---------------------------------------+`

func Header() {
	fmt.Println("\033[34m", header, "\033[0m")
}

type Builder struct {
	sb strings.Builder
}

func NewTableBuilder() *Builder {
	b := &Builder{}
	b.sb.WriteString(tableHeader)
	return b
}

func (b *Builder) AddRow(level string, value int64) *Builder {
	b.sb.WriteString(fmt.Sprintf(tableRow, level, value))
	return b
}

func (b *Builder) Build() string {
	// b.sb.WriteString(tableFooter)
	return b.sb.String()
}

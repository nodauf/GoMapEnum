package template

import (
	_ "embed"
)

//go:embed datatables.tpl
var templateDatatables string

func GetTemplateDatatables() string {
	return templateDatatables
}

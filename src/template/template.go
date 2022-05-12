package template

// File that contains all the embeded templates

import (
	_ "embed"
)

//go:embed datatables.tpl
var templateDatatables string

// GetTemplateDatatables returns the content of the datatables template
func GetTemplateDatatables() string {
	return templateDatatables
}

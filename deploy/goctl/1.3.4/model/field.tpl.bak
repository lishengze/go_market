{{if eq .type "sql.NullFloat64"}}
{{.name}} decimal.Decimal {{.tag}} {{if .hasComment}} // {{.comment}}{{end}}
{{else}}
{{.name}} {{.type}} {{.tag}} {{if .hasComment}}// {{.comment}}{{end}}
{{end}}
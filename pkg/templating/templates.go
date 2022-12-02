package templating

var templates = map[string]string{
	"int":     "-?[0-9]+",
	"number":  `-?[0-9]+(.[0-9]+)?`,
	"uuid":    `[0-9A-Fa-f-]{8}-[0-9A-Fa-f-]{4}-[0-9A-Fa-f-]{4}-[0-9A-Fa-f-]{4}-[0-9A-Fa-f-]{12}`,
	"integer": `[0-9]+`,
	"email":   `[^ ]+@[^ ]+[.][^ ]+`,
	"word":    `\b[a-zA-Z0-9-]+\b`,
}

func GetTemplates() map[string]string {
	return templates
}

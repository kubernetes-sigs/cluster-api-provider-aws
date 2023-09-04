# Release notes for Cluster API Provider AWS (CAPA) <RELEASE_VERSION>

[Documentation](https://cluster-api-aws.sigs.k8s.io/)

# Changelog since <PREVIOUS_VERSION>

{{with .NotesWithActionRequired -}}
## Urgent Upgrade Notes 

### (No, really, you MUST read this before you upgrade)

{{range .}}{{println "-" .}} {{end}}
{{end}}

{{- if .Notes -}}
## Changes by Kind
{{ range .Notes}}
### {{.Kind | prettyKind}}

{{range $note := .NoteEntries }}{{println "-" $note}}{{end}}
{{- end -}}
{{- end }}

The images for this release are:
<ADD_IMAGE_HERE>

Thanks to all our contributors.
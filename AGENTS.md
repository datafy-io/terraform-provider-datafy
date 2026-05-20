# AGENTS.md

Guide for AI agents (and humans) working in this repo.

## What this is

`terraform-provider-datafy` — the official Terraform provider for [Datafy](https://docs.datafy.io). Built with [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework). Published to the Terraform Registry as `datafy-io/datafy`.

## Common commands

```
make build       # go build ./...
make install     # go install ./...
make test        # unit tests
make fmt         # gofmt -s -w -e .
make generate    # regenerate docs/ from templates/ + examples/ + Go schemas
make             # fmt + install + generate
```

## Editing docs — read this before touching any doc

**`docs/` is generated. Never edit files under `docs/` directly — `make generate` will overwrite them.**

The pipeline is:

```
templates/*.md.tmpl  +  examples/**/*.tf  +  Go schema Descriptions
                              │
                              ▼
                        make generate
                              │
                              ▼
                          docs/*.md
```

Where each piece of a generated doc comes from:

| Part of the doc                          | Source                                                  |
|------------------------------------------|---------------------------------------------------------|
| Page title, intro prose, Example Usage   | `templates/<kind>/<name>.md.tmpl`                       |
| `## Schema` section (attribute lists)    | `{{ .SchemaMarkdown }}` → `Description` fields in Go    |
| Inlined `terraform { ... }` code blocks  | Hand-written inside the `.tmpl` (preferred here)        |
| Import section                           | Hand-written in the `.tmpl`                             |

`docs/index.md` is special: its Schema section is hand-written in `templates/index.md.tmpl` because `{{ .SchemaMarkdown }}` alphabetises attributes (`endpoint` before `token`), and we want them in the order they're declared. If you add a new provider-level attribute, update both the Go schema in `internal/provider/provider.go` **and** the schema list in `templates/index.md.tmpl`.

### How to make a doc change

1. **Wording in a schema attribute description** → edit the `Description:` field in the relevant `internal/service/<name>/*.go`. Do not touch `docs/`.
2. **Prose / Example Usage / Import block** → edit `templates/<kind>/<name>.md.tmpl`. Do not touch `docs/`.
3. **New resource or data source** → add the package under `internal/service/<name>/`, register it in `internal/provider/provider.go`, add a `templates/<kind>/<name>.md.tmpl`, and add an example `.tf` under `examples/`.
4. **Run `make generate` and commit the resulting `docs/` changes alongside your source/template change in the same commit.**

### Templates are Go templates

`.tmpl` files are rendered by [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs). Useful directives:

- `{{ .SchemaMarkdown | trimspace }}` — auto-generated schema from the Go provider schema.
- `{{ tffile "examples/resources/foo/resource.tf" }}` — inline a Terraform code block from `examples/`.
- `{{ codefile "shell" "examples/resources/foo/import.sh" }}` — inline an arbitrary code block.

Note: `{{` and `}}` are template delimiters. Bare `{` and `}` (including consecutive `}}`) in markdown text are fine as long as they're not paired with a preceding `{{`.

## Always regenerate

**Every change that can affect docs must be followed by `make generate`, and the resulting `docs/` changes must be committed in the same change.** This applies to:

- Any edit to a `Description:` field in `internal/service/**/*.go` or `internal/provider/provider.go`.
- Adding, removing, or renaming a schema attribute.
- Any edit to a file under `templates/`.
- Any edit to a file under `examples/` that is referenced via `{{ tffile }}` or `{{ codefile }}` in a template.

If you forget, `docs/` drifts from the Go schemas and templates, and the Terraform Registry will show stale or wrong information. CI does not currently re-run generation, so this is an enforce-it-yourself rule.

A safe loop while iterating:

```
# edit Go / template / example
make generate
git diff docs/   # sanity-check what changed
```

## Conventions

- Match existing Go style: descriptions are full sentences with terminal periods.
- `Description:` strings in Go can contain Markdown — backticks, bold, links — and tfplugindocs will pass them through to `docs/`.
- Examples in `examples/` should be self-contained and `terraform fmt`-clean (`make generate` runs `terraform fmt -recursive examples/` first).
- Don't add a `CHANGELOG.md` entry unless the user asks — release notes are managed elsewhere.

## When in doubt

- The Go schema `Description` is the source of truth for attribute docs.
- The `.tmpl` file is the source of truth for everything else in a doc.
- `docs/` is never the source of truth.

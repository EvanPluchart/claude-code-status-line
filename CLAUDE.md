# claude-code-status-line — CLAUDE.md

> Ce fichier est le point d'entree pour l'IA. Il reference les conventions du projet et les docs koulia pertinentes.

## Projet

- **Nom** : claude-code-status-line
- **Description** : Statusline multi-ligne entierement configurable pour Claude Code
- **Stack** : Go 1.25, Bubbletea (TUI), GoReleaser, GitHub Actions
- **Repo** : Repo unique
- **Module** : `github.com/EvanPluchart/claude-code-status-line`

## Structure

```
cmd/claude-code-status-line/   # Point d'entree (main.go)
internal/
  ansi/                        # Helpers ANSI (couleurs, styles)
  config/                      # Chargement config YAML
  engine/                      # Moteur de rendu
  i18n/                        # Internationalisation (en, fr)
  parser/                      # Parsing du JSON stdin Claude Code
  themes/                      # Themes de couleurs (6 themes)
  widgets/                     # 18 widgets (model, cost, tokens, git, etc.)
  wizard/                      # Wizard interactif de configuration
docs/                          # Documentation projet
.github/workflows/             # CI (ci.yml) + Release (release.yml)
```

## Conventions

### Communes (koulia)

- Git : `~/koulia/common/git-conventions.md`
- Style : `~/koulia/common/code-style-rules.md`
- Nommage : `~/koulia/common/naming-conventions.md`

### Go (koulia)

- Index : `~/koulia/technos/go/index.md`

### Projet-specifiques

- Code conventions : `docs/CODE_CONVENTIONS.md`
- Plan d'implementation : `docs/IMPLEMENTATION_PLAN.md`

## Regles de code

- **Formatage** : `gofmt` + `goimports` (CI enforced)
- **Linter** : `golangci-lint` (errcheck, govet, staticcheck, unused, ineffassign, gosimple, misspell)
- **Imports** : 3 groupes separes par des lignes vides (stdlib / external / internal)
- **Booleens** : prefixes `is`, `has`, `should`, `can`, `will`
- **Espacement** : lignes vides avant/apres conditions, boucles, et avant `return`
- **Erreurs** : check immediat, early return, wrap avec `fmt.Errorf("context: %w", err)`
- **Early return** : privilegier les retours anticipes pour eviter le nesting

## Commandes

| Commande | Description |
|----------|-------------|
| `make build` | Build le binaire |
| `make test` | Lance les tests |
| `make lint` | Lance golangci-lint |
| `make clean` | Nettoie les artefacts |
| `make install` | Build + installe dans ~/.local/bin |
| `make release-dry` | Dry run GoReleaser |

## Contexte projet

### Branches

- **Production** : `main`
- **Format** : `type/description-courte`
- **Commits** : `Type - Description`

### Types de branches/commits

| Type | Usage |
|------|-------|
| `feature/` | Nouvelle fonctionnalite |
| `bugfix/` | Correction de bug |
| `hotfix/` | Fix urgent |
| `refactor/` | Refactoring |
| `chore/` | Maintenance |

### Points d'attention

- Le binaire doit rester rapide (~5ms) — pas de deps lourdes
- Cross-platform : tester les chemins et comportements sur macOS, Linux, Windows
- Le JSON stdin vient de Claude Code — parser defensivement
- 3 lignes max dans la statusline, 18 widgets disponibles
- Releases via GoReleaser + Homebrew tap

## Skills recommandes

- `/quick-dev` — Fix/feature rapide
- `/code-review` — Analyse de code
- `/feature-planner` — Planification de features
- `/sprint-planner` — Decoupe en sprint
- `/dev-story` — Implementation de stories

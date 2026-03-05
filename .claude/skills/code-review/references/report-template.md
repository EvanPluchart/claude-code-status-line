# Template de rapport — Code Review

Le rapport de code review doit suivre ce format exact.

## Format

```markdown
# Code Review — {source}

Date : {date}
Fichiers analyses : {N}
Stack : {stack detecte}

## Resume

{N} problemes trouves : {critique} critiques, {haute} hauts, {moyenne} moyens, {basse} bas

## Problemes

### CRITIQUE

| # | Fichier | Ligne | Regle | Description | Suggestion |
|---|---------|-------|-------|-------------|------------|
| 1 | src/... | L42   | ...   | ...         | ...        |

### HAUTE

| # | Fichier | Ligne | Regle | Description | Suggestion |
|---|---------|-------|-------|-------------|------------|
| 1 | src/... | L15   | ...   | ...         | ...        |

### MOYENNE

| # | Fichier | Ligne | Regle | Description | Suggestion |
|---|---------|-------|-------|-------------|------------|
| 1 | src/... | L8    | ...   | ...         | ...        |

### BASSE

| # | Fichier | Ligne | Regle | Description | Suggestion |
|---|---------|-------|-------|-------------|------------|
| 1 | src/... | L23   | ...   | ...         | ...        |

## Points positifs
- {ce qui est bien fait}
- {bonnes pratiques respectees}

## Actions recommandees
- [ ] Corriger les {N} problemes critiques/hauts avant merge
- [ ] Utiliser /fix-review-code pour appliquer les corrections automatiques
```

## Regles du rapport

- **{source}** : identifier la source du diff (ex: "commit abc123", "branch feature/xxx", "staged changes", "fichier src/...")
- **{date}** : date du jour au format YYYY-MM-DD
- **{stack detecte}** : le ou les stacks identifies (ex: "Next.js + TypeScript", "Symfony + API Platform")
- Si une categorie de severite n'a aucun probleme, ecrire "Aucun probleme detecte." sous le titre de la categorie
- La colonne **Regle** doit referencer la convention koulia violee (ex: "naming-conventions/camelCase", "symfony/code-style/return-type")
- La colonne **Suggestion** doit etre actionnable et concise
- Les **Points positifs** doivent mentionner au moins un element bien fait
- Les **Actions recommandees** doivent etre des checkboxes markdown

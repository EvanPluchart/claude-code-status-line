---
name: quick-dev
description: Implementer rapidement un petit fix ou une petite feature sans le workflow complet feature-planner / sprint-planner / dev-story.
---

# Quick Dev

## Objectif

Implementer rapidement un petit changement (fix, petite feature, refacto mineure) en un seul aller-retour quand le workflow complet (feature-planner → sprint-planner → dev-story) serait disproportionne.

## Parametres

- `--auto-commit` : Mode B — implementer puis committer automatiquement
- `--dry-run` : Mode C — analyser sans modifier aucun fichier

Sans parametre, le mode par defaut est le Mode A (implementation directe sans commit auto).

## Workflow

Suivre les etapes decrites dans `references/workflow.md`.

### Etape 1 — Comprendre la demande

Lire la description fournie par le user. Identifier clairement :
- Ce qui doit changer
- Ou dans le codebase
- Le type de changement (fix, feature, refacto)

### Etape 2 — Evaluer la complexite

Evaluer rapidement si le changement est adapte a quick-dev :
- Si le changement touche plus de 5-6 fichiers ou implique un nouveau domaine metier → suggerer `/feature-planner` a la place
- Si le changement est localise et bien defini → continuer

### Etape 3 — Adaptive context loading

Charger uniquement les docs koulia pertinentes selon le stack detecte :
- `~/koulia/technos/{stack}/` pour les conventions du stack concerne
- `~/koulia/common/naming-conventions.md` pour le nommage
- `~/koulia/common/code-style-rules.md` pour le style

Ne pas charger tout le repertoire koulia. Etre chirurgical.

### Etape 4 — Scanner les fichiers concernes

Utiliser les conventions du stack (depuis `~/koulia/technos/{stack}/`) pour localiser les fichiers a modifier. Lire ces fichiers pour comprendre le contexte existant.

### Etape 5 — Implementer

- Mode A (defaut) : Implementer les modifications directement avec les outils Edit/Write
- Mode B (`--auto-commit`) : Implementer puis committer automatiquement
- Mode C (`--dry-run`) : Decrire les modifications qui seraient faites sans toucher aux fichiers

Respecter les conventions et patterns koulia du stack concerne.

### Etape 6 — Verifier l'impact

Apres implementation (modes A et B) :
- Verifier que les imports ne sont pas casses
- Verifier la coherence des types
- Verifier que les tests existants ne sont pas impactes negativement
- Si un probleme est detecte, le corriger immediatement

### Etape 7 — Proposer le commit

- Mode A : Proposer un message de commit conventionnel selon `~/koulia/common/git-conventions.md` et attendre validation du user
- Mode B : Creer le commit automatiquement avec un message conventionnel
- Mode C : Aucun commit

## Principes

- **Efficacite maximale** : un seul aller-retour si possible. Pas de plan elabore, pas de stories.
- **Adaptive context loading** : ne charger que les docs koulia strictement necessaires au changement.
- **Atomic commit** : un seul commit propre qui represente le changement complet.
- **Savoir dire non** : si le changement est trop gros pour quick-dev, rediriger vers le bon workflow.

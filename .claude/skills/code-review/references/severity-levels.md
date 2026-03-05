# Niveaux de severite — Code Review

## CRITIQUE

Problemes qui doivent etre corriges **avant tout merge**. Ils representent un risque immediat.

Exemples :
- Bugs logiques (condition inversee, boucle infinie, off-by-one)
- Failles de securite (injection SQL, XSS, CSRF, exposition de secrets)
- Perte de donnees potentielle (suppression sans confirmation, ecrasement silencieux)
- Crash applicatif (null pointer non gere, exception non catchee dans un contexte critique)
- Regression fonctionnelle evidente

## HAUTE

Problemes qui doivent etre corriges **avant merge** sauf justification explicite.

Exemples :
- Violations de types (any implicite, cast dangereux, type manquant sur une API publique)
- API incorrecte (mauvais status code, payload incoherent, contrat non respecte)
- Regression probable (changement de comportement non documente)
- Violation de pattern architectural (contournement du pattern koulia sans justification)
- Probleme de performance evident (requete N+1, boucle inutile sur grande collection)
- Absence de validation sur des donnees utilisateur

## MOYENNE

Problemes qui devraient etre corriges mais ne bloquent pas le merge.

Exemples :
- Violations de conventions de nommage (selon `~/koulia/common/naming-conventions.md`)
- Violations de code style (selon `~/koulia/technos/{stack}/code-style.md`)
- Organisation des imports non conforme
- Fichier mal place dans l'arborescence
- Commentaire obsolete ou trompeur
- Duplication de code mineure

## BASSE

Suggestions d'amelioration. Le code fonctionne correctement mais pourrait etre ameliore.

Exemples :
- Optimisation possible (simplification d'une expression, utilisation d'une methode plus adaptee)
- Meilleure lisibilite (extraction de variable, renommage plus explicite)
- Ajout de commentaire utile
- Refactoring mineur pour mieux coller au pattern koulia
- Amelioration de la couverture de test

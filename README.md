# projet_IA04

## Participants

- LECLERCQ Vivien
- SOTO Louis 
- CONSTANTIN Alexandre 
- DESCAMPS Jean

### Comportements émergents et modélisation de la dynamique d'un ecosysteme

#### Enjeux 

- Rendu artistique 

- Comprendre les méchanismes régissants la dynamique d'un ecosysteme. 

Les questions auxquelles nous souhaitons répondre à travers cette simulation sont les suivantes : 

- Sous quelles conditions le système est-il stable ? 
- Existe-t-il une hétérogénité dans les sous-populations qui survivent ?

Nous utiliserons la musique comme source de perturbations, l'utilisation de celle-ci nous parait intéressante par l'information riche qu'elle procure (projet etageable)

Nous commençons le projet par l'implémentation d'agents réactifs et souhaitons améliorer la complexité du projet en implémentant des agents inductifs. (pour boid par exemple : Est-ce un choix raisonné de rester dans le meme flock)
 
Nous divergeons pour le moment encore sur le choix  pour l'implémentation de base (boids ou slime) pour simuler notre ecosysteme, les deux projets s'influencent mutuellement pour le moment.

#### Boids 
https://fr.wikipedia.org/wiki/Boids

#### Slime
https://uwe-repository.worktribe.com/output/980579

## Architecture du projet 

### Modules 

- modélisation des agents (agents boids, agents perturbateurs)
- modélisation du monde (runtime) -> Game 
- rendu graphique 
- analyse (macro, statistique)

### Etapes 

- Reproduire Boid tel quel (fait) ->   
- refactoriser le code 
- varier les 3 forces -> voir ce qu'on obtient (equilibre ?) 
- injecter les agents perturbateurs -> étager le comportement -> quel comportement/quel ratio d'agents pour les equilibres qu'on a vu au dessus. 
- Analyse 

### A faire 

-refactorisation du code Vivien
- gerer les effets de bord 
- modélisation des états du boid pour tous les ratios de force possible 


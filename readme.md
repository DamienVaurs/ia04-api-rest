# IA04 - API Web pour différentes méthodes de vote

## Installation et lancement

Le projet peut être clôné depuis le dépôt git suivant :

    https://gitlab.utc.fr/milairhu/ia04-api-rest.git

Sans clôner le projet, les différentes fonctions exécutables, située dans le dossieur */restagent/cmd*, peuvent être téléchargée via la commande suivante :

    go get gitlab.utc.fr/milairhu/ia04-api-rest/restagent/cmd/launch-10-generated-agents@latest
où *launch-10-generated-agents* peut être remplacé par le nom de n'importe quel autre dossier de *cmd* exécutable.
L'utilisation de **@latest** n'est obligatoire que si le dossier actif n'est pas dans un module Go.

Pour lancer un programme via *go run*, il suffit de se placer dans le dossier *cmd* et d'exécuter la commande suivante :

    go run launch-10-generated-agents.go

Pour lancer un exécutable installé via *go install*, il suffit de se placer dans le dossier *go/bin*, généralement dans le home directory, et de lancer la commande de la sorte :

    ./launch-10-generated-agents

A noter qu'il se peut que le programme échoue, l'agent serveur n'arrivant pas ponctuellement à se connecter au port 8080. Il faut alors relancer le programme.

## Structure générale (packages)

### Dossier cmd

La dossier *restagent/cmd/* regroupe tous les fichiers exécutables utilisés pour tester l'ensemble du projet. Certains de ces fichiers reprennent les exemples vus en cours.

- *launch-10-generated-agents.go* : lance 10 agents votants générés aléatoirement, afin de tester chaque méthode de vote implémentée et quelques cas limites.
- *launch-x-generated-agents.go* : Similaire au précédent, à la différence qu'il est demandé à l'utilisateur de fournir le nombre de votants, scrutins et alternatives. Pratique pour tester des cas de figures avec un nombre très important d'agents. Attention : aucun cas limite n'est généré (deadline déjà passée, votant qui n'a pas le droit de voter, etc). Les méthodes de vote de chaqe scrutin sont choisies aléatoirement.
- *launch-approval.go*, *launch-condorcet.go* et *launch-stv.go* : permettent de tester la méthode par Approbation, de Condorcet et STV avec et sans besoin de tie-break, leur manipulation étant différente des autres méthodes.
- *launch-rsagt.go* : lance un serveur REST qui gère les requêtes entrantes sur le port 8080. C'est la commande à lancer si l'utilisateur souhaite testr l'API via un outil comme Postman.
- *launch-rcagt.go* : lance un client REST qui envoie des requêtes au serveur REST lancé précédemment. Il lance un agent créateur de scrutin et un agent votant simples.
- les commandes dans les fichiers *launch-chap2-diapX.go* permettent de tester les exemples vus en cours.

### Package comsoc

Ce package (*dossier /restagent/comsoc/*) contient toutes les classes, méthodes et types relatifs à la gestion des scrutins. C'est ici que nous retrouvons les fonctions de calcul des SWF et SCF pour les méthodes de Borda, Condorcet, Copeland, etc.

On y trouve également un certain nombre de fonctions utilitaires regroupées dans le *fichier /comsoc/basics.go*.

Enfin, le *fichier /comsoc/tiebreak.go* contient les fonctions de type **factory** permettant de créer des fonctions de tie-break pour les différentes méthodes. Seules les tie-breaks pour STV et Approbation ont du être implémentées à la main, leur utilisation étant différente des autres méthodes.

### Package endpoints

Endpoints (*dossier /restagent/endpoints/*) est un package composé d'un seul *fichier /endpoints/endpoints.go* dont le but est de définir certaines constantes utilisées dans tout le reste du projet. On y trouve les éléments pour construire les URL des requêtes HTTP.

### Package instances

Le package instances (*dossier restagent/instances/*) contient tous les fichiers d'initialisation des exécutables du dossier *cmd*, et portent les mêmes noms que ces derniers.

Les fichiers *init-....go* contiennent chacun une fonction instanciant les agents votants et gérants de scrutins pour l'exécutable associé.

Le fichier *launch-agents.go* contient l'abstraction de la fonction lançant les agents. Comme à l'origine, pour chaque exécutable, les fonctions **main** étaient à peu de choses identiques, on a décidé de les abstraires dans ce fichier.

### Package restclientagent

Dans ce package (*dossier /restagent/restclientagent/*), on trouve la définition des agents côté client (votants et gérants de scrutins) ainsi que les méthodes utilisées pour réaliser les différents requêtes HTTP.

### Package restserveragent

Ce package (*dossier /restagent/restserveragent/*), similaire dans sa conception au précédent, définit toutes les classes et méthodes côté serveur. Les fonctions de celui-ci permettent à l'agent serveur de communiquer avec les agents clients du package *restclientagent* via requêtes HTTP.

## Package restagent

Le package restagent, situé à la racine du projet, définit un certain nombre de types (*fichier /types.go*) et de constantes (*fichier /rule.go*) utilisées par les agents clients et serveurs.

## Remarques

- Pour la méthode de scrutin par approbation, une nouvelle méthode de tie-break a dû être créée, afin de prendre en compte les seuils (fonction *MakeApprovalRankingWithTieBreak()* dans le fichier */comsoc/tiebreak.go*).
- Même chose pour la méthode par STV, qui nécessite une fonction de tie-break à part car le départage se fait au sein même de la fonction de calcul du SWF, et pas après (fonction *STV_SWF_TieBreak* dans le fichier */comsoc/tiebreak.go*).
- Lors de la création d'un scrutin, on ne fait pas appel à log.Fatal car on souhaite que l'agent continue ses tâches même si une erreur est rencontrée.
- Le fait de vérifier la cohérence des seuils fournis au moment du calcul du résultat (fichier */restserveragent/result.go*) présente un avantage sécuritaire mais pénalise la performance, car ces seuils sont déjà vérifiés à la réception du vote.
- On peut se poser la question de l'intérêt (ou non) de vérifier la présence d'un seuil dans le cadre d'une méthode de scrutin par approbation. Est-ce que l'absence de celui-ci est une erreur ? Ou bien cela signifie-t-il qu'on compte toutes les alternatives, ou aucune? On a décidé de considérer que l'absence de seuil est une erreur.
- Dans le cas où aucun vote n'est soumis (fichier */restserveragent/result.go*), nous prenons la décision de retourner un résultat plutôt qu'une erreur. Ce résultat est déterminé par le tie-break fourni lors de la création du scrutin.
- La méthode de Condorcet ne nécessite pas l'utilisation d'une fonction de gestion de tie-break. Soit il y a un unique vainqueur, soit il n'y en a pas.
- Le nombre d'alternatives dans le *fichier /cmd/launch-rcagt.go* a été arbitrairement fixé à 5, mais est modifiable.
- Le "profil des votants" affiché a la même forme que ceux vu en cours. La première ligne indique le nombre de votants a avoir cet ordre de préférence, indiqué par la colonne en dessous.

Hugo MILAIR,
Damien VAURS,
Semestre A23, 30/10/2023

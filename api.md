# Envoi

Sauf si le contraire est indiqué, la réponse attendue pour toutes les requêtes est:
```json
{
    "temps": 10000 
}
```
avec "temps" le nombre de secondes écoulées depuis minuit

### Première requête, effectuée à la 1ère connection:
```json
{
    "type": "bonjour"
}
```
Réponse attendue si pas de restauration de sauvegarde:
```json
{
    "sauvegarde": false,
    "restos": [
        {
            "temps": 0,
            "acceleration": 60,
            "horaires": [
                [12, 15],
                [19, 22]
            ],
            "entrees": [
                "Feuilleté au crabe",
                "Œufs cocotte",
                "Tarte au thon",
                "Pics apéritifs de roulés de crêpes au saumon",
                "Salades de pâtes au thon",
                "Œufs à la coque",
                "Carotte à l'orange",
                "Pommes de terre surprise",
                "Cornets de saumon fumé",
                "Brochettes melon et jambon"
            ],
            "plats": [
                "Crêpes poulet béchamel",
                "Burger steack bacon",
                "BLANCS DE POULET A LA CREME ET AU MIEL",
                "Burger double steack",
                "Burger saumon",
                "steak frites",
                "Nuggets frites",
                "Spaguetti bolognaise",
                "Fish ans chips",
                "Burger fish"
            ],
            "desserts": [
                "Crêpes",
                "TIRAMISU",
                "MADELEINE AU MIEL",
                "Gateau fondant au chocolat",
                "Framboise et citron",
                "Crème dessert légère",
                "Mousse au chocolat",
                "Pancakes",
                "Salade de fruits d'été ",
                "Glace chocolat"
            ],
            "carres": [
                {
                    "2": 5,
                    "4": 5,
                    "6": 3,
                    "8": 2,
                    "10": 1
                },
                {
                    "2": 5,
                    "4": 5,
                    "6": 2,
                    "8": 3,
                    "10": 1
                }
            ]
        }
    ]
}
```
Sinon: 
```json
{
    "sauvegarde": true,
    "etat": "dump de l'objet restaurant"
}
```

### Sauvegarde: 
```json
{
    "type": "sauvegarde",
    "etat": "dump de l'objet restaurant"
}
```
### Accélération du temps:
```json
{
    "type": "acceleration",
    "val": 60
}
```

### Demande de pause:
```json
{
    "type": "pause",
    "pause": true
}
```
ou "false" pour remettre le resto en marche

### Commandes:
Les clients commandent les entrées, les plats et les desserts en même temps,
mais ils ne sont pas livrés en même temps.
L'ID est nécessaire pour savoir à qui envoyer les plats et les desserts plus tard.
```json
{
    "type": "commande",
    "id": 12,
    "commande": {
        "entrees": [
            "tableau de strings avec le nom des entrées"
        ],
        "plats": [
            "tableau de strings"
        ],
        "desserts": [
            "tableau de strings"
        ]
    }
}

```

### Retour du matériel commun:
L'ID est l'ID de la commande du client qui a utilisé le matériel. 
Par exemple, la cuisine sait qu'elle doit envoyer le plat quand le matériel retourné 
a l'ID d'une commande qui a seulement déjà reçu une entrée.
```json
{
    "type": "materiel",
    "id": 12,
    "materiel": [
        {
            "type": "fourchette",
                "qtt": 12
        },
        {
            "type": "assiettePlate",
            "qtt": 2
        }
    ]
}
```
Liste de tous les retours de matériel commun possible:
* assiettePetite
* assiettePlate
* assieteCreuse
* assietteDessert
* fourchette
* couteau
* cuillereSoupe
* verreEau
* verreVin
* verreChampagne
* setCafe
* serviette
* nappe


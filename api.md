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
    "sauvegarde": "false",
    "restos": [
        {
            "entrees": [
                "tableau de strings contenant toutes les entrées"
            ],
            "plats": [
                "tableau de strings"
            ],
            "desserts": [
                "tableau de strings"
            ]
        }
    ],
    "temps": 0
}
```
Sinon: 
```json
{
    "sauvegarde": "true",
    "temps": 0,
    "etat": " dump de l'objet restaurant"
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
    "val": 60"
}
```

### commandes:
```json
{
    "type": "commande",
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
```json
{
    "type": "materiel",
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

### Demande de pause:
```json
{
    "type": "pause",
    "pause": "true"
}
```
ou "false" pour remettre le resto en marche

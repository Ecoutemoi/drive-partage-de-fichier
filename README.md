# 🚀 Drive File Microservice

Microservice de gestion de fichiers développé en **Go (net/http)** avec une interface **Flutter (MVC + Provider)**.

Projet réalisé dans le cadre d’un stage et conçu pour être intégré dans une architecture microservices.

---

##  Description

Ce service permet :

- La gestion hiérarchique des dossiers
- L’upload et le téléchargement de fichiers
- Le renommage et la suppression
- Le déplacement de fichiers et dossiers
- La génération de liens de partage sécurisés

L’authentification (JWT, login, gestion des sessions) est volontairement externalisée et gérée par un autre service dans l’architecture globale.

---

##  Stack Technique

### Backend
- Go (net/http)
- MySQL
- API REST JSON
- Middleware CORS
- Configuration via variables d’environnement (.env)

### Frontend
- Flutter
- Provider
- Architecture MVC
- Interface Material 3

---

##  Fonctionnalités

###  Dossiers
- Création
- Renommage
- Suppression
- Déplacement (relation parent/enfant)
- Navigation via fil d’Ariane

###  Fichiers
- Upload
- Téléchargement
- Renommage
- Suppression
- Déplacement entre dossiers

###  Partage
- Génération de lien via token
- Accès sécurisé au fichier

---

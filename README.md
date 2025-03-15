# Hulp met github :D

## 1. **Een Repository Instellen**

### Stap 1: Clone de Repository

```bash
git clone https://github.com/zyphcore/BPV-Matcher.git
```

## 2. **Werken met Git**

### Stap 1: Een nieuwe branch maken

Voordat je begint met werken aan een functie of bugfix, is het een goed idee om een nieuwe branch aan te maken. Dit voorkomt dat je per ongeluk direct op de `main` branch werkt.

Maak een nieuwe branch met:

```bash
git checkout -b nieuwe-branch-naam
```

### Stap 2: Wijzigingen maken en committen

Zodra je wijzigingen hebt aangebracht in de bestanden, kun je deze toevoegen aan de staging area en committen:

1. **Bestanden toevoegen** (stage):

```bash
git add .
```

2. **Wijzigingen committen**:

```bash
git commit -m "Beschrijving van de wijziging"
```

### Stap 3: Push naar GitHub

Als je klaar bent met je wijzigingen, kun je ze naar GitHub pushen:

```bash
git push origin nieuwe-branch-naam
```

Dit stuurt je wijzigingen naar de nieuwe branch op de remote repository op GitHub.

### Stap 4: Werken met de `main` branch

Als je terug wilt naar de `main` branch, gebruik je de volgende opdracht:

```bash
git checkout main
```

Om ervoor te zorgen dat je de laatste wijzigingen van de `main` branch hebt, kun je deze updaten met:

```bash
git pull origin main
```

## 3. **Pull Requests (PR's)**

### Stap 1: Maak een Pull Request aan

Wanneer je klaar bent met je wijzigingen en je deze wilt integreren in de `main` branch (of een andere branch), maak je een **pull request** aan op GitHub:

1. Ga naar de **"Pull requests"** sectie van je repository.
2. Klik op **"New Pull Request"**.
3. Kies de branch die je hebt gemaakt en de `main` branch (of andere doelbranch) als de branch waar je naar wilt mergen.
4. Voeg een titel en beschrijving toe aan de pull request.
5. Klik op **"Create Pull Request"**.

### Stap 2: Code Review en Mergen

Een teamlid (of jijzelf) moet de pull request controleren. Als alles goed is, kan de pull request gemerged worden in de `main` branch door op de **"Merge Pull Request"** knop te klikken.

### Stap 3: Verwijder de Branch (optioneel)

Na het mergen van de pull request, kun je de gemaakte branch verwijderen om de repository schoon te houden:

```bash
git branch -d nieuwe-branch-naam  # Verwijder lokaal
git push origin --delete nieuwe-branch-naam  # Verwijder op GitHub
```

## 4. **Samenwerken met Collega's**

Wanneer je samenwerkt met anderen aan een project, zijn er een paar basisprincipes:

- **Pull vaak**: Voordat je begint met werken, haal altijd de laatste versie van de `main` branch op om te zorgen dat je up-to-date bent:
  
  ```bash
  git pull origin main
  ```

- **Werk in een eigen branch**: Maak altijd een nieuwe branch voor elke taak die je uitvoert. Dit houdt de geschiedenis schoon en maakt het makkelijker om samen te werken.

- **Push vaak**: Het is een goed idee om je werk regelmatig naar de remote repository te pushen, zodat anderen je wijzigingen kunnen zien en je geen werk verliest.

## 5. **Problemen oplossen**

### Conflicten oplossen

Als je een conflict hebt bij het mergen, zul je de conflicterende bestanden handmatig moeten bewerken om de veranderingen te combineren. Git markeert de conflicten in de bestanden, zodat je kunt zien waar de problemen zitten.

Na het oplossen van de conflicten, voeg je de gewijzigde bestanden toe en commit je de wijzigingen:

```bash
git add conflicterend-bestand
git commit -m "Conflicten opgelost"
```

### Force push (alleen in specifieke gevallen)

Als je absoluut zeker weet dat je de geschiedenis wilt overschrijven (bijvoorbeeld na een rebase of andere noodzaak), kun je een **force push** doen. Gebruik dit met voorzichtigheid:

```bash
git push --force origin jouw-branch-naam
```
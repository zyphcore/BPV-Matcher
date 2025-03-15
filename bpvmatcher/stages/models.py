from django.db import models
from django.contrib.auth.models import User # Importeer het User model voor gebruikersbeheer

class Student(models.Model):
    # Relatie met het User model voor authenticatie
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    # Extra velden voor studentprofiel
    programmeertalen = models.TextField(blank=True, null=True)
    werkvoorkeuren = models.TextField(blank=True, null=True)
    portfolio = models.URLField(blank=True, null=True)
    ervaring = models.TextField(blank=True, null=True)
    cv = models.FileField(upload_to='cvs/', blank=True, null=True)
    motivatiebrief = models.FileField(upload_to='motivatiebrieven/', blank=True, null=True)

    def __str__(self):
        return self.user.username # Toon de gebruikersnaam in de admin interface

class Bedrijf(models.Model):
    naam = models.CharField(max_length=200)
    locatie = models.CharField(max_length=200)
    technologieen = models.TextField(blank=True, null=True)
    contactpersoon = models.CharField(max_length=200, blank=True, null=True)
    email = models.EmailField(blank=True, null=True)
    telefoonnummer = models.CharField(max_length=20, blank=True, null=True)
    # Extra velden voor status en beoordeling
    is_samenwerkingsbedrijf = models.BooleanField(default=False)
    is_potentieel_bedrijf = models.BooleanField(default=False)
    beoordeling = models.FloatField(blank=True, null=True) # Gemiddelde beoordeling van studenten

    def __str__(self):
        return self.naam

class Stage(models.Model):
    bedrijf = models.ForeignKey(Bedrijf, on_delete=models.CASCADE)
    titel = models.CharField(max_length=200)
    omschrijving = models.TextField()
    technologieen = models.TextField(blank=True, null=True)
    startdatum = models.DateField(blank=True, null=True)
    einddatum = models.DateField(blank=True, null=True)
    beschikbare_plaatsen = models.IntegerField(default=1) # Aantal beschikbare stageplaatsen

    def __str__(self):
        return self.titel

class Benadering(models.Model):
    student = models.ForeignKey(Student, on_delete=models.CASCADE)
    stage = models.ForeignKey(Stage, on_delete=models.CASCADE)
    datum = models.DateField(auto_now_add=True) # Datum van benadering (automatisch ingevuld)
    contactwijze = models.CharField(max_length=100)
    cv = models.FileField(upload_to='benaderingen/cvs/', blank=True, null=True)
    motivatiebrief = models.FileField(upload_to='benaderingen/motivatiebrieven/', blank=True, null=True)
    opmerkingen = models.TextField(blank=True, null=True)
    status = models.CharField(max_length=100, default='Benaderd') # Status van de sollicitatie

    def __str__(self):
        return f'{self.student.user.username} - {self.stage.titel}'

class Coördinator(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    # Extra velden voor coördinatorprofiel (indien nodig)

    def __str__(self):
        return self.user.username

class Mentor(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    studenten = models.ManyToManyField(Student, blank=True) # Mentoren kunnen meerdere studenten begeleiden

    def __str__(self):
        return self.user.username

class Feedback(models.Model):
    student = models.ForeignKey(Student, on_delete=models.CASCADE)
    stage = models.ForeignKey(Stage, on_delete=models.CASCADE)
    beoordeling = models.IntegerField() # Beoordeling van 1 tot 5
    opmerkingen = models.TextField(blank=True, null=True)

    def __str__(self):
        return f'Feedback van {self.student.user.username} voor {self.stage.titel}'
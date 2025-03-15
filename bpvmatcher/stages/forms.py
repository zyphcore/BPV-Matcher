from django import forms
from .models import Stage, Benadering, Bedrijf, Student
from django.contrib.auth.models import User

class UserRegistrationForm(forms.ModelForm):
    password = forms.CharField(widget=forms.PasswordInput)
    password2 = forms.CharField(widget=forms.PasswordInput, label='Herhaal wachtwoord')

    class Meta:
        model = User
        fields = ('username', 'first_name', 'email')

    def clean_password2(self):
        cd = self.cleaned_data
        if cd['password'] != cd['password2']:
            raise forms.ValidationError('Wachtwoorden komen niet overeen.')
        return cd['password2']

class StudentProfileForm(forms.ModelForm):
    class Meta:
        model = Student
        fields = ['programmeertalen', 'werkvoorkeuren']

class StageForm(forms.ModelForm):
    class Meta:
        model = Stage
        fields = ['bedrijf', 'titel', 'omschrijving', 'technologieen', 'startdatum', 'einddatum', 'beschikbare_plaatsen']

class BenaderingForm(forms.ModelForm):
    class Meta:
        model = Benadering
        fields = ['contactwijze', 'cv', 'motivatiebrief', 'opmerkingen', 'status']

class BedrijfForm(forms.ModelForm):
    class Meta:
        model = Bedrijf
        fields = ['naam', 'locatie', 'technologieen', 'contactpersoon', 'email', 'telefoonnummer', 'is_samenwerkingsbedrijf', 'is_potentieel_bedrijf', 'beoordeling']

class StudentForm(forms.ModelForm):
    class Meta:
        model = Student
        fields = ['programmeertalen', 'werkvoorkeuren', 'portfolio', 'ervaring'] 
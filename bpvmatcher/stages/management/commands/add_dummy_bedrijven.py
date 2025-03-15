from django.core.management.base import BaseCommand
from stages.models import Bedrijf

class Command(BaseCommand):
    help = 'Voegt dummy bedrijven toe aan de database'

    def handle(self, *args, **options):
        bedrijven = [
            {
                'naam': 'Tech Innovators BV',
                'locatie': 'Amsterdam',
                'technologieen': 'Python, Django, React',
                'contactpersoon': 'Jan de Vries',
                'email': 'info@techinnovators.nl',
                'telefoonnummer': '020-1234567',
                'is_samenwerkingsbedrijf': True,
                'is_potentieel_bedrijf': False,
                'beoordeling': 4.5,
            },
            {
                'naam': 'Web Solutions NV',
                'locatie': 'Rotterdam',
                'technologieen': 'JavaScript, Angular, Node.js',
                'contactpersoon': 'Piet Jansen',
                'email': 'contact@websolutions.nl',
                'telefoonnummer': '010-9876543',
                'is_samenwerkingsbedrijf': False,
                'is_potentieel_bedrijf': True,
                'beoordeling': 3.8,
            },
            {
                'naam': 'Data Analytics BV',
                'locatie': 'Utrecht',
                'technologieen': 'Python, Pandas, Scikit-learn',
                'contactpersoon': 'Klaas van Dijk',
                'email': 'info@dataanalytics.nl',
                'telefoonnummer': '030-5678901',
                'is_samenwerkingsbedrijf': True,
                'is_potentieel_bedrijf': True,
                'beoordeling': 4.2,
            },
        ]

        for bedrijf_data in bedrijven:
            Bedrijf.objects.create(**bedrijf_data)

        self.stdout.write(self.style.SUCCESS('Dummy bedrijven succesvol toegevoegd!'))
from django.core.management.base import BaseCommand
from stages.models import Stage, Bedrijf
from datetime import date

class Command(BaseCommand):
    help = 'Voegt dummy stages toe aan de database'

    def handle(self, *args, **options):
        # Haal de bedrijven op die je eerder hebt toegevoegd
        bedrijven = Bedrijf.objects.all()

        if not bedrijven:
            self.stdout.write(self.style.WARNING('Geen bedrijven gevonden. Voeg eerst bedrijven toe.'))
            return

        stages = [
            {
                'bedrijf': bedrijven[0], # Gebruik het eerste bedrijf
                'titel': 'Stage Webontwikkeling',
                'omschrijving': 'Ontwikkel een webapplicatie met Django en React.',
                'technologieen': 'Python, Django, React, JavaScript',
                'startdatum': date(2024, 9, 1),
                'einddatum': date(2024, 12, 1),
                'beschikbare_plaatsen': 2,
            },
            {
                'bedrijf': bedrijven[1], # Gebruik het tweede bedrijf
                'titel': 'Stage Frontend Ontwikkeling',
                'omschrijving': 'Bouw een interactieve frontend met Angular.',
                'technologieen': 'JavaScript, Angular, TypeScript',
                'startdatum': date(2024, 10, 1),
                'einddatum': date(2025, 1, 1),
                'beschikbare_plaatsen': 1,
            },
            {
                'bedrijf': bedrijven[2], # Gebruik het derde bedrijf
                'titel': 'Stage Data Analyse',
                'omschrijving': 'Analyseer data met Python en Pandas.',
                'technologieen': 'Python, Pandas, Scikit-learn',
                'startdatum': date(2024, 11, 1),
                'einddatum': date(2025, 2, 1),
                'beschikbare_plaatsen': 3,
            },
        ]

        for stage_data in stages:
            Stage.objects.create(**stage_data)

        self.stdout.write(self.style.SUCCESS('Dummy stages succesvol toegevoegd!'))
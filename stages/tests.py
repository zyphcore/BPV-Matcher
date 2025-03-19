from django.test import TestCase, Client
from .models import Student, Bedrijf, Stage, Benadering, Coördinator, Mentor, Feedback
from django.contrib.auth.models import User
from django.urls import reverse

class ModelTests(TestCase):

    def setUp(self):
        # Maak testdata aan die gebruikt kan worden in de tests
        self.bedrijf = Bedrijf.objects.create(naam='Test Bedrijf', locatie='Test Locatie')
        self.stage = Stage.objects.create(bedrijf=self.bedrijf, titel='Test Stage', omschrijving='Test Omschrijving')
        self.user = User.objects.create_user(username='teststudent', password='testpassword')
        self.student = Student.objects.create(user=self.user)

    def test_bedrijf_aanmaken(self):
        # Test of een bedrijf correct wordt aangemaakt
        self.assertEqual(self.bedrijf.naam, 'Test Bedrijf')
        self.assertEqual(self.bedrijf.locatie, 'Test Locatie')

    def test_stage_aanmaken(self):
        # Test of een stage correct wordt aangemaakt
        self.assertEqual(self.stage.titel, 'Test Stage')
        self.assertEqual(self.stage.bedrijf, self.bedrijf)

    def test_student_aanmaken(self):
        # Test of een student correct wordt aangemaakt
        self.assertEqual(self.student.user.username, 'teststudent')

class ViewTests(TestCase):

    def setUp(self):
        # Maak een client aan om HTTP-verzoeken te simuleren
        self.client = Client()
        self.bedrijf = Bedrijf.objects.create(naam='Test Bedrijf', locatie='Test Locatie')
        self.stage = Stage.objects.create(bedrijf=self.bedrijf, titel='Test Stage', omschrijving='Test Omschrijving')
        self.user = User.objects.create_user(username='teststudent', password='testpassword')
        self.student = Student.objects.create(user=self.user)
        self.client.login(username='teststudent', password='testpassword') # Log in als teststudent

    def test_stage_lijst_view(self):
        # Test of de stage lijst view correct werkt
        response = self.client.get(reverse('stage_lijst'))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, 'Stageoverzicht')

    def test_stage_detail_view(self):
        # Test of de stage detail view correct werkt
        response = self.client.get(reverse('stage_detail', args=[self.stage.id]))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, 'Test Stage')

    def test_stage_benaderen_view(self):
        # Test of de stage benaderen view correct werkt
        response = self.client.get(reverse('stage_benaderen', args=[self.stage.id]))
        self.assertEqual(response.status_code, 200)
        self.assertContains(response, 'Stage Benaderen')
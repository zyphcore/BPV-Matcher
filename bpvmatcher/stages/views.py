from django.shortcuts import render, get_object_or_404, redirect
from django.http import HttpResponse
from .models import Student, Bedrijf, Stage, Benadering, Coördinator, Mentor, Feedback
from django.contrib.auth.decorators import login_required
from django.contrib.auth.models import User
from .forms import BenaderingForm, StageForm, BedrijfForm, StudentForm
from .forms import UserRegistrationForm
from django.contrib.auth import login
from .forms import StudentProfileForm

@login_required
def student_profiel_bewerken(request):
    """Formulier om programmeertalen en werkvoorkeuren aan te passen."""
    student = request.user.student
    if request.method == 'POST':
        form = StudentProfileForm(request.POST, instance=student)
        if form.is_valid():
            form.save()
            return redirect('student_profiel')
    else:
        form = StudentProfileForm(instance=student)
    return render(request, 'stages/student_profiel_bewerken.html', {'form': form})

def register(request):
    if request.method == 'POST':
        form = UserRegistrationForm(request.POST)
        if form.is_valid():
            new_user = form.save(commit=False)
            new_user.set_password(form.cleaned_data['password'])
            new_user.save()
            Student.objects.create(user=new_user)
            login(request, new_user)
            return redirect('home')
    else:
        form = UserRegistrationForm()
    return render(request, 'registration/register.html', {'form': form})

# Algemene views

def home(request):
    """Startpagina van de applicatie."""
    return render(request, 'stages/home.html')

# Student views

@login_required
def student_profiel(request):
    """Profielpagina van de student."""
    student = Student.objects.get(user=request.user)
    return render(request, 'stages/student_profiel.html', {'student': student})

@login_required
def student_profiel_bewerken(request):
    """Formulier om het studentenprofiel te bewerken."""
    student = request.user.student
    if request.method == 'POST':
        form = StudentForm(request.POST, instance=student)
        if form.is_valid():
            form.save()
            return redirect('student_profiel')
    else:
        form = StudentForm(instance=student)
    return render(request, 'stages/student_profiel_bewerken.html', {'form': form})

@login_required
def stage_lijst(request):
    """Lijst van alle stages."""
    stages = Stage.objects.all()
    return render(request, 'stages/stage_lijst.html', {'stages': stages})

@login_required
def stage_detail(request, stage_id):
    """Detailpagina van een stage."""
    stage = get_object_or_404(Stage, pk=stage_id)
    return render(request, 'stages/stage_detail.html', {'stage': stage})

@login_required
def bedrijf_detail(request, bedrijf_id):
    """Detailpagina van een bedrijf."""
    bedrijf = get_object_or_404(Bedrijf, pk=bedrijf_id)
    return render(request, 'stages/bedrijf_detail.html', {'bedrijf': bedrijf})

@login_required
def stage_benaderen(request, stage_id):
    """Formulier om een stage te benaderen."""
    stage = get_object_or_404(Stage, pk=stage_id)
    if request.method == 'POST':
        form = BenaderingForm(request.POST, request.FILES) # Voeg request.FILES toe voor bestandsuploads
        if form.is_valid():
            benadering = form.save(commit=False)
            benadering.student = request.user.student
            benadering.stage = stage
            benadering.save()
            return redirect('stage_detail', stage_id=stage.id)
    else:
        form = BenaderingForm()
    return render(request, 'stages/stage_benaderen.html', {'form': form, 'stage': stage})

# Coördinator views

@login_required
def coordinator_dashboard(request):
    """Dashboard voor de stagecoördinator."""
    if hasattr(request.user, 'coördinator'):
        coördinator = request.user.coördinator
        return render(request, 'stages/coordinator_dashboard.html', {'coordinator': coördinator})
    else:
        return redirect('home') # Redirect naar de home pagina als de gebruiker geen coördinator is

@login_required
def stage_toevoegen(request):
    """Formulier om een stage toe te voegen."""
    if request.user.is_staff:
        if request.method == 'POST':
            form = StageForm(request.POST)
            if form.is_valid():
                form.save()
                return redirect('stage_lijst')
        else:
            form = StageForm()
        return render(request, 'stages/stage_toevoegen.html', {'form': form})
    else:
        return redirect('stage_lijst')
    
@login_required
def bedrijf_lijst(request):
    bedrijven = Bedrijf.objects.all() # Zorg dat Bedrijf geimporteerd is.
    return render(request, 'stages/bedrijf_lijst.html', {'bedrijven': bedrijven})

@login_required
def bedrijf_toevoegen(request):
    """Formulier om een bedrijf toe te voegen."""
    if request.user.is_staff:
        if request.method == 'POST':
            form = BedrijfForm(request.POST)
            if form.is_valid():
                form.save()
                return redirect('bedrijf_lijst') # Zorg ervoor dat je een 'bedrijf_lijst' URL hebt
        else:
            form = BedrijfForm()
        return render(request, 'stages/bedrijf_toevoegen.html', {'form': form})
    else:
        return redirect('bedrijf_lijst')

# Mentor views

from django.shortcuts import render, redirect
from django.contrib.auth.decorators import login_required
from .models import Mentor, Stage, Feedback

@login_required
def mentor_dashboard(request):
    """Dashboard voor de mentor."""
    if hasattr(request.user, 'mentor'):
        mentor = request.user.mentor
        studenten = mentor.studenten.all()

        stages = Stage.objects.filter(benadering__student__in=studenten).distinct()

        feedback_lijst = Feedback.objects.filter(student__in=studenten)

        print(f"Studenten: {studenten}")
        print(f"Stages: {stages}")
        print(f"Feedback: {feedback_lijst}")

        context = {
            'mentor': mentor,
            'studenten': studenten,
            'stages': stages,
            'feedback_lijst': feedback_lijst,
        }
        return render(request, 'stages/mentor_dashboard.html', context)
    else:
        return redirect('home')


# Extra views

@login_required
def feedback_geven(request, stage_id):
    """Formulier om feedback te geven over een stage."""
    stage = get_object_or_404(Stage, pk=stage_id)
    student = Student.objects.get(user=request.user)
    if request.method == 'POST':
        beoordeling = request.POST['beoordeling']
        opmerkingen = request.POST['opmerkingen']
        Feedback.objects.create(student=student, stage=stage, beoordeling=beoordeling, opmerkingen=opmerkingen)
        return redirect('stage_detail', stage_id=stage_id)
    return render(request, 'stages/feedback_geven.html', {'stage': stage})
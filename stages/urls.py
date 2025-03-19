from django.urls import path
from django.contrib.auth import views as auth_views
from . import views

urlpatterns = [
    # Startpagina
    path('', views.home, name='home'),

    # Student URLs
    path('profiel/', views.student_profiel, name='student_profiel'),
    path('stages/', views.stage_lijst, name='stage_lijst'),
    path('stages/<int:stage_id>/', views.stage_detail, name='stage_detail'),
    path('bedrijven/<int:bedrijf_id>/', views.bedrijf_detail, name='bedrijf_detail'),
    path('stages/<int:stage_id>/benaderen/', views.stage_benaderen, name='stage_benaderen'),
    path('stages/<int:stage_id>/feedback/', views.feedback_geven, name='feedback_geven'),
    path('stages/toevoegen/', views.stage_toevoegen, name='stage_toevoegen'),
    path('bedrijven/toevoegen/', views.bedrijf_toevoegen, name='bedrijf_toevoegen'),
    path('profiel/bewerken/', views.student_profiel_bewerken, name='student_profiel_bewerken'),
    path('bedrijven/', views.bedrijf_lijst, name='bedrijf_lijst'), 

    # Coördinator URLs
    path('coordinator/', views.coordinator_dashboard, name='coordinator_dashboard'),

    # Mentor URLs
    path('mentor/', views.mentor_dashboard, name='mentor_dashboard'),

    # Rest URLs
    path('accounts/login/', auth_views.LoginView.as_view(), name='login'),
    path('accounts/logout/', auth_views.LogoutView.as_view(), name='logout'),
    path('accounts/register/', views.register, name='register'),
]
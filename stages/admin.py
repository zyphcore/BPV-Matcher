from django.contrib import admin
from .models import Student, Bedrijf, Stage, Benadering, Coördinator, Mentor, Feedback

# Registreer het Student model en pas de admin interface aan
@admin.register(Student)
class StudentAdmin(admin.ModelAdmin):
    # Velden die worden weergegeven in de lijstweergave
    list_display = ('user', 'programmeertalen', 'werkvoorkeuren')
    # Zoekvelden voor het zoeken naar studenten
    search_fields = ('user__username', 'programmeertalen')
    # Filters voor het filteren van studenten
    list_filter = ('programmeertalen',)

# Registreer het Bedrijf model en pas de admin interface aan
@admin.register(Bedrijf)
class BedrijfAdmin(admin.ModelAdmin):
    list_display = ('naam', 'locatie', 'is_samenwerkingsbedrijf', 'is_potentieel_bedrijf')
    search_fields = ('naam', 'locatie')
    list_filter = ('is_samenwerkingsbedrijf', 'is_potentieel_bedrijf')

# Registreer het Stage model en pas de admin interface aan
@admin.register(Stage)
class StageAdmin(admin.ModelAdmin):
    list_display = ('titel', 'bedrijf', 'startdatum', 'einddatum')
    search_fields = ('titel', 'bedrijf__naam')
    list_filter = ('bedrijf',)

# Registreer het Benadering model
admin.site.register(Benadering)

# Registreer het Coördinator model
admin.site.register(Coördinator)

# Registreer het Mentor model
admin.site.register(Mentor)

# Registreer het Feedback model
admin.site.register(Feedback)
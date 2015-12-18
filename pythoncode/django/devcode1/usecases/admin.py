from django.contrib import admin

# Register your models here.
from usecases.models import Person

admin.site.register(Person)

from django.db import models

class Person(models.Model):
    person_name = models.CharField(max_length=200)
    person_sId = models.CharField(max_length=200, default='nil')
    person_id = models.IntegerField(primary_key=True, default=0)

# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('usecases', '0003_person_person_sid'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='person',
            name='id',
        ),
        migrations.AlterField(
            model_name='person',
            name='person_id',
            field=models.IntegerField(default=0, serialize=False, primary_key=True),
            preserve_default=True,
        ),
    ]

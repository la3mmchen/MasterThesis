# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import models, migrations


class Migration(migrations.Migration):

    dependencies = [
        ('usecases', '0002_auto_20150421_1048'),
    ]

    operations = [
        migrations.AddField(
            model_name='person',
            name='person_sId',
            field=models.CharField(default=b'nil', max_length=200),
            preserve_default=True,
        ),
    ]

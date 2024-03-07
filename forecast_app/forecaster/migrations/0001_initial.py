# Generated by Django 4.1 on 2024-03-07 15:37

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Forecasts',
            fields=[
                ('id', models.BigAutoField(primary_key=True, serialize=False)),
                ('question', models.TextField(blank=True, null=True)),
                ('short_question', models.TextField(blank=True, null=True)),
                ('category', models.TextField(blank=True, null=True)),
                ('creation_date', models.DateTimeField(blank=True, null=True)),
                ('resolution_criteria', models.TextField(blank=True, null=True)),
            ],
            options={
                'db_table': 'forecasts',
            },
        ),
        migrations.CreateModel(
            name='Resolutions',
            fields=[
                ('id', models.AutoField(primary_key=True, serialize=False)),
                ('resolution', models.TextField(blank=True, null=True)),
                ('resolution_date', models.DateTimeField(blank=True, null=True)),
                ('brier_score', models.FloatField(blank=True, null=True)),
                ('logn_score', models.FloatField(blank=True, null=True)),
                ('log2_score', models.FloatField(blank=True, null=True)),
                ('forecast', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='resolutions', to='forecaster.forecasts')),
            ],
            options={
                'db_table': 'resolutions',
            },
        ),
        migrations.CreateModel(
            name='ForecastPoints',
            fields=[
                ('update_id', models.BigAutoField(primary_key=True, serialize=False)),
                ('point_forecast', models.FloatField(blank=True)),
                ('upper_ci', models.FloatField(blank=True, null=True)),
                ('lower_ci', models.FloatField(blank=True, null=True)),
                ('reason', models.TextField(blank=True, null=True)),
                ('date_added', models.DateTimeField(blank=True, null=True)),
                ('forecast', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='forecast_points', to='forecaster.forecasts')),
            ],
            options={
                'db_table': 'forecast_points',
            },
        ),
    ]

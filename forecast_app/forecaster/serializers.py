from rest_framework import serializers
from .models import Forecasts, ForecastPoints, Resolutions

class forecast_points_serializer(serializers.ModelSerializer):
    class Meta:
        model = ForecastPoints
        fields = '__all__'

class resolutions_serializer(serializers.ModelSerializer):
    class Meta:
        model = Resolutions
        fields = '__all__'

class forecasts_serializer(serializers.ModelSerializer):
    class Meta:
        model = Forecasts
        fields = '__all__'
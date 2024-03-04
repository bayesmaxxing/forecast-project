from django.shortcuts import render
from django.http import HttpResponse
from rest_framework import viewsets
from .models import Forecasts, ForecastPoints, Resolutions
from .serializers import forecasts_serializer, forecast_points_serializer, resolutions_serializer

class ForecastsViewSet(viewsets.ModelViewSet):
    serializer_class = forecasts_serializer

    def get_queryset(self):
        queryset = Forecasts.objects.all()

        category = self.request.query_params.get('category', None)
        if category is not None:
            queryset = queryset.filter(category__icontains=category)
        
        resolved = self.request.query_params.get('resolved', None)
        if resolved is not None:
            if resolved.lower() == 'true':
                queryset = queryset.filter(resolutions__forecast_id__isnull=False).distinct()
            elif resolved.lower() == 'false':
                queryset = queryset.filter(resolutions__forecast_id__isnull=True).distinct()
        return queryset

class ForecastPointsViewSet(viewsets.ModelViewSet):
    serializer_class = forecast_points_serializer
    
    def get_queryset(self):
        queryset = ForecastPoints.objects.all()
        forecast = self.request.query_params.get('forecast', None)
        if forecast is not None:
            queryset = queryset.filter(forecast=forecast)
        return queryset

class ResolutionsViewSet(viewsets.ModelViewSet):
    serializer_class = resolutions_serializer
    
    def get_queryset(self):
        queryset = Resolutions.objects.all()
        forecast = self.request.query_params.get('forecast', None)
        if forecast is not None:
            queryset = queryset.filter(forecast=forecast)
        return queryset


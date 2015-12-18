from django.conf.urls import patterns, url

from usecases import views

urlpatterns = patterns('',
    url(r'^test1/(?P<persons_id>\d+)/$', views.requestHandler, name='test1'),
    url(r'^test2/$', views.test2, name='test2'),
    url(r'^test3/$', views.test3, name='test3'),
    url(r'^test4/(?P<persons_id>\d+)/$', views.requestHandler, name='test4'),
    url(r'^index$', views.index, name='index'),
)

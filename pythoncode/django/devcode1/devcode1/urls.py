from django.conf.urls import patterns, include, url
from django.contrib import admin

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'devcode1.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^tests/', include('usecases.urls')),
    url(r'^admin/', include(admin.site.urls)),
)

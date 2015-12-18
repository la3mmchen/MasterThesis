from django.http import HttpResponse
from usecases.models import Person
import json
import random
import time
import sys

'''
  * @desc Testfall 1, Testfall 2, Testfall 3, Testfall 4
  * @needs request.REQUEST["object"]
  * valides Object {"Name": "Hans-Dieter", "id": "42"}
'''
# Index Methode fuer Zugriff auf alle angelegten Objekte
def index(request):
    latest_question_list = Person.objects.order_by('-person_id')

    output = ', '.join([p.person_name for p in latest_question_list])

    return HttpResponse(output, content_type='application/json')

# post /tests/test2/
def test2(request):
    print >>sys.stdout, '*** HTTP.Method=', request.method
    print >>sys.stdout, '*** Parameters=', request.REQUEST["object"]

    # Build object
    data = { }
    data["Person_sId"] = int(json.loads(request.REQUEST["object"]).values()[1])
    data["Person_name"] = json.loads(request.REQUEST["object"]).values()[0]
    if random.randint(1,2) == 1:
        data["Person_id"] = int(time.time()) + random.randint(1, 99999)
    else:
        data["Person_id"] = int(time.time()) - random.randint(10000, 99999)

    HttpResponse.status_code = 201

    return HttpResponse(json.dumps(data), content_type = 'application/json')

# post /tests/test3/
def test3(request):
    print >>sys.stdout, '*** HTTP.Method=', request.method
    print >>sys.stdout, '*** Parameters=', request.REQUEST["object"]

    # Build object
    data = { }
    data["Person_sId"] = int(json.loads(request.REQUEST["object"]).values()[1])
    data["Person_name"] = json.loads(request.REQUEST["object"]).values()[0]
    if random.randint(1,2) == 1:
        data["Person_id"] = int(time.time()) + random.randint(1, 99999)
    else:
        data["Person_id"] = int(time.time()) - random.randint(10000, 99999)

    # Some DB_Call
    p = Person.objects.get(person_id=1)

    HttpResponse.status_code = 201

    return HttpResponse(json.dumps(data), content_type = 'application/json')

# get /tests/test1/
# post /tests/test4/
def requestHandler(request, persons_id):
    # /test1/<id>/
    if request.method == 'GET':
        print >>sys.stdout, '*** HTTP.Method=', request.method
        HttpResponse.status_code = 200

        #p = Person.objects.get(person_id=persons_id)
        data = { }
        data["Person_id"] = "4711"
        data["Person_name"] = "abcd"
        data["Person_sId"] = "1a1a1a"

        return HttpResponse(json.dumps(data), content_type = 'application/json')
    # /test4/<id>/
    elif request.method == 'POST':
        print >>sys.stdout, '*** HTTP.Method=', request.method
        print >>sys.stdout, '*** Parameters=', request.REQUEST["object"]

        # Build object
        data = { }
        data["Person_sId"] = int(json.loads(request.REQUEST["object"]).values()[1])
        data["Person_name"] = json.loads(request.REQUEST["object"]).values()[0]
        if random.randint(1,2) == 1:
            data["Person_id"] = int(time.time()) + random.randint(1, 99999)
        else:
            data["Person_id"] = int(time.time()) - random.randint(10000, 99999)

        # Write to db
        p = Person.objects.create(person_id=data["Person_id"],person_name=data["Person_name"])
        HttpResponse.status_code = 201

        return HttpResponse(json.dumps(data), content_type = 'application/json')
    # Aussortieren der uebrigen HTTP.verbs
    elif request.method == 'PUT':
        HttpResponse.status_code = 405
        return HttpResponse("")
    elif request.method == 'DELETE':
        HttpResponse.status_code = 405
        return HttpResponse("")
    else:
        HttpResponse.status_code = 405
        return HttpResponse("")

#!/usr/bin/python3

# This script is used by a CI to clean up the existing
# unnamed draft releases

import sys
import requests

if (len(sys.argv) != 3):
    print("usage: "+sys.argv[0]+" GitHub_slug GitHub_OAuth_token")
    print("       "+sys.argv[0]+" digabi/naksu 13e74e84727837072a36f84402d8b2005c35185a")
    exit(1)

param_github_slug = sys.argv[1]
param_oauth_token = sys.argv[2]

auth = {'Authorization': 'token '+param_oauth_token}

r = requests.get('https://api.github.com/repos/'+param_github_slug+'/releases', headers=auth)
if (r.status_code != 200):
    print("Getting releases data failed, HTTP status code: "+str(r.status_code))
    exit(1)

releases = r.json()

if (len(releases) == 0):
    print("There are no draft releases in the given GitHub repo")

for this_release in releases:
    if (this_release['draft']):
        print("Deleting draft release: "+str(this_release['id']))
        del_req = requests.delete('https://api.github.com/repos/'+param_github_slug+'/releases/'+str(this_release['id']), headers=auth)
        if (del_req.status_code == 204):
            print("Deleted")
        else:
            print("Failed, HTTP status code: "+str(del_req.status_code))
    else:
        print("Skipping release #%d (%s)" % (this_release['id'], this_release['name']))

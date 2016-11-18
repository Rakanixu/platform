Path to the JSON file (google-cloud-service-account.json) must be set 
on GOOGLE_APPLICATION_CREDENTIALS environment variable.

GOOGLE_APPLICATION_CREDENTIALS=/home/pabloaguirre/Documents/google-cloud-service-account.json

Must be set properly on docker images:
GOOGLE_APPLICATION_CREDENTIALS=ROOT../platform/cmd/kazoup/google-app-credentials/google-cloud-service-account.json

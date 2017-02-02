Lifecycle configuration for Google cloud storage.

Bucket tmp-token stores JWT tokens for users. 
When we asks a third party to authorize a user, we lose the JWT on the redirect.

To avoid this issue, we save the JWT to be retrieve after redirection happens.
We want to delete those objects after 1 hour they have been created.

https://cloud.google.com/storage/docs/lifecycle#configuration

gsutil lifecycle set [LIFECYCLE_CONFIG_FILE] gs://[BUCKET_NAME]

Current rule defined in tmp_tokens_lifecycle.json deletes
JWT tokens after 1 day they've been created.

gsutil lifecycle set tmp_tokens_lifecycle.json gs://tmp-token

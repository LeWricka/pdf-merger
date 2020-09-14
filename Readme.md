**Run function**
$ gclogo run cmd/main.go 

**Upload function**
$ gcloud functions deploy Merge \  
--runtime go113 --trigger-http --allow-unauthenticated

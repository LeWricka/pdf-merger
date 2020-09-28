**PDF files merger**
PDF files merger using Google Cloud Functions (GCF) and Google Cloud Storage(GCS).

**Steps:**
- Create your own GCF
- Upload this code to your function defining BUCKET and PDFS_PATH:
  - $ gcloud functions deploy Merge --se--set-env-vars BUCKET=your_bucket_name,PDFS_PATH=path_in_gcs  \  
    --runtime go113 --trigger-http --allow-unauthenticated
- Call the function with the given endpoint and pass json array:
    -- ["draft1.pdf", "draft3.pdf"]
- If no error is thrown, a pdf with the given files merged on the given order will be returned :)

**Run function**
$ go run cmd/main.go 

 


 # CORE FUNCTIONS: S3


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.s3.buckets.delete](#delete)
 > * [aws.s3.buckets.describe](#describe)
 > * [aws.s3.buckets.create](#create)
 > * [aws.s3.buckets.website](#website)
 > * [aws.s3.buckets.putACL](#putACL)

 > * [aws.s3.objects.delete](#Odelete)
 > * [aws.s3.objects.create](#Ocreate)
 > * [aws.s3.objects.describe](#Odescribe)
 > * [aws.s3.objects.get](#Oget)
 > * [aws.s3.objects.read](#Oread)
 > * [aws.s3.objects.writeInto](#OwriteInto)

 This API allows resource handlers to manipulate S3 buckets and objects.

 ## AWS.S3.BUCKETS.DESCRIBE
 <a name="describe"></a>
 `aws.s3.buckets.describe(region, bucket-name);`

 Get info about the supplied bucket.

 Example:

 ```

 var bucket = aws.s3.buckets.describe("us-east-1", "mithras.io.");

 ```

 ## AWS.S3.BUCKETS.CREATE
 <a name="create"></a>
 `aws.s3.buckets.create(region, config);`

 Create a bucket.

 Example:

 ```

 var bucket = aws.s3.buckets.create("us-east-1",
 {
                 Bucket: "my-bucket"
                 ACL:    "public-read"
       		       LocationConstraint: "us-east-1"
 });

 ```

 ## AWS.S3.BUCKETS.DELETE
 <a name="delete"></a>
 `aws.s3.buckets.delete(region, bucket-name);`

 Delete a bucket

 Example:

 ```

 aws.s3.buckets.delete("us-east-1", "my-bucket");

 ```

 ## AWS.S3.BUCKETS.WEBSITE
 <a name="website"></a>
 `aws.s3.buckets.website(region, config);

 Set up a bucket to serve a static website.

 Example:

 ```

 aws.s3.buckets.website("us-east-1",
         website: {
             Bucket: bucketName
             WebsiteConfiguration: {
                 ErrorDocument: {
                     Key: "error.html"
                 }
                 IndexDocument: {
                     Suffix: "index.html"
                 }
              }
          });

 ```

 ## AWS.S3.BUCKETS.PUTACL
 <a name="putACL"></a>
 `aws.s3.buckets.putACL(region, config);

 Set up bucket access control config.

 Example:

 ```

 aws.s3.buckets.putACL("us-east-1",
 {
 	   	Bucket: bucketName
 	   	ACL:    "BucketCannedACL"
 	   	AccessControlPolicy: {
 	               Grants: [
 	   		{
 	                       Grantee: {
 	   			Type:         "Type"
 	   			DisplayName:  "DisplayName"
 	   			EmailAddress: "EmailAddress"
 	   			ID:           "ID"
 	   			URI:          "URI"
 	                       }
 	                       Permission: "Permission"
 	   		}
 	               ]
 	               Owner: {
 	   		DisplayName: "DisplayName"
 	   		ID:          "ID"
 	               }
 	   	}
 	   	GrantFullControl: "GrantFullControl"
 	   	GrantRead:        "GrantRead"
 	   	GrantReadACP:     "GrantReadACP"
 	   	GrantWrite:       "GrantWrite"
 	   	GrantWriteACP:    "GrantWriteACP"
 	});

 ```

 ## AWS.S3.OBJECTS.DELETE
 <a name="Odelete"></a>
 `aws.s3.objects.delete(region, bucket, key);`

 Delete an object in a bucket.

 Example:

 ```

 aws.s3.objects.delete("us-east-1", "my-bucket", "index.html");

 ```

 ## AWS.S3.OBJECTS.CREATE
 <a name="Ocreate"></a>
 `aws.s3.objects.create(region, config);`

 Create an object

 Example:

 ```

 aws.s3.objects.create("us-east-1",
 {
                         Bucket:             "my-bucket"
                         Key:                "index.html"
                         ACL:                "public-read"
                         Body:               "contents"
                         ContentType:        type
 });

 ```

 ## AWS.S3.OBJECTS.DESCRIBE
 <a name="Odescribe"></a>
 `aws.s3.objects.describe(region, bucket, prefix);`

 Get object info.

 Example:

 ```

 aws.s3.objects.create("us-east-1", "my-bucket", "index.html");

 ```

 ## AWS.S3.OBJECTS.GET
 <a name="Oget"></a>
 `aws.s3.objects.get(region, bucket, key);`

 Get an object in a bucket.

 Example:

 ```

 aws.s3.objects.get("us-east-1", "my-bucket", "index.html");

 ```

 ## AWS.S3.OBJECTS.READ
 <a name="Oread"></a>
 `aws.s3.objects.read(region, bucket, key);`

 Get object content.

 Example:

 ```

 var bytes = aws.s3.objects.read("us-east-1", "my-bucket", "index.html");

 ```

 ## AWS.S3.OBJECTS.WRITEINTO
 <a name="OwriteInto"></a>
 `aws.s3.objects.writeInto(region, bucket, key, path, mode);`

 Get object content and write it into a file at `path`, with permissions `mode`.

 Example:

 ```

 aws.s3.objects.writeInto("us-east-1", "my-bucket", "index.html", "/tmp/foo", 0644);

 ```



 


 # CORE FUNCTIONS: S3


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.s3.buckets.delete](#delete)
 > * [aws.s3.buckets.get](#get)
 > * [aws.s3.buckets.describe](#describe)
 > * [aws.s3.buckets.create](#create)
 > * [aws.s3.buckets.website](#website)

 > * [aws.s3.objects.delete](#Odelete)
 > * [aws.s3.objects.create](#Ocreate)
 > * [aws.s3.objects.describe](#Odescribe)

 This API allows resource handlers to manipulate DNS records in Route53.

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

 ## AWS.S3.BUCKETS.GET
 <a name="get"></a>
 `aws.s3.buckets.get(region, bucket, key);`

 Get an object in a bucket.

 Example:

 ```

 aws.s3.buckets.get("us-east-1", "my-bucket", "index.html");

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


 


 # CORE FUNCTIONS: TAG


 

 This package exports several entry points into the JS environment,
 including:

 > * [aws.tags.create](#create)

 This API allows resource handlers to manipulate tags on AWS resources.

 ## AWS.TAGS.CREATE
 <a name="create"></a>
 `aws.tags.create(region, id, tags);`

 Tag an AWS resource.

 Example:

 ```

 tags.create("us-east-1", "vpc-abc", { Name: "foo" });

 ```



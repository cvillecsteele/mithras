 
 
 # s3
 
 S3 is resource handler for working with AWS's S3.
 
 This module exports:
 
 > * `init` Initialization function, registers itself as a resource
 >   handler with `mithras.modules.handlers` for resources with a
 >   module value of `"s3"`
 
 Usage:
 
 `var s3 = require("s3").init();`
 
  ## Example Bucket Resource
 
 ```javascript
 var bucket = {
     name: "s3bucket"
     module: "s3"
     params: {
         ensure: ensure
         region: defaultRegion
         bucket: {
             Bucket: bucketName
             ACL:    "public-read"
             LocationConstraint: defaultRegion
         }
         acls: [ 								 
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
 	       }								 
 	   ]                                                                     
         notification: {
                  Bucket: bucketName                                  
                  NotificationConfiguration: {                          
                      QueueConfigurations: [                            
                          {                                             
                              Events: [                                 
                                  "Event"                               
                              ]                                         
                              QueueArn: "QueueArn"                      
                              Filter: {                                 
                                  Key: {                                
                                      FilterRules: [                    
                                          {                             
                                              Name:  "FilterRuleName"   
                                              Value: "FilterRuleValue"  
                                          }                             
                                      ]                                 
                                  }                                     
                              }                                         
                              Id: "NotificationId"                      
                          }                                             
                      ]                                                 
                  }                                                     
              }                                                         
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
          } // website
      } // params
 };
 ```
 
  ## Example Object Resource
 
 ```javascript
 var thing = {
     name: "thingInS3"
     module: "s3"
     dependsOn: [bucket.name]
     params: {
       ensure: "latest"
       region: defaultRegion
       stat: fs.stat(path)
       object: {
            Bucket:             bucketName
            Key:                "some/thing.html" 
            ACL:                "public-read"
            Body:               "..."
            ContentType:        "text/html"
       }
     }
 };
 ```
 
 ## Parameter Properties
 
 ### `ensure`

 * Required: true
 * Allowed Values: "present", "absent" or "latest" (for objects)

 If `"present"`, the bucket/object will be created if it doesn't
 already exist.  If `"absent"`, the bucket/object will be removed if
 it is present.  If `"latest"`, the resource specifies a `stat`
 property with the results of `fs.stat` on the source path, and if
 the object is S3 is older than the one on local disk, the object in
 S3 is updated.
 
 ### `stat`

 * Required: false
 * Allowed Values: results of `fs.stat()` on the source file

 If operating on an object in S3, and the object is S3 has a
 modification time before the `ModTime` of the stat'd file, the
 object in S3 will be updated.

 ### `object`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutObjectInput)

 Specifies the parameters for the object.

 ### `bucket`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-CreateBucketInput)

 Specifies the parameters for the bucket.

 ### `website`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutBucketWebsiteInput)

 ### `notification`

 * Required: false
 * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutBucketNotificationConfigurationInput)

 Configure the bucket to send notification events.



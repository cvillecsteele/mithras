// MITHRAS: Javascript configuration management tool for AWS.
// Copyright (C) 2016, Colin Steele
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//   the Free Software Foundation, either version 3 of the License, or
//                  (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//              GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// @public
// 
// # s3
// 
// S3 is resource handler for working with AWS's S3.
// 
// This module exports:
// 
// > * `init` Initialization function, registers itself as a resource
// >   handler with `mithras.modules.handlers` for resources with a
// >   module value of `"s3"`
// 
// Usage:
// 
// `var s3 = require("s3").init();`
// 
//  ## Example Bucket Resource
// 
// ```javascript
// var bucket = {
//     name: "s3bucket"
//     module: "s3"
//     params: {
//         ensure: ensure
//         region: defaultRegion
//         bucket: {
//             Bucket: bucketName
//             ACL:    "public-read"
//             LocationConstraint: defaultRegion
//         }
//         acls: [ 								 
// 	       {								 
// 	   	Bucket: bucketName
// 	   	ACL:    "BucketCannedACL"
// 	   	AccessControlPolicy: {
// 	               Grants: [
// 	   		{
// 	                       Grantee: {
// 	   			Type:         "Type"
// 	   			DisplayName:  "DisplayName"
// 	   			EmailAddress: "EmailAddress"
// 	   			ID:           "ID"			 
// 	   			URI:          "URI"		 
// 	                       }						 
// 	                       Permission: "Permission"
// 	   		}							 
// 	               ]
// 	               Owner: {					 
// 	   		DisplayName: "DisplayName"
// 	   		ID:          "ID"
// 	               }
// 	   	}
// 	   	GrantFullControl: "GrantFullControl"
// 	   	GrantRead:        "GrantRead"
// 	   	GrantReadACP:     "GrantReadACP"
// 	   	GrantWrite:       "GrantWrite"
// 	   	GrantWriteACP:    "GrantWriteACP"
// 	       }								 
// 	   ]                                                                     
//         notification: {
//                  Bucket: bucketName                                  
//                  NotificationConfiguration: {                          
//                      QueueConfigurations: [                            
//                          {                                             
//                              Events: [                                 
//                                  "Event"                               
//                              ]                                         
//                              QueueArn: "QueueArn"                      
//                              Filter: {                                 
//                                  Key: {                                
//                                      FilterRules: [                    
//                                          {                             
//                                              Name:  "FilterRuleName"   
//                                              Value: "FilterRuleValue"  
//                                          }                             
//                                      ]                                 
//                                  }                                     
//                              }                                         
//                              Id: "NotificationId"                      
//                          }                                             
//                      ]                                                 
//                  }                                                     
//              }                                                         
//         website: {
//             Bucket: bucketName
//             WebsiteConfiguration: {
//                 ErrorDocument: {
//                     Key: "error.html"
//                 }
//                 IndexDocument: {
//                     Suffix: "index.html"
//                 }
//              }
//          } // website
//      } // params
// };
// ```
// 
//  ## Example Object Resource
// 
// ```javascript
// var thing = {
//     name: "thingInS3"
//     module: "s3"
//     dependsOn: [bucket.name]
//     params: {
//       ensure: "latest"
//       region: defaultRegion
//       stat: fs.stat(path)
//       object: {
//            Bucket:             bucketName
//            Key:                "some/thing.html" 
//            ACL:                "public-read"
//            Body:               "..."
//            ContentType:        "text/html"
//       }
//     }
// };
// ```
// 
// ## Parameter Properties
// 
// ### `ensure`
//
// * Required: true
// * Allowed Values: "present", "absent" or "latest" (for objects)
//
// If `"present"`, the bucket/object will be created if it doesn't
// already exist.  If `"absent"`, the bucket/object will be removed if
// it is present.  If `"latest"`, the resource specifies a `stat`
// property with the results of `fs.stat` on the source path, and if
// the object is S3 is older than the one on local disk, the object in
// S3 is updated.
// 
// ### `stat`
//
// * Required: false
// * Allowed Values: results of `fs.stat()` on the source file
//
// If operating on an object in S3, and the object is S3 has a
// modification time before the `ModTime` of the stat'd file, the
// object in S3 will be updated.
//
// ### `object`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutObjectInput)
//
// Specifies the parameters for the object.
//
// ### `bucket`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-CreateBucketInput)
//
// Specifies the parameters for the bucket.
//
// ### `website`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutBucketWebsiteInput)
//
// ### `notification`
//
// * Required: false
// * Allowed Values: JSON corresponding to the structure found [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3.html#type-PutBucketNotificationConfigurationInput)
//
// Configure the bucket to send notification events.
//
(function (root, factory){
    if (typeof module === 'object' && typeof module.exports === 'object') {
        module.exports = factory();
    }
})(this, function() {
    
    var sprintf = require("sprintf.js").sprintf;
    var moment = require('moment');

    var handler = {
        moduleNames: ["s3"]
        preflight: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            return [null, true];
        }
        handleBucket: function(catalog, resource) {
            if (!resource.params.bucket) {
                return;
            }
            var buckets = aws.s3.buckets.describe(resource.params.region, "*");
            var bucket = _.findWhere(buckets,
                                     {"Name": resource.params.bucket.Bucket});
            if (bucket) {
                if (resource.params.ensure === 'absent') {
                    if (mithras.verbose) {
                        log(sprintf("Deleting bucket '%s'",
                                    resource.params.bucket.Bucket));
                    }
                    aws.s3.buckets.delete(resource.params.region,
					  resource.params.bucket.Bucket);
                }
            } else {
                if (resource.params.ensure === 'present') {
                    if (mithras.verbose) {
                        log(sprintf("Creating bucket '%s'",
                                    resource.params.bucket.Bucket));
                    }
                    var res = aws.s3.buckets.create(resource.params.region,
						    resource.params.bucket);
                    if (resource.params.website) {
                        if (mithras.verbose) {
                            log(sprintf("Adding website config to bucket '%s'",
                                        resource.params.bucket.Bucket));
                        }
                        aws.s3.buckets.website(resource.params.region,
                                               resource.params.website);
                    } 
                    if (resource.params.notification) {
                        if (mithras.verbose) {
                            log(sprintf("Adding notification config to bucket '%s'",
                                        resource.params.bucket.Bucket));
                        }
                        aws.s3.buckets.notification(resource.params.region,
                                                    resource.params.notification);
                    } 
                   if (resource.params.acls) {
			_.each(resource.params.acls, function(acl) {
                            if (mithras.verbose) {
				log(sprintf("Adding ACL to bucket '%s'",
                                            resource.params.bucket.Bucket));
                            }
                            aws.s3.buckets.putACL(resource.params.region, acl);
			});
                    }
                }
            }
        }
        runObject: function (params) {
            var sprintf = require("sprintf.js").sprintf;
        
            var objects = aws.s3.objects.describe(params.region,
                                                  params.object.Bucket,
                                                  params.object.Key);

            var obj = _.findWhere(objects,
                                  {"Key": params.object.Key});
            if (obj) {
                if (params.ensure === 'absent') {
                    if (mithras.verbose) {
                        log(sprintf("Deleting object '%s'",
                                    params.object.Key));
                    }
                    aws.s3.objects.delete(params.region,
					  params.object.Bucket,
                                          params.object.Key);
                } else if (params.ensure === "latest") {
                    var m1 = moment(obj.LastModified.String(),
                                    "YYYY-MM-DD HH:mm:ss Z UTC",
                                    true);
                    var m2 = moment(params.stat.ModTime);
                    if (m2.isAfter(m1)) {
                        if (mithras.verbose) {
                            log(sprintf("Updating object '%s'",
                                        params.object.Key));
                        }
                        var res = aws.s3.objects.create(params.region,
							params.object);
                    } else {
                        if (mithras.verbose) {
                            log(sprintf("No action taken.", params.object.Key));
                        }
		    }
                }
            } else {
                if ((params.ensure === 'present') ||
                    (params.ensure === 'latest')) {
                    if (mithras.verbose) {
                        log(sprintf("Creating object '%s'",
                                    params.object.Key));
                    }
                    var res = aws.s3.objects.create(params.region,
						    params.object);
                }
            }
        }
        handleObject: function(catalog, resource) {
            if (!resource.params.object) {
                return;
            }
            var params = resource.params;
            if (Array.isArray(params.hosts)) {
                var js = sprintf("var run = function() {\n (%s)(%s); };\n",
                                 handler.runObject.toString(),
                                 JSON.stringify(_.omit(params, 'hosts')));
                for (var i in params.hosts) {
                    var instance = params.hosts[i];
                    var result = mithras.remote.mithras(instance,
                                                        mithras.sshUserForInstance(resource, instance),
                                                        mithras.sshKeyPathForInstance(resource, instance),
                                                        js,
                                                        params.become,
                                                        params.becomeUser,
                                                        params.becomeMethod);
                    if (result[3] == 0) {
                        log(sprintf("S3 object '%s' %s",
                                    params.object.Key,
                                    params.ensure));
                    }
                }
            } else {
                handler.runObject(params);
            }
        }
        handle: function(catalog, resources, resource) {
            if (!_.find(handler.moduleNames, function(m) { 
                return resource.module === m; 
            })) {
                return [null, false];
            }
            var updated = mithras.updateResource(resource,
                                                 catalog,
                                                 resources,
                                                 resource.name);
            var updatedParams = updated.params;
	    
            if (updatedParams.skip == true) {
                log("Skipped.");
            } else {            
		handler.handleBucket(catalog, updated);
		handler.handleObject(catalog, updated);
	    }
            return [null, true];
        }
    };
    
    handler.init = function () {
        _.each(handler.moduleNames, function(name) {
            mithras.modules.handlers.register(name, handler.handle);
        });
    };

    handler.contentTypeMap = {
        "3g2": "video/3gpp2",
        "3gp": "video/3gpp",
        "3gp2": "video/3gpp2",
        "3gpp": "video/3gpp",
        "aa": "audio/audible",
        "aac": "audio/vnd.dlna.adts",
        "aax": "audio/vnd.audible.aax",
        "addin": "text/xml",
        "adt": "audio/vnd.dlna.adts",
        "adts": "audio/vnd.dlna.adts",
        "ai": "application/postscript",
        "aif": "audio/aiff",
        "aifc": "audio/aiff",
        "aiff": "audio/aiff",
        "application": "application/x-ms-application",
        "asax": "application/xml",
        "ascx": "application/xml",
        "asf": "video/x-ms-asf",
        "ashx": "application/xml",
        "asmx": "application/xml",
        "aspx": "application/xml",
        "asx": "video/x-ms-asf",
        "au": "audio/basic",
        "avi": "video/avi",
        "bmp": "image/bmp",
        "btapp": "application/x-bittorrent-app",
        "btinstall": "application/x-bittorrent-appinst",
        "btkey": "application/x-bittorrent-key",
        "btsearch": "application/x-bittorrentsearchdescription+xml",
        "btskin": "application/x-bittorrent-skin",
        "cat": "application/vnd.ms-pki.seccat",
        "cd": "text/plain",
        "cer": "application/x-x509-ca-cert",
        "config": "application/xml",
        "contact": "text/x-ms-contact",
        "crl": "application/pkix-crl",
        "crt": "application/x-x509-ca-cert",
        "cs": "text/plain",
        "csproj": "text/plain",
        "css": "text/css",
        "csv": "application/vnd.ms-excel",
        "datasource": "application/xml",
        "der": "application/x-x509-ca-cert",
        "dib": "image/bmp",
        "dll": "application/x-msdownload",
        "doc": "application/msword",
        "docm": "application/vnd.ms-word.document.macroEnabled.12",
        "docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "dot": "application/msword",
        "dotm": "application/vnd.ms-word.template.macroEnabled.12",
        "dotx": "application/vnd.openxmlformats-officedocument.wordprocessingml.template",
        "dtd": "application/xml-dtd",
        "dtsconfig": "text/xml",
        "eps": "application/postscript",
        "exe": "application/x-msdownload",
        "fdf": "application/vnd.fdf",
        "fif": "application/fractals",
        "gif": "image/gif",
        "group": "text/x-ms-group",
        "hdd": "application/x-virtualbox-hdd",
        "hqx": "application/mac-binhex40",
        "hta": "application/hta",
        "htc": "text/x-component",
        "htm": "text/html",
        "html": "text/html",
        "hxa": "application/xml",
        "hxc": "application/xml",
        "hxd": "application/octet-stream",
        "hxe": "application/xml",
        "hxf": "application/xml",
        "hxh": "application/octet-stream",
        "hxi": "application/octet-stream",
        "hxk": "application/xml",
        "hxq": "application/octet-stream",
        "hxr": "application/octet-stream",
        "hxs": "application/octet-stream",
        "hxt": "application/xml",
        "hxv": "application/xml",
        "hxw": "application/octet-stream",
        "ico": "image/x-icon",
        "ics": "text/calendar",
        "ipa": "application/x-itunes-ipa",
        "ipg": "application/x-itunes-ipg",
        "ipsw": "application/x-itunes-ipsw",
        "iqy": "text/x-ms-iqy",
        "iss": "text/plain",
        "ite": "application/x-itunes-ite",
        "itlp": "application/x-itunes-itlp",
        "itls": "application/x-itunes-itls",
        "itms": "application/x-itunes-itms",
        "itpc": "application/x-itunes-itpc",
        "jfif": "image/jpeg",
        "jnlp": "application/x-java-jnlp-file",
        "jpe": "image/jpeg",
        "jpeg": "image/jpeg",
        "jpg": "image/jpeg",
        "js": "application/javascript",
        "latex": "application/x-latex",
        "library-ms": "application/windows-library+xml",
        "m1v": "video/mpeg",
        "m2t": "video/vnd.dlna.mpeg-tts",
        "m2ts": "video/vnd.dlna.mpeg-tts",
        "m2v": "video/mpeg",
        "m3u": "audio/mpegurl",
        "m3u8": "audio/x-mpegurl",
        "m4a": "audio/m4a",
        "m4b": "audio/m4b",
        "m4p": "audio/m4p",
        "m4r": "audio/x-m4r",
        "m4v": "video/x-m4v",
        "magnet": "application/x-magnet",
        "man": "application/x-troff-man",
        "master": "application/xml",
        "mht": "message/rfc822",
        "mhtml": "message/rfc822",
        "mid": "audio/mid",
        "midi": "audio/mid",
        "mod": "video/mpeg",
        "mov": "video/quicktime",
        "mp2": "audio/mpeg",
        "mp2v": "video/mpeg",
        "mp3": "audio/mpeg",
        "mp4": "video/mp4",
        "mp4v": "video/mp4",
        "mpa": "video/mpeg",
        "mpe": "video/mpeg",
        "mpeg": "video/mpeg",
        "mpf": "application/vnd.ms-mediapackage",
        "mpg": "video/mpeg",
        "mpv2": "video/mpeg",
        "mts": "video/vnd.dlna.mpeg-tts",
        "odc": "text/x-ms-odc",
        "odg": "application/vnd.oasis.opendocument.graphics",
        "odm": "application/vnd.oasis.opendocument.text-master",
        "odp": "application/vnd.oasis.opendocument.presentation",
        "ods": "application/vnd.oasis.opendocument.spreadsheet",
        "odt": "application/vnd.oasis.opendocument.text",
        "otg": "application/vnd.oasis.opendocument.graphics-template",
        "oth": "application/vnd.oasis.opendocument.text-web",
        "ots": "application/vnd.oasis.opendocument.spreadsheet-template",
        "ott": "application/vnd.oasis.opendocument.text-template",
        "ova": "application/x-virtualbox-ova",
        "ovf": "application/x-virtualbox-ovf",
        "oxt": "application/vnd.openofficeorg.extension",
        "p10": "application/pkcs10",
        "p12": "application/x-pkcs12",
        "p7b": "application/x-pkcs7-certificates",
        "p7c": "application/pkcs7-mime",
        "p7m": "application/pkcs7-mime",
        "p7r": "application/x-pkcs7-certreqresp",
        "p7s": "application/pkcs7-signature",
        "pcast": "application/x-podcast",
        "pdf": "application/pdf",
        "pdfxml": "application/vnd.adobe.pdfxml",
        "pdx": "application/vnd.adobe.pdx",
        "pfx": "application/x-pkcs12",
        "pko": "application/vnd.ms-pki.pko",
        "pls": "audio/scpls",
        "png": "image/png",
        "pot": "application/vnd.ms-powerpoint",
        "potm": "application/vnd.ms-powerpoint.template.macroEnabled.12",
        "potx": "application/vnd.openxmlformats-officedocument.presentationml.template",
        "ppa": "application/vnd.ms-powerpoint",
        "ppam": "application/vnd.ms-powerpoint.addin.macroEnabled.12",
        "pps": "application/vnd.ms-powerpoint",
        "ppsm": "application/vnd.ms-powerpoint.slideshow.macroEnabled.12",
        "ppsx": "application/vnd.openxmlformats-officedocument.presentationml.slideshow",
        "ppt": "application/vnd.ms-powerpoint",
        "pptm": "application/vnd.ms-powerpoint.presentation.macroEnabled.12",
        "pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
        "prf": "application/pics-rules",
        "ps": "application/postscript",
        "psc1": "application/PowerShell",
        "pwz": "application/vnd.ms-powerpoint",
        "py": "text/plain",
        "pyw": "text/plain",
        "rat": "application/rat-file",
        "rc": "text/plain",
        "rc2": "text/plain",
        "rct": "text/plain",
        "rdlc": "application/xml",
        "resx": "application/xml",
        "rmi": "audio/mid",
        "rmp": "application/vnd.rn-rn_music_package",
        "rqy": "text/x-ms-rqy",
        "rtf": "application/msword",
        "sct": "text/scriptlet",
        "settings": "application/xml",
        "shtml": "text/html",
        "sit": "application/x-stuffit",
        "sitemap": "application/xml",
        "skin": "application/xml",
        "sldm": "application/vnd.ms-powerpoint.slide.macroEnabled.12",
        "sldx": "application/vnd.openxmlformats-officedocument.presentationml.slide",
        "slk": "application/vnd.ms-excel",
        "sln": "text/plain",
        "slupkg-ms": "application/x-ms-license",
        "snd": "audio/basic",
        "snippet": "application/xml",
        "spc": "application/x-pkcs7-certificates",
        "sst": "application/vnd.ms-pki.certstore",
        "stc": "application/vnd.sun.xml.calc.template",
        "std": "application/vnd.sun.xml.draw.template",
        "stl": "application/vnd.ms-pki.stl",
        "stw": "application/vnd.sun.xml.writer.template",
        "svg": "image/svg+xml",
        "sxc": "application/vnd.sun.xml.calc",
        "sxd": "application/vnd.sun.xml.draw",
        "sxg": "application/vnd.sun.xml.writer.global",
        "sxw": "application/vnd.sun.xml.writer",
        "tga": "image/targa",
        "thmx": "application/vnd.ms-officetheme",
        "tif": "image/tiff",
        "tiff": "image/tiff",
        "torrent": "application/x-bittorrent",
        "ts": "video/vnd.dlna.mpeg-tts",
        "tts": "video/vnd.dlna.mpeg-tts",
        "txt": "text/plain",
        "user": "text/plain",
        "vb": "text/plain",
        "vbox": "application/x-virtualbox-vbox",
        "vbox-extpack": "application/x-virtualbox-vbox-extpack",
        "vbproj": "text/plain",
        "vcf": "text/x-vcard",
        "vdi": "application/x-virtualbox-vdi",
        "vdp": "text/plain",
        "vdproj": "text/plain",
        "vhd": "application/x-virtualbox-vhd",
        "vmdk": "application/x-virtualbox-vmdk",
        "vor": "application/vnd.stardivision.writer",
        "vscontent": "application/xml",
        "vsi": "application/ms-vsi",
        "vspolicy": "application/xml",
        "vspolicydef": "application/xml",
        "vspscc": "text/plain",
        "vsscc": "text/plain",
        "vssettings": "text/xml",
        "vssscc": "text/plain",
        "vstemplate": "text/xml",
        "vsto": "application/x-ms-vsto",
        "wal": "interface/x-winamp3-skin",
        "wav": "audio/wav",
        "wave": "audio/wav",
        "wax": "audio/x-ms-wax",
        "wbk": "application/msword",
        "wdp": "image/vnd.ms-photo",
        "website": "application/x-mswebsite",
        "wiz": "application/msword",
        "wlz": "interface/x-winamp-lang",
        "wm": "video/x-ms-wm",
        "wma": "audio/x-ms-wma",
        "wmd": "application/x-ms-wmd",
        "wmv": "video/x-ms-wmv",
        "wmx": "video/x-ms-wmx",
        "wmz": "application/x-ms-wmz",
        "wpl": "application/vnd.ms-wpl",
        "wsc": "text/scriptlet",
        "wsdl": "application/xml",
        "wsz": "interface/x-winamp-skin",
        "wvx": "video/x-ms-wvx",
        "xaml": "application/xaml+xml",
        "xbap": "application/x-ms-xbap",
        "xdp": "application/vnd.adobe.xdp+xml",
        "xdr": "application/xml",
        "xfdf": "application/vnd.adobe.xfdf",
        "xht": "application/xhtml+xml",
        "xhtml": "application/xhtml+xml",
        "xla": "application/vnd.ms-excel",
        "xlam": "application/vnd.ms-excel.addin.macroEnabled.12",
        "xld": "application/vnd.ms-excel",
        "xlk": "application/vnd.ms-excel",
        "xll": "application/vnd.ms-excel",
        "xlm": "application/vnd.ms-excel",
        "xls": "application/vnd.ms-excel",
        "xlsb": "application/vnd.ms-excel.sheet.binary.macroEnabled.12",
        "xlsm": "application/vnd.ms-excel.sheet.macroEnabled.12",
        "xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        "xlt": "application/vnd.ms-excel",
        "xltm": "application/vnd.ms-excel.template.macroEnabled.12",
        "xltx": "application/vnd.openxmlformats-officedocument.spreadsheetml.template",
        "xlw": "application/vnd.ms-excel",
        "xml": "text/xml",
        "xrm-ms": "text/xml",
        "xsc": "application/xml",
        "xsd": "application/xml",
        "xsl": "text/xml",
        "xslt": "application/xml",
        "xss": "application/xml",

	"svg": "image/svg+xml"
        "json": "application/json",
    };

    return handler;
});

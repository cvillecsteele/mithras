 


 # CORE FUNCTIONS: LOG


 

 This package exports several entry points into the JS environment,
 including:

 > * [log](#log)
 > * [log0](#log0)
 > * [log1](#log1)
 > * [log2](#log2)
 > * [debug](#debug)
 > * [debug](#info)
 > * [warning](#warning)
 > * [error](#error)
 > * [fatal](#fatal)
 > * [panic](#panic)
 > * [withFields](#withFields)
 > * [setLevel](#setLevel)
 > * [jsonFormatter](#jsonFormatter)

 This API allows resource handlers to log.

 ## LOG
 <a name="log"></a>
 `log(message);`

 Logs a message in the default format

 Example:

 ```

  log("hello");

 ```

 ## LOG
 <a name="log"></a>
 `log(message);`

 Logs a message in the default format

 Example:

 ```

  log("hello");

 ```

 ## LOG0
 <a name="log0"></a>
 `log0(message);`

 Logs a more important message.

 Example:

 ```

  log0("hello");

 ```

 ## LOG1
 <a name="log1"></a>
 `log1(message);`

 Logs a less important message

 Example:

 ```

  log1("hello");

 ```

 ## LOG2
 <a name="log1"></a>
 `log1(message);`

 Logs an even less important message

 Example:

 ```

  log2("hello");

 ```


 
 ## setLevel
 <a name="setLevel"></a>
 `setLevel(logLevel);`

 Only log messages at this severity or above.  Values may be:

 > * "debug"
 > * "info"
 > * "warning"
 > * "error"
 > * "fatal"
 > * "panic"


 
 ## jsonFormatter
 <a name="jsonFormatter"></a>
 `log.jsonFormatter(format);`

 The `format` argument may be omitted, but if set, should be a string as documented [here](https://github.com/Sirupsen/logrus)


 
 ## debug
 <a name="debug"></a>
 `log.debug(message);`

 Log a message at debug level.


 
 ## info
 <a name="info"></a>
 `log.info(message);`

 Log a message at info level.


 
 ## warn
 <a name="warn"></a>
 `log.warn(message);`

 Log a message at warning level.


 
 ## error
 <a name="error"></a>
 `log.error(message);`

 Log a message at error level.


 
 ## fatal
 <a name="fatal"></a>
 `log.fatal(message);`

 Log a message at fatal level.


 
 ## panic
 <a name="panic"></a>
 `log.panic(message);`

 Log a message at panic level.


 
 ## withFields
 <a name="withFields"></a>
 `log.withFields(fieldObj);`

 Create a set of fields to be logged with a message.  Returns an
 object with `.debug()`, `.info()`, etc, function properties.

 For example:

    log.withFields({omg: "Yeah!"}).info("this is awesome");



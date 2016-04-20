 

 # CORE FUNCTIONS: CHANNELS


 
 This package exports several entry points into the JS environment,
 including:

 > * [chan.make](#make)
 > * [chan.snd](#snd)
 > * [chan.rcv](#rcv)
 > * [chan.range](#range)
 > * [chan.select](#select)

 This API allows callers to create a Go channel.

 ## CHAN.MAKE
 <a name="make"></a>
 `chan.make(name);`

 Creates a channel identified by `name`.

 Example:

 ```

  chan.make("test_channel");

 ```

 ## CHAN.SND
 <a name="snd"></a>
 `chan.snd(name, value);`

 Send `value` to channel identified by `name`.

 Example:

 ```

  chan.snd("test_channel", "hello, world!");

 ```

 ## CHAN.RCV
 <a name="rcv"></a>
 `chan.rcv(name);`

 Read a value from the channel specified by `name`.

 Example:

 ```

 var fromChannel = chan.rcv("test_channel");

 ```

 ## CHAN.RANGE
 <a name="range"></a>
 `chan.range(name, callback);`

 Read succeessive values from the channel `name`, calling `callback`
 with each received value.  Callbacks cease to be invoked when the
 channel is closed.

 Example:

 ```

  chan.range({"test_channel": function(value) {
    console.log("Got value " + value);
  }});

 ```

 ## CHAN.SELECT
 <a name="select"></a>
 `chan.select(selectMap);`

 Perform a select operation on multiple channels, as per Go's `select`.

 Example:

 ```

  chan.select({"test_channel": function(channelName, value, ok) {
    console.log("Got value " + value + " on channel " + channelName);
  }});

 ```



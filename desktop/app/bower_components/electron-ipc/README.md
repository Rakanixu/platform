# electron-ipc

An element used to communicate through the `ipc` module for electron. It can
send or receive data from a single or multiple events at once.

By default the element will send data to the specified events whenever it's
required to do so, or it automatically send the data if `auto` is set.

```HTML
<electron-ipc id="target" events="myEvent" data="myData"></electron-ipc>
<script>
    document.querySelector('#target').send();
</script>

<electron-ipc events="myEvent" data="myData" auto></electron-ipc>
```

If you need to listen to events instead of sending them, just transform the
element into a `receiver` and it will keep `data` up to date whenever it
receives new information.

```HTML
<electron-ipc id="target" receiver></electron-ipc>

<script>
    var el = document.querySelector('#target');

    el.events = ['firstEvent', 'secondEvent'];

    /**
     * From the main process:
     * ipc.send('firstEvent', 'Hello!');
     */

    console.log(el.data); // 'Hello!'

    /**
     * From the main process:
     * ipc.fire('secondEvent', 'World!');
     */

    console.log(el.data); // 'World!'
</script>
```

## Properties

#### auto
Default: `false`

If set to `true` a the element will automatically send a message when `data` is changed.

---
#### data
The element's data, either received from events or set by the user.

---
#### events
The events associated to the element, either for listening or sending messages.

---
#### receiver
Default: `false`

If set to `true` the element becomes a `receiver` and is able to receive messages.

---
## Events
#### ipc-data-received

Fired when data is received through ipc.
 
**Detail:**

 1. **firer {Element}** The element itself
 2. **data {*}** The data received
 3. **event {string}** The name of the event that generated the event

---
#### ipc-data-received
Fired when data is sent through ipc.
 
**Detail:**

 1. **firer {Element}** The element itself.
 2. **data {*}** The data sent.
 3. **event {string}** The name of the event that generated the event.

---
## Methods
#### send(data)
Sends an ipc message using data. If no data is passed then the `data` set on the element will be sent.
 
**Parameters:**

 1. **[data] {*}** Optional data to be sent.

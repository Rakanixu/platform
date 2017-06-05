# Documentation for create a new service

##  Create a new folder underneath platform:
```
platform/example
```

## Create the service folder structure

```
Service entry point:
platform/example/srv/main.go

Protobuffer:
platform/example/srv/proto/example/example.proto

Handler (synchronous):
platform/example/srv/handler/handler.go
platform/example/srv/handler/handler_test.go

Subscribers (asynchronous):
platform/example/subscriber/announce.go
platform/example/subscriber/announce_test.go
platform/example/subscriber/subscriber.go
platform/example/subscriber/subscriber_test.go
```


## Create service deployment file for continous integration, gitignore, dockerfiles, make files, etc:
   example can be composed by srv and web. We will omit web for simplicity.
   
General to the whole service, srv and web (if exists):
```
platform/example/.gitignore 
platform/example/circle.yml         (Deploy on continuos integration server)
platform/example/Makefile           (Use to get all go dependencies, build binary and generate protobuff source code from a proto file)
platform/example/README.md
```

Specific to srv:
```
platform/example/srv/Dockerfile     (Dockerfile to build the service)
platform/example/srv/Makefile
platform/example/srv/README.md
```

## Define your service:
   We can start by definining the proto file (platform/example/srv/proto/example/example.proto):
   
   We are going to define the example service and a handler that manages the creation of an abstract entity of data:
   
```
syntax = "proto3";

package proto.example;

service Service {
 rpc Create(CreateRequest) returns (ECreateResponse) {}
}

message CreateRequest {
 string entity_id = 1;
 string entity_data_1 = 2;
 bool   entity_data_2 = 3;
 int64  entity_data_3 = 4;
}

message CreateResponse {

}
```

   Once proto is defined, let's build the source code:
```
$ platform/example make protoc
```

A new file should be created on platform/example/srv/proto/example/example.pb.go

## Define the handler (platform/example/srv/handler/handler.go):
```
package handler

// Have a look on github.com/kazoup/platform/lib/quota, this is a interface that will work because on main.go we will import the plugins, 
// where the implementation resides, those plugins will autoregister themselfs.
import (
	"errors"
	"github.com/kazoup/platform/example/srv/proto/example"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"golang.org/x/net/context"
)

const (
	QUOTA_EXCEEDED_MSG = "Quota for example exceeded."
)

type Service struct {

}

func (s *Service) Create(ctx context.Context, req *proto_example.CreateRequest, rsp *proto_example.CreateResponse) error {
    if err := validate.Exists(req.EntityId, req.EntityData_1, req.EntityData_2,...); err != nil {
		return err
	}
	
	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	_, _, rate, _, quota, ok := quota.Check(ctx, globals.example_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.example_SERVICE_NAME, "Create", "", 403, errors.New("quota.Check"))
	}

	// Quota exceded, respond sync and do not initiate go routines
	if rate-quota > 0 {
		rsp.Info = QUOTA_EXCEEDED_MSG
		return nil
	}

	return nil
}
```
Notice that your protobuffer previously generated is an interface with the methods we defined previously (Create). 
Our Service struct implement this interface, so we always need to import the proto source code.
As you can see this handler does not do much, just validates data and check quota. Quota can be omitted depending what we want to achieve with the service.
This is because we do not have to respond anything to the client, we will respond 200 or an error if something went wrong, for example, client sent invalid data.
At this point the creation /  save of the data did not happen yet, but don't worry, it will happen on the subscribers.


### Create announce subscriber (platform/example/subscriber/announce.go)

Here we will subscribe to the messages we are intereted on. 
We want to process an announce message every time on Create handler is call. The only thing we have to do is create the subscriber, the message will be broadcasted automatically by a service wrapper.


```
package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/example/srv/proto/example"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	create "github.com/kazoup/platform/lib/protomsg/create"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type AnnounceHandler struct{}

// OnCreate reacts to Create handler calls
func (a *AnnounceHandler) OnCreate(ctx context.Context, msg *announce.AnnounceMessage) error {
	// We are just interested on those messages related to the Create handler
	// Notice that HANDLER_example_CREATE probably does not exist, have a look on other consants to follow the same naming pattern.
	if globals.HANDLER_example_CREATE == msg.Handler {
	    // We'll need the service to access to the client. In this example, srv is the example instance
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

        // msg.Data can be unmarssal to the handler request type; which is  proto_example.Create
        // This is deeply related with the automatic annonce system.
		var cr *proto_example.Create
		if err := json.Unmarshal([]byte(msg.Data), &cr); err != nil {
			return err
		}

        // In this example we want to publish a new message; CreateAbstractEntity. We will require a new topic for it
        // We will need to include the new topic and define the new message:
        // globals.CreateAbstractEntityTopic
        // "github.com/kazoup/platform/lib/protomsg/create"
        if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.CreateAbstractEntityTopic, &create.CreateAbstractEntityMessage{
            EntityID:       cr.EntityID,
            EntityData1:    cr.EntityData1,
            EntityData2:    cr.EntityData2,
            EntityData3:    cr.EntityData3,
        })); err != nil {
            log.Println("ERROR publishing CreateAbstractEntityTopic", err)
        }
	}

	return nil
}
```

As you can see, this subscriber is feels like just a proxy that receives messages over the announce topic and publish a specific message, with the same data into a speceific topic.
Now, think that you just want to handle the Creation of abstract entities at night. To achieve this new requierement we can add logic here that check what time is it and 
therefore publish the CreateAbstractEntityMessage if is at night.

The design of this system allows us to think in an action / reaction chain system, for example, a crawler has finish, and afterwards we want to extract the content from written documents.
We only have to subscribe to the crawler finish announce messages, therefore we can publish a new message EnrichDocument. The subscriber would hold the logic to publish such a message when
the file is a written document only. (LINK ARQUITECTURE DIAGRAM)


### Create subscriber (platform/example/subscriber/subscriber.go)
   
On the previous step we defined the annnounce subscriber, and we decided that we are going to publish CreateAbstractEntityMessage over CreateAbstractEntityTopic topic
    
```
type ExampleServiceTaskHandler struct {}

func (t *ExampleServiceTaskHandler) LogCreateAbstractEntityMessage(ctx context.Context, createMsg *create.CreateAbstractEntityMessage) error {
    // So far we just want to log the messages.
    // A normal user case could be to bulk index these AbstractEntities.

    log.Println(createMsg.EntityID)
    log.Println(createMsg.EntityData1)
    log.Println(createMsg.EntityData2)
    log.Println(createMsg.EntityData3)

	return nil
}
```


## Create the entry point (main.go):
   Usually is composed by the main function
   
```
package main

import (

	"github.com/kazoup/platform/example/srv/handler"
	"github.com/kazoup/platform/example/srv/proto/example"
	"github.com/kazoup/platform/example/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
)

func main() {
	var m monitor.Monitor

    // Instantiate the service. NewKazoupService will create a new service with default wrappers 
    // Those wrappers adds behavior like automatic logging, quota checks, authentication and automatic action broadcasting
    // Wrappers can be understand as onion layers, eg)
    
    // extend context (adds the service to the request context) ->
    //     authentication ->
    //          service handler (the handler in itself, for this example will be new(handler.Service).Create) ->
    //          <- automatic publishing (announce will be pusblish)
    //     <- quota
    // <- logging
    
	service := wrappers.NewKazoupService("example", m)

	// example monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Attach handler
	// In this example only a Create handler will be attached
	proto_audio.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Attach subscriber
	// On this subscriber we will log the data we are receiving
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.CreateAbstractEntityTopic,
			new(subscriber.ExampleServiceTaskHandler),
			server.SubscriberQueue("example"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	// This subscriber will pick up all annunce messages that are publish automatically by a wrapper
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-example"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Every service uses the same wrappers for handlers and subscribers in terms of logic. Explicit types are the only differences between HandlerWrappers and SubscriberWrappers.


## The flow
Now, let's go deep inside the flow of a request, and what is being fired underneath the this example service.

Imagine we deploy this service, our first steep will be to call ExampleService.Create handler from our frontend, Curl, Postman or whatever.

-  Define the request:
https://web.kazoup.io:8082/rpc
```
-HEADERS
Content-Type:application/json
Authorization:VALID_JSON_WEB_TOKEN
-BODY
{
    "service":"com.kazoup.srv.example",
    "method":"Service.Create",
    "request":{
        "entity_id":"entity_id",
        "entity_data_1":"entity_data_1",
        "entity_data_2":true,
        "entity_data_3":100
    }
}
```

- Context wrapper is fired, a copy of the service is added to the context request.

- Auth wrapper is fired, Authorization header is extracted from the context, and JWT is validated. Bail out if token is not valid.

- Handler is fired. Our handler.Create method is executed.

- AnnounceTopic is published into the broker. in this ocasion, msg.Handler is equal to the service plus handler name

   msg.Data is a marshal copy of our request data: 
   ```
   {
       "entity_id":"entity_id",
       "entity_data_1":"entity_data_1",
       "entity_data_2":true,
       "entity_data_3":100
   }
   ```
   Two flows will occur in paralel now, A & B.
   
- Quota is fired. Counter for action and user will be increased if  required.

- Loging is fired. A message like this will be seen:

```
    example-srv_1     | time="2017-06-05T12:27:38Z" level=info msg=OK handler=Service.Create service=com.kazoup.srv.example 
```
    
- Client receives HTTP response, in this case 200 Ok with empty body.

- Announce subscriber will be fired; AnnounceHandler.OnCreate.

- Subscriber will be fired; ExampleServiceTaskHandler.LogCreateAbstractEntityMessage

- AnnounceTopic is published into the broker. in this ocasion, msg.Handler is equal to the topic; CreateAbstractEntityTopic
    msg.Data is a marshal copy of our msg data: 
	```
    &create.CreateAbstractEntityMessage{
                EntityID:       "entity_id",
                EntityData1:    "entity_data_1",
                EntityData2:    true,
                EntityData3:    100,
    }
	```
                                                  
- New announcer can de defined to do something when ExampleServiceTaskHandler.LogCreateAbstractEntityMessage has been succesful
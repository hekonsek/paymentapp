# Payment API

This is an application written in Go implementing simple payment API gateway.

## Running

The easiest way to run the application is to do it via Docker:

    docker run --net=host -it hekonsek/paymentapp

To verify that your app is up and running, you can list payments via CURL:

    $ curl http://localhost:8080/payments
    {"data":[]}

## Building

You don't really have to build the application as it is already available on DockerHub 
([hekonsek/paymentapp](https://hub.docker.com/r/hekonsek/paymentapp)), but if you really want to, please follow
 this guide ;) . 

In order to build the application, execute the following command:

    make build

Build process includes the following steps:
- Go code formatting 
- executing unit and REST API tests 
- generating Pact contract file
- validating Pact contract (consumer, provider)
- building actual application binary

## Persistence

Application uses pluggable persistence provider architecture, where provided persistent stores are 
[AWS Document DB](https://aws.amazon.com/documentdb) using Mongo 3.x client compatibility and (for testing purposes
and to demonstrate persistence pluggability) in-memory persistence provider. By default application starts with
in-memory provider.

You can specify persistence provider using command line:

    docker run --net=host -it hekonsek/paymentapp --persistence=awsdocdb
    docker run --net=host -it hekonsek/paymentapp --persistence=mem
    
All options provided by commandline can be overridden via environment variables to make Docker/Kubernetes configuration
easy and idiomatic. For example following command...

    docker run --net=host -it hekonsek/paymentapp --persistence=awsdocdb
 
...can be replaced with...

    docker run -e PERSISTENCE=awsdocdb --net=host -it hekonsek/paymentapp

    
Also see command help for all configuration options and more details:

    docker run --net=host -it hekonsek/paymentapp --help

As our application is targeted to be executed in Kubernetes environment, we rely on Kubernetes env-based service
discovery to detect AWS DocumentDB connection details (`AWSDOCDB_SERVICE_HOST` and `AWSDOCDB_SERVICE_PORT` environment variables)
where `awsdocdb` is expected to be registered as external service in Kubernetes. If `AWSDOCDB_SERVICE_*` environment 
variables are not present, application falls back to `localhost:27017` connection (as we assume that you want to 
test application locally against MongoDB emulating AWS DocumentDB). You can also override connection settings via
those environment variables:

    docker run -e AWSDOCDB_SERVICE_HOST=some-mongo-or-awsdocdb-address.example.com --net=host -it hekonsek/paymentapp --persistence=awsdocdb
    
## Design notes

Application is a simple microservice written in Go v1.11.2.

Project relies on Go Modules for dependency management. Modules become de-facto standard for this purpose so I decided
not to use vet, dep and other older dependency management system.

Project is a command line tool based on [Cobra](https://github.com/spf13/cobra) which brings nice cli experience 
including clear help for commands and neat support for configuration flags.

### API

Application exposes single HTTP endpoint serving RESTful API. I used [Gin](https://github.com/gin-gonic/gin) as web
framework because it is one of the most popular web frameworks for Go and because I like how nicely it adds extra
abstraction layer on the top of standard go `http` and `json` packages.

The API has been designed to follow REST principles in:
- URL design convention around payments resources
- using appropriate HTTP methods for CRUD operations (POST => CREATE, PUT => UPDATE, etc.)
- using appropriate HTTP status code to indicate results of operations

#### Consumer-Driven Contracts

Project uses [Pact](https://github.com/pact-foundation/pact-go) to provide sample contracts for the API. Due to the limited
time I provided only two interaction contracts, but those can give you idea how Pact can be used to verify that current API
fulfills contract expected by the consumer. In real-life application we should create extensive set of contracts tests
to ensure our API is not broken from consumer point of view.

#### OpenAPI

I decided to skip implementation of this one, but for real world application I would expose REST API via OpenAPI 3. Right
now there is no OpenAPI 3 library for Go which is good enough IMHO, so I would use 
[Go Rice](https://github.com/GeertJohan/go.rice) and Go Templates and just maintain OpenAPI specification manually. It would
be nice to have it integrated in [gin-swagger](https://github.com/swaggo/gin-swagger) -like fashion, but the latter
library is not stable enough from my experiences. Also maintaining OpenAPi specification manually is not that hard - been
there, done that ;) .

#### API versioning

In this application we have only single version of the API. In real world we should have introduce versioning of the API, as
it evolves over time. I'm personally fan of URL based versioning without minor versions like:

    http://myapi.com/v1/payments
    http://myapi.com/v2/payments
    
Some people says it is wrong, because you should detect API version from HTTP header, but I like URL driven style because
almost all the world uses this convention.

As we introduce versioning into our application, we should also move away from keeping JSON formatting annotations in the model 
([see](https://github.com/hekonsek/paymentapp/blob/master/payments/payments.go#L14)) and maintain service model without any 
JSON-awareness plus dedicated DTOs for each version of API which should contains JSON specific information.

#### SDK

There is no SDK for this project. After introducing OpenAPI into it we could auto-generate SDK using Swagger. However from
my experiences Swagger autogenerated APIS are not as developer-friendly as nice SDK handcrafted with love and maintained 
by API development team ;) .

### Gateway

As you can see there is not authentication for the service. It is because it is assumed to be deployed behind 
[API gateway](https://microservices.io/patterns/apigateway.html) which is responsible for enforcing security.

Gateway responsibility would be to make sure that API requests coming from public networks are holding valid
JWT OAuth token. Gateway would be also responsible for a token validation. After requests is authenticated, HTTP request
is forward to the proper microservice, for example our paymentapp.

Gateway should also handle some traffic logic like throttling / rate limiting, premium tenant handling, etc.

### Multi-tenancy

This is huge topic. Huuuge ;) For this simple application I assume that you're using shared-all multi tenancy architecture
typical for SaaS applications where data for all customers (tenants) are co-located together in the same database and 
uses the same applications for handling HTTP load.

In real world that would mean that gateway after successful token validation should set tenant id in the API request header
forwarded to the application. Then we should filter out database queries against certain field (for example `payment.organisation`).
Also REST operations creating data should set that field to the proepr tenant value.

Multi-tenancy implementation is highly dependent on the balance between your customer needs (like physical separation of
data for goverment customers, etc) and cost effectiveness.

### Observability

Application exposes standard metrics via Prometheus. We use the same Gin HTTP endpoint for serving REST API and for
serving metrics. We could add custom metrics in the code as well.

### Containers first

Application is assumed to be packaged as Docker container and deployed into Kubernetes. However nothing prevents it from
being deployed into AWS ECS as well. Just keep in mind that service discovery for AWS Document DB relies on Kubernetes-like
conventions (i.e. env-based service discovery `AWSDOCDB_SERVICE_HOST`, etc).

### Wanna know more?

We can talk about architecture and design all day long. Including how we could deploy it into AWS in the most maintainable
and cost-effective fashion. Wanna know more - let's talk :) .
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
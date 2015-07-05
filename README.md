# Comodoro [![Docker Pulls](https://img.shields.io/docker/pulls/mfellner/comodoro.svg)](https://registry.hub.docker.com/u/mfellner/comodoro) [![Travis](https://img.shields.io/travis/mfellner/comodoro.svg?style=flat-square)](https://travis-ci.org/mfellner/comodoro) [![Coveralls](https://img.shields.io/coveralls/mfellner/comodoro.svg?style=flat-square)](https://coveralls.io/r/mfellner/comodoro)

*Work in progress*

## Usage

    ./comodoro

* `-db="/tmp/comodoro.db"`: Path to the BoltDB file
* `-log="info"`: Log level (debug|info)
* `-port=3030`: Port to listen on

## Tests

**With Go test:**

    got test ./...

**With [GoConvey](http://goconvey.co/)**:

    $GOPATH/bin/goconvey

## Debug on CoreOS

You'll need [docker-1.6.2](https://docs.docker.com/installation/binaries/)
and [coreos-vagrant](https://github.com/coreos/coreos-vagrant/):

    git clone https://github.com/coreos/coreos-vagrant.git
    wget https://get.docker.com/builds/Darwin/x86_64/docker-1.6.2
    chmod +x docker-1.6.2

Enable port forwarding of the Docker TCP socket in the config.rb of coreos-vagrant
and also forward Comodoro's APP_PORT (e.g. 3030):

    $expose_docker_tcp=2375
    $forwarded_ports = { 3030 => 3030 }

Then configure Docker to use the CoreOS server:

    export DOCKER_HOST='tcp://127.0.0.1:2375'

After booting CoreOS using `vagrant up` should be able to build the
Docker container:

    ./docker-1.6.2 build -t comodoro .

Don't forget to map the fleet socket into the container:

docker run -d --name comodoro -p 3030:3030 \
  -e APP_LOGLEVEL=debug \
  -v /var/run/fleet.sock:/var/run/fleet.sock \
  comodoro

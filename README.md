# Docker Image Authorization Plugin with DCT

The Docker Engine image authorization plugin allows docker images only from a
predefined authorized Notary to be used by the docker engine. This implements
DCT on the side of DE. The implementation uses `docker pull` with enabled DCT to
check the presence of the image metadata in the Notary and its validity. For
additional information, please refer to [docker
documentation](https://docs.docker.com/engine/extend/) on plugins.


### Running tests for the image authorization plugin

_NOTE: You require *sudo privileges as root* on your machine to run the plugin tests._

#### Testing the plugin with docker engine v1.11
```
# Create the plugin tests image
docker build -t plugin-tests-1.11.2 -f Dockerfile.test --build-arg DOCKER_VERSION=1.11.2 --rm .

# Start the plugin tests container
docker run --privileged -d --name plugin-tests-1.11.2 plugin-tests-1.11.2

# Run the plugin tests
docker exec plugin-tests-1.11.2 python tests.py

# Remove the plugin tests container
docker rm -f plugin-tests-1.11.2

# Remove the plugin tests image
docker rmi -f plugin-tests-1.11.2
```

#### Testing the plugin with docker engine v1.12
```
# Create the plugin tests image
docker build -t plugin-tests-1.12.6 -f Dockerfile.test --build-arg DOCKER_VERSION=1.12.6 --rm .

# Start the plugin tests container
docker run --privileged -d --name plugin-tests-1.12.6 plugin-tests-1.12.6

# Run the plugin tests
docker exec plugin-tests-1.12.6 python tests.py

# Remove the plugin tests container
docker rm -f plugin-tests-1.12.6

# Remove the plugin tests image
docker rmi -f plugin-tests-1.12.6
```

### Build and install the plugin
```
# Create the build tools docker image
docker build -t plugin-build-tools:latest .

# Build the image authorization plugin from sources
docker run --rm -v `pwd`:`pwd` -w `pwd` -e GOPATH=`pwd` plugin-build-tools:latest make

# Generate the plugin service units configuration
# IMPORTANT: Please note that the make config command supports generation of systemd units for CentOS/RHEL only
docker run --rm -v `pwd`:`pwd` -w `pwd` -e GOPATH=`pwd` plugin-build-tools:latest \
  make config NOTARY=<https://my-notary[:port]>

# Add notary root CA certificate.
# The certificate should be placed to ~/.docker/tls of the user on behalf
# of whom the service will be running.
~/.docker/tls/<my-notary[:port]>/root-ca.crt

# Install the plugin service
docker run --rm -v `pwd`:`pwd` -w `pwd` -e GOPATH=`pwd` -v /usr/libexec:/usr/libexec \
  -v /usr/lib/systemd/system:/usr/lib/systemd/system plugin-build-tools:latest \
  make install

# Start the plugin service
systemctl daemon-reload
systemctl enable img-authz-plugin
systemctl start img-authz-plugin
```

### Enable the authorization plugin on docker engine
##### Step-1: Add authorization plugin to the docker engine configuration
Please add the following cmdline flag to your docker engine (e.g. ExecStart line /usr/lib/systemd/system/docker.service)
```
--authorization-plugin img-authz-plugin
```
##### Step-2: Restart docker engine
```
systemctl daemon-reload
systemctl restart docker
```

### Stop and uninstall the plugin
NOTE: Before doing below, remove the authorization-plugin configuration created above and restart the docker daemon.
```
# Stop the plugin service
systemctl stop img-authz-plugin
systemctl disable img-authz-plugin

# Uninstall the plugin service units
docker run --rm -v `pwd`:`pwd` -w `pwd` -e GOPATH=`pwd` -v /usr/libexec:/usr/libexec \
  -v /usr/lib/systemd/system:/usr/lib/systemd/system plugin-build-tools:latest \
  make uninstall

```

### To remove the generated artifacts
```
docker run --rm -v `pwd`:`pwd` -w `pwd` -e GOPATH=`pwd` plugin-build-tools:latest make clean
```

### Access plugin logs
```
journalctl -xe -u img-authz-plugin -f
```

### Contact
For further queries on the plugin, please reach out to me at cpdevws@gmail.com
or post an issue in the repo. Also, pull requests welcome for extending the
plugin for other linux distributions and useful features!

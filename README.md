# FlatFS
Flat File System With Multiple Attributes

## To run it with docker
### Build
docker build -t flatfs .
### Run
docker run -it --name FlatFS --privileged --cap-add SYS_ADMIN --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined flatfs
### Execute FlatFS
Attach another bash to your docker instance with (as you need two to use FlatFS)
docker exec -it FlatFS bash

With one instance run
cd usr/gopath/src/github.com/sarpk/FlatFS/
./main /tmp/mountpoint/

On the other bash instance now you can access to /tmp/mountpoint/ with
cd /tmp/mountpoint/

If you want to exit from the FlatFS please use
fusermount -u /tmp/mountpoint/
command instead of just sending SIGINT or SIGTERM (CTRL + C)

# FlatFS
Flat File System With Multiple Attributes

## To run it with docker
### Build
`docker build -t flatfs .`
### Run
`docker run -it --name FlatFS --privileged --cap-add SYS_ADMIN --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined flatfs`

If you get an error saying that FlatFS is already in use then remove FlatFS container first and then run the above command again:
`docker rm -f FlatFS`

### Execute FlatFS
Attach another bash to your docker instance with (as you need two to use FlatFS):

`docker exec -it FlatFS bash`

With one instance run:

`cd usr/gopath/src/github.com/sarpk/FlatFS/`

`./main /tmp/mountpoint/ /tmp/flatDir default`

instead of default, sqlite can also be used

On the other bash instance now you can access to /tmp/mountpoint/ with:

`cd /tmp/mountpoint/`

Now for a simple test, you can do:
`echo test >> 'foo=hello,bar=world' && ls -l ?foo=hello && tail -n 10 'foo=hello,bar=world'`

If you want to exit from the FlatFS please make sure you are not in `/tmp/mountpoint/` directory and use
`fusermount -u /tmp/mountpoint/`
command instead of just sending SIGINT or SIGTERM (CTRL + C)

# FlatFS
Flat File System With Multiple Attributes

## To run it with docker
### Build
docker build -t flatfs .
### Run
docker run -it --privileged --cap-add SYS_ADMIN --cap-add MKNOD --device=/dev/fuse --security-opt apparmor:unconfined flatfs

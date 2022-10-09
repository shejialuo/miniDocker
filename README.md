# Mini Docker

This is my code for writing a mini docker with [自己动手写Docker](https://book.douban.com/subject/27082348/).

## Environment Setup

The code uses the following environment, you should be consistent with that.
Because we are learning docker, this is the point. Don't waste your time.
However, I refuse to use the low version of go. Although there would be some
change, but it is easy to handle.

+ Ubuntu 14.04
+ Kernel: 3.13.0-83-generic
+ Go: 1.19.1

### Ubuntu IOS

Download from tuna:

```sh
wget https://mirrors.tuna.tsinghua.edu.cn/ubuntu-releases/trusty/ubuntu-14.04.6-server-amd64.iso
```

### Change The Kernel

```sh
sudo apt-get update
sudo apt-get install linux-image-3.13.0-83-generic
sudo apt-get install linux-image-extra-$(uname -r)
sudo modprobe aufs
```

Change the configuration of the `grub` to start the corresponding kernel.

### Download Go

```sh
wget https://golang.google.cn/dl/go1.19.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz
```

Next configure the environment.

```sh
export PATH=$PATH:/usr/local/go/bin
```

### Other things

Due to the upgrade of the Go, when using cgo, the linker is
a bit out-of-date. So we need use another tool:

```sh
sudo apt-get update
sudo apt-get install binutils-2.26
export PATH="/usr/lib/binutils-2.26/bin:$PATH"
```

---

At now, I have studied all the parts except the network.
Well, it is a really wonderful tutorial. What I cannot build,
I do not understand.

## Purpose of clone

Rename the repository, so that binary in $GOPATH/bin will be *go-get-shallow-git* instead of *git*. I wanted to add this workarround permanently using  ~/.bashrc hack without polluting gobal git command.


# How to use shallow clones for `go get`

This is a little hack to use shallow clones for new git checkouts with `go get`. Unfortunately for Gophers, [this has been an open issue for three years counting](https://github.com/golang/go/issues/13078) without a workable solution aside from patching the go toolchain yourself. This solution utilizes a *git* wrapper that determines if a pull/clone is happening and then makes sure it is shallow. 

To install it you do 

```
$ go get github.com/arunmj/go-get-shallow-git
```
add a symlink `$GOPATH/bin/go-get-shallow-git` to a location not covered in $PATH

replace <TARGETDIR> with convinient directory path
  
```
$ ln -s $GOPATH/bin/go-get-shallow-git <TARGETDIR>/git
```

The *git*-wrapper tool is symlinked to *"git"* in \<TARGETDIR\>, so that you can add an alias to *go* command with <TARGETDIR> prepended to the PATH and then the *git*-wrapper will be substituted for the real *git* (`/usr/bin/git`). 

```
$ alias go="PATH=<TARGETDIR>:\$PATH go" 
```
OR use *go_get* instead of *go get*

```
$ alias go_get="PATH=<TARGETDIR>:\$PATH go get" 
```

which you can add to your `.bashrc` files if you want it to be permanent. This way, the wrapper will aways be used and the wrapper will force cloning to be shallow.

# Benchmarks

Here's a benchmark showing a 50% reduction in disk usage and thus a 50% reduction in time taken for a `go get`. You'll not get that much for smaller repositories, but its not bad.

## normal `go get`


```
% docker run -it golang:1.10 /bin/bash
root@d9208178f1fa:/go# time go get github.com/juju/juju/...
real    7m35.631s
user    1m40.059s
sys     0m45.436s
root@d9208178f1fa:/go# du -sh .
1.1G
```

## shallow `go get`

```
% docker run -it golang:1.10 /bin/bash
root@68135fb64a3e:/go# go get github.com/schollz/git
root@68135fb64a3e:/go# export PATH=$GOPATH/bin:$PATH
root@68135fb64a3e:/go# time go get github.com/juju/juju/...
real    3m0.335s
user    0m29.192s
sys     0m17.253s
root@d9208178f1fa:/go# du -sh .
499M    .
```

# Acknowledgements

Thanks [tscholl2](https://github.com/tscholl2) for the idea.

Thanks [schollz](https://github.com/schollz/git) for the code

# License

MIT

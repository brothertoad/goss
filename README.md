# goss

A simple static site generate written in Go.

I wrote this (very simple) static site generator to create a few web sites.  I had previously used [metalsmith](https://metalsmith.io), which is
very powerful, but is a Node application, and after a while dealing with dependabot alerts got annoying.  Since I wrote goss for my own use, it is
not documented yet, but there is a [sample project](https://github.com/brothertoad/uulists) here on github that you can look at and/or clone if you
want to see how goss works.

Note that unlike most static site generators, goss is not aimed at blogs.  In fact, it has no blog-like functionality.

To install goss, just clone to a local directory, go into that directory, and type `go get` followed by `go install`.  You will need to have the Go
compiler installed.

To create a web site, go into the directory where the source files for your site are and simply type `goss`.  You can optionally supply a config file
(defaults to goss.yaml).

Note that uses the text/template module, rather than the html/template module, so no escaping of input is done.

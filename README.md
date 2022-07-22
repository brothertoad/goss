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

As of this writing (22 July 2022), goss supports both generic Golang templates (processed using text/template, rather than html/template), or jinja2
templates (using the
[gonja](https://github.com/noirbizarre/gonja)
module).  The
[jinja2 templates](https://jinja.palletsprojects.com/en/3.1.x/templates/)
are much more powerful and easier to use.

As of this writing, there are no releases, although I hope to do one soon.  (I did have some tags, but I deleted them, as I belatedly realized that
there were significant issues).

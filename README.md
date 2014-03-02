hypergraphdb
============

A generic Hypergraph Database

To install, make sure you have activated a Go workspace:

    echo $GOPATH

If not, create an empty directory somewhere and activate it:

    mkdir -p some/dir/anywhere
    cd some/dir/anywhere
    export GOPATH=`pwd`

Once you have a Go workspace, run:

    go get github.com/yet-another-project/hypergraphdb/cmd/graphdb

Every time you want to activate the workspace, export GOPATH as above. It's
helpful to also adjust your PATH:

    export PATH=$PATH:$GOPATH/bin

Then run the console:

    graphdb

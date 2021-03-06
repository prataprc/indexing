#!/bin/bash


build_indexer(){

    echo "Building Indexer..."
    cd cmd/indexer/
    go build -o indexer
    cp indexer ../../bin/
    cd ../..
    echo "Done"
    echo "Indexer binary under bin/"
}

clean_indexer(){

    cd cmd/indexer/
    go clean
    rm -f indexer
    cd ../..
    rm -f bin/indexer
}

build_projector(){

    echo "Building Projector..."
    cd cmd/projector/
    go build -o projector
    cp projector ../../bin/
    cd ../..
    echo "Done"
    echo "Projector binary under bin/"
}

clean_projector(){

    cd cmd/projector/
    go clean
    rm -f projector
    cd ../..
    rm -f bin/projector
}

build_protobuf(){

    echo "Building Protobuf..."
    cd protobuf
    make
    cd ..
    echo "Done"
}

if [ -z "$1" ]
    then
    build_protobuf
    build_indexer
    build_projector
elif [ $1 == "indexer" ]
    then
    build_indexer
elif [ $1 == "projector" ]
    then
    build_projector
elif [ $1 == "protobuf" ]
    then
    build_protobuf
elif [ $1 == "clean" ]
    then
    echo "Cleaning..."
    clean_indexer
    clean_projector
    echo "Done"
else
    echo "Unknown build option"
fi

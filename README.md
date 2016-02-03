About
=====

Example of Event Sourcing and CQRS in golang. Based on:

* https://github.com/pinballjs/event-sourcing-cqrs
* https://github.com/cer/event-sourcing-examples

Instructions
============

Required:

* go ver. 1.5.x
* gb https://github.com/constabulary/gb

```
    git clone git@github.com:mcveat/event-sourcing-cqrs.git
    cd event-sourcing-cqrs
    git submodule init
    git submodule update
    gb build all
    ./bin/cqrs
```

To test:

```
    gb test
```

# Demo p2p chat

# Architecture & dev principles
- __clean architecture__ (aka ports & adapters)
    - core data structures are in [domain](./backend/domain/)
    - usecases (business logic) in [uc](./backend/uc/)
    - driven ports (calling usecases) in [driven](./backend/driven/)
- __CQRS__ :
  - if a function modifies the state of the system, it returns nothing (errors don't count)
  - if a function returns something, it doesn't modify the system
- __Least privilege principle__ : a function has access only to what it needs to do its job
  - this one is not really easy in go : the idiomatic interfaces are too open (acces to all methods when less are needed) so I used closures instead with function signatures as types... the only way to group them logically is using `structs` instead of `interfaces` but the idea is the same.
- __BDD__ with [goConvey](http://goconvey.co/)
- a single repo and a shared codebase is used for both central server and clients, it's just for the convenience in order to speed up a little bit the dev of this project

## Global architecture
![architecture](./img/archi.png)
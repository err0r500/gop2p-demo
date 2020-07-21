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



### Example 
bob registers
```$xslt
curl -v -X POST localhost:3001/sessions/ -H 'Content-Type: application/json' -d '{"login": "bob", "password": "pass", "address": "client1:4000"}' -H 'user:bob'
```

alice registers
```$xslt
curl -v -X POST localhost:3002/sessions/ -H 'Content-Type: application/json' -d '{"login": "alice", "password": "pass", "address": "client2:4000"}' -H 'user:alice'
```

alice sends a message to bob
```$xslt
curl -v -X POST localhost:3002/messages/ -H 'Content-Type: application/json' -d '{"message":"salut bob, c est alice", "To": "bob"}' -H "user: alice"
```

bob replies
```$xslt
curl -v -X POST localhost:3001/messages/ -H 'Content-Type: application/json' -d '{"message":"salut alice !", "To": "alice"}' -H "user: bob"
```

alice checks messages from bob
```$xslt
curl -v localhost:3002/conversations/bob
```

bob checks messages from alice
```$xslt
curl -v localhost:3001/conversations/alice
```

## Interactions :

caller -> api handling the request

### User -> Client (front)
Create new session (bob registers), NB : this call is the exact same as from client -> server 
```$xslt
curl -v -X POST <clientAddress>/sessions/ -H 'Content-Type: application/json' -d '{"login": "bob", "password": "pass", "address":"127.0.0.1:12345"}'
```

### Client -> Server

Create new session (bob registers)
```$xslt
curl -v -X POST <serverAddress>/sessions/ -H 'Content-Type: application/json' -d '{"login": "bob", "password": "pass", "address":"127.0.0.1:12345"}'
```

Get session details (bob wants to know how to join alice)
```$xslt
curl -v <serverAddress>/sessions/alice -H 'user: bob'
```


### Client -> Client (p2p)


## Security flaws
- users are only authenticated between them with their username as a header, this can easily be spoofed
- everything is transmitted in plain text
- clients don't authenticate between each other

Remediation, example PKI :
the server API can be publicly authenticated with a known root CA (to mitigate mim attacks between client -> server)
the server has a key pair provided at startup (sPK, sSK) to allows the system to work with several servers (using the same keys)
the server self-signs its own certificate (sC) using sSK

when a client opens a new session with the server, it receives : 
    - the self-signed certificate (sC) containing sPK 
    - its own certificate signed with sSK (c1C)
the client installs sC
the client serves the p2p API using c1C

when another client (c2) having followed the same steps above attempts to call another on c1 :
the request is made using its certificate (clientCertificate : c2C)
c1 can check if c2C has been signed with sC
c2 can check if c1C has been signed with sC
then mTLS secures the communication between the peers

## TODO
- change all the log.Println to logger.Log
- add the tests to remaining usecases, at least happy cases, server is fine for the rest

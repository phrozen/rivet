# Rivet

**Rivet is in early stages and the API can change.**

Rivet is a SaaS key/value store written in Go using [Echo](https://github.com/labstack/echo) framework and [BoltDB](https://github.com/boltdb/bolt) and it is designed with simplicity and performance in mind.

### Features

+ Fast
+ Simple to use
+ No external dependecies
+ Simple configuration
+ User based storage
+ Basic Authentication (preferably over SSL or private networking)
+ Sessions based on token headers
+ CORS configuration
+ CRUD operations on keys with binary string values
+ Database backups and Snapshots

### Usage

Rivet provides a simple way to do CRUD operations on keys via a RESTFUL API.

It provides authentication based on login and a ```Session Token``` which can be securely passed to client side apps, this token can expire or be invalidated by the user any time. It provides user namespace which support **route like** keys.

To start using Rivet (after user has been created) the user must log in first to the endpoint:

**[GET]** - **/login** (credentials via Basic Authentication)

This will return a ```Session Token``` which is used to make requests to the data store for as long as it is valid, token can expire, or can be invalidated logging out like:

**[GET]** - **/logout** (credentials via Basic Authentication)

Once logged in the ```Session Token``` is used to query the RESTFUL API and must be present in every request as a ```X-Session-Token``` Header. This can be done directly in the client as the credentials are used (at best) once to gain access to the data. The token can also be passed around to services to access the store, and can be invalidated at any moment when the session is closed.

As a general rule, if you get a **[401] Unauthorized** Status Code, means your App should login and get a new ```Session Token``` from the endpoint.

#### LIST
```
[GET] - /store/<user>?limit=&offset=
```
Lists the keys on the user namespace up to ```Limit```, you can optionally pass ```Limit``` and ```Offset``` as URL variables, ```Offset``` is used for pagination and should be the last value of a previous call.

### CRUD

CRUD operations are done by calling the corresponding endpoints  to a key, **Rivet** is designed in a way that naturally uses route like keys for organization. For example, you can organize your data like:

```
/store/<user>/profile
/store/<user>/profile/bio
/store/<user>/profile/address/billing
...
```

It is noted that with key/value data stores, *POST* and *PUT* operations realize the same operation. The only difference is that **Rivet** will respond with a **[201] Created** Status Code on a *POST* request if the key did not exist previously. This is useful if you need to keep track of your data in any way. *POST* and *PUT* requests expect the data to be in RAW format on the Body of the request.

#### READ
```
[GET] - /store/<user>/<key>
```

#### CREATE
```
[POST] - /store/<user>/<key>
```

#### UPDATE
```
[PUT] - /store/<user>/<key>
```

#### DELETE
```
[DELETE] - /store/<user>/<key>
```

#### WEBSOCKET
*(Not implemented yet)*

### Todo
+ Key Search/Prefix Search
+ Websocket support
+ Multiple Session Tokens/Read Only Tokens
+ Admin web interface
+ Database Stats/Logs
+ Hot Swap

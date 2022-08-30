# MegaSDK-REST 
MegaSDK downloading functionality as a rest server. This project is jaskaranSM's [megasdkgo](https://github.com/jaskaranSM/megasdkgo) project wrapped in a rest server.

## Documentation
Documentation is divided into two categories.

### Server usage
Read output of --help command

### Client usage
There are basically five endpoints on server.
```
POST /login
POST /adddl
POST /canceldl
GET /dlinfo/{gid}
GET /ping
```

**Login**: 
Do a POST on /login with payload as
``` 
{"email":email,"password":password}
```
Process the response.

**Add download**:
Do a POST on /adddl with payload as 
```
{"link":link,"dir":directory}
```
The directory in the context is the directory on the machine where the server is running, the response will return a random string called gid which you are supposed to store for later use.

**Cancel download**:
Do a POST on /canceldl with payload as 
```
{"gid":gidToCancel}
```
process the response.

**Get current info of the download by gid**:
Do a GET on 
```
/dlinfo/{gid}
``` 
The gid here is variable, if the server have that dl in its storage then it will return details of the download, if not then will return empty details with friendly message about what went wrong.

**Ping**: 
Just for testing if the server is up or not.

**NOTE**: 
Only supports GNU/Linux at the moment.

### Credits
- [jaskaranSM](https://github.com/jaskaranSM) - Created the project
- [Viswanath](https://github.com/ViswanathBalusu) - Added some features

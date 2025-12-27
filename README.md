Below is a minimal but practical RESTful API in Go for managing FreeBSD jails.
It exposes endpoints to list, create, start, stop, and delete jails by shelling out to native FreeBSD tools (jls, jail, jexec).

Important

This API must run as root.

Exposing jail control over HTTP is dangerous unless protected (TLS, auth, firewall).

This is a reference implementation, not production-hardened.

| Method | Endpoint              | Description             |
| ------ | --------------------- | ----------------------- |
| GET    | `/jails`              | List jails              |
| POST   | `/jails`              | Create & start a jail   |
| POST   | `/jails/{name}/start` | Start jail              |
| POST   | `/jails/{name}/stop`  | Stop jail               |
| DELETE | `/jails/{name}`       | Remove jail             |
| POST   | `/jails/{name}/exec`  | Run command inside jail |

Example Jail JSON

,,,,

{
  "name": "testjail",
  "path": "/usr/jails/testjail",
  "ip": "192.168.1.50",
  "host_hostname": "testjail.local"
}

,,,,

Example Requests

List jails

,,,
curl http://localhost:8080/jails
,,,

Create jail

,,,
curl -X POST http://localhost:8080/jails \
  -H "Content-Type: application/json" \
  -d '{"name":"testjail","path":"/usr/jails/testjail","ip":"192.168.1.50","host_hostname":"testjail.local"}'
,,,

Execute command

,,,
curl -X POST http://localhost:8080/jails/testjail/exec \
  -H "Content-Type: application/json" \
  -d '{"command":["/bin/ls","/"]}'
,,,


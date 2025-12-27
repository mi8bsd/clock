Below is a minimal but practical RESTful API in Go for managing FreeBSD jails.
It exposes endpoints to list, create, start, stop, and delete jails by shelling out to native FreeBSD tools (jls, jail, jexec).

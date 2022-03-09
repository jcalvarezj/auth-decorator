# Middleware Auth Example

This is a short example of authorization validation through middleware. The AllowRoles decorator includes an input string of roles into the handler's request context, which will be validated against the current role (sent through headers for a quick example) in the AuthorizeRoles middleware

Use `go run main.go` to execute.

This server responds to requests on the following resources:

- **/ [GET]** for anyone with an existing role

- **/ [POST]** for admin only

- **/free-resource [GET]** for anyone (not enforcing auth.)

- **/foo/bar [GET]** for roles "admin", "treasury", and "lawyer"

- **/foo/bar [POST]** for roles "admin" and "lawyer"

- **/foo/bar [PUT]** for roles "admin" and "treasury"

- **/foo/bar [PATCH]** for roles "admin" and "treasury"

Example:

```curl -X POST -H "Role: admin" localhost:8080/foo/bar```
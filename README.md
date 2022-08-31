# Ticketing-system

## A simple open-source ticketing system written in go.

### Admin user
 `admin:admin // change this password`

----
**POST** `/api/register`
```
Cookie: session=<sessionid>
```

```
{
    "username":"John Doe",
    "password":"secretpassword123",
    "role":"programmer"
}
```
----
**POST** `/api/login`

```
{
    "username":"John Doe",
    "password":"secretpassword123"
}
```
----
**POST** `/api/changepassword`
```
Cookie: session=<sessionid>
```

```
{
    "username":"John Doe",
    "password":"secretpassword123"
    "newpassword":"password345"
}
```
----
**DELETE** `/api/logout`

```
Cookie: session=<sessionid>
```
----
**POST** `/api/tickets`
```
Cookie: session=<sessionid>
```
```
{
    "title":"This is a ticket",
    "content":"Lorem Ipsum"
}
```


----
**PUT** `/api/tickets/{id}`
```
Cookie: session=<sessionid>
```
```
{
    "title":"This is a ticket update",
    "content":"Lorem Ipsum Lorem Ipsum"
}
```
----
**PUT** `/api/tickets/{id}/{status}`

{status}
- inprog
- closed
```
Cookie: session=<sessionid>
```
----
**GET** `/api/tickets?user={user1}&user={user2}`

----
**DELETE** `/api/tickets/{ticketid}`
```
Cookie: session=<sessionid>
```
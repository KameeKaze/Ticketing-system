# Ticketing-system

## A simple open-source ticketing system written in go.

### Admin user
 `admin:admin // change this password`

----
**POST** `/register`
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
**POST** `/login`

```
{
    "username":"John Doe",
    "password":"secretpassword123"
}
```
----
**POST** `/changepassword`
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
**DELETE** `/logout`

```
Cookie: session=<sessionid>
```
----
**PUT** `/tickets/create`
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
**POST** `/tickets/{id}`
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
**POST** `/tickets/closed/{id}`
```
Cookie: session=<sessionid>
```
----
**POST** `/tickets/inprog/{id}`
```
Cookie: session=<sessionid>
```
----
**GET** `/tickets?user={user1}&user={user2}`

----
**DELETE** `/tickets/{ticketid}`
```
Cookie: session=<sessionid>
```
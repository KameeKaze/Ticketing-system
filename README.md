# Ticketing-system

## A simple open-source ticketing system written in go.
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
**GET** `/tickets?user={user1}&user={user2}`

----
**DELETE** `/tickets/{ticketid}`
```
Cookie: session=<sessionid>
```
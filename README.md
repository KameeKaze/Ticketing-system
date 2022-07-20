# Ticketing-system

## A simple open-source ticketing system written in go.

**POST** `/register`

```
{
    "username":"John Doe",
    "password":"secretpassword123",
    "role":"programmer"
}
```

**POST** `/login`

```
{
    "username":"John Doe",
    "password":"secretpassword123"
}
```

**DELETE** `/logout`

```
Cookie: session=<sessionid>
```

**POST** `/tickets/create`

```
{
    "issuer":"John Doe",
    "title":"This is a ticket",
    "content":"Lorem Ipsum"
}
```

**GET** `/tickets?user={user1}&user={user2}`

**DELETE** `/tickets/{ticketid}`

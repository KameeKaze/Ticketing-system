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

**POST** `/ticket/create`

```
{
    "issuer":"John Doe",
    "title":"This is a ticket",
    "content":"Lorem Ipsum"
}
```
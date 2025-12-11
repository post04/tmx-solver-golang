# get cookies/sessionID normal
POST x.x.x.x:8081/api/get-cookies
{
    "apiKey": "xxx",
    "proxy": "http://user:pass@ip:port",
    "site": "exampleSite",
    "uuid": "session_id",
    "url": "https://example.com/checkout?something=1",
    "mobile": false
}

response:
{
    "userAgent": "xxx",
    "sessionID": "the sessionID you input or a randomly generated one is returned back to you here",
    "error": "conditional, if an error happens, it returns here."
}

# get cookies/sessionID polling
POST x.x.x.x:8081/api/get-cookies-polling
{
    "apiKey": "xxx",
    "proxy": "http://user:pass@ip:port",
    "site": "exampleSite",
    "uuid": "session_id",
    "url": "https://example.com/checkout?something=1",
    "mobile": false
}

response:
{
    "id": "xxx",
    "status": "pending",
    "createdAt": 123
}

POST x.x.x.x:8081/api/get-cookies-polling-status
{
    "apiKey": "xxx",
    "taskID": "xxx"
}

response:
{
    "userAgent": "xxx",
    "sessionID": "the sessionID you input or a randomly generated one is returned back to you here",
    "error": "conditional, if an error happens, it returns here."
}

# websocket system
wss://x.x.x.x:8081/ws

your first message should be:
{
    "action": "auth",
    "data": {
        "apiKey": "xxx"
    }
}

Then it responds with `authorized`

After that, you send solves like this:
{
    "action": "solve",
    "data": {
        "apiKey": "xxx",
        "proxy": "http://user:pass@ip:port",
        "site": "exampleSite",
        "uuid": "session_id",
        "url": "https://example.com/checkout?something=1",
        "mobile": false
    }
}

Then it will respond with:
{
    "proxy": "http://user:pass@ip:port",
    "userAgent": "xxxx",
    "sessionID": "xxxx"
}

You use the proxy it returns to match your solve to a session. You should always be using sticky proxies
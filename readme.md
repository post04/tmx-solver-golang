# Solver
This solver is written in golang, of course. It uses a mixture of a custom system and [go-fast](https://github.com/T14Raptor/go-fAST) + regex to gather dynamic keys/URLs.

This will never be updated, please do not ask.

This solver supports 3 methods, 1 request (send 1 request, solver does it's thing, that's it), polling (send 1 request, get an ID back, then request again to see if it's solved), then websocket which is what I prefer personally, it's way better in my opinion.

# How to use
It's super simple, load up fingerprints into your mongodb database, run it with `go run .`. Next use the admin API to generate a key, set uses to `-1` to have infinite uses and time to 999 for "infinite" days. Then just read the docs.md for the information on how to use the API.

# Notes
The mobile solver doesn't work.
The TLS is technically not 1:1.
Since it's outdated, on some sites it will not be able to even attempt a solve because the URL/keys it's trying to grab are removed. I noticed this on bestbuy I think.

# Credits
Me, however, please check out my blog at [antibot.blog](https://antibot.blog). This project won't be a blog post, however, I do try to upload blogs a few times a year :)
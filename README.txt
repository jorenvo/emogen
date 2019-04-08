This is a link shortener that shortens links using UTF-8 emojis

API
---
/r POST with link=https://example.com

/r/🐯🐁🙃 GET 301 redirects to https://example.com

Docs
----
List of UTF-8 emojis: https://www.unicode.org/emoji/charts/full-emoji-modifiers.html

There's around ~1000 we can use. 1 emoji holds 10 bits of
information. Concatenating 3 together holds 30 bits of information,
should be enough for our purposes.

This is a link shortener that shortens links using UTF-8 emojis.

API
---
/r POST with {link: "https://example.com"} returns {link: "/🙃🐰🦊"}

/r/🙃🐰🦊 GET 301 redirects to https://example.com

Docs
----
List of UTF-8 emojis: https://www.unicode.org/emoji/charts/full-emoji-modifiers.html

There's 1431 we can use (without skin tones, those could be added
later). 1 emoji holds >10 bits of information. Concatenating 3
together contains ~32 bits of information, this should be enough for
our purposes.

Redis is used as the database because it's easy to setup and our
database needs are simple. Redis' save option is sufficient for
persistence and backups. With Redis' default configuration a maximum
of 15 minutes of data could be lost. This is acceptable because we're
just shortening links and this is a toy project.

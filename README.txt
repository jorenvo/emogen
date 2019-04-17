This is a link shortener that shortens links using UTF-8 emojis. A
live demo with emogen-frontend is available at https://www.jvo.sh/r/.

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

Generate one number that will be used to find 3 emoji's for a
URL. It's range is [0, 1431**3]. For a number n the first emoji is n %
1431, second is (n // 1431) % 1431 and the last is (n // 1431 // 1431)
% 1431.

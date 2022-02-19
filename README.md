# Wordle over DNS

_Why not? It's all the rage._

I wrote [Wikipedia over DNS](https://dgl.cx/wikipedia-dns) a long time ago.
It's amazing how much can be put into TXT records. They are sadly, although in
keeping with their name designed for text and characters outside the ASCII
range will be escaped by most tools. However it's easy to write a shell
function to reverse that and do unescaping. Then we can serve that command over
DNS too, so it's easy to get hold of!

So, here is Wordle over DNS:

    $ host -t txt wd.ip.wtf
    wd.ip.wtf descriptive text "Welcome to Wordle over DNS! Today's puzzle is #1: <guess>.1.wd.ip.wtf"
    wd.ip.wtf descriptive text "This shell function makes it easier to play" "wd() { dig +short txt $1.1.wd.ip.wtf | perl -pe's/\\([0-9]{1,3})/chr$1/eg'; }"

For the sake of this example I've replaced the `1` with `example` which is one
I prepared earlier so I'm not revealing an actual answer!

    $ wd() { dig +short txt $1.example.wd.ip.wtf | perl -pe's/\\([0-9]{1,3})/chr$1/eg'; }
    $ wd crane
    "⬜⬜🟨🟨🟨"
    $ wd reads
    "⬜🟨🟨⬜🟩"
    $ wd sense
    "🟨🟨🟨⬜⬜"
    $ wd names
    "🟩🟩🟩🟩🟩"

[Note crane [was considered](https://www.youtube.com/watch?v=fRed0Xmc2Wg) one
of the best openers, but this is using its own list of words, so that's not
really true here.]

## Implementation

This is implemented as a very simple standalone DNS server in Go using [miekg's
DNS library](https://github.com/miekg/dns). Code is at
https://github.com/dgl/wordle-dns.

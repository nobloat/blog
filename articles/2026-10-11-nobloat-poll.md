# ][✓ nobloat/poll

So, I wanted to send out a poll for our yearly "zeitkapsl Punsch" event, where we invite people who helped [zeitkapsl](https://zeitkapsl.eu) become what it is today. Early customers, former development team members, etc. Well...

> "... I opened up Doodle, trying to create a new poll for an event. I was immediately greeted with flashing popups stating I have to upgrade to pro now. I kept searching for where to delete my existing polls. It was not easy. I finally figured it out and was able to create the poll for my punsch event. But it left me wondering: how is everybody else tolerating this?"

![A hand-drawn poll: "How to name the office plant?" Options: go fern(), /dev/leaf, Monstera Monolith and Photosynthesis as a Service](../images/office-plant.svg)

## TL;DR

[nobloat/poll](https://github.com/nobloat/poll) is a self-hostable single-binary Doodle alternative written in plain `go` with [no further dependencies](https://github.com/nobloat/poll/blob/main/go.mod). It doesn't require a database, job queue, npm or node. It works without signups. Every poll is a single JSON file, the whole thing ships as one 3.5 MB binary, and a page loads ~27 KB (6.9 kB + 20 kB uncompressed JavaScript and CSS). You can use it at [poll.nobloat.org](https://poll.nobloat.org). There is no "pro version" of it.

## There was Doodle

I usually used Doodle for that, but apparently a lot has happened there over the past year, and the path to [enshi***ation has proceeded](https://www.youtube.com/watch?v=T4Upf_B9RLQ) a lot further:

- You have to leave your e-mail address even as a participant.
- You can no longer have more than one poll in total in the "free account", even if the old one is already closed.
- You can't export your participant list.
- Loading the poll page takes around 3 seconds, and downloads 4.7 MB JavaScript.
- You are flashed by so many pop-ups.

There is of course [nuudel](https://nuudel.digitalcourage.de/) and a couple other [acceptable alternatives](https://www.heise.de/tipps-tricks/Doodle-Alternativen-Die-besten-Tools-zur-Terminfindung-4964023.html), but I wanted something I can easily self-host on my quite dated own cabinet "server" (an old ThinkPad X260 of course).

![My cabinet "server", a ThinkPad X260](../images/thinkpad-x260.png)

It's been over a year since I last built something *nobloat-style*, so I thought this is the perfect opportunity to spend some hours and build something I personally **enjoy building and using**.


## My requirements checklist

The goal was to have a tool ready by the end of the weekend that allows me to finally send out a poll for when to do the "Punsch" event.

### Functional

- Optionally send out e-mails to participants (if they choose to by offering an optional e-mail address)
- Allow public poll creation without a register/login dance
- Basic spam protection via a simple proof of work performed by the browser
- [i18n](https://github.com/nobloat/poll/tree/main/i18n) support
- [basic theming](https://github.com/nobloat/poll/tree/main/static/themes) to allow for branding and make it appropriate for different scenarios (birthday parties, dinner invite, business meetings, etc.)

### Non-functional

- Having fun in the process and create something I can enjoy looking at
- Single statically linked cross-compilable binary I can basically deploy on anything that offers a filesystem.
- As few external build-time and run-time dependencies as possible, so this stays low maintenance in the upcoming years. 
- Keep the state in plain text files so it's easy to reason about, backup, migrate, modify, bulk edit/create/delete.
- Absolutely no **node, npm, yarn** madness. I usually build most of my web frontends with [Svelte](https://svelte.dev/) or [SvelteKit](https://svelte.dev/docs/kit/introduction), but this requires npm and a ton of other node modules. Again, keeping this low-maintenance without [renovate-bot](https://docs.renovatebot.com/) yelling at me every other day.
- Simple configuration that works inside and outside of containers


[As it turns out](https://github.com/nobloat/poll/blob/main/go.mod), **plain go std-lib** is completely sufficient to build a simple replacement for Doodle.

```
cat go.mod package.json
module github.com/nobloat/poll

go 1.27.0
cat: package.json: No such file or directory
```

## Bypassing signup and login

To avoid making users register and log in, every poll gets a secret owner token when it is created. The link containing it is all you need to manage the poll.

If the owner chooses to provide a mail address, this unique link is sent to them. If not, they can copy/paste or bookmark it.

```
Hi Peter,

your poll "Punsch Event" is live.

Share this link with participants:
http://localhost:8737/p/d7fminnmxnvygaxg

Manage it (see answers, edit, close, delete):
http://localhost:8737/p/d7fminnmxnvygaxg/m/TMP6BZXB2SE4TOP3PVZODGRS6U

The manage link is personal. Anyone with it can edit, close or delete the poll, so don't share it.
```

The same happens when answering a poll. Every answer gets its own token, and the link allows you to view, edit or delete your answer.

```
Hi Peter,

thanks for answering "Punsch Event". Here is your personal link to change or remove it:

http://localhost:8737/p/d7fminnmxnvygaxg/a/CGYZDMGL7UJ5WA4RR4ON66HBEN

You'll be notified by email when the final decision is made.

This link is personal, don't share it. If you didn't answer this poll, ignore this email or use the link to remove the answer.
```

All mails are localized, just like the web interface. The language is remembered per participant, so everyone gets their mails in the language they used when answering, even if the poll owner speaks a different one.

## Creating a poll or event

You can create a new poll [here](https://poll.nobloat.org/new?lang=en)

![The poll creation form](../images/screenshot-poll-new.png)

It asks you to enter a title and a description.
Options are simply a `<textarea>` where each line represents an option.
You can also switch between *single choice* and *multiple choice*.
For better privacy among voters, you can also choose whether participants can see each other's answers.
There is no separate "event" poll type. If an option looks like a date, the poll becomes an event poll by itself (more on that [below](#calendar-invites)). Since options are just lines of text, you can also copy and paste a whole list of them in one go.

## Answering a poll

![Editing an answer via the personal link](../images/screenshot-poll-answer.png)

If the poll creator supplied an e-mail address, they receive an update mail when answers are coming in.

```
New answer from Grace: Pizza Friday

# Answers

01 Margherita  [2] ■■
   Ada Lovelace, Grace Hopper
02 Diavola  [2] ■■
   Linus, Grace Hopper
03 Quattro Formaggi  [2] ■■
   Ada Lovelace, Grace Hopper
04 Hawaii  [1] ■
   Rubber Duck

http://localhost:8737/p/pizza-friday/m/ASEIRUQD7POJBVG6VITTHGCIIK
```

## Calendar invites

Most of my polls are really about finding a date. So when an option looks like a date, the poll treats it as one. All of these work:

```
2026-12-12
12.12.2026
12.12.
Sat, 12.12.2026, 18:00–22:00
12.12. um 18:00 Uhr
```

A missing year means the next time that date comes up. A date without a time becomes an all-day event, and a start time without an end time gets one hour. If the end time is earlier than the start, the event runs past midnight. Anything that isn't a valid date, like `31.02.`, simply stays a normal text option.

When the poll is closed and the winning option is a date, everybody who left an e-mail address gets the result with an `invite.ics` attached:

```
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//nobloat//poll//EN
METHOD:PUBLISH
BEGIN:VEVENT
UID:zeitkapsl-punsch@poll.nobloat.org
DTSTAMP:20261205T091500Z
DTSTART:20261212T170000Z
DTEND:20261212T210000Z
SUMMARY:zeitkapsl Punsch
DESCRIPTION:Punsch\, Kekse and old stories.
END:VEVENT
END:VCALENDAR
```

One click and it's in their calendar. The whole format is about 40 lines of Go: escaping, folding lines at 75 bytes, and joining them with `\r\n`. There's no need for a library.


## Keyboard shortcuts

All actions can be reached via keyboard shortcuts so you don't have to use the mouse at all. Each shortcut is shown as a small `<kbd>` hint next to where it applies: `1`–`9` toggle the options, `n` jumps to the name field, `e` to the e-mail field, `s` submits, and `Ctrl+Enter` submits from inside any field. On the manage page, `c` copies the share link and `1`–`9` pick the final option.

Under the hood, it's just a `data-key` attribute in the HTML templates and a single `keydown` listener in `app.js`.

## Spam protection without a captcha

Since anybody can create polls on [poll.nobloat.org](https://poll.nobloat.org), I needed some basic protection against bots, but without captchas or third-party services.

Before the server accepts a form, the browser has to burn a bit of CPU. For one human submitting one poll, that's unnoticeable. For a bot submitting thousands, it adds up.

**1. Get a challenge.** On submit, the browser fetches `GET /pow`. The server answers with the difficulty and a challenge:

```
16 1760180000.K3J5ZQ7XW2M4N6P8R2T4V6X8YA.9f2c…e41a
```

The challenge is a timestamp, a random value and an HMAC over both, signed with a key that only lives in memory. Because of the signature, the server doesn't have to store issued challenges. It can recognize its own later.

**2. Solve it.** The browser counts up a number (the nonce) until `SHA-256(challenge + nonce)` starts with 16 zero bits:

```javascript
const zeroBits = (buf) => {
  const b = new Uint8Array(buf);
  let n = 0;
  for (const x of b) {
    if (x === 0) { n += 8; continue; }
    return n + Math.clz32(x) - 24;
  }
  return n;
};

async function solve(challenge, bits) {
  const enc = new TextEncoder();
  const batch = 64;
  for (let n = 0; ; n += batch) {
    const hashes = await Promise.all(Array.from({ length: batch },
      (_, i) => crypto.subtle.digest('SHA-256', enc.encode(challenge + (n + i)))));
    const i = hashes.findIndex((h) => zeroBits(h) >= bits);
    if (i >= 0) return String(n + i);
  }
}
```

Each extra zero bit halves the chance that a random hash qualifies, so 16 bits means about 2¹⁶ ≈ 65,000 tries on average. `crypto.subtle.digest` is built into every browser, so no library is needed. It's async, which is why the hashes are computed in batches of 64 with `Promise.all` rather than awaited one by one.

**3. Submit with proof.** A global `submit` listener holds the form back, solves the challenge and adds both values as hidden fields before sending it off:

```javascript
async function proof() {
  const [bits, pow] = (await (await fetch('/pow')).text()).split(' ');
  return { pow, nonce: await solve(pow, +bits) };
}

document.addEventListener('submit', async (e) => {
  const form = e.target;
  if (form.method !== 'post') return;
  if (form.dataset.proved) { delete form.dataset.proved; return; }
  e.preventDefault();
  // …
  const p = await proof();
  for (const k in p) form.insertAdjacentHTML('beforeend', `<input type="hidden" name="${k}" value="${p[k]}" data-pow>`);
  form.dataset.proved = 1;
  // …
  form.requestSubmit(e.submitter);
});
```

The forms themselves are plain HTML forms. They don't know anything about this.

**4. Verify.** Every `POST` on the server goes through a middleware that checks the proof. Verifying costs a single hash, compared to the client's ~65,000:

- Does the HMAC match? Otherwise the challenge is forged.
- Is it younger than 2 hours?
- Does `SHA-256(challenge + nonce)` really start with 16 zero bits?
- Has this challenge been used before? Used challenges are kept in memory until they expire, so a solved one can't be replayed.

This won't stop a determined attacker, but it makes mass submissions expensive enough that dumb bots move on.


## Skipping the database

Postgres was out of the question as a runtime dependency for this project for obvious reasons.
SQLite would have been fine for me, but the common Go driver requires `CGO`, which introduces a whole other set of problems when cross compiling. Pure-Go ports like `modernc.org/sqlite` exist, but they are yet another external dependency. So I decided to keep it even more minimal and just store the polls in simple `.json` files.

All modifications go through a single global lock. That's the bottleneck if thousands of people vote at the same time, which is not a problem I have. The poll definition as well as the answers from voters are stored in the same file.

The file is named after the generated poll ID: `data/polls/<id>.json`

```
cat data/polls/d7fminnmxnvygaxg.json
{
  "title": "Punsch Event",
  "description": "",
  "multiple": true,
  "resultsVisible": false,
  "options": [
    { "text": "12.10.2026", "start": "2026-10-12T00:00:00+02:00", "end": "2026-10-13T00:00:00+02:00", "allDay": true },
    { "text": "13.10.2026", "start": "2026-10-13T00:00:00+02:00", "end": "2026-10-14T00:00:00+02:00", "allDay": true },
    { "text": "15.10.2026", "start": "2026-10-15T00:00:00+02:00", "end": "2026-10-16T00:00:00+02:00", "allDay": true }
  ],
  "answers": [
    { "token": "CGYZDMGL7UJ5WA4RR4ON66HBEN", "name": "Peter", "email": "", "lang": "en", "choices": [1], "created": "2026-10-11T13:00:27.373786142+02:00" }
  ],
  "closed": false,
  "outcome": 0,
  "created": "2026-10-11T12:59:34.699813479+02:00",
  "owner_name": "Peter",
  "owner_email": "",
  "owner_lang": "en",
  "owner_token": "TMP6BZXB2SE4TOP3PVZODGRS6U",
  "theme": "minimal"
}
```

But what happens if the process crashes or the power goes out in the middle of writing a poll? With a plain `os.WriteFile` you could end up with a truncated, half-written `.json` file, and the poll would be gone.

To avoid that, [every write is atomic](https://github.com/nobloat/poll/blob/main/store.go): the new content goes into a temporary file in the same directory, which is flushed to disk and then renamed over the old file. A `rename` within the same filesystem is atomic, so the poll file is always either the old or the new version.

```go
f, _ := os.CreateTemp(s.dir, ".tmp-*")
f.Write(b)
f.Sync()                            // content is on disk
f.Close()
os.Rename(f.Name(), s.path(p.ID))   // atomic swap
s.syncDir()                         // the rename itself is on disk
```

The last `fsync` on the directory makes sure the rename survives a power loss too. If we crash before the rename, only a stale `.tmp-*` file is left behind, and those get cleaned up on the next startup.

## Reliably sending mails

When a poll is closed, every participant who left an e-mail address gets the outcome, so mails must not get lost. Therefore, we need some kind of reliable e-mail sending mechanism.

Go ships `SMTP` support out of the box with `net/smtp`. The easiest way is to just start one `goroutine` per recipient and call it a day. But SMTP can fail or we can decide to restart the application while mails are sent out. 

To avoid such issues, one usually uses an asynchronous job queue like `redis` or [river queue](https://riverqueue.com/) for postgres. Of course neither is an option for this project.

So we are keeping a [persisted outbox](https://github.com/nobloat/poll/blob/main/outbox.go) under `./data/mail` that stores each mail that needs to be sent out on disk.

On application startup we also check for not yet sent messages.
This gives us an [at-least-once](https://en.wikipedia.org/wiki/Message_queue#Delivery_guarantees) delivery guarantee, with duplicates only if we crash mid-send.

Shutdown of the application [happens gracefully](https://github.com/nobloat/poll/blob/820d66036f974ffbc0d8e270ce57fe09ad9985a2/main.go#L94-L104) and waits for the mail currently being sent.

## Deployment

Since there is no `CGO`, the binary is built fully static:

```
CGO_ENABLED=0 go build -tags timetzdata,nethttpomithttp2 -trimpath -ldflags="-s -w" -o poll .
upx -q --best poll
```

`timetzdata` embeds the time zone database, so the binary doesn't even need `/usr/share/zoneinfo` on the target system. `nethttpomithttp2` drops the HTTP/2 implementation, which isn't needed behind a reverse proxy anyway.

This allows the [Docker image](https://github.com/nobloat/poll/blob/main/Dockerfile) to be built `FROM scratch`: it contains nothing but the CA certificates (for SMTP via TLS) and the binary itself. It runs as an unprivileged user, and the binary refuses to start as `root` altogether.

Configuration is done via environment variables only, which works the same way inside a container, in a systemd unit or in a shell on my ThinkPad. On startup all settings are validated and every error is reported at once.

## 9.1 MB → 3.5 MB binary size

9.1 MB is not as tiny for a poll tool as I hoped for, so I looked at which packages take up the space.

```
CGO_ENABLED=0 go build -tags timetzdata,nethttpomithttp2 -trimpath -o poll-sym .
go tool nm -size poll-sym \
  | awk '$3 !~ /[Bb]/ { match($4, /^[^.(]*(\/[^.(]*)*/); s[substr($4, 1, RLENGTH)] += $2 }
         END { for (p in s) printf "%8.1f KB  %s\n", s[p]/1024, p }' \
  | sort -rn
```

Grouped a bit further:

```
crypto/*                         1.2 MB
runtime, reflect, internal/*     0.9 MB
type information                 0.6 MB
net, net/http, net/smtp          0.5 MB
encoding/json                    0.3 MB
vendor/golang.org/x/*            0.2 MB
html/template, text/template     0.2 MB
nobloat/poll                     0.1 MB
```

Most of this is expected, we need the whole `crypto` package for the TLS stack for SMTP via `STARTTLS`. 

There is an option to omit the HTTP/2 support, which we don't need because we run behind a reverse proxy anway.

```
for tags in timetzdata,nethttpomithttp2 nethttpomithttp2 timetzdata; do
  CGO_ENABLED=0 go build -tags "$tags" -trimpath -ldflags="-s -w" -o poll-test .
  echo "$(stat -c %s poll-test) $tags"
done
du -cb static templates i18n mail | tail -1
```

- `timetzdata` adds **~400 KB**: the zipped time zone database
- `nethttpomithttp2` saves **~410 KB**: the HTTP/2 implementation
- All embedded templates, CSS, JS and translations together are only **~78 KB**


### Squeezing it with UPX

Since most of the binary can't be removed, the last step is to compress it with [UPX](https://upx.github.io/). UPX packs the executable and prepends a small stub that decompresses it into memory on startup:

```
upx -q --best poll
```

```
9,584,800 bytes  poll (stripped)
3,624,388 bytes  poll (upx --best)
```

That's **3.5 MB**, a reduction of about 62%. It doesn't come for free though. Startup time increases quite a lot, but this doesn't really bother me.

```
                 startup   peak memory
poll             2 ms      12.4 MB
poll (upx)       32 ms     14.7 MB
```

## Summary

It is now Sunday afternoon, 14:00 and this was fun to build, even though it turned out to be not as minimal as I hoped for. But I do like that it doesn't use any external dependencies except for the SMTP gateway to send out emails, and even these are optional.

**I have now a poll tool I enjoy using.** - mission accomplished.

I could still leave out the JavaScript stuff itself, but I like it for the shortcuts, proof-of-work and interactivity improvements it delivers. 

Everybody else [is welcome to use](https://poll.nobloat.org) it as well. No service guarantees though, I tend to reboot my X260 from time to time.

It's licensed under [AGPL-3.0](https://github.com/nobloat/poll/blob/main/LICENSE).

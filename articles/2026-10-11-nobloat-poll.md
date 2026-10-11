# ][✓ nobloat/poll

So, I wanted to send out a poll for our yearly "zeitkapsl Punsch" event, were we invite people who helped [zeitkapsl](https://zeitkapsl.eu) becoming what it is todady. Early customers, former development team members, etc. Well...

> "... I opened up doodle, trying to create a new poll for an event. I was immediately greeted with flashing popups stating you have to upgrade to pro now. Turns out you are now allowed to have only one event, even if that is already closed. I kept searching for where to delete the existing polls. It was not easy. I finally figured it out and was able to create the poll for my punsch event. But it left me wondering, how is everybody else tolerating this."

## TL;DR

[nobloat/poll](https://github.com/nobloat/poll) is a self-hostable single binary Doodle alternative written in plain `go` with [no further dependencies](https://github.com/nobloat/poll/blob/main/go.mod). It doesnt require a database, job queue, npm or node. It works without signups. Every poll is a single JSON file, the whole thing ships as one 3.5 MB binary, and a page loads ~27 KB (6.9 kB + 20 kB uncompressed JavaScript and CSS). You can use it at [poll.nobloat.org](https://poll.nobloat.org). There is no "pro version" of it.

## There was doodle

Usually I used Doodle for that, but apparently a lot has happened there over the past year, and the path to enshi***ation has proceeded a lot further:

- You have to leave your e-mail address even as a participant
- You can no longer have more than one poll in total in the "free account".
- You can't export your participant list
- Loading the poll page takes around 3 seconds, and downloads 4.7 MB in resources.
- You are flashed by so many pop-ups.

There is of course [nuudel](https://nuudel.digitalcourage.de/) and a couple other [acceptable alternatives](https://www.heise.de/tipps-tricks/Doodle-Alternativen-Die-besten-Tools-zur-Terminfindung-4964023.html), but I wanted something I can easily self-host on my quite dated own cabinet "server" (an old Thinkad X260 of course).

![My cabinet "server", a ThinkPad X260](../images/thinkpad-x260.png)

Since the last time I did something *nobloat-style* has been over a year ago, I thought this i the perfect opportunity to spend some hours and build something I personally **enjoy building and using**.


## My requirements checklist

The goal was to have a tool ready by the end of the weekend that allows me to finally send out a poll for when to do the "Punsch" event.

### Functional

- Optinally send out e-mails to participants (if they chose to by offering an optional e-mail address)
- Allow public poll creation without a register/login dance
- Basic spam protection via a simple proof of work performed by the browser
- [i18n](https://github.com/nobloat/poll/tree/main/i18n) support
- [basic theming](https://github.com/nobloat/poll/tree/main/static/themes) to allow for branding and make it appropriate for different scenarios (birthday parties, dinner invite, business meetings, etc.)

### Non-functional

- Having fun in the process and create something I can enjoy looking at
- Single statically linked cross-compilable binary I can basically deploy on anything that offers a filesystem.
- As little external build-time and run-time dependencies as possible, so this stays low maintenance in the upcoming years. 
- Keep the state in plain text files so it's easy to reason about, backup, migrate, modify, bulk edit/create/delete.
- Absolutely no **node, npm, yarn** madness. I usually build most of my web frontends with [Svelte](https://svelte.dev/) or [SvelteKit](https://svelte.dev/docs/kit/introduction), but this requires npm and a ton of other node modules. Again, keeping this low-maintenance without [renovate-bot](https://docs.renovatebot.com/) bot yelling at me every other day.
- Simple configuration that works inside and outside of containers


[As it turns out](https://github.com/nobloat/poll/blob/main/go.mod), **plain go std-lib** is completely sufficient to build a simple replacement for doodle.

```
cat go.mod package.json
module github.com/nobloat/poll

go 1.27.0
cat: package.json: No such file or directory
```

## Bypassing singnup and login

To avoid having users to login and register, we create unique identifiers for them, that get created on poll creation.

If they chose to provide a mail address, this unique Link is sent to them. If not, they can copy/paste bookmark it.

```
Hi Peter,

thanks for answering "Punsch Event". Here is your personal link to change or remove it:

http://localhost:8737/p/d7fminnmxnvygaxg/a/CGYZDMGL7UJ5WA4RR4ON66HBEN

You'll be notified by email when the final decision is made.

This link is personal, don't share it. If you didn't answer this poll, ignore this email or use the link to remove the answer.
```

Same happens on answer creation. The link allows to view, edit or delete your answer.

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

![The poll creation form](../images/poll-new.png)

It asks you to enter a title and a descriptions.
Options are simply a `<textarea>` where each line represents an option.
You can also switch between *single choice* and *multiple choice*.
For better privacy among voters, you can also enable or disable if other participants can view each other's answers.
To avoid having complex different poll types, I decided to have just one form but make it a bit smarter in case options are a date.
If we detect the chosen answer is a date, the application automatically attaches an `.ics` file with the selected date/time and event description. This way, choosing a date, or choosing a pizza together can be done the same way. It also allows for easy bulk creation of answers since you can just copy/paste them in.

## Answering a poll

![Editing an answer via the personal link](../images/poll-answer.png)

If the poll creator supplied an e-mail address, s:he receives an update mail when answers are coming in.

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


## Keyboard short cuts

All actions can be reached via keyboard shortcuts so you don't have to use the mouse at all. Each shortcut is indicated next to where it applies.


## Skipping the database

Postgres was out of question as a runtime dependency for this project for obvious reasons. 
Sqlite would have been fine for me, but it requires `CGO` which introduces a whole lot other set of problems when cross compiling. So I decided to keep it even more minimal and just store the polls in simple `.json` files.

All modifications go through a single global lock. That's the bottleneck if thousands of people vote at the same time, which is not a problem I have. The poll definition as well as the answer's from voters are stored in the same file. 

The file is named according to the generated event ID: ``

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
}⏎
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

## Relyably sending mails

In order to not lose any outgoing messages (e.g. when the poll is closed), all recipients are contacted with the result, we need some kind of relyable e-mail sending mechanism.

Go does bring out of the box `SMTP` support. The easiest way is to just start one `goroutine` per recipient and call it a day. But SMTP can fail or we can decide to restart the application while mails are sent out. 

To avoid such issues, one usually uses an asynchronous job queue like `redis` or [river queue](https://riverqueue.com/) for postgres.

Of course neither is an option for this project.

So we are keeping a [peristed outbox](https://github.com/nobloat/poll/blob/main/outbox.go) under `./data/mail` that stores each mail on disk, that needs to be sent out. 

On application startup we also check for not yet sent messages. 
This gives us [at least once] delivery gurantee, with re-sends happening only mid-crashes.

Shutdown of the application [happens gracefully](https://github.com/nobloat/poll/blob/main/main.go#L96) and awaits at least the currently sent out message.

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

Stripped binaries have no symbols, so this needs a build without `-s -w`. `go tool nm -size` lists every symbol with its size, and a bit of `awk` sums them up per package (skipping `B`/`b` symbols, which are zero-initialized memory and take up no space in the file):

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

Comparing builds with and without the build tags:

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

That's **3.5 MB**, a reduction of about 62%. It doesn't come for free though. Startup time increases quite alot, but this doesnt really bother me.

```
                 startup   peak memory
poll             2 ms      12.4 MB
poll (upx)       32 ms     14.7 MB
```

## Summary

It is now sunday afternoon, 14:00 and this was fun to build, even though it turned out to be not as minimal as I hoped for. But I do like that it doesnt use any external dependencies except for the SMTP gateway to send out emails, and even these are optional.


**I have now a poll tool I enjoy using.** - mission accomplished.


I could still leave out the JavaScript stuff itself, but I like it for the shortcuts and interactivity improvements it delivers. 

Everybody else [is welcome to use](https://poll.nobloat.org) it as well. No service guarantees though, I tend to reboot my X260 from time to time.

It's licesned under [AGPL-3.0](https://github.com/nobloat/poll/blob/main/LICENSE).

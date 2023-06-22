<h1 align="center">
  <img src="static/terminal.png" alt="terminal" width="900px">
  <br>
</h1>

<p align="center">
  <a href="#notes">Notes</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#running-subchase">Running subchase</a> •
  <a href="#to-do-functionality">To-Do functionality</a>
</p>

`subchase` is a subdomain discovery tool that returns (almost always) valid subdomains for websites by analyzing search results from Google and Yandex search engines. The goal of `subchase` is not to find all subdomains, but to find a few subdomains that were not found by other tools.

# Notes

- There are false positives in the results. Methods to filter results have not yet been implemented.
- The results may vary from run to run.
    - This is usually due to captchas that cannot be bypassed, and the frequency of which cannot be predicted.

# Installation

```sh
go install -v github.com/tokiakasu/subchase/cmd/subchase@latest
```

# Usage

```sh
$ subchase -h

Usage of subchase:
  -d string
        Specify the domain whose subdomains to look for (ex: -d google.com)
  -silent
        Remove startup banner
```

# Running subchase

To run the tool on a target, just use the following command.

```console
$ subchase -d google.com
               __         __
   _______  __/ /_  _____/ /_  ____ _________
  / ___/ / / / __ \/ ___/ __ \/ __ `/ ___/ _ \  
 (__  ) /_/ / /_/ / /__/ / / / /_/ (__  )  __/
/____/\__,_/_.___/\___/_/ /_/\__,_/____/\___/  v0.1.0

earthengine.google.com
meet.google.com
classroom.google.com
passwords.google.com
cloud.google.com
jibe.google.com
books.google.com
messages.google.com
adsense.google.com
sites.google.com
images.google.com
support.google.com
careers.google.com
ads.google.com
store.google.com
checks.google.com
asia.google.com
firebase.google.com
accounts.google.com
mydevices.google.com
myactivity.google.com
mymaps.google.com
atap.google.com
forms.google.com
admin.sites.google.com
about.artsandculture.google.com
ipv4.google.com
ipv6.google.com
assistant.google.com
fonts.google.com
```

# To-Do functionality

- [ ] Add option to output content-length along with domains
- [ ] Add option to output results in JSON

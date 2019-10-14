# csp

Given a list of hosts, this small utility fetches all whitelisted domains from the hosts' CSPs. I use this for reconnaissance
purposes while bug bounty hunting.

<a href="https://www.buymeacoffee.com/edoverflow" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

# Usage

```
$ cat hosts.txt
http://example.com/
$ cat hosts.txt | csp
example.com
subdomain.example.com
...
```

Set concurrency level using the `-c` flag.

```
$ csp -h
Usage of csp:
  -c int
    	set the concurrency level (default 20)
$ cat hosts.txt | csp -c 2
...
```

# Installation

```
$ go get -u github.com/edoverflow/csp
```

You can also [download a binary](https://github.com/EdOverflow/csp/releases) and put it in your `$PATH` (e.g. in `/usr/bin/`).

# Contributing

I welcome contributions from the public.

### Using the issue tracker 💡

The issue tracker is the preferred channel for bug reports and features requests.

### Issues and labels 🏷

The bug tracker utilizes several labels to help organize and identify issues.

### Guidelines for bug reports 🐛

Use the GitHub issue search — check if the issue has already been reported.

# Credit

Thank you to [@TomNomNom](https://github.com/tomnomnom), [@jimen0](https://github.com/jimen0), and [@003random](https://github.com/003random) for their help.

# aboxes

Run one or more commands on multiple remote hosts via SSH.
This is most commonly used for retrieving ad-hoc information without too much fuzz.

## Usage

```bash
$ aboxes run -t theta,testalex.h -c hostname
INFO[0000] theta: theta
INFO[0000] testalex.h: j6yt29-testalex-magweb-do.nodes.hypernode.io

# Format output with Go template syntax
$ aboxes run -t theta,testalex.h -c hostname --format "{{.Target}} -> {{.Stdout}}"
INFO[0000] theta -> theta
INFO[0000] testalex.h -> j6yt29-testalex-magweb-do.nodes.hypernode.io

# Prevent shell escaping hell by placing commands in scripts
$ cat testscript.sh
#!/usr/bin/env bash
ip a | grep eth0 | awk '{print $2}' | awk -F '/' '{print $1}'
$ aboxes run -t theta,testalex.h --file ./testscript.sh
INFO[0001] theta: eth0:
1.2.3.4
INFO[0001] testalex.h: eth0:
2.3.4.5
```

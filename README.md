# filterschedule

## Name

*filterschedule* - checks and filters requests that match some predefined criterias.

## Description

*filterschedule* checks and filters a request that matches some predefined rules based on:
* the requested name
* optionally, the IP address of the client
* optionally, the date/time when filtering should be applied

## Syntax

~~~
filterschedule ./path/to/filterschedule.yaml
~~~

## Example

~~~ corefile
. {
    filterschedule /etc/coredns/filterschedule.yaml
    forward . 8.8.8.8 9.9.9.9
}
~~~

The configuration file is written using YAML.
~~~ yaml
- filter:
    sites: ["darkweb"]
  for:
    all: true
  when:
    always: true

- filter:
    sites: ["discord", "instagram", "twitch", "twich"]
  for:
    all: false
    hosts:
      - 192.168.1.100
      - 192.168.1.101
  when:
    always: false
    schedule:
      - days: ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
        hours:
          - from: 00:00
            to: 20:00
          - from: 21:00
            to: 23:59
      - days: ["Saturday", "Sunday"]
        hours:
          - from: 00:00
            to: 07:00
          - from: 21:00
            to: 23:59
~~~

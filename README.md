# lavoro

An experiment in running a job until completion with the Kubernetes API

## Running the example:

```bash
go build && ./lavoro 
2020/10/27 22:33:37 Accepted: {ping default {alpine:3.12 [ping -c 5 google.com]}}
ping-1603838017
2020/10/27 22:33:37 ping-1603838017
2020/10/27 22:33:39 Pending
2020/10/27 22:33:40 Running
2020/10/27 22:33:45 Succeeded
PING google.com (216.58.206.142): 56 data bytes
64 bytes from 216.58.206.142: seq=0 ttl=116 time=2.410 ms
64 bytes from 216.58.206.142: seq=1 ttl=116 time=1.763 ms
64 bytes from 216.58.206.142: seq=2 ttl=116 time=1.347 ms
64 bytes from 216.58.206.142: seq=3 ttl=116 time=1.237 ms
64 bytes from 216.58.206.142: seq=4 ttl=116 time=1.221 ms

--- google.com ping statistics ---
5 packets transmitted, 5 packets received, 0% packet loss
round-trip min/avg/max = 1.221/1.595/2.410 ms
2020/10/27 22:33:45 Deleting job ping-1603838017..OK.
```

With a web-scrape:

```bash
go build && ./lavoro 2020/10/27 22:22:44 Accepted: {check-inlets default {alexellis2/check-inlets [mocha --recursive ./integration-tests]}}
check-inlets-1603837364
2020/10/27 22:22:44 check-inlets-1603837364
2020/10/27 22:22:44 Pending
2020/10/27 22:22:45 Pending
2020/10/27 22:22:46 Pending
2020/10/27 22:22:47 Running
2020/10/27 22:22:48 Running
2020/10/27 22:22:49 Succeeded
Started HeadlessChrome/85.0.4182.0
  inlets blog
Promise { <pending> }
    ✓ renders OK (1863ms)
  1 passing (2s)
2020/10/27 22:22:50 Deleting job check-inlets-1603837364..OK.
```

## Example workloads

Web scraping with:
* [puppeteer](https://github.com/puppeteer/puppeteer)
* [scrapy](https://scrapy.org)
* [Selenium](https://www.selenium.dev)

Long-running database queries and index-management tasks.

VM provisioning / garbage collection.

## Kubernetes Jobs API

* [Jobs API](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
* [Alternatives to Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/job/#alternatives)

## See also

* ["Jobs" runner for Docker Swarm](https://github.com/alexellis/jaas)
* [Stefan's kjob CLI](https://github.com/stefanprodan/kjob)
* [OpenFaaS Research item: Long running batch-jobs / workflow #657](https://github.com/openfaas/faas/issues/657)

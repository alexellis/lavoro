# lavoro (work/job)

An experimental CLI which runs a long-running container, then prints out its logs.

This could be easier than remembering the syntax for `kubectl run` for when you want to debug a service using `curl` - or just to run a long-running batch-job, getting the result at the end.

## Status

lavoro uses the Kubernetes API and Go client to schedule a Kubernetes job.

Future work would include the use of informers and caches for efficiency.

## Examples

### Debug a network issue with curl

```bash
kubectl run nginx --image nginx --port 8080
kubectl expose pod nginx --type ClusterIP --port 8080 --target-port 80

lavoro run --name curl1 \
  --image curlimages/curl:latest \
  --command "curl -s http://nginx:8080"

2020/10/30 16:02:35 Accepted: {curl1 default {curlimages/curl:latest [curl -s http://nginx:8080]}}
curl1-1604073755
2020/10/30 16:02:35 curl1-1604073755
2020/10/30 16:02:36 Pending
2020/10/30 16:02:37 Succeeded
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
...
```

### Ping

```bash
lavoro run --image alpine:3.12 \
  --command "ping -c 4 google.com"

2020/10/30 15:50:14 Accepted: {job default {alpine:3.12 [ping -c 4 google.com]}}
job-1604073014
2020/10/30 15:50:14 job-1604073014
2020/10/30 15:50:15 Pending
2020/10/30 15:50:16 Running
2020/10/30 15:50:17 Running
2020/10/30 15:50:18 Running
2020/10/30 15:50:19 Succeeded
PING google.com (216.58.212.238): 56 data bytes
64 bytes from 216.58.212.238: seq=0 ttl=36 time=16.759 ms
64 bytes from 216.58.212.238: seq=1 ttl=36 time=20.031 ms
64 bytes from 216.58.212.238: seq=2 ttl=36 time=13.693 ms
64 bytes from 216.58.212.238: seq=3 ttl=36 time=14.424 ms

--- google.com ping statistics ---
4 packets transmitted, 4 packets received, 0% packet loss
round-trip min/avg/max = 13.693/16.226/20.031 ms
2020/10/30 15:50:19 Deleting job job-1604073014..OK.
```

### Web-scraping with Pupeteer

This runs integration tests with `mocha` which visit the inlets blog to check it's up.

```bash
lavoro run --name inlets-check \
  --image alexellis2/check-inlets \
  --command "mocha --recursive ./integration-tests"

2020/10/30 15:52:08 Accepted: {inlets-check default {alexellis2/check-inlets [mocha --recursive ./integration-tests]}}
inlets-check-1604073128
2020/10/30 15:52:08 inlets-check-1604073128
2020/10/30 15:52:08 Pending
2020/10/30 15:52:09 Pending
2020/10/30 15:52:10 Running
2020/10/30 15:52:11 Running
2020/10/30 15:52:12 Running
2020/10/30 15:52:13 Succeeded


Started HeadlessChrome/85.0.4182.0
  inlets blog
Promise { <pending> }
    ✓ renders OK (2722ms)


  1 passing (3s)

2020/10/30 15:52:13 Deleting job inlets-check-1604073128..OK.
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

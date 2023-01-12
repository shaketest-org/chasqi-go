# A runtime that performs load or smoke tests for you (or any other traffic tests)

## What is it?

An E2E testing tool similar to k6 (https://k6.io/) only more lightweight as it's written in Go. The main difference is
that it is intended to be
used more as a run-time upon which you can build your tooling, such as your own testing tool with
advanced statistics, real-time monitoring, etc.

# Motivation
I'm building an E2E Testkit for a certain nieche in the industry and this will serve as core of the system. I decided to make it open source, maybe it could be useful to someone else.

## How does it work ?

That runtime exposes a local web server and listens to a port. You then simply post your User Journeys and Configs (e.g.
tolerated failure rate or how many users should be emulated, for Smoke Tests 1 would be enough).

## Still WIP
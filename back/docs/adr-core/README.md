# ADR — Back core (Go service baseline)

These Architecture Decision Records define the **reusable baseline of a Go service**:
folder layout, standard library usage, generated server, errors, timeouts, lint,
configuration, scaffolding. They apply the common conventions of
[`../../../docs/adr-common/`](../../../docs/adr-common/) to Go.

Decisions specific to this back live in [`../adr-functional/`](../adr-functional/); its build
and run in [`../adr-devops/`](../adr-devops/). Read them in order.

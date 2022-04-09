# TraceView

A light-weight and cross-platform trace capturing and viewing tool _for developers_. Receives traces by UDP, and and they are then viewable/searchable using a web interface. Inspired by Sysinternal's [DebugView](https://docs.microsoft.com/en-us/sysinternals/downloads/debugview).

## CURRENTLY WORK IN PROGRESS

There are no releases at the moment but the FAQ below should give you an idea what this project is about.

## FAQ

### In this day and age, tracing is a solved problem. We have robust distributed tracing systems capable of handling many nodes and applications. Why would I ever need this toy tracing tool?

The strength of distributed tracing systems have become their downfall for development scenarios: they are beasts with many knobs to configure and as a result difficult to set up when you need a simple tracing. TraceView does not target production or distributed scenarios and is intended for running mostly on developer's machines while debugging - hence trivial to get started with.

### So now do I have to have two sets of tracing code, one for my production distributed scenario the other for my pet debugging?

No! The idea is that there will be language/framework-dependant adapters that will _tee_ (fork) your existing tracing output and send them to TraceView too. If you already have tracing code, TraceView will piggyback on that using these libraries.

### If this is intended for development environment scenarios only, why persisted storage?

The idea is that once you collect a trace, you might want to review it later. With an in-memory storage, as soon as you stop the docker container all traces will be lost.

### What are typical scenarios where I might need it?

Think of a scenario when you are debugging an issue but typically due to concurrency, it is impossible to walk through the code as threads work left, right and centre and interfere with proper debugging. Or a server that is receiving various requests but you need to trace a particular request. Or when the issue only occurs when you run the code and not while debugging (observing interferes with the behaviour) If you have never had these scenarios then you are lucky but for most non-trivial systems, these cases occur.

### How do I get started [Future, not yet]?

1. Use language/framework-dependant libraries ([supported frameworks](...)) to send output of your existing/added tracing code to TraceView
2. Download the docker image
3. Run the docker image
4. Start debugging

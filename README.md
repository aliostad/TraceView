# TraceView

A light-weight and cross-platform trace capturing and viewing tool _for developers_. Receives traces by UDP, and and they are then viewable/searchable using a web interface. Inspired by Sysinternal's [DebugView](https://docs.microsoft.com/en-us/sysinternals/downloads/debugview).

## FAQ's

### There are so many tracing/logging tools out there. Why adding this?

Because they are targeted for production environments and as a result have many knobs and are pretty difficult to set up and configure. TraceView does not target such scenarios and is intended for running mostly on developer's machines while debugging - hence trivial to get started with.

### What are typical scenarios where I might need it?

Think of a scenario when you are debugging an issue but typically due to concurrency, it is impossible to walk through the code as threads work left, right and centre and interfere with proper debugging. Or a server that is receiving various requests but you need to trace a particular request. Or when the issue only occurs when you run the code and not while debugging (observing interferes with the behaviour) If you have never had these scenarios then you are lucky but for most non-trivial code, these cases occur.

### How do I get started?

1. Add the tracing code to your code using library ([supported frameworks](...))
2. Download the docker image
3. Run the docker image

### What happens if I keep the tracing code in the production?

The idea is that you would! There will be language/framework-related switches to turn it off in production. But the idea is that with choosing UDP, the performance impact of leaving it there should not be significant if not completely negligible.

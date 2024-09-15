# Build Tooling In Java Sucks
*A Manifesto*

Java is a great language. It's been running swaths of the internet for decades. It's stood the test of time,
even if it's acquired some baggage over those 25/30 years. It's fast, it's incredibly secure, it's featureful, and the tooling
is in my opinion the best of any language out there. Spring is largely a joy to write full stack applications
with. There's just one problem, the build tooling. Maven and Gradle are needleslly complex for those who
don't need a 20 foot wide toolbox filled to the brim with Snap On. Most of the time, when I'm working with
Java, it's in Spring and the microservices are deployed largely the same way. An embedded Tomcat server
with the Spring Web Starter. Espresso aims to fill the void of simple build tooling.


## What we're doing different

### We're Opinionated

Espresso is relatively opinionated. We have one good way (at least we think) of how everything should be done.
And that's not changing. If you want a custom feature, fork it! It's MIT licensed!

### Wrapper-less

Every project that uses Espresso ships with the executable of Espresso it was designed to use. There's no
runtimes or wrappers to manage, it's compiled code designed for your system.

### Don't Over Engineer From The Beginning

Developers LOVE to over-engineer solutions from the beginning, building things from the beginning to
support hundreds of thousands of users that 99% of applications will never come close to achieving.
Why should the build system reflect that? I should be able to init a new project, add some endpoints,
build an uber JAR and ship it. Not worry about acquiring 5 different plugins for Spring, Shade/Shadow,
Compiler Plugins, etc. If you build with Espresso, and your application DOES grow to need a build system
with more features (congrats!), switching to Maven or Gradle is possible. You'll probably need to refactor
away some tech debt anyways.

### 
I spend more time fucking with Gradle and friends than I spend actually writing code.

Java is a pretty solid language, it makes programming easy without compromising on safety and speed (too much).
In a previous life, I used Python mainly and found myself implementing libraries like Pydantic to account for
the dynamicness of the language. But this isn't about Python, but Java. I feel like Java's build tooling
has constantly been a pissing contest of who can make the tool that gets the most users, the tool that can
hit the most nails. Espresso isn't that. 

Espresso is aimed at developers who have their main class always named Main. Those who always want an uber JAR, not some
stripped down shit JAR. I shouldn't have to specify in my build tool and editor where my main class is,
what my Gradle project name is, if X dependency is an annotation processor, a regular dependency, and a compile
time dependency. It should just fucking know.

Don't even get me started on the DSL crap and trying to figure out how plugins work in Gradle and Maven,
cause no one could tell you! The variables and arguments change so fucking constantly. Whose bright
idea was it to use Groovy as a configuration language..? With Espresso, we use TOML. We maintain


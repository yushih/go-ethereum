
OUTPUT := blockchain.jar
javas := $(shell find . -name '*.java')
classes := $(patsubst %.java,%.class,$(javas))

%.class: %.java
	javac $<

$(OUTPUT): $(classes)
	jar cf $(OUTPUT) $(classes)


default: $(OUTPUT)

FROM anapsix/alpine-java

MAINTAINER francescofrontera

ENV JAR_TO_EXECUTE fakeJar.jar
ENV MAIN_CLASS fake_main_class

CMD ["sh", "-c", "java -cp /jar/${JAR_TO_EXECUTE} ${MAIN_CLASS}"]
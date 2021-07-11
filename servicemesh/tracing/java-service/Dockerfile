FROM openjdk:8-jdk-slim

ADD target/rest.jar /opt

ENTRYPOINT ["java", "-jar", "/opt/rest.jar"]
FROM openjdk:11
RUN mkdir ./app
COPY ./currencies-0.0.1-SNAPSHOT.jar ./app
ENTRYPOINT ["java","-jar","./app/currencies-0.0.1-SNAPSHOT.jar"]


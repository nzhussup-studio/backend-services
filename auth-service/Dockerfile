FROM maven:3.9.10-eclipse-temurin-21 AS build
LABEL authors="nzhussup"
WORKDIR /app
COPY pom.xml .
RUN mvn dependency:go-offline -B
COPY src ./src
RUN mvn clean package -DskipTests

FROM eclipse-temurin:21-jre AS runtime
WORKDIR /app
COPY --from=build /app/target/*.jar ./app.jar
EXPOSE 8083
ENTRYPOINT ["java", "-jar", "app.jar"]

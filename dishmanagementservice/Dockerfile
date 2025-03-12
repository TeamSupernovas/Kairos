# Use the official OpenJDK 21 slim image as base
FROM openjdk:21-jdk-slim

# Set working directory inside the container
WORKDIR /app


# Copy the built Spring Boot JAR file
COPY target/ dishmanagementservice.jar

# Expose the application port (use default 8080)
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["java", "-jar", "dishmanagementservice.jar"]

#ENTRYPOINT ["sh", "-c", "java -jar dishmanagementservice_1.jar & tail -f /dev/null"]
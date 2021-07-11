#!/bin/sh

# Get home directory
APP_HOME=$(cd $(dirname $0); cd ..; pwd -P)

# Project Information
ARTIFACT_ID=rest-service
VERSION=0.0.1-SNAPSHOT

# Set entry
MAIN_JAR_DIR=lib
MAIN_JAR=$ARTIFACT_ID-$VERSION.jar

# Set log directory for log4j
LOG_DIR=/export/Logs

# Set Java environment
#JAVA_ENV="-Dlog-dir=$LOG_DIR"

# Run
java $JAVA_ENV -jar $APP_HOME/$MAIN_JAR_DIR/$MAIN_JAR
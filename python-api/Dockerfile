# Use an official Python runtime as a parent image
FROM python:3.9-slim

# Set the working directory in the container
WORKDIR /usr/src/app

# Set environment variables
# Prevents Python from writing pyc files to disc (equivalent to python -B option)
ENV PYTHONDONTWRITEBYTECODE 1
# Prevents Python from buffering stdout and stderr (equivalent to python -u option)
ENV PYTHONUNBUFFERED 1

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# Copy the current directory contents into the container at /usr/src/app
COPY . /usr/src/app

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Make port 8082 available to the world outside this container
EXPOSE 8002

# Define environment variable
ENV FLASK_APP=api.py
ENV FLASK_RUN_HOST=127.0.0.1
ENV QUEUE_HOST=rabbitmq

# Run api.py when the container launches
CMD ["flask", "run", "--host=0.0.0.0", "--port=8002"]

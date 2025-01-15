import json
import pika
import os

QUEUE_HOST = os.getenv("QUEUE_HOST")
QUEUE_PORT = os.getenv("QUEUE_PORT")

async def rabbitmq_listener():
    def callback(ch, method, properties, body):
        message = json.loads(body)
        print(f"Received message: {message}")
        # Add custom logic for processing the message

    connection = pika.BlockingConnection(pika.ConnectionParameters(QUEUE_HOST, QUEUE_PORT))
    channel = connection.channel()
    channel.queue_declare(queue=QUEUE_HOST)

    channel.basic_consume(
        queue=QUEUE_HOST,
        on_message_callback=callback,
        auto_ack=True,
    )

    print("Starting RabbitMQ listener...")
    channel.start_consuming()
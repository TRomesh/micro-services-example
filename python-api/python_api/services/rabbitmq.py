import json
import pika

QUEUE_HOST = 'rabbitmq'
QUEUE_PORT = 5672

async def rabbitmq_listener():
    def callback(ch, method, properties, body):
        message = json.loads(body)
        print(f"Received message: {message}")
        # Add custom logic for processing the message

    connection = pika.BlockingConnection(pika.ConnectionParameters(QUEUE_HOST, QUEUE_PORT))
    channel = connection.channel()
    channel.queue_declare(queue="payments_queue")

    channel.basic_consume(
        queue="payments_queue",
        on_message_callback=callback,
        auto_ack=True,
    )

    print("Starting RabbitMQ listener...")
    channel.start_consuming()
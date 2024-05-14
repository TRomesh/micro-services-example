import { connect, Connection, Channel } from "amqplib";
import { QUEUE_HOST, QUEUE_PASSWORD, QUEUE_USER } from "../config";

class RabbitMQ {
  private conn: Connection | null = null;
  private channel: Channel | null = null;

  async connect() {
    this.conn = await connect(
      `amqp://${QUEUE_USER}:${QUEUE_PASSWORD}@${QUEUE_HOST}`
    );
    this.channel = await this.conn.createChannel();
  }

  async sendToQueue(queue: string, message: string) {
    await this.channel?.assertQueue(queue, { durable: false });
    this.channel?.sendToQueue(queue, Buffer.from(message));
    console.log(" [x] Sent '%s'", message);
  }

  async consume(queue: string, callback: (msg: string) => void) {
    await this.channel?.assertQueue(queue, { durable: false });
    this.channel?.consume(
      queue,
      (msg) => {
        if (msg !== null) {
          console.log(" [x] Received '%s'", msg.content.toString());
          callback(msg.content.toString());
          this.channel?.ack(msg);
        }
      },
      { noAck: false }
    );
  }
}

export const rabbitMQ = new RabbitMQ();

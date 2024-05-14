import type { Channel } from "amqplib/callback_api";

export const sendMessageToQueue = (
  channel: Channel,
  queue: string,
  message: string
) => {
  channel.sendToQueue(queue, Buffer.from(message));
  console.log(" [x] Sent %s", message);
};

export const receiveMessageFromQueue = (channel: Channel, queue: string) => {
  channel.consume(
    queue,
    (message) => {
      console.log(" [x] Received %s", message?.content.toString());
    },
    {
      noAck: true,
    }
  );
};

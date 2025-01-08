import app from "./app";
import { APP_PORT, QUEUE_NAME } from "./config";
import { rabbitMQ } from "./services/rabbitMQ";

app.listen(APP_PORT, async () => {
  await rabbitMQ.connect();
  console.log(`🚀 Server is running on port ${APP_PORT}`);
  rabbitMQ.consume(QUEUE_NAME, (msg: string) => {
    console.log(`Consumed message: ${msg}`);
  });
});

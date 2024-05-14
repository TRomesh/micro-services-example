import app from "./app";
import { APP_PORT, USER_QUEUE } from "./config";
import { rabbitMQ } from "./services/rabbitMQ";

app.listen(APP_PORT, async () => {
  await rabbitMQ.connect();
  console.log(`🚀 Server is running on port ${APP_PORT}`);
  rabbitMQ.consume(USER_QUEUE, (msg: string) => {
    console.log(`Consumed message: ${msg}`);
  });
});

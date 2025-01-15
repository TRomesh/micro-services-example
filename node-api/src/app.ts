import bodyParser from "body-parser";
import cors from "cors";
import express, { Request, Response } from "express";
import helmet from "helmet";
import morgan from "morgan";

import { ORIGIN } from "./config";
import { rateLimiter } from "./middleware/rate-limiter";
import routes from "./routes";
import db from "./db";

const app = express();

app.use(morgan("dev"));
app.use(helmet());
app.use(express.json());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(
  cors({
    credentials: true,
    origin: ORIGIN,
    methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
  })
);

app.use(rateLimiter);

app.use("/ping", (_: Request, res: Response) => {
  res.status(200).json({ message: "node-api pong 🏓" });
});

app.locals.db = db;

app.use("/api", routes);

export default app;

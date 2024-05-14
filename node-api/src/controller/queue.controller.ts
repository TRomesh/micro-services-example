import { Request, Response } from "express";
import { rabbitMQ } from "../services/rabbitMQ";

export const sendToQueue = async (req: Request, res: Response) => {
  try {
    const { queue, message } = req.body;
    await rabbitMQ.sendToQueue(queue, message);
    res.status(200).json(`Message sent to queue ${queue}!`);
  } catch (error) {
    res.status(500).json({ error: error });
  }
};

import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { Channel } from "amqplib";
import { Request } from "express";

export const getDBConnection = (
  req: Request
): NodePgDatabase<Record<string, never>> => {
  const db: NodePgDatabase<Record<string, never>> = req.app.locals.db;
  if (!db) {
    throw new Error("No database connection instance found!");
  }
  return db;
};

export const getChannelConnection = (req: Request): Channel => {
  const channel: Channel = req.app.locals.channel;
  if (!channel) {
    throw new Error("No queue connection instance found!");
  }
  return channel;
};

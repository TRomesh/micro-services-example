import { users } from "../db/schema";
import { Request, Response } from "express";
import { eq } from "drizzle-orm";
import { getDBConnection } from "../utils";

export const createUser = async (req: Request, res: Response) => {
  try {
    const db = getDBConnection(req);
    const data = req.body;
    const user = await db.insert(users).values(data).returning({
      fullName: users.fullName,
      userName: users.userName,
      email: users.email,
      phone: users.phone,
    });
    res.status(200).json(user);
  } catch (error) {
    res.status(500).json({ error: error });
  }
};

export const getUsers = async (req: Request, res: Response) => {
  try {
    const db = getDBConnection(req);
    const user = await db.select().from(users);
    res.status(200).json(user);
  } catch (error) {
    res.status(500).json({ error: error });
  }
};

export const getUser = async (req: Request, res: Response) => {
  try {
    const db = getDBConnection(req);
    const id = req.body.id;
    const user = await db.select().from(users).where(eq(users.id, id));
    res.status(200).json(user);
  } catch (error) {
    res.status(500).json({ error: error });
  }
};

export const updateUser = async (req: Request, res: Response) => {
  try {
    const data = req.body;
    const db = getDBConnection(req);
    const user = await db
      .update(users)
      .set(data)
      .where(eq(users.id, data.id))
      .returning({ id: users.id });
    res.status(200).json(user);
  } catch (error) {
    res.status(500).json({ error: error });
  }
};

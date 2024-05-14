import { Router } from "express";
import { createUser, getUser, getUsers } from "../controller/user.controller";
import { sendToQueue } from "../controller/queue.controller";

const router = Router();

router.get("/user", getUser);

router.get("/users", getUsers);

router.post("/user", createUser);

router.post("/send", sendToQueue);

export default router;

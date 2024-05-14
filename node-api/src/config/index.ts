import dotenv from "dotenv";
import { Secret } from "jsonwebtoken";

dotenv.config();

/* App Config */
export const APP_PORT = process.env.APP_PORT || 8080;
export const JWT_SECRET_KEY = process.env.JWT_SECRET_KEY as Secret;
export const ORIGIN = process.env.ORIGIN || "*";

/* DB Config */
export const DB_HOST = process.env.DB_HOST || "localhost";
export const DB_PORT = Number(process.env.DB_PORT);
export const DB_USER = process.env.DB_USER;
export const DB_PASSWORD = process.env.DB_PASSWORD;
export const DB_DATABASE = process.env.DB_DATABASE || "main";
export const QUEUE_HOST = process.env.QUEUE_HOST || "localhost";
export const QUEUE_USER = process.env.QUEUE_USER || "guest";
export const QUEUE_PASSWORD = process.env.QUEUE_PASSWORD || "guest";
export const USER_QUEUE = process.env.USER_QUEUE || "DEFAULT_QUEUE";

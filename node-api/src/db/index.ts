import { drizzle } from "drizzle-orm/node-postgres";
import { Pool } from "pg";

import { DB_DATABASE, DB_HOST, DB_PASSWORD, DB_PORT, DB_USER } from "../config";
import schema from "./schema";

const client = new Pool({
  host: DB_HOST,
  port: DB_PORT,
  user: DB_USER,
  password: DB_PASSWORD,
  database: DB_DATABASE,
});

const db = drizzle(client, { schema });

export default db;

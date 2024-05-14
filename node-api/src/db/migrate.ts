import { migrate } from "drizzle-orm/node-postgres/migrator";
import db from ".";

// this will automatically run needed migrations on the database
// "generate-migration": "npx drizzle-kit generate:pg --out src/db/migrations --schema src/db/schema.ts"
migrate(db, { migrationsFolder: "./drizzle" })
  .then(() => {
    console.log("Migrations complete!");
    process.exit(0);
  })
  .catch((err) => {
    console.error("Migrations failed!", err);
    process.exit(1);
  });

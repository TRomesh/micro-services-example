import { pgTable, serial, text, timestamp, varchar } from "drizzle-orm/pg-core";

export const users = pgTable("users", {
  id: serial("id").primaryKey(),
  fullName: text("full_name"),
  userName: text("user_name").unique().notNull(),
  email: text("email").unique().notNull(),
  password: text("password"),
  phone: varchar("phone", { length: 256 }),
  createdAt: timestamp("created_at").notNull().defaultNow(),
});

export type UserType = typeof users.$inferSelect; // return type when queried
export type NewUserType = typeof users.$inferInsert; // insert type

export default { users };

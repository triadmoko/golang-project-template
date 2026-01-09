---
name: database-migrations
description: Create and manage database migrations using golang-migrate. Use when creating tables, modifying schemas, adding columns, or managing database structure changes.
allowed-tools: Read, Write, Bash, Glob
---

# Database Migrations

This skill guides you through creating and managing database migrations using golang-migrate following this project's specific patterns.

## ⚠️ IMPORTANT PROJECT PATTERNS

### Table Structure Pattern
```sql
-- Standard column order for ALL tables:
CREATE TABLE IF NOT EXISTS table_name (
    id VARCHAR(36) PRIMARY KEY,           -- ID first, VARCHAR(36) not UUID!
    created_at TIMESTAMP NULL,            -- Timestamp fields
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    -- Then other columns
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    -- etc...
);
```

### Critical Rules
- ✅ Use `VARCHAR(36)` for IDs (NOT `UUID` type!)
- ✅ ID always PRIMARY KEY, always first column
- ✅ Timestamps: `created_at`, `updated_at`, `deleted_at` - all `TIMESTAMP NULL`
- ❌ NO Foreign Keys (FK)
- ❌ NO References
- ❌ NO Constraints (except PRIMARY KEY and UNIQUE)
- ❌ NO CASCADE operations
- ✅ Use indexes for relationships instead of FKs

## Quick Reference

```bash
make migration-create name=xxx    # Create new migration
make migration-up                 # Apply all pending migrations
make migration-down               # Rollback last migration
make migration-version            # Show current version
make migration-force version=N    # Force to specific version
```

## Creating New Migrations

### 1. Create Migration Files

```bash
make migration-create name=create_users_table
```

This creates two files in `migration/`:
- `000001_create_users_table.up.sql` - Applied when migrating up
- `000001_create_users_table.down.sql` - Applied when migrating down

### 2. Write the UP Migration

The `.up.sql` file should create or modify database structures:

```sql
-- migration/000001_create_users_table.up.sql

-- Standard table structure with correct column order
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    email VARCHAR(255) NOT NULL,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NULL,
    status VARCHAR(50) NULL DEFAULT 'active',
    birth_date DATE NULL,
    gender VARCHAR(10) NULL,
    role VARCHAR(50) NULL DEFAULT 'user',
    provider VARCHAR(50) NULL,
    is_active BOOLEAN NULL DEFAULT TRUE
);

-- Create UNIQUE indexes (NO UNIQUE CONSTRAINT!)
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);

-- Regular indexes for frequently queried columns
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

**IMPORTANT:**
- ✅ `VARCHAR(36)` for ID
- ✅ ID, timestamps first
- ✅ Use `UNIQUE INDEX` instead of `UNIQUE CONSTRAINT`
- ✅ All timestamps are `TIMESTAMP NULL`
- ❌ NO `DEFAULT gen_random_uuid()` - handled by application
- ❌ NO `REFERENCES` or `FOREIGN KEY`

### 3. Write the DOWN Migration

The `.down.sql` file must **reverse** the UP migration:

```sql
-- migration/000001_create_users_table.down.sql

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

**Important**: DOWN migrations should always use `IF EXISTS` to avoid errors.

## Common Migration Patterns

### Adding a New Table

**UP:**
```sql
-- Products table with relationship to users (NO FK!)
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    price DECIMAL(10, 2) NOT NULL,
    user_id VARCHAR(36) NOT NULL,  -- NO FOREIGN KEY!
    category_id VARCHAR(36) NULL,
    status VARCHAR(50) NULL DEFAULT 'active'
);

-- Index on user_id for queries (replaces FK)
CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
```

**DOWN:**
```sql
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_user_id;
DROP TABLE IF EXISTS products;
```

### Adding a Column

**UP:**
```sql
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(500) NULL;
ALTER TABLE users ADD COLUMN bio TEXT NULL;
```

**DOWN:**
```sql
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
```

### Adding a UNIQUE Index

**UP:**
```sql
-- Use UNIQUE INDEX, NOT UNIQUE CONSTRAINT
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
```

**DOWN:**
```sql
DROP INDEX IF EXISTS idx_users_phone;
```

### Renaming a Column

**UP:**
```sql
ALTER TABLE users RENAME COLUMN full_name TO display_name;
```

**DOWN:**
```sql
ALTER TABLE users RENAME COLUMN display_name TO full_name;
```

### Adding an Enum Type

**UP:**
```sql
-- Create enum type
CREATE TYPE user_role AS ENUM ('admin', 'user', 'guest');

-- Add column using enum
ALTER TABLE users ADD COLUMN role user_role DEFAULT 'user';
```

**DOWN:**
```sql
ALTER TABLE users DROP COLUMN IF EXISTS role;
DROP TYPE IF EXISTS user_role;
```

### Modifying Column Type

**UP:**
```sql
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(320);
```

**DOWN:**
```sql
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(255);
```

### Adding Default Value

**UP:**
```sql
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
```

**DOWN:**
```sql
ALTER TABLE users ALTER COLUMN status DROP DEFAULT;
```

### Creating Junction Table (Many-to-Many) - NO FK!

**UP:**
```sql
-- Junction table for user-role relationship
CREATE TABLE IF NOT EXISTS user_roles (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    user_id VARCHAR(36) NOT NULL,     -- NO FOREIGN KEY!
    role_id VARCHAR(36) NOT NULL,     -- NO FOREIGN KEY!
    assigned_at TIMESTAMP NULL
);

-- Composite unique index to prevent duplicates
CREATE UNIQUE INDEX idx_user_roles_user_role ON user_roles(user_id, role_id);

-- Individual indexes for queries
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_deleted_at ON user_roles(deleted_at);
```

**DOWN:**
```sql
DROP INDEX IF EXISTS idx_user_roles_deleted_at;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_user_role;
DROP TABLE IF EXISTS user_roles;
```

### Table with Self-Reference (NO FK!)

**UP:**
```sql
-- Categories with parent-child relationship (NO FK!)
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(36) NULL,  -- NO FOREIGN KEY to self!
    slug VARCHAR(255) NOT NULL,
    description TEXT NULL
);

-- Index on parent_id for hierarchy queries
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE UNIQUE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
```

**DOWN:**
```sql
DROP INDEX IF EXISTS idx_categories_deleted_at;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP TABLE IF EXISTS categories;
```

## Complete Example: Orders and Order Items

**UP Migration:**
```sql
-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    user_id VARCHAR(36) NOT NULL,     -- NO FK!
    order_number VARCHAR(50) NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NULL DEFAULT 'pending',
    notes TEXT NULL
);

CREATE UNIQUE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(36) PRIMARY KEY,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    order_id VARCHAR(36) NOT NULL,    -- NO FK!
    product_id VARCHAR(36) NOT NULL,  -- NO FK!
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_deleted_at ON order_items(deleted_at);
```

**DOWN Migration:**
```sql
DROP INDEX IF EXISTS idx_order_items_deleted_at;
DROP INDEX IF EXISTS idx_order_items_product_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP TABLE IF EXISTS order_items;

DROP INDEX IF EXISTS idx_orders_deleted_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_order_number;
DROP TABLE IF EXISTS orders;
```

## Running Migrations

### Apply All Pending Migrations

```bash
make migration-up
```

### Apply Specific Number of Migrations

```bash
make migration-up version=2  # Apply next 2 migrations
```

### Rollback Last Migration

```bash
make migration-down
```

### Rollback Specific Number of Migrations

```bash
make migration-down version=2  # Rollback last 2 migrations
```

### Check Current Version

```bash
make migration-version
```

### Force to Specific Version (Recovery)

If migrations are in a bad state:

```bash
make migration-force version=3  # Force to version 3
```

**Warning**: Only use `migration-force` for recovery. It doesn't run migrations.

## Best Practices

### 1. Always Test Migrations

Test both UP and DOWN migrations:

```bash
# Apply migration
make migration-up

# Verify database state
# Check tables, columns, indexes

# Test rollback
make migration-down

# Verify clean rollback
# Confirm everything was removed
```

### 2. Make Migrations Idempotent

Always use `IF EXISTS` and `IF NOT EXISTS`:

```sql
-- Good
CREATE TABLE IF NOT EXISTS users (...);
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_email;

-- Bad (will fail if already exists/doesn't exist)
CREATE TABLE users (...);
DROP TABLE users;
DROP INDEX idx_users_email;
```

### 3. Never Edit Applied Migrations

Once a migration is applied in production:
- **Never modify it**
- Create a new migration to make changes
- This ensures version history is consistent

### 4. Use Indexes Instead of Foreign Keys

Since we don't use FK:
- Always add indexes on columns used in joins
- Index columns that reference other tables
- This maintains query performance

```sql
-- Product references user
user_id VARCHAR(36) NOT NULL,

-- Add index for queries
CREATE INDEX idx_products_user_id ON products(user_id);
```

### 5. Standard Index Naming

Follow naming convention:
- `idx_{table}_{column}` - Regular index
- `idx_{table}_{col1}_{col2}` - Composite index

### 6. Always Index deleted_at

For soft delete queries:

```sql
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

### 7. Document Complex Migrations

Add comments to explain non-obvious changes:

```sql
-- This migration adds soft delete support
-- We keep deleted records for audit purposes
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

## Column Order Template

Always follow this order:

```sql
CREATE TABLE IF NOT EXISTS table_name (
    -- 1. ID (always first)
    id VARCHAR(36) PRIMARY KEY,

    -- 2. Timestamp columns (always these three)
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,

    -- 3. Required business columns
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,

    -- 4. Reference columns (NO FK!)
    user_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NULL,

    -- 5. Optional business columns
    description TEXT NULL,
    status VARCHAR(50) NULL DEFAULT 'active',

    -- 6. Boolean flags
    is_active BOOLEAN NULL DEFAULT TRUE,
    is_verified BOOLEAN NULL DEFAULT FALSE
);
```

## Troubleshooting

### Migration Failed Midway

Check current version:
```bash
make migration-version
```

If stuck in "dirty" state, force to last good version:
```bash
make migration-force version=N
```

Then fix the problematic migration and retry.

### Can't Connect to Database

Check `.env` file for correct database credentials:
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASS`
- `DB_NAME`

### Permission Errors

Ensure database user has sufficient privileges:
```sql
GRANT ALL PRIVILEGES ON DATABASE db_name TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
```

### Index Already Exists

If you see "index already exists" error:
- Use `IF NOT EXISTS` in CREATE INDEX (PostgreSQL 9.5+)
- Or use `DROP INDEX IF EXISTS` before creating

## Integration with GORM

After creating migrations, update GORM entities in `internal/shared/domain/entity/`:

```go
type User struct {
    ID        string         `json:"id" gorm:"type:varchar(36);primaryKey"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    Email     string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Username  string `json:"username" gorm:"type:varchar(100);uniqueIndex;not null"`
    Password  string `json:"-" gorm:"type:varchar(255);not null"`
}
```

**Important**:
- GORM struct tags should match your migration schema
- Migrations are the source of truth
- Use migrations, NOT GORM AutoMigrate
- ID is `string` type, not `uuid.UUID`

## Migration Checklist

Before creating migration:
- [ ] ID is VARCHAR(36) PRIMARY KEY (first column)
- [ ] Timestamps: created_at, updated_at, deleted_at (TIMESTAMP NULL)
- [ ] NO Foreign Keys, References, Constraints (except PK, UNIQUE)
- [ ] Use indexes for relationships
- [ ] UNIQUE INDEX instead of UNIQUE CONSTRAINT
- [ ] Index on deleted_at for soft deletes
- [ ] IF EXISTS / IF NOT EXISTS for idempotency
- [ ] DOWN migration reverses UP completely
- [ ] Tested UP and DOWN before committing

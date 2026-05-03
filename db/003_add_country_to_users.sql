-- Migration: Add country column to users table
ALTER TABLE users ADD COLUMN country VARCHAR(100);

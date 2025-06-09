-- OpenAPI Demo Database Schema
-- PostgreSQL database schema for Todo and Category management

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing tables if they exist (for clean setup)
DROP TABLE IF EXISTS todos CASCADE;
DROP TABLE IF EXISTS categories CASCADE;

-- Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL CHECK (LENGTH(name) >= 1),
    description VARCHAR(255),
    color VARCHAR(7) CHECK (color ~ '^#[0-9A-Fa-f]{6}$') DEFAULT '#6c757d',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Todos table
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for better query performance
CREATE INDEX idx_todos_category_id ON todos(category_id);
CREATE INDEX idx_todos_completed ON todos(completed);
CREATE INDEX idx_todos_created_at ON todos(created_at);
CREATE INDEX idx_categories_name ON categories(name);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for categories table
CREATE TRIGGER update_categories_updated_at 
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for todos table
CREATE TRIGGER update_todos_updated_at 
    BEFORE UPDATE ON todos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data (optional)
-- Uncomment the following lines to insert sample data

/*
-- Sample categories
INSERT INTO categories (name, description, color) VALUES
    ('仕事', '仕事関連のタスク', '#007bff'),
    ('プライベート', '個人的なタスク', '#28a745'),
    ('買い物', '買い物リスト', '#ffc107'),
    ('学習', '勉強・学習関連', '#6f42c1');

-- Sample todos
INSERT INTO todos (title, description, completed, category_id) VALUES
    ('プロジェクト計画書の作成', '次四半期のプロジェクト計画を立てる', false, (SELECT id FROM categories WHERE name = '仕事')),
    ('食材の買い出し', '今週末の食材を購入する', false, (SELECT id FROM categories WHERE name = '買い物')),
    ('Go言語の学習', 'Go言語のチュートリアルを完了する', true, (SELECT id FROM categories WHERE name = '学習'));
*/
-- Seed data for OpenAPI Demo database
-- Sample categories and todos for testing and development

-- Insert sample categories
INSERT INTO categories (name, description, color) VALUES
    ('仕事', '仕事関連のタスクとプロジェクト', '#007bff'),
    ('プライベート', '個人的なタスクと予定', '#28a745'),
    ('買い物', '買い物リストとお使い', '#ffc107'),
    ('学習', '勉強・学習・スキルアップ関連', '#6f42c1'),
    ('健康', '運動・健康管理関連', '#fd7e14'),
    ('家事', '家庭での作業やメンテナンス', '#20c997');

-- Insert sample todos
INSERT INTO todos (title, description, completed, category_id) VALUES
    -- 仕事関連
    ('プロジェクト計画書の作成', '次四半期のプロジェクト計画を立てて、ステークホルダーに共有する', false, (SELECT id FROM categories WHERE name = '仕事')),
    ('週次レポートの提出', '今週の進捗をまとめて上司に提出する', true, (SELECT id FROM categories WHERE name = '仕事')),
    ('クライアント会議の準備', '来週のクライアント会議用の資料を準備する', false, (SELECT id FROM categories WHERE name = '仕事')),
    
    -- プライベート
    ('友人との食事の予約', '今度の週末に友人と食事する店を予約する', false, (SELECT id FROM categories WHERE name = 'プライベート')),
    ('映画チケットの購入', '今度公開される映画のチケットを事前購入する', true, (SELECT id FROM categories WHERE name = 'プライベート')),
    
    -- 買い物
    ('食材の買い出し', '今週末の食材を購入する（野菜、肉、調味料）', false, (SELECT id FROM categories WHERE name = '買い物')),
    ('新しいヘッドフォンの購入', 'ワイヤレスヘッドフォンを比較検討して購入する', false, (SELECT id FROM categories WHERE name = '買い物')),
    ('文房具の補充', 'ペン、ノート、付箋などを補充する', true, (SELECT id FROM categories WHERE name = '買い物')),
    
    -- 学習
    ('Go言語の学習', 'Go言語のチュートリアルを完了し、サンプルアプリを作成する', true, (SELECT id FROM categories WHERE name = '学習')),
    ('PostgreSQLの復習', 'データベース設計とクエリ最適化について学習する', false, (SELECT id FROM categories WHERE name = '学習')),
    ('英語の勉強', '毎日30分の英語学習を継続する', false, (SELECT id FROM categories WHERE name = '学習')),
    
    -- 健康
    ('ジムでの筋トレ', '週3回のジム通いを継続する', false, (SELECT id FROM categories WHERE name = '健康')),
    ('定期健康診断の予約', '年次健康診断の予約を取る', false, (SELECT id FROM categories WHERE name = '健康')),
    
    -- 家事
    ('部屋の大掃除', 'クローゼットの整理と不要なものの処分', false, (SELECT id FROM categories WHERE name = '家事')),
    ('エアコンのフィルター掃除', 'エアコンのフィルターを掃除して交換する', true, (SELECT id FROM categories WHERE name = '家事'));
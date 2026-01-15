-- 测试种子数据
-- 用于 E2E 测试的测试账号和基础数据

-- 清理旧的测试数据
DELETE FROM users WHERE username IN ('test_user', 'committee_user', 'student_user', 'admin_user');

-- 1. 创建普通学生用户
-- 密码: test123456 (bcrypt hash)
INSERT INTO users (username, password_hash, email, real_name, role, status, created_at, updated_at)
VALUES (
    'test_user',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'test@example.com',
    '测试用户',
    'student',
    'active',
    NOW(),
    NOW()
);

-- 2. 创建学委用户（有上传权限）
-- 密码: test123456
INSERT INTO users (username, password_hash, email, real_name, role, status, created_at, updated_at)
VALUES (
    'committee_user',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'committee@example.com',
    '学委用户',
    'committee',
    'active',
    NOW(),
    NOW()
);

-- 3. 创建另一个学委用户（用于测试）
-- 密码: test123456
INSERT INTO users (username, password_hash, email, real_name, role, status, created_at, updated_at)
VALUES (
    'student_user',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student2@example.com',
    '学生用户2',
    'student',
    'active',
    NOW(),
    NOW()
);

-- 4. 创建管理员用户
-- 密码: admin123456
INSERT INTO users (username, password_hash, email, real_name, role, status, created_at, updated_at)
VALUES (
    'admin_user',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'admin@example.com',
    '管理员',
    'admin',
    'active',
    NOW(),
    NOW()
);

-- 获取创建的用户ID（用于后续插入）
DO $$
DECLARE
    v_test_user_id INT;
    v_committee_user_id INT;
    v_admin_user_id INT;
BEGIN
    SELECT id INTO v_test_user_id FROM users WHERE username = 'test_user';
    SELECT id INTO v_committee_user_id FROM users WHERE username = 'committee_user';
    SELECT id INTO v_admin_user_id FROM users WHERE username = 'admin_user';

    -- 插入一些测试资料
    INSERT INTO materials (title, description, category, course_name, uploader_id, file_name, file_size, file_key, status, created_at, updated_at)
    VALUES
        ('测试资料1', '这是一个测试资料', 'note', '软件工程', v_committee_user_id, 'test1.pdf', 1024000, 'test/test1.pdf', 'approved', NOW(), NOW()),
        ('测试资料2', '这是另一个测试资料', 'exam', '数据结构', v_committee_user_id, 'test2.pdf', 2048000, 'test/test2.pdf', 'approved', NOW(), NOW()),
        ('测试资料3', '待审核的测试资料', 'note', '计算机网络', v_committee_user_id, 'test3.pdf', 512000, 'test/test3.pdf', 'pending', NOW(), NOW());

    -- 添加一些下载记录
    INSERT INTO download_records (user_id, material_id, created_at)
    VALUES
        (v_test_user_id, (SELECT id FROM materials WHERE title = '测试资料1' LIMIT 1), NOW()),
        (v_test_user_id, (SELECT id FROM materials WHERE title = '测试资料2' LIMIT 1), NOW());

    -- 添加一些收藏
    INSERT INTO favorites (user_id, material_id, created_at)
    VALUES
        (v_test_user_id, (SELECT id FROM materials WHERE title = '测试资料1' LIMIT 1), NOW());

    RAISE NOTICE 'Test seed data created successfully';
END $$;

-- 显示创建的用户
SELECT username, email, role, status FROM users WHERE username IN ('test_user', 'committee_user', 'student_user', 'admin_user');

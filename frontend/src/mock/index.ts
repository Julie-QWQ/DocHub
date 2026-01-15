/**
 * Mock 数据
 * 用于开发环境的前端调试
 */

// 用户 Mock 数据
export const mockUsers = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@study-upc.com',
    role: 'admin',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z'
  },
  {
    id: 2,
    username: 'student1',
    email: 'student1@study-upc.com',
    role: 'student',
    created_at: '2026-01-02T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z'
  },
  {
    id: 3,
    username: 'committee1',
    email: 'committee1@study-upc.com',
    role: 'committee',
    created_at: '2026-01-03T00:00:00Z',
    updated_at: '2026-01-03T00:00:00Z'
  }
]

// 资料分类 Mock 数据
export const mockCategories = [
  { label: '课件资料', value: 'courseware' },
  { label: '实验指导', value: 'experiment' },
  { label: '试卷习题', value: 'exam' },
  { label: '参考文献', value: 'reference' },
  { label: '其他资料', value: 'other' }
]

// 资料 Mock 数据
export const mockMaterials = [
  {
    id: 1,
    title: '高等数学上册课件',
    description: '包含第一章到第五章的完整课件内容',
    category: 'courseware',
    status: 'approved',
    file_name: '高等数学上册课件.pptx',
    file_size: 15728640,
    file_url: 'https://example.com/files/1',
    uploader_id: 3,
    uploader_name: 'committee1',
    download_count: 128,
    view_count: 456,
    tags: ['高等数学', '课件', '上册'],
    created_at: '2026-01-05T10:30:00Z',
    updated_at: '2026-01-05T10:30:00Z'
  },
  {
    id: 2,
    title: '大学物理实验指导书',
    description: '2024-2025学年第一学期物理实验指导',
    category: 'experiment',
    status: 'approved',
    file_name: '大学物理实验指导书.pdf',
    file_size: 5242880,
    file_url: 'https://example.com/files/2',
    uploader_id: 3,
    uploader_name: 'committee1',
    download_count: 89,
    view_count: 234,
    tags: ['物理', '实验', '指导书'],
    created_at: '2026-01-06T14:20:00Z',
    updated_at: '2026-01-06T14:20:00Z'
  },
  {
    id: 3,
    title: 'C语言程序设计期末试卷',
    description: '2024年秋季学期期末考试试卷及答案',
    category: 'exam',
    status: 'approved',
    file_name: 'C语言程序设计期末试卷.pdf',
    file_size: 3145728,
    file_url: 'https://example.com/files/3',
    uploader_id: 3,
    uploader_name: 'committee1',
    download_count: 256,
    view_count: 678,
    tags: ['C语言', '期末试卷', '答案'],
    created_at: '2026-01-07T09:15:00Z',
    updated_at: '2026-01-07T09:15:00Z'
  },
  {
    id: 4,
    title: '数据结构与算法参考文献',
    description: '经典算法书籍推荐和参考资料',
    category: 'reference',
    status: 'pending',
    file_name: '数据结构参考文献.pdf',
    file_size: 2097152,
    file_url: 'https://example.com/files/4',
    uploader_id: 3,
    uploader_name: 'committee1',
    download_count: 0,
    view_count: 12,
    tags: ['数据结构', '算法', '参考'],
    created_at: '2026-01-08T16:45:00Z',
    updated_at: '2026-01-08T16:45:00Z'
  },
  {
    id: 5,
    title: '计算机网络复习提纲',
    description: '期末考试重点知识整理',
    category: 'other',
    status: 'approved',
    file_name: '计算机网络复习提纲.docx',
    file_size: 1048576,
    file_url: 'https://example.com/files/5',
    uploader_id: 3,
    uploader_name: 'committee1',
    download_count: 178,
    view_count: 345,
    tags: ['计算机网络', '复习', '提纲'],
    created_at: '2026-01-09T11:00:00Z',
    updated_at: '2026-01-09T11:00:00Z'
  }
]

// 通知 Mock 数据
export const mockNotifications = [
  {
    id: 1,
    user_id: 2,
    type: 'material_approved',
    title: '资料审核通过',
    content: '您上传的"高等数学上册课件"已通过审核',
    is_read: false,
    created_at: '2026-01-10T10:00:00Z'
  },
  {
    id: 2,
    user_id: 2,
    type: 'system_notice',
    title: '系统维护通知',
    content: '系统将于今晚22:00-23:00进行维护',
    is_read: true,
    created_at: '2026-01-10T09:00:00Z'
  },
  {
    id: 3,
    user_id: 2,
    type: 'download_reminder',
    title: '下载提醒',
    content: '您下载的"大学物理实验指导书"有更新版本',
    is_read: false,
    created_at: '2026-01-10T08:30:00Z'
  }
]

// 下载记录 Mock 数据
export const mockDownloadRecords = [
  {
    id: 1,
    user_id: 2,
    material_id: 1,
    material_title: '高等数学上册课件',
    created_at: '2026-01-10T14:30:00Z'
  },
  {
    id: 2,
    user_id: 2,
    material_id: 2,
    material_title: '大学物理实验指导书',
    created_at: '2026-01-10T13:20:00Z'
  },
  {
    id: 3,
    user_id: 2,
    material_id: 3,
    material_title: 'C语言程序设计期末试卷',
    created_at: '2026-01-10T12:10:00Z'
  }
]

// 收藏 Mock 数据
export const mockFavorites = [
  {
    id: 1,
    user_id: 2,
    material_id: 1,
    material: mockMaterials[0],
    created_at: '2026-01-09T10:00:00Z'
  },
  {
    id: 2,
    user_id: 2,
    material_id: 3,
    material: mockMaterials[2],
    created_at: '2026-01-09T11:30:00Z'
  }
]

// 统计数据 Mock
export const mockStatistics = {
  overview: {
    total_users: 1250,
    total_materials: 3568,
    total_downloads: 45678,
    today_active: 234
  },
  materials: {
    by_category: [
      { name: '课件资料', value: 1234 },
      { name: '实验指导', value: 856 },
      { name: '试卷习题', value: 678 },
      { name: '参考文献', value: 456 },
      { name: '其他资料', value: 344 }
    ],
    by_status: [
      { name: '已通过', value: 3124 },
      { name: '待审核', value: 256 },
      { name: '已拒绝', value: 188 }
    ],
    upload_trend: [
      { date: '2026-01-04', count: 45 },
      { date: '2026-01-05', count: 56 },
      { date: '2026-01-06', count: 67 },
      { date: '2026-01-07', count: 78 },
      { date: '2026-01-08', count: 89 },
      { date: '2026-01-09', count: 92 },
      { date: '2026-01-10', count: 102 }
    ]
  },
  downloads: {
    daily_trend: [
      { date: '2026-01-04', count: 1234 },
      { date: '2026-01-05', count: 1456 },
      { date: '2026-01-06', count: 1678 },
      { date: '2026-01-07', count: 1890 },
      { date: '2026-01-08', count: 2134 },
      { date: '2026-01-09', count: 2345 },
      { date: '2026-01-10', count: 2567 }
    ],
    top_materials: [
      { id: 3, title: 'C语言程序设计期末试卷', count: 256 },
      { id: 1, title: '高等数学上册课件', count: 128 },
      { id: 5, title: '计算机网络复习提纲', count: 178 }
    ]
  }
}

// 分页响应
export function mockPaginate<T>(list: T[], page: number, size: number) {
  const start = (page - 1) * size
  const end = start + size
  const data = list.slice(start, end)

  return {
    code: 0,
    message: 'success',
    data: {
      total: list.length,
      page,
      size,
      list: data
    }
  }
}

// 成功响应
export function mockSuccess<T>(data: T) {
  return {
    code: 0,
    message: 'success',
    data
  }
}

// 失败响应
export function mockFail(code: number, message: string) {
  return {
    code,
    message,
    data: null
  }
}

<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon users">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ overviewStats?.users.total || 0 }}</div>
          <div class="stat-label">总用户数</div>
          <div class="stat-trend">+{{ overviewStats?.users.today || 0 }} 今日新增</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon materials">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ overviewStats?.materials.total || 0 }}</div>
          <div class="stat-label">总资料数</div>
          <div class="stat-trend">+{{ overviewStats?.materials.today || 0 }} 今日新增</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon downloads">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ overviewStats?.downloads.total || 0 }}</div>
          <div class="stat-label">总下载次数</div>
          <div class="stat-trend">{{ overviewStats?.downloads.today || 0 }} 今日下载</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon visits">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
            <circle cx="12" cy="12" r="3"></circle>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ overviewStats?.visits.total || 0 }}</div>
          <div class="stat-label">总访问次数</div>
          <div class="stat-trend">{{ overviewStats?.visits.today || 0 }} 今日访问</div>
        </div>
      </div>
    </div>

    <!-- 待处理事项 -->
    <div class="pending-section">
      <h2 class="section-title">待处理事项</h2>
      <div class="pending-cards">
        <div class="pending-card" @click="$router.push('/admin/materials')">
          <div class="pending-left">
            <div class="pending-badge material">资料</div>
            <div class="pending-info">
              <div class="pending-count">{{ overviewStats?.materials.pending || 0 }}</div>
              <div class="pending-text">待审核资料</div>
            </div>
          </div>
          <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </div>

        <div class="pending-card" @click="$router.push('/admin/applications')">
          <div class="pending-left">
            <div class="pending-badge application">申请</div>
            <div class="pending-info">
              <div class="pending-count">{{ overviewStats?.applications.pending || 0 }}</div>
              <div class="pending-text">待审核申请</div>
            </div>
          </div>
          <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </div>

        <div class="pending-card" @click="$router.push('/admin/reports')">
          <div class="pending-left">
            <div class="pending-badge report">举报</div>
            <div class="pending-info">
              <div class="pending-count">{{ overviewStats?.applications.pending || 0 }}</div>
              <div class="pending-text">待处理举报</div>
            </div>
          </div>
          <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <div class="chart-card large">
        <div class="chart-header">
          <h3>用户增长</h3>
          <div class="time-selector">
            <button
              v-for="days in [7, 30, 90]"
              :key="days"
              :class="['time-btn', { active: userTrendDays === days }]"
              @click="userTrendDays = days; loadUserTrend()"
            >
              {{ days }}天
            </button>
          </div>
        </div>
        <div ref="userTrendChartRef" class="chart-body"></div>
      </div>

      <div class="chart-card large">
        <div class="chart-header">
          <h3>资料上传</h3>
          <div class="time-selector">
            <button
              v-for="days in [7, 30, 90]"
              :key="days"
              :class="['time-btn', { active: materialTrendDays === days }]"
              @click="materialTrendDays = days; loadMaterialTrend()"
            >
              {{ days }}天
            </button>
          </div>
        </div>
        <div ref="materialTrendChartRef" class="chart-body"></div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>下载趋势</h3>
          <div class="time-selector">
            <button
              v-for="days in [7, 30, 90]"
              :key="days"
              :class="['time-btn', { active: downloadTrendDays === days }]"
              @click="downloadTrendDays = days; loadDownloadTrend()"
            >
              {{ days }}天
            </button>
          </div>
        </div>
        <div ref="downloadTrendChartRef" class="chart-body"></div>
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <h3>访问趋势</h3>
          <div class="time-selector">
            <button
              v-for="days in [7, 30, 90]"
              :key="days"
              :class="['time-btn', { active: visitTrendDays === days }]"
              @click="visitTrendDays = days; loadVisitTrend()"
            >
              {{ days }}天
            </button>
          </div>
        </div>
        <div ref="visitTrendChartRef" class="chart-body"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import {
  getOverviewStatistics,
  getUserTrend,
  getMaterialTrend,
  getDownloadTrend,
  getVisitTrend
} from '@/api/statistics'
import type { OverviewStatistics, TrendData } from '@/api/statistics'

const overviewStats = ref<OverviewStatistics>()
const userTrendDays = ref(30)
const materialTrendDays = ref(30)
const downloadTrendDays = ref(30)
const visitTrendDays = ref(30)

// 图表实例
const userTrendChartRef = ref<HTMLElement>()
const materialTrendChartRef = ref<HTMLElement>()
const downloadTrendChartRef = ref<HTMLElement>()
const visitTrendChartRef = ref<HTMLElement>()

let userTrendChart: ECharts | null = null
let materialTrendChart: ECharts | null = null
let downloadTrendChart: ECharts | null = null
let visitTrendChart: ECharts | null = null

// 自动刷新定时器
let refreshTimer: ReturnType<typeof setInterval> | null = null

// 加载概览统计
async function loadOverviewStatistics() {
  try {
    const { data } = await getOverviewStatistics()
    overviewStats.value = data
  } catch (error: any) {
    ElMessage.error(error.message || '获取统计数据失败')
  }
}

// 加载用户趋势
async function loadUserTrend() {
  try {
    const { data } = await getUserTrend(userTrendDays.value)
    renderUserTrendChart(data)
  } catch (error: any) {
    ElMessage.error(error.message || '获取用户趋势失败')
  }
}

// 加载资料趋势
async function loadMaterialTrend() {
  try {
    const { data } = await getMaterialTrend(materialTrendDays.value)
    renderMaterialTrendChart(data)
  } catch (error: any) {
    ElMessage.error(error.message || '获取资料趋势失败')
  }
}

// 加载下载趋势
async function loadDownloadTrend() {
  try {
    const { data } = await getDownloadTrend(downloadTrendDays.value)
    renderDownloadTrendChart(data)
  } catch (error: any) {
    ElMessage.error(error.message || '获取下载趋势失败')
  }
}

// 加载访问趋势
async function loadVisitTrend() {
  try {
    const { data } = await getVisitTrend(visitTrendDays.value)
    renderVisitTrendChart(data)
  } catch (error: any) {
    ElMessage.error(error.message || '获取访问趋势失败')
  }
}

// 渲染用户趋势图表
function renderUserTrendChart(data: TrendData[]) {
  if (!userTrendChartRef.value) return

  if (!userTrendChart) {
    userTrendChart = echarts.init(userTrendChartRef.value)
  }

  if (!data || data.length === 0) {
    userTrendChart.clear()
    return
  }

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '新增用户',
        type: 'line',
        smooth: true,
        data: data.map(item => item.count),
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
          ])
        },
        lineStyle: {
          color: '#409EFF'
        },
        itemStyle: {
          color: '#409EFF'
        }
      }
    ]
  }

  userTrendChart.setOption(option)
}

// 渲染资料趋势图表
function renderMaterialTrendChart(data: TrendData[]) {
  if (!materialTrendChartRef.value) return

  if (!materialTrendChart) {
    materialTrendChart = echarts.init(materialTrendChartRef.value)
  }

  if (!data || data.length === 0) {
    materialTrendChart.clear()
    return
  }

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '新增资料',
        type: 'line',
        smooth: true,
        data: data.map(item => item.count),
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
            { offset: 1, color: 'rgba(103, 194, 58, 0.05)' }
          ])
        },
        lineStyle: {
          color: '#67C23A'
        },
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }

  materialTrendChart.setOption(option)
}

// 渲染下载趋势图表
function renderDownloadTrendChart(data: TrendData[]) {
  if (!downloadTrendChartRef.value) return

  if (!downloadTrendChart) {
    downloadTrendChart = echarts.init(downloadTrendChartRef.value)
  }

  if (!data || data.length === 0) {
    downloadTrendChart.clear()
    return
  }

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '下载次数',
        type: 'line',
        smooth: true,
        data: data.map(item => item.count),
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(230, 162, 60, 0.3)' },
            { offset: 1, color: 'rgba(230, 162, 60, 0.05)' }
          ])
        },
        lineStyle: {
          color: '#E6A23C'
        },
        itemStyle: {
          color: '#E6A23C'
        }
      }
    ]
  }

  downloadTrendChart.setOption(option)
}

// 渲染访问趋势图表
function renderVisitTrendChart(data: TrendData[]) {
  if (!visitTrendChartRef.value) return

  if (!visitTrendChart) {
    visitTrendChart = echarts.init(visitTrendChartRef.value)
  }

  if (!data || data.length === 0) {
    visitTrendChart.clear()
    return
  }

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.map(item => item.date),
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '访问次数',
        type: 'line',
        smooth: true,
        data: data.map(item => item.count),
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245, 108, 108, 0.3)' },
            { offset: 1, color: 'rgba(245, 108, 108, 0.05)' }
          ])
        },
        lineStyle: {
          color: '#F56C6C'
        },
        itemStyle: {
          color: '#F56C6C'
        }
      }
    ]
  }

  visitTrendChart.setOption(option)
}

// 窗口大小改变时重新渲染图表
function handleResize() {
  userTrendChart?.resize()
  materialTrendChart?.resize()
  downloadTrendChart?.resize()
  visitTrendChart?.resize()
}

onMounted(() => {
  // 初始加载所有数据
  loadOverviewStatistics()
  loadUserTrend()
  loadMaterialTrend()
  loadDownloadTrend()
  loadVisitTrend()
  window.addEventListener('resize', handleResize)

  // 设置自动刷新（每 30 秒刷新一次概览统计）
  refreshTimer = setInterval(() => {
    loadOverviewStatistics()
  }, 30000) // 30 秒
})

onUnmounted(() => {
  // 清理定时器
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }

  // 清理图表实例
  userTrendChart?.dispose()
  materialTrendChart?.dispose()
  downloadTrendChart?.dispose()
  visitTrendChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped lang="scss">
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

// 统计卡片
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;

  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  transition: all 0.2s ease;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.users {
    background: #eff6ff;
    color: #3b82f6;
  }

  &.materials {
    background: #f0fdf4;
    color: #22c55e;
  }

  &.downloads {
    background: #fef3c7;
    color: #f59e0b;
  }

  &.visits {
    background: #fef2f2;
    color: #ef4444;
  }
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  line-height: 1.2;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.stat-trend {
  font-size: 12px;
  color: #9ca3af;
}

// 待处理事项
.pending-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 16px 0;
}

.pending-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
  }
}

.pending-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: #111827;
    background: #f9fafb;
  }

  .arrow-icon {
    color: #9ca3af;
    transition: transform 0.2s ease;
  }

  &:hover .arrow-icon {
    transform: translateX(4px);
    color: #111827;
  }
}

.pending-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.pending-badge {
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 600;
  border-radius: 8px;

  &.material {
    background: #dcfce7;
    color: #15803d;
  }

  &.application {
    background: #dbeafe;
    color: #1d4ed8;
  }

  &.report {
    background: #fee2e2;
    color: #b91c1c;
  }
}

.pending-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.pending-count {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  line-height: 1;
}

.pending-text {
  font-size: 13px;
  color: #6b7280;
}

// 图表区域
.charts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
  }
}

.chart-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;

  &.large {
    grid-column: span 1;

    @media (max-width: 900px) {
      grid-column: span 1;
    }
  }
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;

  h3 {
    font-size: 15px;
    font-weight: 600;
    color: #111827;
    margin: 0;
  }
}

.time-selector {
  display: flex;
  gap: 4px;
  background: #f3f4f6;
  padding: 3px;
  border-radius: 8px;
}

.time-btn {
  padding: 6px 12px;
  font-size: 13px;
  font-weight: 500;
  border: none;
  background: transparent;
  color: #6b7280;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    color: #111827;
  }

  &.active {
    background: #ffffff;
    color: #111827;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }
}

.chart-body {
  height: 250px;
  width: 100%;

  @media (max-width: 768px) {
    height: 200px;
  }
}
</style>

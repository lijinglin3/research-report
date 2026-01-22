<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import './App.css'

// 常量定义
const API_CONFIG = {
  BASE_URL: 'https://reportapi.eastmoney.com/report/list2',
  BK_API_URL: 'https://reportapi.eastmoney.com/report/bk?bkCode=016'
}

// 每页显示数量选项
const PAGE_SIZE_OPTIONS = [20, 50, 100, 200]

const REPORT_TYPES = [
  { value: 0, label: '个股研报' },
  { value: 1, label: '行业研报' },
  { value: 2, label: '宏观研报' }
]

// LocalStorage配置
const STORAGE_KEYS = {
  Q_TYPE: 'researchReport_qType',
  BEGIN_TIME: 'researchReport_beginTime',
  END_TIME: 'researchReport_endTime',
  PAGE_SIZE: 'researchReport_pageSize',
  BK_CODE: 'researchReport_bkCode'
}

// 报告数据状态
const reports = ref([])
const loading = ref(true)
const total = ref(0)

// API请求参数
const qType = ref(0)
const beginTime = ref('')
const endTime = ref('')
const pageNo = ref(1)
const pageSize = ref(PAGE_SIZE_OPTIONS[1])
const bkCode = ref('')
const bkOptions = ref([{ value: '', label: '全部板块' }])

// 防止初始加载时的重复请求标志
const isInitialLoad = ref(true)

// 监听时间参数变化，禁止删除
watch([beginTime, endTime], ([newBegin, newEnd]) => {
  if (!newBegin || !newEnd) {
    setDefaultDateRange()
  }
})

// 工具函数：格式化日期为YYYY-MM-DD（使用本地时间）
const formatDate = date => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 工具函数：截断字符串并添加省略号
const truncateString = (str, maxLength = 4) => {
  if (!str || typeof str !== 'string') return '-'
  return str.length > maxLength ? `${str.substring(0, maxLength)}...` : str
}

// 设置默认时间范围（默认查询一年数据）
const setDefaultDateRange = () => {
  const today = new Date()
  const oneYearAgo = new Date(today.getFullYear() - 1, today.getMonth(), today.getDate())
  endTime.value = formatDate(today)
  beginTime.value = formatDate(oneYearAgo)
}

// ==============================
// LocalStorage操作
// ==============================
/**
 * 从localStorage加载保存的查询参数
 */
const loadSavedParams = () => {
  try {
    // 加载类型
    const savedQType = localStorage.getItem(STORAGE_KEYS.Q_TYPE)
    if (savedQType !== null) qType.value = parseInt(savedQType)

    // 加载每页显示数量
    const savedPageSize = localStorage.getItem(STORAGE_KEYS.PAGE_SIZE)
    if (savedPageSize !== null) {
      const size = parseInt(savedPageSize)
      pageSize.value = PAGE_SIZE_OPTIONS.includes(size) ? size : PAGE_SIZE_OPTIONS[0]
    }

    // 加载板块代码
    const savedBkCode = localStorage.getItem(STORAGE_KEYS.BK_CODE)
    if (savedBkCode !== null) bkCode.value = savedBkCode

    // 加载日期范围
    const savedBeginTime = localStorage.getItem(STORAGE_KEYS.BEGIN_TIME)
    const savedEndTime = localStorage.getItem(STORAGE_KEYS.END_TIME)
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/

    if (savedBeginTime && savedEndTime && dateRegex.test(savedBeginTime) && dateRegex.test(savedEndTime)) {
      beginTime.value = savedBeginTime
      endTime.value = savedEndTime
    } else {
      setDefaultDateRange()
    }
  } catch (err) {
    console.error('Failed to load saved params:', err)
    setDefaultDateRange()
  }
}

/**
 * 获取板块数据
 */
const fetchBK = async () => {
  try {
    const response = await fetch(API_CONFIG.BK_API_URL, {
      headers: { 'Content-Type': 'application/json' }
    })

    if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)

    const data = await response.json()
    bkOptions.value = [{ value: '', label: '全部板块' }]

    data.data.forEach(bk => {
      if (bk && typeof bk === 'object') {
        bkOptions.value.push({ value: bk.bkCode, label: bk.bkName })
      }
    })
  } catch (err) {
    console.error('Error fetching BK:', err)
    bkOptions.value = [{ value: '', label: '全部板块' }]
    ElMessage.error('获取板块数据失败')
  }
}

/**
 * 保存查询参数到localStorage
 */
const saveParamsToLocalStorage = () => {
  try {
    localStorage.setItem(STORAGE_KEYS.Q_TYPE, qType.value.toString())
    localStorage.setItem(STORAGE_KEYS.BEGIN_TIME, beginTime.value)
    localStorage.setItem(STORAGE_KEYS.END_TIME, endTime.value)
    localStorage.setItem(STORAGE_KEYS.PAGE_SIZE, pageSize.value.toString())
    localStorage.setItem(STORAGE_KEYS.BK_CODE, bkCode.value)
  } catch (err) {
    console.error('Failed to save params:', err)
  }
}

// ==============================
// 表格相关工具函数
// ==============================
/**
 * 生成PDF链接
 */
const getPdfUrl = row => {
  if (!row || !row.infoCode) return ''
  return `https://pdf.dfcfw.com/pdf/H3_${row.infoCode}_1.pdf`
}

/**
 * 格式化表格中的日期显示
 */
const formatTableCellDate = row => {
  const date = row.publishDate || row.reportDate
  if (!date) return '-'
  if (typeof date === 'string') return date.substring(0, 10)
  if (date instanceof Date && !isNaN(date.getTime())) return formatDate(date)
  return '-'
}

// ==============================
// API请求相关
// ==============================
/**
 * 构建API请求参数
 */
const buildApiParams = () => ({
  qType: qType.value,
  beginTime: beginTime.value,
  endTime: endTime.value,
  pageNo: pageNo.value,
  pageSize: pageSize.value,
  ...(bkCode.value && { industryCode: bkCode.value })
})

/**
 * 从API获取研究报告数据
 */
const fetchReports = async () => {
  loading.value = true
  try {
    const response = await fetch(API_CONFIG.BASE_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(buildApiParams())
    })

    const data = await response.json()
    reports.value = data.data || []
    total.value = data.hits || 0
    if (reports.value.length === 0) ElMessage.info('暂无数据')
  } catch (err) {
    ElMessage.error('获取数据失败，请稍后重试')
    console.error('Error fetching reports:', err)
  } finally {
    loading.value = false
  }
}

// ==============================
// 事件处理函数
// ==============================
/**
 * 切换到指定页码
 */
const handleCurrentChange = page => {
  pageNo.value = page
  fetchReports()
}

/**
 * 重新查询时重置页码到第一页
 */
const resetAndFetch = () => {
  pageNo.value = 1
  fetchReports()
}

/**
 * 重置参数为默认值
 */
const resetParams = () => {
  qType.value = REPORT_TYPES[0].value
  setDefaultDateRange()
  pageSize.value = PAGE_SIZE_OPTIONS[1]
  bkCode.value = ''
  pageNo.value = 1
  saveParamsToLocalStorage()
  fetchReports()
}

// ==============================
// 事件监听
// ==============================
// 监听参数变化，保存到localStorage并重新查询
watch([qType, beginTime, endTime, pageSize, bkCode], () => {
  saveParamsToLocalStorage()
  // 初始加载时不触发重复请求
  if (isInitialLoad.value) {
    isInitialLoad.value = false
    return
  }
  resetAndFetch()
}, { immediate: false })

// 组件挂载时加载保存的参数并获取数据
onMounted(async () => {
  loadSavedParams()
  await fetchBK()
  fetchReports()
})
</script>

<template>
  <div class="app-container">
    <el-container>
      <el-main>
        <!-- 查询参数表单 -->
        <el-form :inline="true" class="search-form" label-position="left">
          <el-form-item label="类型">
            <el-select v-model="qType" placeholder="请选择类型" style="width: 100px;">
              <el-option v-for="type in REPORT_TYPES" :key="type.value" :label="type.label"
                :value="type.value"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="板块">
            <el-select v-model="bkCode" placeholder="请选择板块" filterable clearable style="width: 140px;">
              <el-option v-for="bk in bkOptions" :key="bk.value" :label="bk.label" :value="bk.value"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="开始时间">
            <el-date-picker v-model="beginTime" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD"
              placeholder="选择开始时间" style="width: 120px;" :clearable="false"></el-date-picker>
          </el-form-item>

          <el-form-item label="结束时间">
            <el-date-picker v-model="endTime" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD"
              placeholder="选择结束时间" style="width: 120px;" :clearable="false"></el-date-picker>
          </el-form-item>

          <el-form-item label="每页显示">
            <el-select v-model="pageSize" placeholder="请选择每页显示数量" style="width: 80px;">
              <el-option v-for="size in PAGE_SIZE_OPTIONS" :key="size" :label="size" :value="size"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="default" @click="resetParams">重置</el-button>
          </el-form-item>
        </el-form>

        <!-- 报告表格 -->
        <el-table v-loading="loading" :data="reports" stripe border style="width: 100%; margin-bottom: 20px;">
          <!-- 日期列 -->
          <el-table-column prop="reportDate" label="日期" width="120">
            <template #default="scope">
              {{ formatTableCellDate(scope.row) }}
            </template>
          </el-table-column>

          <!-- 券商列 -->
          <el-table-column prop="orgSName" label="券商" width="100">
            <template #default="scope">
              <el-tooltip :content="scope.row.orgSName || scope.row.orgName || '-'" placement="top">
                <span>{{ truncateString(scope.row.orgSName || scope.row.orgName) }}</span>
              </el-tooltip>
            </template>
          </el-table-column>

          <!-- 行业列（非宏观研报显示） -->
          <el-table-column v-if="qType !== 2" prop="industryName" label="行业" width="120">
            <template #default="scope">
              {{ scope.row.indvInduName || scope.row.industryName || '-' }}
            </template>
          </el-table-column>

          <!-- 股票列（个股研报显示） -->
          <el-table-column v-if="qType === 0" prop="stockName" label="股票" width="100">
            <template #default="scope">
              {{ scope.row.stockName || '-' }}
            </template>
          </el-table-column>

          <!-- 页面列 -->
          <el-table-column prop="attachPages" label="页面" width="60">
            <template #default="scope">
              {{ scope.row.attachPages || scope.row.AttachPages || '-' }}
            </template>
          </el-table-column>

          <!-- 标题列（带超链接） -->
          <el-table-column prop="title" label="标题" min-width="300">
            <template #default="scope">
              <a v-if="getPdfUrl(scope.row)" :href="getPdfUrl(scope.row)" target="_blank" rel="noopener noreferrer"
                class="report-title-link">
                {{ scope.row.title || '-' }}
              </a>
              <span v-else>{{ scope.row.title || '-' }}</span>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页控件 -->
        <div class="pagination-container">
          <el-pagination background layout="total, prev, pager, next, jumper" :current-page="pageNo"
            :page-size="pageSize" :total="total" @current-change="handleCurrentChange"></el-pagination>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

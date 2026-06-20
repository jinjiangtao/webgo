<template>
  <div class="search-page">
    <div class="search-hero">
      <div class="hero-title">
        <h1>智能关键词检索</h1>
        <p>基于分词匹配算法，精准定位内容，提升检索效率</p>
      </div>

      <div class="search-box-wrapper">
        <div class="search-box" :class="{ 'is-focused': isFocused }">
          <el-icon class="search-icon" :size="20"><Search /></el-icon>
          <input
            ref="searchInput"
            v-model="searchQuery"
            type="text"
            class="search-input"
            placeholder="输入关键词开始搜索..."
            @input="handleInput"
            @focus="handleFocus"
            @keydown="handleKeydown"
            @blur="handleBlur"
            autocomplete="off"
          />
          <el-button
            v-if="searchQuery"
            class="clear-btn"
            type="text"
            circle
            @click="clearSearch"
          >
            <el-icon><Close /></el-icon>
          </el-button>
          <el-button
            type="primary"
            class="search-btn"
            @click="doSearch"
            :loading="loading"
          >
            搜索
          </el-button>

          <transition name="dropdown">
            <div
              v-if="showSuggestions && suggestions.length"
              class="suggestions-dropdown"
            >
              <div
                v-for="(item, idx) in suggestions"
                :key="idx"
                class="suggestion-item"
                :class="{ active: suggestionIndex === idx }"
                @mousedown.prevent="selectSuggestion(item)"
                @mouseenter="suggestionIndex = idx"
              >
                <el-icon><Search /></el-icon>
                <span v-html="highlightText(item, searchQuery)"></span>
              </div>
            </div>
          </transition>
        </div>
      </div>

      <div v-if="hotKeywords.length" class="hot-keywords">
        <span class="hot-label">
          <el-icon><HotWater /></el-icon>
          热门搜索：
        </span>
        <span
          v-for="item in hotKeywords"
          :key="item.keyword"
          class="hot-tag"
          @click="quickSearch(item.keyword)"
        >
          {{ item.keyword }}
        </span>
      </div>
    </div>

    <div class="search-results-section">
      <div v-if="hasSearched" class="results-header">
        <div class="results-summary">
          <span v-if="searchQuery">
            找到 <strong>{{ total }}</strong> 条与 "<em>{{ searchQuery }}</em>" 相关的结果
            <span v-if="tokens.length" class="tokens-info">
              （分词：
              <span v-for="(t, i) in tokens.slice(0, 8)" :key="i" class="token-tag">{{ t }}</span>
              <span v-if="tokens.length > 8">...</span>
              ）
            </span>
          </span>
          <span v-else>共 <strong>{{ total }}</strong> 条内容</span>
          <span v-if="selectedCategory" class="filter-info">
            · 分类：{{ categoryNameMap[selectedCategory] }}
          </span>
        </div>

        <div class="sort-control">
          <el-radio-group v-model="sortOrder" size="small" @change="doSearch">
            <el-radio-button value="desc">相关度优先</el-radio-button>
            <el-radio-button value="asc">相关度从低到高</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <div class="results-body">
        <div class="filter-sidebar">
          <div class="filter-card">
            <div class="filter-title">
              <el-icon><Menu /></el-icon>
              分类筛选
            </div>
            <div class="category-list">
              <div
                class="category-item"
                :class="{ active: !selectedCategory }"
                @click="filterByCategory(0)"
              >
                <span>全部分类</span>
              </div>
              <div
                v-for="cat in categories"
                :key="cat.id"
                class="category-item"
                :class="{ active: selectedCategory === cat.id }"
                @click="filterByCategory(cat.id)"
              >
                <span>{{ cat.name }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="results-list-wrapper">
          <div v-if="loading" class="loading-wrapper">
            <el-icon class="loading-icon" :size="40"><Loading /></el-icon>
            <p>正在搜索中...</p>
          </div>

          <div v-else-if="results.length === 0 && hasSearched" class="empty-wrapper">
            <el-empty description="未找到相关内容，试试其他关键词吧">
              <el-button type="primary" @click="clearSearch">清除搜索条件</el-button>
            </el-empty>
          </div>

          <div v-else-if="results.length" class="results-list">
            <div
              v-for="item in results"
              :key="item.keyword.id"
              class="result-card"
              @click="openDetail(item)"
            >
              <div class="result-header">
                <h3 class="result-title" v-html="item.title_highlighted || item.keyword.title"></h3>
                <div class="result-meta">
                  <el-tag
                    v-if="item.keyword.category"
                    :type="getCategoryTagType(item.keyword.category_id)"
                    size="small"
                    effect="light"
                  >
                    {{ item.keyword.category.name }}
                  </el-tag>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ item.keyword.view_count }}
                  </span>
                  <span class="meta-item">
                    <el-icon><Star /></el-icon>
                    {{ item.score.toFixed(1) }}
                  </span>
                </div>
              </div>
              <p class="result-content" v-html="item.content_highlighted || item.keyword.content"></p>
              <div v-if="item.keyword.tags" class="result-tags">
                <span
                  v-for="(tag, i) in parseTags(item.keyword.tags)"
                  :key="i"
                  class="tag-item"
                >
                  #{{ tag }}
                </span>
              </div>
            </div>
          </div>

          <div v-if="total > pageSize" class="pagination-wrapper">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handlePageChange"
            />
          </div>
        </div>
      </div>
    </div>

    <el-dialog
      v-model="detailVisible"
      :title="currentDetail?.title"
      width="700px"
      destroy-on-close
    >
      <div v-if="currentDetail" class="detail-content">
        <div class="detail-tags" v-if="currentDetail.tags">
          <el-tag
            v-for="(tag, i) in parseTags(currentDetail.tags)"
            :key="i"
            size="default"
            effect="plain"
            class="mr-5"
          >
            {{ tag }}
          </el-tag>
        </div>
        <div class="detail-info">
          <span v-if="currentDetail.category">
            <strong>分类：</strong>{{ currentDetail.category.name }}
          </span>
          <span><strong>浏览：</strong>{{ currentDetail.view_count }} 次</span>
          <span><strong>创建时间：</strong>{{ formatDate(currentDetail.created_at) }}</span>
        </div>
        <div class="detail-body">
          {{ currentDetail.content }}
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  search,
  getSuggestions,
  getHotKeywords,
  listCategories,
  incrementView,
  getKeywordDetail
} from '@/api'

const searchInput = ref(null)
const searchQuery = ref('')
const isFocused = ref(false)
const showSuggestions = ref(false)
const suggestions = ref([])
const suggestionIndex = ref(-1)
const hotKeywords = ref([])
const categories = ref([])
const categoryNameMap = ref({})

const loading = ref(false)
const hasSearched = ref(false)
const results = ref([])
const total = ref(0)
const tokens = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const sortOrder = ref('desc')
const selectedCategory = ref(0)

const detailVisible = ref(false)
const currentDetail = ref(null)

let suggestTimer = null

const debounceSuggest = () => {
  if (suggestTimer) clearTimeout(suggestTimer)
  suggestTimer = setTimeout(async () => {
    if (!searchQuery.value.trim()) {
      suggestions.value = []
      return
    }
    try {
      const res = await getSuggestions({
        q: searchQuery.value,
        category_id: selectedCategory.value || undefined,
        limit: 10
      })
      suggestions.value = res.data || []
    } catch (e) {
      suggestions.value = []
    }
  }, 200)
}

const handleInput = () => {
  suggestionIndex.value = -1
  debounceSuggest()
}

const handleFocus = () => {
  isFocused.value = true
  showSuggestions.value = true
  if (searchQuery.value.trim()) debounceSuggest()
}

const handleBlur = () => {
  isFocused.value = false
  setTimeout(() => {
    showSuggestions.value = false
  }, 200)
}

const handleKeydown = (e) => {
  if (!suggestions.value.length) {
    if (e.key === 'Enter') doSearch()
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    suggestionIndex.value = (suggestionIndex.value + 1) % suggestions.value.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    suggestionIndex.value = suggestionIndex.value <= 0 ? suggestions.value.length - 1 : suggestionIndex.value - 1
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (suggestionIndex.value >= 0) {
      selectSuggestion(suggestions.value[suggestionIndex.value])
    } else {
      doSearch()
    }
  } else if (e.key === 'Escape') {
    showSuggestions.value = false
    suggestionIndex.value = -1
  }
}

const selectSuggestion = (item) => {
  searchQuery.value = item
  showSuggestions.value = false
  suggestionIndex.value = -1
  doSearch()
}

const quickSearch = (keyword) => {
  searchQuery.value = keyword
  doSearch()
}

const clearSearch = () => {
  searchQuery.value = ''
  selectedCategory.value = 0
  suggestions.value = []
  nextTick(() => searchInput.value?.focus())
  doSearch()
}

const filterByCategory = (catId) => {
  selectedCategory.value = catId
  currentPage.value = 1
  doSearch()
}

const doSearch = async () => {
  loading.value = true
  hasSearched.value = true
  try {
    const sessionId = localStorage.getItem('session_id') || ''
    const res = await search({
      q: searchQuery.value.trim(),
      category_id: selectedCategory.value || undefined,
      page: currentPage.value,
      page_size: pageSize.value,
      sort_order: sortOrder.value,
      session_id: sessionId
    })
    results.value = res.data.list || []
    total.value = res.data.total || 0
    tokens.value = res.data.tokens || []
  } catch (e) {
    ElMessage.error('搜索失败')
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  doSearch()
  window.scrollTo({ top: 300, behavior: 'smooth' })
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  doSearch()
}

const parseTags = (tagsStr) => {
  if (!tagsStr) return []
  return tagsStr.split(',').map(t => t.trim()).filter(Boolean)
}

const highlightText = (text, keyword) => {
  if (!keyword) return text
  const kw = keyword.trim()
  if (!kw) return text
  const reg = new RegExp(`(${kw.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  return text.replace(reg, '<mark class="search-highlight">$1</mark>')
}

const getCategoryTagType = (id) => {
  const types = ['primary', 'success', 'warning', 'info', 'danger']
  return types[id % types.length]
}

const formatDate = (d) => {
  if (!d) return '-'
  const date = new Date(d)
  return date.toLocaleString('zh-CN')
}

const openDetail = async (item) => {
  try {
    await incrementView(item.keyword.id)
  } catch (e) {}
  try {
    const res = await getKeywordDetail(item.keyword.id)
    currentDetail.value = res.data
    detailVisible.value = true
  } catch (e) {
    ElMessage.error('获取详情失败')
  }
}

const loadCategories = async () => {
  try {
    const res = await listCategories()
    categories.value = res.data || []
    categories.value.forEach(c => {
      categoryNameMap.value[c.id] = c.name
    })
  } catch (e) {}
}

const loadHotKeywords = async () => {
  try {
    const res = await getHotKeywords({ limit: 8 })
    hotKeywords.value = res.data || []
  } catch (e) {}
}

onMounted(() => {
  loadCategories()
  loadHotKeywords()
  nextTick(() => {
    searchInput.value?.focus()
    doSearch()
  })
})
</script>

<style scoped>
.search-page {
  padding-bottom: 40px;
}
.search-hero {
  text-align: center;
  padding: 60px 0 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  margin-bottom: 32px;
  color: #fff;
  position: relative;
  overflow: hidden;
}
.search-hero::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 60%);
}
.hero-title {
  position: relative;
  z-index: 1;
}
.hero-title h1 {
  font-size: 36px;
  margin-bottom: 12px;
  letter-spacing: 2px;
}
.hero-title p {
  opacity: 0.9;
  font-size: 16px;
}
.search-box-wrapper {
  position: relative;
  max-width: 700px;
  margin: 30px auto 0;
  z-index: 1;
}
.search-box {
  display: flex;
  align-items: center;
  background: #fff;
  border-radius: 50px;
  padding: 6px 6px 6px 24px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  position: relative;
}
.search-box.is-focused {
  box-shadow: 0 10px 40px rgba(64, 158, 255, 0.3);
  transform: translateY(-2px);
}
.search-icon {
  color: #c0c4cc;
  margin-right: 12px;
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  padding: 12px 0;
  background: transparent;
  color: #303133;
  min-width: 0;
}
.clear-btn {
  margin-right: 8px;
}
.search-btn {
  border-radius: 50px;
  padding: 12px 28px;
  font-weight: 600;
  flex-shrink: 0;
}
.suggestions-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  overflow: hidden;
  z-index: 200;
  text-align: left;
}
.suggestion-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  color: #303133;
  transition: background 0.15s;
}
.suggestion-item:hover,
.suggestion-item.active {
  background: #f5f7fa;
}
.suggestion-item .el-icon {
  color: #c0c4cc;
  font-size: 14px;
}
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
.hot-keywords {
  margin-top: 20px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 8px;
  position: relative;
  z-index: 1;
}
.hot-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  opacity: 0.85;
  font-size: 14px;
  margin-right: 4px;
}
.hot-tag {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  padding: 4px 14px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.hot-tag:hover {
  background: rgba(255, 255, 255, 0.35);
  transform: translateY(-1px);
}
.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border-radius: 12px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
.results-summary {
  color: #606266;
  font-size: 14px;
}
.results-summary strong {
  color: #409EFF;
  margin: 0 4px;
}
.results-summary em {
  color: #f56c6c;
  font-style: normal;
  font-weight: 600;
  padding: 0 4px;
}
.tokens-info {
  margin-left: 8px;
  color: #909399;
}
.token-tag {
  display: inline-block;
  background: #ecf5ff;
  color: #409EFF;
  padding: 1px 8px;
  border-radius: 4px;
  font-size: 12px;
  margin: 0 2px;
}
.filter-info {
  margin-left: 8px;
  color: #909399;
}
.results-body {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}
.filter-sidebar {
  width: 200px;
  flex-shrink: 0;
  position: sticky;
  top: 88px;
}
.filter-card {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
.filter-title {
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
}
.category-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.category-item {
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  transition: all 0.2s;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.category-item:hover {
  background: #f5f7fa;
  color: #303133;
}
.category-item.active {
  background: #ecf5ff;
  color: #409EFF;
  font-weight: 500;
}
.results-list-wrapper {
  flex: 1;
  min-width: 0;
}
.loading-wrapper,
.empty-wrapper {
  background: #fff;
  border-radius: 12px;
  padding: 60px 20px;
  text-align: center;
}
.loading-icon {
  color: #409EFF;
  animation: rotate 1s linear infinite;
  margin-bottom: 12px;
  display: inline-block;
}
@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.results-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.result-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  cursor: pointer;
  transition: all 0.25s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  border: 1px solid transparent;
}
.result-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
  border-color: #e4e7ed;
}
.result-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
  gap: 16px;
}
.result-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  line-height: 1.5;
}
.result-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}
.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 13px;
}
.result-content {
  color: #606266;
  line-height: 1.8;
  font-size: 14px;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.result-tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.tag-item {
  color: #909399;
  font-size: 12px;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
}
.pagination-wrapper {
  margin-top: 24px;
  display: flex;
  justify-content: center;
  padding: 16px;
  background: #fff;
  border-radius: 12px;
}
.detail-content {
  padding: 8px 0;
}
.detail-tags {
  margin-bottom: 16px;
}
.detail-info {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 14px;
  color: #606266;
}
.detail-body {
  line-height: 2;
  color: #303133;
  font-size: 15px;
  white-space: pre-wrap;
}
.mr-5 {
  margin-right: 8px;
  margin-bottom: 4px;
}
@media (max-width: 768px) {
  .search-hero {
    padding: 40px 16px 32px;
  }
  .hero-title h1 {
    font-size: 26px;
  }
  .results-body {
    flex-direction: column;
  }
  .filter-sidebar {
    width: 100%;
    position: static;
  }
  .category-list {
    flex-direction: row;
    flex-wrap: wrap;
  }
  .category-item {
    flex: 0 0 auto;
  }
}
</style>

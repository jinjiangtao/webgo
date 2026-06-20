<template>
  <div class="admin-page">
    <el-tabs v-model="activeTab" class="admin-tabs">
      <el-tab-pane label="关键词库" name="keywords">
        <div class="panel-card">
          <div class="toolbar">
            <div class="toolbar-left">
              <el-input
                v-model="searchKw"
                placeholder="搜索标题/内容/标签"
                style="width: 260px"
                clearable
                @keyup.enter="loadKeywords"
                @clear="loadKeywords"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
              <el-select
                v-model="filterCategory"
                placeholder="选择分类"
                clearable
                style="width: 160px; margin-left: 12px"
                @change="loadKeywords"
                @clear="loadKeywords"
              >
                <el-option
                  v-for="cat in categories"
                  :key="cat.id"
                  :label="cat.name"
                  :value="cat.id"
                />
              </el-select>
              <el-select
                v-model="filterStatus"
                placeholder="状态"
                clearable
                style="width: 120px; margin-left: 12px"
                @change="loadKeywords"
                @clear="loadKeywords"
              >
                <el-option :label="启用" :value="1" />
                <el-option :label="禁用" :value="0" />
              </el-select>
              <el-button type="primary" @click="loadKeywords" style="margin-left: 12px">
                <el-icon><Search /></el-icon>查询
              </el-button>
            </div>
            <div class="toolbar-right">
              <el-button @click="showBatchImport">
                <el-icon><Upload /></el-icon>批量导入
              </el-button>
              <el-button type="primary" @click="showAddDialog">
                <el-icon><Plus /></el-icon>新增关键词
              </el-button>
            </div>
          </div>

          <el-table :data="keywordList" v-loading="loading" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="title" label="标题" min-width="220">
              <template #default="{ row }">
                <span class="title-text">{{ row.title }}</span>
              </template>
            </el-table-column>
            <el-table-column label="分类" width="120" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.category" size="small" :type="getCatType(row.category_id)">
                  {{ row.category.name }}
                </el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="tags" label="标签" width="180">
              <template #default="{ row }">
                <div class="tags-cell" v-if="row.tags">
                  <el-tag
                    v-for="(t, i) in parseTags(row.tags).slice(0, 3)"
                    :key="i"
                    size="small"
                    effect="plain"
                    style="margin-right: 4px"
                  >{{ t }}</el-tag>
                </div>
                <span v-else class="empty-text">-</span>
              </template>
            </el-table-column>
            <el-table-column prop="sort" label="排序" width="80" align="center" />
            <el-table-column prop="view_count" label="浏览" width="80" align="center" />
            <el-table-column label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handleStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="170" align="center">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="180" align="center" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="primary" link @click="showEditDialog(row)">
                  编辑
                </el-button>
                <el-button size="small" type="danger" link @click="handleDelete(row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="keywordPage"
              v-model:page-size="keywordPageSize"
              :total="keywordTotal"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="loadKeywords"
              @current-change="loadKeywords"
            />
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="分类管理" name="categories">
        <div class="panel-card">
          <div class="toolbar">
            <div class="toolbar-left">
              <span class="hint">共 {{ categories.length }} 个分类</span>
            </div>
            <div class="toolbar-right">
              <el-button type="primary" @click="showAddCategory">
                <el-icon><Plus /></el-icon>新增分类
              </el-button>
            </div>
          </div>

          <el-table :data="categories" v-loading="catLoading" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column prop="name" label="分类名称" width="180">
              <template #default="{ row }">
                <span class="cat-name">{{ row.name }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="300">
              <template #default="{ row }">
                {{ row.description || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-switch
                  v-model="row.status"
                  :active-value="1"
                  :inactive-value="0"
                  @change="handleCatStatusChange(row)"
                />
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="170" align="center">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" align="center">
              <template #default="{ row }">
                <el-button size="small" type="primary" link @click="showEditCategory(row)">
                  编辑
                </el-button>
                <el-button size="small" type="danger" link @click="handleDeleteCategory(row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog
      v-model="keywordDialogVisible"
      :title="isEdit ? '编辑关键词' : '新增关键词'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="keywordFormRef"
        :model="keywordForm"
        :rules="keywordRules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="keywordForm.title" maxlength="200" show-word-limit placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="keywordForm.category_id" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="keywordForm.content"
            type="textarea"
            :rows="6"
            maxlength="5000"
            show-word-limit
            placeholder="请输入详细内容"
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-input
            v-model="keywordForm.tags"
            placeholder="多个标签用逗号分隔，如：go,后端,教程"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="keywordForm.sort" :min="0" :max="9999" />
          <span class="hint">数值越大越靠前</span>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="keywordForm.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="keywordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitKeyword">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="categoryDialogVisible"
      :title="isCatEdit ? '编辑分类' : '新增分类'"
      width="480px"
      destroy-on-close
    >
      <el-form
        ref="categoryFormRef"
        :model="categoryForm"
        :rules="categoryRules"
        label-width="100px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="categoryForm.name" maxlength="100" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="categoryForm.description"
            type="textarea"
            :rows="3"
            maxlength="500"
            placeholder="请输入分类描述"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch
            v-model="categoryForm.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategory">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="importDialogVisible"
      title="批量导入关键词"
      width="520px"
      destroy-on-close
    >
      <el-alert
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      >
        <template #title>CSV 格式说明</template>
        <p>1. 第一行为表头：标题,内容,分类ID,标签,状态,排序</p>
        <p>2. 标题、内容为必填项</p>
        <p>3. 多个标签用英文逗号分隔</p>
        <p>4. 状态：1=启用，0=禁用，默认为1</p>
      </el-alert>
      <el-upload
        drag
        :auto-upload="false"
        :limit="1"
        accept=".csv"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        :file-list="fileList"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          拖拽CSV文件到此处，或 <em>点击选择文件</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">只能上传 CSV 文件</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="doImport">
          开始导入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  listKeywords,
  createKeyword,
  updateKeyword,
  deleteKeyword,
  setKeywordStatus,
  importCSV,
  listCategories,
  createCategory,
  updateCategory,
  deleteCategory,
  setCategoryStatus
} from '@/api'

const activeTab = ref('keywords')

const searchKw = ref('')
const filterCategory = ref(null)
const filterStatus = ref(null)
const keywordPage = ref(1)
const keywordPageSize = ref(20)
const keywordTotal = ref(0)
const keywordList = ref([])
const loading = ref(false)

const categories = ref([])
const catLoading = ref(false)

const keywordDialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref(null)
const keywordFormRef = ref(null)
const keywordForm = reactive({
  title: '',
  category_id: null,
  content: '',
  tags: '',
  sort: 0,
  status: 1
})
const keywordRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }]
}

const categoryDialogVisible = ref(false)
const isCatEdit = ref(false)
const editingCatId = ref(null)
const categoryFormRef = ref(null)
const categoryForm = reactive({
  name: '',
  description: '',
  status: 1
})
const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const importDialogVisible = ref(false)
const fileList = ref([])
const currentFile = ref(null)
const importing = ref(false)

const loadKeywords = async () => {
  loading.value = true
  try {
    const params = {
      page: keywordPage.value,
      page_size: keywordPageSize.value,
      sort_by: 'sort',
      sort_order: 'desc'
    }
    if (searchKw.value) params.q = searchKw.value
    if (filterCategory.value) params.category_id = filterCategory.value
    if (filterStatus.value !== null && filterStatus.value !== undefined) {
      params.status = filterStatus.value
    }
    const res = await listKeywords(params)
    keywordList.value = res.data.list || []
    keywordTotal.value = res.data.total || 0
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  catLoading.value = true
  try {
    const res = await listCategories({ all: true })
    categories.value = res.data || []
  } catch (e) {
    ElMessage.error('加载分类失败')
  } finally {
    catLoading.value = false
  }
}

const showAddDialog = () => {
  isEdit.value = false
  editingId.value = null
  Object.assign(keywordForm, {
    title: '',
    category_id: null,
    content: '',
    tags: '',
    sort: 0,
    status: 1
  })
  keywordDialogVisible.value = true
  nextTick(() => keywordFormRef.value?.clearValidate())
}

const showEditDialog = (row) => {
  isEdit.value = true
  editingId.value = row.id
  Object.assign(keywordForm, {
    title: row.title,
    category_id: row.category_id || null,
    content: row.content,
    tags: row.tags || '',
    sort: row.sort,
    status: row.status
  })
  keywordDialogVisible.value = true
  nextTick(() => keywordFormRef.value?.clearValidate())
}

const submitKeyword = async () => {
  await keywordFormRef.value?.validate()
  try {
    const payload = { ...keywordForm }
    if (!payload.category_id) delete payload.category_id
    if (isEdit.value) {
      await updateKeyword(editingId.value, payload)
      ElMessage.success('更新成功')
    } else {
      await createKeyword(payload)
      ElMessage.success('创建成功')
    }
    keywordDialogVisible.value = false
    loadKeywords()
  } catch (e) {}
}

const handleStatusChange = async (row) => {
  try {
    await setKeywordStatus(row.id, row.status)
    ElMessage.success('状态更新成功')
  } catch (e) {
    row.status = row.status === 1 ? 0 : 1
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确定要删除"${row.title}"吗？`, '确认删除', {
    type: 'warning',
    confirmButtonText: '确定删除',
    cancelButtonText: '取消'
  }).then(async () => {
    await deleteKeyword(row.id)
    ElMessage.success('删除成功')
    loadKeywords()
  }).catch(() => {})
}

const showAddCategory = () => {
  isCatEdit.value = false
  editingCatId.value = null
  Object.assign(categoryForm, { name: '', description: '', status: 1 })
  categoryDialogVisible.value = true
  nextTick(() => categoryFormRef.value?.clearValidate())
}

const showEditCategory = (row) => {
  isCatEdit.value = true
  editingCatId.value = row.id
  Object.assign(categoryForm, {
    name: row.name,
    description: row.description || '',
    status: row.status
  })
  categoryDialogVisible.value = true
  nextTick(() => categoryFormRef.value?.clearValidate())
}

const submitCategory = async () => {
  await categoryFormRef.value?.validate()
  try {
    if (isCatEdit.value) {
      await updateCategory(editingCatId.value, categoryForm)
      ElMessage.success('分类更新成功')
    } else {
      await createCategory(categoryForm)
      ElMessage.success('分类创建成功')
    }
    categoryDialogVisible.value = false
    loadCategories()
    if (activeTab.value === 'keywords') loadKeywords()
  } catch (e) {}
}

const handleCatStatusChange = async (row) => {
  try {
    await setCategoryStatus(row.id, row.status)
    ElMessage.success('状态更新成功')
  } catch (e) {
    row.status = row.status === 1 ? 0 : 1
  }
}

const handleDeleteCategory = (row) => {
  ElMessageBox.confirm(`确定要删除分类"${row.name}"吗？若有关联关键词将无法删除。`, '确认删除', {
    type: 'warning',
    confirmButtonText: '确定删除',
    cancelButtonText: '取消'
  }).then(async () => {
    try {
      await deleteCategory(row.id)
      ElMessage.success('删除成功')
      loadCategories()
    } catch (e) {}
  }).catch(() => {})
}

const showBatchImport = () => {
  importDialogVisible.value = true
  fileList.value = []
  currentFile.value = null
}

const handleFileChange = (file) => {
  if (!file.name.endsWith('.csv')) {
    ElMessage.warning('请上传CSV文件')
    return
  }
  fileList.value = [file]
  currentFile.value = file.raw
}

const handleFileRemove = () => {
  fileList.value = []
  currentFile.value = null
}

const doImport = async () => {
  if (!currentFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  importing.value = true
  try {
    const res = await importCSV(currentFile.value)
    ElMessage.success(`成功导入 ${res.data.imported_count} 条关键词`)
    importDialogVisible.value = false
    loadKeywords()
  } catch (e) {
  } finally {
    importing.value = false
  }
}

const parseTags = (tagsStr) => {
  if (!tagsStr) return []
  return tagsStr.split(',').map(t => t.trim()).filter(Boolean)
}

const formatDate = (d) => {
  if (!d) return '-'
  return new Date(d).toLocaleString('zh-CN')
}

const getCatType = (id) => {
  const types = ['primary', 'success', 'warning', 'info', 'danger']
  return types[id % types.length]
}

onMounted(() => {
  loadCategories().then(() => loadKeywords())
})
</script>

<style scoped>
.admin-page {
  min-height: 100%;
}
.admin-tabs :deep(.el-tabs__nav-wrap::after) {
  background-color: #e4e7ed;
}
.admin-tabs :deep(.el-tabs__item) {
  font-size: 15px;
  height: 48px;
  line-height: 48px;
}
.panel-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}
.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0;
}
.hint {
  color: #909399;
  font-size: 14px;
}
.title-text {
  font-weight: 500;
}
.cat-name {
  font-weight: 500;
}
.tags-cell {
  max-width: 180px;
  overflow: hidden;
}
.empty-text {
  color: #c0c4cc;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

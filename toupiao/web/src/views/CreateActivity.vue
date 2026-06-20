<template>
  <div class="create-activity page-container">
    <h2 class="page-title">
      <el-icon><Edit /></el-icon>
      创建投票活动
    </h2>

    <el-card class="form-card card-shadow">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="活动标题" prop="title">
          <el-input
            v-model="formData.title"
            placeholder="请输入活动标题"
            maxlength="100"
            show-word-limit
          >
          </el-input>
        </el-form-item>

        <el-form-item label="活动描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入活动描述"
            maxlength="500"
            show-word-limit
          >
          </el-input>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始时间" prop="start_time">
              <el-date-picker
                v-model="formData.start_time"
                type="datetime"
                placeholder="选择开始时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DDTHH:mm:ss"
                style="width: 100%"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束时间" prop="end_time">
              <el-date-picker
                v-model="formData.end_time"
                type="datetime"
                placeholder="选择结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DDTHH:mm:ss"
                style="width: 100%"
              >
              </el-date-picker>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="投票类型" prop="vote_type">
              <el-radio-group v-model="formData.vote_type">
                <el-radio value="single">单选投票</el-radio>
                <el-radio value="multiple">多选投票</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item
              label="最多可选"
              prop="max_choices"
              v-if="formData.vote_type === 'multiple'"
            >
              <el-input-number
                v-model="formData.max_choices"
                :min="2"
                :max="10"
                style="width: 100%"
              >
              </el-input-number>
              <span style="margin-left: 8px; color: #909399;">项</span>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="投票选项" prop="options">
          <div class="options-list">
            <div
              v-for="(option, index) in formData.options"
              :key="index"
              class="option-item"
            >
              <span class="option-index">{{ index + 1 }}</span>
              <el-input
                v-model="formData.options[index]"
                :placeholder="`选项 ${index + 1}`"
              >
              </el-input>
              <el-button
                type="danger"
                :icon="Delete"
                circle
                size="small"
                @click="removeOption(index)"
                :disabled="formData.options.length <= 2"
              >
              </el-button>
            </div>
            <el-button type="primary" :icon="Plus" @click="addOption">
              添加选项
            </el-button>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" class="btn-primary" @click="handleSubmit" :loading="submitting">
            <el-icon><Check /></el-icon>
            创建活动
          </el-button>
          <el-button @click="handleReset">
            <el-icon><RefreshRight /></el-icon>
            重置表单
          </el-button>
          <el-button @click="goBack">
            <el-icon><Back /></el-icon>
            返回列表
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Delete, Check, RefreshRight, Back } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createActivity } from '../api/activity'

const router = useRouter()
const formRef = ref(null)
const submitting = ref(false)

const formData = reactive({
  title: '',
  description: '',
  start_time: '',
  end_time: '',
  vote_type: 'single',
  max_choices: 2,
  options: ['', '']
})

const validateOptions = (rule, value, callback) => {
  const validOptions = value.filter(opt => opt.trim() !== '')
  if (validOptions.length < 2) {
    callback(new Error('至少需要2个有效选项'))
  } else {
    callback()
  }
}

const validateEndTime = (rule, value, callback) => {
  if (!value) {
    callback(new Error('请选择结束时间'))
  } else if (formData.start_time && new Date(value) <= new Date(formData.start_time)) {
    callback(new Error('结束时间必须晚于开始时间'))
  } else {
    callback()
  }
}

const rules = {
  title: [
    { required: true, message: '请输入活动标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在2到100个字符', trigger: 'blur' }
  ],
  start_time: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  end_time: [
    { required: true, validator: validateEndTime, trigger: 'change' }
  ],
  vote_type: [
    { required: true, message: '请选择投票类型', trigger: 'change' }
  ],
  options: [
    { required: true, validator: validateOptions, trigger: 'blur' }
  ]
}

const addOption = () => {
  if (formData.options.length < 10) {
    formData.options.push('')
  } else {
    ElMessage.warning('最多添加10个选项')
  }
}

const removeOption = (index) => {
  if (formData.options.length > 2) {
    formData.options.splice(index, 1)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    const validOptions = formData.options.filter(opt => opt.trim() !== '')
    if (validOptions.length < 2) {
      ElMessage.error('至少需要2个有效选项')
      return
    }

    await ElMessageBox.confirm(
      '确认创建该投票活动吗？',
      '创建确认',
      {
        confirmButtonText: '确认创建',
        cancelButtonText: '取消',
        type: 'info'
      }
    )

    submitting.value = true

    const data = {
      ...formData,
      options: validOptions,
      max_choices: formData.vote_type === 'single' ? 1 : formData.max_choices
    }

    await createActivity(data)

    ElMessage.success('活动创建成功！')
    router.push('/')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('创建失败：' + (error.message || '未知错误'))
    }
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  formData.options = ['', '']
  formData.max_choices = 2
  formData.vote_type = 'single'
}

const goBack = () => {
  router.push('/')
}
</script>

<style scoped>
.create-activity {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title .el-icon {
  color: #667eea;
}

.form-card {
  padding: 30px;
}

.options-list {
  width: 100%;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.option-index {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.option-item .el-input {
  flex: 1;
}
</style>

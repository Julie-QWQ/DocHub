<template>
  <div v-if="dialogVisible" class="dialog-overlay" @click.self="handleClose">
    <div class="dialog">
      <div class="dialog-header">
        <h2>审核学委申请</h2>
        <button class="close-btn" @click="handleClose">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <div class="dialog-body">
        <form @submit.prevent="handleSubmit">
          <div class="info-section">
            <div class="info-row">
              <label class="info-label">申请人</label>
              <div class="info-value">{{ application?.user?.real_name || application?.user?.username }}</div>
            </div>

            <div class="info-row">
              <label class="info-label">申请理由</label>
              <div class="info-value text-content">{{ application?.reason }}</div>
            </div>

            <div class="info-row">
              <label class="info-label">申请时间</label>
              <div class="info-value">{{ formatDate(application?.created_at) }}</div>
            </div>
          </div>

          <div class="divider"></div>

          <div class="form-section">
            <div class="form-group">
              <label class="form-label">审核结果</label>
              <div class="radio-group">
                <label class="radio-option" :class="{ selected: formData.approved === true }">
                  <input type="radio" v-model="formData.approved" :value="true" />
                  <span class="radio-text">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                    通过
                  </span>
                </label>
                <label class="radio-option" :class="{ selected: formData.approved === false }">
                  <input type="radio" v-model="formData.approved" :value="false" />
                  <span class="radio-text">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <line x1="18" y1="6" x2="6" y2="18"></line>
                      <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                    拒绝
                  </span>
                </label>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">
                {{ formData.approved ? '审核备注' : '拒绝原因' }}
              </label>
              <textarea
                v-model="formData.comment"
                class="form-textarea"
                :rows="formData.approved ? 3 : 4"
                :placeholder="formData.approved ? '请输入审核备注(选填)' : '请输入拒绝原因(必填)'"
                maxlength="500"
              ></textarea>
              <div class="char-count">{{ formData.comment.length }}/500</div>
            </div>
          </div>

          <div class="dialog-footer">
            <button type="button" class="btn btn-secondary" @click="handleClose">
              取消
            </button>
            <button type="submit" class="btn btn-primary" :disabled="loading">
              {{ loading ? '提交中...' : '确认审核' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useCommitteeStore } from '@/stores/committee'
import type { CommitteeApplication } from '@/types'
import { formatDate } from '@/utils/format'

interface Props {
  modelValue: boolean
  application: CommitteeApplication | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const committeeStore = useCommitteeStore()

const dialogVisible = ref(false)
const loading = ref(false)

const formData = reactive({
  approved: true,
  comment: ''
})

watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
  }
)

watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})

const handleClose = () => {
  dialogVisible.value = false
  formData.approved = true
  formData.comment = ''
}

const handleSubmit = async () => {
  if (!props.application) {
    ElMessage.error('申请信息不存在')
    return
  }

  // 验证: 如果是拒绝,必须填写原因
  if (!formData.approved && !formData.comment.trim()) {
    ElMessage.error('拒绝时必须填写原因')
    return
  }

  loading.value = true
  try {
    await committeeStore.reviewApplication(
      props.application!.id,
      {
        approved: formData.approved,
        comment: formData.comment
      }
    )

    ElMessage.success('审核完成')
    emit('success')
    handleClose()
  } catch (error: any) {
    ElMessage.error(error.message || '审核失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.dialog {
  background: #ffffff;
  border: 1px solid #000000;
  border-radius: 8px;
  width: 100%;
  max-width: 560px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e5e5;

  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #000000;
  }

  .close-btn {
    background: none;
    border: none;
    padding: 4px;
    cursor: pointer;
    color: #666666;
    transition: color 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
      color: #000000;
    }
  }
}

.dialog-body {
  padding: 24px;
}

.info-section {
  margin-bottom: 20px;
}

.info-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }
}

.info-label {
  min-width: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #666666;
  margin-top: 2px;
}

.info-value {
  flex: 1;
  font-size: 14px;
  color: #000000;
  line-height: 1.5;

  &.text-content {
    white-space: pre-wrap;
    word-break: break-word;
  }
}

.divider {
  height: 1px;
  background: #e5e5e5;
  margin: 20px 0;
}

.form-section {
  margin-top: 20px;
}

.form-group {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #000000;
  margin-bottom: 8px;
}

.radio-group {
  display: flex;
  gap: 12px;
}

.radio-option {
  position: relative;
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  background: #ffffff;

  input[type="radio"] {
    position: absolute;
    opacity: 0;
    pointer-events: none;
  }

  .radio-text {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 14px;
    color: #333333;
    font-weight: 500;

    svg {
      flex-shrink: 0;
    }
  }

  &:hover {
    border-color: #999999;
    background: #fafafa;
  }

  &.selected {
    border-color: #000000;
    background: #f5f5f5;

    .radio-text {
      color: #000000;
    }
  }
}

.form-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #cccccc;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  line-height: 1.5;
  resize: vertical;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: #000000;
  }

  &::placeholder {
    color: #999999;
  }
}

.char-count {
  margin-top: 6px;
  font-size: 12px;
  color: #999999;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #e5e5e5;
}

.btn {
  padding: 10px 20px;
  border: 1px solid #000000;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  background: #ffffff;
  color: #000000;

  &:hover:not(:disabled) {
    background: #f5f5f5;
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &.btn-secondary {
    border-color: #cccccc;
    color: #666666;

    &:hover:not(:disabled) {
      border-color: #999999;
      background: #fafafa;
    }
  }

  &.btn-primary {
    background: #000000;
    color: #ffffff;

    &:hover:not(:disabled) {
      background: #333333;
    }
  }
}
</style>

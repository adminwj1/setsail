<template>
  <div class="requirement-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>需求列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新建需求
          </el-button>
        </div>
      </template>

      <!-- 筛选 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="产品">
          <el-select v-model="filters.product_id" clearable placeholder="请选择产品" style="width: 200px">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.type" clearable placeholder="请选择类型" style="width: 150px">
            <el-option :value="1" label="业务需求" />
            <el-option :value="2" label="用户需求" />
            <el-option :value="3" label="研发需求" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" clearable placeholder="请选择状态" style="width: 150px">
            <el-option :value="0" label="待评审" />
            <el-option :value="1" label="已采纳" />
            <el-option :value="2" label="开发中" />
            <el-option :value="3" label="已完成" />
            <el-option :value="4" label="已拒绝" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="需求标题" show-overflow-tooltip />
        <el-table-column prop="product_name" label="所属产品" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag>{{ getTypeName(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)">
              {{ getPriorityName(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creator_nickname" label="创建人" width="120" />
        <el-table-column prop="assignee_nickname" label="负责人" width="120" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </el-card>

    <!-- 创建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" @close="handleDialogClose">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="所属产品" prop="product_id">
          <el-select v-model="form.product_id" placeholder="请选择产品" style="width: 100%">
            <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="需求标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入需求标题" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型" style="width: 100%">
            <el-option :value="1" label="业务需求" />
            <el-option :value="2" label="用户需求" />
            <el-option :value="3" label="研发需求" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="form.priority" placeholder="请选择优先级" style="width: 100%">
            <el-option :value="1" label="高" />
            <el-option :value="2" label="中" />
            <el-option :value="3" label="低" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
            <el-option :value="0" label="待评审" />
            <el-option :value="1" label="已采纳" />
            <el-option :value="2" label="开发中" />
            <el-option :value="3" label="已完成" />
            <el-option :value="4" label="已拒绝" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRequirementList, createRequirement, updateRequirement, deleteRequirement } from '@/api/requirement'
import { getProductList } from '@/api/product'

const loading = ref(false)
const tableData = ref([])
const products = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('新建需求')
const isEdit = ref(false)
const formRef = ref()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const filters = reactive({
  product_id: null,
  type: null,
  status: null
})

const form = reactive({
  id: null,
  product_id: null,
  title: '',
  type: 1,
  priority: 2,
  status: 0,
  description: ''
})

const rules = {
  product_id: [{ required: true, message: '请选择产品', trigger: 'change' }],
  title: [{ required: true, message: '请输入需求标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const typeMap = { 1: '业务需求', 2: '用户需求', 3: '研发需求' }
const priorityMap = { 1: { name: '高', type: 'danger' }, 2: { name: '中', type: 'warning' }, 3: { name: '低', type: 'info' } }
const statusMap = {
  0: { name: '待评审', type: 'info' },
  1: { name: '已采纳', type: 'success' },
  2: { name: '开发中', type: 'warning' },
  3: { name: '已完成', type: 'success' },
  4: { name: '已拒绝', type: 'danger' }
}

const getTypeName = (type) => typeMap[type] || '未知'
const getPriorityName = (p) => priorityMap[p]?.name || '未知'
const getPriorityType = (p) => priorityMap[p]?.type || 'info'
const getStatusName = (s) => statusMap[s]?.name || '未知'
const getStatusType = (s) => statusMap[s]?.type || 'info'

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      product_id: filters.product_id || undefined,
      type: filters.type || undefined,
      status: filters.status ?? undefined
    }
    const res = await getRequirementList(params)
    tableData.value = (res.data.list || []).map(r => ({
      ...r,
      product_name: r.product?.name || '-'
    }))
    pagination.total = res.data.total || 0
  } catch (error) {
    console.error('Failed to fetch requirements:', error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res = await getProductList({ page: 1, pageSize: 100 })
    products.value = res.data.list || []
  } catch (error) {
    console.error('Failed to fetch products:', error)
  }
}

const handleFilter = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  Object.assign(filters, { product_id: null, type: null, status: null })
  handleFilter()
}

const handleCreate = () => {
  dialogTitle.value = '新建需求'
  isEdit.value = false
  Object.assign(form, { id: null, product_id: null, title: '', type: 1, priority: 2, status: 0, description: '' })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑需求'
  isEdit.value = true
  Object.assign(form, {
    id: row.id,
    product_id: row.product_id,
    title: row.title,
    type: row.type,
    priority: row.priority,
    status: row.status,
    description: row.description
  })
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定删除该需求吗？', '提示', { type: 'warning' })
    await deleteRequirement(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') console.error('Failed to delete:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (isEdit.value) {
          await updateRequirement(form.id, form)
          ElMessage.success('更新成功')
        } else {
          await createRequirement(form)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        console.error('Failed to submit:', error)
      }
    }
  })
}

const handleDialogClose = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchData()
  fetchProducts()
})
</script>

<style scoped lang="scss">
.requirement-list {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .filter-form {
    margin-bottom: 16px;
  }
}
</style>

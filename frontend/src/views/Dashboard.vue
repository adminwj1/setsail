<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #409eff">
              <el-icon><Goods /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.productCount }}</div>
              <div class="stat-label">产品数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #67c23a">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.requirementCount }}</div>
              <div class="stat-label">需求数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #e6a23c">
              <el-icon><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">用户数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-icon" style="background: #f56c6c">
              <el-icon><Key /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.roleCount }}</div>
              <div class="stat-label">角色数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新需求</span>
          </template>
          <el-table :data="recentRequirements" style="width: 100%">
            <el-table-column prop="title" label="需求标题" />
            <el-table-column prop="type_name" label="类型" width="100" />
            <el-table-column prop="status_name" label="状态" width="100" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新产品</span>
          </template>
          <el-table :data="recentProducts" style="width: 100%">
            <el-table-column prop="name" label="产品名称" />
            <el-table-column prop="code" label="产品代号" width="120" />
            <el-table-column prop="status_name" label="状态" width="100" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProductList } from '@/api/product'
import { getRequirementList } from '@/api/requirement'

const stats = ref({
  productCount: 0,
  requirementCount: 0,
  userCount: 0,
  roleCount: 0
})

const recentProducts = ref([])
const recentRequirements = ref([])

const statusMap = {
  0: { name: '规划中', class: 'info' },
  1: { name: '开发中', class: 'warning' },
  2: { name: '上线', class: 'success' },
  3: { name: '下线', class: 'danger' }
}

const typeMap = {
  1: '业务需求',
  2: '用户需求',
  3: '研发需求'
}

const reqStatusMap = {
  0: '待评审',
  1: '已采纳',
  2: '开发中',
  3: '已完成',
  4: '已拒绝'
}

onMounted(async () => {
  try {
    const [productRes, reqRes] = await Promise.all([
      getProductList({ page: 1, pageSize: 5 }),
      getRequirementList({ page: 1, pageSize: 5 })
    ])

    stats.value.productCount = productRes.data.total
    recentProducts.value = (productRes.data.list || []).map(p => ({
      ...p,
      status_name: statusMap[p.status]?.name || '未知'
    }))

    stats.value.requirementCount = reqRes.data.total
    recentRequirements.value = (reqRes.data.list || []).map(r => ({
      ...r,
      type_name: typeMap[r.type] || '未知',
      status_name: reqStatusMap[r.status] || '未知'
    }))
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  }
})
</script>

<style scoped lang="scss">
.dashboard {
  padding: 20px;
}

.stat-card {
  display: flex;
  align-items: center;

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 16px;

    .el-icon {
      font-size: 28px;
      color: #fff;
    }
  }

  .stat-info {
    .stat-value {
      font-size: 24px;
      font-weight: bold;
      color: #303133;
    }

    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 4px;
    }
  }
}
</style>

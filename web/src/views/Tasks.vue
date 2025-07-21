<template>
  <div class="tasks-container">
    <div class="mb-4 flex justify-between items-center">
      <h1 class="text-2xl font-bold">任务管理</h1>
      <n-button @click="refreshTasks" secondary>
        <template #icon>
          <div class="i-heroicons:arrow-path"></div>
        </template>
        刷新
      </n-button>
    </div>

    <n-card>
      <n-data-table
        :columns="columns"
        :data="tasks"
        :loading="loading"
        :pagination="pagination"
        :bordered="false"
        striped
      />
    </n-card>
  </div>
</template>

<script setup lang="tsx">
  import { ref, onMounted, h } from 'vue';
  import { useMessage } from 'naive-ui';
  import { NDataTable, NButton, NCard } from 'naive-ui';
  import { taskApi } from '@/api';
  import type { Task } from '@/types';

  const tasks = ref<Task[]>([]);
  const loading = ref(true);
  const message = useMessage();

  const pagination = {
    pageSize: 10,
  };

  // 格式化时间
  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }

  // 获取状态标签类型
  function getStatusType(status: string) {
    switch (status) {
      case 'pending':
        return 'default';
      case 'running':
        return 'info';
      case 'paused':
        return 'warning';
      case 'completed':
        return 'success';
      case 'failed':
        return 'error';
      default:
        return 'default';
    }
  }

  // 获取状态中文名称
  function getStatusName(status: string) {
    switch (status) {
      case 'pending':
        return '等待中';
      case 'running':
        return '运行中';
      case 'paused':
        return '已暂停';
      case 'completed':
        return '已完成';
      case 'failed':
        return '失败';
      default:
        return '未知';
    }
  }

  // 停止任务
  async function stopTask(id: string) {
    try {
      await taskApi.stopTask(id);
      message.success('任务已停止');
      await fetchTasks();
    } catch (error: any) {
      message.error(error.message || '停止任务失败');
    }
  }

  const columns = [
    {
      title: '任务名称',
      key: 'name',
      width: 200,
    },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render(row: Task) {
        return (
          <n-tag type={getStatusType(row.status)} size="small" round>
            {getStatusName(row.status)}
          </n-tag>
        );
      },
    },
    {
      title: '进度',
      key: 'progress',
      width: 150,
      render(row: Task) {
        return (
          <n-progress
            type="line"
            percentage={row.progress}
            indicatorPlacement="inside"
            status={row.status === 'failed' ? 'error' : 'default'}
            showIndicator
            height={12}
          />
        );
      },
    },
    {
      title: '结果',
      key: 'result',
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: '错误信息',
      key: 'error',
      ellipsis: {
        tooltip: true,
      },
      render(row: Task) {
        return row.error ? <span style="color: #f5222d">{row.error}</span> : '';
      },
    },
    {
      title: '创建时间',
      key: 'createdAt',
      width: 180,
      render(row: Task) {
        return formatDate(row.createdAt);
      },
    },
    {
      title: '更新时间',
      key: 'updatedAt',
      width: 180,
      render(row: Task) {
        return formatDate(row.updatedAt);
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render(row: Task) {
        if (row.status === 'running' || row.status === 'pending') {
          return (
            <n-popconfirm
              onPositiveClick={() => stopTask(row.id)}
              negativeText="取消"
              positiveText="确定"
            >
              {{
                trigger: () => (
                  <n-button size="small" type="warning">
                    停止
                  </n-button>
                ),
                default: () => '确定要停止该任务吗？',
              }}
            </n-popconfirm>
          );
        }
        return null;
      },
    },
  ];

  // 获取任务列表
  async function fetchTasks() {
    loading.value = true;
    try {
      tasks.value = await taskApi.fetchTasks();
    } catch (error: any) {
      message.error(error.message || '获取任务列表失败');
    } finally {
      loading.value = false;
    }
  }

  // 刷新任务列表
  function refreshTasks() {
    fetchTasks();
  }

  onMounted(() => {
    fetchTasks();
  });
</script>

<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="overflow-hidden rounded-2xl border border-primary-100 bg-gradient-to-br from-primary-600 via-blue-600 to-indigo-700 p-6 text-white shadow-sm dark:border-primary-900/50">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <p class="text-sm font-medium text-primary-100">Model Marketplace</p>
            <h1 class="mt-1 text-2xl font-bold">模型广场</h1>
            <p class="mt-2 max-w-3xl text-sm text-primary-50/90">
              查看管理员后台已开放、且当前账号可用的模型列表、模型 ID、供应商、分组和计费价格。复制模型 ID 后即可在 OpenAI 兼容客户端中调用。
            </p>
          </div>
          <button
            @click="loadChannels"
            :disabled="loading"
            class="inline-flex items-center justify-center gap-2 rounded-xl bg-white/15 px-4 py-2 text-sm font-medium text-white backdrop-blur transition hover:bg-white/25 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
            刷新
          </button>
        </div>

        <div class="mt-6 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-xl bg-white/12 p-4 backdrop-blur">
            <div class="text-xs text-primary-100">可用模型</div>
            <div class="mt-1 text-2xl font-semibold">{{ modelCards.length }}</div>
          </div>
          <div class="rounded-xl bg-white/12 p-4 backdrop-blur">
            <div class="text-xs text-primary-100">供应商</div>
            <div class="mt-1 text-2xl font-semibold">{{ platformOptions.length }}</div>
          </div>
          <div class="rounded-xl bg-white/12 p-4 backdrop-blur">
            <div class="text-xs text-primary-100">可用分组</div>
            <div class="mt-1 text-2xl font-semibold">{{ groupCount }}</div>
          </div>
          <div class="rounded-xl bg-white/12 p-4 backdrop-blur">
            <div class="text-xs text-primary-100">渠道来源</div>
            <div class="mt-1 text-2xl font-semibold">{{ channelCount }}</div>
          </div>
        </div>
      </div>

      <div class="card p-4">
        <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <div class="relative w-full lg:max-w-md">
            <Icon
              name="search"
              size="md"
              class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
            />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索模型 ID、供应商、分组或渠道..."
              class="input pl-10"
            />
          </div>

          <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <select v-model="selectedPlatform" class="input min-w-[160px]">
              <option value="all">全部供应商</option>
              <option v-for="platform in platformOptions" :key="platform" :value="platform">
                {{ platformLabel(platform) }}
              </option>
            </select>

            <select v-model="selectedBillingMode" class="input min-w-[160px]">
              <option value="all">全部计费类型</option>
              <option value="token">按量计费</option>
              <option value="per_request">按次计费</option>
              <option value="image">图片计费</option>
              <option value="none">暂无价格</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="loading" class="card py-16 text-center">
        <Icon name="refresh" size="lg" class="mx-auto animate-spin text-gray-400" />
        <p class="mt-3 text-sm text-gray-500 dark:text-gray-400">正在加载模型广场...</p>
      </div>

      <div v-else-if="filteredModels.length === 0" class="card py-16 text-center">
        <Icon name="inbox" size="xl" class="mx-auto mb-3 h-12 w-12 text-gray-400" />
        <p class="text-sm text-gray-500 dark:text-gray-400">暂无匹配模型</p>
      </div>

      <div v-else class="grid gap-4 lg:grid-cols-2 xl:grid-cols-3">
        <article
          v-for="model in filteredModels"
          :key="model.key"
          class="group flex h-full flex-col overflow-hidden rounded-2xl border bg-white shadow-sm transition hover:-translate-y-0.5 hover:shadow-md dark:bg-dark-900"
          :class="platformBorderClass(model.platform)"
        >
          <div class="flex items-start justify-between gap-3 border-b border-gray-100 p-4 dark:border-dark-700">
            <div class="min-w-0">
              <div class="mb-2 flex flex-wrap items-center gap-2">
                <span
                  class="inline-flex items-center rounded-md border px-2 py-0.5 text-[11px] font-semibold uppercase"
                  :class="platformBadgeClass(model.platform)"
                >
                  {{ platformLabel(model.platform) }}
                </span>
                <span class="rounded-md bg-gray-100 px-2 py-0.5 text-[11px] text-gray-600 dark:bg-dark-700 dark:text-gray-300">
                  {{ billingModeLabel(model.pricing?.billing_mode ?? null) }}
                </span>
              </div>
              <h2 class="truncate text-lg font-semibold text-gray-900 dark:text-white" :title="model.name">
                {{ model.name }}
              </h2>
              <p class="mt-1 line-clamp-2 text-xs text-gray-500 dark:text-gray-400">
                {{ modelDescription(model) }}
              </p>
            </div>

            <button
              type="button"
              class="flex-shrink-0 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 transition hover:border-primary-300 hover:text-primary-600 dark:border-dark-600 dark:text-gray-300 dark:hover:border-primary-500 dark:hover:text-primary-400"
              @click="copyModelId(model.name)"
            >
              {{ copiedModel === model.name ? '已复制' : '复制 ID' }}
            </button>
          </div>

          <div class="flex flex-1 flex-col gap-4 p-4">
            <div class="rounded-xl bg-gray-50 p-3 dark:bg-dark-800/70">
              <div class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">模型 ID</div>
              <code class="break-all rounded-lg bg-white px-2 py-1 text-sm font-semibold text-gray-900 dark:bg-dark-900 dark:text-gray-100">
                {{ model.name }}
              </code>
            </div>

            <div class="grid gap-2 text-sm sm:grid-cols-2">
              <div
                v-for="row in priceRows(model)"
                :key="row.label"
                class="rounded-xl border border-gray-100 p-3 dark:border-dark-700"
              >
                <div class="text-xs text-gray-500 dark:text-gray-400">{{ row.label }}</div>
                <div class="mt-1 font-semibold text-gray-900 dark:text-white">{{ row.value }}</div>
              </div>
            </div>

            <div>
              <div class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">可用分组</div>
              <div class="flex flex-wrap gap-1.5">
                <span
                  v-for="group in model.groups"
                  :key="group.id"
                  class="inline-flex items-center gap-1 rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-700 dark:bg-dark-700 dark:text-gray-300"
                  :title="`费率 ${group.rateMultiplier}x${group.isExclusive ? ' · 专属分组' : ''}`"
                >
                  <span v-if="group.isExclusive">🔒</span>
                  {{ group.name }}
                  <span class="text-gray-400">{{ group.rateMultiplier }}x</span>
                </span>
              </div>
            </div>

            <div>
              <div class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">渠道来源</div>
              <div class="flex flex-wrap gap-1.5">
                <span
                  v-for="channel in model.channels"
                  :key="channel"
                  class="rounded-md border border-gray-200 px-2 py-1 text-xs text-gray-600 dark:border-dark-600 dark:text-gray-300"
                >
                  {{ channel }}
                </span>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import userChannelsAPI, {
  type UserAvailableChannel,
  type UserSupportedModel,
  type UserSupportedModelPricing,
} from '@/api/channels'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { formatScaled } from '@/utils/pricing'
import { platformBadgeClass, platformBorderClass, platformLabel } from '@/utils/platformColors'
import {
  BILLING_MODE_IMAGE,
  BILLING_MODE_PER_REQUEST,
  BILLING_MODE_TOKEN,
  type BillingMode,
} from '@/constants/channel'

interface MarketplaceGroup {
  id: number
  name: string
  rateMultiplier: number
  isExclusive: boolean
}

interface MarketplaceModel {
  key: string
  name: string
  platform: string
  pricing: UserSupportedModelPricing | null
  channels: string[]
  groups: MarketplaceGroup[]
}

const { t } = useI18n()
const appStore = useAppStore()

const channels = ref<UserAvailableChannel[]>([])
const loading = ref(false)
const searchQuery = ref('')
const selectedPlatform = ref('all')
const selectedBillingMode = ref<'all' | BillingMode | 'none'>('all')
const copiedModel = ref('')

const modelCards = computed<MarketplaceModel[]>(() => {
  const models = new Map<string, MarketplaceModel>()

  for (const channel of channels.value) {
    for (const section of channel.platforms) {
      for (const model of section.supported_models) {
        const key = `${section.platform}:${model.name}`
        const existing = models.get(key) ?? createMarketplaceModel(key, section.platform, model)

        if (!existing.channels.includes(channel.name)) {
          existing.channels.push(channel.name)
        }

        for (const group of section.groups) {
          if (!existing.groups.some((g) => g.id === group.id)) {
            existing.groups.push({
              id: group.id,
              name: group.name,
              rateMultiplier: group.rate_multiplier,
              isExclusive: group.is_exclusive,
            })
          }
        }

        models.set(key, existing)
      }
    }
  }

  return Array.from(models.values()).sort((a, b) => {
    const platformCompare = a.platform.localeCompare(b.platform)
    return platformCompare !== 0 ? platformCompare : a.name.localeCompare(b.name)
  })
})

const platformOptions = computed(() =>
  Array.from(new Set(modelCards.value.map((model) => model.platform))).sort(),
)

const groupCount = computed(() => {
  const ids = new Set<number>()
  for (const model of modelCards.value) {
    model.groups.forEach((group) => ids.add(group.id))
  }
  return ids.size
})

const channelCount = computed(() => new Set(channels.value.map((channel) => channel.name)).size)

const filteredModels = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return modelCards.value.filter((model) => {
    const platformHit = selectedPlatform.value === 'all' || model.platform === selectedPlatform.value
    const billingMode = model.pricing?.billing_mode ?? 'none'
    const billingHit = selectedBillingMode.value === 'all' || billingMode === selectedBillingMode.value
    const textHit =
      !q ||
      model.name.toLowerCase().includes(q) ||
      platformLabel(model.platform).toLowerCase().includes(q) ||
      model.platform.toLowerCase().includes(q) ||
      model.channels.some((channel) => channel.toLowerCase().includes(q)) ||
      model.groups.some((group) => group.name.toLowerCase().includes(q))

    return platformHit && billingHit && textHit
  })
})

function createMarketplaceModel(key: string, platform: string, model: UserSupportedModel): MarketplaceModel {
  return {
    key,
    name: model.name,
    platform: model.platform || platform,
    pricing: model.pricing,
    channels: [],
    groups: [],
  }
}

function billingModeLabel(mode: BillingMode | null): string {
  switch (mode) {
    case BILLING_MODE_TOKEN:
      return '按量计费'
    case BILLING_MODE_PER_REQUEST:
      return '按次计费'
    case BILLING_MODE_IMAGE:
      return '图片计费'
    default:
      return '暂无价格'
  }
}

function priceRows(model: MarketplaceModel): Array<{ label: string; value: string }> {
  const pricing = model.pricing
  if (!pricing) {
    return [{ label: '价格', value: '暂无价格' }]
  }

  if (pricing.billing_mode === BILLING_MODE_PER_REQUEST) {
    return [{ label: '单次价格', value: `${formatScaled(pricing.per_request_price, 1)} / 次` }]
  }

  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return [{ label: '图片价格', value: `${formatScaled(pricing.image_output_price, 1)} / 次` }]
  }

  return [
    { label: '输入价格', value: `${formatScaled(pricing.input_price, 1_000_000)} / 1M Tokens` },
    { label: '输出价格', value: `${formatScaled(pricing.output_price, 1_000_000)} / 1M Tokens` },
    { label: '缓存读取', value: `${formatScaled(pricing.cache_read_price, 1_000_000)} / 1M Tokens` },
    ...(pricing.image_output_price != null && pricing.image_output_price > 0
      ? [{ label: '图片输出', value: `${formatScaled(pricing.image_output_price, 1_000_000)} / 1M Tokens` }]
      : []),
  ]
}

function modelDescription(model: MarketplaceModel): string {
  const name = model.name.toLowerCase()
  if (name.includes('image')) return '图像生成 / 多模态模型，适合生图、图像理解和素材生成。'
  if (name.includes('mini') || name.includes('flash')) return '轻量高速模型，适合高频问答、摘要、分类和低成本调用。'
  if (name.includes('code') || name.includes('codex')) return '编程与代码任务模型，适合代码生成、调试、重构和自动化开发。'
  if (name.includes('opus') || name.includes('gpt-5') || name.includes('sonnet')) return '高性能通用模型，适合复杂问答、写作、分析和多步任务。'
  return `${platformLabel(model.platform)} 模型，适合通用对话、内容生成和工具调用。`
}

async function copyModelId(modelId: string) {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(modelId)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = modelId
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    copiedModel.value = modelId
    window.setTimeout(() => {
      if (copiedModel.value === modelId) copiedModel.value = ''
    }, 1600)
  } catch (err) {
    console.error('Failed to copy model id:', err)
  }
}

async function loadChannels() {
  loading.value = true
  try {
    channels.value = await userChannelsAPI.getAvailable()
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  } finally {
    loading.value = false
  }
}

onMounted(loadChannels)
</script>

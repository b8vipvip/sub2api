<template>
  <AppLayout>
    <div class="grid gap-6 lg:grid-cols-[260px_1fr]">
      <aside class="space-y-6">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-bold text-gray-900 dark:text-white">筛选</h2>
          <button class="btn btn-secondary btn-sm" @click="resetFilters">重置</button>
        </div>

        <section class="space-y-3">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">供应商</h3>
          <button
            class="filter-pill w-full justify-between"
            :class="selectedPlatform === 'all' ? 'filter-pill-active' : ''"
            @click="selectedPlatform = 'all'"
          >
            <span>全部供应商</span>
            <span>{{ modelCards.length }}</span>
          </button>
          <button
            v-for="platform in platformOptions"
            :key="platform"
            class="filter-pill w-full justify-between"
            :class="selectedPlatform === platform ? 'filter-pill-active' : ''"
            @click="selectedPlatform = platform"
          >
            <span>{{ platformLabel(platform) }}</span>
            <span>{{ countByPlatform(platform) }}</span>
          </button>
        </section>

        <section class="space-y-3">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">可用令牌分组</h3>
          <button
            class="filter-pill w-full justify-between"
            :class="selectedGroup === 'all' ? 'filter-pill-active' : ''"
            @click="selectedGroup = 'all'"
          >
            <span>全部分组</span>
            <span>{{ groupOptions.length }}</span>
          </button>
          <button
            v-for="group in groupOptions"
            :key="group"
            class="filter-pill w-full justify-between"
            :class="selectedGroup === group ? 'filter-pill-active' : ''"
            @click="selectedGroup = group"
          >
            <span>{{ group }}</span>
            <span>{{ countByGroup(group) }}</span>
          </button>
        </section>

        <section class="space-y-3">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">计费类型</h3>
          <button
            class="filter-pill w-full justify-between"
            :class="selectedBillingMode === 'all' ? 'filter-pill-active' : ''"
            @click="selectedBillingMode = 'all'"
          >
            <span>全部类型</span>
            <span>{{ modelCards.length }}</span>
          </button>
          <button
            v-for="item in billingOptions"
            :key="item.value"
            class="filter-pill w-full justify-between"
            :class="selectedBillingMode === item.value ? 'filter-pill-active' : ''"
            @click="selectedBillingMode = item.value"
          >
            <span>{{ item.label }}</span>
            <span>{{ countByBilling(item.value) }}</span>
          </button>
        </section>

        <section class="space-y-3">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">端点类型</h3>
          <button class="filter-pill w-full justify-between filter-pill-active">
            <span>全部端点</span>
            <span>{{ modelCards.length }}</span>
          </button>
          <button class="filter-pill w-full justify-between">
            <span>openai</span>
            <span>{{ modelCards.length }}</span>
          </button>
        </section>
      </aside>

      <main class="space-y-4">
        <div class="overflow-hidden rounded-2xl border border-primary-100 bg-gradient-to-br from-primary-600 via-blue-600 to-indigo-700 p-6 text-white shadow-sm dark:border-primary-900/50">
          <div class="flex items-center justify-between gap-4">
            <div>
              <div class="flex items-center gap-3">
                <h1 class="text-2xl font-bold">全部供应商</h1>
                <span class="rounded-full bg-white/20 px-3 py-1 text-xs font-medium">共 {{ modelCards.length }} 个模型</span>
              </div>
              <p class="mt-2 max-w-3xl text-sm text-primary-50/90">
                查看管理员后台调度开关已开启账号的模型白名单。这里展示的模型 ID 可直接用于 OpenAI 兼容接口调用。
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
        </div>

        <div class="card p-4">
          <div class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
            <div class="relative w-full xl:max-w-2xl">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="模糊搜索模型名称"
                class="input pl-10"
              />
            </div>
            <div class="flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
              <button class="btn btn-secondary btn-sm" @click="copyVisibleModelIds">复制</button>
              <span>充值价格显示</span>
              <span class="inline-flex h-5 w-9 items-center rounded-full bg-gray-200 p-0.5 dark:bg-dark-700"><span class="h-4 w-4 rounded-full bg-white shadow" /></span>
              <span>倍率</span>
              <span class="inline-flex h-5 w-9 items-center rounded-full bg-gray-200 p-0.5 dark:bg-dark-700"><span class="h-4 w-4 rounded-full bg-white shadow" /></span>
              <span class="rounded-md bg-primary-50 px-2 py-1 font-medium text-primary-700 dark:bg-primary-950/30 dark:text-primary-300">表格视图</span>
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
          <p class="mt-2 text-xs text-gray-400 dark:text-gray-500">请检查管理员后台账号是否为“正常”状态、调度开关是否开启、模型白名单是否已填写。</p>
        </div>

        <div v-else class="grid gap-4 xl:grid-cols-2 2xl:grid-cols-3">
          <article
            v-for="model in filteredModels"
            :key="model.key"
            class="group relative flex min-h-[220px] flex-col overflow-hidden rounded-2xl border border-gray-100 bg-white p-4 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md dark:border-dark-700 dark:bg-dark-900"
          >
            <div class="flex items-start gap-3">
              <div class="flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-xl bg-gray-100 text-xl dark:bg-dark-700">
                {{ providerIcon(model.platform) }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex min-w-0 items-center gap-2">
                  <h2 class="truncate text-lg font-semibold text-gray-900 dark:text-white" :title="model.name">
                    {{ model.name }}
                  </h2>
                  <button
                    type="button"
                    class="flex-shrink-0 rounded-lg p-1.5 text-gray-400 transition hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400"
                    :title="copiedModel === model.name ? '已复制' : '复制模型id'"
                    aria-label="复制模型id"
                    @click.stop="copyModelId(model.name)"
                  >
                    <Icon name="copy" size="sm" />
                  </button>
                </div>
                <div class="mt-2 space-y-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">
                  <p v-for="row in priceRows(model)" :key="row.label">
                    {{ row.label }} <span class="font-semibold text-gray-900 dark:text-white">{{ row.value }}</span>
                  </p>
                </div>
              </div>
            </div>

            <p class="mt-4 line-clamp-3 text-sm leading-relaxed text-gray-500 dark:text-gray-400">
              {{ modelDescription(model) }}
            </p>

            <div class="mt-auto flex flex-wrap items-center gap-2 pt-4">
              <span class="rounded-full bg-violet-50 px-2.5 py-1 text-xs font-medium text-violet-700 dark:bg-violet-950/30 dark:text-violet-300">
                {{ billingModeLabel(model.pricing?.billing_mode ?? null) }}
              </span>
              <span
                v-for="group in model.groups"
                :key="group.name"
                class="rounded-full bg-gray-100 px-2.5 py-1 text-xs text-gray-600 dark:bg-dark-700 dark:text-gray-300"
              >
                {{ group.name }}
              </span>
            </div>
          </article>
        </div>
      </main>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
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

type BillingFilter = 'all' | BillingMode | 'none'

interface ModelCatalogEntry {
  pricing?: UserSupportedModelPricing
  description: string
}

const { t } = useI18n()
const appStore = useAppStore()

const channels = ref<UserAvailableChannel[]>([])
const loading = ref(false)
const searchQuery = ref('')
const selectedPlatform = ref('all')
const selectedGroup = ref('all')
const selectedBillingMode = ref<BillingFilter>('all')
const copiedModel = ref('')

const billingOptions: Array<{ value: BillingMode | 'none'; label: string }> = [
  { value: BILLING_MODE_TOKEN, label: '按量计费' },
  { value: BILLING_MODE_PER_REQUEST, label: '按次计费' },
  { value: BILLING_MODE_IMAGE, label: '图片计费' },
  { value: 'none', label: '暂无价格' },
]

const MODEL_CATALOG: Record<string, ModelCatalogEntry> = {
  'gpt-5.4': { pricing: pToken(0.025, 0.15, 0.003), description: '高速通用旗舰模型，适合复杂问题和高质量输出。深度问答、内容创作、复杂分析、代码效果最好，推理和表达都很强。' },
  'gpt-5.4-mini': { pricing: pToken(0.007, 0.045, 0.001), description: '高性价比版本，适合大多数用户。日常聊天、问答、文案、轻量代码便宜、速度快、性价比高。' },
  'gpt-5.5': { pricing: pToken(0.05, 0.3, 0.005), description: '高性能通用大模型，适合复杂问答、写作、代码、分析、总结和多步骤任务，综合能力强，适合日常主力使用。' },
  'gpt-image-2': { pricing: pRequest(0.9, BILLING_MODE_IMAGE), description: '生图模型，看文档教程接入。' },
  'gpt-image-2-2k': { pricing: pRequest(0.9, BILLING_MODE_IMAGE), description: '生图模型，看文档教程接入。' },
  'gpt-image-2-4k': { pricing: pRequest(0.9, BILLING_MODE_IMAGE), description: '生图模型，看文档教程接入。' },
  'claude-haiku-4-5-20251001': { pricing: pToken(0.16, 0.8, 0.02, 0.16), description: '轻量快速模型，适合低成本高频使用。简单问答、总结、改写、基础文案、速度快、成本低。' },
  'claude-opus-4-7': { pricing: pToken(0.8, 4.0, 1.0), description: '强力 Claude Opus 模型，适合高难推理、长文分析、复杂代码和严肃写作。' },
  'claude-opus-4-8': { pricing: pToken(0.8, 4.0, 1.0), description: '高性能 Claude Opus 模型，适合深度研究、复杂规划、代码架构和高质量内容生成。' },
  'claude-sonnet-4-6': { pricing: pToken(0.4, 2.4, 0.06, 0.48), description: '更新版 Sonnet，适合大多数通用场景。日常问答、内容生成、分析、代码等综合表现好，适合主力替换。' },
  'codex-auto-review': { pricing: pToken(0.75, 4.5, 0.075), description: '面向代码审查、自动评审和开发辅助的模型，适合代码理解、审查建议、重构和自动化开发流程。' },
  nanobanana: { pricing: pRequest(0.042, BILLING_MODE_PER_REQUEST), description: '谷歌香蕉生图模型，支持 2k-4k。' },
}

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

        const catalog = catalogForModel(model.name)
        existing.pricing = model.pricing ?? catalog?.pricing ?? existing.pricing

        for (const group of section.groups) {
          if (!existing.groups.some((g) => g.id === group.id || g.name === group.name)) {
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

  return Array.from(models.values()).sort((a, b) => a.name.localeCompare(b.name))
})

const platformOptions = computed(() =>
  Array.from(new Set(modelCards.value.map((model) => model.platform))).sort(),
)

const groupOptions = computed(() =>
  Array.from(new Set(modelCards.value.flatMap((model) => model.groups.map((group) => group.name)))).sort(),
)

const filteredModels = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return modelCards.value.filter((model) => {
    const platformHit = selectedPlatform.value === 'all' || model.platform === selectedPlatform.value
    const groupHit = selectedGroup.value === 'all' || model.groups.some((group) => group.name === selectedGroup.value)
    const billingMode = model.pricing?.billing_mode ?? 'none'
    const billingHit = selectedBillingMode.value === 'all' || billingMode === selectedBillingMode.value
    const textHit =
      !q ||
      model.name.toLowerCase().includes(q) ||
      platformLabel(model.platform).toLowerCase().includes(q) ||
      model.platform.toLowerCase().includes(q) ||
      modelDescription(model).toLowerCase().includes(q) ||
      model.channels.some((channel) => channel.toLowerCase().includes(q)) ||
      model.groups.some((group) => group.name.toLowerCase().includes(q))

    return platformHit && groupHit && billingHit && textHit
  })
})

function pToken(input: number, output: number, cacheRead: number, cacheWrite = 0): UserSupportedModelPricing {
  return {
    billing_mode: BILLING_MODE_TOKEN,
    input_price: input,
    output_price: output,
    cache_write_price: cacheWrite,
    cache_read_price: cacheRead,
    image_output_price: null,
    per_request_price: null,
    intervals: [],
  }
}

function pRequest(price: number, mode: BillingMode): UserSupportedModelPricing {
  return {
    billing_mode: mode,
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    image_output_price: mode === BILLING_MODE_IMAGE ? price : null,
    per_request_price: mode === BILLING_MODE_PER_REQUEST ? price : null,
    intervals: [],
  }
}

function createMarketplaceModel(key: string, platform: string, model: UserSupportedModel): MarketplaceModel {
  const catalog = catalogForModel(model.name)
  return {
    key,
    name: model.name,
    platform: model.platform || platform,
    pricing: model.pricing ?? catalog?.pricing ?? null,
    channels: [],
    groups: [],
  }
}

function catalogForModel(name: string): ModelCatalogEntry | undefined {
  const lower = name.toLowerCase()
  if (MODEL_CATALOG[lower]) return MODEL_CATALOG[lower]
  if (lower.startsWith('gpt-image-2')) return MODEL_CATALOG['gpt-image-2']
  return undefined
}

function billingModeLabel(mode: BillingMode | null): string {
  switch (mode) {
    case BILLING_MODE_TOKEN:
      return '按量计费'
    case BILLING_MODE_PER_REQUEST:
      return '按次计费'
    case BILLING_MODE_IMAGE:
      return '按次计费'
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
    return [{ label: '模型价格', value: `${money(pricing.per_request_price)} / 次` }]
  }

  if (pricing.billing_mode === BILLING_MODE_IMAGE) {
    return [{ label: '模型价格', value: `${money(pricing.image_output_price)} / 次` }]
  }

  const rows = [
    { label: '输入价格', value: `${money(pricing.input_price)} / 1M Tokens` },
    { label: '补全价格', value: `${money(pricing.output_price)} / 1M Tokens` },
  ]
  if (pricing.cache_write_price != null && pricing.cache_write_price > 0) {
    rows.push({ label: '缓存写入价格', value: `${money(pricing.cache_write_price)} / 1M Tokens` })
  }
  rows.push({ label: '缓存读取价格', value: `${money(pricing.cache_read_price)} / 1M Tokens` })
  return rows
}

function modelDescription(model: MarketplaceModel): string {
  const catalog = catalogForModel(model.name)
  if (catalog?.description) return catalog.description
  const lower = model.name.toLowerCase()
  if (lower.includes('image') || lower.includes('banana')) return '生图模型，看文档教程接入。'
  if (lower.includes('haiku')) return '轻量快速模型，适合低成本高频使用。简单问答、总结、改写、基础文案、速度快、成本低。'
  if (lower.includes('mini') || lower.includes('flash')) return '高性价比版本，适合大多数用户。日常聊天、问答、文案、轻量代码便宜、速度快、性价比高。'
  if (lower.includes('codex') || lower.includes('code')) return '编程与代码任务模型，适合代码生成、调试、重构、评审和自动化开发。'
  if (lower.includes('opus') || lower.includes('sonnet') || lower.includes('gpt-5')) return '高性能通用大模型，适合复杂问答、写作、代码、分析、总结和多步骤任务，综合能力强。'
  return `${platformLabel(model.platform)} 模型，适合通用对话、内容生成和工具调用。`
}

function money(value: number | null | undefined): string {
  const n = Number(value ?? 0)
  return `$${n.toFixed(4)}`
}

function platformLabel(platform: string): string {
  switch (platform) {
    case 'openai':
      return 'OpenAI'
    case 'anthropic':
      return 'Anthropic'
    case 'gemini':
      return 'Google'
    case 'antigravity':
      return 'Antigravity'
    default:
      return platform || '未知供应商'
  }
}

function providerIcon(platform: string): string {
  switch (platform) {
    case 'anthropic':
      return '✺'
    case 'gemini':
      return '✦'
    case 'openai':
      return '◎'
    default:
      return 'AI'
  }
}

function countByPlatform(platform: string): number {
  return modelCards.value.filter((model) => model.platform === platform).length
}

function countByGroup(group: string): number {
  return modelCards.value.filter((model) => model.groups.some((g) => g.name === group)).length
}

function countByBilling(mode: BillingMode | 'none'): number {
  return modelCards.value.filter((model) => (model.pricing?.billing_mode ?? 'none') === mode).length
}

function resetFilters() {
  searchQuery.value = ''
  selectedPlatform.value = 'all'
  selectedGroup.value = 'all'
  selectedBillingMode.value = 'all'
}

async function copyModelId(modelId: string) {
  try {
    await copyText(modelId)
    copiedModel.value = modelId
    window.setTimeout(() => {
      if (copiedModel.value === modelId) copiedModel.value = ''
    }, 1600)
  } catch (err) {
    console.error('Failed to copy model id:', err)
  }
}

async function copyVisibleModelIds() {
  const ids = filteredModels.value.map((model) => model.name).join('\n')
  await copyText(ids)
}

async function copyText(text: string) {
  try {
    if (navigator.clipboard?.writeText && window.isSecureContext) {
      await navigator.clipboard.writeText(text)
      return
    }
  } catch {
    // Fall back to textarea copy below for browsers that expose clipboard but reject it on HTTP.
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'readonly')
  textarea.style.position = 'fixed'
  textarea.style.left = '-9999px'
  textarea.style.top = '0'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  const ok = document.execCommand('copy')
  document.body.removeChild(textarea)
  if (!ok) {
    throw new Error('copy failed')
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

function forcePageLabels() {
  document.title = `模型广场 - ${appStore.siteName || 'CN2-API'}`
  nextTick(() => {
    const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
    const nodes: Text[] = []
    while (walker.nextNode()) {
      const node = walker.currentNode as Text
      if (node.nodeValue?.trim() === '可用渠道') nodes.push(node)
    }
    nodes.forEach((node) => {
      node.nodeValue = node.nodeValue?.replace('可用渠道', '模型广场') ?? node.nodeValue
    })
  })
}

onMounted(() => {
  forcePageLabels()
  loadChannels().finally(forcePageLabels)
})
</script>

<style scoped>
.filter-pill {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(229, 231, 235, 0.9);
  background: rgba(255, 255, 255, 0.82);
  padding: 0.625rem 0.875rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: rgb(55, 65, 81);
  transition: all 0.15s ease;
}
.filter-pill:hover,
.filter-pill-active {
  border-color: rgba(20, 184, 166, 0.4);
  background: rgba(20, 184, 166, 0.12);
  color: rgb(13, 148, 136);
}
:global(.dark) .filter-pill {
  border-color: rgba(55, 65, 81, 0.9);
  background: rgba(17, 24, 39, 0.8);
  color: rgb(209, 213, 219);
}
:global(.dark) .filter-pill:hover,
:global(.dark) .filter-pill-active {
  border-color: rgba(45, 212, 191, 0.45);
  background: rgba(20, 184, 166, 0.16);
  color: rgb(94, 234, 212);
}
:global(a[href='/monitor']) {
  display: none !important;
}
</style>

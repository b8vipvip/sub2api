<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-6">
      <section class="overflow-hidden rounded-3xl border border-primary-100 bg-gradient-to-br from-primary-600 via-blue-600 to-indigo-700 p-8 text-white shadow-sm dark:border-primary-900/50">
        <div class="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <p class="text-sm font-semibold uppercase tracking-[0.25em] text-primary-100/80">CN2-API Docs</p>
            <h1 class="mt-3 text-3xl font-bold md:text-4xl">开发文档</h1>
            <p class="mt-3 max-w-3xl text-sm leading-7 text-primary-50/90 md:text-base">
              通过 OpenAI 兼容协议快速接入 CN2-API。先创建 API 密钥，再从模型广场复制模型 ID，即可在 cURL、Python、TypeScript 或常见 AI 工具中调用。
            </p>
          </div>
          <div class="rounded-2xl bg-white/15 p-4 text-sm backdrop-blur">
            <div class="text-primary-50/80">OpenAI Base URL</div>
            <div class="mt-1 font-mono text-lg font-semibold">{{ openAIBaseUrl }}</div>
          </div>
        </div>
      </section>

      <div class="grid gap-6 lg:grid-cols-[260px_1fr]">
        <aside class="hidden space-y-3 lg:block">
          <a
            v-for="item in toc"
            :key="item.href"
            :href="item.href"
            class="block rounded-xl border border-gray-100 bg-white px-4 py-3 text-sm font-medium text-gray-700 transition hover:border-primary-200 hover:bg-primary-50 hover:text-primary-700 dark:border-dark-700 dark:bg-dark-900 dark:text-gray-300 dark:hover:border-primary-800 dark:hover:bg-primary-950/30 dark:hover:text-primary-300"
          >
            {{ item.label }}
          </a>
        </aside>

        <main class="space-y-6">
          <DocSection id="quickstart" title="快速开始" subtitle="3 步完成第一次调用">
            <div class="grid gap-4 md:grid-cols-3">
              <StepCard index="1" title="创建 API 密钥" text="进入左侧 API 密钥，创建并复制 sk- 开头的密钥。" />
              <StepCard index="2" title="选择模型 ID" text="进入模型广场，复制可用模型 ID，例如 gpt-5.4、gpt-5.5。" />
              <StepCard index="3" title="调用接口" text="将 Base URL 改为 CN2-API 地址，使用 OpenAI SDK 或 HTTP 请求。" />
            </div>

            <CodeBlock title="cURL" :code="curlChatSnippet" @copy="copySnippet" />
            <CodeBlock title="Python / OpenAI SDK" :code="pythonSnippet" @copy="copySnippet" />
            <CodeBlock title="TypeScript / OpenAI SDK" :code="typescriptSnippet" @copy="copySnippet" />
          </DocSection>

          <DocSection id="base-url" title="Base URL 与认证" subtitle="同一把 CN2-API Key 可以用于 OpenAI 兼容接口">
            <div class="overflow-hidden rounded-2xl border border-gray-100 dark:border-dark-700">
              <table class="min-w-full divide-y divide-gray-100 text-sm dark:divide-dark-700">
                <thead class="bg-gray-50 text-left text-gray-600 dark:bg-dark-800 dark:text-gray-300">
                  <tr>
                    <th class="px-4 py-3 font-semibold">协议</th>
                    <th class="px-4 py-3 font-semibold">Base URL</th>
                    <th class="px-4 py-3 font-semibold">认证 Header</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 bg-white dark:divide-dark-700 dark:bg-dark-900">
                  <tr>
                    <td class="px-4 py-3 font-medium text-gray-900 dark:text-white">OpenAI 兼容</td>
                    <td class="px-4 py-3 font-mono text-gray-700 dark:text-gray-300">{{ openAIBaseUrl }}</td>
                    <td class="px-4 py-3 font-mono text-gray-700 dark:text-gray-300">Authorization: Bearer sk-xxx</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p class="text-sm leading-7 text-gray-600 dark:text-gray-400">
              页面示例会自动使用当前站点域名。如果你后续给 API 配置独立域名，只需要把 SDK 中的 baseURL 改成新的 API 域名即可。
            </p>
          </DocSection>

          <DocSection id="endpoints" title="常用端点" subtitle="OpenAI 兼容协议入口">
            <div class="grid gap-3 md:grid-cols-2">
              <EndpointCard method="POST" path="/v1/chat/completions" text="对话补全，兼容大多数聊天、代码和工具调用场景。" />
              <EndpointCard method="POST" path="/v1/responses" text="Responses API，适合新版本 OpenAI SDK 或 Agent 工作流。" />
              <EndpointCard method="GET" path="/v1/models" text="获取当前密钥可调用的模型列表。" />
              <EndpointCard method="POST" path="/v1/images/generations" text="图像生成接口，适合 gpt-image 系列等生图模型。" />
            </div>
          </DocSection>

          <DocSection id="models" title="模型 ID 与模型广场" subtitle="模型广场展示当前可调用模型和价格说明">
            <div class="rounded-2xl border border-primary-100 bg-primary-50 p-5 dark:border-primary-900/50 dark:bg-primary-950/20">
              <h3 class="font-semibold text-primary-900 dark:text-primary-200">模型 ID 怎么填？</h3>
              <p class="mt-2 text-sm leading-7 text-primary-800/80 dark:text-primary-100/80">
                在左侧点击“模型广场”，复制模型卡片标题旁边的模型 ID，然后填到请求体的 <code class="rounded bg-white/70 px-1 py-0.5 dark:bg-dark-900/70">model</code> 字段中。
                例如：<code class="rounded bg-white/70 px-1 py-0.5 dark:bg-dark-900/70">gpt-5.5</code>、<code class="rounded bg-white/70 px-1 py-0.5 dark:bg-dark-900/70">gpt-image-2</code>。
              </p>
            </div>
            <CodeBlock title="列出模型" :code="modelsSnippet" @copy="copySnippet" />
          </DocSection>

          <DocSection id="streaming" title="流式响应" subtitle="用于边生成边输出的聊天界面">
            <p class="text-sm leading-7 text-gray-600 dark:text-gray-400">
              Chat Completions 支持 OpenAI 风格的 <code class="rounded bg-gray-100 px-1 py-0.5 dark:bg-dark-800">stream: true</code>。前端聊天框、命令行工具、编辑器插件通常都建议开启流式输出。
            </p>
            <CodeBlock title="Streaming cURL" :code="streamSnippet" @copy="copySnippet" />
          </DocSection>

          <DocSection id="errors" title="错误处理" subtitle="常见返回和排查方向">
            <div class="grid gap-3 md:grid-cols-2">
              <InfoCard title="401 Unauthorized" text="API Key 缺失、填错或已禁用。请检查 Authorization Header。" />
              <InfoCard title="403 Forbidden" text="当前密钥无权限调用该模型，或账户分组未开通对应模型。" />
              <InfoCard title="429 Too Many Requests" text="请求过快或并发过高，可以降低频率后重试。" />
              <InfoCard title="5xx Upstream Error" text="上游模型服务异常或账号不可用，可稍后重试或联系管理员。" />
            </div>
          </DocSection>
        </main>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const siteOrigin = computed(() => {
  if (typeof window === 'undefined') return 'https://sub.cn2.io'
  return window.location.origin.replace(/\/$/, '')
})

const openAIBaseUrl = computed(() => `${siteOrigin.value}/v1`)

const toc = [
  { href: '#quickstart', label: '快速开始' },
  { href: '#base-url', label: 'Base URL 与认证' },
  { href: '#endpoints', label: '常用端点' },
  { href: '#models', label: '模型 ID' },
  { href: '#streaming', label: '流式响应' },
  { href: '#errors', label: '错误处理' },
]

const curlChatSnippet = computed(() => `curl ${openAIBaseUrl.value}/chat/completions \\
  -H "Authorization: Bearer <你的 CN2_API_KEY>" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-5.5",
    "messages": [
      {"role": "user", "content": "你好，介绍一下 CN2-API"}
    ]
  }'`)

const pythonSnippet = computed(() => `from openai import OpenAI

client = OpenAI(
    base_url="${openAIBaseUrl.value}",
    api_key="<你的 CN2_API_KEY>",
)

response = client.chat.completions.create(
    model="gpt-5.5",
    messages=[{"role": "user", "content": "你好，介绍一下 CN2-API"}],
)

print(response.choices[0].message.content)`)

const typescriptSnippet = computed(() => `import OpenAI from 'openai'

const client = new OpenAI({
  baseURL: '${openAIBaseUrl.value}',
  apiKey: '<你的 CN2_API_KEY>',
})

const response = await client.chat.completions.create({
  model: 'gpt-5.5',
  messages: [{ role: 'user', content: '你好，介绍一下 CN2-API' }],
})

console.log(response.choices[0]?.message?.content)`)

const modelsSnippet = computed(() => `curl ${openAIBaseUrl.value}/models \\
  -H "Authorization: Bearer <你的 CN2_API_KEY>"`)

const streamSnippet = computed(() => `curl ${openAIBaseUrl.value}/chat/completions \\
  -H "Authorization: Bearer <你的 CN2_API_KEY>" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-5.5",
    "stream": true,
    "messages": [
      {"role": "user", "content": "写一段 100 字的产品介绍"}
    ]
  }'`)

async function copySnippet(code: string) {
  try {
    if (navigator.clipboard?.writeText && window.isSecureContext) {
      await navigator.clipboard.writeText(code)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = code
      textarea.setAttribute('readonly', 'readonly')
      textarea.style.position = 'fixed'
      textarea.style.left = '-9999px'
      document.body.appendChild(textarea)
      textarea.focus()
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    appStore.showSuccess('复制成功', 3000)
  } catch {
    appStore.showError('复制失败，请手动复制代码')
  }
}
</script>

<script lang="ts">
import { defineComponent, h } from 'vue'

export const DocSection = defineComponent({
  name: 'DocSection',
  props: {
    id: { type: String, required: true },
    title: { type: String, required: true },
    subtitle: { type: String, default: '' },
  },
  setup(props, { slots }) {
    return () => h('section', { id: props.id, class: 'scroll-mt-24 rounded-3xl border border-gray-100 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900' }, [
      h('div', { class: 'mb-5' }, [
        h('h2', { class: 'text-2xl font-bold text-gray-900 dark:text-white' }, props.title),
        props.subtitle ? h('p', { class: 'mt-2 text-sm text-gray-500 dark:text-gray-400' }, props.subtitle) : null,
      ]),
      h('div', { class: 'space-y-4' }, slots.default?.()),
    ])
  },
})

export const StepCard = defineComponent({
  name: 'StepCard',
  props: {
    index: { type: String, required: true },
    title: { type: String, required: true },
    text: { type: String, required: true },
  },
  setup(props) {
    return () => h('div', { class: 'rounded-2xl border border-gray-100 bg-gray-50 p-5 dark:border-dark-700 dark:bg-dark-800' }, [
      h('div', { class: 'flex h-9 w-9 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white' }, props.index),
      h('h3', { class: 'mt-4 font-semibold text-gray-900 dark:text-white' }, props.title),
      h('p', { class: 'mt-2 text-sm leading-6 text-gray-600 dark:text-gray-400' }, props.text),
    ])
  },
})

export const CodeBlock = defineComponent({
  name: 'CodeBlock',
  props: {
    title: { type: String, required: true },
    code: { type: String, required: true },
  },
  emits: ['copy'],
  setup(props, { emit }) {
    return () => h('div', { class: 'overflow-hidden rounded-2xl border border-gray-900/10 bg-gray-950 shadow-sm dark:border-dark-700' }, [
      h('div', { class: 'flex items-center justify-between border-b border-white/10 px-4 py-3' }, [
        h('span', { class: 'text-sm font-semibold text-gray-200' }, props.title),
        h('button', {
          class: 'rounded-lg bg-white/10 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-white/20',
          type: 'button',
          onClick: () => emit('copy', props.code),
        }, '复制'),
      ]),
      h('pre', { class: 'overflow-x-auto p-4 text-sm leading-7 text-gray-100' }, [h('code', props.code)]),
    ])
  },
})

export const EndpointCard = defineComponent({
  name: 'EndpointCard',
  props: {
    method: { type: String, required: true },
    path: { type: String, required: true },
    text: { type: String, required: true },
  },
  setup(props) {
    return () => h('div', { class: 'rounded-2xl border border-gray-100 bg-white p-5 dark:border-dark-700 dark:bg-dark-900' }, [
      h('div', { class: 'flex items-center gap-2' }, [
        h('span', { class: 'rounded-md bg-primary-50 px-2 py-1 text-xs font-bold text-primary-700 dark:bg-primary-950/30 dark:text-primary-300' }, props.method),
        h('code', { class: 'text-sm font-semibold text-gray-900 dark:text-white' }, props.path),
      ]),
      h('p', { class: 'mt-3 text-sm leading-6 text-gray-600 dark:text-gray-400' }, props.text),
    ])
  },
})

export const InfoCard = defineComponent({
  name: 'InfoCard',
  props: {
    title: { type: String, required: true },
    text: { type: String, required: true },
  },
  setup(props) {
    return () => h('div', { class: 'rounded-2xl border border-gray-100 bg-gray-50 p-5 dark:border-dark-700 dark:bg-dark-800' }, [
      h('h3', { class: 'font-semibold text-gray-900 dark:text-white' }, props.title),
      h('p', { class: 'mt-2 text-sm leading-6 text-gray-600 dark:text-gray-400' }, props.text),
    ])
  },
})
</script>

<script setup lang="ts">
import { RouterView, useRouter, useRoute } from 'vue-router'
import { nextTick, onMounted, onBeforeUnmount, watch } from 'vue'
import Toast from '@/components/common/Toast.vue'
import NavigationProgress from '@/components/common/NavigationProgress.vue'
import AdminComplianceDialog from '@/components/admin/AdminComplianceDialog.vue'
import { resolveRouteDocumentTitle } from '@/router/title'
import AnnouncementPopup from '@/components/common/AnnouncementPopup.vue'
import { useAppStore, useAuthStore, useSubscriptionStore, useAnnouncementStore, useAdminComplianceStore, useAdminSettingsStore } from '@/stores'
import { getSetupStatus } from '@/api/setup'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()
const announcementStore = useAnnouncementStore()
const adminComplianceStore = useAdminComplianceStore()
const adminSettingsStore = useAdminSettingsStore()

function ensureDeveloperDocsRoute() {
  if (router.hasRoute('DeveloperDocs')) return
  router.addRoute({
    path: '/developer-docs',
    name: 'DeveloperDocs',
    component: () => import('@/views/user/DeveloperDocsView.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: false,
      title: '开发文档',
    },
  })
}

ensureDeveloperDocsRoute()

function developerDocsIconSvg() {
  return `
    <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5" class="h-5 w-5 flex-shrink-0">
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25A8.966 8.966 0 0118 3.75c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" />
    </svg>
  `
}

function updateDeveloperDocsLinkState() {
  const active = route.path === '/developer-docs'
  document.querySelectorAll<HTMLAnchorElement>('a[data-cn2-developer-docs-nav="true"]').forEach((link) => {
    link.classList.toggle('router-link-active', active)
    link.classList.toggle('router-link-exact-active', active)
    link.classList.toggle('sidebar-link-active', active)
    link.setAttribute('aria-current', active ? 'page' : 'false')
  })
}

function injectDeveloperDocsLinks() {
  if (!authStore.isAuthenticated || appStore.backendModeEnabled) return

  const marketplaceLinks = Array.from(document.querySelectorAll<HTMLAnchorElement>('a[href="/available-channels"]'))
  marketplaceLinks.forEach((marketplaceLink) => {
    const parent = marketplaceLink.parentElement
    if (!parent || parent.querySelector('a[data-cn2-developer-docs-nav="true"]')) return

    const label = marketplaceLink.querySelector<HTMLElement>('.sidebar-label')
    const labelClass = label?.className || 'sidebar-label'
    const ariaHidden = label?.getAttribute('aria-hidden') || 'false'
    const link = document.createElement('a')
    link.href = '/developer-docs'
    link.dataset.cn2DeveloperDocsNav = 'true'
    link.className = marketplaceLink.className
      .replace(/\brouter-link-active\b/g, '')
      .replace(/\brouter-link-exact-active\b/g, '')
      .replace(/\bsidebar-link-active\b/g, '')
      .replace(/\s+/g, ' ')
      .trim()
    link.title = '开发文档'
    link.innerHTML = `${developerDocsIconSvg()}<span class="${labelClass}" aria-hidden="${ariaHidden}">开发文档</span>`
    link.addEventListener('click', (event) => {
      event.preventDefault()
      router.push('/developer-docs')
    })
    marketplaceLink.insertAdjacentElement('afterend', link)
  })

  updateDeveloperDocsLinkState()
}

function normalizeMarketplaceLabels() {
  nextTick(() => {
    ensureDeveloperDocsRoute()
    const marketplaceLinks = document.querySelectorAll<HTMLAnchorElement>('a[href="/available-channels"]')
    marketplaceLinks.forEach((link) => {
      const label = link.querySelector<HTMLElement>('.sidebar-label')
      if (label) {
        label.textContent = '模型广场'
      }
      if (link.title === '可用渠道' || link.title === 'Available Channels') {
        link.title = '模型广场'
      }
    })
    injectDeveloperDocsLinks()
  })
}

function updateDocumentTitle() {
  if (route.path === '/available-channels') {
    document.title = `模型广场 - ${appStore.siteName || 'CN2-API'}`
    normalizeMarketplaceLabels()
    return
  }

  if (route.path === '/developer-docs') {
    document.title = `开发文档 - ${appStore.siteName || 'CN2-API'}`
    normalizeMarketplaceLabels()
    return
  }

  const customMenuItems = [
    ...(appStore.cachedPublicSettings?.custom_menu_items ?? []),
    ...(authStore.isAdmin ? adminSettingsStore.customMenuItems : []),
  ]
  document.title = resolveRouteDocumentTitle(route, appStore.siteName, customMenuItems)
  normalizeMarketplaceLabels()
}

/**
 * Update favicon dynamically
 * @param logoUrl - URL of the logo to use as favicon
 */
function updateFavicon(logoUrl: string) {
  // Find existing favicon link or create new one
  let link = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.type = logoUrl.endsWith('.svg') ? 'image/svg+xml' : 'image/x-icon'
  link.href = logoUrl
}

// Watch for site settings changes and update favicon/title
watch(
  () => appStore.siteLogo,
  (newLogo) => {
    if (newLogo) {
      updateFavicon(newLogo)
    }
  },
  { immediate: true }
)

watch(
  [
    () => route.fullPath,
    () => route.meta.title,
    () => route.meta.titleKey,
    () => appStore.siteName,
    () => appStore.cachedPublicSettings?.custom_menu_items,
    () => authStore.isAdmin,
    () => adminSettingsStore.customMenuItems,
  ],
  updateDocumentTitle,
  { deep: true }
)

// Watch for authentication state and manage subscription data + announcements
function onVisibilityChange() {
  if (document.visibilityState === 'visible' && authStore.isAuthenticated) {
    announcementStore.fetchAnnouncements()
  }
}

function onAdminComplianceRequired(event: Event) {
  const detail = (event as CustomEvent<Record<string, string>>).detail || {}
  adminComplianceStore.requireAcknowledgement(detail)
}

watch(
  () => authStore.isAuthenticated,
  (isAuthenticated, oldValue) => {
    if (isAuthenticated) {
      if (authStore.isAdmin) {
        adminComplianceStore.fetchStatus().catch((error) => {
          console.error('Failed to fetch admin compliance status:', error)
        })
      }

      // User logged in: preload subscriptions and start polling
      subscriptionStore.fetchActiveSubscriptions().catch((error) => {
        console.error('Failed to preload subscriptions:', error)
      })
      subscriptionStore.startPolling()

      // Announcements: new login vs page refresh restore
      if (oldValue === false) {
        // New login: delay 3s then force fetch
        setTimeout(() => announcementStore.fetchAnnouncements(true), 3000)
      } else {
        // Page refresh restore (oldValue was undefined)
        announcementStore.fetchAnnouncements()
      }

      // Register visibility change listener
      document.addEventListener('visibilitychange', onVisibilityChange)
      normalizeMarketplaceLabels()
    } else {
      // User logged out: clear data and stop polling
      subscriptionStore.clear()
      announcementStore.reset()
      adminComplianceStore.reset()
      document.removeEventListener('visibilitychange', onVisibilityChange)
    }
  },
  { immediate: true }
)

// Route change trigger (throttled by store)
router.afterEach(() => {
  if (authStore.isAuthenticated) {
    announcementStore.fetchAnnouncements()
  }
  updateDocumentTitle()
})

onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
  window.removeEventListener('admin-compliance-required', onAdminComplianceRequired)
})

onMounted(async () => {
  window.addEventListener('admin-compliance-required', onAdminComplianceRequired)

  ensureDeveloperDocsRoute()
  if (route.path === '/developer-docs' && route.name !== 'DeveloperDocs') {
    router.replace(route.fullPath)
  }

  // Check if setup is needed
  try {
    const status = await getSetupStatus()
    if (status.needs_setup && route.path !== '/setup') {
      router.replace('/setup')
      return
    }
  } catch {
    // If setup endpoint fails, assume normal mode and continue
  }

  // Load public settings into appStore (will be cached for other components)
  await appStore.fetchPublicSettings()

  // Re-resolve document title now that site settings are available
  updateDocumentTitle()
  normalizeMarketplaceLabels()
})
</script>

<template>
  <NavigationProgress />
  <RouterView />
  <Toast />
  <AnnouncementPopup />
  <AdminComplianceDialog />
</template>

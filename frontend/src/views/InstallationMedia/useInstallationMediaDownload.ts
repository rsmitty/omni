// Copyright (c) 2026 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.
import { ref, type Ref } from 'vue'

import { ManagementService } from '@/api/omni/management/management.pb'
import { useImageFactoryAuth } from '@/methods/useImageFactoryAuth'

// useInstallationMediaDownload returns helpers to download installation media.
//
// When the image factory uses basic-auth credentials (public or basic-auth enterprise),
// `useTokenAuth` is false and the template renders a plain <a :href="link"> directly to
// the factory URL.
//
// When no credentials are available (Auth0 / token-based enterprise), `useTokenAuth` is
// true. On click, `downloadViaToken` calls GetInstallationMediaDownloadURL to get a
// presigned factory URL (HMAC-signed, short-lived), then navigates the browser there.
// No bytes flow through Omni.
export function useInstallationMediaDownload() {
  const auth = useImageFactoryAuth()

  // useTokenAuth returns true when the factory requires token-based auth (no basic-auth creds).
  const useTokenAuth = () => !auth.value?.username || !auth.value?.password

  const downloading: Ref<Record<string, boolean>> = ref({})
  const downloadError: Ref<string | undefined> = ref(undefined)

  async function downloadViaToken(
    schematicId: string,
    talosVersion: string,
    filename: string,
  ): Promise<void> {
    const key = `${schematicId}/${talosVersion}/${filename}`
    if (downloading.value[key]) return

    downloading.value[key] = true
    downloadError.value = undefined

    try {
      const { url } = await ManagementService.GetInstallationMediaDownloadURL({
        schematic_id: schematicId,
        talos_version: talosVersion,
        filename,
      })

      if (!url) throw new Error('no download URL returned')

      window.location.href = url
    } catch (e: unknown) {
      downloadError.value = e instanceof Error ? e.message : String(e)
    } finally {
      downloading.value[key] = false
    }
  }

  function isDownloading(schematicId: string, talosVersion: string, filename: string): boolean {
    return !!downloading.value[`${schematicId}/${talosVersion}/${filename}`]
  }

  return { useTokenAuth, downloadViaToken, isDownloading, downloadError }
}

// Copyright (c) 2025 Sidero Labs, Inc.
//
// Use of this software is governed by the Business Source License
// included in the LICENSE file.

import { createApp } from 'vue';
import "@/index.css";
import App from '@/App.vue'
import VueClipboard from 'vue3-clipboard';
import AppUnavailable from '@/AppUnavailable.vue'
import router from '@/router';


import { initState, ResourceService, Resource } from "@/api/grpc";
import { AuthConfigID, AuthConfigType, DefaultNamespace } from "@/api/resources";
import { Runtime } from "@/api/common/omni.pb";
import { AuthConfigSpec } from "@/api/omni/specs/auth.pb";
import { AuthType, authType, suspended } from "@/methods";
import { createAuth0 } from "@auth0/auth0-vue";
import { withRuntime } from "./api/options";
import vClickOutside from "click-outside-vue3";
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
import yamlWorker from 'monaco-yaml/yaml.worker?worker';

if (process.env.NODE_ENV !== 'development') {
  (self as any).MonacoEnvironment = {
    getWorker(_: string, label: string) {
      if (label === 'yaml') {
        return new yamlWorker();
      }

      return new editorWorker();
    }
  };
}

const setupApp = async () => {
  let authConfigSpec: AuthConfigSpec | undefined = undefined

  try {
    const authConfig: Resource<AuthConfigSpec> = await ResourceService.Get({
      namespace: DefaultNamespace,
      type: AuthConfigType,
      id: AuthConfigID,
    }, withRuntime(Runtime.Omni));

    authConfigSpec = authConfig?.spec as AuthConfigSpec | undefined;
  } catch (e) {
    console.error("failed to get auth parameters", e)
    createApp(AppUnavailable).mount('#app');
    return
  }

  suspended.value = authConfigSpec?.suspended ?? false;

  if (authConfigSpec?.saml?.enabled) {
    authType.value = AuthType.SAML;
  } else if (authConfigSpec?.auth0?.enabled) {
    authType.value = AuthType.Auth0;
  }

  let app = createApp(App)
    .use(router)
    .use(VueClipboard, {
      autoSetContainer: true,
      appendToBody: true,
    })

    if (authType.value === AuthType.Auth0) {
      app = app.use(createAuth0({
        domain: authConfigSpec!.auth0!.domain!,
        clientId: authConfigSpec!.auth0!.client_id!,
        authorizationParams: {
          redirect_uri: window.location.origin,
        },
        useFormData: !!authConfigSpec!.auth0?.useFormData,
      }))
    }

  app.use(vClickOutside);
  app.mount('#app');
}

initState();

setupApp();

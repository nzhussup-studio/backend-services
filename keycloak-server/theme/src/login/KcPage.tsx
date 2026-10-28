import { lazy, Suspense } from "react";
import DefaultPage from "keycloakify/login/DefaultPage";
import type { ClassKey } from "keycloakify/login";
import type { KcContext } from "./KcContext";
import Template from "./Template";
import { useI18n } from "./i18n";

const UserProfileFormFields = lazy(() => import("keycloakify/login/UserProfileFormFields"));

export default function KcPage({ kcContext }: { kcContext: KcContext }) {
  const { i18n } = useI18n({ kcContext });

  return (
    <Suspense>
      <DefaultPage
        kcContext={kcContext}
        i18n={i18n}
        classes={classes}
        Template={Template}
        doUseDefaultCss={false}
        UserProfileFormFields={UserProfileFormFields}
        doMakeUserConfirmPassword={true}
      />
    </Suspense>
  );
}

const classes = {} satisfies { [key in ClassKey]?: string };

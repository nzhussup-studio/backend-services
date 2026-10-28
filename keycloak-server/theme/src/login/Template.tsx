import { useEffect } from "react";
import type { TemplateProps } from "keycloakify/login/TemplateProps";
import { kcSanitize } from "keycloakify/lib/kcSanitize";
import { useSetClassName } from "keycloakify/tools/useSetClassName";
import logoDark from "../assets/nz-dark.svg";
import logoLight from "../assets/nz-light.svg";
import { useAdaptiveFavicon } from "../useAdaptiveFavicon";
import type { KcContext } from "./KcContext";
import type { I18n } from "./i18n";
import "./theme.css";

export default function Template(props: TemplateProps<KcContext, I18n>) {
  const {
    displayInfo = false,
    displayMessage = true,
    displayRequiredFields = false,
    headerNode,
    socialProvidersNode,
    infoNode,
    documentTitle,
    bodyClassName,
    kcContext,
    i18n,
    children
  } = props;

  const { msg, msgStr } = i18n;
  const { realm, auth, url, message, isAppInitiatedAction } = kcContext;

  useEffect(() => {
    document.title = documentTitle ?? msgStr("loginTitle", realm.displayName || realm.name);
  }, [documentTitle, msgStr, realm.displayName, realm.name]);

  useSetClassName({ qualifiedName: "html", className: "nz-auth-html" });
  useSetClassName({ qualifiedName: "body", className: bodyClassName ?? "nz-auth-body" });
  useAdaptiveFavicon();

  const shouldShowMessage =
    displayMessage &&
    message !== undefined &&
    (message.type !== "warning" || !isAppInitiatedAction);

  return (
    <main className="auth-shell">
      <section className="auth-panel" aria-label="Account access">
        <div className="auth-panel-inner">
          <a className="auth-brand" href="https://nzhussup.dev" aria-label="Nurzhanat Zhussup home">
            <picture>
              <source srcSet={logoDark} media="(prefers-color-scheme: dark)" />
              <img src={logoLight} alt="NZ." />
            </picture>
          </a>

          <header className="auth-form-header">
            {auth?.showUsername && !auth.showResetCredentials ? (
              <div className="attempted-username">
                <span>{auth.attemptedUsername}</span>
                <a href={url.loginRestartFlowUrl}>{msg("restartLoginTooltip")}</a>
              </div>
            ) : (
              <h1 id="kc-page-title">{headerNode}</h1>
            )}

            {displayRequiredFields && <p className="required-note">* {msg("requiredFields")}</p>}
          </header>

          {shouldShowMessage && (
            <div className={`auth-alert auth-alert-${message.type}`} role="alert">
              <span dangerouslySetInnerHTML={{ __html: kcSanitize(message.summary) }} />
            </div>
          )}

          <div id="kc-content" className="auth-content">
            {children}

            {auth?.showTryAnotherWayLink && (
              <form id="kc-select-try-another-way-form" action={url.loginAction} method="post">
                <input type="hidden" name="tryAnotherWay" value="on" />
                <button className="auth-text-button" type="submit">{msg("doTryAnotherWay")}</button>
              </form>
            )}

            {socialProvidersNode}

            {displayInfo && infoNode && <div id="kc-info" className="auth-info">{infoNode}</div>}
          </div>
        </div>
      </section>
    </main>
  );
}

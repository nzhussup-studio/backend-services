import { createGetKcContextMock } from "keycloakify/login/KcContext";
import type { DeepPartial } from "keycloakify/tools/DeepPartial";
import { kcEnvDefaults, themeNames } from "../kc.gen";
import type { KcContext, KcContextExtension, KcContextExtensionPerPage } from "./KcContext";
import KcPage from "./KcPage";

const kcContextExtension: KcContextExtension = {
  themeName: themeNames[0],
  properties: { ...kcEnvDefaults }
};

export const { getKcContextMock } = createGetKcContextMock({
  kcContextExtension,
  kcContextExtensionPerPage: {} satisfies KcContextExtensionPerPage,
  overrides: {},
  overridesPerPage: {}
});

export function createKcPageStory<PageId extends KcContext["pageId"]>({ pageId }: { pageId: PageId }) {
  function KcPageStory({ kcContext: overrides }: { kcContext?: DeepPartial<Extract<KcContext, { pageId: PageId }>> }) {
    return <KcPage kcContext={getKcContextMock({ pageId, overrides })} />;
  }

  return { KcPageStory };
}

import { i18nBuilder } from "keycloakify/login";
import type { ThemeName } from "../kc.gen";

// The builder exposes this value specifically so its inferred shape can become our I18n type.
// eslint-disable-next-line @typescript-eslint/no-unused-vars
const { useI18n, ofTypeI18n } = i18nBuilder.withThemeName<ThemeName>().build();

export { useI18n };
export type I18n = typeof ofTypeI18n;

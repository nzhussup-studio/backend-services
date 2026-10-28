import { useEffect } from "react";
import logoDark from "./assets/nz-dark.svg";
import logoLight from "./assets/nz-light.svg";

export function useAdaptiveFavicon() {
  useEffect(() => {
    const colorScheme = window.matchMedia("(prefers-color-scheme: dark)");
    let favicon = document.querySelector<HTMLLinkElement>("link[rel~='icon']");

    if (!favicon) {
      favicon = document.createElement("link");
      favicon.rel = "icon";
      favicon.type = "image/svg+xml";
      document.head.appendChild(favicon);
    }

    const faviconElement = favicon;
    const update = () => {
      faviconElement.href = colorScheme.matches ? logoDark : logoLight;
    };

    update();
    colorScheme.addEventListener("change", update);

    return () => {
      colorScheme.removeEventListener("change", update);
    };
  }, []);
}
